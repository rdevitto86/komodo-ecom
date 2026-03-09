# Komodo E-Commerce Monorepo

Production-style e-commerce platform. Independently deployable Go microservices, a SvelteKit SSR frontend, shared SDKs, and local AWS infrastructure via LocalStack.

---

## Repo Layout

```
komodo-ecom/
├── apis/           # Go microservices, SDKs, and localstack infrastructure
│   ├── komodo-*-api/       # Independent Go services (each has its own go.mod)
│   ├── komodo-forge-sdk-go/ # Shared internal Go SDK
│   ├── komodo-forge-sdk-ts/ # Shared internal TypeScript SDK
│   └── localstack/          # Local AWS emulation (DynamoDB, S3, Secrets Manager, Redis)
├── ui/             # SvelteKit 5 frontend (SSR, adapter-node)
├── db/             # DB migrations and seed scripts (planned)
└── deploy/         # Shared deployment scripts and CloudFormation templates (planned)
```

## Getting Started

**Prerequisites:** Docker >= 24.x, Go 1.26+, Bun 1.2+ (for SvelteKit/TS SDK), Make

```bash
# Start local backing services (DynamoDB, S3, Secrets Manager, Redis)
cd apis/localstack && docker compose up -d

# Run a backend service (docker-compose context is apis/)
cd apis/komodo-<service> && docker compose up --build

# Run the frontend
cd ui && bun dev
# or via docker:
cd ui && docker compose --profile local up --build
```

---

## Conventions

| Concern | Convention |
|---------|-----------|
| API style | JSON over HTTP, REST |
| Routing | Go 1.26 `net/http` ServeMux (`GET /path/{id}`) |
| Auth | JWT (RS256) via `komodo-auth-api`. Service-to-service via client credentials. |
| Logging | `slog` structured JSON. `tint` locally, JSON in staging/prod. |
| Errors | RFC 7807 Problem+JSON |
| Schema | `docs/openapi.yaml` per service |
| Tracing | OpenTelemetry OTLP (planned) |
| Secrets | AWS Secrets Manager via `komodo-forge-sdk-go` at startup |
| Networking | All services share `komodo-network` (created by `apis/localstack` compose) |

---

## Services

| Port | Service | Location | Domain | Status |
|------|---------|----------|--------|--------|
| 7001 | `ui` | `ui/` | Frontend shell (SvelteKit host app) | Active |
| 7003 | `komodo-ssr-engine-svelte` | `apis/` | SSR backend — pre-renders/caches component trees, delivers HTML fragments to `ui` as slots | Active |
| 7011 | `komodo-auth-api` | `apis/` | Identity & Security | Active |
| 7021 | `komodo-core-entitlements-api` | `apis/` | Core Platform | Stub |
| 7022 | `komodo-core-features-api` | `apis/` | Core Platform | Stub |
| 7031 | `komodo-address-api` | `apis/` | Address & Geo | Active |
| 7041 | `komodo-shop-items-api` | `apis/` | Commerce & Catalog | Active |
| 7042 | `komodo-search-api` | `apis/` | Commerce & Catalog | Stub |
| 7051 | `komodo-user-api` | `apis/` | User & Profile | Active |
| 7061 | `komodo-order-api` | `apis/` | Orders | Scaffolded |
| 7071 | `komodo-payments-api` | `apis/` | Payments | Scaffolded |
| 7081 | `komodo-communications-api` | `apis/` | Communications | Scaffolded |
| 7091 | `komodo-loyalty-api` | `apis/` | Loyalty & Social | Scaffolded |
| 7092 | `komodo-reviews-api` | `apis/` | Loyalty & Social | Scaffolded |
| 7101 | `komodo-support-api` | `apis/` | Support & CX | Scaffolded |
| 7111 | `komodo-analytics-collector-api` | `apis/` | Analytics | Stub |

> Port override: set the `PORT` env var on any service.

### Service Status
- **Active** — Implemented and running
- **Scaffolded** — Directory structure exists, not yet implemented
- **Stub** — Empty module or `main.go` only

---

## Shared Libraries

**`komodo-forge-sdk-go`** (`apis/komodo-forge-sdk-go`) — Internal Go SDK. AWS clients (DynamoDB, S3, Secrets Manager, Redis, Aurora), HTTP middleware stack, JWT/JWKS crypto, structured logging, concurrency utilities.

**`komodo-forge-sdk-ts`** (`apis/komodo-forge-sdk-ts`) — Internal TypeScript SDK. Domain types, API client utilities, frontend helpers. Backend modules (logging, telemetry) are stubs.

---

## Infrastructure

LocalStack (`apis/localstack/`) emulates AWS locally:

| Service | Purpose |
|---------|---------|
| Secrets Manager | Service secrets (DB passwords, API keys, JWT keys) |
| S3 | Product data, content, file storage |
| DynamoDB | User data (NoSQL) |
| RDS | Aurora-compatible relational DB (planned) |
| Redis | Sessions and caching (standalone container, port 6379) |

Init scripts in `localstack/init/` seed all services on startup.
