package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"inventory/internal/model"
	"inventory/internal/repository/part"

	"github.com/ZanDattSu/star-factory/shared/pkg/interceptor"
	inventoryv1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/inventory/v1"
)

const (
	httpPort = 8081
	grpcPort = 50051

	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second

	apiRelativePath = "../shared/api"
)

// InventoryService реализует gRPC-сервис управления деталями.
type InventoryService struct {
	inventoryv1.UnimplementedInventoryServiceServer
	Storage *part.repository
}

func NewInventoryService(storage *part.repository) *InventoryService {
	return &InventoryService{Storage: storage}
}

func (is *InventoryService) GetPart(_ context.Context, req *inventoryv1.GetPartRequest) (*model.Part, error) {
	part, ok := is.Storage.GetPart(req.Uuid)
	if !ok {
		return nil, status.Errorf(codes.NotFound, "part %s not found", req.Uuid)
	}

	return part, nil
}

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.LoggerInterceptor(),
			interceptor.ValidationInterceptor(),
		),
	)

	partStorage := part.NewRepository()
	partStorage.InitTestData()

	service := NewInventoryService(partStorage)

	inventoryv1.RegisterInventoryServiceServer(server, service)

	reflection.Register(server)

	go func() {
		log.Printf("🚀 gRPC server listening on %d\n", grpcPort)
		err := server.Serve(listener)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	// Запускаем HTTP сервер с gRPC Gateway
	var gatewayServer *http.Server
	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		mux := runtime.NewServeMux()

		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

		err := inventoryv1.RegisterInventoryServiceHandlerFromEndpoint(
			ctx,
			mux,
			fmt.Sprintf("localhost:%d", grpcPort),
			opts,
		)
		if err != nil {
			log.Printf("Failed to register gateway: %v\n", err)
			return
		}

		// Создаем файловый сервер для Swagger UI
		fileServer := http.FileServer(http.Dir(apiRelativePath))

		httpMux := http.NewServeMux()

		httpMux.Handle("/api/", mux)

		// Swagger UI эндпоинты
		httpMux.Handle("/swagger-ui.html", fileServer)
		httpMux.Handle("/inventory/v1/inventory.swagger.json", fileServer)

		// Редирект для swagger.json
		httpMux.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, apiRelativePath+"/inventory/v1/inventory.swagger.json")
		})

		// Редирект с корня на Swagger UI
		httpMux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path == "/" {
				http.Redirect(w, req, "/swagger-ui.html", http.StatusMovedPermanently)
				return
			}
			fileServer.ServeHTTP(w, req)
		})

		gatewayServer = &http.Server{
			Addr:              fmt.Sprintf(":%d", httpPort),
			Handler:           httpMux,
			ReadHeaderTimeout: readHeaderTimeout,
		}

		log.Printf("🌐 HTTP server with gRPC-Gateway listening on %d\n", httpPort)
		err = gatewayServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Failed to serve HTTP: %v\n", err)
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
			log.Printf("HTTP server shutdown error: %v", err)
		}

		log.Println("✅ HTTP server stopped")
	}

	log.Println("🛑 Shutting down gRPC server...")
	if closeErr := listener.Close(); closeErr != nil {
		log.Printf("failed to close listener: %v\n", closeErr)
	}

	server.GracefulStop()
	log.Println("✅ Server stopped")
}
