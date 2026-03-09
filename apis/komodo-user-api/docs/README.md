# komodo-user-api

User profile, address, payment method, and preference management for the Komodo platform.

---

## Ports

| Server   | Port | Env Var         |
|----------|------|-----------------|
| Public   | 7051 | `PORT`          |
| Internal | 7052 | `INTERNAL_PORT` |

---

## Routes

### Public (`PORT`) — JWT required, user identity from token subject

| Method   | Path                    | Handler           | Description                        |
|----------|-------------------------|-------------------|------------------------------------|
| `GET`    | `/health`               | HealthHandler     | Liveness check                     |
| `GET`    | `/me/profile`           | GetProfile        | Get authenticated user's profile   |
| `POST`   | `/me/profile`           | CreateUser        | Create user record on registration |
| `PUT`    | `/me/profile`           | UpdateProfile     | Update authenticated user's profile|
| `DELETE` | `/me/profile`           | DeleteProfile     | Delete authenticated user's account|
| `GET`    | `/me/addresses`         | GetAddresses      | List all addresses                 |
| `POST`   | `/me/addresses`         | AddAddress        | Add a new address                  |
| `PUT`    | `/me/addresses/{id}`    | UpdateAddress     | Update an address by ID            |
| `DELETE` | `/me/addresses/{id}`    | DeleteAddress     | Delete an address by ID            |
| `GET`    | `/me/payments`          | GetPayments       | List saved payment methods         |
| `PUT`    | `/me/payments`          | UpsertPayment     | Add or update a payment method     |
| `DELETE` | `/me/payments/{id}`     | DeletePayment     | Remove a payment method by ID      |
| `GET`    | `/me/preferences`       | GetPreferences    | Get user preferences               |
| `PUT`    | `/me/preferences`       | UpdatePreferences | Update user preferences            |
| `DELETE` | `/me/preferences`       | DeletePreferences | Delete user preferences            |

### Internal (`INTERNAL_PORT`) — service-to-service, `svc:` scoped JWT required

| Method | Path                          | Handler        | Description                         |
|--------|-------------------------------|----------------|-------------------------------------|
| `GET`  | `/health`                     | HealthHandler  | Liveness check                      |
| `GET`  | `/users/{id}`                 | GetProfile     | Get profile by user ID              |
| `GET`  | `/users/{id}/addresses`       | GetAddresses   | Get addresses for a user            |
| `GET`  | `/users/{id}/preferences`     | GetPreferences | Get preferences for a user          |
| `GET`  | `/users/{id}/payments`        | GetPayments    | Get payment methods for a user      |

---

## Environment Variables

### Process env (set at container/process level)

| Variable            | Required | Description                                          |
|---------------------|----------|------------------------------------------------------|
| `APP_NAME`          | Yes      | Service name for logging (`komodo-user-api`)         |
| `ENV`               | Yes      | Runtime environment (`local`, `dev`, `staging`, `prod`) |
| `LOG_LEVEL`         | Yes      | Log verbosity (`debug`, `info`, `error`)             |
| `PORT`              | Yes      | Public server port (default: `7051`)                 |
| `INTERNAL_PORT`     | Yes      | Internal server port (default: `7052`)               |
| `VERSION`           | No       | Deployed version tag                                 |
| `AWS_REGION`        | Yes      | AWS region (e.g. `us-east-1`)                        |
| `AWS_ENDPOINT`      | Yes      | AWS/LocalStack endpoint URL                          |
| `AWS_SECRET_PREFIX` | Yes      | Secrets Manager path prefix (e.g. `komodo-user-api/local`) |
| `AWS_SECRET_BATCH`  | Yes      | Batch secret path (e.g. `/all-secrets`)              |
| `EVAL_RULES_PATH`   | No       | Path to validation rules file                        |

### Secrets (resolved from AWS Secrets Manager at startup)

| Key                      | Description |
|--------------------------|-------------|
| `DYNAMODB_ENDPOINT`      | DynamoDB endpoint URL |
| `DYNAMODB_ACCESS_KEY`    | DynamoDB AWS access key |
| `DYNAMODB_SECRET_KEY`    | DynamoDB AWS secret key |
| `DYNAMODB_TABLE`         | DynamoDB table name (`komodo-users`) |
| `USER_API_CLIENT_ID`     | Service client ID (used to obtain tokens from auth-api) |
| `USER_API_CLIENT_SECRET` | Service client secret |
| `JWT_PUBLIC_KEY`         | RSA public key (PEM) from auth-api — used to validate incoming tokens |
| `JWT_PRIVATE_KEY`        | RSA private key (PEM) — required by `InitializeKeys()`; not used for signing |
| `JWT_KID`                | Key ID (`test-kid` locally) |
| `JWT_ISSUER`             | Expected issuer claim (`test-issuer` locally) |
| `JWT_AUDIENCE`           | Expected audience claim (`test-audience` locally) |
| `IP_WHITELIST`           | Comma-separated allowed IPs (empty = allow all) — public only |
| `IP_BLACKLIST`           | Comma-separated blocked IPs — public only |
| `MAX_CONTENT_LENGTH`     | Max request body bytes — public only |
| `IDEMPOTENCY_TTL_SEC`    | Idempotency key TTL in seconds — public only |
| `RATE_LIMIT_RPS`         | Token bucket rate (requests/sec) — public only |
| `RATE_LIMIT_BURST`       | Token bucket burst capacity — public only |
| `BUCKET_TTL_SECOND`      | Rate limiter bucket TTL in seconds — public only |

---

## Local Development (no Docker for the app)

### Prerequisites

- LocalStack running (provides SM, DynamoDB): `just up` from repo root
- LocalStack initialized: init scripts in `infra/local/localstack/init/` run automatically

### Run (public + internal in separate terminals)

```bash
cd komodo-user-api

# Terminal 1 — public API on :7051
source .env.local && go run ./cmd/public

# Terminal 2 — internal API on :7052
source .env.local && go run ./cmd/internal
```

### cURL Examples

```bash
# Get a token from auth-api (must be running on :7011)
TOKEN=$(curl -s -X POST http://localhost:7011/oauth/token \
  -H "Content-Type: application/json" \
  -d '{"clientId":"test-client","clientSecret":"test-secret","grantType":"client_credentials"}' | jq -r .accessToken)

# Health check (no auth)
curl http://localhost:7051/health

# Profile (returns 501 until implemented)
curl -H "Authorization: Bearer $TOKEN" http://localhost:7051/me/profile

# Internal — get user by ID (requires svc: scoped token)
SVC_TOKEN=$(curl -s -X POST http://localhost:7011/oauth/token \
  -H "Content-Type: application/json" \
  -d '{"clientId":"test-client","clientSecret":"test-secret","grantType":"client_credentials","scope":"svc:user-api"}' | jq -r .accessToken)
curl -H "Authorization: Bearer $SVC_TOKEN" http://localhost:7052/users/usr_123
```

---

## Commands

```bash
# Run all tests
go test ./...

# Start via monorepo (preferred)
just up api          # starts infra + user-api (if enabled in services.jsonc)

# Docker (standalone — requires komodo-network, run just up first)
cd apis/komodo-user-api
docker compose up --build
```

---

## Middleware Stacks

**Public read** (`GET` routes): RequestID → Telemetry → RateLimiter → CORS → SecurityHeaders → Auth → CSRF → Normalization → RuleValidation → Sanitization

**Public write** (`POST/PUT/DELETE` routes): same as read + Idempotency at the end

**Internal** (`/users/*` routes): RequestID → Telemetry → Auth → Scope (`svc:` prefix required via `RequireServiceScope`)

---

## Lambda Deployment (future)

`cmd/public` and `cmd/internal` are Lambda-ready via `server.Run()` in forge-sdk. When `AWS_LAMBDA_FUNCTION_NAME` is present, `Run()` switches to `lambda.Start(httpadapter.NewV2(...))` automatically. Deploy as two separate Lambda functions with separate IAM roles. IAM replaces the port-based internal/public split used in Docker/Fargate.
