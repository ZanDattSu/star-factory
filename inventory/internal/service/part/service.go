package part

import (
	"inventory/internal/repository"
	srvc "inventory/internal/service"
)

// Компиляторная проверка: убеждаемся, что *service реализует интерфейс PartService.
var _ srvc.PartService = (*service)(nil)

type service struct {
	repository repository.PartRepository
}

func NewService(repository repository.PartRepository) *service {
	return &service{
		repository: repository,
	}
}
