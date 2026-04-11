package models

import "time"

// Role represents who sent a chat message
type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

// EscalationReason indicates why a chat was escalated to a human agent
type EscalationReason string

const (
	EscalationReasonUserRequest  EscalationReason = "user_request"
	EscalationReasonKeyword      EscalationReason = "keyword_detected"
	EscalationReasonModelSignal  EscalationReason = "model_signal"
	EscalationReasonMaxRetries   EscalationReason = "max_retries"
)

// ChatSession tracks an ongoing support conversation
type ChatSession struct {
	SessionID       string    `json:"session_id"`
	UserID          string    `json:"user_id,omitempty"`     // empty for anonymous
	IsAuthenticated bool      `json:"is_authenticated"`
	IsEscalated     bool      `json:"is_escalated"`
	CreatedAt       time.Time `json:"created_at"`
	ExpiresAt       time.Time `json:"expires_at"`
	MessageCount    int       `json:"message_count"`
}

// ChatMessage is a single turn in a conversation
type ChatMessage struct {
	MessageID   string           `json:"message_id"`
	SessionID   string           `json:"session_id"`
	Role        Role             `json:"role"`
	Content     string           `json:"content"`
	CreatedAt   time.Time        `json:"created_at"`
	IsEscalated bool             `json:"is_escalated,omitempty"`
}

// --- Request/Response types ---

type SendMessageRequest struct {
	Message string `json:"message"`
}

type SendMessageResponse struct {
	MessageID   string           `json:"message_id"`
	SessionID   string           `json:"session_id"`
	Response    string           `json:"response"`
	IsEscalated bool             `json:"is_escalated"`
	EscalatedBy EscalationReason `json:"escalated_by,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
}

type ChatHistoryResponse struct {
	SessionID string        `json:"session_id"`
	Messages  []ChatMessage `json:"messages"`
	Total     int           `json:"total"`
}

type CreateSessionResponse struct {
	SessionID string    `json:"session_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

type EscalateRequest struct {
	Reason string `json:"reason,omitempty"`
}

type EscalateResponse struct {
	SessionID string           `json:"session_id"`
	Reason    EscalationReason `json:"reason"`
	// TODO: ticket_id once async ticket creation is implemented
}
