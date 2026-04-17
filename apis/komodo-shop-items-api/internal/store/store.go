// Package store fetches catalog and inventory data from S3.
// All JSON objects are stored under keys like "products/{sku}.json",
// "services/{sku}.json", and "inventory.json" in the configured bucket.
package store

import (
	"context"
	"fmt"

	"github.com/rdevitto86/komodo-forge-sdk-go/aws/s3"

	"komodo-shop-items-api/internal/models"
)

// FetchProductBySKU retrieves a product from S3 by SKU.
// Key pattern: products/{sku}.json
func FetchProductBySKU(ctx context.Context, bucket, sku string) (*models.Product, error) {
	key := "products/" + sku + ".json"
	var product models.Product
	if err := s3.GetObjectAs(ctx, bucket, key, &product); err != nil {
		return nil, fmt.Errorf("store.FetchProductBySKU: %w", err)
	}
	return &product, nil
}

// FetchServiceBySKU retrieves a service from S3 by SKU.
// Key pattern: services/{sku}.json
func FetchServiceBySKU(ctx context.Context, bucket, sku string) (*models.Service, error) {
	key := "services/" + sku + ".json"
	var svc models.Service
	if err := s3.GetObjectAs(ctx, bucket, key, &svc); err != nil {
		return nil, fmt.Errorf("store.FetchServiceBySKU: %w", err)
	}
	return &svc, nil
}

// FetchInventory retrieves the full inventory response from S3.
// Key pattern: inventory.json
func FetchInventory(ctx context.Context, bucket string) (*models.InventoryResponse, error) {
	key := "inventory.json"
	var inv models.InventoryResponse
	if err := s3.GetObjectAs(ctx, bucket, key, &inv); err != nil {
		return nil, fmt.Errorf("store.FetchInventory: %w", err)
	}
	return &inv, nil
}
