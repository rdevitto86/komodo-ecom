package handlers

import (
	"encoding/json"
	"net/http"

	httpErr "github.com/rdevitto86/komodo-forge-sdk-go/http/errors"
	// logger "github.com/rdevitto86/komodo-forge-sdk-go/logging/runtime"
)

// TODO: Replace any with proper types from a provider sdk package
// For example: provider.AddressInput, provider.AddressValidationResult, etc.

// POST /addresses/validate
func Validate(p any) http.HandlerFunc {
	return func(wtr http.ResponseWriter, req *http.Request) {
		var input any
		if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
			httpErr.SendError(wtr, req, httpErr.Global.BadRequest, httpErr.WithDetail("invalid request body"))
			return
		}

		// result, err := p.Validate(req.Context(), input)
		// if err != nil {
		// 	logger.Error("validate: provider error", err, logger.FromContext(req.Context())...)
		// 	httpErr.SendError(wtr, req, httpErr.Global.BadGateway, httpErr.WithDetail("address provider unavailable"))
		// 	return
		// }

		wtr.Header().Set("Content-Type", "application/json")
		wtr.WriteHeader(http.StatusOK)
		json.NewEncoder(wtr).Encode(map[string]any{"address": input})
	}
}

// POST /addresses/normalize
func Normalize(p any) http.HandlerFunc {
	return func(wtr http.ResponseWriter, req *http.Request) {
		var input any
		if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
			httpErr.SendError(wtr, req, httpErr.Global.BadRequest, httpErr.WithDetail("invalid request body"))
			return
		}

		// normalized, err := p.Normalize(req.Context(), input)
		// if err != nil {
		// 	logger.Error("normalize: provider error", err, logger.FromContext(req.Context())...)
		// 	httpErr.SendError(wtr, req, httpErr.Global.BadGateway, httpErr.WithDetail("address provider unavailable"))
		// 	return
		// }

		wtr.Header().Set("Content-Type", "application/json")
		wtr.WriteHeader(http.StatusOK)
		json.NewEncoder(wtr).Encode(map[string]any{"address": input})
	}
}

// POST /addresses/geocode
func Geocode(p any) http.HandlerFunc {
	return func(wtr http.ResponseWriter, req *http.Request) {
		var input any
		if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
			httpErr.SendError(wtr, req, httpErr.Global.BadRequest, httpErr.WithDetail("invalid request body"))
			return
		}

		// geocoded, err := p.Geocode(req.Context(), input)
		// if err != nil {
		// 	logger.Error("geocode: provider error", err, logger.FromContext(req.Context())...)
		// 	httpErr.SendError(wtr, req, httpErr.Global.BadGateway, httpErr.WithDetail("address provider unavailable"))
		// 	return
		// }

		wtr.Header().Set("Content-Type", "application/json")
		wtr.WriteHeader(http.StatusOK)
		json.NewEncoder(wtr).Encode(map[string]any{"address": input})
	}
}
