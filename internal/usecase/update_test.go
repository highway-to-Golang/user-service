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

func TestUpdateUser(t *testing.T) {
	tests := []struct {
		name          string
		userID        string
		req           domain.UpdateUserRequest
		setupMocks    func(*MockRepository, *MockEventSink)
		expectedError bool
	}{
		{
			name:   "successful update",
			userID: "test-id",
			req: domain.UpdateUserRequest{
				Name:  stringPtr("Updated Name"),
				Email: stringPtr("updated@example.com"),
				Role:  "admin",
			},
			setupMocks: func(repo *MockRepository, sink *MockEventSink) {
				existingUser := domain.User{
					ID:        "test-id",
					Name:      "Old Name",
					Email:     "old@example.com",
					Role:      "user",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				updatedUser := domain.User{
					ID:        "test-id",
					Name:      "Updated Name",
					Email:     "updated@example.com",
					Role:      "admin",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				repo.On("GetByID", mock.Anything, "test-id").Return(existingUser, nil).Once()
				repo.On("Update", mock.Anything, "test-id", mock.MatchedBy(func(u domain.User) bool {
					return u.Name == "Updated Name" && u.Email == "updated@example.com" && u.Role == "admin"
				})).Return(nil)
				repo.On("GetByID", mock.Anything, "test-id").Return(updatedUser, nil).Once()
			},
			expectedError: false,
		},
		{
			name:   "user not found",
			userID: "non-existent",
			req: domain.UpdateUserRequest{
				Name: stringPtr("Updated Name"),
			},
			setupMocks: func(repo *MockRepository, sink *MockEventSink) {
				repo.On("GetByID", mock.Anything, "non-existent").Return(domain.User{}, domain.ErrNotFound)
			},
			expectedError: true,
		},
		{
			name:   "partial update",
			userID: "test-id",
			req: domain.UpdateUserRequest{
				Role: "admin",
			},
			setupMocks: func(repo *MockRepository, sink *MockEventSink) {
				existingUser := domain.User{
					ID:        "test-id",
					Name:      "John Doe",
					Email:     "john@example.com",
					Role:      "user",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				updatedUser := domain.User{
					ID:        "test-id",
					Name:      "John Doe",
					Email:     "john@example.com",
					Role:      "admin",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				repo.On("GetByID", mock.Anything, "test-id").Return(existingUser, nil).Once()
				repo.On("Update", mock.Anything, "test-id", mock.MatchedBy(func(u domain.User) bool {
					return u.Role == "admin" && u.Name == "John Doe"
				})).Return(nil)
				repo.On("GetByID", mock.Anything, "test-id").Return(updatedUser, nil).Once()
			},
			expectedError: false,
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

			user, err := uc.UpdateUser(context.Background(), tt.userID, tt.req)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Equal(t, domain.ErrNotFound, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.userID, user.ID)
				if tt.req.Name != nil {
					assert.Equal(t, *tt.req.Name, user.Name)
				}
				if tt.req.Email != nil {
					assert.Equal(t, *tt.req.Email, user.Email)
				}
				if tt.req.Role != "" {
					assert.Equal(t, tt.req.Role, user.Role)
				}
			}

			repo.AssertExpectations(t)
		})
	}
}

func stringPtr(s string) *string {
	return &s
}
