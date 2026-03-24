package httperrors

import "fmt"

// Schema range registry. Each service owns a block of 1000 IDs (xxx001–xxx999).
// Register a new service here before defining its error codes — this is the
// single source of truth preventing cross-service ID collisions.
//
// ID format: "<range><3-digit offset>", e.g. RangeUser + offset 1 → "30001"
// Use CodeID(range, offset) to construct IDs programmatically, or hardcode the
// string directly — both are fine as long as the range block is respected.
const (
	RangeGlobal         = 10 // forge-sdk-go — generic HTTP errors
	RangeDB             = 11 // forge-sdk-go — database errors
	RangeAuth           = 20 // forge-sdk-go + komodo-auth-api
	RangeEntitlements   = 21 // komodo-entitlements-api
	RangeFeatures       = 22 // komodo-features-api
	RangeUser           = 30 // komodo-user-api
	RangeAddress        = 31 // komodo-address-api
	RangeOrder          = 40 // komodo-order-api
	RangeOrderItem      = 41 // komodo-order-api (line items)
	RangeReturns        = 42 // komodo-returns-api
	RangeCart           = 43 // komodo-cart-api
	RangeInventory      = 44 // komodo-inventory-api
	RangePayment        = 50 // komodo-payments-api
	RangeShopItem       = 60 // komodo-shop-items-api
	RangeSearch         = 61 // komodo-search-api
	RangeCommunications = 70 // komodo-communications-api
	RangeEvents         = 71 // komodo-event-bus-api
	RangeAnalytics      = 80 // reserved — analytics range
	RangeSupport        = 81 // komodo-support-api
	RangeLoyalty        = 90 // komodo-loyalty-api
	RangeReviews        = 91 // komodo-reviews-api
)

// CodeID constructs a string error code ID from a range root and a 1-based offset.
// Example: CodeID(RangeUser, 1) → "30001"
func CodeID(rangeRoot, offset int) string { return fmt.Sprintf("%d%03d", rangeRoot, offset) }
