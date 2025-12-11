package hasher

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/ZanDattSu/star-factory/auth/internal/model"
)

type bcryptHasher struct{}

func NewBcryptHasher() *bcryptHasher {
	return &bcryptHasher{}
}

func (h *bcryptHasher) HashAndSalt(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (h *bcryptHasher) ComparePasswords(hashedPassword, plainPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {
		return model.ErrInvalidCredentials
	}

	return nil
}
