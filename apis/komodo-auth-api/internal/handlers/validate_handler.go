package handlers

import (
	"encoding/json"
	"net/http"

	"komodo-forge-sdk-go/crypto/jwt"
	logger "komodo-forge-sdk-go/logging/runtime"
)

type ValidateRequest struct {
	Token string `json:"token"`
}

type ValidateResponse struct {
	Valid    bool     `json:"valid"`
	Subject  string   `json:"sub,omitempty"`
	Scopes   []string `json:"scopes,omitempty"`
	Error    string   `json:"error,omitempty"`
}

// ValidateTokenHandler verifies a JWT and returns its claims.
// This is an internal-only endpoint — no Bearer token required.
// The internal network (VPC/IAM) is the access control layer.
func ValidateTokenHandler(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")
	wtr.Header().Set("Cache-Control", "no-store")

	var body ValidateRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil || body.Token == "" {
		wtr.WriteHeader(http.StatusOK)
		json.NewEncoder(wtr).Encode(ValidateResponse{Valid: false, Error: "missing or unparseable token"})
		return
	}

	valid, err := jwt.ValidateToken(body.Token)
	if !valid || err != nil {
		wtr.WriteHeader(http.StatusOK)
		json.NewEncoder(wtr).Encode(ValidateResponse{Valid: false, Error: err.Error()})
		return
	}

	claims, err := jwt.ParseClaims(body.Token)
	if err != nil {
		wtr.WriteHeader(http.StatusOK)
		json.NewEncoder(wtr).Encode(ValidateResponse{Valid: false, Error: "failed to parse claims"})
		return
	}

	logger.Info("token validated for subject: " + claims.Subject)
	wtr.WriteHeader(http.StatusOK)
	json.NewEncoder(wtr).Encode(ValidateResponse{
		Valid:   true,
		Subject: claims.Subject,
		Scopes:  claims.Scopes,
	})
}
