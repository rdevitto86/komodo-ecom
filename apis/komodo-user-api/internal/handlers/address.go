package handlers

import (
	"encoding/json"
	"net/http"

	httpErr "github.com/rdevitto86/komodo-forge-sdk-go/http/errors"

	"komodo-user-api/internal/service"
	"komodo-user-api/pkg/v1/models"
)

// GetAddresses returns all addresses for the authenticated user.
func GetAddresses(wtr http.ResponseWriter, req *http.Request) {
	userID := resolveUserID(req)
	if userID == "" {
		httpErr.SendError(wtr, req, httpErr.Global.Unauthorized)
		return
	}

	addrs, err := service.GetAddresses(req.Context(), userID)
	if err != nil {
		sendUserError(wtr, req, err)
		return
	}

	wtr.Header().Set("Content-Type", "application/json")
	wtr.WriteHeader(http.StatusOK)
	json.NewEncoder(wtr).Encode(addrs)
}

// AddAddress adds a new address for the authenticated user.
func AddAddress(wtr http.ResponseWriter, req *http.Request) {
	userID := resolveUserID(req)
	if userID == "" {
		httpErr.SendError(wtr, req, httpErr.Global.Unauthorized)
		return
	}

	var input models.Address
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		httpErr.SendError(wtr, req, httpErr.Global.BadRequest)
		return
	}

	if err := service.AddAddress(req.Context(), userID, &input); err != nil {
		sendUserError(wtr, req, err)
		return
	}

	wtr.Header().Set("Content-Type", "application/json")
	wtr.WriteHeader(http.StatusCreated)
	json.NewEncoder(wtr).Encode(input)
}

// UpdateAddress updates an address by ID for the authenticated user.
func UpdateAddress(wtr http.ResponseWriter, req *http.Request) {
	userID := resolveUserID(req)
	if userID == "" {
		httpErr.SendError(wtr, req, httpErr.Global.Unauthorized)
		return
	}

	addressID := req.PathValue("id")
	if addressID == "" {
		httpErr.SendError(wtr, req, httpErr.Global.BadRequest)
		return
	}

	var input models.Address
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		httpErr.SendError(wtr, req, httpErr.Global.BadRequest)
		return
	}

	if err := service.UpdateAddress(req.Context(), userID, addressID, &input); err != nil {
		sendUserError(wtr, req, err)
		return
	}

	wtr.WriteHeader(http.StatusOK)
}

// DeleteAddress removes an address by ID for the authenticated user.
func DeleteAddress(wtr http.ResponseWriter, req *http.Request) {
	userID := resolveUserID(req)
	if userID == "" {
		httpErr.SendError(wtr, req, httpErr.Global.Unauthorized)
		return
	}

	addressID := req.PathValue("id")
	if addressID == "" {
		httpErr.SendError(wtr, req, httpErr.Global.BadRequest)
		return
	}

	if err := service.DeleteAddress(req.Context(), userID, addressID); err != nil {
		sendUserError(wtr, req, err)
		return
	}

	wtr.WriteHeader(http.StatusNoContent)
}
