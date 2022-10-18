package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type AuthConfig struct {
	TokenSymmetricKey   string        `env:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `env:"ACCESS_TOKEN_DURATION"`
}

func NewAuthConfig() (*AuthConfig, error) {
	cfg := &AuthConfig{}

	err := cleanenv.ReadConfig("./config/auth.env", cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
