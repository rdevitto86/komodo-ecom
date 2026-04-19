package models

import (
	"net/http"

	httpErr "github.com/rdevitto86/komodo-forge-sdk-go/http/errors"
)

// CartError wraps an ErrorCode so it can be returned as an error from service functions.
// Handlers unwrap it back to ErrorCode via the As method for SendError.
type CartError struct {
	Code httpErr.ErrorCode
}

func (e CartError) Error() string { return e.Code.Message }

// 43xxx — komodo-cart-api (see forge-sdk ranges.go)
type CartAPIErrors struct {
	NotFound       CartError
	ItemNotInCart  CartError
	Expired        CartError
	CheckoutFailed CartError
	OutOfStock     CartError
}

var Err = CartAPIErrors{
	NotFound:       CartError{Code: httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeCart, 1), Status: http.StatusNotFound, Message: "Cart not found"}},
	ItemNotInCart:  CartError{Code: httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeCart, 2), Status: http.StatusNotFound, Message: "Item not in cart"}},
	Expired:        CartError{Code: httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeCart, 3), Status: http.StatusGone, Message: "Cart expired"}},
	CheckoutFailed: CartError{Code: httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeCart, 4), Status: http.StatusConflict, Message: "Checkout failed"}},
	OutOfStock:     CartError{Code: httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeCart, 5), Status: http.StatusConflict, Message: "Item out of stock"}},
}
