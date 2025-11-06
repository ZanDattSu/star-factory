package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	partV1Api "inventory/internal/api/v1/part"
	partRepo "inventory/internal/repository/part"
	"inventory/internal/servers"
	partService "inventory/internal/service/part"
)

const (
	grpcPort = 50051
	httpPort = 8081

	shutdownTimeout = 10 * time.Second
)

func main() {
	repo := partRepo.NewRepository()
	repo.InitTestData()

	service := partService.NewService(repo)
	api := partV1Api.NewApi(service)

	gRPCServer, err := servers.NewGRPCServer(grpcPort, api)
	if err != nil {
		log.Println(err)
		return
	}

	gatewayServer, err := servers.NewHTTPServer(httpPort)
	if err != nil {
		log.Println(err)
		return
	}

	// Запускаем gRPC сервер
	go func() {
		if err := gRPCServer.Serve(); err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	// Запускаем HTTP сервер с gRPC Gateway
	go func() {
		if err := gatewayServer.Serve(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("failed to serve HTTP: %s\n", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Сначала останавливаем HTTP сервер
	if gatewayServer != nil {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err := gatewayServer.Shutdown(shutdownCtx); err != nil {
			log.Println(err)
			return
		}
	}

	if err = gRPCServer.Shutdown(); err != nil {
		log.Println(err)
		return
	}
}
