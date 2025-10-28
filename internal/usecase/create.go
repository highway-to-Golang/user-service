package usecase

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/highway-to-Golang/user-service/internal/domain"
)

func (uc *UseCase) CreateUser(ctx context.Context, req domain.CreateUserRequest) (domain.User, error) {
	if req.Email == "" || req.Name == "" {
		return domain.User{}, fmt.Errorf("email and name are required")
	}

	if req.Role == "" {
		req.Role = "user"
	}

	user, err := domain.NewUser(req.Name, req.Email, req.Role)
	if err != nil {
		slog.Error("failed to create user", "error", err)
		return domain.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	if err := uc.repository.Create(ctx, user); err != nil {
		slog.Error("failed to save user", "error", err)
		return domain.User{}, fmt.Errorf("failed to save user: %w", err)
	}

	if err := uc.eventSink.Publish(ctx, "create"); err != nil {
		slog.Warn("failed to publish event", "error", err, "method", "create")
	}

	slog.Info("user created successfully", "user_id", user.ID, "email", user.Email)

	return user, nil
}
