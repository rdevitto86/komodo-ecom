# komodo-order-returns-api

Returns and RMA (Return Merchandise Authorization) service. Manages the full return lifecycle from customer request through refund and restock.

**Status: scaffolded, not deployed.** Returns logic lives in `order-api` until volume and complexity justify migration. Deploy when return seams become painful — the integration points are documented here to guide that extraction.

---

## Ports

| Server | Port | Env Var |
|--------|------|---------|
| Public | 7062 | `PORT`  |

---

## Deployment Target

**AWS Lambda** — infrequent, event-driven. DynamoDB for RMA state.

---

## Routes

| Method   | Path                          | Auth | Description |
|----------|-------------------------------|------|-------------|
| `GET`    | `/health`                     | None | Liveness check |
| `POST`   | `/me/returns`                 | JWT  | Initiate a return request |
| `GET`    | `/me/returns`                 | JWT  | List user's return requests |
| `GET`    | `/me/returns/{returnId}`      | JWT  | Get return status |
| `DELETE` | `/me/returns/{returnId}`      | JWT  | Cancel a pending return request |
| `GET`    | `/returns/{returnId}`         | Internal | Get return (admin/service use) |
| `PUT`    | `/returns/{returnId}/approve` | Internal | Approve return, trigger refund |
| `PUT`    | `/returns/{returnId}/receive` | Internal | Mark items received, trigger restock |
| `PUT`    | `/returns/{returnId}/reject`  | Internal | Reject return with reason |

---

## RMA Lifecycle

```
Customer: POST /me/returns
  → validate: order exists, belongs to user, within return window
  → create RMA record (status: requested)

Admin/Service: PUT /returns/{id}/approve
  → POST payments-api /refunds (partial or full)
  → status: approved
  → communications-api: "Your return has been approved"

Customer ships items back
Admin/Service: PUT /returns/{id}/receive
  → POST inventory-api /stock/{sku}/restock (reason: return_processed)
  → POST loyalty-api /me/points/reverse (if points were earned on original order)
  → status: processed
  → communications-api: "Your refund is on its way"

PUT /returns/{id}/reject
  → status: rejected
  → communications-api: "Your return request was declined"
```

---

## DynamoDB — Returns Table

| Attribute | Type | Notes |
|-----------|------|-------|
| `PK` | S | `RETURN#<userId>` |
| `SK` | S | `RMA#<returnId>` |
| `return_id` | S | UUID |
| `order_id` | S | Reference to order-api |
| `status` | S | `requested`, `approved`, `received`, `processed`, `rejected`, `cancelled` |
| `items` | L | Items being returned (item_id, sku, quantity, reason) |
| `refund_amount_cents` | N | Approved refund amount |
| `refund_id` | S | Reference to payments-api refund (set on approval) |
| `return_window_expires` | S | ISO 8601 — deadline for customer to ship |
| `created_at` | S | ISO 8601 |
| `updated_at` | S | ISO 8601 |

**GSI:** `OrderReturnsIndex` — `PK: order_id`, `SK: return_id` — for order-api to look up returns by order.

---

## Integration Points

All downstream calls use internal service JWTs (`client_credentials`).

| Service | Call | When |
|---------|------|------|
| `order-api` | `GET /orders/{orderId}` | On initiation — validate order ownership and return window |
| `payments-api` | `POST /refunds` | On approval |
| `inventory-api` | `POST /stock/{sku}/restock` | On receipt, per item |
| `loyalty-api` | `POST /me/points/reverse` | On receipt — reverse points earned on original order |
| `communications-api` | `POST /messages` | On each status transition |

---

## Environment Variables

### Process env

| Variable            | Required | Description |
|---------------------|----------|-------------|
| `APP_NAME`          | Yes | `komodo-returns-api` |
| `ENV`               | Yes | `local`, `dev`, `staging`, `prod` |
| `LOG_LEVEL`         | Yes | `debug`, `info`, `error` |
| `PORT`              | Yes | e.g. `:7062` |
| `AWS_REGION`        | Yes | e.g. `us-east-1` |
| `AWS_ENDPOINT`      | Yes | LocalStack endpoint or empty for real AWS |
| `AWS_SECRET_PREFIX` | Yes | Secrets Manager path prefix |
| `AWS_SECRET_BATCH`  | Yes | Batch secret name |

### Secrets (resolved from AWS Secrets Manager at startup)

| Key | Description |
|-----|-------------|
| `DYNAMODB_RETURNS_TABLE` | Returns table name |
| `ORDER_API_URL` | Internal URL for order validation |
| `PAYMENTS_API_URL` | Internal URL for refund initiation |
| `INVENTORY_API_URL` | Internal URL for restock |
| `LOYALTY_API_URL` | Internal URL for points reversal |
| `COMMUNICATIONS_API_URL` | Internal URL for customer notifications |
| `RETURN_WINDOW_DAYS` | Default return window in days (e.g. `30`) |
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
cd apis/komodo-returns-api
source .env.local
go run ./cmd/server
```

### Run tests
```bash
go test ./...
```
