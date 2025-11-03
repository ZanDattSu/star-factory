package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	inventoryv1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/inventory/v1"
)

const grpcPort = 50052

type PartStorage struct {
	parts map[string]*inventoryv1.Part
	mu    sync.RWMutex
}

func NewPartStorage() *PartStorage {
	return &PartStorage{
		parts: make(map[string]*inventoryv1.Part),
	}
}

func (s *PartStorage) GetPart(uuid string) (*inventoryv1.Part, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	part, ok := s.parts[uuid]
	return part, ok
}

func (s *PartStorage) PutPart(uuid string, part *inventoryv1.Part) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.parts[uuid] = part
}

type InventoryService struct {
	inventoryv1.UnimplementedInventoryServiceServer
	storage *PartStorage
}

func NewInventoryService(storage *PartStorage) *InventoryService {
	return &InventoryService{storage: storage}
}

func (is *InventoryService) GetPart(_ context.Context, req *inventoryv1.GetPartRequest) (*inventoryv1.GetPartResponse, error) {
	part, ok := is.storage.GetPart(req.Uuid)
	if !ok {
		return nil, status.Errorf(codes.NotFound, "part %s not found", req.Uuid)
	}

	return &inventoryv1.GetPartResponse{
		Part: part,
	}, nil
}

func (is *InventoryService) ListParts(ctx context.Context, req *inventoryv1.ListPartsRequest) (*inventoryv1.ListPartsResponse, error) {
	var parts []*inventoryv1.Part

	for _, part := range is.storage.parts {
		parts = append(parts, part)
	}

	return &inventoryv1.ListPartsResponse{
		Parts: parts,
	}, nil
}

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}

	server := grpc.NewServer()

	partStorage := NewPartStorage()

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

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down gRPC server...")
	if closeErr := listener.Close(); closeErr != nil {
		log.Printf("failed to close listener: %v\n", closeErr)
	}

	server.GracefulStop()
	log.Println("✅ Server stopped")
}
