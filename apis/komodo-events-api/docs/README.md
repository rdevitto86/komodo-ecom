# komodo-events-api

Internal event relay service. Accepts domain events from producer services and fans them out to consumers via AWS SNS + SQS.

**Status: built, not deployed.** Deploy when `order.placed` has 3+ consumers and direct HTTP fan-out becomes a coupling liability.

---

## Ports

| Server | Port | Env Var |
|--------|------|---------|
| Public | 7002 | `PORT`  |

---

## Deployment Target

**Deferred.** AWS target when deployed: SNS (one topic per event type) + SQS (per-consumer queue with DLQ). This matches the existing EC2/Fargate stack with no new infrastructure primitives. For local dev, runs as a standard HTTP server with in-process fan-out.

---

## Routes

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/health` | None | Liveness check |
| `POST` | `/events` | Internal JWT | Publish an event |
| `GET` | `/events/types` | Internal JWT | List registered event types |

All routes except `/health` require an internal service JWT (`client_credentials` scope).

---

## Event Envelope

Every event published to `POST /events` must conform to this envelope. **Schema versioning is enforced from day one** — adding fields is backwards-compatible; renaming or removing fields is a breaking change that requires a version bump.

```json
{
  "id": "uuid-v4",
  "type": "order.placed",
  "version": "1",
  "source": "komodo-order-api",
  "occurred_at": "2026-03-10T14:23:01Z",
  "payload": { }
}
```

| Field | Type | Notes |
|-------|------|-------|
| `id` | UUID | Idempotency key — consumers must deduplicate |
| `type` | string | `<domain>.<action>` — see event catalogue below |
| `version` | string | Numeric string. Increment when payload shape changes |
| `source` | string | Publishing service name |
| `occurred_at` | ISO 8601 | When the event happened (not when it was published) |
| `payload` | object | Event-specific data — typed per event type |

---

## Event Catalogue

| Event Type | Producer | Consumers | Payload (key fields) |
|------------|----------|-----------|----------------------|
| `order.placed` | order-api | loyalty-api, communications-api, inventory-api | `order_id`, `user_id`, `items[]`, `total_cents` |
| `order.cancelled` | order-api | inventory-api, communications-api | `order_id`, `user_id`, `reason` |
| `order.fulfilled` | order-api | communications-api | `order_id`, `user_id`, `tracking_number` |
| `payment.confirmed` | payments-api | order-api, communications-api | `payment_id`, `order_id`, `amount_cents` |
| `payment.failed` | payments-api | cart-api, communications-api | `payment_id`, `order_id`, `reason` |
| `stock.low` | inventory-api | communications-api | `sku`, `available_qty`, `threshold` |
| `user.registered` | user-api | loyalty-api, communications-api | `user_id`, `email` |
| `return.approved` | returns-api | inventory-api, payments-api, loyalty-api, communications-api | `return_id`, `order_id`, `items[]` |

---

## AWS Architecture (when deployed)

```
Producer service
  → POST /events (events-api)
      → SNS topic: komodo-events-<type>-<env>
          → SQS queue per consumer (komodo-events-<consumer>-<type>-<env>)
              → DLQ for each queue (retry 3x, then dead-letter)
          → Consumer Lambda / ECS service polls its own queue
```

**Topic naming:** `komodo-events-order-placed-prod`
**Queue naming:** `komodo-events-loyalty-order-placed-prod`
**DLQ naming:** `komodo-events-loyalty-order-placed-prod-dlq`

CloudFormation templates for SNS/SQS resources belong in `infra/deploy/cfn/`.

---

## Environment Variables

### Process env

| Variable            | Required | Description |
|---------------------|----------|-------------|
| `APP_NAME`          | Yes | `komodo-events-api` |
| `ENV`               | Yes | `local`, `dev`, `staging`, `prod` |
| `LOG_LEVEL`         | Yes | `debug`, `info`, `error` |
| `PORT`              | Yes | e.g. `:7002` |
| `AWS_REGION`        | Yes | e.g. `us-east-1` |
| `AWS_ENDPOINT`      | Yes | LocalStack endpoint or empty for real AWS |
| `AWS_SECRET_PREFIX` | Yes | Secrets Manager path prefix |
| `AWS_SECRET_BATCH`  | Yes | Batch secret name |

### Secrets (resolved from AWS Secrets Manager at startup)

| Key | Description |
|-----|-------------|
| `SNS_TOPIC_ARN_PREFIX` | ARN prefix for SNS topics (e.g. `arn:aws:sns:us-east-1:123:komodo-events-`) |
| `JWT_PUBLIC_KEY` | RSA public key for internal JWT validation |
| `JWT_ISSUER` | Expected issuer claim |
| `JWT_AUDIENCE` | Expected audience claim |
| `MAX_CONTENT_LENGTH` | Max request body bytes |
| `RATE_LIMIT_RPS` | Rate limit (requests/sec) |
| `RATE_LIMIT_BURST` | Rate limit burst capacity |

---

## Local Development

### Prerequisites
- LocalStack running: `just up` from repo root (LocalStack supports SNS/SQS)

### Run
```bash
cd apis/komodo-events-api
source .env.local
go run ./cmd/server
```

### Run tests
```bash
go test ./...
```
