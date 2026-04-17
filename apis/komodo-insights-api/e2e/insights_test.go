//go:build e2e

package e2e_test

import (
	"testing"
)

func TestHealth(t *testing.T) {
	res := get(t, "/health", nil)
	defer res.Body.Close()
	checkStatus(t, res, 200)
}

// TestGetItemSummary_Stub verifies the route exists and returns 404 (not found)
// until the LLM provider is wired. Update to 200 + response shape assertion
// once GetItemSummary is implemented.
func TestGetItemSummary_Stub(t *testing.T) {
	res := get(t, "/insights/items/test-item-123/summary", authHeader(t))
	defer res.Body.Close()
	checkStatus(t, res, 404)
}

// TestGetItemSentiment_Stub verifies the route exists and returns 404 (not found)
// until the LLM provider is wired.
func TestGetItemSentiment_Stub(t *testing.T) {
	res := get(t, "/insights/items/test-item-123/sentiment", authHeader(t))
	defer res.Body.Close()
	checkStatus(t, res, 404)
}

// TestGetTrending_Stub verifies the route exists and returns 500 (service not yet
// implemented) until trending signal sources are wired.
func TestGetTrending_Stub(t *testing.T) {
	res := get(t, "/insights/trending", authHeader(t))
	defer res.Body.Close()
	checkStatus(t, res, 500)
}
