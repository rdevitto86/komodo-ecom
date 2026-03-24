package client

import (
	"context"
	"fmt"

	awsS3 "github.com/rdevitto86/komodo-forge-sdk-go/aws/s3"
	logger "github.com/rdevitto86/komodo-forge-sdk-go/logging/runtime"

	"komodo-shop-items-api/pkg/v1/models"
)

// FetchProductBySKU fetches a product by SKU from the S3 items bucket.
func FetchProductBySKU(ctx context.Context, bucket string, sku string) (*models.Product, error) {
	key := fmt.Sprintf("products/%s.json", sku)
	var product models.Product
	if err := awsS3.GetObjectAs(ctx, bucket, key, &product); err != nil {
		logger.Error("failed to fetch product from s3: "+sku, err)
		return nil, err
	}
	return &product, nil
}

// FetchServiceBySKU fetches a service by SKU from the S3 items bucket.
func FetchServiceBySKU(ctx context.Context, bucket string, sku string) (*models.Service, error) {
	key := fmt.Sprintf("services/%s.json", sku)
	var service models.Service
	if err := awsS3.GetObjectAs(ctx, bucket, key, &service); err != nil {
		logger.Error("failed to fetch service from s3: "+sku, err)
		return nil, err
	}
	return &service, nil
}

// FetchInventory fetches the full inventory manifest from S3.
func FetchInventory(ctx context.Context, bucket string) (*models.InventoryResponse, error) {
	key := "inventory/manifest.json"
	var inventory models.InventoryResponse
	if err := awsS3.GetObjectAs(ctx, bucket, key, &inventory); err != nil {
		logger.Error("failed to fetch inventory from s3", err)
		return nil, err
	}
	return &inventory, nil
}
