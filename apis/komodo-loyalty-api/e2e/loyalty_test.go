//go:build e2e

package e2e_test

import (
	"net/http"
	"testing"
)

// TestHealth verifies the service is reachable and returns 200 on /health.
// TODO: remove the t.Skip once handlers are scaffolded and the service is running.
func TestHealth(t *testing.T) {
	t.Skip("service not yet implemented — scaffold routes before enabling e2e tests")
	resp := get(t, "/health", nil)
	defer resp.Body.Close()
	checkStatus(t, resp, http.StatusOK)
}
