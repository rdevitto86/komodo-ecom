package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/rdevitto86/komodo-forge-sdk-go/config"
	httpErr "github.com/rdevitto86/komodo-forge-sdk-go/http/errors"
	logger "github.com/rdevitto86/komodo-forge-sdk-go/logging/runtime"

	"komodo-support-api/internal/repository"
	"komodo-support-api/pkg/v1/models"
)

const SessionCookieName = "komodo_chat_sid"
const defaultSessionTTLDays = 30

type SessionHandler struct {
	repo repository.ChatRepository
}

func NewSessionHandler(repo repository.ChatRepository) *SessionHandler {
	return &SessionHandler{repo: repo}
}

// CreateSession creates a new anonymous chat session and sets the session cookie.
// If the request already carries a valid session cookie, returns the existing session.
func (h *SessionHandler) CreateSession(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")

	// Return existing session if cookie is present and valid
	if cookie, err := req.Cookie(SessionCookieName); err == nil {
		if existing, err := h.repo.GetSession(req.Context(), cookie.Value); err == nil {
			setSessionCookie(wtr, existing.SessionID, existing.ExpiresAt)
			wtr.WriteHeader(http.StatusOK)
			json.NewEncoder(wtr).Encode(models.CreateSessionResponse{
				SessionID: existing.SessionID,
				ExpiresAt: existing.ExpiresAt,
			})
			return
		}
	}

	ttlDays := sessionTTLDays()
	now := time.Now().UTC()
	session := &models.ChatSession{
		SessionID: uuid.NewString(),
		CreatedAt: now,
		ExpiresAt: now.Add(time.Duration(ttlDays) * 24 * time.Hour),
	}

	if err := h.repo.CreateSession(req.Context(), session); err != nil {
		logger.Error("failed to create chat session", err)
		httpErr.SendError(wtr, req, httpErr.Global.Internal)
		return
	}

	setSessionCookie(wtr, session.SessionID, session.ExpiresAt)
	wtr.WriteHeader(http.StatusCreated)
	json.NewEncoder(wtr).Encode(models.CreateSessionResponse{
		SessionID: session.SessionID,
		ExpiresAt: session.ExpiresAt,
	})
}

// GetSession validates the session cookie and returns session metadata.
func (h *SessionHandler) GetSession(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")

	cookie, err := req.Cookie(SessionCookieName)
	if err != nil {
		httpErr.SendError(wtr, req, httpErr.Global.Unauthorized, httpErr.WithDetail("no session cookie"))
		return
	}

	session, err := h.repo.GetSession(req.Context(), cookie.Value)
	if err != nil {
		httpErr.SendError(wtr, req, httpErr.Global.Unauthorized, httpErr.WithDetail("invalid or expired session"))
		return
	}

	wtr.WriteHeader(http.StatusOK)
	json.NewEncoder(wtr).Encode(session)
}

func setSessionCookie(wtr http.ResponseWriter, sessionID string, expires time.Time) {
	http.SetCookie(wtr, &http.Cookie{
		Name:     SessionCookieName,
		Value:    sessionID,
		Expires:  expires,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}

func sessionTTLDays() int {
	if v := config.GetConfigValue("CHAT_SESSION_TTL_DAYS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return defaultSessionTTLDays
}
