# komodo-auth-api

OAuth 2.0 authorization server for the Komodo platform. Issues and validates RS256 JWTs for M2M (`client_credentials`) and user flows.

---

## Ports

| Server | Port | Env Var |
|--------|------|---------|
| Public | 7011 | `PORT`  |

---

## Routes

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET`  | `/health` | None | Liveness check |
| `GET`  | `/.well-known/jwks.json` | None | Public RSA key set (JWKS) |
| `POST` | `/oauth/token` | None | Issue access token (client_credentials, refresh_token, authorization_code) |
| `GET`  | `/oauth/authorize` | None | Authorization code redirect — requires login UI (not yet implemented) |
| `POST` | `/oauth/introspect` | Client token | Token introspection (RFC 7662) |
| `POST` | `/oauth/revoke` | Client token | Token revocation (RFC 7009) |

**Note:** `/oauth/introspect` and `/oauth/revoke` require a valid `client_credentials` Bearer token in the `Authorization` header (`ClientTypeMiddleware` + `AuthMiddleware`).

---

## Environment Variables

### Process env (set at container/process level)

| Variable            | Required | Description |
|---------------------|----------|-------------|
| `APP_NAME`          | Yes | Service name (`komodo-auth-api`) |
| `ENV`               | Yes | Runtime environment (`local`, `dev`, `staging`, `prod`) |
| `LOG_LEVEL`         | Yes | Log verbosity (`debug`, `info`, `error`) |
| `PORT`              | Yes | Public server port (e.g. `:7011`) |
| `AWS_REGION`        | Yes | AWS region (e.g. `us-east-1`) |
| `AWS_ENDPOINT`      | Yes | LocalStack endpoint (`http://localhost:4566`) or empty for real AWS |
| `AWS_SECRET_PREFIX` | Yes | Secrets Manager path prefix (e.g. `komodo-auth-api/local/`) |
| `AWS_SECRET_BATCH`  | Yes | Batch secret name (e.g. `all-secrets`) |

### Secrets (resolved from AWS Secrets Manager at startup)

| Key | Description |
|-----|-------------|
| `AWS_ELASTICACHE_ENDPOINT` | Redis endpoint (e.g. `localhost:6379`) |
| `AWS_ELASTICACHE_PASSWORD` | Redis password (empty for local) |
| `AWS_ELASTICACHE_DB`       | Redis DB index (e.g. `0`) |
| `JWT_PUBLIC_KEY`           | RSA public key (PEM) for JWKS and token validation |
| `JWT_PRIVATE_KEY`          | RSA private key (PEM) for token signing |
| `JWT_KID`                  | Key ID for key rotation (`test-kid` locally) |
| `JWT_ISSUER`               | Token issuer claim (`test-issuer` locally) |
| `JWT_AUDIENCE`             | Token audience claim (`test-audience` locally) |
| `IP_WHITELIST`             | Comma-separated allowed IPs (empty = allow all) |
| `IP_BLACKLIST`             | Comma-separated blocked IPs |
| `MAX_CONTENT_LENGTH`       | Max request body bytes (e.g. `4096`) |
| `IDEMPOTENCY_TTL_SEC`      | Idempotency key TTL in seconds |
| `RATE_LIMIT_RPS`           | Token bucket rate (requests/sec) |
| `RATE_LIMIT_BURST`         | Token bucket burst capacity |
| `BUCKET_TTL_SECOND`        | Rate limiter bucket TTL in seconds |

---

## Local Development (no Docker for the app)

### Prerequisites

- LocalStack running (provides SM, Redis): `just up` from repo root
- LocalStack initialized: init scripts in `infra/local/localstack/init/` run automatically on first start

### Run

```bash
cd komodo-auth-api
source .env.local
go run ./cmd/public
```

### cURL Examples

```bash
# Health check
curl http://localhost:7011/health

# JWKS (public keys)
curl http://localhost:7011/.well-known/jwks.json

# Get a client_credentials token (test credentials seeded by LocalStack)
curl -s -X POST http://localhost:7011/oauth/token \
  -H "Content-Type: application/json" \
  -d '{"clientId":"test-client","clientSecret":"test-secret","grantType":"client_credentials","scope":"svc:user-api"}' | jq

# Store the token
TOKEN=$(curl -s -X POST http://localhost:7011/oauth/token \
  -H "Content-Type: application/json" \
  -d '{"clientId":"test-client","clientSecret":"test-secret","grantType":"client_credentials"}' | jq -r .accessToken)

# Introspect the token (requires auth)
curl -s -X POST http://localhost:7011/oauth/introspect \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" | jq

# Revoke the token
curl -s -X POST http://localhost:7011/oauth/revoke \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"token\":\"$TOKEN\"}" | jq
```

---

## Middleware Stack

**All `/oauth/*` routes:** RequestID → Telemetry → RateLimiter → IPAccess → SecurityHeaders → Normalization → Sanitization → RuleValidation

**Protected routes** (`/oauth/introspect`, `/oauth/revoke`): same as above + ClientType → Auth

---

## Grant Types

| Grant | Status | Use Case |
|-------|--------|----------|
| `client_credentials` | Implemented | M2M — service-to-service authentication |
| `refresh_token` | Implemented | Token renewal from a valid refresh token |
| `authorization_code` | Not implemented | Requires SvelteKit login UI (future) |

---

## Lambda Deployment (future)

The `cmd/public` binary is Lambda-ready via `komodo-forge-sdk-go/http/server.Run()`. When `AWS_LAMBDA_FUNCTION_NAME` is set (injected automatically by the Lambda runtime), `Run()` switches from `ListenAndServe` to `lambda.Start(httpadapter.NewV2(...))`. No code changes required.

IAM replaces the public/internal port split: the auth-api has only a public function. Internal callers (user-api) acquire a service token via `client_credentials` and call the protected `/oauth/introspect` route.
