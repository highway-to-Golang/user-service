package usecase

import (
	"context"
	"time"

	"github.com/highway-to-Golang/user-service/internal/domain"
)

type Repository interface {
	Create(ctx context.Context, user domain.User) error
	GetByID(ctx context.Context, id string) (domain.User, error)
	GetAll(ctx context.Context) ([]domain.User, error)
	Update(ctx context.Context, id string, user domain.User) error
	Delete(ctx context.Context, id string) error
}

type EventSink interface {
	Publish(ctx context.Context, method string) error
}

type IdempotencyStorage interface {
	GetResult(ctx context.Context, key string) ([]byte, error)
	SaveResult(ctx context.Context, key string, value []byte, ttl time.Duration) error
	AcquireLock(ctx context.Context, key string, ttl time.Duration) (bool, error)
	ReleaseLock(ctx context.Context, key string) error
}

type UseCase struct {
	repository         Repository
	eventSink          EventSink
	idempotencyStorage IdempotencyStorage

	locksTTL       time.Duration
	idempotencyTTL time.Duration
}

func New(repository Repository, eventSink EventSink, idempotencyStorage IdempotencyStorage) *UseCase {
	return &UseCase{
		repository:         repository,
		eventSink:          eventSink,
		idempotencyStorage: idempotencyStorage,
		locksTTL:           30 * time.Second,
		idempotencyTTL:     24 * time.Hour,
	}
}
