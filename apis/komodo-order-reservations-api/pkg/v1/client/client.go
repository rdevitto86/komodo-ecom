package client

// TODO: implement external client SDK for komodo-reservations-api.
// Other services (e.g. order-api, cart-api) will use this package to
// check slot availability or read booking status without calling the HTTP API directly.
//
// Planned exports:
//   - GetAvailableSlots(ctx, date, zone) ([]models.Slot, error)
//   - GetBooking(ctx, bookingID) (*models.Booking, error)
