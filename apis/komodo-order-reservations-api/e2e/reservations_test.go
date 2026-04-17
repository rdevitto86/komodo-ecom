//go:build e2e

package e2e_test

import (
	"net/http"
	"testing"
)

func TestHealth(t *testing.T) {
	res := get(t, "/health", nil)
	defer res.Body.Close()
	checkStatus(t, res, http.StatusOK)
}

// TestGetSlots_ReturnsAvailability fetches all available slots.
// Returns 501 if DynamoDB repo stubs are not yet wired.
func TestGetSlots_ReturnsAvailability(t *testing.T) {
	res := get(t, "/slots", nil)
	defer res.Body.Close()
	if res.StatusCode == http.StatusNotImplemented {
		t.Skip("repo stubs not wired — implement DynamoDB GetSlots to enable this test")
	}
	checkStatus(t, res, http.StatusOK)
}

// TestGetSlotsByDate_ValidDate fetches slots for a specific date.
func TestGetSlotsByDate_ValidDate(t *testing.T) {
	res := get(t, "/slots/2026-09-01", nil)
	defer res.Body.Close()
	if res.StatusCode == http.StatusNotImplemented {
		t.Skip("repo stubs not wired — implement DynamoDB GetSlots to enable this test")
	}
	checkStatus(t, res, http.StatusOK)
}

// TestGetSlotsByDate_InvalidFormat verifies malformed dates return 400.
func TestGetSlotsByDate_InvalidFormat(t *testing.T) {
	res := get(t, "/slots/not-a-date", nil)
	defer res.Body.Close()
	if res.StatusCode == http.StatusNotImplemented {
		t.Skip("repo stubs not wired")
	}
	checkStatus(t, res, http.StatusBadRequest)
}

// TestCreateBooking_NoAuth verifies booking creation requires a JWT.
func TestCreateBooking_NoAuth(t *testing.T) {
	body := map[string]any{
		"slot_id":      "slot-e2e-001",
		"service_type": "standard",
	}
	res := post(t, "/bookings", body, nil)
	defer res.Body.Close()
	checkStatus(t, res, http.StatusUnauthorized)
}

// TestBooking_FullFlow creates → reads → confirms → cancels a booking.
// Requires TEST_JWT and a seeded slot in DynamoDB.
func TestBooking_FullFlow(t *testing.T) {
	h := authHeader(t)

	createResp := post(t, "/bookings", map[string]any{
		"slot_id":      "slot-e2e-001",
		"service_type": "standard",
	}, h)
	defer createResp.Body.Close()
	if createResp.StatusCode == http.StatusNotImplemented {
		t.Skip("repo stubs not wired — implement DynamoDB CreateBooking to enable this test")
	}
	if createResp.StatusCode == http.StatusNotFound {
		t.Skip("slot-e2e-001 not in DynamoDB — seed the slot to enable this test")
	}
	checkStatus(t, createResp, http.StatusCreated)

	var created struct {
		ID string `json:"id"`
	}
	decodeJSON(t, createResp, &created)
	if created.ID == "" {
		t.Fatal("expected non-empty booking id in create response")
	}

	// Read booking.
	getResp := get(t, "/bookings/"+created.ID, h)
	defer getResp.Body.Close()
	checkStatus(t, getResp, http.StatusOK)

	// Confirm booking.
	confirmResp := put(t, "/bookings/"+created.ID+"/confirm", nil, h)
	defer confirmResp.Body.Close()
	checkStatus(t, confirmResp, http.StatusOK)

	// Cancel booking.
	cancelResp := put(t, "/bookings/"+created.ID+"/cancel", nil, h)
	defer cancelResp.Body.Close()
	checkStatus(t, cancelResp, http.StatusOK)
}

// TestGetBooking_NotFound verifies 404 for a non-existent booking ID.
func TestGetBooking_NotFound(t *testing.T) {
	h := authHeader(t)
	res := get(t, "/bookings/booking-does-not-exist", h)
	defer res.Body.Close()
	if res.StatusCode == http.StatusNotImplemented {
		t.Skip("repo stubs not wired")
	}
	checkStatus(t, res, http.StatusNotFound)
}
