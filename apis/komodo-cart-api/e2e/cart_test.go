//go:build e2e

package e2e_test

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestHealth(t *testing.T) {
	resp := get(t, "/health", nil)
	defer resp.Body.Close()
	checkStatus(t, resp, http.StatusOK)
}

// --- Guest cart ---

// TestGuestCart_CreateAndGet creates a guest cart and reads it back by ID.
func TestGuestCart_CreateAndGet(t *testing.T) {
	cartID, _ := createGuestCart(t)

	resp := get(t, "/cart/"+cartID, nil)
	defer resp.Body.Close()
	checkStatus(t, resp, http.StatusOK)
}

// TestGuestCart_AddAndRemoveItem adds an item to a guest cart then removes it.
// Skips if the test SKU is not present in the shop-items-api catalog.
func TestGuestCart_AddAndRemoveItem(t *testing.T) {
	cartID, sessionID := createGuestCart(t)
	h := map[string]string{"X-Session-ID": sessionID}

	addResp := post(t, "/cart/"+cartID+"/items",
		map[string]any{"sku": "TEST-SKU-E2E", "quantity": 1},
		h,
	)
	defer addResp.Body.Close()
	if addResp.StatusCode == http.StatusNotFound || addResp.StatusCode == http.StatusUnprocessableEntity {
		t.Skip("TEST-SKU-E2E not in shop-items-api — seed the item in S3 to enable this test")
	}
	checkStatus(t, addResp, http.StatusOK)

	var updated struct {
		Items []struct {
			ItemID string `json:"item_id"`
		} `json:"items"`
	}
	decodeJSON(t, addResp, &updated)
	if len(updated.Items) == 0 {
		t.Fatal("expected at least one item in cart after add")
	}

	itemID := updated.Items[0].ItemID
	delResp := del(t, "/cart/"+cartID+"/items/"+itemID, h)
	defer delResp.Body.Close()
	checkStatus(t, delResp, http.StatusOK)
}

// TestGuestCart_UpdateItemQuantity adds then updates an item's quantity.
func TestGuestCart_UpdateItemQuantity(t *testing.T) {
	cartID, sessionID := createGuestCart(t)
	h := map[string]string{"X-Session-ID": sessionID}

	addResp := post(t, "/cart/"+cartID+"/items",
		map[string]any{"sku": "TEST-SKU-E2E", "quantity": 1},
		h,
	)
	defer addResp.Body.Close()
	if addResp.StatusCode == http.StatusNotFound || addResp.StatusCode == http.StatusUnprocessableEntity {
		t.Skip("TEST-SKU-E2E not in shop-items-api — seed the item to enable this test")
	}
	checkStatus(t, addResp, http.StatusOK)

	var updated struct {
		Items []struct {
			ItemID string `json:"item_id"`
		} `json:"items"`
	}
	decodeJSON(t, addResp, &updated)
	if len(updated.Items) == 0 {
		t.Fatal("expected at least one item")
	}

	itemID := updated.Items[0].ItemID
	putResp := put(t, "/cart/"+cartID+"/items/"+itemID,
		map[string]any{"quantity": 3},
		h,
	)
	defer putResp.Body.Close()
	checkStatus(t, putResp, http.StatusOK)
}

// TestGuestCart_Clear creates and then clears a guest cart.
func TestGuestCart_Clear(t *testing.T) {
	cartID, sessionID := createGuestCart(t)
	h := map[string]string{"X-Session-ID": sessionID}

	resp := del(t, "/cart/"+cartID, h)
	defer resp.Body.Close()
	checkStatus(t, resp, http.StatusOK)
}

// TestGuestCart_GetUnknown verifies 404 for a non-existent cart ID.
func TestGuestCart_GetUnknown(t *testing.T) {
	resp := get(t, "/cart/cart-does-not-exist", nil)
	defer resp.Body.Close()
	checkStatus(t, resp, http.StatusNotFound)
}

// --- Authenticated cart ---

// TestAuthCart_NoToken verifies authenticated cart routes reject missing JWTs.
func TestAuthCart_NoToken(t *testing.T) {
	resp := get(t, "/me/cart", nil)
	defer resp.Body.Close()
	checkStatus(t, resp, http.StatusUnauthorized)
}

// TestAuthCart_GetCart fetches the authenticated user's cart.
func TestAuthCart_GetCart(t *testing.T) {
	resp := get(t, "/me/cart", authHeader(t))
	defer resp.Body.Close()
	checkStatus(t, resp, http.StatusOK)
}

// TestAuthCart_AddAndRemoveItem adds an item to the authenticated cart then removes it.
func TestAuthCart_AddAndRemoveItem(t *testing.T) {
	h := authHeader(t)

	addResp := post(t, "/me/cart/items",
		map[string]any{"sku": "TEST-SKU-E2E", "quantity": 1},
		h,
	)
	defer addResp.Body.Close()
	if addResp.StatusCode == http.StatusNotFound || addResp.StatusCode == http.StatusUnprocessableEntity {
		t.Skip("TEST-SKU-E2E not in shop-items-api — seed the item to enable this test")
	}
	checkStatus(t, addResp, http.StatusOK)

	var updated struct {
		Items []struct {
			ItemID string `json:"item_id"`
		} `json:"items"`
	}
	decodeJSON(t, addResp, &updated)
	if len(updated.Items) == 0 {
		t.Fatal("expected at least one item after add")
	}

	itemID := updated.Items[0].ItemID
	delResp := del(t, "/me/cart/items/"+itemID, h)
	defer delResp.Body.Close()
	checkStatus(t, delResp, http.StatusOK)
}

// TestAuthCart_ClearCart clears the authenticated user's cart.
func TestAuthCart_ClearCart(t *testing.T) {
	h := authHeader(t)
	resp := del(t, "/me/cart", h)
	defer resp.Body.Close()
	checkStatus(t, resp, http.StatusOK)
}

// createGuestCart is a helper that creates a guest cart and returns (cartID, sessionID).
func createGuestCart(t *testing.T) (string, string) {
	t.Helper()
	resp := post(t, "/cart", nil, nil)
	defer resp.Body.Close()
	checkStatus(t, resp, http.StatusCreated)

	sessionID := resp.Header.Get("X-Session-ID")

	var cart struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&cart); err != nil {
		t.Fatalf("decode create-cart response: %v", err)
	}
	if cart.ID == "" {
		t.Fatal("expected non-empty cart id in response")
	}
	return cart.ID, sessionID
}
