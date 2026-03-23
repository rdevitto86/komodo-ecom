package models

import (
	"net/http"

	httpErr "komodo-forge-sdk-go/http/errors"
)

// 91xxx — komodo-reviews-api (see forge-sdk ranges.go)
type ReviewsAPIErrors struct {
	NotFound        httpErr.ErrorCode
	AlreadyReviewed httpErr.ErrorCode
	NotEligible     httpErr.ErrorCode
}

var Err = ReviewsAPIErrors{
	NotFound:        httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeReviews, 1), Status: http.StatusNotFound, Message: "Review not found"},
	AlreadyReviewed: httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeReviews, 2), Status: http.StatusConflict, Message: "Already reviewed"},
	NotEligible:     httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeReviews, 3), Status: http.StatusForbidden, Message: "Not eligible to review — purchase required"},
}
