package main

import (
	"context"
	"fmt"
	orderv1 "github.com/ZanDattSu/star-factory/shared/pkg/openapi/order/v1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-faster/errors"
	"github.com/google/uuid"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"shared/pkg/openapi/order/v1"
	"sync"
	"syscall"
	"time"
)

const (
	httpPort = "8080"
	// Таймауты для HTTP-сервера
	responseTimeout   = 10 * time.Second
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

type Order struct {
	orderv1.GetOrderResponse
}

type OrderStorage struct {
	orders map[uuid.UUID]*Order
	mu     sync.RWMutex
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[uuid.UUID]*Order),
	}
}

func (s *OrderStorage) GetOrder(uuid uuid.UUID) (*Order, bool) {
	s.mu.Lock()
	defer s.mu.RUnlock()
	order, ok := s.orders[uuid]
	return order, ok
}

func (s *OrderStorage) PutOrder(uuid uuid.UUID, order *Order) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.orders[uuid] = order
}

type OrderHandler struct {
	storage *OrderStorage
}

func NewOrderHandler(storage *OrderStorage) *OrderHandler {
	return &OrderHandler{
		storage: storage,
	}
}

func (o *OrderHandler) CreateOrder(_ context.Context, req *order_v1.CreateOrderRequest) (order_v1.CreateOrderRes, error) {
	//TODO implement me
	panic("implement me")
}

func paymentServicePayOrderFake(userUUID uuid.UUID, orderUUID uuid.UUID, method orderv1.PaymentMethod) uuid.UUID {
	return uuid.New()
}

func (o *OrderHandler) OrderPay(_ context.Context, req *order_v1.PayOrderRequest, params order_v1.OrderPayParams) (order_v1.OrderPayRes, error) {
	orderUUID := params.OrderUUID
	order, ok := o.storage.GetOrder(orderUUID)
	if !ok {
		return orderNotFoundError(orderUUID)
	}

	order.TransactionUUID.SetTo(
		//TODO убрать fake
		paymentServicePayOrderFake(
			order.UserUUID,
			order.OrderUUID,
			order.PaymentMethod.Value,
		),
	)

	order.SetStatus(orderv1.OrderStatusPAID)
	order.PaymentMethod.Value = req.PaymentMethod

	return &orderv1.PayOrderResponse{
		TransactionUUID: order.GetTransactionUUID().Value,
	}, nil
}

func (o *OrderHandler) GetOrder(_ context.Context, params order_v1.GetOrderParams) (order_v1.GetOrderRes, error) {
	orderUUID := params.OrderUUID
	order, ok := o.storage.GetOrder(orderUUID)
	if !ok {
		return orderNotFoundError(orderUUID)
	}
	return order, nil
}

func (o *OrderHandler) CancelOrder(_ context.Context, params order_v1.CancelOrderParams) (order_v1.CancelOrderRes, error) {
	orderUUID := params.OrderUUID

	order, ok := o.storage.GetOrder(orderUUID)
	if !ok {
		return orderNotFoundError(orderUUID)
	}

	var resp order_v1.CancelOrderRes

	switch order.GetStatus() {
	case orderv1.OrderStatusPENDINGPAYMENT:
		order.SetStatus(orderv1.OrderStatusCANCELLED)
		resp = &orderv1.CancelOrderNoContent{}
	case orderv1.OrderStatusPAID:
		resp = &orderv1.ConflictError{
			Code:    409,
			Message: "Cannot cancel a paid order",
		}
	case orderv1.OrderStatusCANCELLED:
		resp = &orderv1.ConflictError{
			Code:    409,
			Message: "Cannot cancel a canceled order",
		}
	}

	return resp, nil
}

func (o *OrderHandler) NewError(_ context.Context, err error) *order_v1.GenericErrorStatusCode {
	return &orderv1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: order_v1.GenericError{
			Code:    orderv1.NewOptInt(http.StatusInternalServerError),
			Message: orderv1.NewOptString(err.Error()),
		},
	}
}

func orderNotFoundError(orderUUID uuid.UUID) (*orderv1.NotFoundError, error) {
	return &orderv1.NotFoundError{
		Code:    404,
		Message: fmt.Sprintf("Order %d not found", orderUUID),
	}, nil
}

func main() {
	orderStorage := NewOrderStorage()

	log.Println("Создаем обработчик API погоды")
	orderHandler := NewOrderHandler(orderStorage)

	log.Println("Создаем OpenAPI сервер")
	orderServer, err := orderv1.NewServer(orderHandler)
	if err != nil {
		log.Fatalf("Ошибка создания сервера OpenAPI: %v", err)
	}

	r := chi.NewRouter()

	// Добавляем middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(responseTimeout))

	// Монтируем обработчики OpenAPI
	r.Mount("/", orderServer)

	server := http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout, // Защита от Slowloris атак - тип DDoS-атаки, при которой
		// атакующий умышленно медленно отправляет HTTP-заголовки, удерживая соединения открытыми и истощая
		// пул доступных соединений на сервере. ReadHeaderTimeout принудительно закрывает соединение,
		// если клиент не успел отправить все заголовки за отведенное время.
	}

	go func() {
		log.Printf("HTTP-сервер запущен на порту %s\n", httpPort)
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Ошибка запуска сервера: %v\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	// SIGTERM - "вежливая" просьба завершиться,
	// SIGINT - прерывание с клавиатуры (Ctrl+C)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("Завершение работы сервера...")

	// Создаем контекст с таймаутом для остановки сервера
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}
