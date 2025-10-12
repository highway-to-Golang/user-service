package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		PG PG
	}
	PG struct {
		Host     string `env:"PG_HOST" env-default:"localhost"`
		Port     string `env:"PG_PORT" env-default:"5432"`
		User     string `env:"PG_USER" env-default:"postgres"`
		Password string `env:"PG_PASSWORD" env-default:"postgres"`
	}
)

func NewConfig() (*Config, error) {
	var config Config
	_ = cleanenv.ReadConfig(".env", &config)

	if err := cleanenv.ReadEnv(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
