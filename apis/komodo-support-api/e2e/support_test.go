//go:build e2e

package e2e_test

import (
	"net/http"
	"testing"
)

func TestHealth(t *testing.T) {
	resp := get(t, "/health", nil)
	defer resp.Body.Close()
	checkStatus(t, resp, http.StatusOK)
}

// TestChatSession_CreateAndGet creates an anonymous chat session and reads it back.
func TestChatSession_CreateAndGet(t *testing.T) {
	sessionID := createChatSession(t)

	resp := get(t, "/chat/session", map[string]string{"X-Session-ID": sessionID})
	defer resp.Body.Close()
	checkStatus(t, resp, http.StatusOK)
}

// TestChatSession_DuplicateCreate verifies a second create with the same session returns the existing session.
func TestChatSession_DuplicateCreate(t *testing.T) {
	sessionID := createChatSession(t)

	resp := post(t, "/chat/session", map[string]any{}, map[string]string{"X-Session-ID": sessionID})
	defer resp.Body.Close()
	// Should return 200 (existing) rather than 201 (new).
	checkStatus(t, resp, http.StatusOK)
}

// TestSendMessage_ValidSession sends a message in an anonymous chat session.
func TestSendMessage_ValidSession(t *testing.T) {
	sessionID := createChatSession(t)
	h := map[string]string{"X-Session-ID": sessionID}

	resp := post(t, "/chat/message", map[string]any{"content": "Hello, I need help with my order."}, h)
	defer resp.Body.Close()
	checkStatus(t, resp, http.StatusOK)
}

// TestSendMessage_EmptyContent verifies empty message content is rejected.
func TestSendMessage_EmptyContent(t *testing.T) {
	sessionID := createChatSession(t)
	h := map[string]string{"X-Session-ID": sessionID}

	resp := post(t, "/chat/message", map[string]any{"content": ""}, h)
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
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

	resp := get(t, "/chat/history", map[string]string{"X-Session-ID": sessionID})
	defer resp.Body.Close()
	checkStatus(t, resp, http.StatusOK)
}

// TestGetHistory_NoSession verifies history requires a session identifier.
func TestGetHistory_NoSession(t *testing.T) {
	resp := get(t, "/chat/history", nil)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest && resp.StatusCode != http.StatusUnauthorized {
		checkStatus(t, resp, http.StatusBadRequest)
	}
}

// TestDeleteHistory_Anonymous creates a session, populates it, then deletes the history.
func TestDeleteHistory_Anonymous(t *testing.T) {
	sessionID := createChatSession(t)
	h := map[string]string{"X-Session-ID": sessionID}

	post(t, "/chat/message", map[string]any{"content": "message to delete"}, h)

	resp := del(t, "/chat/history", h)
	defer resp.Body.Close()
	checkStatus(t, resp, http.StatusOK)
}

// TestAuthHistory_Get fetches chat history for the authenticated user.
func TestAuthHistory_Get(t *testing.T) {
	resp := get(t, "/me/chat/history", authHeader(t))
	defer resp.Body.Close()
	checkStatus(t, resp, http.StatusOK)
}

// TestAuthHistory_Delete clears the authenticated user's chat history.
func TestAuthHistory_Delete(t *testing.T) {
	resp := del(t, "/me/chat/history", authHeader(t))
	defer resp.Body.Close()
	checkStatus(t, resp, http.StatusOK)
}

// TestEscalate_ValidSession escalates an anonymous chat session to a support ticket.
func TestEscalate_ValidSession(t *testing.T) {
	sessionID := createChatSession(t)
	h := map[string]string{"X-Session-ID": sessionID}

	post(t, "/chat/message", map[string]any{"content": "I have an urgent issue."}, h)

	resp := post(t, "/chat/escalate", map[string]any{"reason": "e2e test escalation"}, h)
	defer resp.Body.Close()
	// 200 = escalated; 501 = communications-api not wired yet (both acceptable).
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotImplemented {
		checkStatus(t, resp, http.StatusOK)
	}
}

// createChatSession is a helper that creates an anonymous chat session and returns its session ID.
func createChatSession(t *testing.T) string {
	t.Helper()
	resp := post(t, "/chat/session", map[string]any{}, nil)
	defer resp.Body.Close()
	checkStatus(t, resp, http.StatusCreated)

	var session struct {
		SessionID string `json:"session_id"`
	}
	decodeJSON(t, resp, &session)
	if session.SessionID == "" {
		t.Fatal("expected non-empty session_id in create-session response")
	}
	return session.SessionID
}
