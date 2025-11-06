package servers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	paymentV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/payment/v1"
)

const (
	apiRelativePath   = "../shared/api"
	readHeaderTimeout = 5 * time.Second
)

type HTTPServer struct {
	server *http.Server
	port   int
}

func NewHTTPServer(httpPort int) (*HTTPServer, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := paymentV1.RegisterPaymentServiceHandlerFromEndpoint(
		ctx,
		mux,
		fmt.Sprintf("localhost:%d", httpPort),
		opts,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register gateway: %w", err)
	}

	httpMux := registerSwaggerMux(mux)

	gatewayServer := &http.Server{
		Addr:              fmt.Sprintf(":%d", httpPort),
		Handler:           httpMux,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	return &HTTPServer{
		server: gatewayServer,
		port:   httpPort,
	}, nil
}

func (s *HTTPServer) Serve() error {
	log.Printf("🌐 HTTP server with gRPC-Gateway listening on %d\n", s.port)
	return s.server.ListenAndServe()
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("HTTP server shutdown error: %w", err)
	}
	log.Println("✅ HTTP server stopped")
	return nil
}

func registerSwaggerMux(mux *runtime.ServeMux) *http.ServeMux {
	// Создаем файловый сервер для Swagger UI
	fileServer := http.FileServer(http.Dir(apiRelativePath))

	httpMux := http.NewServeMux()

	httpMux.Handle("/api/", mux)

	// Swagger UI эндпоинты
	httpMux.Handle("/swagger-ui.html", fileServer)
	httpMux.Handle("/payment/v1/payment.swagger.json", fileServer)

	// Редирект для swagger.json
	httpMux.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, apiRelativePath+"/payment/v1/payment.swagger.json")
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
