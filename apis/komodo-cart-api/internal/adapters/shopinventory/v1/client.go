// Package shopinventory provides a cart-api-local HTTP adapter for komodo-shop-inventory-api.
// Types are derived from komodo-shop-inventory-api/openapi.yaml (version 1.0.0).
// Transport is provided by forge-sdk-go/http/client.
package shopinventory

import (
	"context"
	"errors"
	"fmt"

	httpc "github.com/rdevitto86/komodo-forge-sdk-go/http/client"
)

// OutOfStockError is returned when available_qty is below the requested quantity,
// or when the upstream responds with a non-200 status. StatusCode lets callers
// distinguish a true OOS (409) from a network/gateway error without parsing the body.
type OutOfStockError struct {
	StatusCode int
	Detail     string
}

func (e *OutOfStockError) Error() string {
	return fmt.Sprintf("shopinventory: stock check failed with status %d: %s", e.StatusCode, e.Detail)
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

// stockLevel maps to the StockLevel schema in shop-inventory-api's openapi.yaml.
type stockLevel struct {
	SKU          string `json:"sku"`
	AvailableQty int    `json:"available_qty"`
}

// CheckStock calls GET /stock/{sku} and returns nil if available_qty >= qty.
// Returns *OutOfStockError with StatusCode 409 when stock is insufficient.
// Returns *OutOfStockError with the upstream status for any other non-200 response.
func (c *Client) CheckStock(ctx context.Context, sku string, qty int) error {
	url := c.baseURL + "/stock/" + sku
	level, err := httpc.GetJSON[stockLevel](c.httpClient, ctx, url)

	if err != nil {
		var httpErr *httpc.HTTPError
		if errors.As(err, &httpErr) {
			return &OutOfStockError{StatusCode: httpErr.StatusCode, Detail: string(httpErr.Body)}
		}
		return fmt.Errorf("stock check failed: %w", err)
	}

	if level.AvailableQty < qty {
		return &OutOfStockError{
			StatusCode: 409,
			Detail: fmt.Sprintf("available: %d, requested: %d", level.AvailableQty, qty),
		}
	}
	return nil
}
