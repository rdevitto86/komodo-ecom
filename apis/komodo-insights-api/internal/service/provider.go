package service

import (
	"context"
	"errors"

	"komodo-insights-api/internal/models"
)

// ErrNotFound is returned when the requested entity has no data to summarise.
var ErrNotFound = errors.New("entity not found")

// SummaryProvider is the LLM abstraction layer. Implementations may target
// Anthropic, AWS Bedrock, or any OpenAI-compatible on-prem endpoint.
// The concrete provider is injected at startup and swapped without changing
// handler or service code.
type SummaryProvider interface {
	// Summarize returns a natural-language summary for the given request corpus.
	// Implementations must be safe to call concurrently.
	Summarize(ctx context.Context, req models.SummaryRequest) (string, error)
}
