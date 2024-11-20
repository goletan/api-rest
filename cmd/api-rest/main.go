// /api-rest/cmd/api-rest/main.go
package main

import (
	"github.com/goletan/api-rest/internal/server"
	observability "github.com/goletan/observability/pkg"
	"go.uber.org/zap"
)

func main() {
	// Initialize observability with the configuration
	obs, err := observability.NewObserver()
	if err != nil {
		panic(err)
	}

	// Create a new REST server instance
	restServer := server.NewRESTServer(obs)
	// Initialize the REST server
	if err := restServer.Initialize(); err != nil {
		obs.Logger.Error("Failed to initialize REST server: %v", zap.Error(err))
	}

	// Start the service
	if err := restServer.Start(); err != nil {
		panic(err)
	}
}
