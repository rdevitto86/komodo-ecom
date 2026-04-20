package handlers

import (
	"testing"

	"komodo-shop-items-api/internal/models"
)

// --- primarySKU ---

func TestPrimarySKU_DefaultVariant(t *testing.T) {
	p := models.Product{
		ID: "prod-1",
		Variants: []models.Variant{
			{SKU: "sku-a", IsDefault: false},
			{SKU: "sku-b", IsDefault: true},
		},
	}
	if got := primarySKU(p); got != "sku-b" {
		t.Errorf("want sku-b, got %s", got)
	}
}

func TestPrimarySKU_FirstVariantFallback(t *testing.T) {
	p := models.Product{
		ID: "prod-1",
		Variants: []models.Variant{
			{SKU: "sku-a", IsDefault: false},
		},
	}
	if got := primarySKU(p); got != "sku-a" {
		t.Errorf("want sku-a, got %s", got)
	}
}

func TestPrimarySKU_ProductIDFallback(t *testing.T) {
	p := models.Product{
		ID:       "prod-fallback",
		Variants: []models.Variant{},
	}
	if got := primarySKU(p); got != "prod-fallback" {
		t.Errorf("want prod-fallback, got %s", got)
	}
}

func TestPrimarySKU_VariantWithEmptySKUFallsToID(t *testing.T) {
	p := models.Product{
		ID:       "prod-x",
		Variants: []models.Variant{{SKU: "", IsDefault: true}},
	}
	if got := primarySKU(p); got != "prod-x" {
		t.Errorf("want prod-x, got %s", got)
	}
}

// --- isOutOfStock ---

func TestIsOutOfStock(t *testing.T) {
	tests := []struct {
		code models.StockCode
		want bool
	}{
		{models.StockOutOfStock, true},
		{models.StockSoldOut, true},
		{models.StockDiscontinued, true},
		{models.StockInStock, false},
		{models.StockLimitedSupply, false},
		{models.StockPreOrder, false},
		{models.StockBackorder, false},
		{models.StockTemporarilyUnavail, false},
	}
	for _, tc := range tests {
		if got := isOutOfStock(tc.code); got != tc.want {
			t.Errorf("isOutOfStock(%q): want %v, got %v", tc.code, tc.want, got)
		}
	}
}

// --- applyExcludeFilter ---

func TestApplyExcludeFilter_LimitRespected(t *testing.T) {
	products := []models.Product{
		{ID: "a", Variants: []models.Variant{{SKU: "sku-a", IsDefault: true}}},
		{ID: "b", Variants: []models.Variant{{SKU: "sku-b", IsDefault: true}}},
		{ID: "c", Variants: []models.Variant{{SKU: "sku-c", IsDefault: true}}},
	}
	result := applyExcludeFilter(products, map[string]struct{}{}, 2)
	if len(result) != 2 {
		t.Errorf("want 2 results, got %d", len(result))
	}
}

func TestApplyExcludeFilter_ExcludedSKUsOmitted(t *testing.T) {
	products := []models.Product{
		{ID: "a", Variants: []models.Variant{{SKU: "sku-a", IsDefault: true}}},
		{ID: "b", Variants: []models.Variant{{SKU: "sku-b", IsDefault: true}}},
		{ID: "c", Variants: []models.Variant{{SKU: "sku-c", IsDefault: true}}},
	}
	excluded := map[string]struct{}{"sku-b": {}}
	result := applyExcludeFilter(products, excluded, 10)
	if len(result) != 2 {
		t.Errorf("want 2 results after exclusion, got %d", len(result))
	}
	for _, p := range result {
		if primarySKU(p) == "sku-b" {
			t.Error("sku-b should have been excluded")
		}
	}
}

func TestApplyExcludeFilter_EmptyProductsReturnsEmpty(t *testing.T) {
	result := applyExcludeFilter([]models.Product{}, map[string]struct{}{}, 10)
	if len(result) != 0 {
		t.Errorf("want empty, got %d", len(result))
	}
}

func TestApplyExcludeFilter_AllExcludedReturnsEmpty(t *testing.T) {
	products := []models.Product{
		{ID: "a", Variants: []models.Variant{{SKU: "sku-a", IsDefault: true}}},
	}
	excluded := map[string]struct{}{"sku-a": {}}
	result := applyExcludeFilter(products, excluded, 10)
	if len(result) != 0 {
		t.Errorf("want empty after all excluded, got %d", len(result))
	}
}

// --- ranking helpers (tested via table-driven scenarios) ---

// rankByStockQty exercises the core sort logic used in GetSuggestions.
// Products with lower stockQty (higher sell-through) should rank first.
// Products without inventory data should sort last.
func TestRankingLowStockFirst(t *testing.T) {
	type scored struct {
		sku      string
		stockQty int
		hasInv   bool
	}

	candidates := []scored{
		{sku: "high-stock", stockQty: 100, hasInv: true},
		{sku: "low-stock", stockQty: 5, hasInv: true},
		{sku: "no-inv", stockQty: 0, hasInv: false},
		{sku: "mid-stock", stockQty: 40, hasInv: true},
	}

	// Mirror the sort logic from GetSuggestions.
	less := func(i, j int) bool {
		ci, cj := candidates[i], candidates[j]
		if ci.hasInv != cj.hasInv {
			return ci.hasInv
		}
		if !ci.hasInv {
			return false
		}
		return ci.stockQty < cj.stockQty
	}

	n := len(candidates)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if less(j+1, j) {
				candidates[j], candidates[j+1] = candidates[j+1], candidates[j]
			}
		}
	}

	want := []string{"low-stock", "mid-stock", "high-stock", "no-inv"}
	for i, c := range candidates {
		if c.sku != want[i] {
			t.Errorf("rank[%d]: want %s, got %s", i, want[i], c.sku)
		}
	}
}
