package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	httpErr "github.com/rdevitto86/komodo-forge-sdk-go/http/errors"
	shopitems "komodo-cart-api/internal/adapters/shopitems/v1"
	"komodo-cart-api/internal/models"
	"komodo-cart-api/internal/repo"
)

// errForbidden is a package-level sentinel for session mismatch that satisfies error.
var errForbidden = models.CartError{Code: httpErr.ErrorCode{
	ID:      httpErr.CodeID(httpErr.RangeGlobal, 4),
	Status:  http.StatusForbidden,
	Message: "Forbidden",
}}

// GuestCartService manages unauthenticated guest carts stored in Redis.
type GuestCartService struct {
	guestTTL  int64
	shopItems *shopitems.Client
}

// NewGuestCartService constructs a GuestCartService.
func NewGuestCartService(ttl int64, shopItems *shopitems.Client) *GuestCartService {
	return &GuestCartService{guestTTL: ttl, shopItems: shopItems}
}

// Create generates a new empty guest cart and persists it to Redis.
// Returns (cart, sessionID, error).
func (s *GuestCartService) Create(ctx context.Context) (*models.Cart, string, error) {
	cartID := uuid.NewString()
	sessionID := uuid.NewString()

	cart := models.Cart{
		ID:        cartID,
		Items:     []models.CartItem{},
		UpdatedAt: time.Now().UTC(),
	}

	if err := repo.CreateGuestCart(cart, sessionID, s.guestTTL); err != nil {
		return nil, "", fmt.Errorf("service.GuestCart.Create: %w", err)
	}
	return &cart, sessionID, nil
}

// Get retrieves a guest cart after validating the session token.
func (s *GuestCartService) Get(ctx context.Context, cartID, sessionID string) (*models.Cart, error) {
	record, err := repo.GetGuestCart(cartID)
	if err != nil {
		return nil, fmt.Errorf("service.GuestCart.Get: %w", err)
	}
	if err := s.validateSession(record, sessionID); err != nil {
		return nil, err
	}
	return &record.Cart, nil
}

// AddItem adds an item to a guest cart, fetching a fresh product snapshot first.
// If the item already exists its quantity is incremented additively.
func (s *GuestCartService) AddItem(ctx context.Context, cartID, sessionID string, req models.AddItemRequest) (*models.Cart, error) {
	record, err := repo.GetGuestCart(cartID)
	if err != nil {
		return nil, fmt.Errorf("service.GuestCart.AddItem: %w", err)
	}
	if err := s.validateSession(record, sessionID); err != nil {
		return nil, err
	}

	snapshot, err := s.shopItems.GetItem(ctx, req.ItemID, req.SKU)
	if err != nil {
		return nil, fmt.Errorf("service.GuestCart.AddItem: fetch item snapshot: %w", err)
	}

	found := false
	for i, existing := range record.Cart.Items {
		if existing.ItemID == req.ItemID && existing.SKU == req.SKU {
			record.Cart.Items[i].Quantity += req.Quantity
			found = true
			break
		}
	}
	if !found {
		snapshot.Quantity = req.Quantity
		record.Cart.Items = append(record.Cart.Items, *snapshot)
	}

	record.Cart.UpdatedAt = time.Now().UTC()
	computeTotals(&record.Cart)

	if err := repo.SaveGuestCart(record, s.guestTTL); err != nil {
		return nil, fmt.Errorf("service.GuestCart.AddItem: save: %w", err)
	}
	return &record.Cart, nil
}

// UpdateItem sets the quantity of an existing item. qty==0 removes the item.
func (s *GuestCartService) UpdateItem(ctx context.Context, cartID, sessionID, itemID string, req models.UpdateItemRequest) (*models.Cart, error) {
	record, err := repo.GetGuestCart(cartID)
	if err != nil {
		return nil, fmt.Errorf("service.GuestCart.UpdateItem: %w", err)
	}
	if err := s.validateSession(record, sessionID); err != nil {
		return nil, err
	}

	found := false
	if req.Quantity == 0 {
		// Remove item when qty is 0.
		updated := record.Cart.Items[:0]
		for _, item := range record.Cart.Items {
			if item.ItemID != itemID {
				updated = append(updated, item)
			} else {
				found = true
			}
		}
		record.Cart.Items = updated
	} else {
		for i, item := range record.Cart.Items {
			if item.ItemID == itemID {
				record.Cart.Items[i].Quantity = req.Quantity
				found = true
				break
			}
		}
	}

	if !found {
		return nil, models.Err.ItemNotInCart // CartError implements error
	}

	record.Cart.UpdatedAt = time.Now().UTC()
	computeTotals(&record.Cart)

	if err := repo.SaveGuestCart(record, s.guestTTL); err != nil {
		return nil, fmt.Errorf("service.GuestCart.UpdateItem: save: %w", err)
	}
	return &record.Cart, nil
}

// RemoveItem removes a single item from a guest cart.
func (s *GuestCartService) RemoveItem(ctx context.Context, cartID, sessionID, itemID string) error {
	record, err := repo.GetGuestCart(cartID)
	if err != nil {
		return fmt.Errorf("service.GuestCart.RemoveItem: %w", err)
	}
	if err := s.validateSession(record, sessionID); err != nil {
		return err
	}

	updated := record.Cart.Items[:0]
	for _, item := range record.Cart.Items {
		if item.ItemID != itemID {
			updated = append(updated, item)
		}
	}
	record.Cart.Items = updated
	record.Cart.UpdatedAt = time.Now().UTC()
	computeTotals(&record.Cart)

	if err := repo.SaveGuestCart(record, s.guestTTL); err != nil {
		return fmt.Errorf("service.GuestCart.RemoveItem: save: %w", err)
	}
	return nil
}

// Clear resets a guest cart's items to empty while keeping the key alive.
func (s *GuestCartService) Clear(ctx context.Context, cartID, sessionID string) error {
	record, err := repo.GetGuestCart(cartID)
	if err != nil {
		return fmt.Errorf("service.GuestCart.Clear: %w", err)
	}
	if err := s.validateSession(record, sessionID); err != nil {
		return err
	}

	record.Cart.Items = []models.CartItem{}
	record.Cart.UpdatedAt = time.Now().UTC()
	computeTotals(&record.Cart)

	if err := repo.SaveGuestCart(record, s.guestTTL); err != nil {
		return fmt.Errorf("service.GuestCart.Clear: save: %w", err)
	}
	return nil
}

// MergeIntoAuthCart fetches a guest cart without session validation (server-initiated merge).
// Returns the items for the auth cart service to merge. Does NOT delete the Redis key —
// the auth service deletes it after a successful merge.
func (s *GuestCartService) MergeIntoAuthCart(ctx context.Context, cartID string) ([]models.CartItem, error) {
	record, err := repo.GetGuestCart(cartID)
	if err != nil {
		return nil, fmt.Errorf("service.GuestCart.MergeIntoAuthCart: %w", err)
	}
	if record == nil {
		return nil, models.Err.NotFound // CartError implements error
	}
	return record.Cart.Items, nil
}

// validateSession checks that the record exists and the session token matches.
func (s *GuestCartService) validateSession(record *repo.GuestCartRecord, sessionID string) error {
	if record == nil {
		return models.Err.NotFound
	}
	if record.SessionID != sessionID {
		return errForbidden
	}
	return nil
}

// computeTotals sets SubtotalCents and ItemCount by summing the cart's Items.
func computeTotals(cart *models.Cart) {
	subtotal := 0
	count := 0
	for _, item := range cart.Items {
		subtotal += item.UnitPriceCents * item.Quantity
		count += item.Quantity
	}
	cart.SubtotalCents = subtotal
	cart.ItemCount = count
}
