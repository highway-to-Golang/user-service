package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/highway-to-Golang/user-service/internal/domain"
	"github.com/highway-to-Golang/user-service/internal/errors"
)

func (uc *UseCase) CreateUser(ctx context.Context, idempotencyKey string, req domain.CreateUserRequest) (domain.User, error) {
	if idempotencyKey != "" {
		cached, err := uc.idempotencyStorage.GetResult(ctx, idempotencyKey)
		if err != nil {
			slog.ErrorContext(ctx, "uc.idempotencyStorage.GetResult", "err", err, "key", idempotencyKey)
		}

		if cached != nil {
			var user domain.User
			if err := json.Unmarshal(cached, &user); err != nil {
				slog.ErrorContext(ctx, "json.Unmarshal", "err", err)
			} else {
				return user, nil
			}
		}

		lockAcquired, err := uc.idempotencyStorage.AcquireLock(ctx, idempotencyKey, uc.locksTTL)
		if err != nil {
			slog.ErrorContext(ctx, "uc.idempotencyStorage.AcquireLock", "err", err, "key", idempotencyKey)
		}

		if !lockAcquired {
			return domain.User{}, errors.ErrRequestAlreadyInProgress
		}

		defer func() {
			err = uc.idempotencyStorage.ReleaseLock(ctx, idempotencyKey)
			if err != nil {
				slog.ErrorContext(ctx, "uc.idempotencyStorage.ReleaseLock", "err", err, "key", idempotencyKey)
			}
		}()

		cached, err = uc.idempotencyStorage.GetResult(ctx, idempotencyKey)
		if err != nil {
			slog.ErrorContext(ctx, "uc.idempotencyStorage.GetResult (second check)", "err", err, "key", idempotencyKey)
		}

		if cached != nil {
			var user domain.User
			if err := json.Unmarshal(cached, &user); err != nil {
				slog.ErrorContext(ctx, "json.Unmarshal (second check)", "err", err)
			} else {
				return user, nil
			}
		}
	}

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

	if idempotencyKey != "" {
		data, err := json.Marshal(user)
		if err != nil {
			slog.ErrorContext(ctx, "json.Marshal", "err", err)
		} else {
			if err := uc.idempotencyStorage.SaveResult(ctx, idempotencyKey, data, uc.idempotencyTTL); err != nil {
				slog.ErrorContext(ctx, "uc.idempotencyStorage.SaveResult", "err", err, "key", idempotencyKey)
			}
		}
	}

	if err := uc.eventSink.Publish(ctx, "create"); err != nil {
		slog.Warn("failed to publish event", "error", err, "method", "create")
	}

	slog.Info("user created successfully", "user_id", user.ID, "email", user.Email)

	return user, nil
}
