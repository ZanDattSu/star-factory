package session

import (
	"context"
	"errors"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/google/uuid"

	"github.com/ZanDattSu/star-factory/auth/internal/model"
	"github.com/ZanDattSu/star-factory/auth/internal/repository/converter"
	repoModel "github.com/ZanDattSu/star-factory/auth/internal/repository/model"
)

func (s *sessionRepository) Get(ctx context.Context, sessionUUID uuid.UUID) (*model.Session, error) {
	cacheKey := s.getCacheKey(sessionUUID.String())

	values, err := s.cache.HGetAll(ctx, cacheKey)
	if err != nil {
		if errors.Is(err, redigo.ErrNil) {
			return nil, model.ErrSessionNotFound
		}
		return nil, err
	}

	if len(values) == 0 {
		return nil, model.ErrSessionNotFound
	}

	var session repoModel.Session
	err = redigo.ScanStruct(values, &session)
	if err != nil {
		return nil, err
	}

	return converter.SessionToModel(session)
}
