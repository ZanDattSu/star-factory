package converter

import (
	"github.com/ZanDattSu/star-factory/auth/internal/model"
	repoModel "github.com/ZanDattSu/star-factory/auth/internal/repository/model"
)

func UserToRepoModel(user model.User) repoModel.User {
	repoUser := repoModel.User{
		UUID:      user.UUID,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	repoUser.Info = repoModel.UserInfo{
		Login: user.Info.Login,
		Email: user.Info.Email,
	}

	// Добавляем методы уведомлений
	if user.Info.NotificationMethods != nil {
		repoUser.Info.NotificationMethods = make([]repoModel.NotificationMethod, len(user.Info.NotificationMethods))

		for i, nm := range user.Info.NotificationMethods {
			repoUser.Info.NotificationMethods[i] = repoModel.NotificationMethod{
				ProviderName: nm.ProviderName,
				Target:       nm.Target,
			}
		}
	}

	return repoUser
}
