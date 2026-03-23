package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// AnthropicProvider implements LLMProvider using the Anthropic Messages API.
// Model defaults to claude-haiku-4-5-20251001 (fast, cost-effective for support).
// Override with a different anthropic.Model constant for experimentation.
type AnthropicProvider struct {
	client    anthropic.Client
	model     anthropic.Model
	maxTokens int64
}

func NewAnthropicProvider(apiKey string) LLMProvider {
	return &AnthropicProvider{
		client:    anthropic.NewClient(option.WithAPIKey(apiKey)),
		model:     anthropic.ModelClaudeHaiku4_5_20251001,
		maxTokens: 1024,
	}
}

func (p *AnthropicProvider) Complete(ctx context.Context, systemPrompt string, history []ChatTurn) (string, error) {
	msgs := make([]anthropic.MessageParam, 0, len(history))
	for _, turn := range history {
		if turn.Role == "user" {
			msgs = append(msgs, anthropic.NewUserMessage(anthropic.NewTextBlock(turn.Content)))
		} else {
			msgs = append(msgs, anthropic.NewAssistantMessage(anthropic.NewTextBlock(turn.Content)))
		}
	}

	resp, err := p.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     p.model,
		MaxTokens: p.maxTokens,
		System: []anthropic.TextBlockParam{
			{Text: systemPrompt},
		},
		Messages: msgs,
	})
	if err != nil {
		return "", fmt.Errorf("anthropic complete: %w", err)
	}

	var parts []string
	for _, block := range resp.Content {
		switch b := block.AsAny().(type) {
		case anthropic.TextBlock:
			parts = append(parts, b.Text)
		}
	}
	return strings.Join(parts, "\n"), nil
}
