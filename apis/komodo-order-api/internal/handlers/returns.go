package handlers

import (
	"net/http"

	httpErr "github.com/rdevitto86/komodo-forge-sdk-go/http/errors"
)

// ListReturns handles GET /me/orders/returns.
// Stubbed — RMA list retrieval not yet implemented.
func ListReturns() http.HandlerFunc {
	return func(wtr http.ResponseWriter, req *http.Request) {
		httpErr.SendError(wtr, req, httpErr.Global.NotImplemented)
	}
}

// CreateReturn handles POST /me/orders/returns.
// Stubbed — RMA creation, order eligibility validation, and return window enforcement not yet implemented.
func CreateReturn() http.HandlerFunc {
	return func(wtr http.ResponseWriter, req *http.Request) {
		httpErr.SendError(wtr, req, httpErr.Global.NotImplemented)
	}
}

// GetReturn handles GET /me/orders/returns/{returnId}.
// Stubbed — RMA retrieval not yet implemented.
func GetReturn() http.HandlerFunc {
	return func(wtr http.ResponseWriter, req *http.Request) {
		httpErr.SendError(wtr, req, httpErr.Global.NotImplemented)
	}
}

// CancelReturn handles DELETE /me/orders/returns/{returnId}.
// Stubbed — RMA cancellation not yet implemented.
func CancelReturn() http.HandlerFunc {
	return func(wtr http.ResponseWriter, req *http.Request) {
		httpErr.SendError(wtr, req, httpErr.Global.NotImplemented)
	}
}

// GetReturnInternal handles GET /internal/returns/{returnId}.
// Internal route for service-to-service lookups (payments-api, etc.).
// Stubbed — not yet implemented.
func GetReturnInternal() http.HandlerFunc {
	return func(wtr http.ResponseWriter, req *http.Request) {
		httpErr.SendError(wtr, req, httpErr.Global.NotImplemented)
	}
}

// ApproveReturn handles PUT /internal/returns/{returnId}/approve.
// Triggers refund via payments-api on approval.
// Stubbed — not yet implemented.
func ApproveReturn() http.HandlerFunc {
	return func(wtr http.ResponseWriter, req *http.Request) {
		httpErr.SendError(wtr, req, httpErr.Global.NotImplemented)
	}
}

// ReceiveReturn handles PUT /internal/returns/{returnId}/receive.
// Triggers restock via shop-inventory-api and loyalty reversal via loyalty-api on receipt.
// Stubbed — not yet implemented.
func ReceiveReturn() http.HandlerFunc {
	return func(wtr http.ResponseWriter, req *http.Request) {
		httpErr.SendError(wtr, req, httpErr.Global.NotImplemented)
	}
}

// RejectReturn handles PUT /internal/returns/{returnId}/reject.
// Stubbed — not yet implemented.
func RejectReturn() http.HandlerFunc {
	return func(wtr http.ResponseWriter, req *http.Request) {
		httpErr.SendError(wtr, req, httpErr.Global.NotImplemented)
	}
}
