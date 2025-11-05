package part

import (
	def "inventory/internal/repository"
	repoModel "inventory/internal/repository/model"
	"sync"
)

// Компиляторная проверка: убеждаемся, что *repository реализует интерфейс PartRepository.
var _ def.PartRepository = (*repository)(nil)

type repository struct {
	parts map[string]*repoModel.Part
	mu    sync.RWMutex
}

func NewRepository() *repository {
	return &repository{
		parts: make(map[string]*repoModel.Part),
	}
}
