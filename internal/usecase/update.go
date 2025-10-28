package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/highway-to-Golang/user-service/internal/domain"
)

func (uc *UseCase) UpdateUser(ctx context.Context, id string, req domain.UpdateUserRequest) (domain.User, error) {
	existingUser, err := uc.repository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.User{}, domain.ErrNotFound
		}
		slog.Error("failed to get user for update", "error", err, "user_id", id)
		return domain.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	if req.Email != nil {
		existingUser.Email = *req.Email
	}
	if req.Name != nil {
		existingUser.Name = *req.Name
	}
	if req.Role != "" {
		existingUser.Role = req.Role
	}

	slog.Info("updating user", "id", id, "email", existingUser.Email, "name", existingUser.Name)

	if err := uc.repository.Update(ctx, id, existingUser); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.User{}, domain.ErrNotFound
		}
		slog.Error("failed to update user", "error", err, "user_id", id)
		return domain.User{}, fmt.Errorf("failed to update user: %w", err)
	}

	updatedUser, err := uc.repository.GetByID(ctx, id)
	if err != nil {
		slog.Error("failed to get updated user", "error", err, "user_id", id)
		return domain.User{}, fmt.Errorf("failed to get updated user: %w", err)
	}

	if err := uc.eventSink.Publish(ctx, "update"); err != nil {
		slog.Warn("failed to publish event", "error", err, "method", "update")
	}

	return updatedUser, nil
}
