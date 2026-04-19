// Package shopitems provides a cart-api-local HTTP adapter for komodo-shop-items-api.
// Types are derived from komodo-shop-items-api/openapi.yaml (version 0.1.0).
// Transport is provided by forge-sdk-go/http/client.
package shopitems

import (
	"context"
	"fmt"

	httpc "github.com/rdevitto86/komodo-forge-sdk-go/http/client"

	"komodo-cart-api/internal/models"
)

// Client is the cart-api's local adapter for shop-items-api.
type Client struct {
	baseURL    string
	httpClient *httpc.Client
}

// NewClient constructs a Client for the given shop-items-api base URL.
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: httpc.NewClient(),
	}
}

// itemImage captures the image fields we need from both Product and Service image lists.
type itemImage struct {
	URL       string `json:"url"`
	IsPrimary bool   `json:"isPrimary"`
}

// itemVariant captures the variant fields we need for price and image fallback.
type itemVariant struct {
	IsDefault bool        `json:"isDefault"`
	Price     float64     `json:"price"`
	Images    []itemImage `json:"images,omitempty"`
}

// itemResponse is the minimal shape decoded from GET /item/{sku}.
// The upstream contract returns either a Product or a Service; both share
// the fields needed to populate a CartItem snapshot.
type itemResponse struct {
	ID       string        `json:"id"`
	Name     string        `json:"name"`
	Price    float64       `json:"price"`
	Variants []itemVariant `json:"variants,omitempty"`
	Images   []itemImage   `json:"images,omitempty"`
}

// GetItem fetches the catalog item for the given itemID and SKU from shop-items-api,
// returning a CartItem snapshot populated with name, price, and primary image URL.
// The returned snapshot has Quantity == 0; callers must set it before persisting.
func (c *Client) GetItem(ctx context.Context, itemID, sku string) (*models.CartItem, error) {
	item, err := httpc.GetJSON[itemResponse](c.httpClient, ctx, c.baseURL + "/item/" + sku)
	if err != nil {
		return nil, fmt.Errorf("shopitems.GetItem: %w", err)
	}

	// Resolve price: use top-level price; fall back to the default variant price.
	price := item.Price
	if price == 0 {
		for _, v := range item.Variants {
			if v.IsDefault {
				price = v.Price
				break
			}
		}
	}

	// Resolve primary image: service-level images first, then the default variant's images.
	imageURL := primaryImageFrom(item.Images)
	if imageURL == "" {
		for _, v := range item.Variants {
			if v.IsDefault {
				imageURL = primaryImageFrom(v.Images)
				break
			}
		}
	}

	return &models.CartItem{
		ItemID:         itemID,
		SKU:            sku,
		Name:           item.Name,
		UnitPriceCents: int(price * 100),
		ImageURL:       imageURL,
	}, nil
}

// primaryImageFrom returns the URL of the first image marked isPrimary,
// or the first image URL if none is marked, or "" for an empty slice.
func primaryImageFrom(images []itemImage) string {
	for _, img := range images {
		if img.IsPrimary { return img.URL }
	}
	if len(images) > 0 { return images[0].URL }
	return ""
}
