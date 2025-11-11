package mongodb

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"inventory/internal/model"
)

func (r *repository) ListParts(ctx context.Context) ([]*model.Part, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return []*model.Part{}, fmt.Errorf("error finding cursor: %w", err)
	}

	defer func() {
		if cerr := cursor.Close(ctx); cerr != nil {
			log.Printf("closing cursor error: %v\n", cerr)
		}
	}()

	var parts []*model.Part
	for cursor.Next(ctx) {
		var p model.Part
		if err := cursor.Decode(&p); err != nil {
			return nil, fmt.Errorf("decode part: %w", err)
		}
		parts = append(parts, &p)
	}

	return parts, nil
}
