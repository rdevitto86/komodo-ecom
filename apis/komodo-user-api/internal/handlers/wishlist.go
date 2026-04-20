package handlers

import (
	"encoding/json"
	"net/http"

	httpErr "github.com/rdevitto86/komodo-forge-sdk-go/http/errors"

	"komodo-user-api/internal/models"
)

// GetWishlist returns all wishlist items for the authenticated user.
func GetWishlist(wtr http.ResponseWriter, req *http.Request) {
	userID := resolveUserID(req)
	if userID == "" {
		httpErr.SendError(wtr, req, httpErr.Global.Unauthorized)
		return
	}

	// TODO: wire service.GetWishlist once wishlist domain is implemented
	wtr.Header().Set("Content-Type", "application/json")
	wtr.WriteHeader(http.StatusNotImplemented)
}

// AddWishlistItem adds an item to the authenticated user's wishlist.
func AddWishlistItem(wtr http.ResponseWriter, req *http.Request) {
	userID := resolveUserID(req)
	if userID == "" {
		httpErr.SendError(wtr, req, httpErr.Global.Unauthorized)
		return
	}

	var input models.AddWishlistItemRequest
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		httpErr.SendError(wtr, req, httpErr.Global.BadRequest)
		return
	}

	// TODO: wire service.AddWishlistItem once wishlist domain is implemented
	wtr.WriteHeader(http.StatusNotImplemented)
}

// RemoveWishlistItem removes an item from the authenticated user's wishlist.
func RemoveWishlistItem(wtr http.ResponseWriter, req *http.Request) {
	userID := resolveUserID(req)
	if userID == "" {
		httpErr.SendError(wtr, req, httpErr.Global.Unauthorized)
		return
	}

	itemID := req.PathValue("itemId")
	if itemID == "" {
		httpErr.SendError(wtr, req, httpErr.Global.BadRequest)
		return
	}

	// TODO: wire service.RemoveWishlistItem once wishlist domain is implemented
	wtr.WriteHeader(http.StatusNotImplemented)
}

// GetWishlistAvailability returns stock availability for all items in the authenticated user's wishlist.
func GetWishlistAvailability(wtr http.ResponseWriter, req *http.Request) {
	userID := resolveUserID(req)
	if userID == "" {
		httpErr.SendError(wtr, req, httpErr.Global.Unauthorized)
		return
	}

	// TODO: wire service.GetWishlistAvailability once wishlist domain is implemented
	wtr.WriteHeader(http.StatusNotImplemented)
}

// MoveWishlistToCart moves one or more wishlist items into the cart.
func MoveWishlistToCart(wtr http.ResponseWriter, req *http.Request) {
	userID := resolveUserID(req)
	if userID == "" {
		httpErr.SendError(wtr, req, httpErr.Global.Unauthorized)
		return
	}

	var input models.MoveToCartRequest
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		httpErr.SendError(wtr, req, httpErr.Global.BadRequest)
		return
	}

	// TODO: wire service.MoveWishlistToCart once wishlist domain is implemented
	wtr.WriteHeader(http.StatusNotImplemented)
}
