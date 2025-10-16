package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		PG   PG
		HTTP HTTP
	}
	PG struct {
		Host     string `env:"PG_HOST" env-default:"localhost"`
		Port     string `env:"PG_PORT" env-default:"5432"`
		User     string `env:"PG_USER" env-default:"postgres"`
		Password string `env:"PG_PASSWORD" env-default:"postgres"`
		Database string `env:"PG_DATABASE" env-default:"user_service"`
	}
	HTTP struct {
		Host string `env:"HTTP_HOST" env-default:"localhost"`
		Port string `env:"HTTP_PORT" env-default:"8080"`
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
