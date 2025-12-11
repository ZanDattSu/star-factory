package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/ZanDattSu/star-factory/auth/internal/model"
)

func (u *usersService) Get(ctx context.Context, userUUID uuid.UUID) (*model.User, error) {
	if userUUID == uuid.Nil {
		return nil, model.ErrInvalidUserUUID
	}

	filter := model.Filter{
		UserUUID: lo.ToPtr(userUUID.String()),
	}

	user, err := u.usersRepository.Get(ctx, filter)
	if err != nil {
		return nil, err
	}

	return user, nil
}
