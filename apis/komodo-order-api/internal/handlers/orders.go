package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"komodo-order-api/internal/models"
	"komodo-order-api/internal/service"

	ctxKeys "github.com/rdevitto86/komodo-forge-sdk-go/http/context"
	httpErr "github.com/rdevitto86/komodo-forge-sdk-go/http/errors"
)

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
// Extracts userId from JWT context only — never from the URL.
// Returns 404 for both missing orders and orders owned by a different user,
// preventing callers from inferring orderId existence via status-code differences.
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
// Supports cursor-based pagination via ?limit=<n>&cursor=<token> query params.
// Extracts userId from JWT context; query is always scoped to the authenticated user.
func ListOrders(svc *service.OrderService) http.HandlerFunc {
	return func(wtr http.ResponseWriter, req *http.Request) {
		userID, ok := req.Context().Value(ctxKeys.USER_ID_KEY).(string)
		if !ok || userID == "" {
			httpErr.SendError(wtr, req, httpErr.Global.Unauthorized)
			return
		}

		q := req.URL.Query()
		limit := 0
		if s := q.Get("limit"); s != "" {
			n, err := strconv.Atoi(s)
			if err != nil || n < 0 {
				httpErr.SendError(wtr, req, httpErr.Global.BadRequest, httpErr.WithDetail("limit must be a positive integer"))
				return
			}
			limit = n
		}
		cursor := q.Get("cursor")

		orders, nextCursor, err := svc.ListOrders(req.Context(), userID, limit, cursor)
		if err != nil {
			sendOrderError(wtr, req, err)
			return
		}

		resp := models.OrderListResponse{
			Orders:     orders,
			NextCursor: nextCursor,
		}

		wtr.Header().Set("Content-Type", "application/json")
		wtr.WriteHeader(http.StatusOK)
		json.NewEncoder(wtr).Encode(resp)
	}
}

// GetOrderInternal handles GET /internal/orders/{orderId}.
// No user ownership check — protected by scope-checked JWT on the private middleware stack.
// Used by internal services: payments-api, returns-api.
func GetOrderInternal(svc *service.OrderService) http.HandlerFunc {
	return func(wtr http.ResponseWriter, req *http.Request) {
		orderID := req.PathValue("orderId")
		if orderID == "" {
			httpErr.SendError(wtr, req, httpErr.Global.BadRequest, httpErr.WithDetail("orderId path parameter is required"))
			return
		}

		order, err := svc.GetOrderInternal(req.Context(), orderID)
		if err != nil {
			sendOrderError(wtr, req, err)
			return
		}

		wtr.Header().Set("Content-Type", "application/json")
		wtr.WriteHeader(http.StatusOK)
		json.NewEncoder(wtr).Encode(order)
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
func CancelOrder(svc *service.OrderService) http.HandlerFunc {
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

		order, err := svc.CancelOrder(req.Context(), userID, orderID)
		if err != nil {
			sendOrderError(wtr, req, err)
			return
		}

		wtr.Header().Set("Content-Type", "application/json")
		wtr.WriteHeader(http.StatusOK)
		json.NewEncoder(wtr).Encode(order)
	}
}

// sendOrderError maps domain errors to RFC 7807 responses.
func sendOrderError(wtr http.ResponseWriter, req *http.Request, err error) {
	switch {
	case errors.Is(err, models.ErrNotFound):
		httpErr.SendError(wtr, req, models.Err.NotFound)
	case errors.Is(err, models.ErrForbidden):
		httpErr.SendError(wtr, req, httpErr.Global.Forbidden)
	case errors.Is(err, models.ErrAlreadyCancelled):
		httpErr.SendError(wtr, req, models.Err.AlreadyCancelled)
	case errors.Is(err, models.ErrNotCancellable):
		httpErr.SendError(wtr, req, models.Err.NotCancellable)
	case errors.Is(err, models.ErrInvalidTransition):
		httpErr.SendError(wtr, req, models.Err.InvalidTransition)
	default:
		httpErr.SendError(wtr, req, httpErr.Global.Internal)
	}
}
