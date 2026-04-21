package models

import "fmt"

// OrderType identifies the nature of an order entity.
type OrderType string

const (
	OrderTypeOrder    OrderType = "ORD"
	OrderTypeReturn   OrderType = "RTN"
	OrderTypeExchange OrderType = "EXC"
)

// OrderStatus represents the lifecycle state of an order.
type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusConfirmed  OrderStatus = "confirmed"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
	OrderStatusRefunded   OrderStatus = "refunded"
)

// OrderNumber is the machine identity for an order entity.
// Format: {type}-{zero-padded sequence}
// Examples: ORD-001234, RTN-000089, EXC-000017
//
// This is an internal PK — never used as a display ID in customer-facing comms.
// Use DisplayID() for customer-facing output.
type OrderNumber struct {
	Type     OrderType
	Sequence uint64
}

// String returns the internal identifier, e.g. "ORD-001234".
func (n OrderNumber) String() string {
	return fmt.Sprintf("%s-%06d", n.Type, n.Sequence)
}

// DisplayID returns the customer-facing label for an order entity.
//
// For the root order it is identical to the internal number ("001234").
// For derivatives the parent sequence is used as the base, with a type
// suffix and a 1-based child index:
//
//	root order  → "001234"
//	first RTN   → "001234-R1"
//	second RTN  → "001234-R2"
//	first EXC   → "001234-X1"
//
// The display ID is a formatting convention — it is never stored as a key
// or used for database lookups.
func (n OrderNumber) DisplayID(parentSeq uint64, childIndex int) string {
	switch n.Type {
	case OrderTypeReturn:
		return fmt.Sprintf("%06d-R%d", parentSeq, childIndex)
	case OrderTypeExchange:
		return fmt.Sprintf("%06d-X%d", parentSeq, childIndex)
	default:
		return fmt.Sprintf("%06d", n.Sequence)
	}
}

// Order is the root purchase entity.
type Order struct {
	ID          string       `json:"id"`              // internal PK, e.g. "ORD-001234"
	DisplayID   string       `json:"displayId"`       // customer-facing label, e.g. "001234"
	UserID      string       `json:"userId"`          // USER#<userId> for registered, GUEST#<uuid> for guests
	Email       string       `json:"email,omitempty"` // universal key; required for guest order lookup
	Status      OrderStatus  `json:"status"`
	Items       []OrderItem  `json:"items"`
	Address     OrderAddress `json:"address"`
	Payment     OrderPayment `json:"payment"`
	Totals      OrderTotals  `json:"totals"`
	CreatedAt   string       `json:"createdAt"`
	UpdatedAt   string       `json:"updatedAt"`
}

// Return is a derivative of an Order representing a customer return request.
// It has its own independent sequence and internal ID (e.g. "RTN-000089").
// The relationship back to the originating order is carried by ParentOrderID.
// The customer-facing DisplayID is derived at render time: "001234-R1".
type Return struct {
	ID            string      `json:"id"`            // internal PK, e.g. "RTN-000089"
	DisplayID     string      `json:"displayId"`     // customer-facing label, e.g. "001234-R1"
	ParentOrderID string      `json:"parentOrderId"` // FK → Order.ID
	UserID        string      `json:"userId"`
	Status        OrderStatus `json:"status"`
	Items         []OrderItem `json:"items"`
	Reason        string      `json:"reason"`
	CreatedAt     string      `json:"createdAt"`
	UpdatedAt     string      `json:"updatedAt"`
}

// Exchange is a derivative of an Order representing a product exchange.
// Same identity model as Return — independent sequence, explicit FK to parent.
// Customer-facing DisplayID example: "001234-X1".
type Exchange struct {
	ID            string      `json:"id"`            // internal PK, e.g. "EXC-000017"
	DisplayID     string      `json:"displayId"`     // customer-facing label, e.g. "001234-X1"
	ParentOrderID string      `json:"parentOrderId"` // FK → Order.ID
	UserID        string      `json:"userId"`
	Status        OrderStatus `json:"status"`
	ReturnItems   []OrderItem `json:"returnItems"`
	ExchangeItems []OrderItem `json:"exchangeItems"`
	PriceDelta    float64     `json:"priceDelta"`
	CreatedAt     string      `json:"createdAt"`
	UpdatedAt     string      `json:"updatedAt"`
}

// OrderItem is a line item within any order entity.
type OrderItem struct {
	ID        string  `json:"id"`
	ProductID string  `json:"productId"`
	VariantID string  `json:"variantId,omitempty"`
	SKU       string  `json:"sku"`
	Name      string  `json:"name"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unitPrice"`
	Total     float64 `json:"total"`
}

type OrderAddress struct {
	Line1      string `json:"line1"`
	Line2      string `json:"line2,omitempty"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postalCode"`
	Country    string `json:"country"`
}

type OrderPayment struct {
	Method        string  `json:"method"`
	TransactionID string  `json:"transactionId,omitempty"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
}

type OrderTotals struct {
	Subtotal float64 `json:"subtotal"`
	Tax      float64 `json:"tax"`
	Shipping float64 `json:"shipping"`
	Discount float64 `json:"discount,omitempty"`
	Total    float64 `json:"total"`
}
