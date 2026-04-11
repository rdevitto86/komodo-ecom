package models

// TokenRequest is the body sent to POST /oauth/token.
type TokenRequest struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	GrantType    string `json:"grantType"`
	Scope        string `json:"scope,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
	Code         string `json:"code,omitempty"`
	RedirectURI  string `json:"redirectUri,omitempty"`
}

// TokenResponse is returned by POST /oauth/token on success.
type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	TokenType    string `json:"tokenType"`
	ExpiresIn    int    `json:"expiresIn"`
	RefreshToken string `json:"refreshToken,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

// IntrospectResponse is returned by POST /oauth/introspect.
// Active is false for expired, revoked, or unrecognized tokens.
type IntrospectResponse struct {
	Active    bool   `json:"active"`
	Scope     string `json:"scope,omitempty"`
	ClientID  string `json:"clientId,omitempty"`
	TokenType string `json:"tokenType,omitempty"`
	Exp       int64  `json:"exp,omitempty"`
	Iat       int64  `json:"iat,omitempty"`
	Sub       string `json:"sub,omitempty"`
	Aud       string `json:"aud,omitempty"`
}

// RevokeRequest is the body sent to POST /oauth/revoke.
type RevokeRequest struct {
	Token         string `json:"token"`
	TokenTypeHint string `json:"token_type_hint,omitempty"` // "access_token" or "refresh_token"
}

// JWK is a single JSON Web Key as returned by GET /.well-known/jwks.json.
type JWK struct {
	Kty string `json:"kty"`
	Use string `json:"use"`
	Kid string `json:"kid"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

// JWKS is the JSON Web Key Set returned by the JWKS endpoint.
type JWKS struct {
	Keys []JWK `json:"keys"`
}

// RegisteredClient is the public view of an OAuth client — secret is never included.
type RegisteredClient struct {
	ClientID      string   `json:"client_id"`
	Name          string   `json:"name"`
	AllowedScopes []string `json:"allowed_scopes"`
}

// ValidateRequest is the body sent to POST /internal/token/validate.
type ValidateRequest struct {
	Token string `json:"token"`
}

// ValidateResponse is returned by POST /internal/token/validate.
type ValidateResponse struct {
	Valid    bool     `json:"valid"`
	Subject  string   `json:"sub,omitempty"`
	Scopes   []string `json:"scopes,omitempty"`
	Error    string   `json:"error,omitempty"`
}
