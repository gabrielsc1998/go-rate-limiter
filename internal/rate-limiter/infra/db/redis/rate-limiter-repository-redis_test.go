package rate_limiter_repo_redis

import (
	"testing"
	"time"

	"github.com/gabrielsc1998/go-rate-limiter/internal/common/infra/db/redis"
	rate_limiter "github.com/gabrielsc1998/go-rate-limiter/internal/rate-limiter/domain"
)

var repo RateLimiterRepositoryRedis

func init() {
	config := redis.RedisClientConfig{Addr: "redis:6379", Password: "", DB: 0}
	client := redis.NewRedisClient(config)
	repo = RateLimiterRepositoryRedis{
		client: client,
	}
}

func createRateLimiter() *rate_limiter.RateLimiter {
	return rate_limiter.NewRateLimiter(rate_limiter.RateLimiterConfig{
		IP:            "192.168.0.1",
		MaxRequests:   10,
		Blocked:       false,
		BlockTime:     5,
		TotalRequests: 0,
	})
}

func TestSave(t *testing.T) {
	err := repo.Save(createRateLimiter())
	if err != nil {
		t.Errorf("Expected error to be nil, got '%s'", err.Error())
	}
}

func TestGetByIP(t *testing.T) {
	rateLimiter := createRateLimiter()
	err := repo.Save(rateLimiter)
	if err != nil {
		t.Errorf("Expected error to be nil, got '%s'", err.Error())
	}

	rateLimiterFound := repo.GetByIP(rateLimiter.IP)
	if rateLimiterFound == nil {
		t.Errorf("Expected rateLimiterFound to be not nil")
	}
}

func TestGetByToken(t *testing.T) {
	rateLimiter := createRateLimiter()
	rateLimiter.Token = "token"

	err := repo.Save(rateLimiter)
	if err != nil {
		t.Errorf("Expected error to be nil, got '%s'", err.Error())
	}

	rateLimiterFound := repo.GetByToken(rateLimiter.Token)
	if rateLimiterFound == nil {
		t.Errorf("Expected rateLimiterFound to be not nil")
	}
}

func TestGetByIPNotFound(t *testing.T) {
	rateLimiterFound := repo.GetByIP("255.255.255.255")
	if rateLimiterFound != nil {
		t.Errorf("Expected rateLimiterFound to be nil")
	}
}

func TestGetByTokenNotFound(t *testing.T) {
	rateLimiterFound := repo.GetByToken("7f5bea77-a80e-4d3a-b036-ab87d3df332d")
	if rateLimiterFound != nil {
		t.Errorf("Expected rateLimiterFound to be nil")
	}
}

func TestSaveWithBlockTime(t *testing.T) {
	rateLimiter := createRateLimiter()
	rateLimiter.Blocked = true
	rateLimiter.BlockTime = 1

	err := repo.Save(rateLimiter)
	if err != nil {
		t.Errorf("Expected error to be nil, got '%s'", err.Error())
	}

	rateLimiterFound := repo.GetByIP(rateLimiter.IP)
	if rateLimiterFound == nil {
		t.Errorf("Expected rateLimiterFound to be not nil")
	}

	time.Sleep(2 * time.Second)

	rateLimiterFound = repo.GetByIP(rateLimiter.IP)
	if rateLimiterFound != nil {
		t.Errorf("Expected rateLimiterFound to be nil")
	}
}
