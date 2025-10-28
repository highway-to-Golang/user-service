package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/highway-to-Golang/user-service/internal/domain"
)

func (uc *UseCase) DeleteUser(ctx context.Context, id string) error {
	slog.Info("deleting user", "id", id)

	if err := uc.repository.Delete(ctx, id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			slog.Warn("user not found for deletion", "user_id", id)
			return domain.ErrNotFound
		}
		slog.Error("failed to delete user", "error", err, "user_id", id)
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
