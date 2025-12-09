package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderApi "github.com/ZanDattSu/star-factory/order/internal/api/v1"
	ordApi "github.com/ZanDattSu/star-factory/order/internal/api/v1/order"
	gRPCClient "github.com/ZanDattSu/star-factory/order/internal/client/grpc"
	inventoryService "github.com/ZanDattSu/star-factory/order/internal/client/grpc/inventory/v1"
	paymentService "github.com/ZanDattSu/star-factory/order/internal/client/grpc/payment/v1"
	"github.com/ZanDattSu/star-factory/order/internal/config"
	kafkaDecoder "github.com/ZanDattSu/star-factory/order/internal/converter/kafka"
	"github.com/ZanDattSu/star-factory/order/internal/converter/kafka/decoder"
	orderRepo "github.com/ZanDattSu/star-factory/order/internal/repository"
	"github.com/ZanDattSu/star-factory/order/internal/repository/order/postgresql"
	orderService "github.com/ZanDattSu/star-factory/order/internal/service"
	"github.com/ZanDattSu/star-factory/order/internal/service/consumer/order_consumer"
	ordService "github.com/ZanDattSu/star-factory/order/internal/service/order"
	"github.com/ZanDattSu/star-factory/order/internal/service/produser/order_producer"
	"github.com/ZanDattSu/star-factory/platform/pkg/closer"
	wrappedKafka "github.com/ZanDattSu/star-factory/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/ZanDattSu/star-factory/platform/pkg/kafka/consumer"
	wrappedKafkaProducer "github.com/ZanDattSu/star-factory/platform/pkg/kafka/producer"
	"github.com/ZanDattSu/star-factory/platform/pkg/logger"
	kafkaMiddleware "github.com/ZanDattSu/star-factory/platform/pkg/middleware/kafka"
	inventoryV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	// API
	orderApi orderApi.OrderApi

	// Services
	orderService            orderService.OrderService
	assemblyConsumerService orderService.ConsumerService
	orderProducerService    orderService.OrderProducerService

	// Repository
	orderRepository orderRepo.OrderRepository

	// gRPC Clients
	paymentClient   gRPCClient.PaymentClient
	inventoryClient gRPCClient.InventoryClient

	// PostgreSQL
	postgreSQLPool *pgxpool.Pool

	// Kafka Decoder
	assemblyDecoder kafkaDecoder.ShipAssembledDecoder

	// Kafka Infrastructure
	consumerGroup    sarama.ConsumerGroup
	assemblyConsumer wrappedKafka.Consumer
	orderProducer    wrappedKafka.Producer
	syncProducer     sarama.SyncProducer
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
			d.OrderProducerService(),
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

func (d *diContainer) AssemblyConsumerService(ctx context.Context) orderService.ConsumerService {
	if d.assemblyConsumerService == nil {
		d.assemblyConsumerService = order_consumer.NewService(
			d.AssemblyConsumer(),
			d.AssemblyDecoder(),
			d.OrderService(ctx),
			d.OrderRepository(ctx),
		)
	}
	return d.assemblyConsumerService
}

func (d *diContainer) AssemblyConsumer() wrappedKafka.Consumer {
	if d.assemblyConsumer == nil {
		d.assemblyConsumer = wrappedKafkaConsumer.NewConsumer(
			d.ConsumerGroup(),
			[]string{
				config.AppConfig().AssemblyConsumer.Topic(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}
	return d.assemblyConsumer
}

func (d *diContainer) ConsumerGroup() sarama.ConsumerGroup {
	if d.consumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().AssemblyConsumer.GroupID(),
			config.AppConfig().AssemblyConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create consumer group: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka consumer group", func(ctx context.Context) error {
			return d.consumerGroup.Close()
		})

		d.consumerGroup = consumerGroup
	}
	return d.consumerGroup
}

func (d *diContainer) AssemblyDecoder() kafkaDecoder.ShipAssembledDecoder {
	if d.assemblyDecoder == nil {
		d.assemblyDecoder = decoder.NewAssemblyDecoder()
	}
	return d.assemblyDecoder
}

func (d *diContainer) OrderProducerService() orderService.OrderProducerService {
	if d.orderProducerService == nil {
		d.orderProducerService = order_producer.NewService(d.OrderProducer())
	}
	return d.orderProducerService
}

func (d *diContainer) OrderProducer() wrappedKafka.Producer {
	if d.orderProducer == nil {
		d.orderProducer = wrappedKafkaProducer.NewProducer(
			d.SyncProducer(),
			config.AppConfig().OrderProducer.Topic(),
			logger.Logger(),
		)
	}
	return d.orderProducer
}

func (d *diContainer) SyncProducer() sarama.SyncProducer {
	if d.syncProducer == nil {
		p, err := sarama.NewSyncProducer(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderProducer.Config(),
		)
		if err != nil {
			panic("failed to create sync producer: " + err.Error())
		}

		closer.AddNamed("Kafka sync producer", func(ctx context.Context) error {
			return p.Close()
		})

		d.syncProducer = p
	}
	return d.syncProducer
}
