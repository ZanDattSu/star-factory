package session

import (
	"context"

	"github.com/ZanDattSu/star-factory/auth/internal/model"
	"github.com/ZanDattSu/star-factory/auth/internal/repository/converter"
)

func (s *sessionRepository) Create(ctx context.Context, session model.Session) error {
	cacheKey := s.getCacheKey(session.UUID.String())

	repoSession, err := converter.SessionToRepoModel(session)
	if err != nil {
		return err
	}

	err = s.cache.HashSet(ctx, cacheKey, repoSession)
	if err != nil {
		return err
	}

	cacheTTL := session.ExpiresAt.Sub(session.CreatedAt)

	return s.cache.Expire(ctx, cacheKey, cacheTTL)
}
