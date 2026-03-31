package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rdevitto86/komodo-forge-sdk-go/aws/elasticache"
	"github.com/rdevitto86/komodo-forge-sdk-go/config"
	httpErr "github.com/rdevitto86/komodo-forge-sdk-go/http/errors"
	logger "github.com/rdevitto86/komodo-forge-sdk-go/logging/runtime"
	"komodo-cart-api/internal/repo"
	"komodo-cart-api/pkg/v1/client"
	"komodo-cart-api/pkg/v1/models"
)

// errBadRequest and errBadGateway are CartError sentinels for generic HTTP errors
// returned from the service layer. Handlers unwrap them to the correct HTTP status.
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
	holdTTL   int64
	shopItems *client.ShopItemsClient
	inventory *client.InventoryClient
	guestSvc  *GuestCartService
}

// NewCartService constructs a CartService.
func NewCartService(holdTTL int64, shopItems *client.ShopItemsClient, inv *client.InventoryClient, guest *GuestCartService) *CartService {
	return &CartService{
		holdTTL:   holdTTL,
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
// If the item already exists its quantity is incremented additively.
func (s *CartService) AddItem(ctx context.Context, userID string, req models.AddItemRequest) (*models.Cart, error) {
	snapshot, err := s.shopItems.GetItem(ctx, req.ItemID, req.SKU)
	if err != nil {
		return nil, fmt.Errorf("service.Cart.AddItem: fetch snapshot: %w", err)
	}

	exists, err := repo.ItemExists(ctx, userID, req.ItemID)
	if err != nil {
		return nil, fmt.Errorf("service.Cart.AddItem: check exists: %w", err)
	}

	if exists {
		// Fetch current quantity to increment it.
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
	UserID    string            `json:"user_id"`
	CartID    string            `json:"cart_id"`
	HoldIDs   map[string]string `json:"hold_ids"` // sku → holdId
	ExpiresAt time.Time         `json:"expires_at"`
}

// Checkout places stock holds for every item in the cart, issues a checkout token,
// and stores the token payload in Redis with TTL = holdTTL.
func (s *CartService) Checkout(ctx context.Context, userID string) (*models.CheckoutResponse, error) {
	cart, err := repo.GetCart(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("service.Cart.Checkout: get cart: %w", err)
	}
	computeTotals(cart)

	if len(cart.Items) == 0 {
		return nil, errBadRequest
	}

	isLocal := config.GetConfigValue("ENV") == "local"
	holdIDs := make(map[string]string, len(cart.Items))
	var conflictSKUs []string
	expiresAt := time.Now().UTC().Add(time.Duration(s.holdTTL) * time.Second)

	for _, item := range cart.Items {
		holdID, holdExpiry, err := s.inventory.PlaceHold(ctx, item.SKU, cart.ID, item.Quantity)
		if err != nil {
			var holdErr *client.HoldError
			if he, ok := err.(*client.HoldError); ok && he.StatusCode == 409 {
				holdErr = he
				_ = holdErr
				conflictSKUs = append(conflictSKUs, item.SKU)
				continue
			}
			// Network/unexpected error.
			if isLocal {
				// Local dev: synthesise a hold ID so checkout can proceed without inventory-api.
				logger.Warn("service.Cart.Checkout: inventory hold failed in local mode, using synthetic holdID")
				holdIDs[item.SKU] = uuid.NewString()
				continue
			}
			return nil, fmt.Errorf("service.Cart.Checkout: place hold for sku %s: %w", item.SKU, errBadGateway)
		}
		if holdExpiry.Before(expiresAt) {
			expiresAt = holdExpiry
		}
		holdIDs[item.SKU] = holdID
	}

	if len(conflictSKUs) > 0 {
		detail := fmt.Sprintf("insufficient stock for SKUs: %v", conflictSKUs)
		return nil, models.CartError{Code: httpErr.ErrorCode{
			ID:      models.Err.CheckoutFailed.Code.ID,
			Status:  models.Err.CheckoutFailed.Code.Status,
			Message: models.Err.CheckoutFailed.Code.Message,
			Detail:  detail,
		}}
	}

	token := uuid.NewString()
	record := checkoutRecord{
		UserID:    userID,
		CartID:    cart.ID,
		HoldIDs:   holdIDs,
		ExpiresAt: expiresAt,
	}
	data, err := json.Marshal(record)
	if err != nil {
		return nil, fmt.Errorf("service.Cart.Checkout: marshal checkout record: %w", err)
	}
	if err := elasticache.Set("checkout:"+token, string(data), s.holdTTL); err != nil {
		return nil, fmt.Errorf("service.Cart.Checkout: store checkout token: %w", err)
	}

	return &models.CheckoutResponse{
		CheckoutToken: token,
		ExpiresAt:     expiresAt,
		Items:         cart.Items,
	}, nil
}
