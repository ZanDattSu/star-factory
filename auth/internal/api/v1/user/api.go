package user

import (
	"github.com/ZanDattSu/star-factory/auth/internal/service"
	userV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/user/v1"
)

type UserApi struct {
	userV1.UnimplementedUserServiceServer
	userService service.UserService
}

func NewUserApi(userService service.UserService) *UserApi {
	return &UserApi{userService: userService}
}
