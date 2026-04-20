package handlers

import (
	"encoding/json"
	"net/http"

	"komodo-auth-api/internal/oauth/clients"
	httpErr "github.com/rdevitto86/komodo-forge-sdk-go/http/errors"
	logger "github.com/rdevitto86/komodo-forge-sdk-go/logging/runtime"
)

// GetClientHandler returns metadata for a registered OAuth client by ID.
// The secret is never included in the response.
func GetClientHandler(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")

	id := req.PathValue("id")
	rec, ok := clients.Get(id)
	if !ok {
		httpErr.SendError(wtr, req, httpErr.Auth.InvalidClientCredentials, httpErr.WithDetail("client not found: "+id))
		return
	}

	logger.Info("client lookup: " + id)
	wtr.WriteHeader(http.StatusOK)
	json.NewEncoder(wtr).Encode(map[string]any{
		"client_id":      id,
		"name":           rec.Name,
		"allowed_scopes": rec.AllowedScopes,
	})
}

// ListClientsHandler returns all registered OAuth clients with secrets redacted.
func ListClientsHandler(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")

	clients := clients.List()
	result := make([]map[string]any, 0, len(clients))
	for id, rec := range clients {
		result = append(result, map[string]any{
			"client_id":      id,
			"name":           rec.Name,
			"allowed_scopes": rec.AllowedScopes,
		})
	}

	wtr.WriteHeader(http.StatusOK)
	json.NewEncoder(wtr).Encode(result)
}
