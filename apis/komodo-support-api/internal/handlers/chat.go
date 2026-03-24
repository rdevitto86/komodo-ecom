package handlers

import (
	"encoding/json"
	"net/http"

	ctxKeys "github.com/rdevitto86/komodo-forge-sdk-go/http/context"
	httpErr "github.com/rdevitto86/komodo-forge-sdk-go/http/errors"
	logger "github.com/rdevitto86/komodo-forge-sdk-go/logging/runtime"

	"komodo-support-api/internal/repository"
	"komodo-support-api/internal/service"
	"komodo-support-api/pkg/v1/models"
)

type ChatHandler struct {
	svc  *service.ChatService
	repo repository.ChatRepository
}

func NewChatHandler(svc *service.ChatService, repo repository.ChatRepository) *ChatHandler {
	return &ChatHandler{svc: svc, repo: repo}
}

// resolveSession returns the effective session ID and whether a merge is needed.
// A merge is needed when a JWT-authenticated user also carries an anonymous cookie session —
// the anonymous history should be folded into their user session before continuing.
//
// Returns: (effectiveSessionID, anonSessionID to merge or "", ok)
func resolveSession(req *http.Request) (sessionID string, anonSessionID string, ok bool) {
	userID, hasJWT := req.Context().Value(ctxKeys.USER_ID_KEY).(string)
	hasJWT = hasJWT && userID != ""

	cookieID := ""
	if cookie, err := req.Cookie(SessionCookieName); err == nil && cookie.Value != "" {
		cookieID = cookie.Value
	}

	switch {
	case hasJWT && cookieID != "" && cookieID != userID:
		// Authenticated user carrying an anonymous cookie — merge needed
		return userID, cookieID, true
	case hasJWT:
		return userID, "", true
	case cookieID != "":
		return cookieID, "", true
	default:
		return "", "", false
	}
}

// SendMessage handles POST /chat/message.
// Triggers an anonymous→user session merge if a JWT user has an active cookie session.
func (h *ChatHandler) SendMessage(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")

	sessionID, anonSessionID, ok := resolveSession(req)
	if !ok {
		httpErr.SendError(wtr, req, httpErr.Global.Unauthorized, httpErr.WithDetail("no active session — call POST /chat/session first"))
		return
	}

	// Merge anonymous history into the user's session on first authenticated message
	if anonSessionID != "" {
		if err := h.repo.MergeSession(req.Context(), anonSessionID, sessionID); err != nil {
			logger.Warn("failed to merge anonymous session", "anon_id", anonSessionID, "user_id", sessionID)
			// Non-fatal: proceed with user session; anonymous history may be lost
		} else {
			logger.Info("merged anonymous session into user session", "anon_id", anonSessionID, "user_id", sessionID)
		}
	}

	var body models.SendMessageRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		httpErr.SendError(wtr, req, httpErr.Global.BadRequest, httpErr.WithDetail("invalid request body"))
		return
	}
	if body.Message == "" {
		httpErr.SendError(wtr, req, httpErr.Global.BadRequest, httpErr.WithDetail("message is required"))
		return
	}

	resp, err := h.svc.SendMessage(req.Context(), sessionID, body.Message)
	if err != nil {
		logger.Error("failed to send chat message", err)
		httpErr.SendError(wtr, req, models.Err.ChatError)
		return
	}

	wtr.WriteHeader(http.StatusOK)
	json.NewEncoder(wtr).Encode(resp)
}

// GetHistory handles GET /chat/history and GET /me/chat/history
func (h *ChatHandler) GetHistory(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")

	sessionID, _, ok := resolveSession(req)
	if !ok {
		httpErr.SendError(wtr, req, httpErr.Global.Unauthorized, httpErr.WithDetail("no active session"))
		return
	}

	msgs, err := h.repo.GetHistory(req.Context(), sessionID)
	if err != nil {
		logger.Error("failed to get chat history", err)
		httpErr.SendError(wtr, req, httpErr.Global.Internal)
		return
	}

	wtr.WriteHeader(http.StatusOK)
	json.NewEncoder(wtr).Encode(models.ChatHistoryResponse{
		SessionID: sessionID,
		Messages:  msgs,
		Total:     len(msgs),
	})
}

// DeleteHistory handles DELETE /chat/history and DELETE /me/chat/history.
// User-initiated clear — wipes conversation. History is also auto-expired by TTL (max 30 days).
// TODO: emit an audit event before deletion for the business-side audit trail.
func (h *ChatHandler) DeleteHistory(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")

	sessionID, _, ok := resolveSession(req)
	if !ok {
		httpErr.SendError(wtr, req, httpErr.Global.Unauthorized, httpErr.WithDetail("no active session"))
		return
	}

	if err := h.repo.DeleteHistory(req.Context(), sessionID); err != nil {
		logger.Error("failed to delete chat history", err)
		httpErr.SendError(wtr, req, httpErr.Global.Internal)
		return
	}

	wtr.WriteHeader(http.StatusNoContent)
}

// Escalate handles POST /chat/escalate — explicit user-initiated escalation to a live agent.
//
// Current behaviour: marks session as escalated and publishes a business event stub.
//
// TODO (human agent side — not yet designed):
//   - Define the agent portal (separate login view, customer lookup, agent actions)
//   - Wire escalation event to SQS queue (queue name TBD from business event pipeline)
//   - SNS topic subscription for agent routing service to consume
//   - Agent availability / queue depth logic lives outside this service
func (h *ChatHandler) Escalate(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")

	sessionID, _, ok := resolveSession(req)
	if !ok {
		httpErr.SendError(wtr, req, httpErr.Global.Unauthorized, httpErr.WithDetail("no active session"))
		return
	}

	session, err := h.repo.GetSession(req.Context(), sessionID)
	if err != nil {
		httpErr.SendError(wtr, req, models.Err.SessionNotFound)
		return
	}

	if session.IsEscalated {
		// Already escalated — idempotent response
		wtr.WriteHeader(http.StatusOK)
		json.NewEncoder(wtr).Encode(models.EscalateResponse{
			SessionID: sessionID,
			Reason:    models.EscalationReasonUserRequest,
		})
		return
	}

	session.IsEscalated = true
	if err := h.repo.UpdateSession(req.Context(), session); err != nil {
		logger.Error("failed to update session for escalation", err)
		httpErr.SendError(wtr, req, httpErr.Global.Internal)
		return
	}

	// Publish business event for agent routing.
	// TODO: replace with real SQS publish once the business event pipeline is wired.
	// Event shape (draft):
	//   { "event": "chat.escalated", "session_id": "...", "user_id": "...", "timestamp": "..." }
	// Consumer: agent routing service (queue name TBD)
	publishEscalationEvent(sessionID, session.UserID)

	wtr.WriteHeader(http.StatusOK)
	json.NewEncoder(wtr).Encode(models.EscalateResponse{
		SessionID: sessionID,
		Reason:    models.EscalationReasonUserRequest,
	})
}

// publishEscalationEvent is a stub — replace with SQS/SNS publish when the
// business event pipeline (Dynamo CDC + EventBridge or direct SQS) is decided.
func publishEscalationEvent(sessionID, userID string) {
	logger.Info("chat.escalated event",
		"session_id", sessionID,
		"user_id", userID,
		// TODO: publish to SQS queue for agent routing
	)
}
