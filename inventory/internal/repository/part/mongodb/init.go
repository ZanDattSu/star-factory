package mongodb

import (
	"context"

	"inventory/internal/model"
	"inventory/internal/service/part"
)

func (r *repository) InitTestData() {
	var parts []*model.Part
	for i := 0; i < 10; i++ {
		parts = append(parts, part.RandomPart())
	}

	for _, p := range parts {
		err := r.PutPart(context.Background(), p.Uuid, p)
		if err != nil {
			continue
		}
	}
}
