package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/ZanDattSu/star-factory/auth/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user model.User) error
	Get(ctx context.Context, filter model.Filter) (*model.User, error)
}

type SessionRepository interface {
	Get(ctx context.Context, sessionUUID uuid.UUID) (*model.Session, error)
	Create(ctx context.Context, session model.Session) error
}
