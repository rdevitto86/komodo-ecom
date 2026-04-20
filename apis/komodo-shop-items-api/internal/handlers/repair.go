package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"komodo-shop-items-api/internal/config"

	httpErr "github.com/rdevitto86/komodo-forge-sdk-go/http/errors"
	logger "github.com/rdevitto86/komodo-forge-sdk-go/logging/runtime"

	"komodo-shop-items-api/internal/models"
	"komodo-shop-items-api/internal/store"
)

const (
	defaultRepairPageLimit = 20
	maxRepairPageLimit     = 100
)

// GetRepairServices handles GET /services/repair.
// Loads all services from S3 (via the inventory manifest), filters to those
// with service_type == "repair", and returns a paginated list.
func GetRepairServices(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")

	bucket := os.Getenv(config.S3_ITEMS_BUCKET)
	if bucket == "" {
		logger.Error("S3_ITEMS_BUCKET not configured", nil)
		httpErr.SendError(wtr, req, models.Err.StorageError, httpErr.WithDetail("storage not configured"))
		return
	}

	page, limit := parsePaginationParams(req)

	all, err := store.FetchAllServices(req.Context(), bucket)
	if err != nil {
		logger.Error("failed to fetch services for repair listing", err)
		httpErr.SendError(wtr, req, models.Err.StorageError, httpErr.WithDetail("failed to retrieve service catalog"))
		return
	}

	// Filter to repair-type services only.
	repairs := make([]models.Service, 0, len(all))
	for _, svc := range all {
		if svc.ServiceType == models.ServiceTypeRepair {
			repairs = append(repairs, svc)
		}
	}

	total := len(repairs)
	start, end := paginateSlice(total, page, limit)

	resp := models.RepairServicesResponse{
		Items: repairs[start:end],
		Total: total,
		Page:  page,
		Limit: limit,
	}

	wtr.WriteHeader(http.StatusOK)
	json.NewEncoder(wtr).Encode(resp)
}

// GetRepairService handles GET /services/repair/{id}.
// Fetches a service by its catalog ID. Returns 404 if the item does not exist
// or if its service_type is not "repair".
func GetRepairService(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")

	id := req.PathValue("id")
	if id == "" {
		httpErr.SendError(wtr, req, models.Err.InvalidSKU, httpErr.WithDetail("id path parameter is required"))
		return
	}

	bucket := os.Getenv(config.S3_ITEMS_BUCKET)
	if bucket == "" {
		logger.Error("S3_ITEMS_BUCKET not configured", nil)
		httpErr.SendError(wtr, req, models.Err.StorageError, httpErr.WithDetail("storage not configured"))
		return
	}

	svc, err := store.FetchServiceByID(req.Context(), bucket, id)
	if err != nil {
		logger.Warn("repair service not found for id: " + id)
		httpErr.SendError(wtr, req, models.Err.ItemNotFound, httpErr.WithDetail("no repair service found for id: "+id))
		return
	}

	if svc.ServiceType != models.ServiceTypeRepair {
		// Item exists but is not a repair service — treat as not found to avoid
		// leaking catalog structure through the repair-specific endpoint.
		logger.Warn("item " + id + " found but service_type is not repair")
		httpErr.SendError(wtr, req, models.Err.ItemNotFound, httpErr.WithDetail("no repair service found for id: "+id))
		return
	}

	wtr.WriteHeader(http.StatusOK)
	json.NewEncoder(wtr).Encode(svc)
}

// parsePaginationParams reads ?page=N&limit=N from the request, applying
// sensible defaults and clamping limit to maxRepairPageLimit.
func parsePaginationParams(req *http.Request) (page, limit int) {
	q := req.URL.Query()

	page = 1
	if p := q.Get("page"); p != "" {
		if n, err := strconv.Atoi(p); err == nil && n > 0 {
			page = n
		}
	}

	limit = defaultRepairPageLimit
	if l := q.Get("limit"); l != "" {
		if n, err := strconv.Atoi(l); err == nil && n > 0 {
			if n > maxRepairPageLimit {
				n = maxRepairPageLimit
			}
			limit = n
		}
	}
	return page, limit
}

// paginateSlice returns safe start/end indices for slicing a slice of length
// total at the given 1-based page and limit.
func paginateSlice(total, page, limit int) (start, end int) {
	start = (page - 1) * limit
	if start >= total {
		return total, total
	}
	end = min(start+limit, total)
	return start, end
}
