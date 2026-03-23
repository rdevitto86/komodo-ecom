package models

import (
	"net/http"

	httpErr "komodo-forge-sdk-go/http/errors"
)

// 31xxx — komodo-address-api (see forge-sdk ranges.go)
type AddressAPIErrors struct {
	InvalidFormat     httpErr.ErrorCode
	ProviderTimeout   httpErr.ErrorCode
	ProviderError     httpErr.ErrorCode
	UnsupportedRegion httpErr.ErrorCode
}

var Err = AddressAPIErrors{
	InvalidFormat:     httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeAddress, 1), Status: http.StatusUnprocessableEntity, Message: "Invalid address format"},
	ProviderTimeout:   httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeAddress, 2), Status: http.StatusGatewayTimeout, Message: "Address provider timed out"},
	ProviderError:     httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeAddress, 3), Status: http.StatusBadGateway, Message: "Address provider error"},
	UnsupportedRegion: httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeAddress, 4), Status: http.StatusUnprocessableEntity, Message: "Region not supported"},
}
