# Komodo Monorepo — AI Context

## Project Purpose
Portfolio-grade e-commerce platform. Personal project with a realistic path to a real small business. Architecture decisions should be cost-efficient today with a clear AWS scaling path.

---

## 🚦 Active Mode

| Mode | Trigger | Role | Full Rules |
|------|---------|------|------------|
| **ADVISOR** (default) | No prefix | Senior backend peer — challenge, guide, never implement | `.agents/advisor.md` |
| **SENIOR** | `[SWE]` | Implements with judgment — brief flags, then executes | `.agents/swe.md` (backend + frontend) |

---

## ADVISOR Protocol
See `.agents/advisor.md` for the full role definition. Summary:

| Protocol | Behavior |
|----------|----------|
| Trade-offs First | Lead with non-obvious implications — partition costs, race conditions, scaling ceilings |
| Challenge | Ask "have you considered X?" before approving any design |
| Ask Before Showing | Request an attempt first. If stuck: *"Hint or answer?"* |
| Snippet-Only | No full-file rewrites. Targeted snippets with exact placement |
| Flag, Don't Fix | Surface mistakes; let the developer reason through the fix |
| `[Q]` | Direct answer, no mentorship overhead |

---

## Repo Layout

```
komodo-ecom/
├── apis/           # Go microservices + SDKs
│   ├── komodo-*-api/        # Independently deployable Go services
│   ├── komodo-forge-sdk-go/ # Shared Go SDK (referenced by all Go services)
│   └── komodo-forge-sdk-ts/ # Shared TypeScript SDK
├── infra/          # All infrastructure config
│   ├── local/               # Local dev (docker-compose, LocalStack, services.jsonc)
│   └── deploy/              # AWS deployment (CloudFormation, EC2 compose, scripts)
├── ui/             # SvelteKit 5 frontend (bun runtime)
├── Brewfile        # Repo toolchain (just, jq, go, bun, gh, docker)
└── Justfile        # Local dev task runner
```

> **Docker build context:** All Go service `docker-compose.yaml` files use `context: ../..` (repo root) so the SDK can be `COPY`'d alongside the service. Run `docker compose` from inside `apis/<service>/` for standalone use.

## Local Dev Orchestration

`infra/local/docker-compose.yml` with profiles, driven by `just`. Toggle services in `infra/local/services.jsonc`.

| Command | Services started |
|---------|-----------------|
| `just up` | localstack + redis + auth-api (always) |
| `just up api` | + APIs enabled in `services.jsonc` |
| `just up ui` | + UI enabled in `services.jsonc` |
| `just up api ui` | + APIs + UI |
| `just up api ui support` | everything enabled in `services.jsonc` |
| `just up api dev` | APIs, routed to AWS dev endpoints |
| `just down` | stop all |

Individual service composes (`apis/<service>/docker-compose.yaml`) still work standalone — they reference `komodo-network` as external, so run `just up` first to create the network.

---

## Context Strategy
**Do not pre-load monorepo context.** Discover only what's relevant to the current task.

**Working inside a Go service (`apis/komodo-*`):**
1. Read `apis/<service>/docs/README.md` first — authoritative reference for routes, env vars, port, commands.
2. Pull other `docs/` files only if directly relevant (e.g. `data-model.md` for DynamoDB work, `openapi.yaml` for contract questions).
3. Do not read sibling service directories unless the task explicitly spans services.
4. Fall back to this file only for monorepo-wide conventions.
5. When using forge SDK packages, read the relevant source in `apis/komodo-forge-sdk-go/` — do not guess signatures. Key packages:
   - `http/server` — `server.Run` (Lambda/HTTP universal entrypoint)
   - `http/middleware` — `middleware.Chain`, auth, rate-limiter, request-id, CORS, etc. (see `http/middleware/exports.go`)
   - `http/errors` — `httpErr.SendError`, error codes (`Global`, `Auth`, `User`, `Payment`, etc.)
   - `aws/secrets-manager` — `secretsmanager.Bootstrap`, `GetSecret`, `GetSecrets`
   - `config` — `config.GetConfigValue`, `SetConfigValue`
   - `logging/runtime` — `logger.Info`, `logger.Error`, `logger.Warn`
   - `http/context` — context keys (`USER_ID_KEY`, `SESSION_ID_KEY`, `SCOPES_KEY`, etc.)

**Working inside the frontend (`ui/`):**
1. Read `ui/docs/` for architecture and design context.
2. `ui/CLAUDE.md` is the authoritative context file for that workspace.
3. When using forge SDK types or API clients, read `apis/komodo-forge-sdk-ts/` for exact shapes — do not guess.

**Working at the monorepo root:**
- Use root `README.md` as the service registry.
- Backend services live under `apis/komodo-*`. Frontend lives under `ui/`.
- Shared SDKs live under `apis/komodo-forge-sdk-go` and `apis/komodo-forge-sdk-ts`.
- Root `docs/` contains monorepo-wide and cross-service documentation (audit logging, ADRs, platform decisions). Check here before writing new platform-level docs.

---

## Service `/docs` Standard
Every service should maintain this structure. JUNIOR mode uses it as its primary context source.

| File | Purpose | JUNIOR edits? |
|------|---------|---------------|
| `README.md` | Routes, port, env vars, run commands | Yes |
| `openapi.yaml` | API contract spec | Yes (post-struct) |
| `architecture.md` | Service topology, dependencies, data flow | No |
| `design-decisions.md` | Key decisions with rationale | No |
| `data-model.md` | DynamoDB table design, GSIs, access patterns, cost notes | No |

---

## Tech Stack
- **Go services:** Go 1.26, `net/http` ServeMux — no Chi, no Gin
- **Frontend (`ui/`):** SvelteKit 5 + TypeScript (bun runtime). Currently `adapter-static` for S3/CloudFront demo. Switch to `adapter-node` when wiring real backend.
- **SSR engine (`apis/komodo-ssr-engine-svelte`):** Backend SSR service — pre-renders/caches component trees for performance-critical pages. Not active until real backend is live.
- **Auth:** OAuth 2.0 / JWT RS256 via `komodo-auth-api`
- **Databases:** DynamoDB, S3, Redis (planned)
- **Infra:** Docker + LocalStack locally; EC2 + docker-compose for first production deploy; CloudFormation/ECS Fargate as the scale-up path
- **SDKs:** `komodo-forge-sdk-go` (AWS clients, middleware, crypto, logging, concurrency), `komodo-forge-sdk-ts` (types, API clients)

## Deployment Strategy
> Current state: demo site only. Backend not deployed.

| Service | Compute | Status |
|---------|---------|--------|
| `ui` | S3 + CloudFront (static) | `build:demo` → deploy to S3 |
| `auth-api` | EC2 docker-compose | Ready — `deploy/ec2/` |
| `user-api` | EC2 docker-compose | Ready — `deploy/ec2/` |
| `shop-items-api` | EC2 docker-compose | Ready — `deploy/ec2/` |
| `cart-api` | EC2 docker-compose | Scaffolded — docs complete |
| `inventory-api` | Lambda | Scaffolded — TODO: implement |
| `address-api` | Lambda | TODO: add Lambda handler |
| `order-api` | Lambda | TODO: add Lambda handler |
| `payments-api` | Lambda | TODO: add Lambda handler |
| `communications-api` | Lambda | TODO: add Lambda handler |
| `features-api` | Lambda | TODO: add Lambda handler |
| `entitlements-api` | Lambda | TODO: add Lambda handler |

**Scale-up path:** `infra/deploy/cfn/` templates are ready. When EC2 hits its ceiling, run `deploy-infra.sh` + `deploy-services.sh` to migrate to ECS Fargate. No code changes required.

**GitHub Actions:** All workflow auto-triggers are disabled (manual `workflow_dispatch` only). Re-enable when backend is live — uncomment the `on:` blocks in `ci.yml` and `deploy-dev.yml`.

## Conventions
- **Routing:** `net/http` ServeMux pattern syntax — `GET /me/profile`, `DELETE /me/profile/{id}`
- **Errors:** RFC 7807 Problem+JSON. Wrap: `fmt.Errorf("op: %w", err)`
- **Logging:** `slog` JSON. `KomodoTextHandler` locally (string format + ANSI color), JSON in staging/prod. Use `logger.FromContext(ctx)` to attach requestId/correlationId/userId/sessionId.
- **Auth:** JWT validated via forge SDK middleware on all protected routes
- **Context:** `context.Context` through every layer — handler → service → repo
- **DI:** Constructor functions, accept interfaces, return structs
- **Tests:** `go test ./...` from service root. `*_test.go` adjacent to source

## Port Allocation
> Local dev only. AWS Fargate/Lambda don't use host ports. Port 7000 reserved (macOS conflict).

**Convention within each domain block:**
- `anchor` — primary service, public-facing port
- `anchor+1` — internal-only server for the same service (if it has one)
- `anchor+2` onward — additional services in the domain (each may also claim `+1` for internal)
- Last 2–3 slots in every block are reserved for future growth

| Range | Domain | Assigned | Reserved |
|-------|--------|----------|---------|
| 7001–7010 | Frontend & Infrastructure | 7001 `ui`, 7002 `events-api`, 7003 `ssr-engine-svelte` | 7004–7010 |
| 7011–7020 | Identity & Security | 7011 `auth-api` pub, 7012 `auth-api` int | 7013–7020 |
| 7021–7030 | Core Platform | 7021 `entitlements-api`, 7022 `features-api` | 7023–7030 |
| 7031–7040 | Address & Geo | 7031 `address-api` | 7032–7040 |
| 7041–7050 | Commerce & Catalog | 7041 `shop-items-api`, 7042 `search-api`, 7043 `cart-api`, 7044 `inventory-api` | 7045–7050 |
| 7051–7060 | User & Profile | 7051 `user-api` pub, 7052 `user-api` int | 7053–7060 |
| 7061–7070 | Orders | 7061 `order-api`, 7062 `returns-api` | 7063–7070 |
| 7071–7080 | Payments | 7071 `payments-api` | 7072–7080 |
| 7081–7090 | Communications | 7081 `communications-api` | 7082–7090 |
| 7091–7100 | Loyalty & Social | 7091 `loyalty-api`, 7092 `reviews-api` | 7093–7100 |
| 7101–7110 | Support & CX | 7101 `support-api` | 7102–7110 |
| 7111–7120 | Analytics & Discovery | — | 7111–7120 |
