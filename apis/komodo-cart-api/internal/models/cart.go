package models

import "time"

// Cart is the API response shape for both authenticated and guest carts.
// SubtotalCents and ItemCount are always computed at read time — never stored.
type Cart struct {
	ID            string     `json:"id"`
	UserID        string     `json:"user_id,omitempty"`
	Items         []CartItem `json:"items"`
	SubtotalCents int        `json:"subtotal_cents"`
	ItemCount     int        `json:"item_count"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// CartItem is a snapshot of product data at the time it was added to the cart.
type CartItem struct {
	ItemID         string `json:"item_id"`
	SKU            string `json:"sku"`
	Name           string `json:"name"`
	Quantity       int    `json:"quantity"`
	UnitPriceCents int    `json:"unit_price_cents"`
	ImageURL       string `json:"image_url,omitempty"`
}

// AddItemRequest is the request body for adding an item to a cart.
type AddItemRequest struct {
	ItemID   string `json:"item_id"`
	SKU      string `json:"sku"`
	Quantity int    `json:"quantity"`
}

// UpdateItemRequest is the request body for updating item quantity.
// Quantity == 0 means remove the item.
type UpdateItemRequest struct {
	Quantity int `json:"quantity"`
}

// MergeCartRequest is the request body for merging a guest cart into an authenticated cart.
type MergeCartRequest struct {
	GuestCartID string `json:"guest_cart_id"`
}

// CheckoutResponse is returned by POST /me/cart/checkout.
type CheckoutResponse struct {
	CheckoutToken string     `json:"checkout_token"`
	ExpiresAt     time.Time  `json:"expires_at"`
	Items         []CartItem `json:"items"`
}
