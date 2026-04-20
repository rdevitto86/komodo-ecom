package handlers

import (
	"net/http"
)

// Handles GET /stats/items/{itemId}/recently-bought
func RecentlyBoughtHandler(wtr http.ResponseWriter, req *http.Request) {
	// TODO: implement
}

// Handles GET /stats/items/{itemId}/in-cart
func InCartHandler(wtr http.ResponseWriter, req *http.Request) {
	// TODO: implement
}

// Handles GET /stats/items/{itemId}/frequently-bought-with
func FrequentlyBoughtWithHandler(wtr http.ResponseWriter, req *http.Request) {
	// TODO: implement
}

// Handles GET /stats/items/{itemId}
func ItemStatsHandler(wtr http.ResponseWriter, req *http.Request) {
	// TODO: implement
}
