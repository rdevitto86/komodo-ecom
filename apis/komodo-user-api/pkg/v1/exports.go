package v1

import (
	"komodo-user-api/pkg/v1/adapters"
	"komodo-user-api/pkg/v1/mocks"
)

// Adapter is the typed HTTP client for calling the user-api internal server.
// Consuming services instantiate it via NewAdapter() and inject USER_API_INTERNAL_URL.
type Adapter = adapters.Client

var (
	NewAdapter = adapters.NewClient
	Mocks      = mocks.Mocks
)
