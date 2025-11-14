package app

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	orderApi "order/internal/api/v1"
	ordApi "order/internal/api/v1/order"
	gRPCClient "order/internal/client/grpc"
	inventoryService "order/internal/client/grpc/inventory/v1"
	paymentService "order/internal/client/grpc/payment/v1"
	"order/internal/config"
	orderRepo "order/internal/repository"
	"order/internal/repository/order/postgresql"
	orderService "order/internal/service"
	ordService "order/internal/service/order"
	"platform/pkg/closer"

	inventoryV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	orderApi        orderApi.OrderApi
	orderService    orderService.OrderService
	orderRepository orderRepo.OrderRepository

	paymentClient   gRPCClient.PaymentClient
	inventoryClient gRPCClient.InventoryClient

	postgreSQLPool *pgxpool.Pool
}

func NewDIContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) OrderApi(ctx context.Context) orderApi.OrderApi {
	if d.orderApi == nil {
		d.orderApi = ordApi.NewApi(d.OrderService(ctx))
	}

	return d.orderApi
}

func (d *diContainer) OrderService(ctx context.Context) orderService.OrderService {
	if d.orderService == nil {
		d.orderService = ordService.NewService(
			d.OrderRepository(ctx),
			d.PaymentClient(ctx),
			d.InventoryClient(ctx),
		)
	}

	return d.orderService
}

func (d *diContainer) PaymentClient(_ context.Context) gRPCClient.PaymentClient {
	if d.paymentClient == nil {
		paymentConn, err := newGRPCConnectWithoutSecure(config.AppConfig().Payment.PaymentAddress())
		if err != nil {
			panic(fmt.Sprintf(
				"Failed to connect to payment gRPC service (%s): %v",
				config.AppConfig().Payment.PaymentServicePort(),
				err,
			))
		}

		closer.AddNamed("Payment connection", func(ctx context.Context) error {
			if closeErr := paymentConn.Close(); closeErr != nil {
				return closeErr
			}
			return nil
		})

		paymentClient := paymentV1.NewPaymentServiceClient(paymentConn)

		d.paymentClient = paymentService.NewClient(paymentClient)
	}

	return d.paymentClient
}

func (d *diContainer) InventoryClient(_ context.Context) gRPCClient.InventoryClient {
	if d.inventoryClient == nil {
		inventoryConn, err := newGRPCConnectWithoutSecure(config.AppConfig().Inventory.InventoryAddress())
		if err != nil {
			panic(fmt.Sprintf(
				"Failed to connect to inventory gRPC service (%s): %v",
				config.AppConfig().Inventory.InventoryServicePort(),
				err,
			))
		}
		closer.AddNamed("Inventory connection", func(ctx context.Context) error {
			if closeErr := inventoryConn.Close(); closeErr != nil {
				return closeErr
			}
			return nil
		})

		inventoryClient := inventoryV1.NewInventoryServiceClient(inventoryConn)

		d.inventoryClient = inventoryService.NewClient(inventoryClient)
	}

	return d.inventoryClient
}

func (d *diContainer) OrderRepository(ctx context.Context) orderRepo.OrderRepository {
	if d.orderRepository == nil {
		d.orderRepository = postgresql.NewRepository(d.PostgreSQLPool(ctx))
	}

	return d.orderRepository
}

func (d *diContainer) PostgreSQLPool(ctx context.Context) *pgxpool.Pool {
	if d.postgreSQLPool == nil {
		dbURI := config.AppConfig().Postgres.URI()

		pool, err := pgxpool.New(ctx, dbURI)
		if err != nil {
			panic(fmt.Sprintf("Failed to create pgxpool connect: %s", err))
		}

		err = pool.Ping(ctx)
		if err != nil {
			panic(fmt.Sprintf("Database is unavailable: %s", err))
		}

		closer.AddNamed("PostgreSQL pool", func(ctx context.Context) error {
			pool.Close()
			return nil
		})

		d.postgreSQLPool = pool
	}

	return d.postgreSQLPool
}

func newGRPCConnectWithoutSecure(port string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(
		port,
		grpc.WithTransportCredentials(insecure.NewCredentials()), // отключаем TLS
	)
	return conn, err
}
