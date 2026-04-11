# Komodo Monorepo — AI Context

## Project Purpose
E-commerce platform for a real business. Architecture decisions prioritize correctness, operational reliability, and cost efficiency — not skill demonstration. Initial deployment is low-cost AWS (EC2 + Lambda); the target state is fully production-grade on ECS Fargate with no shortcuts carried forward.

---

## Agent Model

The main chat session is the **orchestrator** (architect role by default). Specialized work is delegated to named agent sessions — either spawned inline by the orchestrator or run manually in dedicated terminal windows.

**Three layers:**

| Layer | How | Role |
|-------|-----|------|
| Orchestrator | Main chat | Plans, coordinates, delegates — does not implement |
| Inline agents | `Agent` tool | Short-lived subagents for contained tasks within the same session |
| Terminal agents | Dedicated `claude` session per agent type | Long-running or parallel work with focused context |

**Terminal agent sessions** are per-role, not per-service. A single SWE session handles all Go work; a single DevOps session handles all infra. Each session loads only the context relevant to its current task.

```
Main chat (orchestrator)
  ├── Agent tool → inline subagent (fast, same context window)
  ├── RemoteTrigger → named terminal session (background, parallel)
  └── OpenClaw MCP → background tasks, cron, cross-session memory
                      (TODO: configure in .claude/settings.json)

Terminal sessions (spin up manually or via orchestrator):
  claude --agent swe          # all implementation work
  claude --agent devops       # all infra / deploy work
  claude --agent qa           # all testing / review work
  claude --agent architect    # deep design sessions
```

Agent definitions live in `.claude/agents/`. The orchestrator routes to them by name.

---

## Repo Layout

```
komodo-ecom/
├── apis/           # Go microservices
│   └── komodo-*-api/        # Independently deployable Go services
├── infra/          # All infrastructure config
│   ├── local/               # Local dev (docker-compose, LocalStack, services.jsonc)
│   └── deploy/              # AWS deployment (CloudFormation, EC2 compose, scripts)
├── ui/             # SvelteKit 5 frontend (bun runtime)
├── Brewfile        # Repo toolchain (just, jq, go, bun, gh, docker)
└── Justfile        # Local dev task runner
```

> **Docker build context:** Individual service `docker-compose.yaml` files use `context: ..` (i.e. `apis/`) so paths like `COPY komodo-<service>-api/` resolve correctly. Run `docker compose` from inside `apis/<service>/` for standalone use.

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

## TODO Tracking

Three `TODO.md` files serve as the project's living backlog. Check the relevant one before starting any task and flag completed items when work is done.

| File | Scope |
|------|-------|
| `apis/TODO.md` | All API services |
| `ui/TODO.md` | Frontend (`ui/`) |
| `infra/TODO.md` | Infrastructure (`infra/`) |

**Rules:**
- Before starting work in any area, read the relevant `TODO.md` and note any items your task touches.
- When a task completes an item, call it out explicitly so the user can check it off.
- When creating new work that isn't already listed, suggest adding it to the relevant `TODO.md`.
- Do not modify `TODO.md` files autonomously — surface the suggestion and let the user decide.

---

## Context Strategy
**Do not pre-load monorepo context.** Discover only what's relevant to the current task.

**Working inside a Go service (`apis/komodo-*`):**
1. Read `apis/<service>/docs/README.md` first — authoritative reference for routes, env vars, port, commands.
2. Pull `openapi.yaml` (service root) for contract questions, or other `docs/` files only if directly relevant (e.g. `data-model.md` for DynamoDB work).
3. Do not read sibling service directories unless the task explicitly spans services.
4. Fall back to this file only for monorepo-wide conventions.
5. Before writing any models, error codes, or adapters, check `apis/<service>/pkg/<version>/` — custom types, error codes (`models/errors.go`), and exported adapters live here. Do not redefine what already exists.
6. When using forge SDK packages, read the source in the `komodo-forge-sdk-go` repo (`github.com/rdevitto86/komodo-forge-sdk-go`) — do not guess signatures. Key packages:
   - `http/server` — `server.Run` (Lambda/HTTP universal entrypoint)
   - `http/middleware` — `middleware.Chain`, auth, rate-limiter, request-id, CORS, etc. (see `http/middleware/exports.go`)
   - `http/errors` — `httpErr.SendError`, error codes (`Global`, `Auth`, `User`, `Payment`, etc.); read `codes.go` for available codes before defining new ones
   - `http/request` — `request.GetQueryParams`, `GetClientKey`, `GetRequestID`, etc.
   - `http/response` — `ResponseWriter` wrapper, `IsSuccess`, `IsError`, etc.
   - `aws/secrets-manager` — `secretsmanager.Bootstrap`, `GetSecret`, `GetSecrets`
   - `config` — `config.GetConfigValue`, `SetConfigValue`
   - `logging/runtime` — `logger.Info`, `logger.Error`, `logger.Warn`, `logger.Attr`
   - `http/context` — context keys (`USER_ID_KEY`, `SESSION_ID_KEY`, `SCOPES_KEY`, etc.)
   - `events` — `events.Event`, `EventType`, `Source`, `EntityType` constants, `events.New`

**Working inside the frontend (`ui/`):**
1. Read `ui/docs/` for architecture and design context.
2. `ui/CLAUDE.md` is the authoritative context file for that workspace.
3. When using forge SDK types or API clients, read the source in the `komodo-forge-sdk-ts` repo (`github.com/rdevitto86/komodo-forge-sdk-ts`) for exact shapes — do not guess.

**Working at the monorepo root:**
- Use root `README.md` as the service registry.
- Backend services live under `apis/komodo-*`. Frontend lives under `ui/`.
- Shared SDKs live in external repos: `github.com/rdevitto86/komodo-forge-sdk-go` (Go) and `github.com/rdevitto86/komodo-forge-sdk-ts` (TypeScript).
- Root `docs/` contains monorepo-wide and cross-service documentation (audit logging, ADRs, platform decisions). Check here before writing new platform-level docs.

---

## Service File Standard
Every service should maintain this structure. JUNIOR mode uses it as its primary context source.

| File | Location | Purpose | JUNIOR edits? |
|------|----------|---------|---------------|
| `openapi.yaml` | service root | API contract spec — machine-readable, consumed by codegen | Yes (post-struct) |
| `README.md` | `docs/` | Routes, port, env vars, run commands | Yes |
| `architecture.md` | `docs/` | Service topology, dependencies, data flow | No |
| `design-decisions.md` | `docs/` | Key decisions with rationale | No |
| `data-model.md` | `docs/` | DynamoDB table design, GSIs, access patterns, cost notes | No |

`openapi.yaml` lives at the service root (not in `docs/`) so codegen tools can reference it without path gymnastics: `../komodo-auth-api/openapi.yaml`.

---

## Tech Stack
- **Go services:** Go 1.26, `net/http` ServeMux — no Chi, no Gin
- **Frontend (`ui/`):** SvelteKit 5 + TypeScript (bun runtime). Currently `adapter-static` for S3/CloudFront (cheap initial hosting). Switch to `adapter-node` when wiring real backend.
- **SSR engine (`apis/komodo-ssr-engine-svelte`):** Backend SSR service — pre-renders/caches component trees for performance-critical pages. Not active until real backend is live.
- **Auth:** OAuth 2.0 / JWT RS256 via `komodo-auth-api`
- **Databases:** DynamoDB, S3, Redis (planned)
- **Infra:** Docker + LocalStack locally; EC2 + docker-compose for first production deploy; CloudFormation/ECS Fargate as the scale-up path
- **SDKs:** `komodo-forge-sdk-go` — `github.com/rdevitto86/komodo-forge-sdk-go` (AWS clients, middleware, crypto, logging, concurrency); `komodo-forge-sdk-ts` — `github.com/rdevitto86/komodo-forge-sdk-ts` (types, API clients)

## Deployment Strategy
> Current state: frontend live (static), backend not yet deployed. EC2 is the low-cost bootstrap path; ECS Fargate is the production target.

| Service | Compute | Status |
|---------|---------|--------|
| `ui` | S3 + CloudFront (static) | `build:demo` → deploy to S3 |
| `auth-api` | EC2 docker-compose | Ready — `deploy/ec2/` |
| `user-api` | EC2 docker-compose | Ready — `deploy/ec2/` |
| `shop-items-api` | EC2 docker-compose | Ready — `deploy/ec2/` |
| `cart-api` | EC2 docker-compose | Scaffolded |
| `shop-inventory-api` | EC2 docker-compose | Scaffolded |
| `event-bus-api` | EC2 docker-compose | Built, not deployed |
| `order-api` | EC2 docker-compose | Scaffolded |
| `order-returns-api` | Lambda | Scaffolded |
| `order-reservations-api` | EC2 docker-compose | Foundation built — TODO: DynamoDB + checkout flow |
| `search-api` | EC2 docker-compose | Foundation built — TODO: Typesense integration |
| `loyalty-api` | EC2 docker-compose | Scaffolded |
| `reviews-api` | EC2 docker-compose | Scaffolded |
| `support-api` | EC2 docker-compose | Implemented (in-memory — wire DynamoDB before prod) |
| `address-api` | Lambda | TODO: add Dockerfile + Lambda handler |
| `payments-api` | Lambda | TODO: add Lambda handler |
| `communications-api` | Lambda | TODO: add Lambda handler |
| `features-api` | Lambda | TODO: add Dockerfile + Lambda handler |
| `entitlements-api` | Lambda | TODO: add Dockerfile + Lambda handler |

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
| 7001–7010 | Frontend & Infrastructure | 7001 `ui`, 7002 `event-bus-api`, 7003 `ssr-engine-svelte` | 7004–7010 |
| 7011–7020 | Identity & Security | 7011 `auth-api` pub, 7012 `auth-api` int | 7013–7020 |
| 7021–7030 | Core Platform | 7021 `entitlements-api`, 7022 `features-api` | 7023–7030 |
| 7031–7040 | Address & Geo | 7031 `address-api` | 7032–7040 |
| 7041–7050 | Commerce & Catalog | 7041 `shop-items-api`, 7042 `search-api`, 7043 `cart-api`, 7044 `shop-inventory-api` | 7045–7050 |
| 7051–7060 | User & Profile | 7051 `user-api` pub, 7052 `user-api` int | 7053–7060 |
| 7061–7070 | Orders | 7061 `order-api`, 7062 `order-returns-api`, 7063 `order-reservations-api` | 7064–7070 |
| 7071–7080 | Payments | 7071 `payments-api` | 7072–7080 |
| 7081–7090 | Communications | 7081 `communications-api` | 7082–7090 |
| 7091–7100 | Loyalty & Social | 7091 `loyalty-api`, 7092 `reviews-api` | 7093–7100 |
| 7101–7110 | Support & CX | 7101 `support-api` | 7102–7110 |
| 7111–7120 | Analytics & Discovery | — | 7111–7120 |

> **Rust variants:** `komodo-payments-api-rust` and `komodo-shop-inventory-api-rust` are V2 experiments — no ports assigned until A/B migration begins.
