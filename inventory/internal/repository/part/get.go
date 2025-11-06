package part

import (
	"context"

	"inventory/internal/model"
	"inventory/internal/repository/converter"
)

// GetPart возвращает деталь по UUID. Потокобезопасно.
func (r *repository) GetPart(_ context.Context, uuid string) (*model.Part, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	part, ok := r.parts[uuid]

	return converter.PartToModel(part), ok
}
