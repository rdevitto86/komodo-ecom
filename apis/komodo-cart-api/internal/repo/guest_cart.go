package repo

import (
	"encoding/json"
	"fmt"
	"komodo-cart-api/internal/models"

	"github.com/rdevitto86/komodo-forge-sdk-go/aws/elasticache"
)

const guestCartKeyPrefix = "cart:guest:"

// guestCartRecord is the Redis envelope for a guest cart.
// Storing session_id alongside the cart makes session validation atomic with the read.
type GuestCartRecord struct {
	SessionID string      `json:"session_id"`
	Cart      models.Cart `json:"cart"`
}

// CreateGuestCart writes a new guest cart envelope to Redis with the given TTL.
func CreateGuestCart(cart models.Cart, sessionID string, ttlSecs int64) error {
	record := GuestCartRecord{SessionID: sessionID, Cart: cart}
	data, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("repo.CreateGuestCart: marshal: %w", err)
	}
	if err := elasticache.Set(guestCartKeyPrefix+cart.ID, string(data), ttlSecs); err != nil {
		return fmt.Errorf("repo.CreateGuestCart: set: %w", err)
	}
	return nil
}

// GetGuestCart retrieves a guest cart envelope from Redis.
// Returns (nil, nil) if the key does not exist (expired or never created).
func GetGuestCart(cartID string) (*GuestCartRecord, error) {
	raw, err := elasticache.Get(guestCartKeyPrefix + cartID)
	if err != nil {
		return nil, fmt.Errorf("repo.GetGuestCart: get: %w", err)
	}
	if raw == "" {
		return nil, nil
	}
	var record GuestCartRecord
	if err := json.Unmarshal([]byte(raw), &record); err != nil {
		return nil, fmt.Errorf("repo.GetGuestCart: unmarshal: %w", err)
	}
	return &record, nil
}

// SaveGuestCart overwrites the guest cart envelope in Redis, refreshing the TTL.
func SaveGuestCart(record *GuestCartRecord, ttlSecs int64) error {
	data, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("repo.SaveGuestCart: marshal: %w", err)
	}
	if err := elasticache.Set(guestCartKeyPrefix+record.Cart.ID, string(data), ttlSecs); err != nil {
		return fmt.Errorf("repo.SaveGuestCart: set: %w", err)
	}
	return nil
}

// DeleteGuestCart removes a guest cart key from Redis.
func DeleteGuestCart(cartID string) error {
	if err := elasticache.Delete(guestCartKeyPrefix + cartID); err != nil {
		return fmt.Errorf("repo.DeleteGuestCart: delete: %w", err)
	}
	return nil
}
