package inmemory

import (
	"context"

	"github.com/ZanDattSu/star-factory/inventory/internal/model"
	"github.com/ZanDattSu/star-factory/inventory/internal/repository/converter"
)

// PutPart сохраняет деталь по UUID. Потокобезопасно.
func (r *repository) PutPart(_ context.Context, uuid string, part *model.Part) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.parts[uuid] = converter.PartToRepoModel(part)
	return nil
}
