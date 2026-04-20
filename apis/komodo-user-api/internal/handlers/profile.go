package handlers

import (
	"encoding/json"
	"net/http"

	ctxKeys "github.com/rdevitto86/komodo-forge-sdk-go/http/context"
	httpErr "github.com/rdevitto86/komodo-forge-sdk-go/http/errors"

	"komodo-user-api/internal/models"
	"komodo-user-api/internal/service"
)

// resolveUserID returns the user ID from the request, preferring the {id} path
// parameter (internal routes: /users/{id}) over the JWT subject in context
// (public routes: /me/profile). Returns empty string if neither is set.
func resolveUserID(req *http.Request) string {
	if id := req.PathValue("id"); id != "" {
		return id
	}
	id, _ := req.Context().Value(ctxKeys.USER_ID_KEY).(string)
	return id
}

// GetProfile returns the authenticated user's profile.
func GetProfile(wtr http.ResponseWriter, req *http.Request) {
	userID := resolveUserID(req)
	if userID == "" {
		httpErr.SendError(wtr, req, httpErr.Global.Unauthorized)
		return
	}

	user, err := service.GetProfile(req.Context(), userID)
	if err != nil {
		sendUserError(wtr, req, err)
		return
	}

	wtr.Header().Set("Content-Type", "application/json")
	wtr.WriteHeader(http.StatusOK)
	json.NewEncoder(wtr).Encode(user)
}

// CreateUser creates a new user record on registration.
func CreateUser(wtr http.ResponseWriter, req *http.Request) {
	userID := resolveUserID(req)
	if userID == "" {
		httpErr.SendError(wtr, req, httpErr.Global.Unauthorized)
		return
	}

	var input models.User
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		httpErr.SendError(wtr, req, httpErr.Global.BadRequest)
		return
	}
	// Ensure the user ID from the token is authoritative — do not let the caller set it.
	input.UserID = userID

	if err := service.CreateUser(req.Context(), &input); err != nil {
		sendUserError(wtr, req, err)
		return
	}

	wtr.Header().Set("Content-Type", "application/json")
	wtr.WriteHeader(http.StatusCreated)
	json.NewEncoder(wtr).Encode(input)
}

// UpdateProfile updates the authenticated user's profile.
func UpdateProfile(wtr http.ResponseWriter, req *http.Request) {
	userID := resolveUserID(req)
	if userID == "" {
		httpErr.SendError(wtr, req, httpErr.Global.Unauthorized)
		return
	}

	var input models.User
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		httpErr.SendError(wtr, req, httpErr.Global.BadRequest)
		return
	}

	updated, err := service.UpdateProfile(req.Context(), userID, &input)
	if err != nil {
		sendUserError(wtr, req, err)
		return
	}

	wtr.Header().Set("Content-Type", "application/json")
	wtr.WriteHeader(http.StatusOK)
	json.NewEncoder(wtr).Encode(updated)
}

// DeleteProfile deletes the authenticated user's account.
func DeleteProfile(wtr http.ResponseWriter, req *http.Request) {
	userID := resolveUserID(req)
	if userID == "" {
		httpErr.SendError(wtr, req, httpErr.Global.Unauthorized)
		return
	}

	if err := service.DeleteProfile(req.Context(), userID); err != nil {
		sendUserError(wtr, req, err)
		return
	}

	wtr.WriteHeader(http.StatusNoContent)
}
