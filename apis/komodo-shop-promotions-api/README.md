# komodo-shop-promotions-api

Promo code validation and automatic discount management for cart and checkout.

**Status: scaffolded, not deployed.** Promo logic should live in cart-api until checkout volume justifies extraction. Integration points are documented here to guide that migration.

---

## Ports

| Server | Port | Env Var |
|--------|------|---------|
| Public | 7045 | `PORT` |

---

## Deployment Target

**EC2 / ECS Fargate** — called on every cart load and at checkout. Too frequent for Lambda cold starts.

---

## Routes

| Method   | Path                              | Auth     | Description |
|----------|-----------------------------------|----------|-------------|
| `GET`    | `/health`                         | None     | Liveness check |
| `POST`   | `/promotions/validate`            | None     | Validate a promo code; returns discount if eligible |
| `GET`    | `/me/promotions`                  | JWT      | List applicable automatic promotions for the user's cart |
| `GET`    | `/internal/promotions`            | Internal | List all promotions (admin) |
| `POST`   | `/internal/promotions`            | Internal | Create a promotion |
| `PUT`    | `/internal/promotions/{promoId}`  | Internal | Update a promotion |
| `DELETE` | `/internal/promotions/{promoId}`  | Internal | Deactivate a promotion |

---

## Promo Types

| Type | Description |
|------|-------------|
| `percentage_off` | Reduce subtotal by a percentage (e.g. 20% off) |
| `fixed_amount_off` | Reduce subtotal by a fixed amount (e.g. $10 off) |
| `free_shipping` | Zero out shipping cost |
| `buy_x_get_y` | Buy X units of a SKU, get Y free |

---

## Promotion Lifecycle

```
Admin: POST /internal/promotions
  → create promo with code, type, discount, conditions, date window, caps

Customer: POST /promotions/validate { code, cart_subtotal_cents }
  → verify: active, within date window, conditions met, per-user cap not hit
  → return: discount_cents, discount_type, description

Order placed (via order-api):
  → record promo_id + discount_cents on order
  → increment redemption_count
  → publish promotion.redeemed to event-bus-api
```

---

## DynamoDB — Promotions Table

| Attribute | Type | Notes |
|-----------|------|-------|
| `PK` | S | `PROMO#<promoId>` |
| `SK` | S | `METADATA` |
| `promo_id` | S | UUID |
| `code` | S | Case-insensitive promo code (also a GSI key) |
| `type` | S | `percentage_off`, `fixed_amount_off`, `free_shipping`, `buy_x_get_y` |
| `discount_value` | N | Percentage (0–100) or cents |
| `conditions` | M | `min_order_cents`, `eligible_skus`, `eligible_categories` |
| `status` | S | `active`, `inactive`, `expired` |
| `start_at` | S | ISO 8601 |
| `end_at` | S | ISO 8601 (optional) |
| `max_redemptions` | N | Global cap (optional) |
| `redemption_count` | N | Current global redemption total |
| `per_user_limit` | N | Max redemptions per user (default 1) |
| `created_at` | S | ISO 8601 |
| `updated_at` | S | ISO 8601 |

Per-user redemption tracking:

| Attribute | Type | Notes |
|-----------|------|-------|
| `PK` | S | `PROMO#<promoId>` |
| `SK` | S | `USER#<userId>` |
| `redemption_count` | N | Number of times this user has used this promo |

**GSI:** `PromoCodeIndex` — `PK: code` — for O(1) code lookups.

---

## Integration Points

All downstream calls use internal service JWTs (`client_credentials`).

| Service | Call | When |
|---------|------|------|
| `order-api` | records `promo_id` + `discount_cents` on order | On order placement |
| `event-bus-api` | `POST /events` (`promotion.redeemed`) | On successful redemption |

---

## Environment Variables

### Process env

| Variable            | Required | Description |
|---------------------|----------|-------------|
| `APP_NAME`          | Yes | `komodo-shop-promotions-api` |
| `ENV`               | Yes | `local`, `dev`, `staging`, `prod` |
| `LOG_LEVEL`         | Yes | `debug`, `info`, `error` |
| `PORT`              | Yes | e.g. `:7045` |
| `AWS_REGION`        | Yes | e.g. `us-east-2` |
| `AWS_ENDPOINT`      | Yes | LocalStack endpoint or empty for real AWS |
| `AWS_SECRET_PREFIX` | Yes | Secrets Manager path prefix |
| `AWS_SECRET_BATCH`  | Yes | Batch secret name |

### Secrets (resolved from AWS Secrets Manager at startup)

| Key | Description |
|-----|-------------|
| `DYNAMODB_PROMOTIONS_TABLE` | Promotions table name |
| `EVENT_BUS_API_URL` | Internal URL for event publishing |
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
cd apis/komodo-shop-promotions-api
source .env.local
go run .
```

### Run tests
```bash
go test ./...
```
