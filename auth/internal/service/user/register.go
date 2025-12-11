package user

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"go.uber.org/zap"

	"github.com/ZanDattSu/star-factory/auth/internal/model"
	"github.com/ZanDattSu/star-factory/platform/pkg/logger"
)

func (u *usersService) Register(ctx context.Context, userInfo model.UserInfo, password string) (*model.User, error) {
	hashedPassword, err := u.passwordHasher.HashAndSalt(password)
	if err != nil {
		logger.Error(ctx, "UserService.Register", zap.Error(err))
		return nil, err
	}

	user := &model.User{
		UUID:      uuid.New(),
		Info:      userInfo,
		Password:  hashedPassword,
		CreatedAt: lo.ToPtr(time.Now()),
	}

	err = u.usersRepository.Create(ctx, *user)
	if err != nil {
		logger.Error(ctx, "UserService.Register", zap.Error(err))
		return nil, err
	}

	return user, nil
}
