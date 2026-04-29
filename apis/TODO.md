# Komodo Platform — TODO

Priority guide: **[H]** = blocking UI simulation · **[M]** = important, not blocking · **[L]** = low priority (docs, testing)

APIs are ordered by how soon the UI needs them to simulate a real backend.

---

## komodo-auth-api
> Status: MVP complete. Token issuing, validation, and revocation (JTI + ElastiCache) work. Auth code flow missing.

- [ ] **[L]** Implement `authorization_code` grant flow (requires SvelteKit login UI to be live first)
- [ ] **[M]** Implement OTP generation and acceptance — generate a time-limited one-time code (e.g. 6-digit TOTP or random), store in Redis with TTL, expose `POST /otp/request` (sends code via communications-api) and `POST /otp/verify` (validates code, returns short-lived JWT or session token); use for passwordless login and email verification flows
- [ ] **[L]** Add unit tests for token signing, validation, and introspection
- [ ] **[L]** Migrate to Rust for version 3.0 — evaluate Axum-based implementation leveraging existing Rust patterns from other services

---

## komodo-user-api
> Status: Fully wired. All sub-item CRUD (addresses, payments, preferences) implemented and wired through service → repo. Internal ownership checks enforced via `resolveUserID`.

- [ ] **[L]** Add integration tests for public + internal handler paths

---

## komodo-shop-items-api
> Status: Fully implemented. S3-backed product catalog works. Suggestions use real inventory-based ranking. Repair routes live here pending relocation decision.

- [ ] **[M]** Relocate `/services/repair` routes to `komodo-order-reservations-api` — repair booking is a time-slot/appointment concern, not a catalog concern; shop-items-api should only expose repair *listings* (`service_type=repair`); the booking flow belongs in reservations-api; coordinate the route split before implementing either handler

---

## komodo-cart-api
> Status: Fully implemented. Guest + authenticated carts, checkout token generation, and stock hold coordination all work.

- [ ] **[M]** Design and implement "Save for Later" feature (separate DynamoDB entity, no TTL) — see README TODO
- [ ] **[L]** Add unit tests for guest cart TTL, merge flow, and checkout token lifecycle (e2e exists; handler/service/repo unit tests are stubs)

---

## komodo-shop-inventory-api
> Status: Stub complete — all layers scaffolded (Axum, models, repo trait, handlers). DynamoDB impl is `todo!()`.
> **Plan: V1 in Go (unblocks checkout hold flow on normal timeline), V2 migrate to Rust. Rust variant remains as a teaching experiment — no ports assigned until migration.**
> Downgraded from [H]: cart-api already coordinates stock holds at token level; manual out-of-stock handling is an acceptable operational gap pre-launch.

- [ ] **[M]** Implement `DynamoInventoryRepo::reserve` — conditional write (`available_qty >= requested`), write HOLD# record with TTL
- [ ] **[M]** Implement `DynamoInventoryRepo::get_stock` + `batch_stock` — GetItem / BatchGetItem for SKU#/STOCK records
- [ ] **[M]** Implement `DynamoInventoryRepo::release_hold` — DeleteItem HOLD# record, restore `available_qty`
- [ ] **[M]** Implement `DynamoInventoryRepo::confirm` — DeleteItem HOLD# record, decrement `reserved_qty`, increment `committed_qty`
- [ ] **[M]** Implement `DynamoInventoryRepo::restock` — UpdateItem `available_qty += qty`
- [ ] **[M]** Wire Secrets Manager bootstrap — populate `Config` secret fields at startup
- [ ] **[M]** Implement JWT RS256 validation in `middleware/auth.rs` — `DecodingKey::from_rsa_pem` + `jsonwebtoken::decode`
- [ ] **[M]** DynamoDB Streams handler (separate Lambda) — listen for TTL expiry events, restore `available_qty` on hold expiry
- [ ] **[M]** Wire `communications_api_url` to fire restock threshold alert when `available_qty` drops below `restock_threshold`
- [ ] **[L]** Implement `common::spawn_app()` in tests + enable integration tests

---

## komodo-order-api
> Status: Routes wired and returning stubbed responses. Unified order submission (`POST /v1/orders`) and lookup (`GET /v1/orders/{orderId}`) implemented with guest + registered identity model. Returns flow stubbed. Private state-transition routes (`PUT /v1/returns/{returnId}/approve|reject|receive`) registered. DynamoDB not yet wired.
>
> **Identity model:** email is the universal key for both guests and registered users. Every order carries an email. At placement, if the email matches a registered account the order is automatically linked to that `USER#<id>` — no separate guest/user routes needed, and guest conversion is zero-work. Order submission and lookup are unified and handle both identity types.

- [ ] **[H]** Wire DynamoDB for order persistence — implement `CreateOrder`, `GetOrder`, `ListOrders` (GSI on userId), `UpdateOrderStatus` in the repo layer; replace all stub responses
- [ ] **[H]** Implement `POST /v1/orders/{orderId}/cancel` — release stock holds via shop-inventory-api, trigger refund via payments-api; enforce cancellable state check (`pending` or `confirmed` only)
- [ ] **[H]** Wire return request persistence — implement `CreateReturn`, `GetReturn`, `ListReturns`, `UpdateReturnStatus` in DynamoDB; replace stub responses on `GET/POST /v1/orders/returns` and `GET/DELETE /v1/orders/returns/{returnId}`
- [ ] **[H]** Implement private return state transitions — wire `PUT /v1/returns/{returnId}/approve` (trigger refund via payments-api), `PUT /v1/returns/{returnId}/reject`, and `PUT /v1/returns/{returnId}/receive` (trigger restock via shop-inventory-api and loyalty reversal)
- [ ] **[M]** Implement order status state machine: `pending → confirmed → shipped → delivered → cancelled` — enforce transition rules in service layer; reject invalid transitions with 409
- [ ] **[M]** Publish `order.placed`, `order.cancelled`, and `order.fulfilled` events to event-bus-api
- [ ] **[L]** Add integration tests for order creation, guest lookup, account auto-link, cancellation flow, and return lifecycle

---

## komodo-payments-api
> Status: Stub complete — all layers scaffolded (Axum, models, repo trait, Stripe provider, handlers). DynamoDB impl and Stripe calls are `todo!()`.

- [ ] **[H]** Implement `DynamoPaymentsRepo::save_charge` / `get_charge` — PK=CHARGE#<uuid>, SK=METADATA
- [ ] **[H]** Implement `DynamoPaymentsRepo::save_refund` — PK=CHARGE#<charge_id>, SK=REFUND#<refund_id>
- [ ] **[H]** Implement `DynamoPaymentsRepo::add_method` / `list_methods` / `delete_method` — PK=USER#<user_id>, SK=METHOD#<id>
- [ ] **[H]** Implement `StripeClient::charge` — POST `/v1/payment_intents` with idempotency key
- [ ] **[H]** Implement `StripeClient::refund` — POST `/v1/refunds`
- [ ] **[H]** Wire Secrets Manager bootstrap — populate `Config` secret fields at startup
- [ ] **[H]** Implement JWT RS256 validation in `middleware/auth.rs`
- [ ] **[M]** Implement `DynamoPaymentsRepo` plan methods — PK=PLAN#<plan_id>, installments as SK=INSTALLMENT#<n>
- [ ] **[M]** Implement `handlers/methods::execute_installment` — find next `Scheduled` installment, call `provider.charge()`, update status
- [ ] **[M]** Implement Stripe webhook validation — verify `Stripe-Signature` header using `STRIPE_WEBHOOK_SECRET`
- [ ] **[M]** Publish `payment.succeeded` / `payment.failed` / `payment.refunded` events to event-bus-api
- [ ] **[M]** Publish payment plan events (`payment.plan.created`, `payment.plan.installment.charged`, etc.) to event-bus-api
- [ ] **[M]** Write `docs/data-model.md` — finalize DynamoDB table schema
- [ ] **[M]** Enforce autopay requires payment method on file — validate that the user has an active bank account or credit card record before enabling autopay or processing any autopay transaction; return a clear error if no method exists; may require an `autopay_enabled` boolean and method-presence check in the DB schema
- [ ] **[L]** Implement `common::spawn_app()` in tests + enable integration tests with Stripe test mode

---

## komodo-address-api
> Status: Implemented but provider calls are stubs. Stateless — no DB needed.

- [ ] **[H]** Wire real address validation provider (SmartyStreets, Google Address Validation, or similar) — replace all 3 stub `TODO` bodies in `internal/provider/address.go`
- [ ] **[M]** Add provider API key secret (`ADDRESS_PROVIDER_API_KEY`) to LocalStack init seed
- [ ] **[L]** Add unit tests for provider error handling and response mapping

---

## komodo-search-api
> Status: Partial — middleware and routes wired, net/http ServeMux in use. All handlers return empty results. Typesense not initialized.

- [ ] **[H]** Initialize Typesense client after secrets load (`TODO(typesense)` in main.go)
- [ ] **[H]** Implement `GET /search` — build query params from request, call Typesense, return results
- [ ] **[H]** Implement `POST /v1/index/sync` — full re-index from shop-items-api S3 data into Typesense (previously referred to as `/internal/index/sync` — path corrected to match ROUTES.md)
- [ ] **[M]** Wire event-bus-api subscriber to listen for `shop_item.created/updated/deleted` → incremental index updates
- [ ] **[M]** Implement `DELETE /v1/index` — drop and recreate Typesense collection for schema migrations (previously referred to as `/internal/index` — path corrected to match ROUTES.md)
- [ ] **[L]** Add integration tests for search query building and index sync

---

## komodo-support-api
> Status: Fully implemented with Anthropic Haiku integration. In-memory storage not production-safe.

- [ ] **[H]** Replace in-memory repository with DynamoDB — design table schema and implement all repo functions (`repository/chat.go:24`)
- [ ] **[M]** Add wildcard `GET /v1/chat/history` for all user messages — when the `session` query param is omitted, return all messages across all sessions for the authenticated user; requires JWT and DynamoDB GSI on `user_id`; currently `session` is required (noted as TODO in openapi.yaml)
- [ ] **[M]** Define audit event schema and destination before wiring deletion audit — DynamoDB audit table or S3 archive (`repository/chat.go:18`)
- [ ] **[M]** Wire escalation (`POST /chat/escalate`) to communications-api for async ticket creation
- [ ] **[M]** Replace SQS publish placeholder in escalation handler with real SQS client once forge SDK has SQS support (`handlers/chat.go:182`)
- [ ] **[M]** Emit audit event to event-bus-api before chat history deletion (compliance trail)
- [ ] **[M]** Complete anonymous → authenticated session merge flow (currently not wired end-to-end)
- [ ] **[L]** Design human agent handoff flow (currently no handoff target exists)

---

## komodo-event-bus-api
> Status: Functional for local dev. DynamoDB persistence, HTTP dispatch, and event type allowlist all work. SNS/SQS wired but not deployed. CDC classifiers incomplete.

- [ ] **[M]** Deploy and test SNS/SQS fan-out path (in-memory fan-out is local-only)
- [ ] **[M]** Wire and deploy CDC Lambda handler for DynamoDB Streams
- [ ] **[M]** Add CDC event classifiers for payments, users, inventory, and cart domains — only orders classifier exists (`cdc/domains/orders.go`)
- [ ] **[M]** Expand CDC order event payload with additional fields: `total_cents`, `item_count`, `customer_id` (`cdc/domains/orders.go:48`)
- [ ] **[M]** Wire EventBridge as the routing layer for CDC events — add EventBridge rules for flexible per-domain fanout (orders → order consumers, payments → payment consumers, etc.)
- [ ] **[M]** Define and publish interaction events (`cart.item_added`, `order.started`, `order.abandoned`) — extend event type catalogue for analytics consumers
- [ ] **[M]** Per-connector publisher workers + DLQ — each outbound sink (SNS, EventBridge, S3 data lake, in-memory fan-out) runs its own goroutine/worker so one slow or failing sink cannot block the others; failed publishes route to a per-connector DLQ (SQS or in-memory bounded queue) with retry policy and max-age eviction
- [ ] **[L]** Emit CloudWatch metric (or fixed-key structured log) on unroutable CDC events (`cdc/handler.go:51`)
- [ ] **[L]** Evaluate gRPC as an additional internal transport — research protobuf schema enforcement, bi-directional streaming, and performance vs complexity tradeoff compared to current HTTP REST; document findings in `docs/design-decisions.md` before any implementation; not a blocking item

---

## komodo-order-reservations-api
> Status: Routes wired, middleware configured, but all repository functions are stubs (15+ TODOs).
>
> **Identity model:** same email-as-universal-key pattern as order-api. `POST /reservations` is a unified route (optional JWT, email always required); if email matches a registered account the booking is auto-linked. `GET /reservations/{id}` accepts optional JWT or `email` query param fallback. `/me/reservations` is JWT-required for account-scoped booking history. Repair reservations follow the same model — guest or registered, single route, email links the record.

- [ ] **[M]** Initialize DynamoDB client in bootstrap (blocked on forge SDK `aws/dynamodb` availability — confirm package path)
- [ ] **[M]** Implement `repo.GetBooking` / `CreateBooking` / `UpdateBooking` (DynamoDB)
- [ ] **[M]** Implement `repo.GetSlots` / `UpdateSlotAvailability` (DynamoDB)
- [ ] **[M]** Migrate identity model — replace `customer_id` JWT extraction with email-based linking: extract email from JWT claims if present, otherwise require email in request body; resolve to `USER#<id>` if registered, `GUEST#<uuid>` if not
- [ ] **[M]** Add ownership check on booking reads/mutations — JWT path checks subject match; email path validates email matches booking record
- [ ] **[M]** Confirm `POST /v1/slots/sync` route is correctly wired to the private middleware stack (ROUTES.md path is `/v1/slots/sync` — old references to `/internal/slots/sync` are stale)
- [ ] **[M]** Decide and implement checkout hold flow (Option A: hold at reservation time vs Option B: hold at order confirm)
- [ ] **[M]** Extend booking model for repair intake — add `repair` booking type with fields: `device_type`, `serial_number`, `reported_issue`, `inbound_shipment_id`; wire to inbound shipping flow once shipping-api exists
- [ ] **[M]** Implement repair status state machine: `intake_pending → received → diagnosing → repairing → quality_check → ready → shipped_back`; emit status change events to event-bus-api on each transition
- [ ] **[L]** Add integration tests for booking lifecycle, guest auto-link, and repair intake flow

---

## komodo-shipping-api (NEW)
> Status: Not yet created. Lambda service, port 7064. Handles both inbound (customer → warehouse: returns, repair intake) and outbound (warehouse → customer: order fulfillment, repaired items) shipment label generation and tracking.
> Routes: `GET /v1/shipments/{shipmentId}` (pub), `POST /v1/labels/outbound` (pub — may move to private once specced), `POST /v1/labels/inbound` (pub — may move to private once specced), `POST /webhooks/carrier` (ext, signature-validated not JWT). See ROUTES.md for authoritative list.

- [ ] **[M]** Scaffold service — Lambda compute, port 7064; bootstrap (logger, secrets, DynamoDB), public middleware stack, ServeMux routes; add Docker + Lambda handler; register port 7064 in port allocation table (currently listed as `shipping-api (planned)`)
- [ ] **[M]** Select and integrate a carrier aggregator (EasyPost, ShipStation, or EasyPost-compatible) — abstract behind a provider interface so carriers are swappable
- [ ] **[M]** Implement `POST /v1/labels/outbound` — generate outbound label for order fulfillment; called by order-api when order transitions to `shipped`; return carrier, tracking number, and label URL; revisit whether this belongs on the private port before implementing
- [ ] **[M]** Implement `POST /v1/labels/inbound` — generate prepaid inbound return/repair label; called by order-api returns flow and order-reservations-api; customer receives label URL to print or QR scan; revisit whether this belongs on the private port before implementing
- [ ] **[M]** Implement `GET /v1/shipments/{shipmentId}` — real-time shipment status; poll carrier API or return latest cached status from DynamoDB
- [ ] **[M]** Implement `POST /webhooks/carrier` — receive status events from carrier (`delivered`, `out_for_delivery`, `exception`, `in_transit`); validate carrier signature (not JWT); update shipment record and publish `shipment.status_updated` event to event-bus-api
- [ ] **[M]** Publish `shipment.label.created`, `shipment.delivered`, `shipment.received.inbound` events — `shipment.received.inbound` triggers inspection/repair workflow in reservations-api; `shipment.delivered` triggers loyalty points and fulfillment confirmation in order-api
- [ ] **[L]** Add integration tests for label generation, status polling, and webhook handling

---

## komodo-communications-api
> Status: Not yet scaffolded — `cmd/private/main.go` is empty.

- [ ] **[L]** Scaffold service: bootstrap, middleware stack, ServeMux routes
- [ ] **[L]** Implement `POST /send/email` — transactional email via provider (SendGrid, SES, etc.)
- [ ] **[L]** Implement `POST /send/sms` — SMS via provider (Twilio, SNS, etc.)
- [ ] **[L]** Implement `POST /send/push` — in-app push notification
- [ ] **[L]** Subscribe to event-bus-api for async trigger events (`order.placed`, `order.shipped`, etc.)
- [ ] **[L]** Store and load transactional email templates from S3 — support per-locale variants; templates managed separately from code
- [ ] **[L]** Add template management for transactional messages
- [ ] **[L]** Add unit tests for provider client and template rendering
- [ ] **[L]** Add integration tests for email/SMS/push sending flows

---

## komodo-loyalty-api
> Status: Routes wired (health + reviews stub). All handlers return NotImplemented. Bootstrap, middleware, and event subscription missing.

- [ ] **[L]** Scaffold service: bootstrap, middleware stack, ServeMux routes (loyalty + reviews endpoints)
- [ ] **[L]** Implement points earn on `order.placed` event (subscribe via event-bus-api)
- [ ] **[L]** Implement points reversal on `order.returned` event
- [ ] **[L]** Implement `GET /me/loyalty` — points balance and tier status
- [ ] **[L]** Implement `POST /me/loyalty/redeem` — apply points discount to a cart/order
- [ ] **[L]** Design tier rules and rewards catalogue
- [ ] **[L]** Implement `POST /me/reviews` — submit a review (require verified purchase check via order-api)
- [ ] **[L]** Implement `GET /items/{itemId}/reviews` — paginated review listing with avg rating + count
- [ ] **[L]** Implement `PUT/DELETE /me/reviews/{reviewId}` — edit and remove reviews
- [ ] **[L]** Add moderation queue for flagged reviews
- [ ] **[L]** Add unit tests for points calculation, tier logic, and rating aggregation
- [ ] **[L]** Add integration tests for event subscription, redemption, and review CRUD

---

## komodo-features-api
> Status: Complete stub — merged with entitlements-api. `cmd/public/main.go` is a Hello World placeholder. Single access evaluation service: feature flags + role-based entitlements backed by S3 + Redis. Private-only for V1.
> S3 schema: `s3://komodo-platform-config/{env}/access.json` — top-level `roles` (role → scopes, inherits) and `flags` (flag key → enabled + allowed roles).
> Evaluation model: `access_granted = flag.enabled AND user.role IN flag.allowed_roles`
>
> **Naming note:** ROUTES.md lists port 7022 under `komodo-platforms-api`. This service is NOT being renamed — it covers the same access-control concern. If a separate platforms/tenant service is needed, create a new `komodo-platforms-api` service.

**V1 route implementation (spec-only — openapi.yaml updated)**
- [ ] **[L]** Implement `GET /v1/access?appId={appId}` — evaluate and return role, toggles, and entitlements for a given appId; backed by S3 `access.json` with Redis cache
- [ ] **[L]** Implement `PUT /v1/access` — write role, toggle, and entitlement overrides for a given appId; finalize request body schema (role, toggles, entitlements) before implementing

**V1 (ecom MVP)**
- [ ] **[L]** Scaffold service: bootstrap (logger, secrets, S3 client, Redis cache), private-only ServeMux, health check
- [ ] **[L]** Load and cache `access.json` from S3 on startup; refresh on Redis TTL expiry (60s)
- [ ] **[L]** Implement `GET /internal/access/evaluate?flag=<key>` — returns allowed/denied given JWT role claim; used by any service needing runtime flag checks
- [ ] **[L]** Implement `GET /internal/access/roles/{role}` — returns scopes for a role; called by auth-api at token issuance to embed scopes in JWT
- [ ] **[L]** Seed `access.json` for local dev in LocalStack S3 init — include all roles (customer, premium_customer, servicing_agent, developer, admin) and placeholder flags
- [ ] **[L]** Add unit tests for flag evaluation logic (role match, flag disabled, wildcard role)

**V2 (enterprise path — not blocking ecom)**
- [ ] **[L]** Add public route `GET /me/access` — returns active flags + role for the authenticated user (UI consumption, replaces script-based AD group checks)
- [ ] **[L]** Per-user entitlement overrides via S3 per-user files (`users/{userId}.json`) — written by order-api on premium purchase, admin action, or onboarding script; read and merged with role defaults at evaluation time
- [ ] **[L]** Subscribe to event-bus-api `order.placed` events — auto-grant `premium_customer` role override when a qualifying plan is purchased
- [ ] **[L]** Audit log: append-only writes to S3 (`audit/{env}/changes.jsonl`) on any role or flag change
- [ ] **[L]** AD/LDAP group sync path — resolve role from external directory at login; pass resolved role to auth-api for JWT embedding
- [ ] **[L]** Add integration tests for S3 config load, cache refresh, and evaluate endpoint

---

## komodo-ai-guardrails-api (Python FastAPI)
> Status: Spec-only — Python FastAPI service scaffolded at `apis/komodo-ai-guardrails-api/` but `POST /moderate` is not implemented. Port is 7023 (Core Platform block) — README incorrectly listed 7113; openapi.yaml corrected.

**V1 route implementation (spec-only — openapi.yaml updated)**
- [ ] **[L]** Implement `POST /v1/moderate` — run selected guardrail checks (pii, injection, deviation, coding, obscenity, toxicity) on input text; return flags, redacted text, and latency

**Service scaffolding**
- [ ] **[L]** Implement provider abstraction — `LocalSLMProvider` (Ollama or similar on-prem) and `BedrockProvider` (AWS Bedrock); toggle via env config (`SLM_PROVIDER=local|bedrock`); both implement the same interface so callers are provider-agnostic
- [ ] **[L]** Wire input sanitization and prompt injection protection on all endpoints — normalize whitespace, strip control characters, enforce max token limits before forwarding to the model
- [ ] **[L]** Add Bedrock IAM role and secret config to infra — `AWS_BEDROCK_MODEL_ID`, region, and credentials; add local stub values to LocalStack secrets init
- [ ] **[L]** Fix port in README — update from 7113 to 7023 to match CLAUDE.md port allocation
- [ ] **[L]** Add integration tests for each task endpoint against a local SLM stub

---

## komodo-statistics-api
> Status: Scaffolded. Public and internal stat routes registered with middleware stacks. SQLite client is a placeholder — no data persists yet.

- [ ] **[M]** Replace SQLite with DynamoDB — design table schema for stat counters (trending scores, in-cart counts, purchase counts, co-purchase pairs); implement DynamoDB repo layer; remove `modernc.org/sqlite` dependency; add table definition to `infra/deploy/cfn/infra.yaml`
- [ ] **[M]** Wire event consumption — subscribe to event-bus-api events (`cart.item_added`, `cart.item_removed`, `order.placed`, `order.fulfilled`, `shop_item.viewed`); the `POST /v1/events` route is wired but handler is a stub; update DynamoDB stat counters on each event
- [ ] **[M]** Implement `GET /v1/stats/items/trending` — query DynamoDB trending-score GSI, return ranked item list with configurable `limit` and `period` window
- [ ] **[M]** Implement `GET /v1/stats/items/{itemId}/in-cart` — read live in-cart counter from DynamoDB; return count and display label
- [ ] **[M]** Implement `GET /v1/stats/items/{itemId}/recently-bought` — read purchase-count counter from DynamoDB for the requested rolling window (7d / 30d / 90d)
- [ ] **[M]** Implement `GET /v1/stats/items/{itemId}/frequently-bought-with` — query co-purchase pair records from DynamoDB; return related items ranked by co-purchase count
- [ ] **[M]** Implement `GET /v1/stats/dashboard` (private, port 7112) — aggregate active carts, orders today/month, top in-cart items, and current trending list; serve to admin consumers via the private server
- [ ] **[M]** Background TTL / data eviction — use DynamoDB native TTL attribute on stat counter records to prevent unbounded growth; remove the SQLite-era TTL cleanup goroutine

---

## komodo-insights-api
> Status: Stub — handler skeletons only. LLM provider not initialized. Backed by Bedrock/Claude API calls.

- [ ] **[M]** Wire LLM provider — initialize Bedrock/Claude client at service startup; pull model ID and credentials from Secrets Manager; return 503 with `Retry-After` header if provider is unavailable
- [ ] **[M]** Implement `GET /v1/items/trending` — call LLM to derive trending item signals from recent purchase/view data; cache result (see caching item below)
- [ ] **[M]** Implement `GET /v1/items/{itemId}/summary` — generate an AI summary for the item using product data as context; cache per itemId
- [ ] **[M]** Implement `GET /v1/items/{itemId}/sentiment` — run sentiment analysis over item reviews using LLM; return sentiment enum and score in [-1.0, 1.0]; cache per itemId
- [ ] **[M]** Add local response cache for Bedrock/Claude API calls — cache by summary/query type with an expiry timestamp; serve cached result if not expired; refresh on weekly/biweekly schedule or when threshold exceeded; prevents repeated expensive LLM calls on every request

---

## komodo-shop-promotions-api
> Status: Not yet scaffolded. Handles promotions, discount logic, and first-order tracking.

- [ ] **[M]** Track first-order flag per account — on `order.placed` event, record whether this is the account's first order (guest or registered); store in promotions DB as a boolean/timestamp record; used to gate first-order discount eligibility
- [ ] **[L]** Add unit tests for discount calculation and first-order eligibility logic
- [ ] **[L]** Add integration tests for promotion CRUD operations

---

## komodo-ssr-engine-svelte
> Status: Spec-only — openapi.yaml complete, no routes implemented. Bun runtime, Fargate, port 7003, private-only.

- [ ] **[M]** Implement `POST /v1/render` — accept component identifier and props, server-side render the specified SvelteKit component tree, return HTML fragment and cache-hit flag; use Bun's native SSR APIs; apply cache key derivation (component + props hash if `cache_key` not supplied)
- [ ] **[M]** Implement `POST /v1/admin/content/invalidate` — accept a list of cache keys and evict the corresponding rendered fragments from the in-process or Redis cache; M2M JWT required
- [ ] **[M]** Implement `POST /v1/admin/content/upsert` — insert or replace a content record in the SSR cache; useful for pre-warming critical pages; M2M JWT required
- [ ] **[M]** Choose and wire cache backend — in-process LRU for local dev; Redis (forge SDK) for production; TTL per render configurable via `ttl_seconds` in the request; default server TTL via env config
- [ ] **[L]** Add integration tests for render, invalidate, and upsert flows using a minimal SvelteKit component fixture

---

## Cross-Cutting

- [ ] **[M]** Add `Retry-After` header support for transient errors across all APIs — return appropriate seconds value on rate limits (429), service unavailability (503), and provider timeouts; signal retryability via standard HTTP semantics instead of error codes
- [ ] **[M]** Establish shared event type catalogue (`event-bus-api` has it defined but downstream services aren't validating against it)
- [ ] **[M]** Wire `support-api` escalation → `communications-api` once communications-api is scaffolded
- [ ] **[L]** Add `docs/data-model.md` to every API that uses DynamoDB (currently missing for all services)
- [ ] **[L]** Write unit + integration tests for all services (`go test ./...` must pass; at minimum happy path + error cases per handler)
- [ ] **[M]** Wire inbound shipping to returns flow — `order-returns-api` calls `shipping-api` `POST /shipping/labels/inbound` on RMA approval and returns the label URL to the customer
- [ ] **[M]** Wire inbound shipping to repair intake — `order-reservations-api` calls `shipping-api` `POST /shipping/labels/inbound` on repair booking confirmation; store `inbound_shipment_id` on the booking record
- [ ] **[M]** Wire outbound shipping to order fulfillment — `order-api` calls `shipping-api` `POST /shipping/labels/outbound` when order transitions to `shipped`; store tracking number and carrier on the order record
- [ ] **[M]** Guest + registered account identity model — email (not phone) is the cross-DB correlation key for both guest and registered accounts; each account still has a unique `account_id` as the primary key; email must be stored consistently across user-api, auth-api, and promotions-api to enable unsubscribe preference management for guest accounts without requiring registration; do not use phone as a linking key (privacy risk)
