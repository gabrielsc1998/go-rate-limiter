package rate_limiter_repo_redis

import (
	"context"
	"encoding/json"

	"github.com/gabrielsc1998/go-rate-limiter/internal/common/infra/db/redis"
	rate_limiter "github.com/gabrielsc1998/go-rate-limiter/internal/rate-limiter/domain"
)

type RateLimiterRepositoryRedis struct {
	client *redis.RedisClient
}

func NewRateLimiterRepositoryRedis(client *redis.RedisClient) *RateLimiterRepositoryRedis {
	return &RateLimiterRepositoryRedis{
		client: client,
	}
}

func (r RateLimiterRepositoryRedis) GetByIP(ip string) *rate_limiter.RateLimiter {
	return r.get(ip)
}

func (r RateLimiterRepositoryRedis) GetByToken(token string) *rate_limiter.RateLimiter {
	return r.get(token)
}

func (r RateLimiterRepositoryRedis) get(key string) *rate_limiter.RateLimiter {
	ctx := context.Background()
	rateLimiterFound, _ := r.client.Get(ctx, key)
	if rateLimiterFound == "" {
		return nil
	}
	var rateLimiter rate_limiter.RateLimiter
	json.Unmarshal([]byte(rateLimiterFound), &rateLimiter)
	return &rateLimiter
}

func (r RateLimiterRepositoryRedis) Save(rateLimiter *rate_limiter.RateLimiter) error {
	ctx := context.Background()
	data, err := json.Marshal(rateLimiter)
	if err != nil {
		return err
	}

	var key string = rateLimiter.IP
	if rateLimiter.Token != "" {
		key = rateLimiter.Token
	}

	// TODO: Refactor this
	if rateLimiter.Blocked {
		return r.client.Set(ctx, redis.Input{
			Key:    key,
			Value:  data,
			Expire: rateLimiter.BlockTime,
		})
	}
	return r.client.Set(ctx, redis.Input{
		Key:   key,
		Value: data,
	})
}
