package models

import (
	"net/http"

	httpErr "komodo-forge-sdk-go/http/errors"
)

// 22xxx — komodo-features-api (see forge-sdk ranges.go)
type FeaturesAPIErrors struct {
	NotFound httpErr.ErrorCode
	Disabled httpErr.ErrorCode
}

var Err = FeaturesAPIErrors{
	NotFound: httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeFeatures, 1), Status: http.StatusNotFound, Message: "Feature flag not found"},
	Disabled: httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeFeatures, 2), Status: http.StatusForbidden, Message: "Feature disabled"},
}
