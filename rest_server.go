// /api-rest/rest_server.go
package rest

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/goletan/config"
	"github.com/goletan/observability/logger"
	"github.com/goletan/security/mtls"
	"github.com/goletan/services"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// RESTServer is an enhanced HTTP server that implements the Service interface.
type RESTServer struct {
	server *http.Server
	name   string
}

type RestConfig struct {
	Address      string        `mapstructure:"address"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
	EnableTLS    bool          `mapstructure:"enable_tls"`
	CertFilePath string        `mapstructure:"cert_file_path"`
	KeyFilePath  string        `mapstructure:"key_file_path"`
}

var cfg *RestConfig

// NewRESTServer creates a new instance of the RESTServer.
func NewRESTServer() services.Service {

	cfg = &RestConfig{}
	err := config.LoadConfig("Rest", cfg, nil)
	if err != nil {
		fmt.Printf("Warning: failed to load REST configuration, using defaults: %v\n", err)
	}

	r := mux.NewRouter()

	// Define middlewares for observability
	r.Use(loggingMiddleware)
	r.Use(metricsMiddleware)

	// Define your REST endpoints here, e.g.:
	r.HandleFunc("/health", healthHandler).Methods("GET")

	// Initialize server
	server := &http.Server{
		Addr:              cfg.Address, // Load from config
		Handler:           r,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Configure TLS if enabled
	if cfg.EnableTLS {
		tlsConfig, tlsErr := mtls.ConfigureMTLS() // Use the security package's TLS configuration
		if tlsErr != nil {
			fmt.Printf("Warning: failed to configure TLS, continuing without TLS: %v\n", tlsErr)
		} else {
			server.TLSConfig = tlsConfig
		}
	}

	return &RESTServer{
		server: server,
		name:   "REST Server",
	}
}

// Name returns the service name.
func (s *RESTServer) Name() string {
	return s.name
}

// Initialize performs any initialization tasks needed by the service.
func (s *RESTServer) Initialize() error {
	logger.Info("Initializing REST server", zap.String("service", s.name))
	return nil
}

// Start starts the REST server.
func (s *RESTServer) Start() error {
	go func() {
		logger.Info("Starting REST server", zap.String("address", s.server.Addr))
		var err error
		if cfg.EnableTLS {
			err = s.server.ListenAndServeTLS(cfg.CertFilePath, cfg.KeyFilePath) // Use configured certificate and key files
		} else {
			err = s.server.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			logger.Error("Failed to start REST server", zap.Error(err))
		}
	}()
	return nil
}

// Stop gracefully stops the REST server.
func (s *RESTServer) Stop() error {
	logger.Info("Stopping REST server", zap.String("service", s.name))
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}

// Middleware for logging requests.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := fmt.Sprintf("%d", time.Now().UnixNano()) // Simplified unique ID
		logger.Info("Incoming request",
			zap.String("method", r.Method),
			zap.String("url", r.URL.String()),
			zap.String("request_id", requestID),
			zap.String("client_ip", r.RemoteAddr),
		)
		next.ServeHTTP(w, r)
	})
}

// Middleware for collecting metrics.
func metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)

		// You would use the metrics library to collect metrics here
		logger.Info("Request processed",
			zap.String("method", r.Method),
			zap.String("url", r.URL.String()),
			zap.Duration("duration", duration),
		)
	})
}

// Health handler for health checks.
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
