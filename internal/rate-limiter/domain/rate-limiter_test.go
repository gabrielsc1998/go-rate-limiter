package rate_limiter

import (
	"testing"
	"time"
)

func TestNewRateLimiter(t *testing.T) {
	config := RateLimiterConfig{IP: "127.0.0.1", Token: "123456", MaxRequests: 10, BlockTime: 5}
	limiter := NewRateLimiter(config)

	if limiter.IP != "127.0.0.1" {
		t.Errorf("Expected IP to be '127.0.0.1', got '%s'", limiter.IP)
	}

	if limiter.Token != "123456" {
		t.Errorf("Expected Token to be '123456', got '%s'", limiter.Token)
	}

	if limiter.MaxRequests != 10 {
		t.Errorf("Expected MaxRequests to be 10, got '%d'", limiter.MaxRequests)
	}

	if limiter.TotalRequests != 1 {
		t.Errorf("Expected TotalRequests to be 1, got '%d'", limiter.TotalRequests)
	}

	if limiter.BlockTime != 5 {
		t.Errorf("Expected BlockTime to be 5, got '%f'", limiter.BlockTime)
	}
}

func TestAddRequest(t *testing.T) {
	config := RateLimiterConfig{IP: "127.0.0.1", Token: "123456", MaxRequests: 10}
	limiter := NewRateLimiter(config)

	limiter.AddRequest()

	if limiter.TotalRequests != 2 {
		t.Errorf("Expected TotalRequests to be 2, got '%d'", limiter.TotalRequests)
	}
}

func TestClearRequests(t *testing.T) {
	config := RateLimiterConfig{IP: "127.0.0.1", Token: "123456", MaxRequests: 10}
	limiter := NewRateLimiter(config)

	limiter.AddRequest()
	limiter.ClearRequests()

	if limiter.TotalRequests != 1 {
		t.Errorf("Expected TotalRequests to be 1, got '%d'", limiter.TotalRequests)
	}
}

func TestSetBlocked(t *testing.T) {
	config := RateLimiterConfig{IP: "127.0.0.1", Token: "123456", MaxRequests: 10}
	limiter := NewRateLimiter(config)

	limiter.SetBlocked()

	if limiter.Blocked != true {
		t.Errorf("Expected Blocked to be true, got '%t'", limiter.Blocked)
	}
}

func TestIsMaxRequestsReached(t *testing.T) {
	config := RateLimiterConfig{IP: "127.0.0.1", Token: "123456", MaxRequests: 2}
	limiter := NewRateLimiter(config)

	limiter.AddRequest()

	if limiter.IsMaxRequestsReached() != true {
		t.Errorf("Expected IsMaxRequestsReached to be true, got '%t'", limiter.IsMaxRequestsReached())
	}
}

func TestIsTimeOfOneSecondFinished(t *testing.T) {
	config := RateLimiterConfig{IP: "127.0.0.1", Token: "123456", MaxRequests: 2}
	limiter := NewRateLimiter(config)

	time.Sleep(2 * time.Second)

	if limiter.IsTimeOfOneSecondFinished() != true {
		t.Errorf("Expected IsTimeOfOneSecondFinished to be true, got '%t'", limiter.IsTimeOfOneSecondFinished())
	}
}
