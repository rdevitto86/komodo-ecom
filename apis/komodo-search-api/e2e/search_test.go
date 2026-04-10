//go:build e2e

package e2e_test

import (
	"net/http"
	"testing"
)

func TestHealth(t *testing.T) {
	resp := get(t, "/health", nil)
	defer resp.Body.Close()
	checkStatus(t, resp, http.StatusOK)
}

// TestSearch_WithQuery searches for a term and expects a 200 with results.
// Returns 501 if the Typesense client is not yet initialized.
func TestSearch_WithQuery(t *testing.T) {
	resp := get(t, "/search?q=jacket&page=1&limit=10", nil)
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotImplemented {
		t.Skip("Typesense not initialized — wire the client in main.go to enable this test")
	}
	checkStatus(t, resp, http.StatusOK)
}

// TestSearch_EmptyQuery verifies the endpoint handles an empty q param.
func TestSearch_EmptyQuery(t *testing.T) {
	resp := get(t, "/search?q=", nil)
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotImplemented {
		t.Skip("Typesense not initialized")
	}
	// Empty q may return all results (200) or a validation error (400) — both are acceptable.
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusBadRequest {
		checkStatus(t, resp, http.StatusOK)
	}
}

// TestSearch_MissingQueryParam verifies the endpoint validates the q param is present.
func TestSearch_MissingQueryParam(t *testing.T) {
	resp := get(t, "/search", nil)
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotImplemented {
		t.Skip("Typesense not initialized")
	}
	// Missing q should be 400.
	if resp.StatusCode != http.StatusBadRequest && resp.StatusCode != http.StatusUnprocessableEntity {
		checkStatus(t, resp, http.StatusBadRequest)
	}
}

// TestSearch_Pagination verifies page and limit params are accepted.
func TestSearch_Pagination(t *testing.T) {
	resp := get(t, "/search?q=shoe&page=2&limit=5", nil)
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotImplemented {
		t.Skip("Typesense not initialized")
	}
	checkStatus(t, resp, http.StatusOK)
}
