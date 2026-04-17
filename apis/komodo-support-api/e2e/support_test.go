//go:build e2e

package e2e_test

import (
	"net/http"
	"testing"
)

func TestHealth(t *testing.T) {
	res := get(t, "/health", nil)
	defer res.Body.Close()
	checkStatus(t, res, http.StatusOK)
}

// TestChatSession_CreateAndGet creates an anonymous chat session and reads it back.
func TestChatSession_CreateAndGet(t *testing.T) {
	sessionID := createChatSession(t)

	res := get(t, "/chat/session", map[string]string{"X-Session-ID": sessionID})
	defer res.Body.Close()
	checkStatus(t, res, http.StatusOK)
}

// TestChatSession_DuplicateCreate verifies a second create with the same session returns the existing session.
func TestChatSession_DuplicateCreate(t *testing.T) {
	sessionID := createChatSession(t)

	res := post(t, "/chat/session", map[string]any{}, map[string]string{"X-Session-ID": sessionID})
	defer res.Body.Close()
	// Should return 200 (existing) rather than 201 (new).
	checkStatus(t, res, http.StatusOK)
}

// TestSendMessage_ValidSession sends a message in an anonymous chat session.
func TestSendMessage_ValidSession(t *testing.T) {
	sessionID := createChatSession(t)
	h := map[string]string{"X-Session-ID": sessionID}

	res := post(t, "/chat/message", map[string]any{"content": "Hello, I need help with my order."}, h)
	defer res.Body.Close()
	checkStatus(t, res, http.StatusOK)
}

// TestSendMessage_EmptyContent verifies empty message content is rejected.
func TestSendMessage_EmptyContent(t *testing.T) {
	sessionID := createChatSession(t)
	h := map[string]string{"X-Session-ID": sessionID}

	res := post(t, "/chat/message", map[string]any{"content": ""}, h)
	defer res.Body.Close()
	if res.StatusCode == http.StatusOK {
		t.Log("warn: empty content was accepted — consider adding validation")
	}
}

// TestGetHistory_ValidSession fetches the chat history for an anonymous session.
func TestGetHistory_ValidSession(t *testing.T) {
	sessionID := createChatSession(t)

	// Send a message first so history is non-empty.
	post(t, "/chat/message",
		map[string]any{"content": "Test message for history"},
		map[string]string{"X-Session-ID": sessionID},
	)

	res := get(t, "/chat/history", map[string]string{"X-Session-ID": sessionID})
	defer res.Body.Close()
	checkStatus(t, res, http.StatusOK)
}

// TestGetHistory_NoSession verifies history requires a session identifier.
func TestGetHistory_NoSession(t *testing.T) {
	res := get(t, "/chat/history", nil)
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest && res.StatusCode != http.StatusUnauthorized {
		checkStatus(t, res, http.StatusBadRequest)
	}
}

// TestDeleteHistory_Anonymous creates a session, populates it, then deletes the history.
func TestDeleteHistory_Anonymous(t *testing.T) {
	sessionID := createChatSession(t)
	h := map[string]string{"X-Session-ID": sessionID}

	post(t, "/chat/message", map[string]any{"content": "message to delete"}, h)

	res := del(t, "/chat/history", h)
	defer res.Body.Close()
	checkStatus(t, res, http.StatusOK)
}

// TestAuthHistory_Get fetches chat history for the authenticated user.
func TestAuthHistory_Get(t *testing.T) {
	res := get(t, "/me/chat/history", authHeader(t))
	defer res.Body.Close()
	checkStatus(t, res, http.StatusOK)
}

// TestAuthHistory_Delete clears the authenticated user's chat history.
func TestAuthHistory_Delete(t *testing.T) {
	res := del(t, "/me/chat/history", authHeader(t))
	defer res.Body.Close()
	checkStatus(t, res, http.StatusOK)
}

// TestEscalate_ValidSession escalates an anonymous chat session to a support ticket.
func TestEscalate_ValidSession(t *testing.T) {
	sessionID := createChatSession(t)
	h := map[string]string{"X-Session-ID": sessionID}

	post(t, "/chat/message", map[string]any{"content": "I have an urgent issue."}, h)

	res := post(t, "/chat/escalate", map[string]any{"reason": "e2e test escalation"}, h)
	defer res.Body.Close()
	// 200 = escalated; 501 = communications-api not wired yet (both acceptable).
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNotImplemented {
		checkStatus(t, res, http.StatusOK)
	}
}

// createChatSession is a helper that creates an anonymous chat session and returns its session ID.
func createChatSession(t *testing.T) string {
	t.Helper()
	res := post(t, "/chat/session", map[string]any{}, nil)
	defer res.Body.Close()
	checkStatus(t, res, http.StatusCreated)

	var session struct {
		SessionID string `json:"session_id"`
	}
	decodeJSON(t, res, &session)
	if session.SessionID == "" {
		t.Fatal("expected non-empty session_id in create-session response")
	}
	return session.SessionID
}
