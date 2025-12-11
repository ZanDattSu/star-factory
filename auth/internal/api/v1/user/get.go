package user

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/ZanDattSu/star-factory/auth/internal/converter"
	"github.com/ZanDattSu/star-factory/platform/pkg/logger"
	userV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/user/v1"
)

func (u *UserApi) GetUser(ctx context.Context, request *userV1.GetUserRequest) (*userV1.GetUserResponse, error) {
	userUUID, err := uuid.Parse(request.GetUserUuid())
	if err != nil {
		logger.Error(ctx, "GetUser", zap.Error(err))
		return nil, err
	}

	user, err := u.userService.Get(ctx, userUUID)
	if err != nil {
		logger.Error(ctx, "GetUser", zap.Error(err))
		return nil, err
	}

	return converter.UserToGetUserResponse(user)
}
