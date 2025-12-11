package session

import (
	"fmt"

	"github.com/ZanDattSu/star-factory/platform/pkg/cache"
)

const cacheKeyPrefix = "session:"

type sessionRepository struct {
	cache cache.RedisClient
}

func NewSessionRepository(redisClient cache.RedisClient) *sessionRepository {
	return &sessionRepository{
		cache: redisClient,
	}
}

func (s *sessionRepository) getCacheKey(uuid string) string {
	return fmt.Sprintf("%s%s", cacheKeyPrefix, uuid)
}
