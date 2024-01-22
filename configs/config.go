package configs

import "github.com/spf13/viper"

type Conf struct {
	RedisAddr                 string   `mapstructure:"REDIS_ADDR"`
	RedisPassword             string   `mapstructure:"REDIS_PASSWORD"`
	RedisDB                   int      `mapstructure:"REDIS_DB"`
	MySQLHost                 string   `mapstructure:"MYSQL_HOST"`
	MySQLPort                 string   `mapstructure:"MYSQL_PORT"`
	MySQLUser                 string   `mapstructure:"MYSQL_USER"`
	MySQLPassword             string   `mapstructure:"MYSQL_PASSWORD"`
	MySQLDatabase             string   `mapstructure:"MYSQL_DATABASE"`
	PersistenceMechanism      string   `mapstructure:"PERSISTENCE_MECHANISM"`
	WebServerPort             string   `mapstructure:"WEB_SERVER_PORT"`
	RateLimiterMaxReqsIP      int      `mapstructure:"RATE_LIMITER_MAX_REQUESTS_PER_SECOND_FOR_IP"`
	RateLimiterBlockTimeIP    float64  `mapstructure:"RATE_LIMITER_BLOCK_TIME_FOR_IP"`
	RateLimiterBlockTimeToken float64  `mapstructure:"RATE_LIMITER_BLOCK_TIME_FOR_TOKEN"`
	RateLimiterTokens         []string `mapstructure:"RATE_LIMITER_TOKENS"`
}

func LoadConfig(path string) (*Conf, error) {
	var cfg *Conf
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg, err
}
