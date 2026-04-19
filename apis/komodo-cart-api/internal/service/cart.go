package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	shopinventory "komodo-cart-api/internal/adapters/shopinventory/v1"
	shopitems "komodo-cart-api/internal/adapters/shopitems/v1"
	"komodo-cart-api/internal/models"
	"komodo-cart-api/internal/repo"

	"github.com/google/uuid"
	"github.com/rdevitto86/komodo-forge-sdk-go/aws/elasticache"
	httpErr "github.com/rdevitto86/komodo-forge-sdk-go/http/errors"
	logger "github.com/rdevitto86/komodo-forge-sdk-go/logging/runtime"
)

var (
	errBadRequest = models.CartError{Code: httpErr.ErrorCode{
		ID:      httpErr.CodeID(httpErr.RangeGlobal, 1),
		Status:  http.StatusBadRequest,
		Message: "Bad request",
	}}
	errBadGateway = models.CartError{Code: httpErr.ErrorCode{
		ID:      httpErr.CodeID(httpErr.RangeGlobal, 12),
		Status:  http.StatusBadGateway,
		Message: "Bad gateway",
	}}
)

// CartService manages authenticated user carts stored in DynamoDB.
type CartService struct {
	tokenTTL  int64
	shopItems *shopitems.Client
	inventory *shopinventory.Client
	guestSvc  *GuestCartService
}

// NewCartService constructs a CartService. tokenTTL is the Redis TTL for checkout tokens (seconds).
func NewCartService(tokenTTL int64, shopItems *shopitems.Client, inv *shopinventory.Client, guest *GuestCartService) *CartService {
	return &CartService{
		tokenTTL:  tokenTTL,
		shopItems: shopItems,
		inventory: inv,
		guestSvc:  guest,
	}
}

// Get fetches the authenticated cart. If guestCartID is non-empty, merges the guest
// cart items into the auth cart first (additive qty for duplicates; auth cart wins on
// price/name conflicts). Deletes the guest cart key after a successful merge.
func (s *CartService) Get(ctx context.Context, userID, guestCartID string) (*models.Cart, error) {
	cart, err := repo.GetCart(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("service.Cart.Get: %w", err)
	}

	if guestCartID != "" {
		guestItems, err := s.guestSvc.MergeIntoAuthCart(ctx, guestCartID)
		if err != nil {
			// Non-fatal: guest cart may be expired. Log and continue with the auth cart.
			logger.Warn("service.Cart.Get: merge guest cart failed, continuing without merge")
		} else {
			for _, guestItem := range guestItems {
				exists := false
				for _, authItem := range cart.Items {
					if authItem.ItemID == guestItem.ItemID && authItem.SKU == guestItem.SKU {
						// Auth cart wins on price/name; quantities are additive.
						if err := repo.UpdateCartItemQuantity(ctx, userID, authItem.ItemID, authItem.Quantity+guestItem.Quantity); err != nil {
							return nil, fmt.Errorf("service.Cart.Get: update merged item qty: %w", err)
						}
						exists = true
						break
					}
				}
				if !exists {
					if err := repo.PutCartItem(ctx, userID, guestItem); err != nil {
						return nil, fmt.Errorf("service.Cart.Get: put merged item: %w", err)
					}
				}
			}

			// Delete guest cart after successful merge.
			if err := repo.DeleteGuestCart(guestCartID); err != nil {
				// Non-fatal: guest cart TTL will clean it up eventually.
				logger.Warn("service.Cart.Get: failed to delete guest cart after merge")
			}

			// Re-fetch to reflect merged state.
			cart, err = repo.GetCart(ctx, userID)
			if err != nil {
				return nil, fmt.Errorf("service.Cart.Get: re-fetch after merge: %w", err)
			}
		}
	}

	computeTotals(cart)
	return cart, nil
}

// AddItem adds an item to the authenticated cart, fetching a product snapshot first.
// Fails with OutOfStock if inventory reports insufficient available_qty.
// If the item already exists its quantity is incremented additively.
func (s *CartService) AddItem(ctx context.Context, userID string, req models.AddItemRequest) (*models.Cart, error) {
	snapshot, err := s.shopItems.GetItem(ctx, req.ItemID, req.SKU)
	if err != nil {
		return nil, fmt.Errorf("service.Cart.AddItem: fetch snapshot: %w", err)
	}

	if os.Getenv("ENV") != "local" {
		if err := s.inventory.CheckStock(ctx, req.SKU, req.Quantity); err != nil {
			if oosErr, ok := err.(*shopinventory.OutOfStockError); ok && oosErr.StatusCode == 409 {
				return nil, models.Err.OutOfStock
			}
			return nil, fmt.Errorf("service.Cart.AddItem: check stock: %w", errBadGateway)
		}
	}

	exists, err := repo.ItemExists(ctx, userID, req.ItemID)
	if err != nil {
		return nil, fmt.Errorf("service.Cart.AddItem: check exists: %w", err)
	}

	if exists {
		cart, err := repo.GetCart(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("service.Cart.AddItem: get cart for qty: %w", err)
		}
		currentQty := req.Quantity
		for _, item := range cart.Items {
			if item.ItemID == req.ItemID {
				currentQty = item.Quantity + req.Quantity
				break
			}
		}
		if err := repo.UpdateCartItemQuantity(ctx, userID, req.ItemID, currentQty); err != nil {
			return nil, fmt.Errorf("service.Cart.AddItem: update qty: %w", err)
		}
	} else {
		snapshot.Quantity = req.Quantity
		if err := repo.PutCartItem(ctx, userID, *snapshot); err != nil {
			return nil, fmt.Errorf("service.Cart.AddItem: put item: %w", err)
		}
	}

	cart, err := repo.GetCart(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("service.Cart.AddItem: get final cart: %w", err)
	}
	computeTotals(cart)
	return cart, nil
}

// UpdateItem sets the quantity of an item. qty==0 removes it.
// Returns ItemNotInCart if the item does not exist.
func (s *CartService) UpdateItem(ctx context.Context, userID, itemID string, req models.UpdateItemRequest) (*models.Cart, error) {
	exists, err := repo.ItemExists(ctx, userID, itemID)
	if err != nil {
		return nil, fmt.Errorf("service.Cart.UpdateItem: check exists: %w", err)
	}
	if !exists {
		return nil, models.Err.ItemNotInCart
	}

	if req.Quantity == 0 {
		if err := repo.DeleteCartItem(ctx, userID, itemID); err != nil {
			return nil, fmt.Errorf("service.Cart.UpdateItem: delete item: %w", err)
		}
	} else {
		if err := repo.UpdateCartItemQuantity(ctx, userID, itemID, req.Quantity); err != nil {
			return nil, fmt.Errorf("service.Cart.UpdateItem: update qty: %w", err)
		}
	}

	cart, err := repo.GetCart(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("service.Cart.UpdateItem: get final cart: %w", err)
	}
	computeTotals(cart)
	return cart, nil
}

// RemoveItem deletes a single item from the cart. Idempotent — no existence check.
func (s *CartService) RemoveItem(ctx context.Context, userID, itemID string) error {
	if err := repo.DeleteCartItem(ctx, userID, itemID); err != nil {
		return fmt.Errorf("service.Cart.RemoveItem: %w", err)
	}
	return nil
}

// Clear removes all items from the authenticated cart.
func (s *CartService) Clear(ctx context.Context, userID string) error {
	if err := repo.ClearCart(ctx, userID); err != nil {
		return fmt.Errorf("service.Cart.Clear: %w", err)
	}
	return nil
}

// checkoutRecord is stored in Redis at checkout:<token>.
// Consumed one-time by order-api (GET + DELETE).
type checkoutRecord struct {
	UserID    string    `json:"user_id"`
	CartID    string    `json:"cart_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

// Checkout performs a final stock check on every cart item, then issues a checkout
// token stored in Redis with TTL = tokenTTL. The stock check guards against OOS items
// that slipped in via bugs or race conditions; it is not a hold.
func (s *CartService) Checkout(ctx context.Context, userID string) (*models.CheckoutResponse, error) {
	cart, err := repo.GetCart(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("service.Cart.Checkout: get cart: %w", err)
	}
	computeTotals(cart)

	if len(cart.Items) == 0 {
		return nil, errBadRequest
	}

	isLocal := os.Getenv("ENV") == "local"
	var conflictSKUs []string

	if !isLocal {
		for _, item := range cart.Items {
			if err := s.inventory.CheckStock(ctx, item.SKU, item.Quantity); err != nil {
				if oosErr, ok := err.(*shopinventory.OutOfStockError); ok && oosErr.StatusCode == 409 {
					conflictSKUs = append(conflictSKUs, item.SKU)
					continue
				}
				return nil, fmt.Errorf("service.Cart.Checkout: check stock for sku %s: %w", item.SKU, errBadGateway)
			}
		}
	}

	if len(conflictSKUs) > 0 {
		return nil, models.CartError{Code: httpErr.ErrorCode{
			ID:      models.Err.CheckoutFailed.Code.ID,
			Status:  models.Err.CheckoutFailed.Code.Status,
			Message: models.Err.CheckoutFailed.Code.Message,
			Detail:  fmt.Sprintf("out of stock: %v", conflictSKUs),
		}}
	}

	expiresAt := time.Now().UTC().Add(time.Duration(s.tokenTTL) * time.Second)
	token := uuid.NewString()
	record := checkoutRecord{
		UserID:    userID,
		CartID:    cart.ID,
		ExpiresAt: expiresAt,
	}
	data, err := json.Marshal(record)
	if err != nil {
		return nil, fmt.Errorf("service.Cart.Checkout: marshal checkout record: %w", err)
	}
	if err := elasticache.Set("checkout:"+token, string(data), s.tokenTTL); err != nil {
		return nil, fmt.Errorf("service.Cart.Checkout: store checkout token: %w", err)
	}

	return &models.CheckoutResponse{
		CheckoutToken: token,
		ExpiresAt:     expiresAt,
		Items:         cart.Items,
	}, nil
}
