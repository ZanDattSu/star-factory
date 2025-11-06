package servers

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	inventoryV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/inventory/v1"
)

const (
	apiRelativePath   = "../shared/api"
	readHeaderTimeout = 5 * time.Second
)

type HTTPServer struct {
	server *http.Server
	port   int
}

func (s *HTTPServer) GetPort() int {
	return s.port
}

func NewHTTPServer(ctx context.Context, httpPort, grpcPort int) (*HTTPServer, error) {
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := inventoryV1.RegisterInventoryServiceHandlerFromEndpoint(
		ctx,
		mux,
		setEndpoint(grpcPort),
		opts,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register gateway: %w", err)
	}

	httpMux := registerSwaggerMux(mux)

	gatewayServer := &http.Server{
		Addr:              setEndpoint(httpPort),
		Handler:           httpMux,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	return &HTTPServer{
		server: gatewayServer,
		port:   httpPort,
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

func registerSwaggerMux(mux *runtime.ServeMux) *http.ServeMux {
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
	return httpMux
}

func setEndpoint(port int) string {
	return net.JoinHostPort("localhost", strconv.Itoa(port))
}
