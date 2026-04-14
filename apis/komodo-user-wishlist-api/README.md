# komodo-user-wishlist-api

Persistent per-user wishlist: save items, check stock availability, and move items to cart.

**Status: scaffolded, not deployed.** Wishlist logic is independent of cart (no TTL, no guest mode). Implement after cart-api is stable.

---

## Ports

| Server | Port | Env Var |
|--------|------|---------|
| Public | 7053 | `PORT` |

---

## Deployment Target

**EC2 / ECS Fargate** — authenticated-only, low write volume, moderate read volume.

---

## Routes

| Method   | Path                            | Auth | Description |
|----------|---------------------------------|------|-------------|
| `GET`    | `/health`                       | None | Liveness check |
| `GET`    | `/me/wishlist`                  | JWT  | Get the user's wishlist |
| `POST`   | `/me/wishlist/items`            | JWT  | Add an item to the wishlist |
| `DELETE` | `/me/wishlist/items/{itemId}`   | JWT  | Remove an item from the wishlist |
| `GET`    | `/me/wishlist/availability`     | JWT  | Check stock status of all wishlist items |
| `POST`   | `/me/wishlist/move-to-cart`     | JWT  | Move one or more wishlist items to cart |

---

## DynamoDB — Wishlist Table

| Attribute | Type | Notes |
|-----------|------|-------|
| `PK` | S | `WISHLIST#<userId>` |
| `SK` | S | `ITEM#<itemId>` |
| `item_id` | S | shop-items-api item ID |
| `sku` | S | Selected variant SKU |
| `name` | S | Denormalized item name (snapshot at add time) |
| `image_url` | S | Denormalized thumbnail (snapshot at add time) |
| `price_cents` | N | Denormalized price (snapshot at add time) |
| `added_at` | S | ISO 8601 |

No TTL — wishlist items persist indefinitely until explicitly removed.

---

## Integration Points

All downstream calls use internal service JWTs (`client_credentials`).

| Service | Call | When |
|---------|------|------|
| `shop-inventory-api` | `GET /stock/{sku}` | On `GET /me/wishlist/availability` |
| `cart-api` | `POST /me/cart/items` | On `POST /me/wishlist/move-to-cart` per item |

---

## Environment Variables

### Process env

| Variable            | Required | Description |
|---------------------|----------|-------------|
| `APP_NAME`          | Yes | `komodo-user-wishlist-api` |
| `ENV`               | Yes | `local`, `dev`, `staging`, `prod` |
| `LOG_LEVEL`         | Yes | `debug`, `info`, `error` |
| `PORT`              | Yes | e.g. `:7053` |
| `AWS_REGION`        | Yes | e.g. `us-east-2` |
| `AWS_ENDPOINT`      | Yes | LocalStack endpoint or empty for real AWS |
| `AWS_SECRET_PREFIX` | Yes | Secrets Manager path prefix |
| `AWS_SECRET_BATCH`  | Yes | Batch secret name |

### Secrets (resolved from AWS Secrets Manager at startup)

| Key | Description |
|-----|-------------|
| `DYNAMODB_WISHLIST_TABLE` | Wishlist table name |
| `INVENTORY_API_URL` | Internal URL for stock checks |
| `CART_API_URL` | Internal URL for cart item creation |
| `JWT_PUBLIC_KEY` | RSA public key for JWT validation |
| `JWT_ISSUER` | Expected issuer claim |
| `JWT_AUDIENCE` | Expected audience claim |
| `MAX_CONTENT_LENGTH` | Max request body bytes |
| `RATE_LIMIT_RPS` | Rate limit (requests/sec) |
| `RATE_LIMIT_BURST` | Rate limit burst capacity |

---

## Local Development

### Prerequisites
- LocalStack running: `just up` from repo root

### Run
```bash
cd apis/komodo-user-wishlist-api
source .env.local
go run .
```

### Run tests
```bash
go test ./...
```
