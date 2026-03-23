package models

import (
	"net/http"

	httpErr "komodo-forge-sdk-go/http/errors"
)

// 21xxx — komodo-entitlements-api (see forge-sdk ranges.go)
type EntitlementsAPIErrors struct {
	NotFound   httpErr.ErrorCode
	Expired    httpErr.ErrorCode
	NotGranted httpErr.ErrorCode
}

var Err = EntitlementsAPIErrors{
	NotFound:   httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeEntitlements, 1), Status: http.StatusNotFound, Message: "Entitlement not found"},
	Expired:    httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeEntitlements, 2), Status: http.StatusForbidden, Message: "Entitlement expired"},
	NotGranted: httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeEntitlements, 3), Status: http.StatusForbidden, Message: "Feature not entitled"},
}
