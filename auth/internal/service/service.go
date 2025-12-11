package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/ZanDattSu/star-factory/auth/internal/model"
)

type AuthService interface {
	Login(ctx context.Context, login, password string) (uuid.UUID, error)
	Whoami(ctx context.Context, session uuid.UUID) (*model.Session, error)
}

type UserService interface {
	Register(ctx context.Context, userInfo model.UserInfo, password string) (*model.User, error)
	Get(ctx context.Context, uuid uuid.UUID) (*model.User, error)
}

type PasswordHasher interface {
	HashAndSalt(plainPassword string) (string, error)
	ComparePasswords(hashedPassword, plainPassword string) error
}
