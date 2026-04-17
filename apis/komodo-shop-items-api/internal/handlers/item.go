package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	httpErr "github.com/rdevitto86/komodo-forge-sdk-go/http/errors"
	logger "github.com/rdevitto86/komodo-forge-sdk-go/logging/runtime"

	"komodo-shop-items-api/internal/models"
	"komodo-shop-items-api/internal/store"
)

// Returns a single item (product or service) by SKU
func GetItemBySKU(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")

	sku := req.PathValue("sku")
	if sku == "" {
		httpErr.SendError(wtr, req, models.Err.InvalidSKU, httpErr.WithDetail("sku path parameter is required"))
		return
	}

	bucket := os.Getenv("S3_ITEMS_BUCKET")
	if bucket == "" {
		logger.Error("S3_ITEMS_BUCKET not configured", nil)
		httpErr.SendError(wtr, req, models.Err.StorageError, httpErr.WithDetail("storage not configured"))
		return
	}

	// Try product first, then fall back to service
	product, err := store.FetchProductBySKU(req.Context(), bucket, sku)
	if err == nil {
		wtr.WriteHeader(http.StatusOK)
		json.NewEncoder(wtr).Encode(product)
		return
	}

	service, err := store.FetchServiceBySKU(req.Context(), bucket, sku)
	if err == nil {
		wtr.WriteHeader(http.StatusOK)
		json.NewEncoder(wtr).Encode(service)
		return
	}

	logger.Warn("item not found for sku: " + sku)
	httpErr.SendError(wtr, req, models.Err.ItemNotFound, httpErr.WithDetail("no product or service found for sku: "+sku))
}
