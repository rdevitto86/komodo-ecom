package models

import (
	"net/http"

	httpErr "komodo-forge-sdk-go/http/errors"
)

// 81xxx — komodo-support-api (see forge-sdk ranges.go)
type SupportAPIErrors struct {
	TicketNotFound     httpErr.ErrorCode
	InvalidTicketState httpErr.ErrorCode
	SessionNotFound    httpErr.ErrorCode
	SessionExpired     httpErr.ErrorCode
	ChatError          httpErr.ErrorCode
	EscalationFailed   httpErr.ErrorCode
}

var Err = SupportAPIErrors{
	TicketNotFound:     httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeSupport, 1), Status: http.StatusNotFound, Message: "Support ticket not found"},
	InvalidTicketState: httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeSupport, 2), Status: http.StatusConflict, Message: "Invalid ticket state transition"},
	SessionNotFound:    httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeSupport, 3), Status: http.StatusNotFound, Message: "Chat session not found"},
	SessionExpired:     httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeSupport, 4), Status: http.StatusUnauthorized, Message: "Chat session expired"},
	ChatError:          httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeSupport, 5), Status: http.StatusInternalServerError, Message: "Chat processing error"},
	EscalationFailed:   httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeSupport, 6), Status: http.StatusInternalServerError, Message: "Failed to escalate chat"},
}
