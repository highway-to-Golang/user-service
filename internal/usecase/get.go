package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/highway-to-Golang/user-service/internal/domain"
)

func (uc *UseCase) GetUser(ctx context.Context, id string) (domain.User, error) {
	slog.Info("getting user", "id", id)

	user, err := uc.repository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			slog.Warn("user not found", "user_id", id)
			return domain.User{}, err
		}
		slog.Error("failed to get user", "error", err, "user_id", id)
		return domain.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}
