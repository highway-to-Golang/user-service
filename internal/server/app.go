package server

import (
	"context"
	"log/slog"
	"time"

	"github.com/highway-to-Golang/user-service/config"
	"github.com/highway-to-Golang/user-service/internal/database"
	"github.com/highway-to-Golang/user-service/internal/grpc"
	"github.com/highway-to-Golang/user-service/internal/http"
	"github.com/highway-to-Golang/user-service/internal/monitoring"
	"github.com/highway-to-Golang/user-service/internal/nats"
	"github.com/highway-to-Golang/user-service/internal/redis"
	"github.com/highway-to-Golang/user-service/internal/repository"
	"github.com/highway-to-Golang/user-service/internal/usecase"
)

func Run(ctx context.Context, cfg *config.Config) error {
	// Инициализация трассировки
	if cfg.Tracing.Enabled {
		if err := monitoring.InitTracer(cfg.Tracing.Service, cfg.Tracing.OTLP.Endpoint); err != nil {
			slog.Warn("Failed to initialize tracer", "error", err)
		} else {
			defer func() {
				if err := monitoring.ShutdownTracer(context.Background()); err != nil {
					slog.Error("Failed to shutdown tracer", "error", err)
				}
			}()
		}
	}

	db, err := database.NewDB(ctx, *cfg)
	if err != nil {
		return err
	}
	defer db.Pool.Close()

	userRepo := repository.NewUserRepository(db)

	var eventSink *nats.EventSink
	if cfg.NATS.Enabled {
		es, err := nats.New(cfg.NATS.URL, cfg.NATS.SubjectPrefix)
		if err != nil {
			return err
		}
		defer es.Close()
		eventSink = es
	}

	var idempotencyStorage *redis.IdempotencyStorage
	if cfg.Redis.URL != "" {
		is, err := redis.NewIdempotencyStorage(cfg.Redis.URL)
		if err != nil {
			return err
		}
		defer func() { _ = is.Close() }()
		idempotencyStorage = is
	}

	userUC := usecase.New(userRepo, eventSink, idempotencyStorage, cfg)

	httpUserHandler := http.NewUserHandler(userUC)
	httpServer := http.NewServer(*cfg, httpUserHandler)

	grpcUserHandler := grpc.NewUserHandler(userUC)
	grpcServer := grpc.NewServer(*cfg, grpcUserHandler)

	go func() {
		if err := httpServer.Start(); err != nil {
			slog.Error("HTTP server error", "error", err)
		}
	}()

	go func() {
		if err := grpcServer.Start(); err != nil {
			slog.Error("gRPC server error", "error", err)
		}
	}()

	<-ctx.Done()
	slog.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var shutdownErr error
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		shutdownErr = err
		slog.Error("HTTP server shutdown error", "error", err)
	}

	if err := grpcServer.Shutdown(shutdownCtx); err != nil {
		shutdownErr = err
		slog.Error("gRPC server shutdown error", "error", err)
	}

	return shutdownErr
}
