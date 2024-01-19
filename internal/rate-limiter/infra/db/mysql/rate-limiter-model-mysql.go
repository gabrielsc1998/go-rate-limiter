package rate_limiter_repo_mysql

import (
	"time"

	"gorm.io/gorm"
)

type RateLimiterModel struct {
	gorm.Model
	ID            int       `gorm:"primaryKey"`
	IP            string    `gorm:"unique"`
	Token         string    `gorm:"unique"`
	TotalRequests int       `gorm:"default:0"`
	Blocked       *bool     `gorm:"default:false"`
	FirstReqTime  time.Time `gorm:"not null"`
	MaxRequests   int       `gorm:"default:0"`
	BlockTime     float64   `gorm:"default:0"`
	BlockedAt     time.Time `gorm:"default:null"`
}
