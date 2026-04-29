package service

import (
	"context"
	"errors"
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

// UserServiceAdapter is the interface used to look up a registered account by
// email at order placement. Stubbed until the HTTP adapter is implemented.
type UserServiceAdapter interface {
	// LookupUserByEmail returns the userId for a registered account, or ("", nil)
	// if no account exists for that email.
	LookupUserByEmail(ctx context.Context, email string) (string, error)
}

// EventBusAdapter is the interface used to publish domain events to event-bus-api.
// Stubbed until the HTTP adapter is implemented.
type EventBusAdapter interface {
	// Publish sends an event of the given type with the provided payload.
	// Callers treat publish failures as non-fatal — log and continue.
	Publish(ctx context.Context, eventType string, payload map[string]any) error
}

// cancellableStatuses is the set of order states from which cancellation is permitted.
var cancellableStatuses = map[models.OrderStatus]bool{
	models.OrderStatusPending:   true,
	models.OrderStatusConfirmed: true,
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
	eventBus  EventBusAdapter
}

// NewOrderService constructs an OrderService with the provided adapters.
// Pass nil for any adapter that is not yet implemented — the service will skip
// the corresponding integration and treat the request conservatively (e.g.
// nil user adapter → treat all placements as guest).
func NewOrderService(cart CartServiceAdapter, inventory InventoryServiceAdapter, user UserServiceAdapter, eventBus EventBusAdapter) *OrderService {
	return &OrderService{
		cart:      cart,
		inventory: inventory,
		user:      user,
		eventBus:  eventBus,
	}
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
// userID is the raw UUID from the JWT. The comparison accounts for the USER#
// prefix stored in order.UserID. Returns ErrNotFound for any miss to prevent
// leaking orderId existence via status-code differences.
func (s *OrderService) GetOrder(ctx context.Context, userID, orderID string) (*models.Order, error) {
	order, err := repo.GetOrder(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("service.GetOrder: %w", models.ErrNotFound)
	}
	if order.UserID != "USER#"+userID {
		return nil, fmt.Errorf("service.GetOrder: ownership mismatch: %w", models.ErrNotFound)
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

// GetOrderInternal retrieves a single order by ID with no user ownership check.
// Used by internal callers (payments-api, returns-api) that operate across users.
// Protected at the transport layer by scope-checked JWT on the private server.
func (s *OrderService) GetOrderInternal(ctx context.Context, orderID string) (*models.Order, error) {
	order, err := repo.GetOrder(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("service.GetOrderInternal: %w", err)
	}
	return order, nil
}

// CancelOrder transitions an order to the cancelled state.
// Only orders in pending or confirmed status can be cancelled; all other
// transitions return ErrNotCancellable (or ErrAlreadyCancelled for orders
// that are already cancelled).
// Ownership is enforced — mismatches return ErrNotFound to prevent callers
// from inferring order existence via status-code differences.
// Event publication to event-bus-api is best-effort: a publish failure is
// logged but does not roll back the cancellation.
func (s *OrderService) CancelOrder(ctx context.Context, userID, orderID string) (*models.Order, error) {
	order, err := repo.GetOrder(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("service.CancelOrder: %w", models.ErrNotFound)
	}

	if order.UserID != "USER#"+userID {
		return nil, fmt.Errorf("service.CancelOrder: ownership mismatch: %w", models.ErrNotFound)
	}

	if order.Status == models.OrderStatusCancelled {
		return nil, fmt.Errorf("service.CancelOrder: %w", models.ErrAlreadyCancelled)
	}
	if !cancellableStatuses[order.Status] {
		return nil, fmt.Errorf("service.CancelOrder: status %s: %w", order.Status, models.ErrNotCancellable)
	}

	if err := repo.UpdateOrderStatus(ctx, orderID, models.OrderStatusCancelled, order.Status); err != nil {
		if errors.Is(err, models.ErrInvalidTransition) {
			// A concurrent request changed the status between our read and write.
			// Surface as NotCancellable — the caller should re-fetch and retry.
			return nil, fmt.Errorf("service.CancelOrder: concurrent transition: %w", models.ErrNotCancellable)
		}
		return nil, fmt.Errorf("service.CancelOrder: update: %w", err)
	}

	// TODO: release inventory holds via shop-inventory-api adapter once wired.
	// TODO: trigger refund via payments-api adapter once wired.

	if s.eventBus != nil {
		if pubErr := s.eventBus.Publish(ctx, "order.cancelled", map[string]any{
			"order_id":     orderID,
			"user_id":      userID,
			"cancelled_at": time.Now().UTC().Format(time.RFC3339),
		}); pubErr != nil {
			logger.Warn("service.CancelOrder: failed to publish order.cancelled event",
				logger.Attr("order_id", orderID),
			)
		}
	}

	order.Status = models.OrderStatusCancelled
	order.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	return order, nil
}

// ListOrders returns a paginated list of orders for the authenticated user.
// limit and cursor are forwarded to the repo layer for DynamoDB-native pagination.
func (s *OrderService) ListOrders(ctx context.Context, userID string, limit int, cursor string) ([]*models.Order, string, error) {
	orders, nextCursor, err := repo.ListOrdersByUser(ctx, userID, limit, cursor)
	if err != nil {
		return nil, "", fmt.Errorf("service.ListOrders: %w", err)
	}
	return orders, nextCursor, nil
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
