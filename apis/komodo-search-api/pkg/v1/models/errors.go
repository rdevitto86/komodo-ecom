package models

import (
	"net/http"

	httpErr "komodo-forge-sdk-go/http/errors"
)

// 61xxx — komodo-search-api (see forge-sdk ranges.go)
type SearchAPIErrors struct {
	SearchFailed     httpErr.ErrorCode
	InvalidQuery     httpErr.ErrorCode
	IndexUnavailable httpErr.ErrorCode
}

var Err = SearchAPIErrors{
	SearchFailed:     httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeSearch, 1), Status: http.StatusInternalServerError, Message: "Search failed"},
	InvalidQuery:     httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeSearch, 2), Status: http.StatusBadRequest, Message: "Invalid search query"},
	IndexUnavailable: httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeSearch, 3), Status: http.StatusServiceUnavailable, Message: "Search index unavailable"},
}
