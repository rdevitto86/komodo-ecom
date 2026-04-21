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

// PlaceOrderUnified handles POST /orders.
// Accepts both authenticated (JWT) and guest callers. When a JWT is present,
// the userID from context is used and any email in the request body is ignored.
// When no JWT is present, the email field is required and is used to look up or
// create a guest identity at the service layer.
func PlaceOrderUnified(svc *service.OrderService) http.HandlerFunc {
	return func(wtr http.ResponseWriter, req *http.Request) {
		// userID may be empty for unauthenticated (guest) callers.
		userID, _ := req.Context().Value(ctxKeys.USER_ID_KEY).(string)

		var body models.UnifiedPlaceOrderRequest
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			httpErr.SendError(wtr, req, httpErr.Global.BadRequest, httpErr.WithDetail("invalid request body"))
			return
		}
		if body.CheckoutToken == "" {
			httpErr.SendError(wtr, req, httpErr.Global.BadRequest, httpErr.WithDetail("checkoutToken is required"))
			return
		}
		// Email is only required when there is no authenticated identity.
		if userID == "" && body.Email == "" {
			httpErr.SendError(wtr, req, httpErr.Global.BadRequest, httpErr.WithDetail("email is required for guest orders"))
			return
		}

		order, err := svc.PlaceOrderUnified(req.Context(), userID, body.Email, body.CheckoutToken)
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
	case errors.Is(err, models.ErrForbidden):
		httpErr.SendError(wtr, req, httpErr.Global.Forbidden)
	default:
		httpErr.SendError(wtr, req, httpErr.Global.Internal)
	}
}
