package handlers

import (
	"encoding/json"
	"net/http"

	httpErr "github.com/rdevitto86/komodo-forge-sdk-go/http/errors"
	logger "github.com/rdevitto86/komodo-forge-sdk-go/logging/runtime"

	"komodo-order-reservations-api/internal/repository"
	"komodo-order-reservations-api/internal/models"
)

// GetAvailableSlots returns all available slots optionally filtered by date and zone.
// GET /slots?date=YYYY-MM-DD&zone=<zone>
func GetAvailableSlots(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")

	// TODO: extract and validate query params (date, zone, technician_id)
	// Use req.URL.Query() or forge SDK request.GetQueryParams when available.
	// Validate date format: YYYY-MM-DD. Return 400 on bad format.

	slots, err := repository.GetAvailableSlots(req.Context())
	if err != nil {
		logger.Error("failed to fetch available slots", err)
		httpErr.SendError(wtr, req, models.Err.SlotNotFound, httpErr.WithDetail("could not retrieve slots"))
		return
	}

	wtr.WriteHeader(http.StatusOK)
	json.NewEncoder(wtr).Encode(slots)
}

// GetSlotsByDate returns available slots for a specific date.
// GET /slots/{date}
func GetSlotsByDate(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")

	date := req.PathValue("date")
	if date == "" {
		httpErr.SendError(wtr, req, models.Err.InvalidSlotDate, httpErr.WithDetail("date path parameter is required"))
		return
	}

	// TODO: validate date format (YYYY-MM-DD), return 400 on bad format
	// TODO: optionally filter by zone query param

	slots, err := repository.GetSlotsByDate(req.Context(), date)
	if err != nil {
		logger.Error("failed to fetch slots for date: "+date, err)
		httpErr.SendError(wtr, req, models.Err.SlotNotFound, httpErr.WithDetail("no slots found for date: "+date))
		return
	}

	wtr.WriteHeader(http.StatusOK)
	json.NewEncoder(wtr).Encode(slots)
}
