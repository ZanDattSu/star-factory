package converter

import (
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ZanDattSu/star-factory/auth/internal/model"
	authV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/auth/v1"
	commonV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/common/v1"
)

func SessionUUIDToLoginResponse(uuid uuid.UUID) *authV1.LoginResponse {
	return &authV1.LoginResponse{
		SessionUuid: uuid.String(),
	}
}

func NotificationMethodToProto(notificationMethods []*model.NotificationMethod) []*commonV1.NotificationMethod {
	methods := make([]*commonV1.NotificationMethod, 0, len(notificationMethods))

	for _, notificationMethod := range notificationMethods {
		methods = append(methods, &commonV1.NotificationMethod{
			ProviderName: notificationMethod.ProviderName,
			Target:       notificationMethod.Target,
		})
	}

	return methods
}

func SessionToWhoamiResponse(session *model.Session) (*authV1.WhoamiResponse, error) {
	if session == nil {
		return nil, model.ErrSessionNotFound
	}

	if session.UUID == uuid.Nil {
		return nil, model.ErrInvalidSessionUUID
	}

	if session.User.UUID == uuid.Nil {
		return nil, model.ErrInvalidUserUUID
	}

	response := &authV1.WhoamiResponse{
		Session: &commonV1.Session{
			Uuid: session.UUID.String(),
		},
		User: &commonV1.User{
			Uuid: session.User.UUID.String(),
		},
	}

	createdAt := timestamppb.New(session.CreatedAt)
	if err := createdAt.CheckValid(); err != nil {
		return nil, err
	}
	response.Session.CreatedAt = createdAt

	updatedAt := timestamppb.New(session.UpdatedAt)
	if err := updatedAt.CheckValid(); err != nil {
		return nil, err
	}
	response.Session.UpdatedAt = updatedAt

	expiresAt := timestamppb.New(session.ExpiresAt)
	if err := expiresAt.CheckValid(); err != nil {
		return nil, err
	}
	response.Session.ExpiresAt = expiresAt

	if session.User.CreatedAt != nil {
		userCreatedAt := timestamppb.New(*session.User.CreatedAt)
		if err := userCreatedAt.CheckValid(); err != nil {
			return nil, err
		}
		response.User.CreatedAt = userCreatedAt
	}

	if session.User.UpdatedAt != nil {
		userUpdatedAt := timestamppb.New(*session.User.UpdatedAt)
		if err := userUpdatedAt.CheckValid(); err != nil {
			return nil, err
		}
		response.User.UpdatedAt = userUpdatedAt
	}

	response.User.Info = &commonV1.UserInfo{
		Login:               session.User.Info.Login,
		Email:               session.User.Info.Email,
		NotificationMethods: NotificationMethodToProto(session.User.Info.NotificationMethods),
	}

	return response, nil
}
