package events

import (
	"time"

	"github.com/google/uuid"
)

// EventType is the canonical event name in <entity>.<verb> format.
// This is the public contract between services — treat as stable once defined.
type EventType string

const (
	// Order events
	EventOrderCreated       EventType = "order.created"
	EventOrderStatusUpdated EventType = "order.status_updated"
	EventOrderCancelled     EventType = "order.cancelled"
	EventOrderFulfilled     EventType = "order.fulfilled"

	// User events
	EventUserCreated        EventType = "user.created"
	EventUserProfileUpdated EventType = "user.profile_updated"
	EventUserDeleted        EventType = "user.deleted"

	// Payment events
	EventPaymentInitiated EventType = "payment.initiated"
	EventPaymentSucceeded EventType = "payment.succeeded"
	EventPaymentFailed    EventType = "payment.failed"
	EventPaymentRefunded  EventType = "payment.refunded"

	// Cart events
	EventCartCheckedOut EventType = "cart.checked_out"

	// Inventory events
	EventInventoryReserved EventType = "inventory.reserved"
	EventInventoryReleased EventType = "inventory.released"
)

// Source identifies which Komodo service emitted the event.
type Source string

const (
	SourceAuthAPI           Source = "komodo-auth-api"
	SourceUserAPI           Source = "komodo-user-api"
	SourceOrderAPI          Source = "komodo-order-api"
	SourceCartAPI           Source = "komodo-cart-api"
	SourceInventoryAPI      Source = "komodo-inventory-api"
	SourcePaymentsAPI       Source = "komodo-payments-api"
	SourceShopItemsAPI      Source = "komodo-shop-items-api"
	SourceCommunicationsAPI Source = "komodo-communications-api"
)

// EntityType identifies the domain entity at the centre of the event.
type EntityType string

const (
	EntityOrder     EntityType = "order"
	EntityUser      EntityType = "user"
	EntityPayment   EntityType = "payment"
	EntityCart      EntityType = "cart"
	EntityInventory EntityType = "inventory"
	EntityProduct   EntityType = "product"
)

// Event is the canonical business event envelope published to SNS FIFO topics
// and consumed via SQS FIFO queues. The SNS/SQS MessageGroupId must be set to
// EntityID so that all events for the same entity are ordered, while events for
// different entities can be processed in parallel.
type Event struct {
	ID         string         `json:"id"`
	Type       EventType      `json:"type"`
	Source     Source         `json:"source"`
	EntityID   string         `json:"entity_id"`
	EntityType EntityType     `json:"entity_type"`
	OccurredAt time.Time      `json:"occurred_at"`
	Version    string         `json:"version"`
	Payload    map[string]any `json:"payload"`
	// CorrelationID traces an event chain back to the originating HTTP request.
	// Populate from the request's X-Correlation-ID header when available.
	CorrelationID string `json:"correlation_id,omitempty"`
}

// New constructs an Event with a generated ID, current UTC timestamp,
// and version "1". Use WithCorrelation to attach a correlation ID.
func New(
	eventType EventType,
	source Source,
	entityID string,
	entityType EntityType,
	payload map[string]any,
) Event {
	return Event{
		ID:         uuid.NewString(),
		Type:       eventType,
		Source:     source,
		EntityID:   entityID,
		EntityType: entityType,
		OccurredAt: time.Now().UTC(),
		Version:    "1",
		Payload:    payload,
	}
}

// WithCorrelation returns a copy of the event with CorrelationID set.
func (e Event) WithCorrelation(correlationID string) Event {
	e.CorrelationID = correlationID
	return e
}

// TODO: add WithCorrelationFromContext(ctx context.Context) Event that reads
// the correlation ID from ctx via ctxKeys.CORRELATION_ID_KEY and calls
// WithCorrelation — lets CDC and HTTP handlers attach correlation in one call.
