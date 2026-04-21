package service

import (
	"context"
	"fmt"
	"os"
	"strconv"
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

// UserServiceAdapter is the interface used to look up a registered account by
// email at order placement. Stubbed until the HTTP adapter is implemented.
type UserServiceAdapter interface {
	// LookupUserByEmail returns the userId for a registered account, or ("", nil)
	// if no account exists for that email.
	LookupUserByEmail(ctx context.Context, email string) (string, error)
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
	user      UserServiceAdapter
}

// NewOrderService constructs an OrderService with the provided adapters.
// Pass nil for any adapter that is not yet implemented — the service will skip
// the corresponding integration and treat the request conservatively (e.g.
// nil user adapter → treat all placements as guest).
func NewOrderService(cart CartServiceAdapter, inventory InventoryServiceAdapter, user UserServiceAdapter) *OrderService {
	return &OrderService{
		cart:      cart,
		inventory: inventory,
		user:      user,
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
		UserID:    "USER#" + userID, // GSI1PK key convention — prefixed for consistency with unified route
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

// PlaceOrderUnified executes the purchase flow for both authenticated and guest
// callers. The email is the universal identity key:
//
//   - If a JWT was validated by middleware, userID is populated and the email
//     from the request body is ignored (the JWT is the authoritative identity).
//   - If no JWT is present, userID is empty and email must be supplied in the
//     request body.
//
// At placement the email is looked up in user-api (if the adapter is wired). A
// match links the order to USER#<userId> and queues an "order added to your
// account" notification. No match results in a GUEST#<uuid> key.
//
// The purchase steps mirror PlaceOrder — idempotency, cart validation (stubbed),
// inventory hold confirmation (stubbed), DynamoDB write, and idempotency cache.
func (s *OrderService) PlaceOrderUnified(ctx context.Context, userID, email, checkoutToken string) (*models.Order, error) {
	// 1. Idempotency check.
	idemKey := idempotencyKey(checkoutToken)
	if existing, err := elasticache.Get(idemKey); err == nil && existing != "" {
		order, fetchErr := repo.GetOrder(ctx, existing)
		if fetchErr != nil {
			logger.Warn("service.PlaceOrderUnified: idempotency hit but GetOrder not implemented; falling through")
		} else if order != nil {
			return order, nil
		}
	}

	// 2. Resolve owner key.
	// When the caller is authenticated (JWT present), trust the userID from the
	// token and use it directly as the GSI key.
	var ownerKey string
	var notifyAccountLink bool
	if userID != "" {
		ownerKey = "USER#" + userID
	} else {
		// Guest path — look up email in user-api to auto-link if account exists.
		if s.user != nil {
			linkedID, err := s.user.LookupUserByEmail(ctx, email)
			if err != nil {
				return nil, fmt.Errorf("service.PlaceOrderUnified: lookup user by email: %w", err)
			}
			if linkedID != "" {
				ownerKey = "USER#" + linkedID
				notifyAccountLink = true
			}
		}
		if ownerKey == "" {
			ownerKey = "GUEST#" + uuid.NewString()
		}
	}

	// 3. Validate checkout token via cart-api adapter (stubbed).
	var snapshot *CheckoutSnapshot
	if s.cart != nil {
		var cartErr error
		snapshot, cartErr = s.cart.ValidateCheckoutToken(ctx, ownerKey, checkoutToken)
		if cartErr != nil {
			return nil, fmt.Errorf("service.PlaceOrderUnified: validate checkout token: %w", cartErr)
		}
	}

	// 4. Build the order.
	now := time.Now().UTC().Format(time.RFC3339)
	orderID := uuid.NewString()

	order := &models.Order{
		ID:        orderID,
		DisplayID: buildDisplayID(orderID),
		UserID:    ownerKey,
		Email:     email,
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

	// 5. Confirm inventory holds (stubbed).
	if s.inventory != nil {
		if err := s.inventory.ConfirmHolds(ctx, orderID, order.Items); err != nil {
			return nil, fmt.Errorf("service.PlaceOrderUnified: confirm holds: %w", err)
		}
	}

	// 6. Persist the order.
	if err := repo.CreateOrder(ctx, order); err != nil {
		return nil, fmt.Errorf("service.PlaceOrderUnified: create order: %w", err)
	}

	// 7. Store idempotency key.
	if err := elasticache.Set(idemKey, orderID, idempotencyTTL); err != nil {
		logger.Warn("service.PlaceOrderUnified: failed to set idempotency key in Redis; order already persisted")
	}

	// 8. Queue account-link notification (non-blocking — best effort).
	// TODO: replace with real call to communications-api once the adapter is wired.
	if notifyAccountLink {
		logger.Info("service.PlaceOrderUnified: guest email matched registered account; notification pending",
			logger.Attr("email", email),
			logger.Attr("order_id", orderID),
		)
	}

	return order, nil
}

// GetOrder retrieves a single order by ID, enforcing user ownership.
// userID is the raw UUID from the JWT (without prefix). The comparison accounts
// for the USER# prefix stored in order.UserID.
func (s *OrderService) GetOrder(ctx context.Context, userID, orderID string) (*models.Order, error) {
	order, err := repo.GetOrder(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("service.GetOrder: fetch: %w", err)
	}
	if order.UserID != "USER#"+userID {
		return nil, fmt.Errorf("service.GetOrder: %w", models.ErrForbidden)
	}
	return order, nil
}

// ListOrders returns all orders for the authenticated user.
// userID is the raw UUID from the JWT; the USER# prefix is applied here before
// the GSI query.
func (s *OrderService) ListOrders(ctx context.Context, userID string) ([]*models.Order, error) {
	orders, err := repo.ListOrdersByUser(ctx, "USER#"+userID)
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
