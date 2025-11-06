package part

import (
	"inventory/internal/model"
	"inventory/internal/repository/converter"
)

// PutPart сохраняет деталь по UUID. Потокобезопасно.
func (r *repository) PutPart(uuid string, part *model.Part) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.parts[uuid] = converter.PartToRepoModel(part)
}
