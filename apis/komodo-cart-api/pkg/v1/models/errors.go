package models

import (
	"net/http"

	httpErr "komodo-forge-sdk-go/http/errors"
)

// 43xxx — komodo-cart-api (see forge-sdk ranges.go)
type CartAPIErrors struct {
	NotFound       httpErr.ErrorCode
	ItemNotInCart  httpErr.ErrorCode
	Expired        httpErr.ErrorCode
	CheckoutFailed httpErr.ErrorCode
}

var Err = CartAPIErrors{
	NotFound:       httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeCart, 1), Status: http.StatusNotFound, Message: "Cart not found"},
	ItemNotInCart:  httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeCart, 2), Status: http.StatusNotFound, Message: "Item not in cart"},
	Expired:        httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeCart, 3), Status: http.StatusGone, Message: "Cart expired"},
	CheckoutFailed: httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeCart, 4), Status: http.StatusConflict, Message: "Checkout failed"},
}
