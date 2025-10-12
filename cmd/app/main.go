package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/highway-to-Golang/user-service/config"
	"github.com/highway-to-Golang/user-service/internal/app"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.NewConfig()
	if err != nil {
		slog.Error("Error loading config:", "error", err.Error())
		os.Exit(1)
	}

	slog.Info("cfg", "cfg", cfg)

	err = app.Run(ctx, cfg)
	if err != nil {
		slog.Error("Error running app", "error", err.Error())
		os.Exit(1)
	}
}
