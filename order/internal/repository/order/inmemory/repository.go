package inmemory

import (
	"sync"

	repo "github.com/ZanDattSu/star-factory/order/internal/repository"
	repoModel "github.com/ZanDattSu/star-factory/order/internal/repository/model"
)

// Компиляторная проверка: убеждаемся, что *repository реализует интерфейс OrderRepository.
var _ repo.OrderRepository = (*repository)(nil)

type repository struct {
	orders map[string]*repoModel.Order
	mu     sync.RWMutex
}

// NewRepository создаёт и возвращает указатель на repository.
// Возврат конкретного типа (*repository), а не интерфейса, позволяет
// избежать потери методов и облегчает тестирование и расширение реализации.
func NewRepository() *repository {
	return &repository{
		orders: make(map[string]*repoModel.Order),
	}
}
