package registry

import (
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/rdevitto86/komodo-forge-sdk-go/config"
)

// ClientRecord is the stored representation of a registered OAuth client.
// The Secret is used only for validation and is never returned in API responses.
type ClientRecord struct {
	Name          string   `json:"name"`
	Secret        string   `json:"secret"`
	AllowedScopes []string `json:"allowed_scopes"`
}

// HasScope reports whether scope is permitted for this client.
// An empty AllowedScopes list grants access to any scope.
func (r ClientRecord) HasScope(scope string) bool {
	if len(r.AllowedScopes) == 0 {
		return true
	}
	for _, s := range strings.Fields(scope) {
		allowed := false
		for _, a := range r.AllowedScopes {
			if a == s {
				allowed = true
				break
			}
		}
		if !allowed {
			return false
		}
	}
	return true
}

type store struct {
	mu      sync.RWMutex
	clients map[string]ClientRecord
}

var global *store
var once sync.Once

// Load parses REGISTERED_CLIENTS from config into the global registry.
// Must be called once after Secrets Manager has been bootstrapped.
func Load() error {
	raw := config.GetConfigValue("REGISTERED_CLIENTS")
	if raw == "" {
		return fmt.Errorf("registry: REGISTERED_CLIENTS not configured")
	}

	var clients map[string]ClientRecord
	if err := json.Unmarshal([]byte(raw), &clients); err != nil {
		return fmt.Errorf("registry: failed to parse REGISTERED_CLIENTS: %w", err)
	}

	once.Do(func() {
		global = &store{clients: clients}
	})
	return nil
}

// Validate returns true if clientID exists and secret matches using a
// constant-time comparison to prevent timing attacks.
func Validate(clientID, secret string) bool {
	if global == nil {
		return false
	}
	global.mu.RLock()
	rec, ok := global.clients[clientID]
	global.mu.RUnlock()
	if !ok {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(rec.Secret), []byte(secret)) == 1
}

// Get returns the ClientRecord for clientID without exposing the secret.
func Get(clientID string) (ClientRecord, bool) {
	if global == nil {
		return ClientRecord{}, false
	}
	global.mu.RLock()
	rec, ok := global.clients[clientID]
	global.mu.RUnlock()
	if !ok {
		return ClientRecord{}, false
	}
	// Strip secret before returning
	return ClientRecord{Name: rec.Name, AllowedScopes: rec.AllowedScopes}, true
}

// List returns all registered clients with secrets redacted.
func List() map[string]ClientRecord {
	if global == nil {
		return nil
	}
	global.mu.RLock()
	defer global.mu.RUnlock()

	result := make(map[string]ClientRecord, len(global.clients))
	for id, rec := range global.clients {
		result[id] = ClientRecord{Name: rec.Name, AllowedScopes: rec.AllowedScopes}
	}
	return result
}
