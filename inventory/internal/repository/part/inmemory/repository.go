package inmemory

import (
	"sync"

	repo "github.com/ZanDattSu/star-factory/inventory/internal/repository"
	repoModel "github.com/ZanDattSu/star-factory/inventory/internal/repository/model"
)

// Компиляторная проверка: убеждаемся, что *repository реализует интерфейс PartRepository.
var _ repo.PartRepository = (*repository)(nil)

type repository struct {
	parts map[string]*repoModel.Part
	mu    sync.RWMutex
}

// NewRepository создаёт и возвращает указатель на repository.
// Возврат конкретного типа (*repository), а не интерфейса, позволяет
// избежать потери методов и облегчает тестирование и расширение реализации.
func NewRepository() *repository {
	return &repository{
		parts: make(map[string]*repoModel.Part),
	}
}
