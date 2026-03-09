package handlers

import (
	httpErr "komodo-forge-sdk-go/http/errors"
	"net/http"
)

// GetAddresses returns all addresses for the authenticated user.
func GetAddresses(wtr http.ResponseWriter, req *http.Request) {
	httpErr.SendError(wtr, req, httpErr.Global.NotImplemented, httpErr.WithDetail("GetAddresses not yet implemented"))
}

// AddAddress adds a new address for the authenticated user.
func AddAddress(wtr http.ResponseWriter, req *http.Request) {
	httpErr.SendError(wtr, req, httpErr.Global.NotImplemented, httpErr.WithDetail("AddAddress not yet implemented"))
}

// UpdateAddress updates an address by ID for the authenticated user.
func UpdateAddress(wtr http.ResponseWriter, req *http.Request) {
	httpErr.SendError(wtr, req, httpErr.Global.NotImplemented, httpErr.WithDetail("UpdateAddress not yet implemented"))
}

// DeleteAddress removes an address by ID for the authenticated user.
func DeleteAddress(wtr http.ResponseWriter, req *http.Request) {
	httpErr.SendError(wtr, req, httpErr.Global.NotImplemented, httpErr.WithDetail("DeleteAddress not yet implemented"))
}
