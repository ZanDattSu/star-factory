package user

import (
	"github.com/ZanDattSu/star-factory/auth/internal/repository"
	"github.com/ZanDattSu/star-factory/auth/internal/service"
)

type usersService struct {
	usersRepository repository.UserRepository
	passwordHasher  service.PasswordHasher
}

func NewUsersService(
	usersRepository repository.UserRepository,
	passwordHasher service.PasswordHasher,
) *usersService {
	return &usersService{
		usersRepository: usersRepository,
		passwordHasher:  passwordHasher,
	}
}
