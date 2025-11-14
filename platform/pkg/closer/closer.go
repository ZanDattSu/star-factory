package closer

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"platform/pkg/logger"
	"sync"
	"time"
)

// Глобальный экземпляр для использования по всему приложению
var (
	globalCloser = NewWithLogger(&logger.NoopLogger{})
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

// Closer управляет процессом graceful shutdown приложения
type Closer struct {
	mu              sync.Mutex                    // Защита от гонки при добавлении функций
	once            sync.Once                     // Гарантия однократного вызова CloseAll
	done            chan struct{}                 // Канал для оповещения о завершении
	funcs           []func(context.Context) error // Зарегистрированные функции закрытия
	logger          Logger                        // Используемый логгер
	shutdownTimeout time.Duration
}

// SetShutdownTimeout устанавливает gracefulShutdownTimeout
func SetShutdownTimeout(shutdownTimeout time.Duration) {
	globalCloser.shutdownTimeout = shutdownTimeout
}

// New создаёт новый экземпляр Closer с дефолтным логгером log.Default()
func New(shutdownTimeout time.Duration, signals ...os.Signal) *Closer {
	return NewWithLogger(logger.Logger(), signals...)
}

// NewWithLogger создаёт новый экземпляр Closer с указанием логгера.
// Если переданы сигналы, Closer начнёт их слушать и вызовет CloseAll при получении.
func NewWithLogger(logger Logger, signals ...os.Signal) *Closer {
	c := &Closer{
		done:   make(chan struct{}),
		logger: logger,
	}

	if len(signals) > 0 {
		go c.handleSignals(signals...)
	}

	return c
}

// SetLogger позволяет установить кастомный логгер для глобального closer'а
func SetLogger(l Logger) {
	globalCloser.SetLogger(l)
}

// Add добавляет функции закрытия в глобальный closer
func Add(f ...func(context.Context) error) {
	globalCloser.Add(f...)
}

// AddNamed добавляет функцию закрытия с именем зависимости для логирования в глобальный closer
func AddNamed(name string, f func(context.Context) error) {
	globalCloser.AddNamed(name, f)
}

// Configure настраивает глобальный closer для обработки системных сигналов
func Configure(signals ...os.Signal) {
	go globalCloser.handleSignals(signals...)
}

// CloseAll инициирует процесс закрытия всех зарегистрированных функций глобального closer'а
func CloseAll(ctx context.Context) error {
	return globalCloser.CloseAll(ctx)
}

// SetLogger устанавливает логгер для Closer
func (c *Closer) SetLogger(l Logger) {
	c.logger = l
}

// SetShutdownTimeout устанавливает gracefulShutdownTimeout
func (c *Closer) SetShutdownTimeout(shutdownTimeout time.Duration) {
	c.shutdownTimeout = shutdownTimeout
}

// Add добавляет одну или несколько функций закрытия
func (c *Closer) Add(f ...func(context.Context) error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.funcs = append(c.funcs, f...)
}

// AddNamed добавляет функцию закрытия с именем зависимости для логирования
func (c *Closer) AddNamed(name string, f func(context.Context) error) {
	c.Add(func(ctx context.Context) error {
		start := time.Now()

		c.logger.Info(ctx, fmt.Sprintf("Закрываем %s...", name))

		err := f(ctx)

		duration := time.Since(start)

		if err != nil {
			c.logger.Error(ctx, fmt.Sprintf("Ошибка при закрытии %s: %v (заняло %s)", name, err, duration))
		} else {
			c.logger.Info(ctx, fmt.Sprintf("%s успешно закрыт за %s", name, duration))
		}
		return err
	})
}

// handleSignals обрабатывает системные сигналы и вызывает CloseAll с shutdown context
func (c *Closer) handleSignals(signals ...os.Signal) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, signals...)
	defer signal.Stop(ch)

	select {
	case <-ch:
		c.logger.Info(context.Background(), "Получен системный сигнал, начинаем graceful shutdown...")

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), c.shutdownTimeout)
		defer shutdownCancel()

		if err := c.CloseAll(shutdownCtx); err != nil {
			c.logger.Error(context.Background(), "Ошибка при закрытии ресурсов", zap.Error(err))
		}
	case <-c.done:
		return
		// CloseAll уже был вызван вручную, просто выходим
	}
}

// CloseAll вызывает все зарегистрированные функции закрытия.
// Возвращает первую возникшую ошибку, если таковая была.
func (c *Closer) CloseAll(ctx context.Context) error {
	var result error

	c.once.Do(func() {
		defer close(c.done)

		c.mu.Lock()
		funcs := c.funcs
		c.funcs = nil // освободим память
		c.mu.Unlock()

		if len(funcs) == 0 {
			c.logger.Info(ctx, "Нет функций для закрытия.")
			return
		}

		c.logger.Info(ctx, "Начинаем процесс graceful shutdown...")

		errCh := make(chan error, len(funcs))
		var wg sync.WaitGroup

		// Выполняем в обратном порядке, потому что последний созданный ресурс часто зависит от предыдущих.
		// При закрытии нужно сначала освободить самые "верхние" зависимые ресурсы, а потом базовые.
		for i := len(funcs) - 1; i >= 0; i-- {

			wg.Add(1)
			go func(f func(context.Context) error) {
				defer wg.Done()

				defer func() {
					if r := recover(); r != nil {
						errCh <- errors.New("panic recovered in closer")
						c.logger.Error(ctx, "Panic в функции закрытия", zap.Any("error", r))
					}
				}()

				if err := f(ctx); err != nil {
					errCh <- err
				}
			}(funcs[i])
		}

		// Закрываем канал ошибок, когда все функции завершатся
		go func() {
			wg.Wait()
			close(errCh)
		}()

		for {
			select {
			case <-ctx.Done():
				c.logger.Info(ctx, "Контекст отменён во время закрытия", zap.Error(ctx.Err()))
				if result == nil {
					result = ctx.Err()
				}
				return
			case err, ok := <-errCh:
				// Канал закрыт - все функции завершились
				if !ok {
					c.logger.Info(ctx, "Все ресурсы успешно закрыты")
					return
				}
				c.logger.Error(ctx, "Ошибка при закрытии", zap.Error(err))
				if result == nil {
					result = err
				}
			}
		}
	})

	return result
}
