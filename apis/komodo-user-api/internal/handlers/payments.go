package handlers

import (
	httpErr "github.com/rdevitto86/komodo-forge-sdk-go/http/errors"
	"net/http"
)

// GetPayments returns all saved payment methods for the authenticated user.
func GetPayments(wtr http.ResponseWriter, req *http.Request) {
	httpErr.SendError(wtr, req, httpErr.Global.NotImplemented, httpErr.WithDetail("GetPayments not yet implemented"))
}

// UpsertPayment adds or updates a payment method for the authenticated user.
func UpsertPayment(wtr http.ResponseWriter, req *http.Request) {
	httpErr.SendError(wtr, req, httpErr.Global.NotImplemented, httpErr.WithDetail("UpsertPayment not yet implemented"))
}

// DeletePayment removes a payment method by ID for the authenticated user.
func DeletePayment(wtr http.ResponseWriter, req *http.Request) {
	httpErr.SendError(wtr, req, httpErr.Global.NotImplemented, httpErr.WithDetail("DeletePayment not yet implemented"))
}
