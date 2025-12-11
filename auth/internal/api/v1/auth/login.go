package auth

import (
	"context"

	"go.uber.org/zap"

	"github.com/ZanDattSu/star-factory/auth/internal/converter"
	"github.com/ZanDattSu/star-factory/platform/pkg/logger"
	auth_v1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/auth/v1"
)

func (a *AuthApi) Login(ctx context.Context, request *auth_v1.LoginRequest) (*auth_v1.LoginResponse, error) {
	sessionUUID, err := a.authService.Login(ctx, request.Login, request.Password)
	if err != nil {
		logger.Error(ctx, "Login", zap.Error(err))
		return nil, err
	}

	return converter.SessionUUIDToLoginResponse(sessionUUID), nil
}
