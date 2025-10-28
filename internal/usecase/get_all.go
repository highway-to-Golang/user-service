package usecase

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/highway-to-Golang/user-service/internal/domain"
)

func (uc *UseCase) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	slog.Info("getting all users")

	users, err := uc.repository.GetAll(ctx)
	if err != nil {
		slog.Error("failed to get users", "error", err)
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	return users, nil
}
