# komodo-cart-api

Shopping cart service. Manages guest and authenticated carts, persists cart state, and coordinates stock holds during checkout initiation.

---

## Ports

| Server | Port | Env Var |
|--------|------|---------|
| Public | 7043 | `PORT`  |

---

## Routes

### Authenticated cart (`/me/cart` — requires JWT)

| Method   | Path                        | Auth | Description                                      |
|----------|-----------------------------|------|--------------------------------------------------|
| `GET`    | `/health`                   | None | Liveness check                                   |
| `GET`    | `/me/cart`                  | JWT  | Get current user's cart                          |
| `POST`   | `/me/cart/merge`            | JWT  | Merge a guest cart into the authenticated cart   |
| `POST`   | `/me/cart/items`            | JWT  | Add item to cart                                 |
| `PUT`    | `/me/cart/items/{itemId}`   | JWT  | Update item quantity                             |
| `DELETE` | `/me/cart/items/{itemId}`   | JWT  | Remove item from cart                            |
| `DELETE` | `/me/cart`                  | JWT  | Clear cart                                       |
| `POST`   | `/me/cart/checkout`         | JWT  | Initiate checkout — places stock holds, returns checkout token |

### Guest cart (`/cart` — session token via `X-Session-ID` header)

| Method   | Path                              | Auth    | Description                       |
|----------|-----------------------------------|---------|-----------------------------------|
| `POST`   | `/cart`                           | None    | Create guest cart, returns cart ID |
| `GET`    | `/cart/{cartId}`                  | Session | Get guest cart                    |
| `POST`   | `/cart/{cartId}/items`            | Session | Add item to guest cart            |
| `PUT`    | `/cart/{cartId}/items/{itemId}`   | Session | Update item quantity              |
| `DELETE` | `/cart/{cartId}/items/{itemId}`   | Session | Remove item from guest cart       |
| `DELETE` | `/cart/{cartId}`                  | Session | Clear guest cart                  |

---

## Cart Lifecycle

```
Guest cart (Redis, TTL 7d)
    → user logs in
    → UI calls POST /me/cart/merge with { guest_cart_id } in body
        → items merged into authenticated cart (DynamoDB, no TTL)
        → quantities additive for duplicate items; auth cart wins on price/name conflict
        → Redis guest cart key deleted
        → merged cart returned in response (saves a round-trip)
    → POST /me/cart/checkout
        → stock holds placed in inventory-api (TTL 15min)
        → checkout token (Redis UUID) returned to order-api
    → order confirmed → holds converted to decrements
    → order failed / timeout → holds released
```

> **TODO: Save for Later** — Amazon-style "save for later" feature: move items out of the
> active cart into a persisted saved-items list (separate DynamoDB entity, no TTL).
> Not in scope for initial build. Design separately before implementing.

---

## Environment Variables

### Process env

| Variable            | Required | Description |
|---------------------|----------|-------------|
| `APP_NAME`          | Yes | `komodo-cart-api` |
| `ENV`               | Yes | `local`, `dev`, `staging`, `prod` |
| `LOG_LEVEL`         | Yes | `debug`, `info`, `error` |
| `PORT`              | Yes | Public port (e.g. `:7043`) |
| `AWS_REGION`        | Yes | e.g. `us-east-1` |
| `AWS_ENDPOINT`      | Yes | LocalStack endpoint or empty for real AWS |
| `AWS_SECRET_PREFIX` | Yes | Secrets Manager path prefix |
| `AWS_SECRET_BATCH`  | Yes | Batch secret name |
| `EVAL_RULES_PATH`   | Yes | Path to validation rules YAML (e.g. `/app/config/validation_rules.yaml`) |

### Secrets (resolved from AWS Secrets Manager at startup)

| Key                          | Description |
|------------------------------|-------------|
| `AWS_ELASTICACHE_ENDPOINT`   | Redis endpoint (e.g. `localhost:6379`) |
| `AWS_ELASTICACHE_PASSWORD`   | Redis password (empty for local) |
| `AWS_ELASTICACHE_DB`         | Redis DB index |
| `DYNAMODB_CARTS_TABLE`       | DynamoDB table for authenticated carts |
| `INVENTORY_API_URL`          | Internal URL for stock hold requests |
| `SHOP_ITEMS_API_URL`         | Internal URL for product name/price lookups at add-item time |
| `CART_GUEST_TTL_SEC`         | Guest cart Redis TTL in seconds (default `604800` / 7d) |
| `CART_HOLD_TTL_SEC`          | Stock hold TTL in seconds (default `900` / 15min) |
| `JWT_PUBLIC_KEY`             | RSA public key for validating incoming JWTs |
| `JWT_PRIVATE_KEY`            | RSA private key for signing checkout tokens |
| `JWT_ISSUER`                 | Expected JWT issuer claim |
| `JWT_AUDIENCE`               | Expected JWT audience claim |
| `JWT_KID`                    | Key ID header for issued tokens |
| `MAX_CONTENT_LENGTH`         | Max request body bytes |
| `RATE_LIMIT_RPS`             | Rate limit (requests/sec) |
| `RATE_LIMIT_BURST`           | Rate limit burst capacity |
| `IDEMPOTENCY_TTL_SEC`        | Idempotency key TTL for write endpoints |

---

## DynamoDB — Authenticated Carts Table

| Attribute        | Type | Row          | Notes                                          |
|------------------|------|--------------|------------------------------------------------|
| `PK`             | S    | Both         | `CART#<userId>`                                |
| `SK`             | S    | Both         | `METADATA` or `ITEM#<itemId>`                  |
| `user_id`        | S    | METADATA     | Denormalised for convenience                   |
| `updated_at`     | S    | Both         | ISO 8601, updated on every write               |
| `item_id`        | S    | ITEM#*       | Product ID from shop-items-api                 |
| `sku`            | S    | ITEM#*       | Variant SKU                                    |
| `name`           | S    | ITEM#*       | Product name snapshot at add-time              |
| `quantity`       | N    | ITEM#*       | Requested quantity                             |
| `unit_price_cents` | N  | ITEM#*       | Price snapshot at add-time (integer cents)     |
| `image_url`      | S    | ITEM#*       | Optional, snapshot at add-time                 |

**No TTL on authenticated cart items** — items persist indefinitely until the user removes
them or completes checkout. There is no abandoned cart expiry; a background cleanup job can
be added later if storage cost becomes a concern.

No GSIs on initial build. Access pattern is always `CART#<userId>`.

### Cart ID

`Cart.ID` in API responses is a **deterministic UUID derived from the user ID**:

```go
cartID := uuid.NewSHA1(uuid.NameSpaceOID, []byte(userID)).String()
```

This is a Version 5 UUID (SHA-1 hash of the namespace + userID). It is always the same
for a given user and is computed on the fly when building the response — never stored in
DynamoDB. The client gets a stable, opaque UUID it can reference without any extra
create/fetch round-trip.

---

## Redis

All cart-related keys share a single Redis DB. Key namespaces:

| Key pattern               | TTL                    | Description                                      |
|---------------------------|------------------------|--------------------------------------------------|
| `cart:guest:<cartId>`     | `CART_GUEST_TTL_SEC`   | Guest cart — JSON-serialised cart + session ID   |
| `checkout:<token>`        | `CART_HOLD_TTL_SEC`    | One-time checkout token — consumed by order-api  |

### Guest carts (`cart:guest:<cartId>`)

Value: JSON envelope `{ session_id, cart }` — session ID is stored alongside the cart so
validation is atomic with the read.

TTL is refreshed on every write. On merge (`POST /me/cart/merge`), the UI sends the
`guest_cart_id` (a UUID it generated and stored client-side) in the request body. The
backend merges items into the authenticated cart and deletes the Redis key.

### Checkout tokens (`checkout:<token>`)

Value: JSON `{ user_id, cart_id, hold_ids: {sku: holdId}, expires_at }`.
Generated by `POST /me/cart/checkout`. Order-api validates via a single Redis GET then
DELETEs the key to consume it (one-time use). TTL matches the stock hold TTL.

---

## Local Development

### Prerequisites
- LocalStack + Redis running: `just up` from repo root

### Run
```bash
cd apis/komodo-cart-api
source .env.local
go run .
```

### Run tests
```bash
go test ./...
```

### cURL Examples

```bash
# ── Setup — get a token from auth-api ────────────────────────────────────────
TOKEN=$(curl -s -X POST http://localhost:7011/oauth/token \
  -H "Content-Type: application/json" \
  -d '{"clientId":"test-client","clientSecret":"test-secret","grantType":"client_credentials"}' \
  | jq -r .accessToken)

# ── Health ────────────────────────────────────────────────────────────────────
curl http://localhost:7043/health

# ── Guest cart ────────────────────────────────────────────────────────────────

# Create guest cart — returns cartId and sessionId
GUEST=$(curl -s -X POST http://localhost:7043/cart)
CART_ID=$(echo $GUEST | jq -r .id)
SESSION_ID=$(echo $GUEST | jq -r .session_id)   # also returned in X-Session-ID response header

# Get guest cart
curl -s http://localhost:7043/cart/$CART_ID \
  -H "X-Session-ID: $SESSION_ID" | jq

# Add item to guest cart
curl -s -X POST http://localhost:7043/cart/$CART_ID/items \
  -H "X-Session-ID: $SESSION_ID" \
  -H "Content-Type: application/json" \
  -d '{"item_id":"item_abc123","sku":"SKU-001","quantity":2}' | jq

# Update item quantity in guest cart
curl -s -X PUT http://localhost:7043/cart/$CART_ID/items/item_abc123 \
  -H "X-Session-ID: $SESSION_ID" \
  -H "Content-Type: application/json" \
  -d '{"quantity":3}' | jq

# Remove item from guest cart
curl -s -X DELETE http://localhost:7043/cart/$CART_ID/items/item_abc123 \
  -H "X-Session-ID: $SESSION_ID"

# Clear guest cart
curl -s -X DELETE http://localhost:7043/cart/$CART_ID \
  -H "X-Session-ID: $SESSION_ID"

# ── Authenticated cart ────────────────────────────────────────────────────────

# Get authenticated cart
curl -s http://localhost:7043/me/cart \
  -H "Authorization: Bearer $TOKEN" | jq

# Merge guest cart into authenticated cart (run after login)
curl -s -X POST http://localhost:7043/me/cart/merge \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"guest_cart_id\":\"$CART_ID\"}" | jq

# Add item to authenticated cart
curl -s -X POST http://localhost:7043/me/cart/items \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"item_id":"item_abc123","sku":"SKU-001","quantity":1}' | jq

# Update item quantity (set to 0 to remove)
curl -s -X PUT http://localhost:7043/me/cart/items/item_abc123 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"quantity":2}' | jq

# Remove item from authenticated cart
curl -s -X DELETE http://localhost:7043/me/cart/items/item_abc123 \
  -H "Authorization: Bearer $TOKEN"

# Clear authenticated cart
curl -s -X DELETE http://localhost:7043/me/cart \
  -H "Authorization: Bearer $TOKEN"

# Initiate checkout — places stock holds, returns checkout_token
curl -s -X POST http://localhost:7043/me/cart/checkout \
  -H "Authorization: Bearer $TOKEN" | jq
```
