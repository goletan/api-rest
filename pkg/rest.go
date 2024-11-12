// /api-rest/pkg/rest.go
package rest

import (
	"github.com/goletan/api-rest/internal/server"
	observability "github.com/goletan/observability/pkg"
	services "github.com/goletan/services/pkg"
	"go.uber.org/zap"
)

// NewRESTService creates a new REST service that implements the Goletan service interface.
func NewRESTService(logger *zap.Logger) services.Service {
	obs, err := observability.NewObserver()
	if err != nil {
		logger.Error("Failed to initialize observability", zap.Error(err))
	}

	return server.NewRESTServer(obs)
}
