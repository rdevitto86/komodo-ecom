package v1

import (
	"komodo-auth-api/pkg/v1/adapters"
	"komodo-auth-api/pkg/v1/models"
)

// Adapter is the auth-api HTTP client. Construct with NewAdapter(baseURL).
type Adapter = adapters.Client

// NewAdapter returns a Client targeting baseURL (e.g. "http://localhost:7011").
var NewAdapter = adapters.NewClient

// Auth model types — stable import path for consumers.
type (
	TokenRequest       = models.TokenRequest
	TokenResponse      = models.TokenResponse
	IntrospectResponse = models.IntrospectResponse
	RevokeRequest      = models.RevokeRequest
	ValidateRequest    = models.ValidateRequest
	ValidateResponse   = models.ValidateResponse
	RegisteredClient   = models.RegisteredClient
	JWK                = models.JWK
	JWKS               = models.JWKS
	AuthAPIError       = models.AuthAPIError
)
