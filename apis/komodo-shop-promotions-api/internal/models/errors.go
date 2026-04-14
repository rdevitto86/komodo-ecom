package models

import (
	"net/http"

	httpErr "github.com/rdevitto86/komodo-forge-sdk-go/http/errors"
)

// 62xxx — komodo-shop-promotions-api
// TODO: register RangePromotions = 62 in forge-sdk ranges.go before production
const rangePromotions = 62

type PromotionsAPIErrors struct {
	NotFound             httpErr.ErrorCode
	Expired              httpErr.ErrorCode
	NotEligible          httpErr.ErrorCode
	AlreadyUsed          httpErr.ErrorCode
	MaxRedemptionsReached httpErr.ErrorCode
	Inactive             httpErr.ErrorCode
}

var Err = PromotionsAPIErrors{
	NotFound:              httpErr.ErrorCode{ID: httpErr.CodeID(rangePromotions, 1), Status: http.StatusNotFound, Message: "Promotion not found"},
	Expired:               httpErr.ErrorCode{ID: httpErr.CodeID(rangePromotions, 2), Status: http.StatusGone, Message: "Promotion has expired"},
	NotEligible:           httpErr.ErrorCode{ID: httpErr.CodeID(rangePromotions, 3), Status: http.StatusConflict, Message: "Cart does not meet promotion conditions"},
	AlreadyUsed:           httpErr.ErrorCode{ID: httpErr.CodeID(rangePromotions, 4), Status: http.StatusConflict, Message: "Promotion already redeemed by this user"},
	MaxRedemptionsReached: httpErr.ErrorCode{ID: httpErr.CodeID(rangePromotions, 5), Status: http.StatusConflict, Message: "Promotion redemption limit reached"},
	Inactive:              httpErr.ErrorCode{ID: httpErr.CodeID(rangePromotions, 6), Status: http.StatusConflict, Message: "Promotion is not active"},
}
