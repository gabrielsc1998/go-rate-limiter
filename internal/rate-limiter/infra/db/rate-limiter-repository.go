package rate_limiter_repo

import (
	"errors"

	"github.com/gabrielsc1998/go-rate-limiter/internal/common/infra/db/mysql"
	"github.com/gabrielsc1998/go-rate-limiter/internal/common/infra/db/redis"
	rate_limiter "github.com/gabrielsc1998/go-rate-limiter/internal/rate-limiter/domain"
	rate_limiter_repo_mysql "github.com/gabrielsc1998/go-rate-limiter/internal/rate-limiter/infra/db/mysql"
	rate_limiter_repo_redis "github.com/gabrielsc1998/go-rate-limiter/internal/rate-limiter/infra/db/redis"
)

type RateLimiterRepository interface {
	GetByIP(ip string) *rate_limiter.RateLimiter
	GetByToken(token string) *rate_limiter.RateLimiter
	Save(rateLimiter *rate_limiter.RateLimiter) error
}

type Config struct {
	Repo   string
	Inject interface{}
}

func NewRateLimiterRepository(config Config) (RateLimiterRepository, error) {
	if config.Repo == "redis" {
		repo := rate_limiter_repo_redis.NewRateLimiterRepositoryRedis(config.Inject.(*redis.RedisClient))
		return repo, nil
	}
	if config.Repo == "mysql" {
		repo := rate_limiter_repo_mysql.NewRateLimiterRepositoryMySQL(config.Inject.(*mysql.MySQLDB))
		return repo, nil
	}
	return nil, errors.New("invalid repository")
}
