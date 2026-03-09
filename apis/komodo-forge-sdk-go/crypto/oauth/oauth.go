package oauth

import (
	"strings"
)

// Define allowed scopes as a map for performance and clarity
var allowedScopes = map[string]bool{
	"read":           true,
	"write":          true,
	"admin":          true,
	"checkout:read":  true,
	"checkout:write": true,
	"orders:read":    true,
	"users:profile":  true,
}

// Checks if the provided scope string is valid according to predefined rules.
// Scopes with a "svc:" prefix are service-to-service scopes and are always valid.
func IsValidScope(scope string) bool {
	if scope == "" { return false }
	for _, s := range strings.Fields(strings.ReplaceAll(scope, ",", " ")) {
		if strings.HasPrefix(s, "svc:") { continue }
		if !allowedScopes[s] { return false }
	}
	return true
}

// Returns a slice of invalid scopes found in the provided scope string.
func GetInvalidScopes(scope string) []string {
	var invalid []string
	for _, s := range strings.Fields(strings.ReplaceAll(scope, ",", " ")) {
		if strings.HasPrefix(s, "svc:") { continue }
		if !allowedScopes[s] { invalid = append(invalid, s) }
	}
	return invalid
}

// Checks if the provided grant type string is valid according to predefined rules.
func IsValidGrantType(grantType string) bool {
	switch grantType {
		case "client_credentials", "authorization_code", "refresh_token":
			return true
		default:
			return false
	}
}
