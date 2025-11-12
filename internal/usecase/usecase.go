package usecase

import (
	"context"
	"time"

	"github.com/highway-to-Golang/user-service/config"
	"github.com/highway-to-Golang/user-service/internal/domain"
	"github.com/highway-to-Golang/user-service/internal/nats"
	"github.com/highway-to-Golang/user-service/internal/redis"
)

type Repository interface {
	Create(ctx context.Context, user domain.User) error
	GetByID(ctx context.Context, id string) (domain.User, error)
	GetAll(ctx context.Context) ([]domain.User, error)
	Update(ctx context.Context, id string, user domain.User) error
	Delete(ctx context.Context, id string) error
}

type UseCase struct {
	repository         Repository
	eventSink          *nats.EventSink
	idempotencyStorage *redis.IdempotencyStorage
	cfg                *config.Config

	locksTTL       time.Duration
	idempotencyTTL time.Duration
}

func New(repository Repository, eventSink *nats.EventSink, idempotencyStorage *redis.IdempotencyStorage, cfg *config.Config) *UseCase {
	return &UseCase{
		repository:         repository,
		eventSink:          eventSink,
		idempotencyStorage: idempotencyStorage,
		cfg:                cfg,
		locksTTL:           30 * time.Second,
		idempotencyTTL:     24 * time.Hour,
	}
}
