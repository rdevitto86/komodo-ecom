package models

import (
	"net/http"

	httpErr "komodo-forge-sdk-go/http/errors"
)

// 60xxx — komodo-shop-items-api (see forge-sdk ranges.go)
type ShopItemErrors struct {
	ItemNotFound         httpErr.ErrorCode
	InventoryUnavailable httpErr.ErrorCode
	InvalidSKU           httpErr.ErrorCode
	SuggestionFailed     httpErr.ErrorCode
	StorageError         httpErr.ErrorCode
}

var Err = ShopItemErrors{
	ItemNotFound:         httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeShopItem, 1), Status: http.StatusNotFound, Message: "Shop item not found"},
	InventoryUnavailable: httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeShopItem, 2), Status: http.StatusServiceUnavailable, Message: "Inventory data unavailable"},
	InvalidSKU:           httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeShopItem, 3), Status: http.StatusBadRequest, Message: "Invalid SKU"},
	SuggestionFailed:     httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeShopItem, 4), Status: http.StatusInternalServerError, Message: "Failed to generate suggestions"},
	StorageError:         httpErr.ErrorCode{ID: httpErr.CodeID(httpErr.RangeShopItem, 5), Status: http.StatusInternalServerError, Message: "Storage retrieval error"},
}
