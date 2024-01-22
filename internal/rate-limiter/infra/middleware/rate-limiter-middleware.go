package rate_limiter_middleware

import (
	"github.com/gabrielsc1998/go-rate-limiter/configs"
	common_errors "github.com/gabrielsc1998/go-rate-limiter/internal/common/errors"
	rate_limiter "github.com/gabrielsc1998/go-rate-limiter/internal/rate-limiter/domain"
	rate_limiter_repo "github.com/gabrielsc1998/go-rate-limiter/internal/rate-limiter/infra/db"
)

type RateLimiterMiddlewareInterface struct {
	configs *configs.Conf
	repo    rate_limiter_repo.RateLimiterRepository
}

func NewRateLimiterMiddleware(configs *configs.Conf, repo rate_limiter_repo.RateLimiterRepository) *RateLimiterMiddlewareInterface {
	return &RateLimiterMiddlewareInterface{
		repo:    repo,
		configs: configs,
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
			maxReqs = r.configs.ConfigTokens[token]
			if maxReqs == 0 {
				maxReqs = 10
			}
		}
		blockTime := r.configs.RateLimiterBlockTimeIP
		if token != "" {
			blockTime = r.configs.RateLimiterBlockTimeToken
		}
		err := r.repo.Save(rate_limiter.NewRateLimiter(rate_limiter.RateLimiterConfig{
			IP:          ip,
			Token:       token,
			MaxRequests: maxReqs,
			BlockTime:   blockTime,
		}))
		if err != nil {
			return err
		}
		return nil
	}

	if rateLimiter.Blocked {
		if rateLimiter.IsTimeOfBlockFinished() {
			rateLimiter.Unblock()
			r.repo.Save(rateLimiter)
			return nil
		}
		return common_errors.ErrTooManyRequests
	}
	if rateLimiter.IsTimeOfOneSecondFinished() {
		rateLimiter.ClearRequests()
		rateLimiter.AddRequest()
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
