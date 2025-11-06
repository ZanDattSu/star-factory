package part

import (
	"inventory/internal/model"
	"inventory/internal/repository/converter"
	repoModel "inventory/internal/repository/model"
)

// Values возвращает все детали. Потокобезопасно.
func (r *repository) Values() []*model.Part {
	r.mu.RLock()
	defer r.mu.RUnlock()

	parts := make([]*repoModel.Part, 0, len(r.parts))
	for _, part := range r.parts {
		parts = append(parts, part)
	}
	return converter.PartsToModel(parts)
}
