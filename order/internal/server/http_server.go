package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	orderV1 "github.com/ZanDattSu/star-factory/shared/pkg/openapi/order/v1"
)

const (
	readHeaderTimeout = 5 * time.Second
	responseTimeout   = 10 * time.Second
)

type HTTPServer struct {
	server *http.Server
}

func NewHTTPServer(address string, api orderV1.Handler) (*HTTPServer, error) {
	// Создаем OpenAPI сервер
	openAPIHandler, err := orderV1.NewServer(api)
	if err != nil {
		return nil, fmt.Errorf("failed to create OpenAPI server: %w", err)
	}

	// Настраиваем роутер
	r := chi.NewRouter()

	// Добавляем middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(responseTimeout))

	// Монтируем обработчики OpenAPI
	r.Mount("/", openAPIHandler)

	// Создаем HTTP сервер
	server := &http.Server{
		Addr:              address,
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout, // Защита от Slowloris атак
		// тип DDoS-атаки, при которой атакующий умышленно медленно отправляет HTTP-заголовки,
		// удерживая соединения открытыми и истощая пул доступных соединений на сервере.
		// ReadHeaderTimeout принудительно закрывает соединение,
		// если клиент не успел отправить все заголовки за отведенное время.
	}

	return &HTTPServer{
		server: server,
	}, nil
}

func (s *HTTPServer) Serve() error {
	return s.server.ListenAndServe()
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("HTTP server shutdown error: %w", err)
	}
	return nil
}
