package rate_limiter_repo

import (
	"errors"

	"github.com/gabrielsc1998/go-rate-limiter/internal/common/infra/db/redis"
	rate_limiter "github.com/gabrielsc1998/go-rate-limiter/internal/rate-limiter/domain"
	rate_limiter_repo_redis "github.com/gabrielsc1998/go-rate-limiter/internal/rate-limiter/infra/db/redis"
)

type RateLimiterRepository interface {
	GetByIP(ip string) *rate_limiter.RateLimiter
	GetByToken(token string) *rate_limiter.RateLimiter
	Save(rateLimiter *rate_limiter.RateLimiter) error
}

type RateLimiterRepositoryConfig struct {
	Repo   string
	Inject interface{}
}

func NewRateLimiterRepository(config RateLimiterRepositoryConfig) (RateLimiterRepository, error) {
	if config.Repo == "redis" {
		repo := rate_limiter_repo_redis.NewRateLimiterRepositoryRedis(config.Inject.(*redis.RedisClient))
		return repo, nil
	}
	return nil, errors.New("invalid repository")
}
