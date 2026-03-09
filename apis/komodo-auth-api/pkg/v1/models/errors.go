package models

import "fmt"

// AuthAPIError is returned by the adapter when auth-api responds with a non-2xx status.
// Callers can type-assert to inspect the status code and machine-readable error code.
type AuthAPIError struct {
	Status  int
	Code    string
	Message string
	Detail  string
}

func (e *AuthAPIError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("auth-api [%s] %s: %s", e.Code, e.Message, e.Detail)
	}
	return fmt.Sprintf("auth-api [%s] %s", e.Code, e.Message)
}
