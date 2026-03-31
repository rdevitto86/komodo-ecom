package handlers

import (
	"encoding/json"
	"net/http"

	"komodo-address-api/internal/provider"

	httpErr "github.com/rdevitto86/komodo-forge-sdk-go/http/errors"
	logger "github.com/rdevitto86/komodo-forge-sdk-go/logging/runtime"
)

// POST /addresses/validate
func Validate(p *provider.Client) http.HandlerFunc {
	return func(wtr http.ResponseWriter, req *http.Request) {
		var input provider.AddressInput
		if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
			httpErr.SendError(wtr, req, httpErr.Global.BadRequest, httpErr.WithDetail("invalid request body"))
			return
		}

		result, err := p.Validate(req.Context(), input)
		if err != nil {
			logger.Error("validate: provider error", err, logger.FromContext(req.Context())...)
			httpErr.SendError(wtr, req, httpErr.Global.BadGateway, httpErr.WithDetail("address provider unavailable"))
			return
		}

		wtr.Header().Set("Content-Type", "application/json")
		wtr.WriteHeader(http.StatusOK)
		json.NewEncoder(wtr).Encode(result)
	}
}

// POST /addresses/normalize
func Normalize(p *provider.Client) http.HandlerFunc {
	return func(wtr http.ResponseWriter, req *http.Request) {
		var input provider.AddressInput
		if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
			httpErr.SendError(wtr, req, httpErr.Global.BadRequest, httpErr.WithDetail("invalid request body"))
			return
		}

		normalized, err := p.Normalize(req.Context(), input)
		if err != nil {
			logger.Error("normalize: provider error", err, logger.FromContext(req.Context())...)
			httpErr.SendError(wtr, req, httpErr.Global.BadGateway, httpErr.WithDetail("address provider unavailable"))
			return
		}

		wtr.Header().Set("Content-Type", "application/json")
		wtr.WriteHeader(http.StatusOK)
		json.NewEncoder(wtr).Encode(map[string]any{"address": normalized})
	}
}

// POST /addresses/geocode
func Geocode(p *provider.Client) http.HandlerFunc {
	return func(wtr http.ResponseWriter, req *http.Request) {
		var input provider.AddressInput
		if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
			httpErr.SendError(wtr, req, httpErr.Global.BadRequest, httpErr.WithDetail("invalid request body"))
			return
		}

		geocoded, err := p.Geocode(req.Context(), input)
		if err != nil {
			logger.Error("geocode: provider error", err, logger.FromContext(req.Context())...)
			httpErr.SendError(wtr, req, httpErr.Global.BadGateway, httpErr.WithDetail("address provider unavailable"))
			return
		}

		wtr.Header().Set("Content-Type", "application/json")
		wtr.WriteHeader(http.StatusOK)
		json.NewEncoder(wtr).Encode(geocoded)
	}
}
