package rate_limiter

import (
	"testing"
	"time"
)

func TestNewRateLimiter(t *testing.T) {
	date := time.Now()
	config := RateLimiterConfig{
		IP:            "127.0.0.1",
		Token:         "123456",
		MaxRequests:   10,
		BlockTime:     5,
		TotalRequests: 5,
		Blocked:       true,
		FirstReqTime:  date,
		BlockedAt:     date.Add(10 * time.Second),
	}
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
	if limiter.TotalRequests != 5 {
		t.Errorf("Expected TotalRequests to be 5, got '%d'", limiter.TotalRequests)
	}
	if limiter.BlockTime != 5 {
		t.Errorf("Expected BlockTime to be 5, got '%f'", limiter.BlockTime)
	}
	if limiter.Blocked != true {
		t.Errorf("Expected Blocked to be true, got '%t'", limiter.Blocked)
	}
	if limiter.FirstReqTime != date {
		t.Errorf("Expected FirstReqTime to be '%s', got '%s'", date, limiter.FirstReqTime)
	}
	if limiter.BlockedAt != date.Add(10*time.Second) {
		t.Errorf("Expected BlockedAt to be '%s', got '%s'", date.Add(10*time.Second), limiter.BlockedAt)
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
	config := RateLimiterConfig{TotalRequests: 5}
	limiter := NewRateLimiter(config)

	if limiter.TotalRequests != 5 {
		t.Errorf("Expected TotalRequests to be 5, got '%d'", limiter.TotalRequests)
	}

	limiter.ClearRequests()

	if limiter.TotalRequests != 0 {
		t.Errorf("Expected TotalRequests to be 0, got '%d'", limiter.TotalRequests)
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

func TestUnblock(t *testing.T) {
	config := RateLimiterConfig{TotalRequests: 10}
	limiter := NewRateLimiter(config)

	limiter.SetBlocked()

	if limiter.Blocked != true {
		t.Errorf("Expected Blocked to be true, got '%t'", limiter.Blocked)
	}

	limiter.Unblock()

	if limiter.Blocked != false {
		t.Errorf("Expected Blocked to be false, got '%t'", limiter.Blocked)
	}
	if limiter.TotalRequests != 0 {
		t.Errorf("Expected TotalRequests to be 0, got '%d'", limiter.TotalRequests)
	}
}

func TestIsMaxRequestsReached(t *testing.T) {
	config := RateLimiterConfig{MaxRequests: 2}
	limiter := NewRateLimiter(config)

	limiter.AddRequest()

	if limiter.IsMaxRequestsReached() != true {
		t.Errorf("Expected IsMaxRequestsReached to be true, got '%t'", limiter.IsMaxRequestsReached())
	}
}

func TestIsTimeOfOneSecondFinished(t *testing.T) {
	config := RateLimiterConfig{IP: "127.0.0.1", Token: "123456", MaxRequests: 2}
	limiter := NewRateLimiter(config)
	limiter.FirstReqTime = time.Now().Add(-2 * time.Second)

	if limiter.IsTimeOfOneSecondFinished() != true {
		t.Errorf("Expected IsTimeOfOneSecondFinished to be true, got '%t'", limiter.IsTimeOfOneSecondFinished())
	}
}

func TestIsTimeOfBlockFinished(t *testing.T) {
	config := RateLimiterConfig{Blocked: true, BlockTime: 1}
	limiter := NewRateLimiter(config)
	limiter.BlockedAt = time.Now().Add(-2 * time.Second)

	if limiter.IsTimeOfBlockFinished() != true {
		t.Errorf("Expected IsTimeOfOneSecondFinished to be true, got '%t'", limiter.IsTimeOfOneSecondFinished())
	}
}
