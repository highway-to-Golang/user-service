package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/highway-to-Golang/user-service/config"
	"github.com/highway-to-Golang/user-service/internal/domain"
	"github.com/highway-to-Golang/user-service/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetUser(t *testing.T) {
	tests := []struct {
		name          string
		userID        string
		setupMocks    func(*MockRepository, *MockEventSink)
		expectedError bool
	}{
		{
			name:   "successful get",
			userID: "test-id",
			setupMocks: func(repo *MockRepository, sink *MockEventSink) {
				user := domain.User{
					ID:        "test-id",
					Name:      "John Doe",
					Email:     "john@example.com",
					Role:      "user",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				repo.On("GetByID", mock.Anything, "test-id").Return(user, nil)
			},
			expectedError: false,
		},
		{
			name:   "user not found",
			userID: "non-existent",
			setupMocks: func(repo *MockRepository, sink *MockEventSink) {
				repo.On("GetByID", mock.Anything, "non-existent").Return(domain.User{}, domain.ErrNotFound)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockRepository)

			tt.setupMocks(repo, new(MockEventSink))

			cfg := &config.Config{
				NATS: config.NATS{Enabled: false},
			}

			uc := usecase.New(repo, nil, nil, cfg)

			user, err := uc.GetUser(context.Background(), tt.userID)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Equal(t, domain.ErrNotFound, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.userID, user.ID)
			}

			repo.AssertExpectations(t)
		})
	}
}

func TestGetAllUsers(t *testing.T) {
	repo := new(MockRepository)

	users := []domain.User{
		{
			ID:        "id1",
			Name:      "User 1",
			Email:     "user1@example.com",
			Role:      "user",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        "id2",
			Name:      "User 2",
			Email:     "user2@example.com",
			Role:      "admin",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	repo.On("GetAll", mock.Anything).Return(users, nil)

	cfg := &config.Config{
		NATS: config.NATS{Enabled: false},
	}

	uc := usecase.New(repo, nil, nil, cfg)

	result, err := uc.GetAllUsers(context.Background())

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, users[0].ID, result[0].ID)
	assert.Equal(t, users[1].ID, result[1].ID)

	repo.AssertExpectations(t)
}
