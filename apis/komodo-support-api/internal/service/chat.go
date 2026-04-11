package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"komodo-support-api/internal/repository"
	"komodo-support-api/internal/models"
)

const (
	maxHistory = 20 // max turns to include in context window

	// The model signals escalation by prefixing its response with this tag.
	escalateFlag = "[ESCALATE]"
)

var systemPrompt = strings.TrimSpace(`
You are a friendly and professional customer service assistant for Komodo, an e-commerce store.

Your role:
- Help customers with questions about orders, products, returns, shipping, and account issues
- Be concise, warm, and solution-focused
- If you cannot resolve an issue, escalate to a human agent

Topics you handle:
- Order status, tracking, and cancellations
- Returns and exchanges
- Product information and availability
- Shipping questions and delays
- Account and billing questions

Rules:
- Only discuss Komodo store topics. Politely redirect off-topic questions.
- Never fabricate order details, tracking numbers, or policies you are unsure about.
- Do not ask for passwords, full credit card numbers, or sensitive authentication info.
- Keep responses under 200 words unless a detailed explanation is genuinely necessary.

Escalation — start your response with exactly "[ESCALATE]" if:
- The customer explicitly asks for a human, manager, or supervisor
- The situation involves a legal threat, fraud claim, or chargeback dispute
- You cannot resolve the issue after two attempts
`)

// escalationKeywords are a client-side safety net.
// The model's self-reported [ESCALATE] prefix is the primary signal.
var escalationKeywords = []string{
	"speak to human", "talk to a person", "real person", "human agent",
	"supervisor", "manager", "this is unacceptable",
	"lawyer", "legal action", "sue", "lawsuit",
	"chargeback", "fraud", "dispute my charge",
}

type ChatService struct {
	llm  LLMProvider
	repo repository.ChatRepository
}

func NewChatService(llm LLMProvider, repo repository.ChatRepository) *ChatService {
	return &ChatService{llm: llm, repo: repo}
}

// SendMessage sends a user message, gets an AI response, persists both turns,
// and returns the assistant response with an escalation flag if triggered.
func (s *ChatService) SendMessage(ctx context.Context, sessionID, content string) (*models.SendMessageResponse, error) {
	session, err := s.repo.GetSession(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("get session: %w", err)
	}

	history, err := s.repo.GetHistory(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("get history: %w", err)
	}

	// Persist user message before calling the LLM
	userMsg := &models.ChatMessage{
		MessageID: uuid.NewString(),
		SessionID: sessionID,
		Role:      models.RoleUser,
		Content:   content,
		CreatedAt: time.Now().UTC(),
	}
	if err := s.repo.SaveMessage(ctx, userMsg); err != nil {
		return nil, fmt.Errorf("save user message: %w", err)
	}

	turns := buildTurns(history, content, maxHistory)

	assistantContent, err := s.llm.Complete(ctx, systemPrompt, turns)
	if err != nil {
		return nil, fmt.Errorf("llm complete: %w", err)
	}

	escalated, reason := detectEscalation(content, assistantContent)
	if escalated && !session.IsEscalated {
		session.IsEscalated = true
		_ = s.repo.UpdateSession(ctx, session)
	}

	// Strip the [ESCALATE] prefix before showing the response to the customer
	displayContent := strings.TrimSpace(strings.TrimPrefix(assistantContent, escalateFlag))

	assistantMsg := &models.ChatMessage{
		MessageID:   uuid.NewString(),
		SessionID:   sessionID,
		Role:        models.RoleAssistant,
		Content:     displayContent,
		CreatedAt:   time.Now().UTC(),
		IsEscalated: escalated,
	}
	if err := s.repo.SaveMessage(ctx, assistantMsg); err != nil {
		return nil, fmt.Errorf("save assistant message: %w", err)
	}

	session.MessageCount += 2
	_ = s.repo.UpdateSession(ctx, session)

	return &models.SendMessageResponse{
		MessageID:   assistantMsg.MessageID,
		SessionID:   sessionID,
		Response:    displayContent,
		IsEscalated: escalated,
		EscalatedBy: reason,
		CreatedAt:   assistantMsg.CreatedAt,
	}, nil
}

func buildTurns(history []models.ChatMessage, newContent string, limit int) []ChatTurn {
	if len(history) > limit {
		history = history[len(history)-limit:]
	}
	turns := make([]ChatTurn, 0, len(history)+1)
	for _, h := range history {
		turns = append(turns, ChatTurn{Role: string(h.Role), Content: h.Content})
	}
	turns = append(turns, ChatTurn{Role: "user", Content: newContent})
	return turns
}

func detectEscalation(userContent, assistantContent string) (bool, models.EscalationReason) {
	if strings.HasPrefix(assistantContent, escalateFlag) {
		return true, models.EscalationReasonModelSignal
	}
	lower := strings.ToLower(userContent)
	for _, kw := range escalationKeywords {
		if strings.Contains(lower, kw) {
			return true, models.EscalationReasonKeyword
		}
	}
	return false, ""
}
