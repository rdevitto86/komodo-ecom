package handlers

import (
	"encoding/json"
	"net/http"

	"komodo-forge-sdk-go/config"
	httpErr "komodo-forge-sdk-go/http/errors"
	shopitems "komodo-forge-sdk-go/http/services/shop_items"
	logger "komodo-forge-sdk-go/logging/runtime"
)

// Returns a single item (product or service) by SKU
func GetItemBySKU(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")

	sku := req.PathValue("sku")
	if sku == "" {
		httpErr.SendError(wtr, req, httpErr.ShopItem.InvalidSKU, httpErr.WithDetail("sku path parameter is required"))
		return
	}

	bucket := config.GetConfigValue("S3_ITEMS_BUCKET")
	if bucket == "" {
		logger.Error("S3_ITEMS_BUCKET not configured", nil)
		httpErr.SendError(wtr, req, httpErr.ShopItem.StorageError, httpErr.WithDetail("storage not configured"))
		return
	}

	// Try product first, then fall back to service
	product, err := shopitems.FetchProductBySKU(req.Context(), bucket, sku)
	if err == nil {
		wtr.WriteHeader(http.StatusOK)
		json.NewEncoder(wtr).Encode(product)
		return
	}

	service, err := shopitems.FetchServiceBySKU(req.Context(), bucket, sku)
	if err == nil {
		wtr.WriteHeader(http.StatusOK)
		json.NewEncoder(wtr).Encode(service)
		return
	}

	logger.Warn("item not found for sku: " + sku)
	httpErr.SendError(wtr, req, httpErr.ShopItem.ItemNotFound, httpErr.WithDetail("no product or service found for sku: "+sku))
}
