package auth

import (
	"github.com/ZanDattSu/star-factory/auth/internal/service"
	authV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/auth/v1"
)

type AuthApi struct {
	authV1.UnimplementedAuthServiceServer
	authService service.AuthService
}

func NewAuthApi(authService service.AuthService) *AuthApi {
	return &AuthApi{authService: authService}
}
