package rate_limiter_repo_mysql

import (
	"github.com/gabrielsc1998/go-rate-limiter/internal/common/infra/db/mysql"
	rate_limiter "github.com/gabrielsc1998/go-rate-limiter/internal/rate-limiter/domain"
)

type RateLimiterRepositoryMySQL struct {
	client *mysql.MySQLDB
}

func NewRateLimiterRepositoryMySQL(client *mysql.MySQLDB) *RateLimiterRepositoryMySQL {
	return &RateLimiterRepositoryMySQL{
		client: client,
	}
}

func (r RateLimiterRepositoryMySQL) GetByIP(ip string) *rate_limiter.RateLimiter {
	return r.get("ip", ip)
}

func (r RateLimiterRepositoryMySQL) GetByToken(token string) *rate_limiter.RateLimiter {
	return r.get("token", token)
}

func (r RateLimiterRepositoryMySQL) get(key string, value string) *rate_limiter.RateLimiter {
	var rateLimiterFound RateLimiterModel
	err := r.client.DB.Where(key+" = ?", value).First(&rateLimiterFound).Error
	if err != nil {
		return nil
	}
	return rate_limiter.NewRateLimiter(rate_limiter.RateLimiterConfig{
		ID:            rateLimiterFound.ID,
		IP:            rateLimiterFound.IP,
		Token:         rateLimiterFound.Token,
		MaxRequests:   rateLimiterFound.MaxRequests,
		Blocked:       *rateLimiterFound.Blocked,
		TotalRequests: rateLimiterFound.TotalRequests,
		BlockTime:     rateLimiterFound.BlockTime,
		FirstReqTime:  rateLimiterFound.FirstReqTime,
		BlockedAt:     rateLimiterFound.BlockedAt,
	})
}

func (r RateLimiterRepositoryMySQL) Save(rateLimiter *rate_limiter.RateLimiter) error {
	if rateLimiter.Token != "" {
		rateLimiter.IP = ""
		rateLimiterFound := r.get("token", rateLimiter.Token)
		if rateLimiterFound != nil {
			rateLimiterFound = rateLimiter
			return r.client.DB.Updates(r.toModel(rateLimiterFound)).Error
		}
		return r.client.DB.Create(r.toModel(rateLimiter)).Error
	}

	rateLimiter.Token = ""

	rateLimiterFound := r.get("ip", rateLimiter.IP)
	if rateLimiterFound != nil {
		rateLimiterFound = rateLimiter
		return r.client.DB.Updates(r.toModel(rateLimiterFound)).Error
	}
	return r.client.DB.Create(r.toModel(rateLimiter)).Error
}

func (r RateLimiterRepositoryMySQL) toModel(rateLimiter *rate_limiter.RateLimiter) *RateLimiterModel {
	return &RateLimiterModel{
		ID:            rateLimiter.ID,
		IP:            rateLimiter.IP,
		Token:         rateLimiter.Token,
		MaxRequests:   rateLimiter.MaxRequests,
		Blocked:       &rateLimiter.Blocked,
		TotalRequests: rateLimiter.TotalRequests,
		BlockTime:     rateLimiter.BlockTime,
		FirstReqTime:  rateLimiter.FirstReqTime,
		BlockedAt:     rateLimiter.BlockedAt,
	}
}
