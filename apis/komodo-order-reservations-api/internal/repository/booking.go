package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"komodo-order-reservations-api/internal/models"
)

// CreateBooking writes a new booking to DynamoDB using a conditional write to prevent double-booking.
//
// TODO(db): implement DynamoDB TransactWriteItems to atomically:
//   1. Write the Booking item (PK: BookingID)
//   2. Update the Slot item: set IsAvailable=false with condition IsAvailable=true
//      (if condition fails → slot was taken concurrently → return BookingConflict error)
//
// TODO(checkout-A): if Option A (pre-payment hold), set ExpiresAt = now + holdTTL and
// write a DynamoDB TTL attribute so the hold auto-expires. holdTTL should be config-driven.
//
// TODO(checkout-B): if Option B (post-order), set status=PENDING, no TTL needed.
//
// Table design (proposed):
//   - PK: BookingID (UUID)
//   - GSI: BookingsByCustomer — PK: CustomerID, SK: CreatedAt
//     (enables GET /me/bookings without scanning)
//   - GSI: BookingsBySlot — PK: SlotDateTime+TechnicianID, SK: Status
//     (enables slot availability check from booking side)
func CreateBooking(ctx context.Context, req models.CreateBookingRequest) (*models.Booking, error) {
	booking := &models.Booking{
		BookingID:    uuid.New().String(),
		TechnicianID: req.TechnicianID,
		ServiceSKU:   req.ServiceSKU,
		SlotDateTime: req.SlotDateTime,
		Status:       models.BookingStatusPending,
		CreatedAt:    time.Now().UTC(),
		OrderID:      req.OrderID,
		// TODO: set CustomerID from context (JWT USER_ID_KEY)
		// TODO: set ExpiresAt for Option A hold
	}

	// TODO: DynamoDB TransactWriteItems (booking + slot availability update)
	// On ConditionalCheckFailedException → return fmt.Errorf("slot unavailable: %w", ErrConflict)

	_ = fmt.Sprintf("TODO: write booking %s to DynamoDB", booking.BookingID)
	return booking, nil
}

// GetBooking retrieves a booking by ID.
//
// TODO(db): DynamoDB GetItem with PK=BookingID
// TODO: return typed not-found error so handler can map to BookingNotFound error code
func GetBooking(ctx context.Context, bookingID string) (*models.Booking, error) {
	// TODO: implement DynamoDB GetItem
	return nil, fmt.Errorf("GetBooking not implemented")
}

// UpdateBookingStatus updates the status of an existing booking.
//
// TODO(db): DynamoDB UpdateItem with condition expression to enforce valid transitions:
//   HELD → CONFIRMED (confirm after payment)
//   HELD | PENDING | CONFIRMED → CANCELLED (cancellation)
//   HELD → RELEASED (TTL expiry or explicit release)
//
// TODO: emit a status-change event via events-api for downstream consumers (order-api, communications-api)
func UpdateBookingStatus(ctx context.Context, bookingID string, status models.BookingStatus) error {
	// TODO: implement DynamoDB UpdateItem with conditional status transition
	return fmt.Errorf("UpdateBookingStatus not implemented")
}
