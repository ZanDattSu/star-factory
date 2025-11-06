package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-faster/errors"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	httpServer "order/internal/server"

	orderV1 "github.com/ZanDattSu/star-factory/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/payment/v1"
)

const (
	orderServicePort     = "8080"
	paymentServicePort   = "50052"
	inventoryServicePort = "50051"

	shutdownTimeout = 10 * time.Second
)

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[string]*orderV1.OrderDto),
	}
}

type OrderStorage struct {
	orders map[string]*orderV1.OrderDto
	mu     sync.RWMutex
}

func (s *OrderStorage) GetOrder(uuid string) (*orderV1.OrderDto, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	order, ok := s.orders[uuid]
	return order, ok
}

func (s *OrderStorage) PutOrder(uuid string, order *orderV1.OrderDto) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.orders[uuid] = order
}

type OrderHandler struct {
	storage             *OrderStorage
	paymentGRPCClient   paymentV1.PaymentServiceClient
	inventoryGRPCClient inventoryV1.InventoryServiceClient
}

func NewOrderHandler(
	storage *OrderStorage,
	paymentClient paymentV1.PaymentServiceClient,
	inventoryClient inventoryV1.InventoryServiceClient,
) *OrderHandler {
	return &OrderHandler{
		storage:             storage,
		paymentGRPCClient:   paymentClient,
		inventoryGRPCClient: inventoryClient,
	}
}

func (oh *OrderHandler) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	parts, err := oh.inventoryGRPCClient.ListParts(
		ctx,
		&inventoryV1.ListPartsRequest{
			Filter: &inventoryV1.PartsFilter{
				Uuids: req.PartUuids,
				// TODO параметры запросов
				Names:                 nil,
				Categories:            nil,
				ManufacturerCountries: nil,
				Tags:                  nil,
			},
		},
	)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			return &orderV1.NotFoundError{
				Code:    404,
				Message: "one or more parts not found",
			}, nil
		} else {
			return &orderV1.InternalServerError{
				Code:    500,
				Message: fmt.Sprintf("failed to list parts from inventory service: %v", err),
			}, nil
		}
	}

	if len(parts.Parts) != len(req.PartUuids) {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "one or more parts not found",
		}, nil
	}

	partUuids := make([]string, 0, len(parts.Parts))
	var totalPrice float64
	for _, part := range parts.Parts {
		partUuids = append(partUuids, part.Uuid)
		totalPrice += part.Price
	}

	orderUUID := uuid.New().String()
	newOrder := &orderV1.OrderDto{
		OrderUUID:  orderUUID,
		UserUUID:   req.UserUUID,
		PartUuids:  partUuids,
		TotalPrice: totalPrice,
		Status:     orderV1.OrderStatusPENDINGPAYMENT,
	}
	oh.storage.PutOrder(orderUUID, newOrder)

	return &orderV1.CreateOrderResponse{
		OrderUUID:  orderUUID,
		TotalPrice: totalPrice,
	}, nil
}

func convertPaymentMethod(method orderV1.PaymentMethod) paymentV1.PaymentMethod {
	switch method {
	case orderV1.PaymentMethodCARD:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CARD
	case orderV1.PaymentMethodSBP:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_SBP
	case orderV1.PaymentMethodCREDITCARD:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case orderV1.PaymentMethodINVESTORMONEY:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}

func (oh *OrderHandler) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	orderUUID := params.OrderUUID

	if orderUUID == "" {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: "empty order UUID",
		}, nil
	}

	order, ok := oh.storage.GetOrder(orderUUID)
	if !ok {
		return OrderNotFoundError(orderUUID), nil
	}

	paymentResponse, err := oh.paymentGRPCClient.PayOrder(
		ctx,
		&paymentV1.PayOrderRequest{
			OrderUuid:     order.OrderUUID,
			UserUuid:      order.UserUUID,
			PaymentMethod: convertPaymentMethod(req.PaymentMethod),
		})
	if err != nil {
		statusCode, ok := status.FromError(err)
		if ok && statusCode.Code() == codes.Internal {
			return &orderV1.InternalServerError{
				Code:    500,
				Message: fmt.Sprintf("payment service internal error: %v", err),
			}, nil
		}
	}

	order.SetStatus(orderV1.OrderStatusPAID)
	order.TransactionUUID.SetTo(paymentResponse.TransactionUuid)
	order.PaymentMethod = orderV1.NewOptPaymentMethod(req.PaymentMethod)

	oh.storage.PutOrder(order.OrderUUID, order)

	return &orderV1.PayOrderResponse{
		TransactionUUID: order.GetTransactionUUID().Value,
	}, nil
}

func (oh *OrderHandler) GetOrder(_ context.Context, params orderV1.GetOrderParams) (orderV1.GetOrderRes, error) {
	orderUUID := params.OrderUUID
	order, ok := oh.storage.GetOrder(orderUUID)
	if !ok {
		return OrderNotFoundError(orderUUID), nil
	}
	return order, nil
}

func (oh *OrderHandler) CancelOrder(_ context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	orderUUID := params.OrderUUID

	order, ok := oh.storage.GetOrder(orderUUID)
	if !ok {
		return OrderNotFoundError(orderUUID), nil
	}

	var resp orderV1.CancelOrderRes

	switch order.GetStatus() {
	case orderV1.OrderStatusPENDINGPAYMENT:
		order.SetStatus(orderV1.OrderStatusCANCELLED)
		oh.storage.PutOrder(order.OrderUUID, order)
		resp = &orderV1.CancelOrderNoContent{}
	case orderV1.OrderStatusPAID:
		resp = &orderV1.ConflictError{
			Code:    409,
			Message: "Cannot cancel a paid order",
		}
	case orderV1.OrderStatusCANCELLED:
		resp = &orderV1.ConflictError{
			Code:    409,
			Message: "Cannot cancel a canceled order",
		}
	}

	return resp, nil
}

func newGRPCConnectWithoutSecure(port string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(
		getAddress(port),
		grpc.WithTransportCredentials(insecure.NewCredentials()), // отключаем TLS
	)
	return conn, err
}

func getAddress(port string) string {
	return net.JoinHostPort("localhost", port)
}

func (oh *OrderHandler) NewError(_ context.Context, err error) *orderV1.GenericErrorStatusCode {
	return &orderV1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: orderV1.GenericError{
			Code:    orderV1.NewOptInt(http.StatusInternalServerError),
			Message: orderV1.NewOptString(err.Error()),
		},
	}
}

func OrderNotFoundError(orderUUID string) *orderV1.NotFoundError {
	return &orderV1.NotFoundError{
		Code:    404,
		Message: fmt.Sprintf("Order %s not found", orderUUID),
	}
}

func main() {
	log.Println("Creating payment gRPC client...")
	paymentConn, err := newGRPCConnectWithoutSecure(paymentServicePort)
	if err != nil {
		log.Printf("❌ Failed to connect to payment gRPC service (%s): %v", inventoryServicePort, err)
		return
	}
	defer func() {
		if closeErr := paymentConn.Close(); closeErr != nil {
			log.Printf("Failed to close payment gRPC connection: %v", closeErr)
		}
	}()

	paymentClient := paymentV1.NewPaymentServiceClient(paymentConn)
	log.Printf("✅ Payment gRPC client created successfully (%s)", paymentServicePort)

	log.Println("======================================")

	log.Println("Creating inventory gRPC client...")
	inventoryConn, err := newGRPCConnectWithoutSecure(inventoryServicePort)
	if err != nil {
		log.Printf("❌ Failed to connect to inventory gRPC service (%s): %v", inventoryServicePort, err)
		return
	}
	defer func() {
		if closeErr := inventoryConn.Close(); closeErr != nil {
			log.Printf("Failed to close inventory gRPC connection: %v", closeErr)
		}
	}()

	inventoryClient := inventoryV1.NewInventoryServiceClient(inventoryConn)
	log.Printf("✅ Inventory gRPC client created successfully (%s)", inventoryServicePort)

	log.Println("======================================")

	orderStorage := NewOrderStorage()

	log.Println("Creating order API handler...")
	orderHandler := NewOrderHandler(orderStorage, paymentClient, inventoryClient)

	log.Println("Creating HTTP server...")
	server, err := httpServer.NewHTTPServer(orderServicePort, orderHandler)
	if err != nil {
		log.Printf("❌ Failed to create HTTP server: %v", err)
		return
	}
	log.Println("✅ HTTP server created successfully")

	go func() {
		log.Printf("🌐 HTTP server listening on :%s\n", orderServicePort)
		if err := server.Serve(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("HTTP server error: %v", err)
			return
		}
	}()

	gracefulShutdown()

	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown error: %v", err)
		return
	}

	log.Println("✅ HTTP server stopped successfully")
}

// gracefulShutdown мягко завершает работу программы
// когда в канал quit поступает один из сисколлов ОС
//
// SIGTERM - "вежливая" просьба завершиться,
// SIGINT - прерывание с клавиатуры (Ctrl+C)
func gracefulShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
}
