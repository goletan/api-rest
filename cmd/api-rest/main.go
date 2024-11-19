// /api-rest/cmd/api-rest/main.go
package main

import (
	rest "github.com/goletan/api-rest/pkg"
	observability "github.com/goletan/observability/pkg"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Load the observability configuration
	cfg, err := observability.LoadObservabilityConfig(logger)
	if err != nil {
		logger.Sugar().Fatalf("Failed to load observability configuration: %v", err)
	}

	// Initialize observability with the configuration
	obs, err := observability.NewObserver(*cfg)
	if err != nil {
		logger.Sugar().Fatalf("Failed to initialize observability: %v", err)
	}

	// Create a new REST service
	service := rest.NewRESTService(obs)

	// Start the service
	if err := service.Start(); err != nil {
		panic(err)
	}
}
