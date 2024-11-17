// /api-rest/pkg/rest.go
package rest

import (
	"github.com/goletan/api-rest/internal/server"
	observability "github.com/goletan/observability/pkg"
	services "github.com/goletan/services/pkg"
)

// NewRESTService creates a new REST service that implements the Goletan service interface.
func NewRESTService(obs *observability.Observability) services.Service {
	return server.NewRESTServer(obs)
}
