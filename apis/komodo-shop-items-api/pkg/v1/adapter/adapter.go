package adapter

import (
	"context"
	"fmt"

	awsS3 "komodo-forge-sdk-go/aws/s3"
	logger "komodo-forge-sdk-go/logging/runtime"
)

// Fetches a product by SKU from the S3 items bucket
func FetchProductBySKU(ctx context.Context, bucket string, sku string) (*Product, error) {
	key := fmt.Sprintf("products/%s.json", sku)
	var product Product
	if err := awsS3.GetObjectAs(ctx, bucket, key, &product); err != nil {
		logger.Error("failed to fetch product from s3: "+sku, err)
		return nil, err
	}
	return &product, nil
}

// Fetches a service by SKU from the S3 items bucket
func FetchServiceBySKU(ctx context.Context, bucket string, sku string) (*Service, error) {
	key := fmt.Sprintf("services/%s.json", sku)
	var service Service
	if err := awsS3.GetObjectAs(ctx, bucket, key, &service); err != nil {
		logger.Error("failed to fetch service from s3: "+sku, err)
		return nil, err
	}
	return &service, nil
}

// Fetches the full inventory manifest from S3
func FetchInventory(ctx context.Context, bucket string) (*InventoryResponse, error) {
	key := "inventory/manifest.json"
	var inventory InventoryResponse
	if err := awsS3.GetObjectAs(ctx, bucket, key, &inventory); err != nil {
		logger.Error("failed to fetch inventory from s3", err)
		return nil, err
	}
	return &inventory, nil
}
