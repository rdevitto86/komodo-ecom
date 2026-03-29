package handlers

import (
	"encoding/json"
	"net/http"

	httpErr "github.com/rdevitto86/komodo-forge-sdk-go/http/errors"
	logger "github.com/rdevitto86/komodo-forge-sdk-go/logging/runtime"

	"komodo-search-api/pkg/v1/models"
)

// Search handles GET /search
//
// Query params:
//   - q        (required) — search query string
//   - type     (optional) — filter by item type: "product" | "service"
//   - category (optional) — filter by category slug
//   - page     (optional) — page number (default 1)
//   - per_page (optional) — results per page (default 20, max 100)
//
// TODO(typesense): replace stub response with real Typesense query.
// Call repository.Search(ctx, params) which wraps the Typesense client.
func Search(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")

	q := req.URL.Query().Get("q")
	if q == "" {
		httpErr.SendError(wtr, req, models.Err.InvalidQuery, httpErr.WithDetail("q parameter is required"))
		return
	}

	// TODO(typesense): build search params from query string
	// params := repository.SearchParams{
	// 	Query:      q,
	// 	ItemType:   req.URL.Query().Get("type"),     // "product" | "service" | "" (all)
	// 	Category:   req.URL.Query().Get("category"),
	// 	Page:       parseIntParam(req, "page", 1),
	// 	PerPage:    parseIntParam(req, "per_page", 20),
	// }

	// TODO(typesense): call repository.Search(req.Context(), params)
	// results, err := repository.Search(req.Context(), params)
	// if err != nil {
	// 	logger.Error("typesense search failed", err)
	// 	httpErr.SendError(wtr, req, models.Err.SearchFailed, httpErr.WithDetail(err.Error()))
	// 	return
	// }

	logger.Info("search query received: " + q)

	// Stub response — replace when Typesense is wired
	wtr.WriteHeader(http.StatusOK)
	json.NewEncoder(wtr).Encode(models.SearchResponse{
		Query:   q,
		Results: []models.SearchResult{},
		Total:   0,
		Page:    1,
		PerPage: 20,
	})
}
