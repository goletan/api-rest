// /api-rest/internal/handlers/metrics.go
package handlers

import (
	"net/http"

	observability "github.com/goletan/observability/pkg"
	"github.com/prometheus/client_golang/prometheus"
)

type RestMetrics struct{}

// HTTP Metrics: Track HTTP requests and errors.
var (
	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "goletan",
			Subsystem: "rest",
			Name:      "http_request_duration_seconds",
			Help:      "Duration of HTTP requests in seconds.",
		},
		[]string{"method", "endpoint", "status"},
	)
)

// Security Tool: Scrub sensitive data
var (
	scrubber = observability.NewScrubber()
)

func InitMetrics() {
	obs, err := observability.NewObserver()
	if err != nil {
		obs.Logger.Error("cannot observe REST API")
		return
	}

	obs.Metrics.Register(&RestMetrics{})
}

func (em *RestMetrics) Register() error {
	if err := prometheus.Register(RequestDuration); err != nil {
		return err
	}

	return nil
}

// ObserveRequestDuration records the duration of HTTP requests.
func ObserveRequestDuration(method, endpoint string, status int, duration float64) {
	scrubbedMethod := scrubber.Scrub(method)
	scrubbedEndpoint := scrubber.Scrub(endpoint)
	RequestDuration.WithLabelValues(scrubbedMethod, scrubbedEndpoint, http.StatusText(status)).Observe(duration)
}
