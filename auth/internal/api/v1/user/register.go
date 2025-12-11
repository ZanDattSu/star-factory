package user

import (
	"context"

	"go.uber.org/zap"

	"github.com/ZanDattSu/star-factory/auth/internal/converter"
	"github.com/ZanDattSu/star-factory/platform/pkg/logger"
	userV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/user/v1"
)

func (u *UserApi) Register(ctx context.Context, request *userV1.RegisterRequest) (*userV1.RegisterResponse, error) {
	user := converter.UserFromRequest(request)

	user, err := u.userService.Register(ctx, user.Info, user.Password)
	if err != nil {
		logger.Error(ctx, "UserImplementation.Register", zap.Error(err))
		return nil, err
	}

	return converter.UserToRegisterResponse(user)
}
