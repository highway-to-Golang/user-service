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

func TestDeleteUser(t *testing.T) {
	tests := []struct {
		name          string
		userID        string
		setupMocks    func(*MockRepository, *MockEventSink)
		expectedError bool
	}{
		{
			name:   "successful deletion",
			userID: "test-id",
			setupMocks: func(repo *MockRepository, sink *MockEventSink) {
				repo.On("Delete", mock.Anything, "test-id").Return(nil)
			},
			expectedError: false,
		},
		{
			name:   "user not found",
			userID: "non-existent",
			setupMocks: func(repo *MockRepository, sink *MockEventSink) {
				repo.On("Delete", mock.Anything, "non-existent").Return(domain.ErrNotFound)
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

			err := uc.DeleteUser(context.Background(), tt.userID)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Equal(t, domain.ErrNotFound, err)
			} else {
				assert.NoError(t, err)
			}

			repo.AssertExpectations(t)
		})
	}
}
