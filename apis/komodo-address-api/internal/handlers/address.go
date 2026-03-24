package handlers

import (
	"net/http"

	httpErr "github.com/rdevitto86/komodo-forge-sdk-go/http/errors"
)

// POST /addresses/validate
func Validate(wtr http.ResponseWriter, req *http.Request) {
	httpErr.SendError(wtr, req, httpErr.Global.NotImplemented)
}

// POST /addresses/normalize
func Normalize(wtr http.ResponseWriter, req *http.Request) {
	httpErr.SendError(wtr, req, httpErr.Global.NotImplemented)
}

// POST /addresses/geocode
func Geocode(wtr http.ResponseWriter, req *http.Request) {
	httpErr.SendError(wtr, req, httpErr.Global.NotImplemented)
}
