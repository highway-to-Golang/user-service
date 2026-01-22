package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/highway-to-Golang/user-service/config"
	"github.com/highway-to-Golang/user-service/internal/domain"
	"github.com/highway-to-Golang/user-service/internal/http"
	"github.com/highway-to-Golang/user-service/internal/usecase"
	usecase_test "github.com/highway-to-Golang/user-service/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		idempotencyKey string
		setupMocks     func(*usecase_test.MockRepository, *usecase_test.MockEventSink, *usecase_test.MockIdempotencyStorage)
		expectedStatus int
	}{
		{
			name: "successful creation",
			requestBody: domain.CreateUserRequest{
				Name:  "John Doe",
				Email: "john@example.com",
				Role:  "user",
			},
			setupMocks: func(repo *usecase_test.MockRepository, sink *usecase_test.MockEventSink, storage *usecase_test.MockIdempotencyStorage) {
				repo.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "invalid request body",
			requestBody: map[string]interface{}{
				"invalid": "data",
			},
			setupMocks: func(repo *usecase_test.MockRepository, sink *usecase_test.MockEventSink, storage *usecase_test.MockIdempotencyStorage) {
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "idempotency key conflict",
			requestBody: domain.CreateUserRequest{
				Name:  "John Doe",
				Email: "john@example.com",
				Role:  "user",
			},
			idempotencyKey: "test-key",
			setupMocks: func(repo *usecase_test.MockRepository, sink *usecase_test.MockEventSink, storage *usecase_test.MockIdempotencyStorage) {
				storage.On("GetResult", mock.Anything, "test-key").Return(nil, nil).Twice()
				storage.On("AcquireLock", mock.Anything, "test-key", mock.Anything).Return(false, nil)
				storage.On("ReleaseLock", mock.Anything, "test-key").Return(nil)
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(usecase_test.MockRepository)
			sink := new(usecase_test.MockEventSink)
			storage := new(usecase_test.MockIdempotencyStorage)
			tt.setupMocks(repo, sink, storage)

			cfg := &config.Config{
				Redis: config.Redis{URL: "redis://localhost:6379"},
				NATS:  config.NATS{Enabled: false},
			}
			uc := usecase.New(repo, sink, storage, cfg)
			handler := http.NewUserHandler(uc)

			var body []byte
			var err error
			if tt.requestBody != nil {
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest("POST", "/api/users", bytes.NewBuffer(body))
			if tt.idempotencyKey != "" {
				req.Header.Set("Idempotency-Key", tt.idempotencyKey)
			}
			rr := httptest.NewRecorder()

			handler.CreateUser(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestGetUser(t *testing.T) {
	repo := new(usecase_test.MockRepository)
	sink := new(usecase_test.MockEventSink)
	storage := new(usecase_test.MockIdempotencyStorage)

	user := domain.User{
		ID:        "test-id",
		Name:      "John Doe",
		Email:     "john@example.com",
		Role:      "user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repo.On("GetByID", mock.Anything, "test-id").Return(user, nil)

	cfg := &config.Config{NATS: config.NATS{Enabled: false}}
	uc := usecase.New(repo, sink, storage, cfg)
	handler := http.NewUserHandler(uc)

	req := httptest.NewRequest("GET", "/api/users/test-id", nil)
	rr := httptest.NewRecorder()

	handler.GetUser(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response domain.User
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, response.ID)
	assert.Equal(t, user.Email, response.Email)

	repo.AssertExpectations(t)
}

func TestGetAllUsers(t *testing.T) {
	repo := new(usecase_test.MockRepository)
	sink := new(usecase_test.MockEventSink)
	storage := new(usecase_test.MockIdempotencyStorage)

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
	sink.On("Publish", mock.Anything, "get_all").Return(nil)

	cfg := &config.Config{NATS: config.NATS{Enabled: false}}
	uc := usecase.New(repo, sink, storage, cfg)
	handler := http.NewUserHandler(uc)

	req := httptest.NewRequest("GET", "/api/users", nil)
	rr := httptest.NewRecorder()

	handler.GetAllUsers(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response["users"])
	assert.Equal(t, float64(2), response["total"])

	repo.AssertExpectations(t)
	sink.AssertExpectations(t)
}

func TestUpdateUser(t *testing.T) {
	repo := new(usecase_test.MockRepository)
	sink := new(usecase_test.MockEventSink)
	storage := new(usecase_test.MockIdempotencyStorage)

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
	repo.On("Update", mock.Anything, "test-id", mock.Anything).Return(nil)
	repo.On("GetByID", mock.Anything, "test-id").Return(updatedUser, nil).Once()

	cfg := &config.Config{NATS: config.NATS{Enabled: false}}
	uc := usecase.New(repo, sink, storage, cfg)
	handler := http.NewUserHandler(uc)

	reqBody := domain.UpdateUserRequest{
		Name:  stringPtr("Updated Name"),
		Email: stringPtr("updated@example.com"),
		Role:  "admin",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/api/users/test-id", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	handler.UpdateUser(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response domain.User
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, updatedUser.ID, response.ID)
	assert.Equal(t, updatedUser.Name, response.Name)

	repo.AssertExpectations(t)
}

func TestDeleteUser(t *testing.T) {
	repo := new(usecase_test.MockRepository)
	sink := new(usecase_test.MockEventSink)
	storage := new(usecase_test.MockIdempotencyStorage)

	repo.On("Delete", mock.Anything, "test-id").Return(nil)

	cfg := &config.Config{NATS: config.NATS{Enabled: false}}
	uc := usecase.New(repo, sink, storage, cfg)
	handler := http.NewUserHandler(uc)

	req := httptest.NewRequest("DELETE", "/api/users/test-id", nil)
	rr := httptest.NewRecorder()

	handler.DeleteUser(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "User deleted successfully", response["message"])

	repo.AssertExpectations(t)
}

func stringPtr(s string) *string {
	return &s
}
