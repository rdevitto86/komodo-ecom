package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"komodo-auth-api/pkg/v1/models"
)

// Client calls the auth-api over HTTP. Construct with NewClient(baseURL).
type Client struct {
	baseURL string
	http    *http.Client
}

// NewClient returns a Client targeting baseURL (e.g. "http://localhost:7011").
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		http:    &http.Client{Timeout: 10 * time.Second},
	}
}

// GetToken issues a client_credentials token on behalf of the caller.
// scope is optional; omit for a scopeless service token.
func (c *Client) GetToken(ctx context.Context, clientID, clientSecret string, scope ...string) (*models.TokenResponse, error) {
	body := models.TokenRequest{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		GrantType:    "client_credentials",
	}
	if len(scope) > 0 {
		body.Scope = scope[0]
	}

	var resp models.TokenResponse
	if err := c.post(ctx, "/oauth/token", "", body, &resp); err != nil {
		return nil, fmt.Errorf("GetToken: %w", err)
	}
	return &resp, nil
}

// IntrospectToken checks whether token is active.
// serviceToken is a valid client_credentials Bearer token used to authorize the call.
func (c *Client) IntrospectToken(ctx context.Context, token, serviceToken string) (*models.IntrospectResponse, error) {
	// Per RFC 7662 the token is sent in the request body, but our handler reads
	// it from the Authorization header via ExtractTokenFromRequest. Send token as
	// both the body field and authorize with serviceToken.
	body := map[string]string{"token": token}

	var resp models.IntrospectResponse
	if err := c.post(ctx, "/oauth/introspect", serviceToken, body, &resp); err != nil {
		return nil, fmt.Errorf("IntrospectToken: %w", err)
	}
	return &resp, nil
}

// RevokeToken revokes the given token.
// serviceToken is a valid client_credentials Bearer token used to authorize the call.
func (c *Client) RevokeToken(ctx context.Context, token, serviceToken string) error {
	body := models.RevokeRequest{Token: token}
	if err := c.post(ctx, "/oauth/revoke", serviceToken, body, nil); err != nil {
		return fmt.Errorf("RevokeToken: %w", err)
	}
	return nil
}

// ValidateToken calls the internal /internal/token/validate endpoint.
// internalBaseURL must point at the internal port (e.g. "http://localhost:7012").
// No Bearer token is required — the internal network is the access control layer.
func (c *Client) ValidateToken(ctx context.Context, token string) (*models.ValidateResponse, error) {
	body := models.ValidateRequest{Token: token}
	var resp models.ValidateResponse
	if err := c.post(ctx, "/internal/token/validate", "", body, &resp); err != nil {
		return nil, fmt.Errorf("ValidateToken: %w", err)
	}
	return &resp, nil
}

// post marshals body, POSTs to path, and unmarshals a 2xx response into out.
// If bearer is non-empty it is sent as Authorization: Bearer <bearer>.
// If out is nil the response body is discarded.
func (c *Client) post(ctx context.Context, path, bearer string, body, out any) error {
	b, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if bearer != "" {
		req.Header.Set("Authorization", "Bearer "+bearer)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiErr struct {
			Status  int    `json:"status"`
			Code    string `json:"code"`
			Message string `json:"message"`
			Detail  string `json:"detail"`
		}
		_ = json.Unmarshal(raw, &apiErr)
		return &models.AuthAPIError{
			Status:  resp.StatusCode,
			Code:    apiErr.Code,
			Message: apiErr.Message,
			Detail:  apiErr.Detail,
		}
	}

	if out != nil {
		if err := json.Unmarshal(raw, out); err != nil {
			return fmt.Errorf("unmarshal response: %w", err)
		}
	}
	return nil
}
