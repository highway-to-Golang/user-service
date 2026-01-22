package http

import (
	"net/http"

	"github.com/highway-to-Golang/user-service/config"
	"github.com/highway-to-Golang/user-service/internal/monitoring"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter(userHandler *UserHandler, cfg config.Config) http.Handler {
	mux := http.NewServeMux()

	// API endpoints
	mux.HandleFunc("GET /api/users", userHandler.GetAllUsers)
	mux.HandleFunc("POST /api/users", userHandler.CreateUser)
	mux.HandleFunc("GET /api/users/{id}", userHandler.GetUser)
	mux.HandleFunc("PUT /api/users/{id}", userHandler.UpdateUser)
	mux.HandleFunc("DELETE /api/users/{id}", userHandler.DeleteUser)

	// Metrics endpoint
	if cfg.Monitoring.Enabled {
		mux.Handle("/metrics", promhttp.Handler())
	}

	// pprof endpoints (только если включен)
	if cfg.Pprof.Enabled {
		monitoring.RegisterPprofHandlers(mux)
	}

	return mux
}
