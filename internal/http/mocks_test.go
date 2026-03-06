package http_test

import (
	"context"
	"time"

	"github.com/highway-to-Golang/user-service/internal/domain"
	"github.com/stretchr/testify/mock"
)

// MockRepository - мок для репозитория
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, user domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockRepository) GetByID(ctx context.Context, id string) (domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return domain.User{}, args.Error(1)
	}
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockRepository) GetAll(ctx context.Context) ([]domain.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, id string, user domain.User) error {
	args := m.Called(ctx, id, user)
	return args.Error(0)
}

func (m *MockRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockEventSink - мок для event sink
type MockEventSink struct {
	mock.Mock
}

func (m *MockEventSink) Publish(ctx context.Context, event string) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEventSink) Close() error {
	args := m.Called()
	return args.Error(0)
}

// MockIdempotencyStorage - мок для idempotency storage
type MockIdempotencyStorage struct {
	mock.Mock
}

func (m *MockIdempotencyStorage) AcquireLock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	args := m.Called(ctx, key, ttl)
	return args.Bool(0), args.Error(1)
}

func (m *MockIdempotencyStorage) ReleaseLock(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockIdempotencyStorage) GetResult(ctx context.Context, key string) ([]byte, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockIdempotencyStorage) SaveResult(ctx context.Context, key string, data []byte, ttl time.Duration) error {
	args := m.Called(ctx, key, data, ttl)
	return args.Error(0)
}

func (m *MockIdempotencyStorage) Close() error {
	args := m.Called()
	return args.Error(0)
}
