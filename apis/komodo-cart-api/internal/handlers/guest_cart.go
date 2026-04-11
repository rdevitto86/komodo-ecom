package handlers

import (
	"encoding/json"
	"net/http"

	httpErr "github.com/rdevitto86/komodo-forge-sdk-go/http/errors"
	"komodo-cart-api/internal/service"
	"komodo-cart-api/internal/models"
)

// CreateGuestCart creates a new guest cart and returns its cart ID and session token.
func CreateGuestCart(svc *service.GuestCartService) http.HandlerFunc {
	return func(wtr http.ResponseWriter, req *http.Request) {
		cart, sessionID, err := svc.Create(req.Context())
		if err != nil {
			sendCartError(wtr, req, err)
			return
		}

		wtr.Header().Set("Content-Type", "application/json")
		wtr.Header().Set("X-Session-ID", sessionID)
		wtr.WriteHeader(http.StatusCreated)
		json.NewEncoder(wtr).Encode(cart)
	}
}

// GetGuestCart returns the guest cart for the given cartId.
func GetGuestCart(svc *service.GuestCartService) http.HandlerFunc {
	return func(wtr http.ResponseWriter, req *http.Request) {
		cartID := req.PathValue("cartId")
		if cartID == "" {
			httpErr.SendError(wtr, req, httpErr.Global.BadRequest, httpErr.WithDetail("cartId path parameter required"))
			return
		}

		sessionID := req.Header.Get("X-Session-ID")
		if sessionID == "" {
			httpErr.SendError(wtr, req, httpErr.Global.BadRequest, httpErr.WithDetail("X-Session-ID header required"))
			return
		}

		cart, err := svc.Get(req.Context(), cartID, sessionID)
		if err != nil {
			sendCartError(wtr, req, err)
			return
		}

		wtr.Header().Set("Content-Type", "application/json")
		wtr.WriteHeader(http.StatusOK)
		json.NewEncoder(wtr).Encode(cart)
	}
}

// AddGuestCartItem adds an item to the specified guest cart.
func AddGuestCartItem(svc *service.GuestCartService) http.HandlerFunc {
	return func(wtr http.ResponseWriter, req *http.Request) {
		cartID := req.PathValue("cartId")
		if cartID == "" {
			httpErr.SendError(wtr, req, httpErr.Global.BadRequest, httpErr.WithDetail("cartId path parameter required"))
			return
		}

		sessionID := req.Header.Get("X-Session-ID")
		if sessionID == "" {
			httpErr.SendError(wtr, req, httpErr.Global.BadRequest, httpErr.WithDetail("X-Session-ID header required"))
			return
		}

		var body models.AddItemRequest
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			httpErr.SendError(wtr, req, httpErr.Global.BadRequest, httpErr.WithDetail("invalid request body"))
			return
		}
		if body.ItemID == "" || body.SKU == "" || body.Quantity < 1 {
			httpErr.SendError(wtr, req, httpErr.Global.BadRequest, httpErr.WithDetail("item_id, sku, and quantity >= 1 are required"))
			return
		}

		cart, err := svc.AddItem(req.Context(), cartID, sessionID, body)
		if err != nil {
			sendCartError(wtr, req, err)
			return
		}

		wtr.Header().Set("Content-Type", "application/json")
		wtr.WriteHeader(http.StatusCreated)
		json.NewEncoder(wtr).Encode(cart)
	}
}

// UpdateGuestCartItem updates the quantity of an item in the specified guest cart.
func UpdateGuestCartItem(svc *service.GuestCartService) http.HandlerFunc {
	return func(wtr http.ResponseWriter, req *http.Request) {
		cartID := req.PathValue("cartId")
		itemID := req.PathValue("itemId")
		if cartID == "" || itemID == "" {
			httpErr.SendError(wtr, req, httpErr.Global.BadRequest, httpErr.WithDetail("cartId and itemId path parameters required"))
			return
		}

		sessionID := req.Header.Get("X-Session-ID")
		if sessionID == "" {
			httpErr.SendError(wtr, req, httpErr.Global.BadRequest, httpErr.WithDetail("X-Session-ID header required"))
			return
		}

		var body models.UpdateItemRequest
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			httpErr.SendError(wtr, req, httpErr.Global.BadRequest, httpErr.WithDetail("invalid request body"))
			return
		}

		cart, err := svc.UpdateItem(req.Context(), cartID, sessionID, itemID, body)
		if err != nil {
			sendCartError(wtr, req, err)
			return
		}

		wtr.Header().Set("Content-Type", "application/json")
		wtr.WriteHeader(http.StatusOK)
		json.NewEncoder(wtr).Encode(cart)
	}
}

// RemoveGuestCartItem removes an item from the specified guest cart.
func RemoveGuestCartItem(svc *service.GuestCartService) http.HandlerFunc {
	return func(wtr http.ResponseWriter, req *http.Request) {
		cartID := req.PathValue("cartId")
		itemID := req.PathValue("itemId")
		if cartID == "" || itemID == "" {
			httpErr.SendError(wtr, req, httpErr.Global.BadRequest, httpErr.WithDetail("cartId and itemId path parameters required"))
			return
		}

		sessionID := req.Header.Get("X-Session-ID")
		if sessionID == "" {
			httpErr.SendError(wtr, req, httpErr.Global.BadRequest, httpErr.WithDetail("X-Session-ID header required"))
			return
		}

		if err := svc.RemoveItem(req.Context(), cartID, sessionID, itemID); err != nil {
			sendCartError(wtr, req, err)
			return
		}

		wtr.WriteHeader(http.StatusNoContent)
	}
}

// ClearGuestCart removes all items from the specified guest cart.
func ClearGuestCart(svc *service.GuestCartService) http.HandlerFunc {
	return func(wtr http.ResponseWriter, req *http.Request) {
		cartID := req.PathValue("cartId")
		if cartID == "" {
			httpErr.SendError(wtr, req, httpErr.Global.BadRequest, httpErr.WithDetail("cartId path parameter required"))
			return
		}

		sessionID := req.Header.Get("X-Session-ID")
		if sessionID == "" {
			httpErr.SendError(wtr, req, httpErr.Global.BadRequest, httpErr.WithDetail("X-Session-ID header required"))
			return
		}

		if err := svc.Clear(req.Context(), cartID, sessionID); err != nil {
			sendCartError(wtr, req, err)
			return
		}

		wtr.WriteHeader(http.StatusNoContent)
	}
}
