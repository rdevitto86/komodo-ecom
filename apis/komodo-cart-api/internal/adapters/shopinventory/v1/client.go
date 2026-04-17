// Package shopinventory provides a cart-api-local HTTP adapter for komodo-shop-inventory-api.
// Types are derived from komodo-shop-inventory-api/openapi.yaml (version 1.0.0).
// Transport is provided by forge-sdk-go/http/client.
package shopinventory

import (
	"context"
	"errors"
	"fmt"
	"time"

	httpc "github.com/rdevitto86/komodo-forge-sdk-go/http/client"
)

// HoldError is returned when the upstream responds with 409 (insufficient stock)
// or any other non-201 HTTP status. StatusCode allows callers to distinguish 409
// conflicts from network/gateway errors without parsing the body.
type HoldError struct {
	StatusCode int
	Detail     string
}

func (e *HoldError) Error() string {
	return fmt.Sprintf("shopinventory: hold failed with status %d: %s", e.StatusCode, e.Detail)
}

// Client is the cart-api's local adapter for shop-inventory-api.
type Client struct {
	baseURL    string
	httpClient *httpc.Client
}

// NewClient constructs a Client for the given shop-inventory-api base URL.
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: httpc.NewClient(),
	}
}

// reserveRequest maps to the ReserveRequest schema in shop-inventory-api's openapi.yaml.
type reserveRequest struct {
	CartID   string `json:"cart_id"`
	Quantity int    `json:"quantity"`
}

// holdResponse maps to the HoldResponse schema in shop-inventory-api's openapi.yaml.
type holdResponse struct {
	HoldID    string    `json:"hold_id"`
	SKU       string    `json:"sku"`
	Quantity  int       `json:"quantity"`
	ExpiresAt time.Time `json:"expires_at"`
}

// PlaceHold calls POST /stock/{sku}/reserve and returns the hold ID and expiry.
// Returns a *HoldError with StatusCode == 409 when stock is insufficient.
// Returns a *HoldError with the upstream status for any other non-201 response.
func (c *Client) PlaceHold(ctx context.Context, sku, cartID string, qty int) (holdID string, holdExpiry time.Time, err error) {
	url := c.baseURL + "/stock/" + sku + "/reserve"
	hold, err := httpc.PostJSON[holdResponse](c.httpClient, ctx, url, reserveRequest{CartID: cartID, Quantity: qty})
	if err != nil {
		var httpErr *httpc.HTTPError
		if errors.As(err, &httpErr) {
			return "", time.Time{}, &HoldError{StatusCode: httpErr.StatusCode, Detail: string(httpErr.Body)}
		}
		return "", time.Time{}, fmt.Errorf("shopinventory.PlaceHold: %w", err)
	}
	return hold.HoldID, hold.ExpiresAt, nil
}
