package relay

import (
	"strings"
	"time"

	komodoEvents "github.com/rdevitto86/komodo-forge-sdk-go/events"
)

// EventEnvelope is the HTTP request body for POST /events.
// Fields mirror komodo-forge-sdk-go/events.Event but omit entity_id/entity_type,
// which are CDC-path concerns. Schema versioning is enforced from day one —
// adding fields is backwards-compatible; removing or renaming fields requires a version bump.
type EventEnvelope struct {
	ID         string                 `json:"id"`
	Type       komodoEvents.EventType `json:"type"`
	Version    string                 `json:"version"`
	Source     komodoEvents.Source    `json:"source"`
	OccurredAt time.Time              `json:"occurred_at"`
	Payload    map[string]any         `json:"payload"`
}

// KnownEventTypes is the authoritative set of event types this bus accepts.
// Must be kept in sync with komodo-forge-sdk-go/events/envelope.go constants.
var KnownEventTypes = map[komodoEvents.EventType]struct{}{
	komodoEvents.EventOrderCreated:       {},
	komodoEvents.EventOrderStatusUpdated: {},
	komodoEvents.EventOrderCancelled:     {},
	komodoEvents.EventOrderFulfilled:     {},
	komodoEvents.EventUserCreated:        {},
	komodoEvents.EventUserProfileUpdated: {},
	komodoEvents.EventUserDeleted:        {},
	komodoEvents.EventPaymentInitiated:   {},
	komodoEvents.EventPaymentSucceeded:   {},
	komodoEvents.EventPaymentFailed:      {},
	komodoEvents.EventPaymentRefunded:    {},
	komodoEvents.EventCartCheckedOut:     {},
	komodoEvents.EventInventoryReserved:  {},
	komodoEvents.EventInventoryReleased:  {},
}

// domainFromType extracts the domain segment from a "<domain>.<action>" event type.
// "order.placed" → "order", "payment.failed" → "payment".
// Falls back to the full string if no dot is present.
func domainFromType(eventType string) string {
	if idx := strings.IndexByte(eventType, '.'); idx > 0 {
		return eventType[:idx]
	}
	return eventType
}
