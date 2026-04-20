package models

import (
	"net/http"

	httpErr "github.com/rdevitto86/komodo-forge-sdk-go/http/errors"
)

// 91xxx — komodo-loyalty-api reviews domain (see forge-sdk ranges.go RangeReviews).
type ReviewsErrors struct {
	NotFound        httpErr.ErrorCode
	AlreadyReviewed httpErr.ErrorCode
	NotEligible     httpErr.ErrorCode
}

var ReviewErr = ReviewsErrors{
	NotFound:        httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeReviews, 1), Status: http.StatusNotFound, Message: "Review not found"},
	AlreadyReviewed: httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeReviews, 2), Status: http.StatusConflict, Message: "Already reviewed"},
	NotEligible:     httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeReviews, 3), Status: http.StatusForbidden, Message: "Not eligible to review — purchase required"},
}
