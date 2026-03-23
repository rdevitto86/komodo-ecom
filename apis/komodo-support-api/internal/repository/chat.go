package repository

import (
	"context"
	"fmt"
	"sync"
	"time"

	"komodo-support-api/pkg/v1/models"
)

// ChatHistoryTTL is the maximum retention period for all chat sessions (anonymous and authenticated).
// After this window, DynamoDB TTL will auto-delete — no explicit cleanup needed.
// User-initiated deletion (DELETE /chat/history) clears immediately regardless of TTL.
//
// Audit note: before deleting messages (user-initiated or TTL), the DynamoDB implementation
// should emit a deletion audit event for the business-side audit trail.
// TODO: define audit event schema and destination (DynamoDB audit table or S3 archive).
const ChatHistoryTTL = 30 * 24 // hours

// ChatRepository defines the storage interface for chat sessions and messages.
// Current implementation: in-memory (swap for DynamoDB — see TODO below).
//
// TODO: DynamoDB table design
//   Table: komodo-chat-sessions
//   PK: session_id (string) — used for both anonymous and authenticated sessions
//   SK: message_id (string) | "SESSION#META" for the session record itself
//   GSI1: user_id (PK) + created_at (SK) — fetch all sessions for a logged-in user
//   TTL: expires_at (epoch) — DynamoDB TTL auto-deletes sessions after ChatHistoryTTL
//   Note: anonymous sessions have no user_id; authenticated sessions populate it
//   at creation time (or lazily on first authenticated request via MergeSession).
type ChatRepository interface {
	// Session management
	CreateSession(ctx context.Context, session *models.ChatSession) error
	GetSession(ctx context.Context, sessionID string) (*models.ChatSession, error)
	UpdateSession(ctx context.Context, session *models.ChatSession) error

	// Message management
	SaveMessage(ctx context.Context, msg *models.ChatMessage) error
	GetHistory(ctx context.Context, sessionID string) ([]models.ChatMessage, error)
	DeleteHistory(ctx context.Context, sessionID string) error

	// Merge: called at login to associate an anonymous session with a user_id.
	// TODO: implement once DynamoDB is wired — update session record + GSI1 user_id.
	MergeSession(ctx context.Context, sessionID, userID string) error
}

// InMemoryChatRepository is a thread-safe in-memory implementation.
// Replace with a DynamoDB implementation before production.
type InMemoryChatRepository struct {
	mu       sync.RWMutex
	sessions map[string]*models.ChatSession
	messages map[string][]models.ChatMessage // keyed by session_id
}

func NewInMemoryChatRepository() ChatRepository {
	return &InMemoryChatRepository{
		sessions: make(map[string]*models.ChatSession),
		messages: make(map[string][]models.ChatMessage),
	}
}

func (r *InMemoryChatRepository) CreateSession(_ context.Context, session *models.ChatSession) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sessions[session.SessionID] = session
	return nil
}

func (r *InMemoryChatRepository) GetSession(_ context.Context, sessionID string) (*models.ChatSession, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.sessions[sessionID]
	if !ok {
		return nil, fmt.Errorf("session not found: %s", sessionID)
	}
	if time.Now().After(s.ExpiresAt) {
		return nil, fmt.Errorf("session expired: %s", sessionID)
	}
	return s, nil
}

func (r *InMemoryChatRepository) UpdateSession(_ context.Context, session *models.ChatSession) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sessions[session.SessionID] = session
	return nil
}

func (r *InMemoryChatRepository) SaveMessage(_ context.Context, msg *models.ChatMessage) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.messages[msg.SessionID] = append(r.messages[msg.SessionID], *msg)
	return nil
}

func (r *InMemoryChatRepository) GetHistory(_ context.Context, sessionID string) ([]models.ChatMessage, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	msgs := r.messages[sessionID]
	result := make([]models.ChatMessage, len(msgs))
	copy(result, msgs)
	return result, nil
}

func (r *InMemoryChatRepository) DeleteHistory(_ context.Context, sessionID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.messages, sessionID)
	return nil
}

// MergeSession folds an anonymous session into a user session.
// Messages are prepended to the user's history so context is preserved.
// The anonymous session and its messages are cleaned up after the merge.
//
// TODO (DynamoDB): TransactWriteItems —
//   1. UpdateItem: set user_id on the session record, add to GSI1
//   2. BatchWriteItem: rekey all message items from anon session_id → user_id
//   3. DeleteItem: remove the anonymous session META record
func (r *InMemoryChatRepository) MergeSession(_ context.Context, anonSessionID, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	anonSession, ok := r.sessions[anonSessionID]
	if !ok {
		// No anonymous session to merge — not an error, just a no-op
		return nil
	}

	// Copy messages from anonymous session, rekeyed to user session
	anonMsgs := r.messages[anonSessionID]
	userMsgs := r.messages[userID]

	merged := make([]models.ChatMessage, 0, len(anonMsgs)+len(userMsgs))
	for _, m := range anonMsgs {
		m.SessionID = userID
		merged = append(merged, m)
	}
	merged = append(merged, userMsgs...)
	r.messages[userID] = merged

	// Upsert user session, preserving the earlier created_at if possible
	if existing, ok := r.sessions[userID]; ok {
		if anonSession.CreatedAt.Before(existing.CreatedAt) {
			existing.CreatedAt = anonSession.CreatedAt
		}
		existing.IsAuthenticated = true
		existing.UserID = userID
		r.sessions[userID] = existing
	} else {
		anonSession.SessionID = userID
		anonSession.UserID = userID
		anonSession.IsAuthenticated = true
		r.sessions[userID] = anonSession
	}

	// Clean up anonymous session
	delete(r.sessions, anonSessionID)
	delete(r.messages, anonSessionID)

	return nil
}
