package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"shared/pkg/interceptor"
	"sync"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	inventoryv1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/inventory/v1"
)

const (
	httpPort = 8081
	grpcPort = 50051

	responseTimeout   = 10 * time.Second
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

type Part = inventoryv1.Part

type PartStorage struct {
	parts map[string]*inventoryv1.Part
	mu    sync.RWMutex
}

func NewPartStorage() *PartStorage {
	return &PartStorage{
		parts: make(map[string]*inventoryv1.Part),
	}
}

// GetPart возвращает деталь по UUID. Потокобезопасно.
func (ps *PartStorage) GetPart(uuid string) (*Part, bool) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	part, ok := ps.parts[uuid]
	return part, ok
}

// PutPart сохраняет деталь по UUID. Потокобезопасно.
func (ps *PartStorage) PutPart(uuid string, part *Part) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.parts[uuid] = part
}

// DeletePart удаляет деталь по UUID. Потокобезопасно.
func (ps *PartStorage) DeletePart(uuid string) bool {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	_, existed := ps.parts[uuid]
	delete(ps.parts, uuid)
	return existed
}

// Values возвращает все детали. Потокобезопасно.
func (ps *PartStorage) Values() []*Part {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	parts := make([]*Part, 0, len(ps.parts))
	for _, part := range ps.parts {
		parts = append(parts, part)
	}
	return parts
}

// InventoryService реализует gRPC-сервис управления деталями.
type InventoryService struct {
	inventoryv1.UnimplementedInventoryServiceServer
	Storage *PartStorage
}

func NewInventoryService(storage *PartStorage) *InventoryService {
	return &InventoryService{Storage: storage}
}

func (is *InventoryService) GetPart(_ context.Context, req *inventoryv1.GetPartRequest) (*inventoryv1.GetPartResponse, error) {
	part, ok := is.Storage.GetPart(req.Uuid)
	if !ok {
		return nil, status.Errorf(codes.NotFound, "part %s not found", req.Uuid)
	}

	return &inventoryv1.GetPartResponse{
		Part: part,
	}, nil
}

// filterIsEmpty проверяет, пустой ли фильтр.
func filterIsEmpty(f *inventoryv1.PartsFilter) bool {
	if f == nil {
		return true
	}
	return len(f.Uuids) == 0 &&
		len(f.Names) == 0 &&
		len(f.Categories) == 0 &&
		len(f.ManufacturerCountries) == 0 &&
		len(f.Tags) == 0
}

// toSet преобразует slice в set для O(1) поиска.
func toSet[T comparable](values []T) map[T]struct{} {
	set := make(map[T]struct{}, len(values))
	for _, v := range values {
		set[v] = struct{}{}
	}
	return set
}

// filterByField возвращает детали, у которых значение поля есть в values.
//
// В отличие от реализации через slices.Contains (O(n²)),
// использует внутренний set на основе map для поиска за O(1),
// что обеспечивает общую сложность O(n + m).
//
// n — количество деталей, m — количество элементов фильтра.
func filterByField[T comparable](
	parts []*Part,
	values []T,
	getField func(*Part) T,
) []*Part {
	if len(values) == 0 {
		return parts
	}

	set := toSet(values)
	result := make([]*Part, 0, len(parts))

	for _, p := range parts {
		if _, exists := set[getField(p)]; exists {
			result = append(result, p)
		}
	}

	return result
}

// filterByTags оставляет детали с хотя бы одним тегом из списка (OR-логика).
func filterByTags(parts []*Part, tags []string) []*Part {
	if len(tags) == 0 {
		return parts
	}

	tagSet := toSet(tags)
	result := make([]*Part, 0, len(parts))

	for _, p := range parts {
		if hasAnyTag(p.Tags, tagSet) {
			result = append(result, p)
		}
	}

	return result
}

// hasAnyTag проверяет наличие хотя бы одного тега из set'а.
func hasAnyTag(partTags []string, tagSet map[string]struct{}) bool {
	for _, tag := range partTags {
		if _, exists := tagSet[tag]; exists {
			return true
		}
	}
	return false
}

// FilterFunc представляет одну стадию фильтрации в pipeline.
type FilterFunc func([]*Part) []*Part

// buildFilterPipeline создаёт цепочку фильтров на основе PartsFilter.
func buildFilterPipeline(f *inventoryv1.PartsFilter) []FilterFunc {
	var pipeline []FilterFunc

	if len(f.Uuids) > 0 {
		pipeline = append(pipeline, func(parts []*Part) []*Part {
			return filterByField(
				parts,
				f.Uuids,
				func(p *Part) string { return p.Uuid },
			)
		})
	}

	if len(f.Names) > 0 {
		pipeline = append(pipeline, func(parts []*Part) []*Part {
			return filterByField(
				parts,
				f.Names,
				func(part *Part) string { return part.Name },
			)
		})
	}

	if len(f.Categories) > 0 {
		pipeline = append(pipeline, func(parts []*Part) []*Part {
			return filterByField(
				parts,
				f.Categories,
				func(p *Part) inventoryv1.Category { return p.Category },
			)
		})
	}

	if len(f.ManufacturerCountries) > 0 {
		pipeline = append(pipeline, func(parts []*Part) []*Part {
			return filterByField(
				parts,
				f.ManufacturerCountries,
				func(p *Part) string { return p.Manufacturer.Country },
			)
		})
	}

	if len(f.Tags) > 0 {
		pipeline = append(pipeline, func(parts []*Part) []*Part {
			return filterByTags(parts, f.Tags)
		})
	}

	return pipeline
}

// applyFilterPipeline последовательно применяет все фильтры.
func applyFilterPipeline(parts []*Part, pipeline []FilterFunc) []*Part {
	for _, filter := range pipeline {
		parts = filter(parts)
		// прерываем если результат пуст
		if len(parts) == 0 {
			return parts
		}
	}
	return parts
}

// ListParts возвращает отфильтрованный список деталей.
// Использует pipeline для последовательной фильтрации с AND-логикой между полями
// и OR-логикой внутри каждого поля.
func (is *InventoryService) ListParts(
	_ context.Context,
	req *inventoryv1.ListPartsRequest,
) (*inventoryv1.ListPartsResponse, error) {
	parts := is.Storage.Values()

	if filterIsEmpty(req.Filter) {
		return &inventoryv1.ListPartsResponse{Parts: parts}, nil
	}

	// Строим и применяем pipeline фильтров
	pipeline := buildFilterPipeline(req.Filter)
	filteredParts := applyFilterPipeline(parts, pipeline)

	return &inventoryv1.ListPartsResponse{Parts: filteredParts}, nil
}

func seedParts(ps *PartStorage) {
	parts := []*inventoryv1.Part{
		{
			Uuid:     "11111111-1111-1111-1111-111111111111",
			Name:     "Fusion Engine Mk.I",
			Price:    4999.90,
			Category: inventoryv1.Category_CATEGORY_ENGINE,
			Manufacturer: &inventoryv1.Manufacturer{
				Name:    "StarWorks",
				Country: "JP",
			},
			Tags: []string{"core", "engine"},
		},
		{
			Uuid:     "22222222-2222-2222-2222-222222222222",
			Name:     "Quantum Hull Plate",
			Price:    799.00,
			Category: inventoryv1.Category_CATEGORY_WING,
			Manufacturer: &inventoryv1.Manufacturer{
				Name:    "OrbitalFoundry",
				Country: "DE",
			},
			Tags: []string{"hull", "shielding"},
		},
		{
			Uuid:     "33333333-3333-3333-3333-333333333333",
			Name:     "Cryo Fuel Pump",
			Price:    249.50,
			Category: inventoryv1.Category_CATEGORY_FUEL,
			Manufacturer: &inventoryv1.Manufacturer{
				Name:    "DeepSpace Ltd",
				Country: "US",
			},
			Tags: []string{"fuel", "pump"},
		},
	}
	for _, p := range parts {
		ps.PutPart(p.Uuid, p)
	}
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

	partStorage := NewPartStorage()
	seedParts(partStorage)

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

		gatewayServer = &http.Server{
			Addr:              fmt.Sprintf(":%d", httpPort),
			Handler:           mux,
			ReadHeaderTimeout: readHeaderTimeout,
		}

		log.Printf("🌐 HTTP server with gRPC-Gateway listening on %d\n", httpPort)
		err = gatewayServer.ListenAndServe()
		if err != nil && !errors.Is(http.ErrServerClosed, err) {
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
