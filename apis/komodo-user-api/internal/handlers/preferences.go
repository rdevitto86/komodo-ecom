package handlers

import (
	"encoding/json"
	"net/http"

	httpErr "github.com/rdevitto86/komodo-forge-sdk-go/http/errors"

	"komodo-user-api/internal/service"
	"komodo-user-api/internal/models"
)

// GetPreferences returns preferences for the authenticated user.
func GetPreferences(wtr http.ResponseWriter, req *http.Request) {
	userID := resolveUserID(req)
	if userID == "" {
		httpErr.SendError(wtr, req, httpErr.Global.Unauthorized)
		return
	}

	prefs, err := service.GetPreferences(req.Context(), userID)
	if err != nil {
		sendUserError(wtr, req, err)
		return
	}

	wtr.Header().Set("Content-Type", "application/json")
	wtr.WriteHeader(http.StatusOK)
	json.NewEncoder(wtr).Encode(prefs)
}

// UpdatePreferences updates preferences for the authenticated user.
func UpdatePreferences(wtr http.ResponseWriter, req *http.Request) {
	userID := resolveUserID(req)
	if userID == "" {
		httpErr.SendError(wtr, req, httpErr.Global.Unauthorized)
		return
	}

	var input models.Preferences
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		httpErr.SendError(wtr, req, httpErr.Global.BadRequest)
		return
	}

	if err := service.UpdatePreferences(req.Context(), userID, &input); err != nil {
		sendUserError(wtr, req, err)
		return
	}

	wtr.Header().Set("Content-Type", "application/json")
	wtr.WriteHeader(http.StatusOK)
	json.NewEncoder(wtr).Encode(input)
}

// DeletePreferences removes preferences for the authenticated user.
func DeletePreferences(wtr http.ResponseWriter, req *http.Request) {
	userID := resolveUserID(req)
	if userID == "" {
		httpErr.SendError(wtr, req, httpErr.Global.Unauthorized)
		return
	}

	if err := service.DeletePreferences(req.Context(), userID); err != nil {
		sendUserError(wtr, req, err)
		return
	}

	wtr.WriteHeader(http.StatusNoContent)
}
