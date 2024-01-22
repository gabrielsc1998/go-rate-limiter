package rate_limiter_repo_mysql

import (
	"testing"

	"github.com/gabrielsc1998/go-rate-limiter/internal/common/infra/db/mysql"
	rate_limiter "github.com/gabrielsc1998/go-rate-limiter/internal/rate-limiter/domain"
)

func connectDb() (*mysql.MySQLDB, error) {
	db := mysql.NewMySQLDBConnection()
	options := mysql.MySQLConnectionOptions{
		Host:     "mysql",
		Port:     "3306",
		User:     "root",
		Password: "root",
		Database: "rate_limiter",
	}
	err := db.Connect(options)
	db.DB.Exec("DELETE FROM rate_limiter_models")
	return db, err
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

func TestNewRateLimiterRepositoryMySQL(t *testing.T) {
	client, err := connectDb()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	repo := NewRateLimiterRepositoryMySQL(client)

	if repo == nil {
		t.Errorf("Expected repo to be not nil")
	}
}

func TestGetByIPNotFound(t *testing.T) {
	client, err := connectDb()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	repo := NewRateLimiterRepositoryMySQL(client)

	limiter := repo.GetByIP("127.0.0.1")

	if limiter != nil {
		t.Errorf("Expected limiter to be nil")
	}
}

func TestGetByTokenNotFound(t *testing.T) {
	client, err := connectDb()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	repo := NewRateLimiterRepositoryMySQL(client)

	limiter := repo.GetByToken("token")

	if limiter != nil {
		t.Errorf("Expected limiter to be nil")
	}
}

func TestSave(t *testing.T) {
	client, err := connectDb()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	repo := NewRateLimiterRepositoryMySQL(client)

	limiter := createRateLimiter()
	err = repo.Save(limiter)
	if err != nil {
		t.Errorf("Expected error to be nil, got '%s'", err.Error())
	}
}

func TestGetByIP(t *testing.T) {
	client, err := connectDb()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	repo := NewRateLimiterRepositoryMySQL(client)

	limiter := createRateLimiter()
	err = repo.Save(limiter)
	if err != nil {
		t.Errorf("Expected error to be nil, got '%s'", err.Error())
	}

	limiterFound := repo.GetByIP(limiter.IP)
	if limiterFound == nil {
		t.Errorf("Expected limiterFound to be not nil")
	}
}

func TestGetByToken(t *testing.T) {
	client, err := connectDb()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	repo := NewRateLimiterRepositoryMySQL(client)

	limiter := createRateLimiter()
	limiter.Token = "token"

	err = repo.Save(limiter)
	if err != nil {
		t.Errorf("Expected error to be nil, got '%s'", err.Error())
	}

	limiterFound := repo.GetByToken(limiter.Token)
	if limiterFound == nil {
		t.Errorf("Expected limiterFound to be not nil")
	}
}

func TestSaveAndUpdate(t *testing.T) {
	client, err := connectDb()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	repo := NewRateLimiterRepositoryMySQL(client)

	limiter := createRateLimiter()
	err = repo.Save(limiter)
	if err != nil {
		t.Errorf("Expected error to be nil, got '%s'", err.Error())
	}

	limiterFound := repo.GetByIP(limiter.IP)
	if limiterFound == nil {
		t.Errorf("Expected limiterFound to be not nil")
	}

	limiterFound.Blocked = true
	err = repo.Save(limiterFound)
	if err != nil {
		t.Errorf("Expected error to be nil, got '%s'", err.Error())
	}

	limiterFound = repo.GetByIP(limiter.IP)
	if limiterFound == nil {
		t.Errorf("Expected limiterFound to be not nil")
	}
	if limiterFound.Blocked != true {
		t.Errorf("Expected limiterFound.Blocked to be true")
	}
}
