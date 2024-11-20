// /api-rest/internal/server/server.go
package server

import (
	"context"
	"net/http"
	"time"

	"github.com/goletan/api-rest/internal/handlers"
	config "github.com/goletan/config/pkg"
	observability "github.com/goletan/observability/pkg"
	security "github.com/goletan/security/pkg"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// RESTServer is an enhanced HTTP server that implements the Service interface.
type RESTServer struct {
	server         *http.Server
	name           string
	observability  *observability.Observability
	securityModule *security.Security
}

// RestConfig represents the REST server configuration.
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
func NewRESTServer(obs *observability.Observability) *RESTServer {
	return &RESTServer{
		name:          "REST Server",
		observability: obs,
	}
}

// Name returns the service name.
func (s *RESTServer) Name() string {
	return s.name
}

// Initialize performs any initialization tasks needed by the service.
func (s *RESTServer) Initialize() error {
	s.observability.Logger.Info("Initializing REST server", zap.String("service", s.name))

	// Load the REST configuration
	cfg := &RestConfig{
		Address:      ":8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	err := config.LoadConfig("Rest", cfg, s.observability.Logger)
	if err != nil {
		s.observability.Logger.Warn("Failed to load REST configuration, using defaults", zap.Error(err))
	}

	// Set up the router
	router := setupRouter(s.observability)
	s.server = &http.Server{
		Addr:              cfg.Address,
		Handler:           router,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Load and configure security module (mTLS)
	s.securityModule = setupSecurityModule(s.observability)

	return nil
}

// Start starts the REST server.
func (s *RESTServer) Start() error {
	go func() {
		s.observability.Logger.Info("Starting REST server", zap.String("address", s.server.Addr))
		var err error
		if cfg.EnableTLS {
			err = s.server.ListenAndServeTLS(cfg.CertFilePath, cfg.KeyFilePath)
		} else {
			err = s.server.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			s.observability.Logger.Error("Failed to start REST server", zap.Error(err))
		}
	}()
	return nil
}

// Stop gracefully stops the REST server.
func (s *RESTServer) Stop() error {
	s.observability.Logger.Info("Stopping REST server", zap.String("service", s.name))
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}

func loadRestConfig(obs *observability.Observability) *RestConfig {
	cfg := &RestConfig{
		Address:      ":8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	err := config.LoadConfig("Rest", cfg, obs.Logger)
	if err != nil {
		obs.Logger.Warn("Failed to load REST configuration, using defaults", zap.Error(err))
	}
	return cfg
}

func setupRouter(obs *observability.Observability) *mux.Router {
	r := mux.NewRouter()

	// Define middlewares for observability
	r.Use(observabilityMiddleware(obs))

	// Define REST endpoints
	r.HandleFunc("/health", handlers.HealthHandler).Methods("GET")
	r.HandleFunc("/status", handlers.StatusHandler).Methods("GET")

	return r
}

func setupSecurityModule(obs *observability.Observability) *security.Security {
	securityModule, secErr := security.NewSecurity(obs.Logger)
	if secErr != nil {
		obs.Logger.Error("Failed to initialize security module", zap.Error(secErr))
	}
	return securityModule
}
