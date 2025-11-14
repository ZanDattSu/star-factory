package postgresql

import (
	"github.com/jackc/pgx/v5/pgxpool"

	repo "github.com/ZanDattSu/star-factory/order/internal/repository"
)

// Компиляторная проверка: убеждаемся, что *repository реализует интерфейс OrderRepository.
var _ repo.OrderRepository = (*repository)(nil)

type repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *repository {
	return &repository{pool: pool}
}
