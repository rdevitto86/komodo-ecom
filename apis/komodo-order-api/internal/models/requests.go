package models

// PlaceOrderRequest is the body for POST /me/orders.
// checkoutToken is issued by cart-api after a successful /me/cart/checkout call.
type PlaceOrderRequest struct {
	CheckoutToken string `json:"checkoutToken"`
}
