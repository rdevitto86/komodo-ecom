package models

import "time"

// WishlistItem represents a single saved item in a user's wishlist.
// Name and price are denormalized snapshots captured at add time.
type WishlistItem struct {
	ItemID     string    `json:"item_id"               dynamodbav:"item_id"`
	SKU        string    `json:"sku"                   dynamodbav:"sku"`
	Name       string    `json:"name"                  dynamodbav:"name"`
	ImageURL   string    `json:"image_url,omitempty"   dynamodbav:"image_url,omitempty"`
	PriceCents int       `json:"price_cents,omitempty" dynamodbav:"price_cents,omitempty"`
	AddedAt    time.Time `json:"added_at"              dynamodbav:"added_at"`
}

// Wishlist is the top-level response for GET /me/wishlist.
type Wishlist struct {
	Items []WishlistItem `json:"items" dynamodbav:"items"`
}

// AddWishlistItemRequest is the request body for POST /me/wishlist/items.
type AddWishlistItemRequest struct {
	ItemID     string `json:"item_id"`
	SKU        string `json:"sku"`
	Name       string `json:"name,omitempty"`
	ImageURL   string `json:"image_url,omitempty"`
	PriceCents int    `json:"price_cents,omitempty"`
}

// ItemAvailability represents the stock status of a single wishlist item.
type ItemAvailability struct {
	ItemID       string `json:"item_id"`
	SKU          string `json:"sku"`
	InStock      bool   `json:"in_stock"`
	AvailableQty int    `json:"available_qty,omitempty"`
}

// WishlistAvailability is the response for GET /me/wishlist/availability.
type WishlistAvailability struct {
	Items []ItemAvailability `json:"items"`
}

// MoveToCartRequest is the request body for POST /me/wishlist/move-to-cart.
type MoveToCartRequest struct {
	ItemIDs []string `json:"item_ids"`
}

// MoveToCartResultItem reports the outcome for a single item in a move-to-cart operation.
type MoveToCartResultItem struct {
	ItemID string `json:"item_id"`
	Moved  bool   `json:"moved"`
	Reason string `json:"reason,omitempty"`
}

// MoveToCartResult is the response for POST /me/wishlist/move-to-cart.
type MoveToCartResult struct {
	Results []MoveToCartResultItem `json:"results"`
}
