// Package store fetches catalog and inventory data from S3.
// All JSON objects are stored under keys like "products/{sku}.json",
// "services/{sku}.json", and "inventory.json" in the configured bucket.
package store

import (
	"context"
	"fmt"

	"github.com/rdevitto86/komodo-forge-sdk-go/aws/s3"
	logger "github.com/rdevitto86/komodo-forge-sdk-go/logging/runtime"

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

// FetchServiceByID retrieves a service from S3 by its catalog ID.
// The catalog convention is that a service's ID matches its SKU, so the key
// pattern is identical to FetchServiceBySKU: services/{id}.json.
func FetchServiceByID(ctx context.Context, bucket, id string) (*models.Service, error) {
	key := "services/" + id + ".json"
	var svc models.Service
	if err := s3.GetObjectAs(ctx, bucket, key, &svc); err != nil {
		return nil, fmt.Errorf("store.FetchServiceByID: %w", err)
	}
	return &svc, nil
}

// FetchAllServices enumerates all service SKUs from the inventory manifest and
// fetches each service object. Items that fail to load are skipped with a
// warning so a single corrupt object cannot degrade the full listing.
// Key pattern: inventory.json → services/{sku}.json
func FetchAllServices(ctx context.Context, bucket string) ([]models.Service, error) {
	inv, err := FetchInventory(ctx, bucket)
	if err != nil {
		return nil, fmt.Errorf("store.FetchAllServices: %w", err)
	}

	services := make([]models.Service, 0, len(inv.Items))
	for _, item := range inv.Items {
		svc, err := FetchServiceBySKU(ctx, bucket, item.SKU)
		if err != nil {
			// A missing or malformed service object should not abort the listing.
			logger.Warn("store.FetchAllServices: skipping sku " + item.SKU + ": " + err.Error())
			continue
		}
		services = append(services, *svc)
	}
	return services, nil
}

// FetchAllProducts enumerates all product SKUs from the inventory manifest and
// fetches each product object. Items that fail to load are skipped with a
// warning so a single corrupt object cannot degrade the full listing.
// Key pattern: inventory.json → products/{sku}.json
func FetchAllProducts(ctx context.Context, bucket string) ([]models.Product, error) {
	inv, err := FetchInventory(ctx, bucket)
	if err != nil {
		return nil, fmt.Errorf("store.FetchAllProducts: %w", err)
	}

	products := make([]models.Product, 0, len(inv.Items))
	for _, item := range inv.Items {
		product, err := FetchProductBySKU(ctx, bucket, item.SKU)
		if err != nil {
			// A missing or malformed product object should not abort the listing.
			logger.Warn("store.FetchAllProducts: skipping sku " + item.SKU + ": " + err.Error())
			continue
		}
		products = append(products, *product)
	}
	return products, nil
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
