package handlers

import (
	"errors"
	"net/http"

	httpErr "github.com/rdevitto86/komodo-forge-sdk-go/http/errors"

	"komodo-user-api/internal/service"
)

// sendUserError maps service-layer errors to HTTP error responses.
// not-found sentinel → 404; everything else → 500.
func sendUserError(wtr http.ResponseWriter, req *http.Request, err error) {
	if errors.Is(err, service.ErrNotFound) {
		httpErr.SendError(wtr, req, httpErr.Global.NotFound)
		return
	}
	httpErr.SendError(wtr, req, httpErr.Global.Internal)
}
