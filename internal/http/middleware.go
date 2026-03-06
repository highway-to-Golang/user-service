package http

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/highway-to-Golang/user-service/internal/monitoring"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		slog.Info("incoming request",
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
		)

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)
		slog.Info("request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"status", wrapped.statusCode,
			"duration_ms", duration.Milliseconds(),
		)
	})
}

// PrometheusMiddleware добавляет метрики Prometheus для HTTP запросов
func PrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)
		status := strconv.Itoa(wrapped.statusCode)
		path := r.URL.Path

		monitoring.HTTPRequestsTotal.WithLabelValues(r.Method, path, status).Inc()
		monitoring.HTTPRequestDuration.WithLabelValues(r.Method, path, status).Observe(duration.Seconds())
	})
}

// TracingMiddleware добавляет трассировку для HTTP запросов
func TracingMiddleware(next http.Handler) http.Handler {
	return otelhttp.NewHandler(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			span := trace.SpanFromContext(ctx)
			if span.IsRecording() {
				span.SetAttributes(
					attribute.String("http.method", r.Method),
					attribute.String("http.path", r.URL.Path),
					attribute.String("http.remote_addr", r.RemoteAddr),
				)
			}
			next.ServeHTTP(w, r)
		}),
		"http.request",
		otelhttp.WithTracerProvider(otel.GetTracerProvider()),
	)
}
