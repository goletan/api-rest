// /api-graphql/test/integration/rest_test.go
package integration_test

import (
	"crypto/tls"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRESTAPIHealth(t *testing.T) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get("https://localhost:8443/health")
	assert.NoError(t, err, "Should be able to call health endpoint")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Health endpoint should return status OK")
}
