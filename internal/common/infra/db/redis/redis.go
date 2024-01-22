package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

type RedisClientConfig struct {
	Addr     string
	Password string
	DB       int
}

func NewRedisClient(config RedisClientConfig) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})
	return &RedisClient{
		client: rdb,
	}
}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisClient) ClearAll(ctx context.Context) error {
	return r.client.FlushAll(ctx).Err()
}

type Input struct {
	Key    string
	Value  interface{}
	Expire float64
}

func (r *RedisClient) Set(ctx context.Context, input Input) error {
	err := r.client.Set(ctx, input.Key, input.Value, 0).Err()
	if err != nil {
		return err
	}
	if input.Expire != 0 {
		err = r.client.Expire(ctx, input.Key, time.Duration(input.Expire)*time.Second).Err()
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *RedisClient) Del(key string, ctx context.Context) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisClient) Close() error {
	return r.client.Close()
}
