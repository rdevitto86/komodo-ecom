package handlers

import (
	httpErr "github.com/rdevitto86/komodo-forge-sdk-go/http/errors"
	"net/http"
)

// GetPreferences returns preferences for the authenticated user.
func GetPreferences(wtr http.ResponseWriter, req *http.Request) {
	httpErr.SendError(wtr, req, httpErr.Global.NotImplemented, httpErr.WithDetail("GetPreferences not yet implemented"))
}

// UpdatePreferences updates preferences for the authenticated user.
func UpdatePreferences(wtr http.ResponseWriter, req *http.Request) {
	httpErr.SendError(wtr, req, httpErr.Global.NotImplemented, httpErr.WithDetail("UpdatePreferences not yet implemented"))
}

// DeletePreferences removes preferences for the authenticated user.
func DeletePreferences(wtr http.ResponseWriter, req *http.Request) {
	httpErr.SendError(wtr, req, httpErr.Global.NotImplemented, httpErr.WithDetail("DeletePreferences not yet implemented"))
}
