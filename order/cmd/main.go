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

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-faster/errors"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	orderV1 "github.com/ZanDattSu/star-factory/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/payment/v1"
)

const (
	orderServicePort     = "8080"
	paymentServicePort   = "50052"
	inventoryServicePort = "50051"

	// Таймауты для HTTP-сервера
	responseTimeout   = 10 * time.Second
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
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
	log.Println("Создаем payment gRPC клиент")
	conn, err := newGRPCConnectWithoutSecure(paymentServicePort)
	if err != nil {
		log.Printf("❌ Ошибка подключения к gRPC (%s): %v", inventoryServicePort, err)
		return
	}
	defer func() {
		if closeErr := conn.Close(); closeErr != nil {
			log.Printf("failed to close connect: %v", closeErr)
		}
	}()

	paymentClient := paymentV1.NewPaymentServiceClient(conn)
	log.Printf("✅ Успешно создан payment gRPC-клиент (%s)", paymentServicePort)

	log.Println("======================================")

	log.Println("Создаем inventory gRPC клиент")
	conn, err = newGRPCConnectWithoutSecure(inventoryServicePort)
	if err != nil {
		log.Printf("❌ Ошибка подключения к gRPC (%s): %v", inventoryServicePort, err)
		return
	}
	defer func() {
		if closeErr := conn.Close(); closeErr != nil {
			log.Printf("failed to close connect: %v", closeErr)
		}
	}()

	inventoryClient := inventoryV1.NewInventoryServiceClient(conn)
	log.Printf("✅ Успешно создан inventory gRPC-клиент (%s)", inventoryServicePort)

	log.Println("======================================")

	orderStorage := NewOrderStorage()

	log.Println("Создаем обработчик Order API")
	orderHandler := NewOrderHandler(orderStorage, paymentClient, inventoryClient)

	log.Println("Создаем OpenAPI сервер")
	orderServer, err := orderV1.NewServer(orderHandler)
	if err != nil {
		log.Printf("Ошибка создания сервера OpenAPI: %v", err)
		return
	}

	r := chi.NewRouter()

	// Добавляем middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(responseTimeout))

	// Монтируем обработчики OpenAPI
	r.Mount("/", orderServer)

	server := http.Server{
		Addr:              getAddress(orderServicePort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout, // Защита от Slowloris атак - тип DDoS-атаки, при которой
		// атакующий умышленно медленно отправляет HTTP-заголовки, удерживая соединения открытыми и истощая
		// пул доступных соединений на сервере. ReadHeaderTimeout принудительно закрывает соединение,
		// если клиент не успел отправить все заголовки за отведенное время.
	}

	go func() {
		log.Printf("HTTP-сервер запущен на порту %s\n", orderServicePort)
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Ошибка запуска сервера: %v\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	// SIGTERM - "вежливая" просьба завершиться,
	// SIGINT - прерывание с клавиатуры (Ctrl+C)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("Завершение работы сервера...")

	// Создаем контекст с таймаутом для остановки сервера
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}
