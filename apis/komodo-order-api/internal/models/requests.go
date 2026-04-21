package models

// PlaceOrderRequest is the body for POST /me/orders.
// checkoutToken is issued by cart-api after a successful /me/cart/checkout call.
type PlaceOrderRequest struct {
	CheckoutToken string `json:"checkoutToken"`
}

// UnifiedPlaceOrderRequest is the body for POST /orders.
// Accepts both authenticated (JWT) and guest requests. Email is required when
// no JWT is present; it is ignored if the JWT already carries an identity.
type UnifiedPlaceOrderRequest struct {
	CheckoutToken string `json:"checkoutToken"`
	Email         string `json:"email,omitempty"` // required if no JWT; ignored if JWT present
}

// OrderListResponse is the paginated response for GET /me/orders.
// nextCursor is empty when no further pages exist.
type OrderListResponse struct {
	Orders     []*Order `json:"orders"`
	NextCursor string   `json:"nextCursor,omitempty"`
}
