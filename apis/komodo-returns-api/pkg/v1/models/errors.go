package models

import (
	"net/http"

	httpErr "github.com/rdevitto86/komodo-forge-sdk-go/http/errors"
)

// 42xxx — komodo-returns-api (see forge-sdk ranges.go)
type ReturnsAPIErrors struct {
	NotFound         httpErr.ErrorCode
	WindowExpired    httpErr.ErrorCode
	NotEligible      httpErr.ErrorCode
	AlreadyProcessed httpErr.ErrorCode
}

var Err = ReturnsAPIErrors{
	NotFound:         httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeReturns, 1), Status: http.StatusNotFound, Message: "Return not found"},
	WindowExpired:    httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeReturns, 2), Status: http.StatusGone, Message: "Return window expired"},
	NotEligible:      httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeReturns, 3), Status: http.StatusConflict, Message: "Order not eligible for return"},
	AlreadyProcessed: httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeReturns, 4), Status: http.StatusConflict, Message: "Return already processed"},
}
