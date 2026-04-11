package handlers

import (
	"encoding/json"
	"net/http"

	httpErr "github.com/rdevitto86/komodo-forge-sdk-go/http/errors"

	"komodo-user-api/internal/service"
	"komodo-user-api/internal/models"
)

// GetPayments returns all saved payment methods for the authenticated user.
func GetPayments(wtr http.ResponseWriter, req *http.Request) {
	userID := resolveUserID(req)
	if userID == "" {
		httpErr.SendError(wtr, req, httpErr.Global.Unauthorized)
		return
	}

	payments, err := service.GetPayments(req.Context(), userID)
	if err != nil {
		sendUserError(wtr, req, err)
		return
	}

	wtr.Header().Set("Content-Type", "application/json")
	wtr.WriteHeader(http.StatusOK)
	json.NewEncoder(wtr).Encode(payments)
}

// UpsertPayment adds or updates a payment method for the authenticated user.
func UpsertPayment(wtr http.ResponseWriter, req *http.Request) {
	userID := resolveUserID(req)
	if userID == "" {
		httpErr.SendError(wtr, req, httpErr.Global.Unauthorized)
		return
	}

	var input models.PaymentMethod
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		httpErr.SendError(wtr, req, httpErr.Global.BadRequest)
		return
	}

	if err := service.UpsertPayment(req.Context(), userID, &input); err != nil {
		sendUserError(wtr, req, err)
		return
	}

	wtr.Header().Set("Content-Type", "application/json")
	wtr.WriteHeader(http.StatusOK)
	json.NewEncoder(wtr).Encode(input)
}

// DeletePayment removes a payment method by ID for the authenticated user.
func DeletePayment(wtr http.ResponseWriter, req *http.Request) {
	userID := resolveUserID(req)
	if userID == "" {
		httpErr.SendError(wtr, req, httpErr.Global.Unauthorized)
		return
	}

	paymentID := req.PathValue("id")
	if paymentID == "" {
		httpErr.SendError(wtr, req, httpErr.Global.BadRequest)
		return
	}

	if err := service.DeletePayment(req.Context(), userID, paymentID); err != nil {
		sendUserError(wtr, req, err)
		return
	}

	wtr.WriteHeader(http.StatusNoContent)
}
