package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"sort"

	"komodo-shop-items-api/internal/config"

	httpErr "github.com/rdevitto86/komodo-forge-sdk-go/http/errors"
	logger "github.com/rdevitto86/komodo-forge-sdk-go/logging/runtime"

	"komodo-shop-items-api/internal/models"
	"komodo-shop-items-api/internal/store"
)

const (
	defaultSuggestionLimit = 10
	maxSuggestionLimit     = 50
)

// GetSuggestions returns product suggestions ranked by inventory sell-through.
//
// Ranking strategy: lower remaining stock relative to peers signals higher
// demand ("low stock = popular"). Items that are fully out of stock are
// excluded. If inventory data is unavailable the handler degrades gracefully
// to catalog insertion order rather than returning an error.
//
// Route: POST /item/suggestion (auth required)
func GetSuggestions(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")

	var reqBody models.SuggestionRequest
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		logger.Error("failed to parse suggestion request body", err)
		httpErr.SendError(wtr, req, httpErr.Global.BadRequest, httpErr.WithDetail("failed to parse request body"))
		return
	}

	limit := reqBody.Limit
	if limit <= 0 {
		limit = defaultSuggestionLimit
	}
	if limit > maxSuggestionLimit {
		limit = maxSuggestionLimit
	}

	// Build a set of SKUs to exclude (items already in cart / recently viewed).
	excluded := make(map[string]struct{}, len(reqBody.Exclude)+len(reqBody.SKUs))
	for _, sku := range reqBody.Exclude {
		excluded[sku] = struct{}{}
	}

	bucket := os.Getenv(config.S3_ITEMS_BUCKET)
	if bucket == "" {
		logger.Error("S3_ITEMS_BUCKET not configured", nil)
		httpErr.SendError(wtr, req, models.Err.StorageError, httpErr.WithDetail("storage not configured"))
		return
	}

	ctx := req.Context()

	// Fetch all products. No products = empty suggestions, not an error.
	products, err := store.FetchAllProducts(ctx, bucket)
	if err != nil {
		logger.Error("failed to fetch products for suggestions", err)
		httpErr.SendError(wtr, req, models.Err.SuggestionFailed, httpErr.WithDetail("failed to retrieve catalog"))
		return
	}

	// Attempt to fetch inventory for sell-through ranking.
	// On failure we degrade gracefully: keep catalog order, skip OOS filter.
	inv, err := store.FetchInventory(ctx, bucket)
	if err != nil {
		logger.Warn("suggestion ranking unavailable, falling back to catalog order: " + err.Error())
		suggestions := applyExcludeFilter(products, excluded, limit)
		writeResponse(wtr, suggestions)
		return
	}

	// Build a lookup of SKU → stockQty from the inventory manifest.
	// Items missing from the manifest are treated as having unknown stock
	// and sorted to the back.
	invBySKU := make(map[string]models.InventoryItem, len(inv.Items))
	for _, item := range inv.Items {
		invBySKU[item.SKU] = item
	}

	// Filter and rank.
	type scored struct {
		product  models.Product
		stockQty int
		hasInv   bool
	}

	candidates := make([]scored, 0, len(products))
	for _, p := range products {
		// Determine the product's primary SKU from its first default (or first) variant.
		sku := primarySKU(p)

		if _, skip := excluded[sku]; skip {
			continue
		}

		invItem, ok := invBySKU[sku]
		if ok {
			// Exclude items that are fully out of stock.
			if isOutOfStock(invItem.StockCode) {
				continue
			}
			candidates = append(candidates, scored{product: p, stockQty: invItem.StockQty, hasInv: true})
		} else {
			// No inventory record — include but push to the back.
			candidates = append(candidates, scored{product: p, stockQty: 0, hasInv: false})
		}
	}

	// Rank: items with known inventory first, sorted ascending by stockQty
	// (lower remaining stock = more demand). Items without inventory data
	// come last in stable order.
	sort.SliceStable(candidates, func(i, j int) bool {
		ci, cj := candidates[i], candidates[j]
		if ci.hasInv != cj.hasInv {
			return ci.hasInv // known inventory sorts before unknown
		}
		if !ci.hasInv {
			return false // both unknown: preserve order
		}
		return ci.stockQty < cj.stockQty
	})

	// Collect up to limit products.
	suggestions := make([]models.Product, 0, limit)
	for _, c := range candidates {
		if len(suggestions) >= limit {
			break
		}
		suggestions = append(suggestions, c.product)
	}

	writeResponse(wtr, suggestions)
}

// primarySKU returns the SKU that identifies this product in the inventory
// manifest. It prefers the default variant's SKU, falls back to the first
// variant, and ultimately falls back to the product ID.
func primarySKU(p models.Product) string {
	for _, v := range p.Variants {
		if v.IsDefault && v.SKU != "" {
			return v.SKU
		}
	}
	if len(p.Variants) > 0 && p.Variants[0].SKU != "" {
		return p.Variants[0].SKU
	}
	return p.ID
}

// isOutOfStock returns true for stock codes that mean the item cannot be
// purchased (out of stock, sold out, discontinued).
func isOutOfStock(code models.StockCode) bool {
	switch code {
	case models.StockOutOfStock, models.StockSoldOut, models.StockDiscontinued:
		return true
	}
	return false
}

// applyExcludeFilter returns up to limit products after removing excluded SKUs.
// Used for the graceful-degradation path when inventory is unavailable.
func applyExcludeFilter(products []models.Product, excluded map[string]struct{}, limit int) []models.Product {
	result := make([]models.Product, 0, limit)
	for _, p := range products {
		if len(result) >= limit {
			break
		}
		sku := primarySKU(p)
		if _, skip := excluded[sku]; skip {
			continue
		}
		result = append(result, p)
	}
	return result
}

// writeResponse encodes the suggestion response and writes it to the wire.
func writeResponse(wtr http.ResponseWriter, suggestions []models.Product) {
	resp := models.SuggestionResponse{
		Suggestions: suggestions,
		Total:       len(suggestions),
	}
	wtr.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(wtr).Encode(resp); err != nil {
		logger.Error("failed to encode suggestion response", err)
	}
}
