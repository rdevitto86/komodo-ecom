# komodo-shop-items-api

Product and service catalog API for the Komodo platform. Serves item data and inventory from S3, with authenticated personalized suggestions.

---

## Port

| Server | Port | Env Var |
|--------|------|---------|
| Public | 7041 | `PORT` |

---

## Routes

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/health` | None | Liveness check |
| `GET` | `/item/inventory` | None | Bulk inventory and stock status |
| `GET` | `/item/{sku}` | None | Single product or service by SKU |
| `POST` | `/item/suggestion` | Bearer JWT | Personalized product suggestions |

---

## Middleware

**Public routes** (`/item/inventory`, `/item/{sku}`):
RequestID → Telemetry → RateLimiter → IPAccess → CORS → SecurityHeaders

**Protected routes** (`/item/suggestion`):
RequestID → Telemetry → RateLimiter → IPAccess → CORS → SecurityHeaders → Auth → CSRF → Normalization → Sanitization → RuleValidation

---

## S3 Bucket Layout

```
s3://<S3_ITEMS_BUCKET>/
├── products/<sku>.json       # Product JSON per SKU
├── services/<sku>.json       # Service JSON per SKU
└── inventory/manifest.json   # Inventory manifest (all tracked items)
```

---

## Environment Variables

### Process env

| Variable | Required | Description |
|----------|----------|-------------|
| `APP_NAME` | Yes | Service name (`komodo-shop-items-api`) |
| `ENV` | Yes | Runtime environment (`local`, `dev`, `staging`, `prod`) |
| `LOG_LEVEL` | Yes | Log verbosity (`debug`, `info`, `error`) |
| `PORT` | Yes | HTTP listen port (default: `7041`) |
| `AWS_REGION` | Yes | AWS region (e.g. `us-east-1`) |
| `AWS_ENDPOINT` | Yes | AWS/LocalStack endpoint URL |
| `AWS_SECRET_PREFIX` | Yes | Secrets Manager path prefix |
| `AWS_SECRET_BATCH` | Yes | Secrets Manager batch secret name |

### Secrets (resolved from AWS Secrets Manager at startup)

| Key | Description |
|-----|-------------|
| `S3_ENDPOINT` | S3 endpoint URL |
| `S3_ACCESS_KEY` | S3 access key |
| `S3_SECRET_KEY` | S3 secret key |
| `S3_ITEMS_BUCKET` | S3 bucket for product/service/inventory JSON |
| `SHOP_ITEMS_API_CLIENT_ID` | Service client ID (for auth-api token requests) |
| `SHOP_ITEMS_API_CLIENT_SECRET` | Service client secret |
| `IP_WHITELIST` | Comma-separated allowed IPs (empty = allow all) |
| `IP_BLACKLIST` | Comma-separated blocked IPs |
| `MAX_CONTENT_LENGTH` | Max request body bytes |
| `IDEMPOTENCY_TTL_SEC` | Idempotency key TTL in seconds |
| `RATE_LIMIT_RPS` | Token bucket rate (requests/sec) |
| `RATE_LIMIT_BURST` | Token bucket burst capacity |
| `BUCKET_TTL_SECOND` | Rate limiter bucket TTL in seconds |

---

## Local Development

### Prerequisites

- LocalStack running: `just up` from repo root

### Run

```bash
cd apis/komodo-shop-items-api
source .env.local
go run .
```

### Docker

```bash
cd apis/komodo-shop-items-api
docker compose up --build
```

### cURL Examples

```bash
# Health check
curl http://localhost:7041/health

# Inventory manifest
curl http://localhost:7041/item/inventory

# Single item by SKU
curl http://localhost:7041/item/SKU-001

# Personalized suggestions (requires JWT)
TOKEN=$(curl -s -X POST http://localhost:7011/oauth/token \
  -H "Content-Type: application/json" \
  -d '{"clientId":"test-client","clientSecret":"test-secret","grantType":"client_credentials"}' | jq -r .accessToken)

curl -X POST http://localhost:7041/item/suggestion \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"userId":"usr_123","limit":5}'
```

---

## Status

**Active**

| Key | Value |
|-----|-------|
| Language | Go 1.26 |
| Router | `net/http` ServeMux |
| Port | 7041 |
| Domain | Commerce & Catalog |
| Data store | AWS S3 |
| SDK | `komodo-forge-sdk-go` |
| Docs | `docs/openapi.yaml` |
