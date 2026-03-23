package service

import "context"

// ChatTurn is a single message in a conversation sent to the LLM provider.
type ChatTurn struct {
	Role    string // "user" or "assistant"
	Content string
}

// LLMProvider is the abstraction over any chat completion provider.
// Swap implementations in main.go — current default: Anthropic.
// Future options: OpenAI, AWS Bedrock, local Ollama, etc.
type LLMProvider interface {
	// Complete sends a conversation to the model and returns the assistant response.
	// systemPrompt is applied once at the start (provider implementations handle
	// the format differences — Anthropic uses a dedicated field, OpenAI uses a
	// system role message, Bedrock varies by model).
	Complete(ctx context.Context, systemPrompt string, history []ChatTurn) (string, error)
}
