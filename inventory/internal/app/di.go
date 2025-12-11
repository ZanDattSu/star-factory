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
	grpcclient "github.com/ZanDattSu/star-factory/platform/pkg/grpc"
	"github.com/ZanDattSu/star-factory/platform/pkg/grpc/interceptor"
	authV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/auth/v1"
	inventoryV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/inventory/v1"
)

type diContainer struct {
	inventoryV1Api inventoryV1.InventoryServiceServer

	authClient      authV1.AuthServiceClient
	authInterceptor *interceptor.AuthInterceptor

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

func (d *diContainer) AuthClient(_ context.Context) authV1.AuthServiceClient {
	if d.authClient == nil {
		authConn, err := grpcclient.NewGRPCConnectWithoutSecure(config.AppConfig().Auth.AuthServiceAddress())
		if err != nil {
			panic(fmt.Sprintf(
				"Failed to connect to Auth gRPC service (%s): %v",
				config.AppConfig().Auth.AuthServicePort(),
				err,
			))
		}

		closer.AddNamed("Auth connection", func(ctx context.Context) error {
			return authConn.Close()
		})

		d.authClient = authV1.NewAuthServiceClient(authConn)
	}

	return d.authClient
}

func (d *diContainer) AuthInterceptor(ctx context.Context) *interceptor.AuthInterceptor {
	if d.authInterceptor == nil {
		d.authInterceptor = interceptor.NewAuthInterceptor(d.AuthClient(ctx))
	}

	return d.authInterceptor
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
