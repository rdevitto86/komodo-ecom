package handlers

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math/big"
	"net/http"

	"komodo-forge-sdk-go/config"
	httpErr "komodo-forge-sdk-go/http/errors"
	logger "komodo-forge-sdk-go/logging/runtime"
)

type JWK struct {
	Kty string `json:"kty"` // Key Type (e.g., "RSA")
	Use string `json:"use"` // Public Key Use (e.g., "sig" for signature)
	Kid string `json:"kid"` // Key ID
	Alg string `json:"alg"` // Algorithm (e.g., "RS256")
	N   string `json:"n"`   // RSA Modulus (base64url encoded)
	E   string `json:"e"`   // RSA Exponent (base64url encoded)
}

type JWKS struct {
	Keys []JWK `json:"keys"`
}

// Public keys for JWT verification
func JWKSHandler(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")
	wtr.Header().Set("Cache-Control", "public, max-age=3600") // Cache for 1 hour

	// Get public key from environment
	publicKeyPEM := config.GetConfigValue("JWT_PUBLIC_KEY")
	if publicKeyPEM == "" {
		logger.Error("JWT_PUBLIC_KEY not configured", fmt.Errorf("JWT_PUBLIC_KEY not configured"))
		httpErr.SendError(
			wtr, req, httpErr.Auth.InvalidKey, httpErr.WithDetail("public key not configured"),
		)
		return
	}

	// Parse PEM block
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		logger.Error("failed to parse PEM block containing public key", fmt.Errorf("failed to parse PEM block containing public key"))
		httpErr.SendError(
			wtr, req, httpErr.Auth.InvalidKey, httpErr.WithDetail("failed to parse public key"),
		)
		return
	}

	// Parse public key
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		logger.Error("failed to parse public key", err)
		httpErr.SendError(
			wtr, req, httpErr.Auth.InvalidKey, httpErr.WithDetail("failed to parse public key"),
		)
		return
	}

	// Cast to RSA public key
	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		logger.Error("public key is not RSA", fmt.Errorf("public key is not RSA"))
		httpErr.SendError(
			wtr, req, httpErr.Auth.InvalidKey, httpErr.WithDetail("public key is invalid"),
		)
		return
	}

	// Convert RSA modulus (N) and exponent (E) to base64url
	nBytes := rsaPub.N.Bytes()
	eBytes := big.NewInt(int64(rsaPub.E)).Bytes()

	// Create JWKS response
	jwks := JWKS{
		Keys: []JWK{
			{
				Kty: "RSA",
				Use: "sig",
				Kid: config.GetConfigValue("JWT_KID"),
				Alg: "RS256",
				N: base64.RawURLEncoding.EncodeToString(nBytes),
				E: base64.RawURLEncoding.EncodeToString(eBytes),
			},
		},
	}

	wtr.WriteHeader(http.StatusOK)
	json.NewEncoder(wtr).Encode(jwks)
}