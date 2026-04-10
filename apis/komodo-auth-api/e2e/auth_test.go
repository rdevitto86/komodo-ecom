//go:build e2e

package e2e_test

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"
)

func TestHealth(t *testing.T) {
	resp := get(t, "/health", nil)
	defer resp.Body.Close()
	checkStatus(t, resp, http.StatusOK)
}

// TestJWKS verifies the public key set is served and contains at least one key.
func TestJWKS(t *testing.T) {
	resp := get(t, "/.well-known/jwks.json", nil)
	defer resp.Body.Close()
	checkStatus(t, resp, http.StatusOK)

	var body struct {
		Keys []map[string]any `json:"keys"`
	}
	decodeJSON(t, resp, &body)
	if len(body.Keys) == 0 {
		t.Fatal("JWKS response contains no keys")
	}
}

// TestOAuthToken_ClientCredentials issues a token using client credentials.
// Requires TEST_CLIENT_ID and TEST_CLIENT_SECRET to be set.
func TestOAuthToken_ClientCredentials(t *testing.T) {
	clientID := os.Getenv("TEST_CLIENT_ID")
	clientSecret := os.Getenv("TEST_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		t.Skip("TEST_CLIENT_ID / TEST_CLIENT_SECRET not set — register a client in LocalStack secrets to enable")
	}

	body := map[string]any{
		"grant_type":    "client_credentials",
		"client_id":     clientID,
		"client_secret": clientSecret,
	}
	resp := post(t, "/oauth/token", body, nil)
	defer resp.Body.Close()
	checkStatus(t, resp, http.StatusOK)

	var tok struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}
	decodeJSON(t, resp, &tok)
	if tok.AccessToken == "" {
		t.Fatal("expected non-empty access_token in token response")
	}
	if tok.TokenType == "" {
		t.Fatal("expected non-empty token_type in token response")
	}
}

// TestOAuthToken_MissingGrantType verifies the endpoint rejects a missing grant_type.
func TestOAuthToken_MissingGrantType(t *testing.T) {
	resp := post(t, "/oauth/token", map[string]any{"client_id": "x"}, nil)
	defer resp.Body.Close()
	checkStatus(t, resp, http.StatusBadRequest)
}

// TestOAuthToken_UnknownGrantType verifies an unsupported grant type is rejected.
func TestOAuthToken_UnknownGrantType(t *testing.T) {
	resp := post(t, "/oauth/token", map[string]any{
		"grant_type": "magic_beans",
		"client_id":  "x",
	}, nil)
	defer resp.Body.Close()
	checkStatus(t, resp, http.StatusBadRequest)
}

// TestOAuthIntrospect_NoAuth verifies introspect requires a client token.
func TestOAuthIntrospect_NoAuth(t *testing.T) {
	resp := post(t, "/oauth/introspect", map[string]any{"token": "fake-token"}, nil)
	defer resp.Body.Close()
	checkStatus(t, resp, http.StatusUnauthorized)
}

// TestOAuthRevoke_NoAuth verifies revoke requires a client token.
func TestOAuthRevoke_NoAuth(t *testing.T) {
	resp := post(t, "/oauth/revoke", map[string]any{"token": "fake-token"}, nil)
	defer resp.Body.Close()
	checkStatus(t, resp, http.StatusUnauthorized)
}

// TestOAuthIntrospect_WithClientToken introspects a token issued by the service.
// Requires TEST_CLIENT_ID, TEST_CLIENT_SECRET, and TEST_JWT to be set.
func TestOAuthIntrospect_WithClientToken(t *testing.T) {
	clientID := os.Getenv("TEST_CLIENT_ID")
	clientSecret := os.Getenv("TEST_CLIENT_SECRET")
	testJWT := os.Getenv("TEST_JWT")
	if clientID == "" || clientSecret == "" || testJWT == "" {
		t.Skip("TEST_CLIENT_ID / TEST_CLIENT_SECRET / TEST_JWT not set")
	}

	// Use client credentials as the client auth token.
	clientToken := issueClientToken(t, clientID, clientSecret)

	resp := post(t, "/oauth/introspect",
		map[string]any{"token": testJWT},
		map[string]string{"Authorization": "Bearer " + clientToken},
	)
	defer resp.Body.Close()
	checkStatus(t, resp, http.StatusOK)

	var result struct {
		Active bool `json:"active"`
	}
	decodeJSON(t, resp, &result)
	if !result.Active {
		t.Fatal("expected active=true for a valid TEST_JWT")
	}
}

// issueClientToken is a test helper that fetches a client_credentials token.
func issueClientToken(t *testing.T, clientID, clientSecret string) string {
	t.Helper()
	body := map[string]any{
		"grant_type":    "client_credentials",
		"client_id":     clientID,
		"client_secret": clientSecret,
	}
	resp := post(t, "/oauth/token", body, nil)
	defer resp.Body.Close()
	checkStatus(t, resp, http.StatusOK)

	var tok struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tok); err != nil {
		t.Fatalf("decode client token: %v", err)
	}
	return tok.AccessToken
}
