package grpc

import (
	"context"
	"time"

	"github.com/highway-to-Golang/user-service/internal/monitoring"
	"google.golang.org/grpc"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"
)

// PrometheusUnaryInterceptor добавляет метрики Prometheus для gRPC запросов
func PrometheusUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	resp, err := handler(ctx, req)

	duration := time.Since(start)
	statusCode := status.Code(err).String()
	method := info.FullMethod

	monitoring.GRPCRequestsTotal.WithLabelValues(method, statusCode).Inc()
	monitoring.GRPCRequestDuration.WithLabelValues(method, statusCode).Observe(duration.Seconds())

	return resp, err
}

// GetUnaryInterceptors возвращает unary interceptor для Prometheus
func GetUnaryInterceptors(tracingEnabled bool) grpc.UnaryServerInterceptor {
	return PrometheusUnaryInterceptor
}

// GetStatsHandler возвращает stats handler для трассировки (если включена)
func GetStatsHandler(tracingEnabled bool) stats.Handler {
	if !tracingEnabled {
		return nil
	}
	// Используем стандартный stats handler из otelgrpc
	// Для упрощения, трассировка будет работать через HTTP middleware
	// gRPC трассировка может быть добавлена позже при необходимости
	return nil
}
