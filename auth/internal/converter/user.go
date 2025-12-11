package converter

import (
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ZanDattSu/star-factory/auth/internal/model"
	commonV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/common/v1"
	userV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/user/v1"
)

func UserToGetUserResponse(user *model.User) (*userV1.GetUserResponse, error) {
	if user.UUID == uuid.Nil {
		return nil, model.ErrInvalidUserUUID
	}

	response := &userV1.GetUserResponse{
		User: &commonV1.User{
			Uuid: user.UUID.String(),
		},
	}

	response.User.Info = &commonV1.UserInfo{
		Login:               user.Info.Login,
		Email:               user.Info.Email,
		NotificationMethods: NotificationMethodToProto(user.Info.NotificationMethods),
	}

	if user.CreatedAt != nil {
		createdAt := timestamppb.New(*user.CreatedAt)
		if err := createdAt.CheckValid(); err != nil {
			return nil, err
		}
		response.User.CreatedAt = createdAt
	}

	if user.UpdatedAt != nil {
		updatedAt := timestamppb.New(*user.UpdatedAt)
		if err := updatedAt.CheckValid(); err != nil {
			return nil, err
		}
		response.User.UpdatedAt = updatedAt
	}

	return response, nil
}

func UserFromRequest(request *userV1.RegisterRequest) *model.User {
	user := &model.User{
		Password: request.GetPassword(),
		Info:     model.UserInfo{},
	}

	if request.Info != nil {
		user.Info.Login = request.Info.Login
		user.Info.Email = request.Info.Email
		user.Info.NotificationMethods = NotificationMethodsFromProto(request.Info.NotificationMethods)
	}

	return user
}

func NotificationMethodsFromProto(notificationMethods []*commonV1.NotificationMethod) []*model.NotificationMethod {
	methods := make([]*model.NotificationMethod, 0, len(notificationMethods))

	for _, notificationMethod := range notificationMethods {
		methods = append(methods, &model.NotificationMethod{
			ProviderName: notificationMethod.ProviderName,
			Target:       notificationMethod.Target,
		})
	}

	return methods
}

func UserToRegisterResponse(user *model.User) (*userV1.RegisterResponse, error) {
	if user.UUID == uuid.Nil {
		return nil, model.ErrInvalidUserUUID
	}

	response := &userV1.RegisterResponse{
		UserUuid: user.UUID.String(),
	}

	return response, nil
}
