package auth

import (
	"time"

	"github.com/ZanDattSu/star-factory/auth/internal/repository"
	"github.com/ZanDattSu/star-factory/auth/internal/service"
)

type authService struct {
	userRepository    repository.UserRepository
	sessionRepository repository.SessionRepository
	cacheTTL          time.Duration
	passwordHasher    service.PasswordHasher
}

func NewAuthService(
	usersRepository repository.UserRepository,
	cache repository.SessionRepository,
	cacheTTL time.Duration,
	passwordHasher service.PasswordHasher,
) *authService {
	return &authService{
		userRepository:    usersRepository,
		sessionRepository: cache,
		cacheTTL:          cacheTTL,
		passwordHasher:    passwordHasher,
	}
}
