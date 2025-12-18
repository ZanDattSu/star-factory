package api

import (
	"github.com/google/uuid"

	"github.com/ZanDattSu/star-factory/notification/internal/model"
	commonV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/common/v1"
)

func UserFromProto(pb *commonV1.User) *model.User {
	if pb == nil {
		return nil
	}

	userUUID, err := uuid.Parse(pb.GetUuid())
	if err != nil {
		return nil
	}

	if pb.GetInfo() == nil {
		return nil
	}

	user := &model.User{
		UUID: userUUID,
		Info: model.UserInfo{
			Login:               pb.GetInfo().GetLogin(),
			Email:               pb.GetInfo().GetEmail(),
			NotificationMethods: NotificationMethodFromProto(pb.GetInfo().GetNotificationMethods()),
		},
	}

	if pb.GetCreatedAt() != nil {
		if err := pb.GetCreatedAt().CheckValid(); err != nil {
			return nil
		}
		t := pb.GetCreatedAt().AsTime()
		user.CreatedAt = &t
	}

	if pb.GetUpdatedAt() != nil {
		if err := pb.GetUpdatedAt().CheckValid(); err != nil {
			return nil
		}
		t := pb.GetUpdatedAt().AsTime()
		user.UpdatedAt = &t
	}

	return user
}

func NotificationMethodFromProto(
	methods []*commonV1.NotificationMethod,
) []*model.NotificationMethod {
	if len(methods) == 0 {
		return nil
	}

	result := make([]*model.NotificationMethod, 0, len(methods))
	for _, m := range methods {
		if m == nil {
			continue
		}

		result = append(result, &model.NotificationMethod{
			ProviderName: m.GetProviderName(),
			Target:       m.GetTarget(),
		})
	}

	return result
}
