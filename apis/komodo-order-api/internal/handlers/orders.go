package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"komodo-order-api/internal/models"
	"komodo-order-api/internal/service"

	ctxKeys "github.com/rdevitto86/komodo-forge-sdk-go/http/context"
	httpErr "github.com/rdevitto86/komodo-forge-sdk-go/http/errors"
)

// PlaceOrder handles POST /me/orders.
// It consumes a checkout token issued by cart-api to create a new order.
func PlaceOrder(svc *service.OrderService) http.HandlerFunc {
	return func(wtr http.ResponseWriter, req *http.Request) {
		userID, ok := req.Context().Value(ctxKeys.USER_ID_KEY).(string)
		if !ok || userID == "" {
			httpErr.SendError(wtr, req, httpErr.Global.Unauthorized)
			return
		}

		var body models.PlaceOrderRequest
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			httpErr.SendError(wtr, req, httpErr.Global.BadRequest, httpErr.WithDetail("invalid request body"))
			return
		}
		if body.CheckoutToken == "" {
			httpErr.SendError(wtr, req, httpErr.Global.BadRequest, httpErr.WithDetail("checkoutToken is required"))
			return
		}

		order, err := svc.PlaceOrder(req.Context(), userID, body.CheckoutToken)
		if err != nil {
			sendOrderError(wtr, req, err)
			return
		}

		wtr.Header().Set("Content-Type", "application/json")
		wtr.WriteHeader(http.StatusCreated)
		json.NewEncoder(wtr).Encode(order)
	}
}

// GetOrder handles GET /me/orders/{orderId}.
func GetOrder(svc *service.OrderService) http.HandlerFunc {
	return func(wtr http.ResponseWriter, req *http.Request) {
		userID, ok := req.Context().Value(ctxKeys.USER_ID_KEY).(string)
		if !ok || userID == "" {
			httpErr.SendError(wtr, req, httpErr.Global.Unauthorized)
			return
		}

		orderID := req.PathValue("orderId")
		if orderID == "" {
			httpErr.SendError(wtr, req, httpErr.Global.BadRequest, httpErr.WithDetail("orderId path parameter is required"))
			return
		}

		order, err := svc.GetOrder(req.Context(), userID, orderID)
		if err != nil {
			sendOrderError(wtr, req, err)
			return
		}

		wtr.Header().Set("Content-Type", "application/json")
		wtr.WriteHeader(http.StatusOK)
		json.NewEncoder(wtr).Encode(order)
	}
}

// ListOrders handles GET /me/orders.
func ListOrders(svc *service.OrderService) http.HandlerFunc {
	return func(wtr http.ResponseWriter, req *http.Request) {
		userID, ok := req.Context().Value(ctxKeys.USER_ID_KEY).(string)
		if !ok || userID == "" {
			httpErr.SendError(wtr, req, httpErr.Global.Unauthorized)
			return
		}

		orders, err := svc.ListOrders(req.Context(), userID)
		if err != nil {
			sendOrderError(wtr, req, err)
			return
		}

		wtr.Header().Set("Content-Type", "application/json")
		wtr.WriteHeader(http.StatusOK)
		json.NewEncoder(wtr).Encode(orders)
	}
}

// GetOrderUnified handles GET /orders/{orderId}.
// Supports both authenticated (JWT) and guest (email query param) access.
// If a JWT is present the userID is extracted and used for ownership validation.
// If no JWT is present the ?email query param is required and validated against
// the email stored on the order. In both cases a missing or mismatched identity
// results in 404 to prevent order ID enumeration.
func GetOrderUnified(svc *service.OrderService) http.HandlerFunc {
	return func(wtr http.ResponseWriter, req *http.Request) {
		// userID is optional — may be empty string when no JWT is provided.
		userID, _ := req.Context().Value(ctxKeys.USER_ID_KEY).(string)

		orderID := req.PathValue("orderId")
		if orderID == "" {
			httpErr.SendError(wtr, req, httpErr.Global.BadRequest, httpErr.WithDetail("orderId path parameter is required"))
			return
		}

		var email string
		if userID == "" {
			email = req.URL.Query().Get("email")
			if email == "" {
				httpErr.SendError(wtr, req, httpErr.Global.BadRequest, httpErr.WithDetail("email query parameter is required for unauthenticated requests"))
				return
			}
		}

		order, err := svc.GetOrderUnified(req.Context(), userID, email, orderID)
		if err != nil {
			sendOrderError(wtr, req, err)
			return
		}

		wtr.Header().Set("Content-Type", "application/json")
		wtr.WriteHeader(http.StatusOK)
		json.NewEncoder(wtr).Encode(order)
	}
}

// CancelOrder handles POST /me/orders/{orderId}/cancel.
// Stubbed — cancellation logic (status transition validation, refund trigger) not yet implemented.
func CancelOrder(_ *service.OrderService) http.HandlerFunc {
	return func(wtr http.ResponseWriter, req *http.Request) {
		httpErr.SendError(wtr, req, httpErr.Global.NotImplemented)
	}
}

// sendOrderError maps domain errors to RFC 7807 responses.
func sendOrderError(wtr http.ResponseWriter, req *http.Request, err error) {
	switch {
	case errors.Is(err, models.ErrNotFound):
		httpErr.SendError(wtr, req, models.Err.NotFound)
	case errors.Is(err, models.ErrForbidden):
		httpErr.SendError(wtr, req, httpErr.Global.Forbidden)
	default:
		httpErr.SendError(wtr, req, httpErr.Global.Internal)
	}
}
