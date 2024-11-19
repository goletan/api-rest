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

	cfg, err := observability.LoadObservabilityConfig(logger)

	// Initialize observability (or pass nil if observability is not configured yet)
	obs, err := observability.NewObserver(cfg)
	if err != nil {
		panic(err)
	}

	// Create a new REST service
	service := rest.NewRESTService(obs)

	// Start the service
	if err := service.Start(); err != nil {
		panic(err)
	}
}
