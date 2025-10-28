package usecase

import (
	"context"

	"github.com/highway-to-Golang/user-service/internal/domain"
)

type Repository interface {
	Create(ctx context.Context, user domain.User) error
	GetByID(ctx context.Context, id string) (domain.User, error)
	GetAll(ctx context.Context) ([]domain.User, error)
	Update(ctx context.Context, id string, user domain.User) error
	Delete(ctx context.Context, id string) error
}

type UseCase struct {
	repository Repository
}

func New(repository Repository) *UseCase {
	return &UseCase{
		repository: repository,
	}
}
