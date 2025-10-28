package app

import (
	"context"
	"log/slog"
	"time"

	"github.com/highway-to-Golang/user-service/config"
	"github.com/highway-to-Golang/user-service/internal/database"
	"github.com/highway-to-Golang/user-service/internal/http"
	"github.com/highway-to-Golang/user-service/internal/nats"
	"github.com/highway-to-Golang/user-service/internal/repository"
	"github.com/highway-to-Golang/user-service/internal/usecase"
)

func Run(ctx context.Context, cfg *config.Config) error {
	db, err := database.NewDB(ctx, *cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)

	var eventSink usecase.EventSink
	if cfg.NATS.Enabled {
		es, err := nats.New(cfg.NATS.URL, cfg.NATS.SubjectPrefix)
		if err != nil {
			return err
		}
		defer es.Close()
		eventSink = es
	} else {
		eventSink = nats.NullEventSink{}
	}

	userUC := usecase.New(userRepo, eventSink)
	userHandler := http.NewUserHandler(userUC)
	server := http.NewServer(*cfg, userHandler)

	go func() {
		if err := server.Start(); err != nil {
			slog.Error("HTTP server error", "error", err)
		}
	}()

	<-ctx.Done()
	slog.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return err
	}

	return nil
}
