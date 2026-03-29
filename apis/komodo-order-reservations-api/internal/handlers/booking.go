package handlers

import (
	"encoding/json"
	"net/http"

	httpErr "github.com/rdevitto86/komodo-forge-sdk-go/http/errors"
	logger "github.com/rdevitto86/komodo-forge-sdk-go/logging/runtime"

	"komodo-order-reservations-api/internal/repository"
	"komodo-order-reservations-api/pkg/v1/models"
)

// CreateBooking reserves a slot for a customer.
// POST /bookings
//
// TODO(checkout): decide between two checkout flows before finalizing this handler:
//
//	Option A — Pre-payment hold:
//	  Slot is held immediately (status=HELD) with a short TTL (e.g. 5min).
//	  On successful payment, order-api calls PUT /bookings/{id}/confirm.
//	  On abandonment or TTL expiry, the hold is released automatically (DynamoDB TTL).
//	  Requires: ExpiresAt field, TTL-based auto-release, confirm endpoint wired to order-api.
//
//	Option B — Post-order scheduling:
//	  Booking created after order is submitted (status=PENDING).
//	  Customer selects slot within a post-order window.
//	  No TTL hold; status stays PENDING until customer schedules or cancels.
//	  Requires: order_id in request body, separate scheduling step in order-api flow.
func CreateBooking(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")

	var body models.CreateBookingRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		httpErr.SendError(wtr, req, models.Err.InvalidBookingOp, httpErr.WithDetail("invalid request body"))
		return
	}

	// TODO: extract customer_id from JWT context (use forge SDK http/context USER_ID_KEY)
	// TODO: validate body fields (technician_id, service_sku, slot_datetime non-zero)
	// TODO: verify service_sku exists and is a service type (call shop-items-api client or cache)

	booking, err := repository.CreateBooking(req.Context(), body)
	if err != nil {
		logger.Error("failed to create booking", err)
		httpErr.SendError(wtr, req, models.Err.BookingConflict, httpErr.WithDetail("slot is unavailable or booking conflict"))
		return
	}

	wtr.WriteHeader(http.StatusCreated)
	json.NewEncoder(wtr).Encode(booking)
}

// GetBooking returns a booking by ID.
// GET /bookings/{id}
func GetBooking(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")

	bookingID := req.PathValue("id")
	if bookingID == "" {
		httpErr.SendError(wtr, req, models.Err.BookingNotFound, httpErr.WithDetail("booking id is required"))
		return
	}

	// TODO: enforce ownership — customer_id from JWT must match booking.CustomerID unless internal caller

	booking, err := repository.GetBooking(req.Context(), bookingID)
	if err != nil {
		logger.Warn("booking not found: " + bookingID)
		httpErr.SendError(wtr, req, models.Err.BookingNotFound, httpErr.WithDetail("booking not found: "+bookingID))
		return
	}

	wtr.WriteHeader(http.StatusOK)
	json.NewEncoder(wtr).Encode(booking)
}

// CancelBooking cancels an existing booking.
// PUT /bookings/{id}/cancel
func CancelBooking(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")

	bookingID := req.PathValue("id")
	if bookingID == "" {
		httpErr.SendError(wtr, req, models.Err.BookingNotFound, httpErr.WithDetail("booking id is required"))
		return
	}

	// TODO: enforce ownership check
	// TODO: define cancellation policy — can CONFIRMED bookings be cancelled? Within what window?
	// TODO: if cancellation releases a slot, emit an event so the schedule read model updates

	if err := repository.UpdateBookingStatus(req.Context(), bookingID, models.BookingStatusCancelled); err != nil {
		logger.Error("failed to cancel booking: "+bookingID, err)
		httpErr.SendError(wtr, req, models.Err.InvalidBookingOp, httpErr.WithDetail("cannot cancel booking in current state"))
		return
	}

	wtr.WriteHeader(http.StatusNoContent)
}

// ConfirmBooking transitions a HELD booking to CONFIRMED after successful payment.
// PUT /bookings/{id}/confirm
//
// TODO(checkout): this endpoint is only relevant for Option A (pre-payment hold).
// It should be called by order-api after payment succeeds, not by the customer directly.
// Consider making this internal-only (behind internal middleware or an internal port).
func ConfirmBooking(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")

	bookingID := req.PathValue("id")
	if bookingID == "" {
		httpErr.SendError(wtr, req, models.Err.BookingNotFound, httpErr.WithDetail("booking id is required"))
		return
	}

	// TODO: verify caller is order-api (internal auth / scopes check)
	// TODO: verify booking is in HELD status — reject if already CONFIRMED, CANCELLED, or RELEASED
	// TODO: clear ExpiresAt TTL on confirmation so DynamoDB doesn't auto-release

	if err := repository.UpdateBookingStatus(req.Context(), bookingID, models.BookingStatusConfirmed); err != nil {
		logger.Error("failed to confirm booking: "+bookingID, err)
		httpErr.SendError(wtr, req, models.Err.InvalidBookingOp, httpErr.WithDetail("cannot confirm booking in current state"))
		return
	}

	wtr.WriteHeader(http.StatusNoContent)
}
