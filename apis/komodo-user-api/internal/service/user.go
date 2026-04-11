package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"komodo-user-api/internal/repo"
	"komodo-user-api/internal/models"
)

// ErrNotFound is the sentinel returned when a requested resource does not exist.
// Handlers inspect this via errors.Is to map it to HTTP 404.
var ErrNotFound = errors.New("not found")

func GetProfile(ctx context.Context, userID string) (*models.User, error) {
	user, err := repo.GetUser(ctx, userID)
	if err != nil {
		if isNotFound(err) {
			return nil, fmt.Errorf("service.GetProfile: %w", ErrNotFound)
		}
		return nil, fmt.Errorf("service.GetProfile: %w", err)
	}
	return user, nil
}

func CreateUser(ctx context.Context, user *models.User) error {
	if user.UserID == "" || user.Email == "" || user.FirstName == "" || user.LastName == "" {
		return fmt.Errorf("service.CreateUser: %w", errors.New("user_id, email, first_name, and last_name are required"))
	}
	now := time.Now().UTC()
	user.CreatedAt = now
	user.UpdatedAt = now
	if err := repo.CreateUser(ctx, user); err != nil {
		return fmt.Errorf("service.CreateUser: %w", err)
	}
	return nil
}

func UpdateProfile(ctx context.Context, userID string, update *models.User) (*models.User, error) {
	updated, err := repo.UpdateUser(ctx, userID, update)
	if err != nil {
		if isNotFound(err) {
			return nil, fmt.Errorf("service.UpdateProfile: %w", ErrNotFound)
		}
		return nil, fmt.Errorf("service.UpdateProfile: %w", err)
	}
	return updated, nil
}

func DeleteProfile(ctx context.Context, userID string) error {
	if err := repo.DeleteUser(ctx, userID); err != nil {
		return fmt.Errorf("service.DeleteProfile: %w", err)
	}
	return nil
}

func GetAddresses(ctx context.Context, userID string) ([]models.Address, error) {
	addrs, err := repo.GetUserAddresses(ctx, userID)
	if err != nil {
		if isNotFound(err) {
			return nil, fmt.Errorf("service.GetAddresses: %w", ErrNotFound)
		}
		return nil, fmt.Errorf("service.GetAddresses: %w", err)
	}
	return addrs, nil
}

func AddAddress(ctx context.Context, userID string, addr *models.Address) error {
	if addr.AddressID == "" {
		raw := strings.ReplaceAll(uuid.NewString(), "-", "")
		addr.AddressID = "addr_" + raw[:12]
	}
	// TODO: replace with a dedicated repo.CreateAddress once the address sub-item write
	// pattern is defined in data-model.md. CreateUser is the only available write today
	// and does not support address items; this call will return an error until the repo
	// is extended.
	user := &models.User{UserID: userID}
	_ = addr // addr will be wired into the dedicated repo call
	if err := repo.CreateUser(ctx, user); err != nil {
		return fmt.Errorf("service.AddAddress: %w", err)
	}
	return nil
}

func UpdateAddress(ctx context.Context, userID, addressID string, update *models.Address) error {
	// TODO: wire to repo.UpdateAddress once that function exists (data-model.md pending).
	_, err := repo.UpdateUser(ctx, userID, &models.User{UserID: userID})
	if err != nil {
		if isNotFound(err) {
			return fmt.Errorf("service.UpdateAddress: %w", ErrNotFound)
		}
		return fmt.Errorf("service.UpdateAddress: %w", err)
	}
	_ = addressID
	_ = update
	return nil
}

func DeleteAddress(ctx context.Context, userID, addressID string) error {
	// TODO: wire to repo.DeleteAddress once that function exists (data-model.md pending).
	if err := repo.DeleteUser(ctx, userID); err != nil {
		return fmt.Errorf("service.DeleteAddress: %w", err)
	}
	_ = addressID
	return nil
}

func GetPayments(ctx context.Context, userID string) ([]models.PaymentMethod, error) {
	// TODO: wire to repo.GetUserPayments once that function exists (data-model.md pending).
	_ = userID
	return nil, fmt.Errorf("service.GetPayments: not implemented")
}

func UpsertPayment(ctx context.Context, userID string, pm *models.PaymentMethod) error {
	if pm.PaymentID == "" {
		raw := strings.ReplaceAll(uuid.NewString(), "-", "")
		pm.PaymentID = "pay_" + raw[:12]
	}
	// TODO: wire to repo.UpsertPayment once that function exists (data-model.md pending).
	_ = userID
	return fmt.Errorf("service.UpsertPayment: not implemented")
}

func DeletePayment(ctx context.Context, userID, paymentID string) error {
	// TODO: wire to repo.DeletePayment once that function exists (data-model.md pending).
	_ = userID
	_ = paymentID
	return fmt.Errorf("service.DeletePayment: not implemented")
}

func GetPreferences(ctx context.Context, userID string) (*models.Preferences, error) {
	prefs, err := repo.GetUserPreferences(ctx, userID)
	if err != nil {
		if isNotFound(err) {
			return nil, fmt.Errorf("service.GetPreferences: %w", ErrNotFound)
		}
		return nil, fmt.Errorf("service.GetPreferences: %w", err)
	}
	return prefs, nil
}

func UpdatePreferences(ctx context.Context, userID string, prefs *models.Preferences) error {
	// TODO: wire to repo.UpdateUserPreferences once that function exists (data-model.md pending).
	_ = userID
	_ = prefs
	return fmt.Errorf("service.UpdatePreferences: not implemented")
}

func DeletePreferences(ctx context.Context, userID string) error {
	// TODO: wire to repo.DeleteUserPreferences once that function exists (data-model.md pending).
	_ = userID
	return fmt.Errorf("service.DeletePreferences: not implemented")
}

// isNotFound checks whether an error originates from a DynamoDB item-not-found condition.
// The repo wraps raw errors with fmt.Errorf so we check the message string; replace with
// errors.Is once the repo exports a typed ErrNotFound sentinel.
func isNotFound(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "not found") ||
		strings.Contains(msg, "ResourceNotFoundException") ||
		strings.Contains(msg, "no item")
}
