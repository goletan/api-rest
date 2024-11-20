// /api-rest/internal/handlers/handlers.go

package handlers

import (
	"net/http"
)

// HealthHandler handles the health check endpoint.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// StatusHandler provides status information of the API.
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("REST API is running smoothly"))
}
