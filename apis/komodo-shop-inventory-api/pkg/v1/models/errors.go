package models

import (
	"net/http"

	httpErr "github.com/rdevitto86/komodo-forge-sdk-go/http/errors"
)

// 44xxx — komodo-inventory-api (see forge-sdk ranges.go)
type InventoryAPIErrors struct {
	InsufficientStock httpErr.ErrorCode
	SKUNotFound       httpErr.ErrorCode
	HoldNotFound      httpErr.ErrorCode
	HoldExpired       httpErr.ErrorCode
}

var Err = InventoryAPIErrors{
	InsufficientStock: httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeInventory, 1), Status: http.StatusConflict, Message: "Insufficient stock"},
	SKUNotFound:       httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeInventory, 2), Status: http.StatusNotFound, Message: "SKU not found"},
	HoldNotFound:      httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeInventory, 3), Status: http.StatusNotFound, Message: "Stock hold not found"},
	HoldExpired:       httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeInventory, 4), Status: http.StatusGone, Message: "Stock hold expired"},
}
