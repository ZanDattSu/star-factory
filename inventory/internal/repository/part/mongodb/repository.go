package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	repo "github.com/ZanDattSu/star-factory/inventory/internal/repository"
)

var _ repo.PartRepository = (*repository)(nil)

type repository struct {
	collection *mongo.Collection
}

func NewRepository(db *mongo.Database) *repository {
	partsCollection := db.Collection("parts")

	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "uuid", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indexUUID, err := partsCollection.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		panic(fmt.Sprintf("Failed to create index %s: %s", indexUUID, err))
	}

	r := &repository{collection: partsCollection}
	r.InitTestData()
	return r
}
