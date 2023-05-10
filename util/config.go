package util

import (
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	DBDriver             string        `mapstructure:"DB_DRIVER"`
	DBSource             string        `mapstructure:"DB_SOURCE"`
	HttpServerAddress    string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	GrpcServerAddress    string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.SetDefault("DB_DRIVER", "postgres")
	viper.SetDefault("DB_SOURCE", "postgresql://root:ab@localhost:5432/simple_bank?sslmode=disable")
	viper.SetDefault("HTTP_SERVER_ADDRESS", "0.0.0.0:8080")
	viper.SetDefault("GRPC_SERVER_ADDRESS", "0.0.0.0:9090")
	viper.SetDefault("TOKEN_SYMMETRIC_KEY", "12345678901234567890123456789012")
	viper.SetDefault("ACCESS_TOKEN_DURATION", 15*time.Minute)
	viper.SetDefault("REFRESH_TOKEN_DURATION", 24*time.Hour)

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
