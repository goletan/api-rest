// /api-rest/internal/server/middleware.go
package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/goletan/api-rest/internal/handlers"
	observability "github.com/goletan/observability/pkg"
	"go.uber.org/zap"
)

func observabilityMiddleware(obs *observability.Observability) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			requestID := fmt.Sprintf("%d", time.Now().UnixNano()) // Simplified unique ID

			obs.Logger.Info("Incoming request",
				zap.String("method", r.Method),
				zap.String("url", r.URL.String()),
				zap.String("request_id", requestID),
				zap.String("client_ip", r.RemoteAddr),
			)

			next.ServeHTTP(w, r)

			duration := time.Since(start)
			handlers.ObserveRequestDuration(r.Method, r.URL.Path, http.StatusOK, duration.Seconds())

			obs.Logger.Info("Request processed",
				zap.String("method", r.Method),
				zap.String("url", r.URL.String()),
				zap.Duration("duration", duration),
			)
		})
	}
}
