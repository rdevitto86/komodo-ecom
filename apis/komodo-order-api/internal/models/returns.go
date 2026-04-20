package models

// ReturnStatus represents the RMA lifecycle state.
// The existing Return entity uses OrderStatus for its status field; this type
// captures the returns-specific states that are not part of the main order lifecycle.
type ReturnStatus string

const (
	ReturnStatusRequested ReturnStatus = "requested"
	ReturnStatusApproved  ReturnStatus = "approved"
	ReturnStatusReceived  ReturnStatus = "received"
	ReturnStatusProcessed ReturnStatus = "processed"
	ReturnStatusRejected  ReturnStatus = "rejected"
	ReturnStatusCancelled ReturnStatus = "cancelled"
)

// ReturnReason is the customer-supplied reason for returning a line item.
type ReturnReason string

const (
	ReturnReasonDefective      ReturnReason = "defective"
	ReturnReasonWrongItem      ReturnReason = "wrong_item"
	ReturnReasonNotAsDescribed ReturnReason = "not_as_described"
	ReturnReasonChangedMind    ReturnReason = "changed_mind"
	ReturnReasonOther          ReturnReason = "other"
)

// ReturnLineItem describes a single item in an RMA request.
// Distinct from OrderItem — carries the return reason and customer notes.
type ReturnLineItem struct {
	ItemID   string       `json:"item_id"`
	SKU      string       `json:"sku"`
	Quantity int          `json:"quantity"`
	Reason   ReturnReason `json:"reason"`
	Notes    string       `json:"notes,omitempty"`
}

// ReturnList is the paginated response for listing RMA requests.
type ReturnList struct {
	Returns    []Return `json:"returns"`
	NextCursor string   `json:"next_cursor,omitempty"`
}

// CreateReturnRequest is the body for POST /me/orders/returns.
type CreateReturnRequest struct {
	OrderID string           `json:"order_id"`
	Items   []ReturnLineItem `json:"items"`
}

// ApproveReturnRequest is the optional body for PUT /internal/returns/{returnId}/approve.
// If RefundAmountCents is nil the full item value is used.
type ApproveReturnRequest struct {
	RefundAmountCents *int64 `json:"refund_amount_cents,omitempty"`
}

// RejectReturnRequest is the body for PUT /internal/returns/{returnId}/reject.
type RejectReturnRequest struct {
	Reason string `json:"reason"`
}
