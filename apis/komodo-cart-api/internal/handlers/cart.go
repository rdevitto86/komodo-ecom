package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"komodo-cart-api/internal/models"
	"komodo-cart-api/internal/service"

	ctxKeys "github.com/rdevitto86/komodo-forge-sdk-go/http/context"
	httpErr "github.com/rdevitto86/komodo-forge-sdk-go/http/errors"
)

// GetMyCart returns the authenticated user's cart.
// Accepts optional X-Guest-Cart-ID header to trigger a merge on first load after login.
func GetMyCart(svc *service.CartService) http.HandlerFunc {
	return func(wtr http.ResponseWriter, req *http.Request) {
		userID, ok := req.Context().Value(ctxKeys.USER_ID_KEY).(string)
		if !ok || userID == "" {
			httpErr.SendError(wtr, req, httpErr.Global.Unauthorized)
			return
		}

		guestCartID := req.Header.Get("X-Guest-Cart-ID")

		cart, err := svc.Get(req.Context(), userID, guestCartID)
		if err != nil {
			sendCartError(wtr, req, err)
			return
		}

		wtr.Header().Set("Content-Type", "application/json")
		wtr.WriteHeader(http.StatusOK)
		json.NewEncoder(wtr).Encode(cart)
	}
}

// MergeGuestCart merges a guest cart into the authenticated user's cart.
func MergeGuestCart(svc *service.CartService) http.HandlerFunc {
	return func(wtr http.ResponseWriter, req *http.Request) {
		userID, ok := req.Context().Value(ctxKeys.USER_ID_KEY).(string)
		if !ok || userID == "" {
			httpErr.SendError(wtr, req, httpErr.Global.Unauthorized)
			return
		}

		var body models.MergeCartRequest
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			httpErr.SendError(wtr, req, httpErr.Global.BadRequest, httpErr.WithDetail("invalid request body"))
			return
		}
		if body.GuestCartID == "" {
			httpErr.SendError(wtr, req, httpErr.Global.BadRequest, httpErr.WithDetail("guest_cart_id is required"))
			return
		}

		cart, err := svc.Get(req.Context(), userID, body.GuestCartID)
		if err != nil {
			sendCartError(wtr, req, err)
			return
		}

		wtr.Header().Set("Content-Type", "application/json")
		wtr.WriteHeader(http.StatusOK)
		json.NewEncoder(wtr).Encode(cart)
	}
}

// AddMyCartItem adds an item to the authenticated user's cart.
func AddMyCartItem(svc *service.CartService) http.HandlerFunc {
	return func(wtr http.ResponseWriter, req *http.Request) {
		userID, ok := req.Context().Value(ctxKeys.USER_ID_KEY).(string)
		if !ok || userID == "" {
			httpErr.SendError(wtr, req, httpErr.Global.Unauthorized)
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

		cart, err := svc.AddItem(req.Context(), userID, body)
		if err != nil {
			sendCartError(wtr, req, err)
			return
		}

		wtr.Header().Set("Content-Type", "application/json")
		wtr.WriteHeader(http.StatusCreated)
		json.NewEncoder(wtr).Encode(cart)
	}
}

// UpdateMyCartItem updates the quantity of an item in the authenticated user's cart.
func UpdateMyCartItem(svc *service.CartService) http.HandlerFunc {
	return func(wtr http.ResponseWriter, req *http.Request) {
		userID, ok := req.Context().Value(ctxKeys.USER_ID_KEY).(string)
		if !ok || userID == "" {
			httpErr.SendError(wtr, req, httpErr.Global.Unauthorized)
			return
		}

		itemID := req.PathValue("itemId")
		if itemID == "" {
			httpErr.SendError(wtr, req, httpErr.Global.BadRequest, httpErr.WithDetail("itemId path parameter required"))
			return
		}

		var body models.UpdateItemRequest
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			httpErr.SendError(wtr, req, httpErr.Global.BadRequest, httpErr.WithDetail("invalid request body"))
			return
		}

		cart, err := svc.UpdateItem(req.Context(), userID, itemID, body)
		if err != nil {
			sendCartError(wtr, req, err)
			return
		}

		wtr.Header().Set("Content-Type", "application/json")
		wtr.WriteHeader(http.StatusOK)
		json.NewEncoder(wtr).Encode(cart)
	}
}

// RemoveMyCartItem removes an item from the authenticated user's cart.
func RemoveMyCartItem(svc *service.CartService) http.HandlerFunc {
	return func(wtr http.ResponseWriter, req *http.Request) {
		userID, ok := req.Context().Value(ctxKeys.USER_ID_KEY).(string)
		if !ok || userID == "" {
			httpErr.SendError(wtr, req, httpErr.Global.Unauthorized)
			return
		}

		itemID := req.PathValue("itemId")
		if itemID == "" {
			httpErr.SendError(wtr, req, httpErr.Global.BadRequest, httpErr.WithDetail("itemId path parameter required"))
			return
		}

		if err := svc.RemoveItem(req.Context(), userID, itemID); err != nil {
			sendCartError(wtr, req, err)
			return
		}

		wtr.WriteHeader(http.StatusNoContent)
	}
}

// ClearMyCart removes all items from the authenticated user's cart.
func ClearMyCart(svc *service.CartService) http.HandlerFunc {
	return func(wtr http.ResponseWriter, req *http.Request) {
		userID, ok := req.Context().Value(ctxKeys.USER_ID_KEY).(string)
		if !ok || userID == "" {
			httpErr.SendError(wtr, req, httpErr.Global.Unauthorized)
			return
		}

		if err := svc.Clear(req.Context(), userID); err != nil {
			sendCartError(wtr, req, err)
			return
		}

		wtr.WriteHeader(http.StatusNoContent)
	}
}

// InitiateCheckout places stock holds via inventory-api and returns a checkout token.
func InitiateCheckout(svc *service.CartService) http.HandlerFunc {
	return func(wtr http.ResponseWriter, req *http.Request) {
		userID, ok := req.Context().Value(ctxKeys.USER_ID_KEY).(string)
		if !ok || userID == "" {
			httpErr.SendError(wtr, req, httpErr.Global.Unauthorized)
			return
		}

		res, err := svc.Checkout(req.Context(), userID)
		if err != nil {
			sendCartError(wtr, req, err)
			return
		}

		wtr.Header().Set("Content-Type", "application/json")
		wtr.WriteHeader(http.StatusOK)
		json.NewEncoder(wtr).Encode(res)
	}
}

// sendCartError maps service errors to HTTP responses.
// All domain errors from the service layer are models.CartError values.
// errors.As walks the chain, so wrapped CartErrors are found correctly.
func sendCartError(wtr http.ResponseWriter, req *http.Request, err error) {
	var cartErr models.CartError
	if errors.As(err, &cartErr) {
		httpErr.SendError(wtr, req, cartErr.Code)
		return
	}
	httpErr.SendError(wtr, req, httpErr.Global.Internal, httpErr.WithDetail(err.Error()))
}
