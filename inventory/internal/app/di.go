package app

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	inventoryV1Api "github.com/ZanDattSu/star-factory/inventory/internal/api/v1/part"
	"github.com/ZanDattSu/star-factory/inventory/internal/config"
	"github.com/ZanDattSu/star-factory/inventory/internal/repository"
	inventoryRepository "github.com/ZanDattSu/star-factory/inventory/internal/repository/part/mongodb"
	"github.com/ZanDattSu/star-factory/inventory/internal/service"
	inventoryService "github.com/ZanDattSu/star-factory/inventory/internal/service/part"
	"github.com/ZanDattSu/star-factory/platform/pkg/closer"
	inventoryV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/inventory/v1"
)

type diContainer struct {
	inventoryV1Api inventoryV1.InventoryServiceServer
	partService    service.PartService
	partRepository repository.PartRepository

	mongoDBClient   *mongo.Client
	mongoDBDatabase *mongo.Database
}

func NewDIContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) InventoryV1Api(ctx context.Context) inventoryV1.InventoryServiceServer {
	if d.inventoryV1Api == nil {
		d.inventoryV1Api = inventoryV1Api.NewApi(d.PartService(ctx))
	}

	return d.inventoryV1Api
}

func (d *diContainer) PartService(ctx context.Context) service.PartService {
	if d.partService == nil {
		d.partService = inventoryService.NewService(d.PartRepository(ctx))
	}

	return d.partService
}

func (d *diContainer) PartRepository(ctx context.Context) repository.PartRepository {
	if d.partRepository == nil {
		d.partRepository = inventoryRepository.NewRepository(d.MongoDBDatabase(ctx)) //nolint:contextcheck
	}

	return d.partRepository
}

func (d *diContainer) MongoDBDatabase(ctx context.Context) *mongo.Database {
	if d.mongoDBDatabase == nil {
		d.mongoDBDatabase = d.MongoDBClient(ctx).Database(config.AppConfig().Mongo.DatabaseName())
	}

	return d.mongoDBDatabase
}

func (d *diContainer) MongoDBClient(ctx context.Context) *mongo.Client {
	if d.mongoDBClient == nil {
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
		if err != nil {
			panic(fmt.Sprintf("failed to connect to MongoDB: %s\n", err.Error()))
		}

		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			panic(fmt.Sprintf("failed to ping MongoDB: %v\n", err))
		}

		closer.AddNamed("MongoDB client", func(ctx context.Context) error { return client.Disconnect(ctx) })

		d.mongoDBClient = client
	}

	return d.mongoDBClient
}
