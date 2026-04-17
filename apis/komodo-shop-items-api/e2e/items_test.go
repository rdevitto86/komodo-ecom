//go:build e2e

package e2e_test

import (
	"net/http"
	"testing"
)

func TestHealth(t *testing.T) {
	res := get(t, "/health", nil)
	defer res.Body.Close()
	checkStatus(t, res, http.StatusOK)
}

// TestGetInventory_ReturnsItemList fetches the full inventory list from S3.
// An empty result is valid if the S3 bucket has not been seeded.
func TestGetInventory_ReturnsItemList(t *testing.T) {
	res := get(t, "/item/inventory", nil)
	defer res.Body.Close()
	checkStatus(t, res, http.StatusOK)
}

// TestGetItemBySKU_NotFound verifies 404 for a SKU that does not exist in S3.
func TestGetItemBySKU_NotFound(t *testing.T) {
	res := get(t, "/item/SKU-DOES-NOT-EXIST-E2E", nil)
	defer res.Body.Close()
	checkStatus(t, res, http.StatusNotFound)
}

// TestGetItemBySKU_Valid fetches a known test item from S3.
// Skips if TEST-SKU-E2E is not present in the S3 bucket.
func TestGetItemBySKU_Valid(t *testing.T) {
	res := get(t, "/item/TEST-SKU-E2E", nil)
	defer res.Body.Close()
	if res.StatusCode == http.StatusNotFound {
		t.Skip("TEST-SKU-E2E not in S3 catalog — upload a test item JSON to S3_ITEMS_BUCKET to enable")
	}
	checkStatus(t, res, http.StatusOK)

	var item struct {
		SKU  string `json:"sku"`
		Name string `json:"name"`
	}
	decodeJSON(t, res, &item)
	if item.SKU == "" {
		t.Fatal("expected non-empty sku in item response")
	}
}

// TestGetSuggestions_RequiresAuth verifies the suggestions endpoint requires a JWT.
func TestGetSuggestions_RequiresAuth(t *testing.T) {
	res := post(t, "/item/suggestion", map[string]any{"limit": 6}, nil)
	defer res.Body.Close()
	checkStatus(t, res, http.StatusUnauthorized)
}

// TestGetSuggestions_Authenticated fetches item suggestions for the authenticated user.
func TestGetSuggestions_Authenticated(t *testing.T) {
	res := post(t, "/item/suggestion", map[string]any{"limit": 6}, authHeader(t))
	defer res.Body.Close()
	checkStatus(t, res, http.StatusOK)
}

// TestGetSuggestions_DefaultLimit verifies suggestions are returned without an explicit limit.
func TestGetSuggestions_DefaultLimit(t *testing.T) {
	res := post(t, "/item/suggestion", map[string]any{}, authHeader(t))
	defer res.Body.Close()
	checkStatus(t, res, http.StatusOK)
}
