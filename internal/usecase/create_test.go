package usecase_test

import (
	"context"
	"testing"

	"github.com/highway-to-Golang/user-service/config"
	"github.com/highway-to-Golang/user-service/internal/domain"
	"github.com/highway-to-Golang/user-service/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name           string
		req            domain.CreateUserRequest
		idempotencyKey string
		setupMocks     func(*MockRepository, *MockEventSink, *MockIdempotencyStorage)
		expectedError  error
	}{
		{
			name: "successful creation",
			req: domain.CreateUserRequest{
				Name:  "John Doe",
				Email: "john@example.com",
				Role:  "user",
			},
			idempotencyKey: "",
			setupMocks: func(repo *MockRepository, sink *MockEventSink, storage *MockIdempotencyStorage) {
				repo.On("Create", mock.Anything, mock.MatchedBy(func(u domain.User) bool {
					return u.Name == "John Doe" && u.Email == "john@example.com" && u.Role == "user"
				})).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "missing email",
			req: domain.CreateUserRequest{
				Name:  "John Doe",
				Email: "",
				Role:  "user",
			},
			setupMocks: func(repo *MockRepository, sink *MockEventSink, storage *MockIdempotencyStorage) {
				// No repository calls expected
			},
			expectedError: assert.AnError,
		},
		{
			name: "default role",
			req: domain.CreateUserRequest{
				Name:  "Jane Doe",
				Email: "jane@example.com",
				Role:  "",
			},
			setupMocks: func(repo *MockRepository, sink *MockEventSink, storage *MockIdempotencyStorage) {
				repo.On("Create", mock.Anything, mock.MatchedBy(func(u domain.User) bool {
					return u.Role == "user"
				})).Return(nil)
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockRepository)
			storage := new(MockIdempotencyStorage)

			tt.setupMocks(repo, new(MockEventSink), storage)

			cfg := &config.Config{
				Redis: config.Redis{URL: ""}, // Отключаем Redis для упрощения тестов
				NATS:  config.NATS{Enabled: false},
			}

			uc := usecase.New(repo, nil, nil, cfg)

			user, err := uc.CreateUser(context.Background(), tt.idempotencyKey, tt.req)

			if tt.expectedError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, user.ID)
				if tt.req.Email != "" {
					assert.Equal(t, tt.req.Email, user.Email)
				}
				if tt.req.Name != "" {
					assert.Equal(t, tt.req.Name, user.Name)
				}
			}

			repo.AssertExpectations(t)
		})
	}
}

// TestCreateUser_WithIdempotencyLock пропущен, так как требует реального storage mock
// Для полного тестирования idempotency нужны интеграционные тесты
