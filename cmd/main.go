package main

import (
	"context"
	"net/http"
	"strings"

	"github.com/gabrielsc1998/go-rate-limiter/configs"
	common_errors "github.com/gabrielsc1998/go-rate-limiter/internal/common/errors"
	"github.com/gabrielsc1998/go-rate-limiter/internal/common/infra/db/mysql"
	"github.com/gabrielsc1998/go-rate-limiter/internal/common/infra/db/redis"
	"github.com/gabrielsc1998/go-rate-limiter/internal/common/infra/webserver"
	rate_limiter_repo "github.com/gabrielsc1998/go-rate-limiter/internal/rate-limiter/infra/db"
	rate_limiter_repo_mysql "github.com/gabrielsc1998/go-rate-limiter/internal/rate-limiter/infra/db/mysql"
	rate_limiter_middleware "github.com/gabrielsc1998/go-rate-limiter/internal/rate-limiter/infra/middleware"
)

func main() {
	config, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	rateLimiterRepository, err := setupRateLimiterRepository(config)
	if err != nil {
		panic(err)
	}

	rateLimiterMiddleware := rate_limiter_middleware.NewRateLimiterMiddleware(config, rateLimiterRepository)

	if err != nil {
		panic(err)
	}

	webserver := webserver.NewWebServer(config.WebServerPort)

	webserver.AddMiddleware(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := ReadUserIP(r)
			token := r.Header.Get("API_KEY")
			err := rateLimiterMiddleware.Handle(ip, token)
			if err != nil {
				if common_errors.Is(err, common_errors.ErrTooManyRequests) {
					w.WriteHeader(http.StatusTooManyRequests)
					w.Write([]byte("you have reached the maximum number of requests or actions allowed within a certain time frame"))
					return
				}
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			h.ServeHTTP(w, r)
		})
	})

	webserver.Get("/teste", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})
	webserver.Start()
}

func ReadUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress != "" {
		return IPAddress
	}

	IPAddress = r.Header.Get("X-Forwarded-For")
	if IPAddress != "" {
		return IPAddress
	}

	IPAddress = r.RemoteAddr
	if IPAddress != "" {
		splittedIp := strings.Split(IPAddress, ":")
		return splittedIp[0]
	}

	return IPAddress
}

func setupRateLimiterRepository(config *configs.Conf) (rate_limiter_repo.RateLimiterRepository, error) {
	// ================== Redis ================== //

	redisClient := redis.NewRedisClient(redis.RedisClientConfig{
		Addr:     config.RedisAddr,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})
	// defer redisClient.Close()
	redisClient.ClearAll(context.Background())

	// ================== MySQL ================== //

	mysqlDB := mysql.NewMySQLDBConnection()
	mysqlDB.Connect(mysql.MySQLConnectionOptions{
		Host:     config.MySQLHost,
		Port:     config.MySQLPort,
		User:     config.MySQLUser,
		Password: config.MySQLPassword,
		Database: config.MySQLDatabase,
	})
	err := mysqlDB.DB.AutoMigrate(
		&rate_limiter_repo_mysql.RateLimiterModel{},
	)
	if err != nil {
		panic(err)
	}
	// defer mysqlDB.Close()

	// ================== Rate Limiter Repository ================== //

	switch config.PersistenceMechanism {
	case "redis":
		return rate_limiter_repo.NewRateLimiterRepository(rate_limiter_repo.Config{
			Repo:   "redis",
			Inject: redisClient,
		})
	case "mysql":
		return rate_limiter_repo.NewRateLimiterRepository(rate_limiter_repo.Config{
			Repo:   "mysql",
			Inject: mysqlDB,
		})
	default:
		panic("invalid persistence mechanism")
	}
}
