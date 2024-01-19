package rate_limiter

import (
	"time"
)

type RateLimiter struct {
	ID            int
	IP            string
	Token         string
	TotalRequests int
	Blocked       bool
	FirstReqTime  time.Time
	MaxRequests   int
	BlockTime     float64
	BlockedAt     time.Time
}

type RateLimiterConfig struct {
	ID            int
	IP            string
	Token         string
	TotalRequests int
	Blocked       bool
	FirstReqTime  time.Time
	MaxRequests   int
	BlockTime     float64
	BlockedAt     time.Time
}

func NewRateLimiter(config RateLimiterConfig) *RateLimiter {
	rateLimiter := &RateLimiter{
		ID:            config.ID,
		IP:            config.IP,
		Token:         config.Token,
		FirstReqTime:  config.FirstReqTime,
		TotalRequests: config.TotalRequests,
		Blocked:       config.Blocked,
		MaxRequests:   config.MaxRequests,
		BlockTime:     config.BlockTime,
		BlockedAt:     config.BlockedAt,
	}
	if config.FirstReqTime.IsZero() {
		rateLimiter.FirstReqTime = time.Now()
	}
	if config.TotalRequests == 0 {
		rateLimiter.TotalRequests = 1
	}
	return rateLimiter
}

func (r *RateLimiter) AddRequest() {
	r.TotalRequests++
}

func (r *RateLimiter) ClearRequests() {
	r.TotalRequests = 1
	r.FirstReqTime = time.Now()
}

func (r *RateLimiter) SetBlocked() {
	r.Blocked = true
	r.BlockedAt = time.Now()
}

func (r *RateLimiter) IsTimeOfBlockFinished() bool {
	if !r.Blocked {
		return true
	}
	return time.Since(r.BlockedAt).Seconds() > r.BlockTime
}

func (r *RateLimiter) Unblock() {
	r.Blocked = false
	r.ClearRequests()
}

func (r *RateLimiter) IsMaxRequestsReached() bool {
	return r.TotalRequests >= r.MaxRequests
}

func (r *RateLimiter) IsTimeOfOneSecondFinished() bool {
	return time.Since(r.FirstReqTime).Milliseconds() > 1000
}
