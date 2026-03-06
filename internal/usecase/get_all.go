package usecase

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/highway-to-Golang/user-service/internal/domain"
	"github.com/highway-to-Golang/user-service/internal/monitoring"
)

func (uc *UseCase) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	slog.Info("getting all users")

	users, err := uc.repository.GetAll(ctx)
	if err != nil {
		slog.Error("failed to get users", "error", err)
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	if uc.cfg.NATS.Enabled && uc.eventSink != nil {
		if err := uc.eventSink.Publish(ctx, "get_all"); err != nil {
			slog.Warn("failed to publish event", "error", err, "method", "get_all")
		}
	}

	monitoring.UsersRetrieved.Add(float64(len(users)))

	return users, nil
}
