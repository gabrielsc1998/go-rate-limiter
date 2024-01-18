package rate_limiter_middleware

import (
	"github.com/gabrielsc1998/go-rate-limiter/configs"
	common_errors "github.com/gabrielsc1998/go-rate-limiter/internal/common/errors"
	rate_limiter "github.com/gabrielsc1998/go-rate-limiter/internal/rate-limiter/domain"
	rate_limiter_repo "github.com/gabrielsc1998/go-rate-limiter/internal/rate-limiter/infra/db"
)

type RateLimiterMiddlewareInterface struct {
	configs      *configs.Conf
	repo         rate_limiter_repo.RateLimiterRepository
	configTokens map[string]int
}

func NewRateLimiterMiddleware(configs *configs.Conf, configTokens map[string]int, repo rate_limiter_repo.RateLimiterRepository) *RateLimiterMiddlewareInterface {
	return &RateLimiterMiddlewareInterface{
		repo:         repo,
		configs:      configs,
		configTokens: configTokens,
	}
}

func (r *RateLimiterMiddlewareInterface) Handle(ip string, token string) error {
	var rateLimiter *rate_limiter.RateLimiter = nil

	if token != "" {
		rateLimiter = r.repo.GetByToken(token)
	} else {
		rateLimiter = r.repo.GetByIP(ip)
	}

	if rateLimiter == nil {
		maxReqs := r.configs.RateLimiterMaxReqsIP
		if token != "" {
			maxReqs = r.configTokens[token]
			if maxReqs == 0 {
				maxReqs = 10
			}
		}
		blockTime := r.configs.RateLimiterBlockTimeIP
		if token != "" {
			blockTime = r.configs.RateLimiterBlockTimeToken
		}
		r.repo.Save(rate_limiter.NewRateLimiter(rate_limiter.RateLimiterConfig{
			IP:          ip,
			Token:       token,
			MaxRequests: maxReqs,
			BlockTime:   blockTime,
		}))
		return nil
	}

	if rateLimiter.Blocked {
		return common_errors.ErrTooManyRequests
	}
	if rateLimiter.IsTimeOfOneSecondFinished() {
		rateLimiter.ClearRequests()
		r.repo.Save(rateLimiter)
	} else if rateLimiter.IsMaxRequestsReached() {
		rateLimiter.SetBlocked()
		r.repo.Save(rateLimiter)
		return common_errors.ErrTooManyRequests
	} else {
		rateLimiter.AddRequest()
		r.repo.Save(rateLimiter)
	}
	return nil
}
