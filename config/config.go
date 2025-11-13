package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		PG    PG
		HTTP  HTTP
		GRPC  GRPC
		NATS  NATS
		Redis Redis
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
	GRPC struct {
		Host string `env:"GRPC_HOST" env-default:"localhost"`
		Port string `env:"GRPC_PORT" env-default:"8081"`
	}
	NATS struct {
		Enabled       bool   `env:"NATS_ENABLED" env-default:"false"`
		URL           string `env:"NATS_URL" env-default:"nats://localhost:4222"`
		SubjectPrefix string `env:"NATS_SUBJECT_PREFIX" env-default:"user.events"`
	}
	Redis struct {
		URL string `env:"REDIS_URL" env-default:"redis://localhost:6379"`
	}
	Client struct {
		Protocol string `env:"CLIENT_PROTOCOL" env-default:"http"`
		HTTPURL  string `env:"CLIENT_HTTP_URL" env-default:"http://localhost:8080"`
		GRPCAddr string `env:"CLIENT_GRPC_ADDR" env-default:"localhost:8081"`
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

func NewClientConfig() (*Client, error) {
	var clientConfig Client
	_ = cleanenv.ReadConfig(".env", &clientConfig)

	if err := cleanenv.ReadEnv(&clientConfig); err != nil {
		return nil, err
	}

	return &clientConfig, nil
}
