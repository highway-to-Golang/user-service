package client

import (
	"context"

	"github.com/highway-to-Golang/user-service/internal/domain"
)

type Client interface {
	CreateUser(ctx context.Context, req domain.CreateUserRequest, idempotencyKey string) (domain.User, error)
	GetUser(ctx context.Context, id string) (domain.User, error)
	GetAllUsers(ctx context.Context) ([]domain.User, error)
	UpdateUser(ctx context.Context, id string, req domain.UpdateUserRequest) (domain.User, error)
	DeleteUser(ctx context.Context, id string) error
	Close() error
}
