package app

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/IBM/sarama"
	"github.com/go-telegram/bot"

	httpClient "github.com/ZanDattSu/star-factory/notification/internal/client/http"
	telegramClient "github.com/ZanDattSu/star-factory/notification/internal/client/http/telegram"
	"github.com/ZanDattSu/star-factory/notification/internal/config"
	kafkaConverter "github.com/ZanDattSu/star-factory/notification/internal/converter/kafka"
	"github.com/ZanDattSu/star-factory/notification/internal/converter/kafka/decoder"
	"github.com/ZanDattSu/star-factory/notification/internal/service"
	orderPaidConsumer "github.com/ZanDattSu/star-factory/notification/internal/service/consumer/order_paid_consumer"
	shipAssembledConsumer "github.com/ZanDattSu/star-factory/notification/internal/service/consumer/ship_assembled_consumer"
	"github.com/ZanDattSu/star-factory/notification/internal/service/telegram"
	"github.com/ZanDattSu/star-factory/platform/pkg/closer"
	wrappedKafka "github.com/ZanDattSu/star-factory/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/ZanDattSu/star-factory/platform/pkg/kafka/consumer"
	"github.com/ZanDattSu/star-factory/platform/pkg/logger"
	kafkaMiddleware "github.com/ZanDattSu/star-factory/platform/pkg/middleware/kafka"
)

type diContainer struct {
	// Services
	notificationService          service.NotificationService
	orderPaidConsumerService     service.OrderPaidConsumerService
	shipAssembledConsumerService service.ShipAssembledConsumerService

	// Converters
	orderPaidDecoder     kafkaConverter.OrderPaidDecoder
	shipAssembledDecoder kafkaConverter.ShipAssembledDecoder

	// telegram
	telegramClient httpClient.TelegramClient
	telegramBot    *bot.Bot

	// Consumer Groups
	shipAssembledConsumerGroup sarama.ConsumerGroup
	orderPaidConsumerGroup     sarama.ConsumerGroup

	// Consumers
	shipAssembledConsumer wrappedKafka.Consumer
	orderPaidConsumer     wrappedKafka.Consumer
}

func NewDIContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) NotificationService() service.NotificationService {
	if d.notificationService == nil {
		d.notificationService = telegram.NewService(d.TelegramClient())
	}
	return d.notificationService
}

func (d *diContainer) OrderPaidConsumerService() service.OrderPaidConsumerService {
	if d.orderPaidConsumerService == nil {
		d.orderPaidConsumerService = orderPaidConsumer.NewService(
			d.OrderPaidConsumer(),
			d.OrderPaidDecoder(),
			d.NotificationService(),
		)
	}
	return d.orderPaidConsumerService
}

func (d *diContainer) ShipAssembledConsumerService() service.ShipAssembledConsumerService {
	if d.shipAssembledConsumerService == nil {
		d.shipAssembledConsumerService = shipAssembledConsumer.NewService(
			d.ShipAssembledConsumer(),
			d.ShipAssembledDecoder(),
			d.NotificationService(),
		)
	}
	return d.shipAssembledConsumerService
}

func (d *diContainer) ShipAssembledDecoder() kafkaConverter.ShipAssembledDecoder {
	if d.shipAssembledDecoder == nil {
		d.shipAssembledDecoder = decoder.NewAssemblyDecoder()
	}
	return d.shipAssembledDecoder
}

func (d *diContainer) OrderPaidDecoder() kafkaConverter.OrderPaidDecoder {
	if d.orderPaidDecoder == nil {
		d.orderPaidDecoder = decoder.NewOrderPaidDecoder()
	}

	return d.orderPaidDecoder
}

func (d *diContainer) TelegramClient() httpClient.TelegramClient {
	if d.telegramClient == nil {
		tgBot, _ := d.TelegramBot() //nolint:gosec
		d.telegramClient = telegramClient.NewClient(tgBot)
	}

	return d.telegramClient
}

func (d *diContainer) TelegramBot() (*bot.Bot, error) {
	if d.telegramBot == nil {
		dialer := &net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}

		transport := &http.Transport{
			DialContext:           dialer.DialContext,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 0,
		}

		client := &http.Client{
			Transport: transport,
		}

		b, err := bot.New(
			config.AppConfig().TelegramBot.Token(),
			bot.WithHTTPClient(0, client),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create telegram bot: %w", err)
		}

		d.telegramBot = b
	}

	return d.telegramBot, nil
}

func (d *diContainer) ShipAssembledConsumerGroup() sarama.ConsumerGroup {
	if d.shipAssembledConsumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().ShipAssembledConsumer.GroupID(),
			config.AppConfig().ShipAssembledConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create ship assembled consumer group: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka ship assembled consumer group", func(ctx context.Context) error {
			return d.shipAssembledConsumerGroup.Close()
		})

		d.shipAssembledConsumerGroup = consumerGroup
	}
	return d.shipAssembledConsumerGroup
}

func (d *diContainer) OrderPaidConsumerGroup() sarama.ConsumerGroup {
	if d.orderPaidConsumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidConsumer.GroupID(),
			config.AppConfig().OrderPaidConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create order paid consumer group: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka order paid consumer group", func(ctx context.Context) error {
			return d.orderPaidConsumerGroup.Close()
		})
		d.orderPaidConsumerGroup = consumerGroup
	}
	return d.orderPaidConsumerGroup
}

func (d *diContainer) ShipAssembledConsumer() wrappedKafka.Consumer {
	if d.shipAssembledConsumer == nil {
		d.shipAssembledConsumer = wrappedKafkaConsumer.NewConsumer(
			d.ShipAssembledConsumerGroup(),
			[]string{
				config.AppConfig().ShipAssembledConsumer.Topic(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}
	return d.shipAssembledConsumer
}

func (d *diContainer) OrderPaidConsumer() wrappedKafka.Consumer {
	if d.orderPaidConsumer == nil {
		d.orderPaidConsumer = wrappedKafkaConsumer.NewConsumer(
			d.OrderPaidConsumerGroup(),
			[]string{
				config.AppConfig().OrderPaidConsumer.Topic(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}
	return d.orderPaidConsumer
}
