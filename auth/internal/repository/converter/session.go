package converter

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/ZanDattSu/star-factory/auth/internal/model"
	repoModel "github.com/ZanDattSu/star-factory/auth/internal/repository/model"
)

func SessionToRepoModel(session model.Session) (*repoModel.Session, error) {
	if session.UUID == uuid.Nil {
		return nil, model.ErrInvalidSessionUUID
	}

	if session.User.UUID == uuid.Nil {
		return nil, model.ErrInvalidUserUUID
	}

	result := &repoModel.Session{
		UUID:        session.UUID.String(),
		UserUUID:    session.User.UUID.String(),
		CreatedAtNs: session.CreatedAt.UnixNano(),
		UpdatedAtNs: session.UpdatedAt.UnixNano(),
		ExpiresAtNs: session.ExpiresAt.UnixNano(),
	}

	result.Login = session.User.Info.Login
	result.Email = session.User.Info.Email

	return result, nil
}

func SessionToModel(session repoModel.Session) (*model.Session, error) {
	sessionUUID, err := uuid.Parse(session.UUID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", model.ErrInvalidSessionUUID, err)
	}

	userUUID, err := uuid.Parse(session.UserUUID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", model.ErrInvalidUserUUID, err)
	}

	result := &model.Session{
		UUID: sessionUUID,
		User: model.User{
			UUID: userUUID,
			Info: model.UserInfo{
				Login: session.Login,
				Email: session.Email,
			},
		},
		CreatedAt: time.Unix(0, session.CreatedAtNs),
		UpdatedAt: time.Unix(0, session.UpdatedAtNs),
		ExpiresAt: time.Unix(0, session.ExpiresAtNs),
	}

	return result, nil
}
