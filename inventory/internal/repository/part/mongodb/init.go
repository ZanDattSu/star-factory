package mongodb

import (
	"context"

	"github.com/ZanDattSu/star-factory/inventory/internal/model"
	"github.com/ZanDattSu/star-factory/inventory/internal/service/part"
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
