package models

import (
	"net/http"

	httpErr "github.com/rdevitto86/komodo-forge-sdk-go/http/errors"
)

// TODO: add RangeReservation = 45 to komodo-forge-sdk-go/http/errors/ranges.go
// alongside RangeCart = 43 and RangeInventory = 44.
const rangeReservation = 45

// ReservationErrors defines error codes for komodo-reservations-api (45xxx range).
type ReservationErrors struct {
	SlotNotFound      httpErr.ErrorCode
	SlotUnavailable   httpErr.ErrorCode
	InvalidSlotDate   httpErr.ErrorCode
	BookingNotFound   httpErr.ErrorCode
	BookingConflict   httpErr.ErrorCode
	BookingExpired    httpErr.ErrorCode
	InvalidBookingOp  httpErr.ErrorCode
	ScheduleSyncError httpErr.ErrorCode
}

var Err = ReservationErrors{
	SlotNotFound:      httpErr.ErrorCode{ID: httpErr.CodeID(rangeReservation, 1), Status: http.StatusNotFound, Message: "Slot not found"},
	SlotUnavailable:   httpErr.ErrorCode{ID: httpErr.CodeID(rangeReservation, 2), Status: http.StatusConflict, Message: "Slot is no longer available"},
	InvalidSlotDate:   httpErr.ErrorCode{ID: httpErr.CodeID(rangeReservation, 3), Status: http.StatusBadRequest, Message: "Invalid slot date"},
	BookingNotFound:   httpErr.ErrorCode{ID: httpErr.CodeID(rangeReservation, 4), Status: http.StatusNotFound, Message: "Booking not found"},
	BookingConflict:   httpErr.ErrorCode{ID: httpErr.CodeID(rangeReservation, 5), Status: http.StatusConflict, Message: "Booking conflict — slot already taken"},
	BookingExpired:    httpErr.ErrorCode{ID: httpErr.CodeID(rangeReservation, 6), Status: http.StatusGone, Message: "Booking hold has expired"},
	InvalidBookingOp:  httpErr.ErrorCode{ID: httpErr.CodeID(rangeReservation, 7), Status: http.StatusBadRequest, Message: "Invalid booking operation for current status"},
	ScheduleSyncError: httpErr.ErrorCode{ID: httpErr.CodeID(rangeReservation, 8), Status: http.StatusInternalServerError, Message: "Failed to sync schedule data"},
}
