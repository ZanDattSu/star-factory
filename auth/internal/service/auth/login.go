package auth

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/ZanDattSu/star-factory/auth/internal/model"
)

func (a *authService) Login(ctx context.Context, login, password string) (uuid.UUID, error) {
	user, err := a.userRepository.Get(ctx, model.Filter{UserLogin: &login})
	if err != nil {
		return uuid.Nil, err
	}

	err = a.passwordHasher.ComparePasswords(user.Password, password)
	if err != nil {
		return uuid.Nil, err
	}

	now := time.Now()

	session := model.Session{
		UUID:      uuid.New(),
		User:      *user,
		CreatedAt: now,
		UpdatedAt: now,
		ExpiresAt: now.Add(a.cacheTTL),
	}

	err = a.sessionRepository.Create(ctx, session)
	if err != nil {
		return uuid.Nil, err
	}

	return session.UUID, nil
}
