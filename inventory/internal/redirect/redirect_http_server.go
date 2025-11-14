package redirect

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/ZanDattSu/star-factory/platform/pkg/path"
	inventoryV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/inventory/v1"
)

const (
	readHeaderTimeout = 5 * time.Second
)

var apiPath = path.GetPathRelativeToRoot("/shared/api")

type HTTPServer struct {
	server *http.Server
}

func NewHTTPServer(ctx context.Context, grpcAddress, httpAddress string) (*HTTPServer, error) {
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := inventoryV1.RegisterInventoryServiceHandlerFromEndpoint(
		ctx,
		mux,
		grpcAddress,
		opts,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register gateway: %w", err)
	}

	httpMux := registerSwaggerMux(mux)

	gatewayServer := &http.Server{
		Addr:              httpAddress,
		Handler:           httpMux,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	return &HTTPServer{
		server: gatewayServer,
	}, nil
}

func (s *HTTPServer) Serve() error {
	return s.server.ListenAndServe()
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func registerSwaggerMux(mux *runtime.ServeMux) *http.ServeMux {
	// Создаем файловый сервер для Swagger UI
	fileServer := http.FileServer(http.Dir(apiPath))

	httpMux := http.NewServeMux()

	httpMux.Handle("/api/", mux)

	// Swagger UI эндпоинты
	httpMux.Handle("/swagger-ui.html", fileServer)
	httpMux.Handle("/inventory/v1/inventory.swagger.json", fileServer)

	// Редирект для swagger.json
	httpMux.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, apiPath+"/inventory/v1/inventory.swagger.json")
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
