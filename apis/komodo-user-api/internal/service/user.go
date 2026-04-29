package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"komodo-user-api/internal/models"
	"komodo-user-api/internal/repo"
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
	// Pass addr by pointer so the repo can write back the generated AddressID.
	if err := repo.CreateAddress(ctx, userID, addr); err != nil {
		return fmt.Errorf("service.AddAddress: %w", err)
	}
	return nil
}

func UpdateAddress(ctx context.Context, userID, addressID string, update *models.Address) error {
	update.AddressID = addressID
	if err := repo.UpdateAddress(ctx, userID, *update); err != nil {
		if isNotFound(err) {
			return fmt.Errorf("service.UpdateAddress: %w", ErrNotFound)
		}
		return fmt.Errorf("service.UpdateAddress: %w", err)
	}
	return nil
}

func DeleteAddress(ctx context.Context, userID, addressID string) error {
	if err := repo.DeleteAddress(ctx, userID, addressID); err != nil {
		if isNotFound(err) {
			return fmt.Errorf("service.DeleteAddress: %w", ErrNotFound)
		}
		return fmt.Errorf("service.DeleteAddress: %w", err)
	}
	return nil
}

func GetPayments(ctx context.Context, userID string) ([]models.PaymentMethod, error) {
	methods, err := repo.ListPayments(ctx, userID)
	if err != nil {
		if isNotFound(err) {
			return nil, fmt.Errorf("service.GetPayments: %w", ErrNotFound)
		}
		return nil, fmt.Errorf("service.GetPayments: %w", err)
	}
	return methods, nil
}

func UpsertPayment(ctx context.Context, userID string, pm *models.PaymentMethod) error {
	// ID generation is delegated to the repo so the assigned PaymentID is
	// reflected back on the caller's pointer after the write succeeds.
	if err := repo.UpsertPayment(ctx, userID, pm); err != nil {
		return fmt.Errorf("service.UpsertPayment: %w", err)
	}
	return nil
}

func DeletePayment(ctx context.Context, userID, paymentID string) error {
	if err := repo.DeletePayment(ctx, userID, paymentID); err != nil {
		if isNotFound(err) {
			return fmt.Errorf("service.DeletePayment: %w", ErrNotFound)
		}
		return fmt.Errorf("service.DeletePayment: %w", err)
	}
	return nil
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
	if err := repo.UpdateUserPreferences(ctx, userID, prefs); err != nil {
		if isNotFound(err) {
			return fmt.Errorf("service.UpdatePreferences: %w", ErrNotFound)
		}
		return fmt.Errorf("service.UpdatePreferences: %w", err)
	}
	return nil
}

func DeletePreferences(ctx context.Context, userID string) error {
	if err := repo.DeleteUserPreferences(ctx, userID); err != nil {
		if isNotFound(err) {
			return fmt.Errorf("service.DeletePreferences: %w", ErrNotFound)
		}
		return fmt.Errorf("service.DeletePreferences: %w", err)
	}
	return nil
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
