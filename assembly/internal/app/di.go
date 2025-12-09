package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"

	"github.com/ZanDattSu/star-factory/assembly/internal/config"
	kafkaConverter "github.com/ZanDattSu/star-factory/assembly/internal/converter/kafka"
	"github.com/ZanDattSu/star-factory/assembly/internal/converter/kafka/decoder"
	"github.com/ZanDattSu/star-factory/assembly/internal/service"
	assemblyService "github.com/ZanDattSu/star-factory/assembly/internal/service/assembly"
	orderPaidConsumer "github.com/ZanDattSu/star-factory/assembly/internal/service/consumer/order_paid_consumer"
	shipAssembledProducer "github.com/ZanDattSu/star-factory/assembly/internal/service/producer/ship_assembled_producer"
	"github.com/ZanDattSu/star-factory/platform/pkg/closer"
	wrappedKafka "github.com/ZanDattSu/star-factory/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/ZanDattSu/star-factory/platform/pkg/kafka/consumer"
	wrappedKafkaProducer "github.com/ZanDattSu/star-factory/platform/pkg/kafka/producer"
	"github.com/ZanDattSu/star-factory/platform/pkg/logger"
	kafkaMiddleware "github.com/ZanDattSu/star-factory/platform/pkg/middleware/kafka"
)

type diContainer struct {
	// Services
	assemblyService              service.AssemblyService
	orderPaidConsumerService     service.OrderPaidConsumerService
	shipAssembledProducerService service.ShipAssembledProducerService

	// Converters
	orderPaidDecoder kafkaConverter.OrderPaidDecoder

	// Kafka infrastructure
	consumerGroup         sarama.ConsumerGroup
	orderPaidConsumer     wrappedKafka.Consumer
	shipAssembledProducer wrappedKafka.Producer
	syncProducer          sarama.SyncProducer
}

func NewDIContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) AssemblyService() service.AssemblyService {
	if d.assemblyService == nil {
		d.assemblyService = assemblyService.NewService(d.ShipAssembledProducerService())
	}
	return d.assemblyService
}

func (d *diContainer) OrderPaidConsumerService() service.OrderPaidConsumerService {
	if d.orderPaidConsumerService == nil {
		d.orderPaidConsumerService = orderPaidConsumer.NewService(
			d.OrderPaidConsumer(),
			d.OrderPaidDecoder(),
			d.AssemblyService(),
		)
	}
	return d.orderPaidConsumerService
}

func (d *diContainer) ShipAssembledProducerService() service.ShipAssembledProducerService {
	if d.shipAssembledProducerService == nil {
		d.shipAssembledProducerService = shipAssembledProducer.NewService(d.ShipAssembledProducer())
	}
	return d.shipAssembledProducerService
}

func (d *diContainer) ConsumerGroup() sarama.ConsumerGroup {
	if d.consumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderConsumer.GroupID(),
			config.AppConfig().OrderConsumer.Config(),
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

func (d *diContainer) OrderPaidConsumer() wrappedKafka.Consumer {
	if d.orderPaidConsumer == nil {
		d.orderPaidConsumer = wrappedKafkaConsumer.NewConsumer(
			d.ConsumerGroup(),
			[]string{
				config.AppConfig().OrderConsumer.Topic(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}
	return d.orderPaidConsumer
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

func (d *diContainer) ShipAssembledProducer() wrappedKafka.Producer {
	if d.shipAssembledProducer == nil {
		d.shipAssembledProducer = wrappedKafkaProducer.NewProducer(
			d.SyncProducer(),
			config.AppConfig().OrderProducer.Topic(),
			logger.Logger(),
		)
	}
	return d.shipAssembledProducer
}

func (d *diContainer) OrderPaidDecoder() kafkaConverter.OrderPaidDecoder {
	if d.orderPaidDecoder == nil {
		d.orderPaidDecoder = decoder.NewOrderPaidDecoder()
	}

	return d.orderPaidDecoder
}
