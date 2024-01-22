package rate_limiter_middleware

import (
	"context"
	"testing"
	"time"

	"github.com/gabrielsc1998/go-rate-limiter/configs"
	common_errors "github.com/gabrielsc1998/go-rate-limiter/internal/common/errors"
	"github.com/gabrielsc1998/go-rate-limiter/internal/common/infra/db/redis"
	rate_limiter_repo_redis "github.com/gabrielsc1998/go-rate-limiter/internal/rate-limiter/infra/db/redis"
)

var repo *rate_limiter_repo_redis.RateLimiterRepositoryRedis

func setupMiddleware() *RateLimiterMiddlewareInterface {
	config := redis.RedisClientConfig{Addr: "redis:6379", Password: "", DB: 0}
	client := redis.NewRedisClient(config)
	client.ClearAll(context.Background())
	repo = rate_limiter_repo_redis.NewRateLimiterRepositoryRedis(client)
	return NewRateLimiterMiddleware(&configs.Conf{
		RateLimiterMaxReqsIP:      3,
		RateLimiterBlockTimeIP:    1,
		RateLimiterBlockTimeToken: 1,
		ConfigTokens:              map[string]int{"token": 3},
	}, repo)
}

func TestHandleFirstReqIP(t *testing.T) {
	middleware := setupMiddleware()
	err := middleware.Handle("127.0.0.1", "")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	rateLimiter := repo.GetByIP("127.0.0.1")
	if rateLimiter == nil {
		t.Errorf("Expected rate limiter to be created")
	}
	if rateLimiter.Blocked {
		t.Errorf("Expected rate limiter to not be blocked")
	}
}

func TestHandleFirstReqToken(t *testing.T) {
	middleware := setupMiddleware()
	err := middleware.Handle("127.0.0.1", "token")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	rateLimiter := repo.GetByIP("127.0.0.1")
	if rateLimiter != nil {
		t.Errorf("Expected rate limiter to be nil")
	}

	rateLimiter = repo.GetByToken("token")
	if rateLimiter == nil {
		t.Errorf("Expected rate limiter to be created")
	}
}

func TestHandleAddReqIP(t *testing.T) {
	middleware := setupMiddleware()

	for i := 0; i < 3; i++ {
		err := middleware.Handle("127.0.0.1", "")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	}

	rateLimiter := repo.GetByIP("127.0.0.1")
	if rateLimiter == nil {
		t.Errorf("Expected rate limiter to be created")
	}
	if rateLimiter.Blocked {
		t.Errorf("Expected rate limiter to not be blocked")
	}
	if rateLimiter.TotalRequests != 3 {
		t.Errorf("Expected rate limiter to have 3 requests, got %d", rateLimiter.TotalRequests)
	}
}

func TestHandleAddReqToken(t *testing.T) {
	middleware := setupMiddleware()

	for i := 0; i < 3; i++ {
		err := middleware.Handle("127.0.0.1", "token")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	}

	rateLimiter := repo.GetByToken("token")
	if rateLimiter == nil {
		t.Errorf("Expected rate limiter to be created")
	}
	if rateLimiter.Blocked {
		t.Errorf("Expected rate limiter to not be blocked")
	}
	if rateLimiter.TotalRequests != 3 {
		t.Errorf("Expected rate limiter to have 3 requests, got %d", rateLimiter.TotalRequests)
	}
}

func TestMaxRequestForIP(t *testing.T) {
	middleware := setupMiddleware()

	for i := 0; i < 4; i++ {
		err := middleware.Handle("127.0.0.1", "")
		if i < 3 {
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
		} else {
			if !common_errors.Is(err, common_errors.ErrTooManyRequests) {
				t.Errorf("Expected error to be %s", common_errors.ErrTooManyRequests)
			}
		}
	}

	rateLimiter := repo.GetByIP("127.0.0.1")
	if rateLimiter == nil {
		t.Errorf("Expected rate limiter to be created")
	}
	if !rateLimiter.Blocked {
		t.Errorf("Expected rate limiter to be blocked")
	}
}

func TestMaxRequestForToken(t *testing.T) {
	middleware := setupMiddleware()

	for i := 0; i < 4; i++ {
		err := middleware.Handle("127.0.0.1", "token")
		if i < 3 {
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
		} else {
			if !common_errors.Is(err, common_errors.ErrTooManyRequests) {
				t.Errorf("Expected error to be %s", common_errors.ErrTooManyRequests)
			}
		}
	}

	rateLimiter := repo.GetByIP("127.0.0.1")
	if rateLimiter != nil {
		t.Errorf("Expected rate limiter to be nil")
	}

	rateLimiter = repo.GetByToken("token")
	if rateLimiter == nil {
		t.Errorf("Expected rate limiter to be created")
	}
	if !rateLimiter.Blocked {
		t.Errorf("Expected rate limiter to be blocked")
	}
}

func TestMaxRequestForIPAfterBlockTime(t *testing.T) {
	middleware := setupMiddleware()

	for i := 0; i < 4; i++ {
		err := middleware.Handle("127.0.0.1", "")
		if i < 3 {
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
		} else {
			if !common_errors.Is(err, common_errors.ErrTooManyRequests) {
				t.Errorf("Expected error to be %s", common_errors.ErrTooManyRequests)
			}
		}
	}

	rateLimiter := repo.GetByIP("127.0.0.1")
	if rateLimiter == nil {
		t.Errorf("Expected rate limiter to be created")
	}
	if !rateLimiter.Blocked {
		t.Errorf("Expected rate limiter to be blocked")
	}

	time.Sleep(1 * time.Second)

	err := middleware.Handle("127.0.0.1", "")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	rateLimiter = repo.GetByIP("127.0.0.1")
	if rateLimiter == nil {
		t.Errorf("Expected rate limiter to be created")
	}
	if rateLimiter.Blocked {
		t.Errorf("Expected rate limiter not to be blocked")
	}
}

func TestMaxRequestForTokenAfterBlockTime(t *testing.T) {
	middleware := setupMiddleware()

	for i := 0; i < 4; i++ {
		err := middleware.Handle("127.0.0.1", "token")
		if i < 3 {
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
		} else {
			if !common_errors.Is(err, common_errors.ErrTooManyRequests) {
				t.Errorf("Expected error to be %s", common_errors.ErrTooManyRequests)
			}
		}
	}

	rateLimiter := repo.GetByToken("token")
	if rateLimiter == nil {
		t.Errorf("Expected rate limiter to be created")
	}
	if !rateLimiter.Blocked {
		t.Errorf("Expected rate limiter to be blocked")
	}

	time.Sleep(1 * time.Second)

	err := middleware.Handle("token", "")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	rateLimiter = repo.GetByToken("token")
	if rateLimiter == nil {
		t.Errorf("Expected rate limiter to be created")
	}
	if rateLimiter.Blocked {
		t.Errorf("Expected rate limiter not to be blocked")
	}
}
