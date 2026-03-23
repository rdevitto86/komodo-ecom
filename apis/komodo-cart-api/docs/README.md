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
| `GET`    | `/me/cart`                  | JWT  | Get current user's cart (merges guest on login)  |
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
    → user logs in → merge into authenticated cart (DynamoDB)
    → POST /me/cart/checkout
        → stock holds placed in inventory-api (TTL 15min)
        → checkout token returned to order-api
    → order confirmed → holds converted to decrements
    → order failed / timeout → holds released
```

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

### Secrets (resolved from AWS Secrets Manager at startup)

| Key                          | Description |
|------------------------------|-------------|
| `AWS_ELASTICACHE_ENDPOINT`   | Redis endpoint (e.g. `localhost:6379`) |
| `AWS_ELASTICACHE_PASSWORD`   | Redis password (empty for local) |
| `AWS_ELASTICACHE_DB`         | Redis DB index |
| `DYNAMODB_CARTS_TABLE`       | DynamoDB table for authenticated carts |
| `INVENTORY_API_URL`          | Internal URL for stock hold requests |
| `CART_GUEST_TTL_SEC`         | Guest cart Redis TTL in seconds (default `604800` / 7d) |
| `CART_HOLD_TTL_SEC`          | Stock hold TTL in seconds (default `900` / 15min) |
| `JWT_PUBLIC_KEY`             | RSA public key for JWT validation |
| `JWT_ISSUER`                 | Expected JWT issuer claim |
| `JWT_AUDIENCE`               | Expected JWT audience claim |
| `MAX_CONTENT_LENGTH`         | Max request body bytes |
| `RATE_LIMIT_RPS`             | Rate limit (requests/sec) |
| `RATE_LIMIT_BURST`           | Rate limit burst capacity |

---

## DynamoDB — Authenticated Carts Table

| Attribute | Type | Notes |
|-----------|------|-------|
| `PK` | S | `CART#<userId>` |
| `SK` | S | `METADATA` (cart summary) or `ITEM#<itemId>` |
| `item_id` | S | Product ID from shop-items-api |
| `sku` | S | Variant SKU |
| `quantity` | N | Requested quantity |
| `unit_price` | N | Price snapshot at add-time (pence/cents) |
| `updated_at` | S | ISO 8601 |
| `ttl` | N | Unix epoch — abandoned cart expiry (90d) |

No GSIs on initial build. Access pattern is always `CART#<userId>`.

---

## Redis — Guest Carts

Key: `cart:<cartId>` (UUID)
Value: JSON-serialised cart (items array + metadata)
TTL: `CART_GUEST_TTL_SEC` (default 7 days), refreshed on every write.

On authenticated login, the UI sends the `cartId` cookie and the backend merges items (quantity additive for duplicates, authenticated cart wins on conflict) then deletes the Redis key.

---

## Local Development

### Prerequisites
- LocalStack + Redis running: `just up` from repo root

### Run
```bash
cd apis/komodo-cart-api
source .env.local
go run ./cmd/server
```

### Run tests
```bash
go test ./...
```
