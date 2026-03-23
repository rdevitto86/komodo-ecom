package models

import (
	"net/http"

	httpErr "komodo-forge-sdk-go/http/errors"
)

// 40xxx — komodo-order-api orders (see forge-sdk ranges.go)
// 41xxx — komodo-order-api line items
type OrderAPIErrors struct {
	NotFound          httpErr.ErrorCode
	AlreadyCancelled  httpErr.ErrorCode
	NotCancellable    httpErr.ErrorCode
	InvalidTransition httpErr.ErrorCode
	ItemNotFound      httpErr.ErrorCode
	ItemUnavailable   httpErr.ErrorCode
	InvalidQuantity   httpErr.ErrorCode
}

var Err = OrderAPIErrors{
	NotFound:          httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeOrder, 1), Status: http.StatusNotFound, Message: "Order not found"},
	AlreadyCancelled:  httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeOrder, 2), Status: http.StatusConflict, Message: "Order already cancelled"},
	NotCancellable:    httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeOrder, 3), Status: http.StatusConflict, Message: "Order cannot be cancelled"},
	InvalidTransition: httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeOrder, 4), Status: http.StatusConflict, Message: "Invalid order state transition"},
	ItemNotFound:      httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeOrderItem, 1), Status: http.StatusNotFound, Message: "Order item not found"},
	ItemUnavailable:   httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeOrderItem, 2), Status: http.StatusConflict, Message: "Item unavailable"},
	InvalidQuantity:   httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeOrderItem, 3), Status: http.StatusBadRequest, Message: "Invalid quantity"},
}
