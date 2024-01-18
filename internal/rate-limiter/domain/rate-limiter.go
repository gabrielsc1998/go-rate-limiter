package rate_limiter

import "time"

type RateLimiter struct {
	IP            string
	Token         string
	TotalRequests int
	Blocked       bool
	FirstReqTime  time.Time
	MaxRequests   int
	BlockTime     float64
}

type RateLimiterConfig struct {
	IP          string
	Token       string
	MaxRequests int
	BlockTime   float64
}

func NewRateLimiter(config RateLimiterConfig) *RateLimiter {
	return &RateLimiter{
		IP:            config.IP,
		Token:         config.Token,
		FirstReqTime:  time.Now(),
		TotalRequests: 1,
		MaxRequests:   config.MaxRequests,
		BlockTime:     config.BlockTime,
	}
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
}

func (r *RateLimiter) IsMaxRequestsReached() bool {
	return r.TotalRequests >= r.MaxRequests
}

func (r *RateLimiter) IsTimeOfOneSecondFinished() bool {
	return time.Since(r.FirstReqTime).Milliseconds() > 1000
}
