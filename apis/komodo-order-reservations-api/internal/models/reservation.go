package models

import "time"

// BookingStatus represents the lifecycle state of a reservation.
type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "PENDING"   // awaiting scheduling (post-order, option B)
	BookingStatusHeld      BookingStatus = "HELD"      // slot reserved during checkout window (option A)
	BookingStatusConfirmed BookingStatus = "CONFIRMED" // payment confirmed, booking locked
	BookingStatusCancelled BookingStatus = "CANCELLED" // cancelled by customer or system
	BookingStatusReleased  BookingStatus = "RELEASED"  // hold expired or order abandoned
)

// Slot represents a single bookable time window for a technician.
type Slot struct {
	TechnicianID     string    `json:"technician_id"`
	SlotDateTime     time.Time `json:"slot_datetime"`
	DurationMinutes  int       `json:"duration_minutes"`
	Zone             string    `json:"zone"`
	IsAvailable      bool      `json:"is_available"`
	SourceScheduleID string    `json:"source_schedule_id"` // reference to external schedule record
}

// Booking represents a customer reservation against a technician slot.
type Booking struct {
	BookingID    string        `json:"booking_id"`
	CustomerID   string        `json:"customer_id"`
	TechnicianID string        `json:"technician_id"`
	ServiceSKU   string        `json:"service_sku"`
	SlotDateTime time.Time     `json:"slot_datetime"`
	Status       BookingStatus `json:"status"`
	CreatedAt    time.Time     `json:"created_at"`
	ExpiresAt    *time.Time    `json:"expires_at,omitempty"` // populated for HELD status (checkout TTL)
	ConfirmedAt  *time.Time    `json:"confirmed_at,omitempty"`
	CancelledAt  *time.Time    `json:"cancelled_at,omitempty"`
	OrderID      string        `json:"order_id,omitempty"` // set after order submission
}

// CreateBookingRequest is the payload for POST /bookings.
type CreateBookingRequest struct {
	TechnicianID string    `json:"technician_id"`
	ServiceSKU   string    `json:"service_sku"`
	SlotDateTime time.Time `json:"slot_datetime"`
	OrderID      string    `json:"order_id,omitempty"` // present for option B (post-order scheduling)
}
