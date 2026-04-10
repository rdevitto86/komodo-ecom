//go:build e2e

package e2e_test

import (
	"net/http"
	"testing"
)

// TestHealth verifies the internal server is reachable.
// The event-bus-api exposes only an internal server (port 7002) — no public routes.
func TestHealth(t *testing.T) {
	resp := get(t, "/health", nil)
	defer resp.Body.Close()
	checkStatus(t, resp, http.StatusOK)
}

// TestPublishEvent_Valid publishes a well-formed event and expects 200 or 202.
func TestPublishEvent_Valid(t *testing.T) {
	event := map[string]any{
		"type":        "order.placed",
		"source":      "e2e-test",
		"entity_type": "order",
		"entity_id":   "order-e2e-001",
		"payload": map[string]any{
			"user_id":     "user-e2e-001",
			"total_cents": 1099,
		},
	}
	resp := post(t, "/events", event, nil)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		checkStatus(t, resp, http.StatusOK) // force fail with body context
	}
}

// TestPublishEvent_CartEvent publishes a cart domain event.
func TestPublishEvent_CartEvent(t *testing.T) {
	event := map[string]any{
		"type":        "cart.item_added",
		"source":      "e2e-test",
		"entity_type": "cart",
		"entity_id":   "cart-e2e-001",
		"payload": map[string]any{
			"user_id":  "user-e2e-001",
			"sku":      "TEST-SKU-E2E",
			"quantity": 1,
		},
	}
	resp := post(t, "/events", event, nil)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		checkStatus(t, resp, http.StatusOK)
	}
}

// TestPublishEvent_EmptyBody verifies the endpoint rejects an empty payload.
func TestPublishEvent_EmptyBody(t *testing.T) {
	resp := post(t, "/events", map[string]any{}, nil)
	defer resp.Body.Close()
	checkStatus(t, resp, http.StatusBadRequest)
}

// TestPublishEvent_MissingType verifies events without a type are rejected.
func TestPublishEvent_MissingType(t *testing.T) {
	event := map[string]any{
		"source":      "e2e-test",
		"entity_type": "order",
		"entity_id":   "order-e2e-001",
	}
	resp := post(t, "/events", event, nil)
	defer resp.Body.Close()
	checkStatus(t, resp, http.StatusBadRequest)
}
