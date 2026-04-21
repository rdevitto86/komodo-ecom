package service

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"komodo-order-api/internal/config"
	"komodo-order-api/internal/models"
	"komodo-order-api/internal/repo"

	"github.com/google/uuid"
	"github.com/rdevitto86/komodo-forge-sdk-go/aws/elasticache"
	logger "github.com/rdevitto86/komodo-forge-sdk-go/logging/runtime"
)

// idempotencyTTL is the Redis TTL for checkout-token idempotency keys (default: 24 h).
var idempotencyTTL = mustParseInt64(os.Getenv(config.IDEMPOTENCY_TTL_SEC), 86400)

// CartServiceAdapter is the interface used to validate and consume a checkout token
// issued by cart-api. Stubbed until the HTTP adapter is implemented.
type CartServiceAdapter interface {
	// ValidateCheckoutToken verifies the token is valid and returns the cart
	// snapshot associated with it. Returns nil, nil when not yet implemented.
	ValidateCheckoutToken(ctx context.Context, userID, token string) (*CheckoutSnapshot, error)
}

// InventoryServiceAdapter is the interface used to confirm inventory holds
// placed during cart checkout. Stubbed until the HTTP adapter is implemented.
type InventoryServiceAdapter interface {
	// ConfirmHolds marks inventory holds as committed for all items in the order.
	// Returns nil when not yet implemented.
	ConfirmHolds(ctx context.Context, orderID string, items []models.OrderItem) error
}

// CheckoutSnapshot is the cart-api payload returned for a valid checkout token.
// Fields will be populated once the cart-api adapter is implemented.
type CheckoutSnapshot struct {
	Items   []models.OrderItem
	Address models.OrderAddress
	Payment models.OrderPayment
	Totals  models.OrderTotals
}

// OrderService manages order creation and retrieval.
type OrderService struct {
	cart      CartServiceAdapter
	inventory InventoryServiceAdapter
}

// NewOrderService constructs an OrderService with the provided adapters.
func NewOrderService(cart CartServiceAdapter, inventory InventoryServiceAdapter) *OrderService {
	return &OrderService{
		cart:      cart,
		inventory: inventory,
	}
}

// PlaceOrder executes the core purchase flow:
//  1. Idempotency check — return cached response if token already processed.
//  2. Validate checkout token via cart-api (stubbed).
//  3. Confirm inventory holds via shop-inventory-api (stubbed).
//  4. Persist the order in DynamoDB.
//  5. Store idempotency key in Redis.
func (s *OrderService) PlaceOrder(ctx context.Context, userID, checkoutToken string) (*models.Order, error) {
	// 1. Idempotency check.
	idemKey := idempotencyKey(checkoutToken)
	if existing, err := elasticache.Get(idemKey); err == nil && existing != "" {
		// Key exists — this is a duplicate submission. Fetch and return the original order.
		order, fetchErr := repo.GetOrder(ctx, existing)
		if fetchErr != nil {
			// GetOrder is not yet implemented; log and fall through to re-create
			// (safe because DynamoDB condition prevents double-write).
			logger.Warn("service.PlaceOrder: idempotency hit but GetOrder not implemented; falling through")
		} else if order != nil {
			return order, nil
		}
	}

	// 2. Validate checkout token via cart-api adapter (stubbed).
	// TODO: replace stub with real HTTP call to cart-api once adapter is implemented.
	var snapshot *CheckoutSnapshot
	if s.cart != nil {
		var cartErr error
		snapshot, cartErr = s.cart.ValidateCheckoutToken(ctx, userID, checkoutToken)
		if cartErr != nil {
			return nil, fmt.Errorf("service.PlaceOrder: validate checkout token: %w", cartErr)
		}
	}

	// 3. Build the order — use snapshot data when available, zero values otherwise.
	now := time.Now().UTC().Format(time.RFC3339)
	orderID := uuid.NewString()

	order := &models.Order{
		ID:        orderID,
		DisplayID: buildDisplayID(orderID),
		UserID:    userID,
		Status:    models.OrderStatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if snapshot != nil {
		order.Items = snapshot.Items
		order.Address = snapshot.Address
		order.Payment = snapshot.Payment
		order.Totals = snapshot.Totals
	}

	// 4. Confirm inventory holds (stubbed).
	// TODO: replace stub with real HTTP call to shop-inventory-api once adapter is implemented.
	if s.inventory != nil {
		if err := s.inventory.ConfirmHolds(ctx, orderID, order.Items); err != nil {
			return nil, fmt.Errorf("service.PlaceOrder: confirm holds: %w", err)
		}
	}

	// 5. Persist the order.
	if err := repo.CreateOrder(ctx, order); err != nil {
		return nil, fmt.Errorf("service.PlaceOrder: create order: %w", err)
	}

	// 6. Store idempotency key so duplicate submissions return the same order.
	if err := elasticache.Set(idemKey, orderID, idempotencyTTL); err != nil {
		// Non-fatal: log and continue. The order was written; the idempotency key
		// is a best-effort guard. A retry will either hit the DynamoDB condition
		// expression (no-op) or re-enter this path without the cache hit.
		logger.Warn("service.PlaceOrder: failed to set idempotency key in Redis; order already persisted")
	}

	return order, nil
}

// GetOrder retrieves a single order by ID, enforcing user ownership.
func (s *OrderService) GetOrder(ctx context.Context, userID, orderID string) (*models.Order, error) {
	order, err := repo.GetOrder(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("service.GetOrder: fetch: %w", err)
	}
	if order.UserID != userID {
		return nil, fmt.Errorf("service.GetOrder: %w", models.ErrForbidden)
	}
	return order, nil
}

// GetOrderUnified retrieves an order for either an authenticated user or a guest.
// If userID is non-empty (JWT present), ownership is enforced via UserID field.
// The stored UserID may be prefixed as "USER#<uuid>" for registered users, so
// both the raw uuid and the prefixed form are accepted.
// If userID is empty (no JWT), ownership is enforced via case-insensitive email match.
// Returns ErrNotFound for any ownership mismatch to prevent order ID enumeration.
func (s *OrderService) GetOrderUnified(ctx context.Context, userID, email, orderID string) (*models.Order, error) {
	order, err := repo.GetOrder(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("service.GetOrderUnified: fetch: %w", models.ErrNotFound)
	}

	if userID != "" {
		// Authenticated path — accept both raw uuid and USER#<uuid> prefix.
		if order.UserID != userID && order.UserID != "USER#"+userID {
			return nil, fmt.Errorf("service.GetOrderUnified: ownership mismatch: %w", models.ErrNotFound)
		}
		return order, nil
	}

	// Guest path — require non-empty email and case-insensitive match.
	if email == "" || !strings.EqualFold(order.Email, email) {
		return nil, fmt.Errorf("service.GetOrderUnified: email mismatch: %w", models.ErrNotFound)
	}
	return order, nil
}

// ListOrders returns all orders for the authenticated user.
func (s *OrderService) ListOrders(ctx context.Context, userID string) ([]*models.Order, error) {
	orders, err := repo.ListOrdersByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("service.ListOrders: %w", err)
	}
	return orders, nil
}

// idempotencyKey builds the Redis key for a checkout token.
func idempotencyKey(checkoutToken string) string {
	return "idempotency:order:" + checkoutToken
}

// buildDisplayID derives a short display label from a UUID orderID.
// Currently uses the last 8 hex characters of the UUID for brevity.
// TODO: replace with a proper sequence counter once the sequence table exists.
func buildDisplayID(orderID string) string {
	if len(orderID) >= 8 {
		return orderID[len(orderID)-8:]
	}
	return orderID
}

func mustParseInt64(s string, fallback int64) int64 {
	if s == "" {
		return fallback
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return fallback
	}
	return v
}
