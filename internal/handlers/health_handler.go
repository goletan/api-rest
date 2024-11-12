// /api-rest/internal/handlers/health_handler.go
package handlers

import (
	"net/http"
)

// HealthHandler handles the health check endpoint.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
