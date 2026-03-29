package repository

import (
	"context"

	"komodo-order-reservations-api/pkg/v1/models"
)

// GetAvailableSlots returns all available slots from the DynamoDB read model.
//
// TODO(db): implement DynamoDB query against the slots table.
// Table design (proposed):
//   - PK: TechnicianID (string)
//   - SK: SlotDateTime (ISO8601 string for range queries)
//   - GSI: SlotsByDate — PK: Date (YYYY-MM-DD), SK: SlotDateTime
//     (enables customer-facing queries by date without knowing technician IDs)
//   - Attributes: Zone, DurationMinutes, IsAvailable, SourceScheduleID
//
// TODO(db): add filter expression for IsAvailable=true
// TODO(db): support pagination (LastEvaluatedKey → cursor token in response)
// TODO(sync): this table is a read model populated by the external schedule system.
// Wire up the push endpoint (POST /internal/slots/sync) once events-api is ready.
func GetAvailableSlots(ctx context.Context) ([]models.Slot, error) {
	// TODO: initialize DynamoDB client (see forge SDK aws/dynamodb when available)
	// TODO: scan or query slots table with IsAvailable=true filter
	return []models.Slot{}, nil
}

// GetSlotsByDate returns available slots for a given date string (YYYY-MM-DD).
//
// TODO(db): query the SlotsByDate GSI with PK=date, filter IsAvailable=true
// TODO(db): optionally accept zone as a secondary filter
func GetSlotsByDate(ctx context.Context, date string) ([]models.Slot, error) {
	// TODO: query SlotsByDate GSI
	return []models.Slot{}, nil
}

// UpsertSlot writes or updates a slot record in the DynamoDB read model.
// Called by the schedule sync endpoint when the external schedule changes.
//
// TODO(db): use DynamoDB PutItem with conditional expression to avoid overwriting
// a slot that has an active HELD or CONFIRMED booking.
// TODO(sync): call this from the schedule push handler (POST /internal/slots/sync)
func UpsertSlot(ctx context.Context, slot models.Slot) error {
	// TODO: implement PutItem
	return nil
}
