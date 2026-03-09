package handlers

import (
	"encoding/json"
	"net/http"

	"komodo-forge-sdk-go/config"
	httpErr "komodo-forge-sdk-go/http/errors"
	shopitems "komodo-forge-sdk-go/http/services/shop_items"
	logger "komodo-forge-sdk-go/logging/runtime"
)

// Returns inventory data for all tracked items
func GetInventory(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")

	bucket := config.GetConfigValue("S3_ITEMS_BUCKET")
	if bucket == "" {
		logger.Error("S3_ITEMS_BUCKET not configured", nil)
		httpErr.SendError(wtr, req, httpErr.ShopItem.StorageError, httpErr.WithDetail("storage not configured"))
		return
	}

	inventory, err := shopitems.FetchInventory(req.Context(), bucket)
	if err != nil {
		logger.Error("failed to fetch inventory", err)
		httpErr.SendError(wtr, req, httpErr.ShopItem.InventoryUnavailable, httpErr.WithDetail("failed to retrieve inventory data"))
		return
	}

	wtr.WriteHeader(http.StatusOK)
	json.NewEncoder(wtr).Encode(inventory)
}
