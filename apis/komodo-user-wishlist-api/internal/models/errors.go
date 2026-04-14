package models

import (
	"net/http"

	httpErr "github.com/rdevitto86/komodo-forge-sdk-go/http/errors"
)

// 32xxx — komodo-user-wishlist-api
// TODO: register RangeWishlist = 32 in forge-sdk ranges.go before production
const rangeWishlist = 32

type WishlistAPIErrors struct {
	NotFound        httpErr.ErrorCode
	AlreadyExists   httpErr.ErrorCode
	ItemUnavailable httpErr.ErrorCode
}

var Err = WishlistAPIErrors{
	NotFound:        httpErr.ErrorCode{ID: httpErr.CodeID(rangeWishlist, 1), Status: http.StatusNotFound, Message: "Item not found in wishlist"},
	AlreadyExists:   httpErr.ErrorCode{ID: httpErr.CodeID(rangeWishlist, 2), Status: http.StatusConflict, Message: "Item already in wishlist"},
	ItemUnavailable: httpErr.ErrorCode{ID: httpErr.CodeID(rangeWishlist, 3), Status: http.StatusConflict, Message: "Item no longer available in catalog"},
}
