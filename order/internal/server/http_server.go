package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	httpmiddleware "github.com/ZanDattSu/star-factory/platform/pkg/middleware/http"
	orderV1 "github.com/ZanDattSu/star-factory/shared/pkg/openapi/order/v1"
	authV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/auth/v1"
)

const (
	readHeaderTimeout = 5 * time.Second
	responseTimeout   = 10 * time.Second
)

type HTTPServer struct {
	server *http.Server
}

func NewHTTPServer(address string, api orderV1.Handler, authClient authV1.AuthServiceClient) (*HTTPServer, error) {
	openAPIHandler, err := orderV1.NewServer(api)
	if err != nil {
		return nil, fmt.Errorf("failed to create OpenAPI server: %w", err)
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(responseTimeout))

	r.Group(func(r chi.Router) {
		auth := httpmiddleware.NewAuthMiddleware(authClient)
		r.Use(auth.Handle)
		r.Mount("/", openAPIHandler)
	})

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
