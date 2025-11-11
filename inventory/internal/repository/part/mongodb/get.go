package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"inventory/internal/model"
)

func (r *repository) GetPart(ctx context.Context, uuid string) (*model.Part, error) {
	part := &model.Part{}
	err := r.collection.FindOne(ctx, bson.M{"uuid": uuid}).Decode(part)
	if err != nil {
		return nil, fmt.Errorf("failed to find part with uuid %s: %w", uuid, err)
	}

	return part, nil
}
