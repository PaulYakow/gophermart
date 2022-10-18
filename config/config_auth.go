package config

import (
	"github.com/spf13/viper"
	"time"
)

type AuthConfig struct {
	TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

func LoadAuthConfig() (config AuthConfig, err error) {
	viper.SetConfigFile("./config/.env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
