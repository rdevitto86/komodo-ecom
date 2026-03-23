# komodo-inventory-api

Stock tracking and reservation service. Manages available inventory per SKU, coordinates stock holds during the cart→checkout transition, and decrements confirmed stock on order confirmation.

Separated from `shop-items-api` because stock levels have fundamentally different write throughput and access patterns than catalog data — every order, reservation, cancellation, and restock is a write here.

---

## Ports

| Server | Port | Env Var |
|--------|------|---------|
| Public | 7044 | `PORT`  |

---

## Deployment Target

**AWS Lambda** — event-driven, bursty write pattern (reservations spike at checkout). Identical binary runs locally as a standard HTTP server.

---

## Routes

| Method   | Path                           | Auth     | Description |
|----------|--------------------------------|----------|-------------|
| `GET`    | `/health`                      | None     | Liveness check |
| `GET`    | `/stock/{sku}`                 | Internal | Get current stock level for a SKU |
| `POST`   | `/stock/{sku}/reserve`         | Internal | Place a stock hold (TTL-based) |
| `DELETE` | `/stock/{sku}/reserve/{holdId}`| Internal | Release a hold early (cancellation/failure) |
| `POST`   | `/stock/{sku}/confirm`         | Internal | Convert hold to confirmed decrement (order confirmed) |
| `POST`   | `/stock/{sku}/restock`         | Internal | Increment stock (purchase received, return processed) |
| `GET`    | `/stock`                       | Internal | Batch stock levels for a list of SKUs |

> **Internal only** — all routes except `/health` require an internal service JWT (`client_credentials` scope). Not exposed via public API gateway.

---

## Stock Hold Pattern

The core concurrency primitive. Prevents overselling under concurrent demand.

```
POST /me/cart/checkout (cart-api)
  → POST /stock/{sku}/reserve for each item
      → DynamoDB conditional write:
          IF available_qty >= requested_qty
          THEN reserved_qty += requested_qty, available_qty -= requested_qty
          ELSE 409 Insufficient Stock
      → hold record written with TTL = HOLD_TTL_SEC (default 900s / 15min)
  ← checkout_token returned to client

order confirmed (order-api)
  → POST /stock/{sku}/confirm?hold_id=<id>
      → DynamoDB: delete hold record, decrement committed_qty

hold expires (DynamoDB TTL)
  → DynamoDB Streams → Lambda → available_qty restored automatically

payment failed / cart abandoned (cart-api or order-api)
  → DELETE /stock/{sku}/reserve/{holdId}
      → DynamoDB: delete hold, restore available_qty
```

---

## DynamoDB — Stock Table

| Attribute | Type | Notes |
|-----------|------|-------|
| `PK` | S | `SKU#<sku>` |
| `SK` | S | `STOCK` (summary) or `HOLD#<holdId>` |
| `available_qty` | N | Units available for new reservations |
| `reserved_qty` | N | Units currently held |
| `committed_qty` | N | Units sold (orders confirmed) |
| `restock_threshold` | N | Alert fires when `available_qty` drops below this |
| `hold_id` | S | UUID — present on `HOLD#` records only |
| `cart_id` | S | Cart that placed this hold |
| `quantity` | N | Units held |
| `ttl` | N | Unix epoch — auto-released by DynamoDB TTL |

**Conditional write example (reserve):**
```
ConditionExpression: available_qty >= :requested
UpdateExpression: SET available_qty = available_qty - :requested,
                      reserved_qty  = reserved_qty  + :requested
```

DynamoDB Streams on this table drive automatic hold release (TTL expiry events) and restock threshold notifications.

---

## Environment Variables

### Process env

| Variable            | Required | Description |
|---------------------|----------|-------------|
| `APP_NAME`          | Yes | `komodo-inventory-api` |
| `ENV`               | Yes | `local`, `dev`, `staging`, `prod` |
| `LOG_LEVEL`         | Yes | `debug`, `info`, `error` |
| `PORT`              | Yes | e.g. `:7044` |
| `AWS_REGION`        | Yes | e.g. `us-east-1` |
| `AWS_ENDPOINT`      | Yes | LocalStack endpoint or empty for real AWS |
| `AWS_SECRET_PREFIX` | Yes | Secrets Manager path prefix |
| `AWS_SECRET_BATCH`  | Yes | Batch secret name |

### Secrets (resolved from AWS Secrets Manager at startup)

| Key | Description |
|-----|-------------|
| `DYNAMODB_STOCK_TABLE` | Stock + holds table name |
| `COMMUNICATIONS_API_URL` | Internal URL for restock threshold alerts |
| `HOLD_TTL_SEC` | Stock hold TTL in seconds (default `900`) |
| `JWT_PUBLIC_KEY` | RSA public key for internal JWT validation |
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
cd apis/komodo-inventory-api
source .env.local
go run ./cmd/server
```

### Run tests
```bash
go test ./...
```
