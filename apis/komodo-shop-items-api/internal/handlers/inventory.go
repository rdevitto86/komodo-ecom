package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"komodo-shop-items-api/internal/config"

	httpErr "github.com/rdevitto86/komodo-forge-sdk-go/http/errors"
	logger "github.com/rdevitto86/komodo-forge-sdk-go/logging/runtime"

	"komodo-shop-items-api/internal/models"
	"komodo-shop-items-api/internal/store"
)

// Returns inventory data for all tracked items
func GetInventory(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")

	bucket := os.Getenv(config.S3_ITEMS_BUCKET)
	if bucket == "" {
		logger.Error("S3_ITEMS_BUCKET not configured", nil)
		httpErr.SendError(wtr, req, models.Err.StorageError, httpErr.WithDetail("storage not configured"))
		return
	}

	inventory, err := store.FetchInventory(req.Context(), bucket)
	if err != nil {
		logger.Error("failed to fetch inventory", err)
		httpErr.SendError(wtr, req, models.Err.InventoryUnavailable, httpErr.WithDetail("failed to retrieve inventory data"))
		return
	}

	wtr.WriteHeader(http.StatusOK)
	json.NewEncoder(wtr).Encode(inventory)
}
