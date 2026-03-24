package handlers

import (
	httpErr "github.com/rdevitto86/komodo-forge-sdk-go/http/errors"
	"net/http"
)

// GetProfile returns the authenticated user's profile.
func GetProfile(wtr http.ResponseWriter, req *http.Request) {
	httpErr.SendError(wtr, req, httpErr.Global.NotImplemented, httpErr.WithDetail("GetProfile not yet implemented"))
}

// CreateUser creates a new user record on registration.
func CreateUser(wtr http.ResponseWriter, req *http.Request) {
	httpErr.SendError(wtr, req, httpErr.Global.NotImplemented, httpErr.WithDetail("CreateUser not yet implemented"))
}

// UpdateProfile updates the authenticated user's profile.
func UpdateProfile(wtr http.ResponseWriter, req *http.Request) {
	httpErr.SendError(wtr, req, httpErr.Global.NotImplemented, httpErr.WithDetail("UpdateProfile not yet implemented"))
}

// DeleteProfile deletes the authenticated user's account.
func DeleteProfile(wtr http.ResponseWriter, req *http.Request) {
	httpErr.SendError(wtr, req, httpErr.Global.NotImplemented, httpErr.WithDetail("DeleteProfile not yet implemented"))
}
