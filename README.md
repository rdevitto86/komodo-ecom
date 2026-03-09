# Komodo E-Commerce Monorepo

Production-style e-commerce platform. Independently deployable Go microservices, a SvelteKit SSR frontend, shared SDKs, and local AWS infrastructure via LocalStack.

---

## Repo Layout

```
komodo-ecom/
├── apis/                     # Go microservices + SDKs
│   ├── komodo-*-api/         # Independently deployable Go services
│   ├── komodo-forge-sdk-go/  # Shared internal Go SDK
│   └── komodo-forge-sdk-ts/  # Shared internal TypeScript SDK
├── ui/                       # SvelteKit 5 frontend (bun runtime, adapter-static for demo)
├── infra/                    # Local AWS emulation (DynamoDB, S3, Secrets Manager) + Redis
│   ├── deploy/               # AWS deployment (CloudFormation, EC2 compose, deploy scripts)
│   └── local/                # Local development setup (including LocalStack)
```

## Getting Started

**Prerequisites:** [Homebrew](https://brew.sh), Docker >= 24.x

```bash
# 1. Install Homebrew (if not already installed)
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# 2. Install just, then bootstrap everything else
brew install just
just bootstrap

# 3. Start services (toggle what's enabled in infra/local/services.jsonc)
just up              # infra only (localstack + redis + auth-api)
just up api          # + enabled APIs
just up api ui       # + enabled APIs + UI
just up api ui dev   # same, but routed to AWS dev endpoints

# Stop everything
just down

# Run the frontend standalone (infra must be running for komodo-network)
cd ui && bun run dev         # live mode (needs backend running)
cd ui && bun run dev:mock    # mock mode (no backend needed)
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
| Networking | All services share `komodo-network` (created by `infra/local/docker-compose.yml`) |

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

LocalStack (`infra/local/localstack/`) emulates AWS locally:

| Service | Purpose |
|---------|---------|
| Secrets Manager | Service secrets (DB passwords, API keys, JWT keys) |
| S3 | Product data, content, file storage |
| DynamoDB | User data (NoSQL) |
| RDS | Aurora-compatible relational DB (planned) |
| Redis | Sessions and caching (standalone container, port 6379) |

Init scripts in `infra/local/localstack/init/` seed all services on startup.
