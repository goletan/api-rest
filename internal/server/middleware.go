// /api-rest/internal/server/middleware.go
package server

import (
	"fmt"
	"net/http"
	"time"

	observability "github.com/goletan/observability/pkg"
	"go.uber.org/zap"
)

// Middleware for logging requests.
func loggingMiddleware(obs *observability.Observability) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := fmt.Sprintf("%d", time.Now().UnixNano()) // Simplified unique ID
			obs.Logger.Info("Incoming request",
				zap.String("method", r.Method),
				zap.String("url", r.URL.String()),
				zap.String("request_id", requestID),
				zap.String("client_ip", r.RemoteAddr),
			)
			next.ServeHTTP(w, r)
		})
	}
}

// Middleware for collecting metrics.
func metricsMiddleware(obs *observability.Observability) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			duration := time.Since(start)

			// Use the metrics library to collect metrics here
			obs.Logger.Info("Request processed",
				zap.String("method", r.Method),
				zap.String("url", r.URL.String()),
				zap.Duration("duration", duration),
			)
		})
	}
}
