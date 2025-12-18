package auth

import (
	"context"

	"github.com/ZanDattSu/star-factory/notification/internal/model"
)

type AuthClient interface {
	GetUser(ctx context.Context, userUUID string) (*model.User, error)
}
