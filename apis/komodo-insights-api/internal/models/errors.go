package models

import (
	"net/http"

	httpErr "github.com/rdevitto86/komodo-forge-sdk-go/http/errors"
)

// TODO: register httpErr.RangeInsights in komodo-forge-sdk-go/http/errors/ranges.go
// and replace the placeholder range below before going to production.
//
// Proposed range: 11xxx — komodo-insights-api

type InsightsAPIErrors struct {
	SummaryFailed   httpErr.ErrorCode
	ProviderTimeout httpErr.ErrorCode
	EntityNotFound  httpErr.ErrorCode
}

var Err = InsightsAPIErrors{
	SummaryFailed:   httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeGlobal, 1), Status: http.StatusInternalServerError, Message: "Failed to generate summary"},
	ProviderTimeout: httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeGlobal, 2), Status: http.StatusGatewayTimeout, Message: "LLM provider timed out"},
	EntityNotFound:  httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeGlobal, 3), Status: http.StatusNotFound, Message: "Entity not found"},
}
