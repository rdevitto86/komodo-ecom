# komodo-payments-api

Payment processing, refunds, payment methods, and installment plan management.

**V1 is Rust (Axum).** The Go scaffold (`komodo-payments-api`) is abandoned.

| Key | Value |
|-----|-------|
| Port | 7071 (public), 7072 (internal, planned) |
| Domain | Payments |
| Status | Stub — implement handlers + DynamoDB repo |
| Language | Rust (Axum 0.7, Tokio) |
| Deployment | AWS Lambda (target) / EC2 docker-compose (bootstrap) |

---

## Routes

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/health` | None | Liveness check |
| `POST` | `/payments/charge` | Internal JWT | Charge a payment method |
| `POST` | `/payments/refund` | Internal JWT | Refund a charge (full or partial) |
| `GET` | `/payments/:charge_id` | Internal JWT | Get charge status |
| `GET` | `/me/payments` | User JWT | List user's payment history |
| `POST` | `/payments/webhook` | Stripe signature | Stripe webhook receiver |
| `GET` | `/me/payments/methods` | User JWT | List saved payment methods |
| `POST` | `/me/payments/methods` | User JWT | Add a payment method |
| `DELETE` | `/me/payments/methods/:method_id` | User JWT | Remove a payment method |
| `POST` | `/payments/plans` | Internal JWT | Create a payment plan (installments) |
| `GET` | `/me/payments/plans` | User JWT | List user's payment plans |
| `GET` | `/me/payments/plans/:plan_id` | User JWT | Get a payment plan |
| `DELETE` | `/me/payments/plans/:plan_id` | User JWT | Cancel a payment plan |
| `POST` | `/internal/payments/plans/:plan_id/execute` | Internal JWT | Execute next installment (scheduled job) |

---

## Environment Variables

### Process env

| Variable | Required | Description |
|----------|----------|-------------|
| `APP_NAME` | Yes | `komodo-payments-api` |
| `ENV` | Yes | `local`, `dev`, `staging`, `prod` |
| `LOG_LEVEL` | Yes | `debug`, `info`, `error` |
| `PORT` | Yes | e.g. `7071` |
| `AWS_REGION` | Yes | e.g. `us-east-1` |
| `AWS_ENDPOINT` | No | LocalStack endpoint or empty for real AWS |
| `AWS_SECRET_PREFIX` | Yes | Secrets Manager path prefix |
| `AWS_SECRET_BATCH` | Yes | Batch secret name |

### Secrets (resolved from AWS Secrets Manager at startup)

| Key | Description |
|-----|-------------|
| `DYNAMODB_PAYMENTS_TABLE` | Payments table name |
| `STRIPE_SECRET_KEY` | Stripe API secret key |
| `STRIPE_WEBHOOK_SECRET` | Stripe webhook signing secret |
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
cd apis/komodo-payments-api
source .env.local
cargo run
```

### Test
```bash
cargo test
cargo test -- --ignored   # run ignored integration tests (requires localstack)
```

---

## DynamoDB Table Design

See `docs/data-model.md` (TODO: write when implementing repo layer).

Primary key pattern: `PK=CHARGE#<uuid>`, `SK=METADATA` for charges.
Payment methods: `PK=USER#<user_id>`, `SK=METHOD#<id>`.
Plans: `PK=PLAN#<plan_id>`, `SK=METADATA` + `SK=INSTALLMENT#<n>`.
