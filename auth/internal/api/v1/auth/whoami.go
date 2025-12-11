package auth

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/ZanDattSu/star-factory/auth/internal/converter"
	"github.com/ZanDattSu/star-factory/platform/pkg/logger"
	authV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/auth/v1"
)

func (a *AuthApi) Whoami(ctx context.Context, request *authV1.WhoamiRequest) (*authV1.WhoamiResponse, error) {
	sessionUUID, err := uuid.Parse(request.SessionUuid)
	if err != nil {
		logger.Error(ctx, "AuthImplementation.Whoami", zap.Error(err))
		return nil, err
	}

	session, err := a.authService.Whoami(ctx, sessionUUID)
	if err != nil {
		logger.Error(ctx, "AuthImplementation.Whoami", zap.Error(err))
		return nil, err
	}

	return converter.SessionToWhoamiResponse(session)
}
