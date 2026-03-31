package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"komodo-cart-api/pkg/v1/models"
	"net/http"
	"time"
)

// ShopItemsClient fetches product snapshots from shop-items-api at add-item time.
type ShopItemsClient struct {
	baseURL string
	http    *http.Client
}

// NewShopItemsClient creates a ShopItemsClient targeting baseURL.
func NewShopItemsClient(baseURL string) *ShopItemsClient {
	return &ShopItemsClient{
		baseURL: baseURL,
		http:    &http.Client{Timeout: 5 * time.Second},
	}
}

// shopItemResponse is the subset of the shop-items-api response we care about.
type shopItemResponse struct {
	Name           string `json:"name"`
	UnitPriceCents int    `json:"unit_price_cents"`
	ImageURL       string `json:"image_url"`
}

// GetItem fetches product metadata for itemID/sku and returns a CartItem snapshot.
// Quantity is left at 0 — the caller sets it.
func (c *ShopItemsClient) GetItem(ctx context.Context, itemID, sku string) (*models.CartItem, error) {
	url := fmt.Sprintf("%s/item/%s", c.baseURL, itemID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("client.ShopItems.GetItem: build request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client.ShopItems.GetItem: do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("client.ShopItems.GetItem: non-2xx response %d: %s", resp.StatusCode, body)
	}

	var item shopItemResponse
	if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
		return nil, fmt.Errorf("client.ShopItems.GetItem: decode response: %w", err)
	}

	return &models.CartItem{
		ItemID:         itemID,
		SKU:            sku,
		Name:           item.Name,
		UnitPriceCents: item.UnitPriceCents,
		ImageURL:       item.ImageURL,
		// Quantity intentionally left at 0 — caller sets it
	}, nil
}

// HoldError is returned by PlaceHold when inventory signals a stock conflict (409).
type HoldError struct {
	StatusCode int
	Detail     string
}

func (e *HoldError) Error() string {
	return fmt.Sprintf("inventory hold failed (status %d): %s", e.StatusCode, e.Detail)
}

// holdRequest is the body sent to inventory-api for a stock hold.
type holdRequest struct {
	CartID   string `json:"cart_id"`
	Quantity int    `json:"quantity"`
}

// holdResponse is the body returned by inventory-api on a successful hold.
type holdResponse struct {
	HoldID    string    `json:"hold_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

// InventoryClient places stock holds via shop-inventory-api.
type InventoryClient struct {
	baseURL string
	http    *http.Client
}

// NewInventoryClient creates an InventoryClient targeting baseURL.
func NewInventoryClient(baseURL string) *InventoryClient {
	return &InventoryClient{
		baseURL: baseURL,
		http:    &http.Client{Timeout: 10 * time.Second},
	}
}

// PlaceHold reserves qty units of sku for cartID.
// Returns (holdID, expiresAt, nil) on success.
// Returns (*HoldError, ...) on 409 stock conflict.
// Returns a wrapped error on network or unexpected failures.
func (c *InventoryClient) PlaceHold(ctx context.Context, sku, cartID string, qty int) (string, time.Time, error) {
	body, err := json.Marshal(holdRequest{CartID: cartID, Quantity: qty})
	if err != nil {
		return "", time.Time{}, fmt.Errorf("client.Inventory.PlaceHold: marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/stock/%s/reserve", c.baseURL, sku)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("client.Inventory.PlaceHold: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("client.Inventory.PlaceHold: do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusConflict {
		detail, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return "", time.Time{}, &HoldError{StatusCode: 409, Detail: string(detail)}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		detail, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return "", time.Time{}, fmt.Errorf("client.Inventory.PlaceHold: non-2xx response %d: %s", resp.StatusCode, detail)
	}

	var hold holdResponse
	if err := json.NewDecoder(resp.Body).Decode(&hold); err != nil {
		return "", time.Time{}, fmt.Errorf("client.Inventory.PlaceHold: decode response: %w", err)
	}
	return hold.HoldID, hold.ExpiresAt, nil
}
