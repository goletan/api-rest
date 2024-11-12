// /api-rest/pkg/rest.go
package rest

import (
	"fmt"

	"github.com/goletan/api-rest/internal/server"
	observability "github.com/goletan/observability/pkg"
	services "github.com/goletan/services/pkg"
)

// NewRESTService creates a new REST service that implements the Goletan service interface.
func NewRESTService() services.Service {
	obs, err := observability.NewObserver()
	if err != nil {
		fmt.Errorf("cannot create new observer")
	}

	return server.NewRESTServer(obs)
}
