package inmemory

import (
	"context"
	"fmt"

	"inventory/internal/model"
	"inventory/internal/repository/converter"
)

// GetPart возвращает деталь по UUID. Потокобезопасно.
func (r *repository) GetPart(_ context.Context, uuid string) (*model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	part, ok := r.parts[uuid]
	if !ok {
		return nil, fmt.Errorf(`part "%s" not found`, uuid)
	}

	return converter.PartToModel(part), nil
}
