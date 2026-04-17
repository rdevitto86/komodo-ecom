# Komodo Platform — TODO

Priority guide: **[H]** = blocking UI simulation · **[M]** = important, not blocking · **[L]** = low priority (docs, testing)

APIs are ordered by how soon the UI needs them to simulate a real backend.

---

## komodo-auth-api
> Status: MVP complete. Token issuing and validation work. Auth code flow and revocation storage missing.

- [ ] **[M]** Persist revoked tokens to ElastiCache (TTL = token expiry) so `/oauth/revoke` is actually effective
- [ ] **[M]** Store token JTI in ElastiCache on issue; check on introspect to detect reuse after revocation
- [ ] **[M]** Check if refresh token is revoked in ElastiCache before issuing a new access token (`oauth_token_handler.go:173`)
- [ ] **[L]** Implement `authorization_code` grant flow (requires SvelteKit login UI to be live first)
- [ ] **[L]** Add unit tests for token signing, validation, and introspection

---

## komodo-user-api
> Status: Handlers and service layer implemented. Sub-item CRUD (addresses, payments, preferences) blocked on DynamoDB schema.

- [ ] **[H]** Finalize DynamoDB single-table key schema in `docs/data-model.md` — unblocks all sub-item operations
- [ ] **[H]** Wire `repo.CreateAddress` / `UpdateAddress` / `DeleteAddress` once schema is finalized
- [ ] **[H]** Wire `repo.UpsertPayment` / `DeletePayment` once schema is finalized
- [ ] **[H]** Wire `repo.UpdatePreferences` / `DeletePreferences` once schema is finalized
- [ ] **[M]** Verify internal server ownership checks — confirm `resolveUserID` correctly rejects cross-user access on internal routes
- [ ] **[L]** Add integration tests for public + internal handler paths

---

## komodo-shop-items-api
> Status: Fully implemented. S3-backed product catalog works. Suggestions endpoint is a stub.

- [ ] **[M]** Replace stub recommendation logic in `GET /suggestions` with real logic (rule-based, ML, or simple bestsellers query)
- [ ] **[M]** Evaluate recommendation automation — assess whether user browsing/purchase history can drive `GET /suggestions` (rule-based first, ML later)
- [ ] **[M]** Add `service_type` field to `ShopItem` model (`product | service | repair`) — repair items carry additional fields: `accepted_device_types`, `estimated_turnaround_days`, `warranty_on_repair`; update S3 schema and `openapi.yaml`
- [ ] **[M]** Add `GET /services/repair` route — filter shop items by `service_type=repair`; return paginated repair service listings
- [ ] **[M]** Add `GET /services/repair/{id}` route — single repair service detail (accepted devices, pricing, turnaround, warranty)
- [ ] **[L]** Add unit tests for S3 fetch and item parsing

---

## komodo-cart-api
> Status: Fully implemented. Guest + authenticated carts, checkout token generation, and stock hold coordination all work.

- [ ] **[M]** Design and implement "Save for Later" feature (separate DynamoDB entity, no TTL) — see README TODO
- [ ] **[L]** Add integration tests for guest cart TTL, merge flow, and checkout token lifecycle

---

## komodo-shop-inventory-api
> Status: Stub complete — all layers scaffolded (Axum, models, repo trait, handlers). DynamoDB impl is `todo!()`.

- [ ] **[H]** Implement `DynamoInventoryRepo::reserve` — conditional write (`available_qty >= requested`), write HOLD# record with TTL
- [ ] **[H]** Implement `DynamoInventoryRepo::get_stock` + `batch_stock` — GetItem / BatchGetItem for SKU#/STOCK records
- [ ] **[H]** Implement `DynamoInventoryRepo::release_hold` — DeleteItem HOLD# record, restore `available_qty`
- [ ] **[H]** Implement `DynamoInventoryRepo::confirm` — DeleteItem HOLD# record, decrement `reserved_qty`, increment `committed_qty`
- [ ] **[H]** Implement `DynamoInventoryRepo::restock` — UpdateItem `available_qty += qty`
- [ ] **[H]** Wire Secrets Manager bootstrap — populate `Config` secret fields at startup
- [ ] **[H]** Implement JWT RS256 validation in `middleware/auth.rs` — `DecodingKey::from_rsa_pem` + `jsonwebtoken::decode`
- [ ] **[M]** DynamoDB Streams handler (separate Lambda) — listen for TTL expiry events, restore `available_qty` on hold expiry
- [ ] **[M]** Wire `communications_api_url` to fire restock threshold alert when `available_qty` drops below `restock_threshold`
- [ ] **[L]** Implement `common::spawn_app()` in tests + enable integration tests

---

## komodo-order-api
> Status: Complete stub — `cmd/public/main.go` is an empty package declaration. Core purchase flow depends on this.

- [ ] **[H]** Scaffold service: bootstrap (logger, secrets, DynamoDB, Redis), middleware stacks, dual-server (public + internal) ServeMux
- [ ] **[H]** Implement `POST /me/orders` — consume checkout token from cart-api, confirm holds, write order to DynamoDB, publish `order.placed` event
- [ ] **[H]** Implement `GET /me/orders` — list authenticated user's orders (paginated)
- [ ] **[H]** Implement `GET /me/orders/{orderId}` — get order detail
- [ ] **[H]** Implement `POST /me/orders/{orderId}/cancel` — cancel order, release holds, trigger refund via payments-api
- [ ] **[H]** Implement internal `GET /internal/orders/{orderId}` — for returns-api and payments-api lookups
- [ ] **[M]** Implement order status state machine: `pending → confirmed → shipped → delivered → cancelled`
- [ ] **[M]** Publish `order.cancelled` and `order.fulfilled` events to event-bus-api
- [ ] **[M]** Add idempotency on `POST /me/orders` to prevent double-order on retry
- [ ] **[L]** Add integration tests for order creation and cancellation flow

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
- [ ] **[L]** Implement `common::spawn_app()` in tests + enable integration tests with Stripe test mode
- [ ] **[M]** Enforce autopay requires payment method on file — validate that the user has an active bank account or credit card record before enabling autopay or processing any autopay transaction; return a clear error if no method exists; may require an `autopay_enabled` boolean and method-presence check in the DB schema

---

## komodo-address-api
> Status: Implemented but provider calls are stubs. Stateless — no DB needed.

- [ ] **[H]** Wire real address validation provider (SmartyStreets, Google Address Validation, or similar) — replace all 3 stub `TODO` bodies in `internal/provider/address.go`
- [ ] **[M]** Add provider API key secret (`ADDRESS_PROVIDER_API_KEY`) to LocalStack init seed
- [ ] **[L]** Add unit tests for provider error handling and response mapping

---

## komodo-search-api
> Status: Partial — middleware and routes wired, but all handlers return empty results. Typesense not initialized.

- [ ] **[H]** Initialize Typesense client after secrets load (`TODO(typesense)` in main.go)
- [ ] **[H]** Implement `GET /search` — build query params from request, call Typesense, return results
- [ ] **[H]** Implement `POST /internal/index/sync` — full re-index from shop-items-api S3 data into Typesense
- [ ] **[M]** Wire event-bus-api subscriber to listen for `shop_item.created/updated/deleted` → incremental index updates
- [ ] **[M]** Implement `DELETE /internal/index` — drop and recreate Typesense collection (used for schema migrations)
- [ ] **[M]** Migrate off `gorilla/mux` to `net/http` ServeMux (platform convention)
- [ ] **[L]** Add integration tests for search query building and index sync

---

## komodo-support-api
> Status: Fully implemented with Anthropic Haiku integration. In-memory storage not production-safe.

- [ ] **[H]** Replace in-memory repository with DynamoDB — design table schema and implement all repo functions (`repository/chat.go:24`)
- [ ] **[M]** Define audit event schema and destination before wiring deletion audit — DynamoDB audit table or S3 archive (`repository/chat.go:18`)
- [ ] **[M]** Wire escalation (`POST /chat/escalate`) to communications-api for async ticket creation
- [ ] **[M]** Replace SQS publish placeholder in escalation handler with real SQS client once forge SDK has SQS support (`handlers/chat.go:182`)
- [ ] **[M]** Emit audit event to event-bus-api before chat history deletion (compliance trail)
- [ ] **[M]** Complete anonymous → authenticated session merge flow (currently not wired end-to-end)
- [ ] **[L]** Design human agent handoff flow (currently no handoff target exists)
- [ ] **[L]** Add integration tests for chat session lifecycle

---

## komodo-event-bus-api
> Status: Functional for local dev (in-memory fan-out). SNS/SQS wired but not deployed.

- [ ] **[M]** Enforce event type allowlist validation on `POST /events` (catalogue is defined but not enforced)
- [ ] **[M]** Deploy and test SNS/SQS fan-out path (in-memory fan-out is local-only)
- [ ] **[M]** Wire and deploy CDC Lambda handler for DynamoDB Streams
- [ ] **[M]** Add CDC event classifiers for payments, users, inventory, and cart domains — only orders classifier exists (`cdc/domains/orders.go:19`)
- [ ] **[M]** Expand CDC order event payload with additional fields: `total_cents`, `item_count`, `customer_id` (`cdc/domains/orders.go:48`)
- [ ] **[M]** Wire EventBridge as the routing layer for CDC events — add EventBridge rules for flexible per-domain fanout (orders → order consumers, payments → payment consumers, etc.)
- [ ] **[M]** Define and publish interaction events (`cart.item_added`, `order.started`, `order.abandoned`) — extend event type catalogue for analytics consumers
- [ ] **[M]** Per-connector publisher workers + DLQ — each outbound sink (SNS, EventBridge, S3 data lake, in-memory fan-out) runs its own goroutine/worker so one slow or failing sink cannot block the others; failed publishes route to a per-connector DLQ (SQS or in-memory bounded queue) with retry policy and max-age eviction
- [ ] **[L]** Emit CloudWatch metric (or fixed-key structured log) on unroutable CDC events (`cdc/handler.go:51`)
- [ ] **[L]** Add integration tests for event publish and consumer routing
- [ ] **[L]** Evaluate gRPC as an additional internal transport — research protobuf schema enforcement, bi-directional streaming, and performance vs complexity tradeoff compared to current HTTP REST; document findings in `docs/design-decisions.md` before any implementation; not a blocking item

---

## komodo-order-reservations-api
> Status: Routes wired, middleware configured, but all repository functions are stubs (15+ TODOs).

- [ ] **[M]** Initialize DynamoDB client in bootstrap (blocked on forge SDK `aws/dynamodb` availability — confirm package path)
- [ ] **[M]** Implement `repo.GetBooking` / `CreateBooking` / `UpdateBooking` (DynamoDB)
- [ ] **[M]** Implement `repo.GetSlots` / `UpdateSlotAvailability` (DynamoDB)
- [ ] **[M]** Extract `customer_id` from JWT context in handlers (currently a TODO)
- [ ] **[M]** Add ownership/authorization check — reject cross-customer booking reads/mutations
- [ ] **[M]** Add `POST /internal/slots/sync` route for slot inventory management
- [ ] **[M]** Decide and implement checkout hold flow (Option A: hold at reservation time vs Option B: hold at order confirm)
- [ ] **[M]** Extend booking model for repair intake — add `repair` booking type with fields: `device_type`, `serial_number`, `reported_issue`, `inbound_shipment_id`; wire to inbound shipping flow once shipping-api exists
- [ ] **[M]** Implement repair status state machine: `intake_pending → received → diagnosing → repairing → quality_check → ready → shipped_back`; emit status change events to event-bus-api on each transition
- [ ] **[L]** Add integration tests for booking lifecycle

---

## komodo-shipping-api (NEW)
> Status: Not yet created. Handles both inbound (customer → warehouse: returns, repair intake) and outbound (warehouse → customer: order fulfillment, repaired items) shipment label generation and tracking.

- [ ] **[M]** Scaffold service: bootstrap (logger, secrets, DynamoDB), middleware stack, ServeMux routes (port TBD — reserve in port allocation table)
- [ ] **[M]** Select and integrate a carrier aggregator (EasyPost, ShipStation, or EasyPost-compatible) — abstract behind a provider interface so carriers are swappable
- [ ] **[M]** Implement `POST /shipping/labels/outbound` — generate outbound label for order fulfillment; called by order-api when order transitions to `shipped`; return carrier, tracking number, and label URL
- [ ] **[M]** Implement `POST /shipping/labels/inbound` — generate prepaid inbound return/repair label; called by order-returns-api and order-reservations-api; customer receives label URL to print or QR scan
- [ ] **[M]** Implement `GET /shipping/{shipmentId}` — real-time shipment status; poll carrier API or return latest cached status
- [ ] **[M]** Add carrier webhook handler — receive status events from carrier (`delivered`, `out_for_delivery`, `exception`, `in_transit`); update shipment record and publish `shipment.status_updated` event to event-bus-api
- [ ] **[M]** Publish `shipment.label.created`, `shipment.delivered`, `shipment.received.inbound` events — `shipment.received.inbound` triggers inspection/repair workflow in reservations-api; `shipment.delivered` triggers loyalty points and fulfillment confirmation in order-api
- [ ] **[L]** Add integration tests for label generation, status polling, and webhook handling

---

## komodo-order-returns-api
> Status: Complete stub — `main.go` is a 27-line comment block describing what to implement.

- [ ] **[M]** Scaffold service: bootstrap, middleware stack, ServeMux routes
- [ ] **[M]** Implement RMA creation (`POST /me/returns`) — validate order eligibility, create return record
- [ ] **[M]** Implement RMA status tracking (`GET /me/returns`, `GET /me/returns/{returnId}`)
- [ ] **[M]** Coordinate refund via payments-api on return approval
- [ ] **[M]** Coordinate inventory restock via shop-inventory-api on return receipt
- [ ] **[L]** Wire points reversal via loyalty-api on refund
- [ ] **[L]** Trigger customer notification via communications-api on status change
- [ ] **[L]** Add integration tests for RMA lifecycle

---

## komodo-communications-api
> Status: Complete stub — directory exists but `main.go` is empty.

- [ ] **[L]** Scaffold service: bootstrap, middleware stack, ServeMux routes
- [ ] **[L]** Implement `POST /send/email` — transactional email via provider (SendGrid, SES, etc.)
- [ ] **[L]** Implement `POST /send/sms` — SMS via provider (Twilio, SNS, etc.)
- [ ] **[L]** Implement `POST /send/push` — in-app push notification
- [ ] **[L]** Subscribe to event-bus-api for async trigger events (`order.placed`, `order.shipped`, etc.)
- [ ] **[L]** Store and load transactional email templates from S3 — support per-locale variants; templates managed separately from code
- [ ] **[L]** Add template management for transactional messages

---

## komodo-loyalty-api
> Status: Complete stub — `main.go` is an empty `func main() {}`.

- [ ] **[L]** Scaffold service: bootstrap, middleware stack, ServeMux routes
- [ ] **[L]** Implement points earn on `order.placed` event (subscribe via event-bus-api)
- [ ] **[L]** Implement points reversal on `order.returned` event
- [ ] **[L]** Implement `GET /me/loyalty` — points balance and tier status
- [ ] **[L]** Implement `POST /me/loyalty/redeem` — apply points discount to a cart/order
- [ ] **[L]** Design tier rules and rewards catalogue

---

## komodo-reviews-api
> Status: Complete stub — `main.go` is a bare package declaration.

- [ ] **[L]** Scaffold service: bootstrap, middleware stack, ServeMux routes
- [ ] **[L]** Implement `POST /me/reviews` — submit a review (require verified purchase check via order-api)
- [ ] **[L]** Implement `GET /items/{itemId}/reviews` — paginated review listing
- [ ] **[L]** Implement `PUT /me/reviews/{reviewId}` / `DELETE /me/reviews/{reviewId}`
- [ ] **[L]** Implement rating aggregation (avg rating + count) — maintain in DynamoDB alongside reviews
- [ ] **[L]** Add moderation queue for flagged reviews

---

## komodo-entitlements-api
> Status: Complete stub — only `go.mod` and empty `main.go` exist.

- [ ] **[L]** Define entitlement model (what is being gated — plans, features, access levels)
- [ ] **[L]** Scaffold service: bootstrap, middleware stack, ServeMux routes
- [ ] **[L]** Implement `GET /me/entitlements` — return active entitlements for JWT subject
- [ ] **[L]** Implement `POST /internal/entitlements` — grant entitlement (called by order-api on purchase)
- [ ] **[L]** Implement `DELETE /internal/entitlements/{id}` — revoke entitlement

---

## komodo-features-api
> Status: Complete stub — only `go.mod`, empty `main.go`, and `openapi.yaml` exist.

- [ ] **[L]** Scaffold service: bootstrap, middleware stack, ServeMux routes
- [ ] **[L]** Implement `GET /features/{key}` — evaluate feature flag for a given context (user, env, percent rollout)
- [ ] **[L]** Implement `GET /me/features` — bulk flag evaluation for authenticated user
- [ ] **[L]** Implement internal CRUD for flag management (`POST/PUT/DELETE /internal/features/{key}`)
- [ ] **[L]** Back flags with DynamoDB; cache evaluated results in Redis with short TTL

---

## komodo-slm-api (NEW — Python FastAPI)
> Status: Not yet created. Python FastAPI service for local SLM inference and Bedrock offloading. Python chosen for native LLM tooling, sanitization, and injection protection — Go and Rust lack sufficient LLM ecosystem support for this layer.

- [ ] **[L]** Scaffold Python FastAPI service — project structure, dependency management (uv or pip), health endpoint, logging, secrets loading; port TBD in 7111–7120 Analytics & Discovery block (reserve 7113)
- [ ] **[L]** Implement provider abstraction — `LocalSLMProvider` (Ollama or similar on-prem) and `BedrockProvider` (AWS Bedrock); toggle via env config (`SLM_PROVIDER=local|bedrock`); both implement the same interface so callers are provider-agnostic
- [ ] **[L]** Implement task endpoints: `POST /summarize`, `POST /insights`, `POST /moderate` (guardrails/content moderation), `POST /support` (customer servicing suggestions); each endpoint validates input, sanitizes against prompt injection, and routes to the active provider
- [ ] **[L]** Wire input sanitization and prompt injection protection on all endpoints — normalize whitespace, strip control characters, enforce max token limits before forwarding to the model
- [ ] **[L]** Add Bedrock IAM role and secret config to infra — `AWS_BEDROCK_MODEL_ID`, region, and credentials; add local stub values to LocalStack secrets init
- [ ] **[L]** Add integration tests for each task endpoint against a local SLM stub

---

## komodo-statistics-api (NEW)
> Status: Not yet created. Real-time ecom stats service — subscribes to event-bus-api and maintains aggregated stats in an in-memory SQLite DB.

- [ ] **[M]** Scaffold service: bootstrap (logger, secrets, SQLite in-memory DB), dual-server (public + private) ServeMux, event-bus-api subscriber client; port 7111 (public), 7112 (private) — Analytics & Discovery block
- [ ] **[M]** Subscribe to relevant event-bus-api events — `cart.item_added`, `cart.item_removed`, `order.placed`, `order.fulfilled`, `shop_item.viewed`; update stat counters in SQLite on each event
- [ ] **[M]** Implement public stat routes — `GET /stats/items/{itemId}/in-cart` ("X users have this in cart"), `GET /stats/items/{itemId}/recently-bought` ("Y people bought this in the last month"), `GET /stats/items/{itemId}/frequently-bought-with` (co-purchase pairings); served to ecom UI via BFF
- [ ] **[M]** Implement internal/admin stat routes — `GET /internal/stats/dashboard` and per-domain aggregates for admin dashboards and inter-API analytics consumers
- [ ] **[M]** Background TTL cleanup worker — goroutine that periodically scans the SQLite DB for rows past their TTL column value and deletes them; prevents unbounded memory growth
- [ ] **[L]** Add integration tests for event subscription, stat accumulation, and TTL cleanup

---

## komodo-insights-api
> Status: Stub — handler skeletons only. Backed by Bedrock/Claude API calls.

- [ ] **[M]** Add local response cache for Bedrock/Claude API calls — cache by summary/query type with an expiry timestamp; serve cached result if not expired; refresh on weekly/biweekly schedule or when threshold exceeded; prevents repeated expensive LLM calls on every request

---

## komodo-shop-promotions-api
> Status: Not yet scaffolded. Handles promotions, discount logic, and first-order tracking.

- [ ] **[M]** Track first-order flag per account — on `order.placed` event, record whether this is the account's first order (guest or registered); store in promotions DB as a boolean/timestamp record; used to gate first-order discount eligibility

---

## Cross-Cutting

- [ ] **[H]** Finalize `user-api` DynamoDB single-table schema — unblocks addresses, payments, and preferences across the UI
- [ ] **[M]** Establish shared event type catalogue (`event-bus-api` has it defined but services aren't validating against it)
- [ ] **[M]** Wire `support-api` escalation → `communications-api` once communications-api is scaffolded
- [ ] **[M]** Typed config constants across Go services — replace sprinkled `os.Getenv("FOO")` calls with a per-service `internal/config` package (or shared SDK helper) exposing typed constants; centralizes the env-var surface, surfaces required-vs-optional vars, prevents typos
- [ ] **[L]** Add `docs/data-model.md` to every API that uses DynamoDB (currently missing for most)
- [ ] **[L]** Standardize `openapi.yaml` across all APIs (several stubs are missing or outdated)
- [ ] **[L]** Write unit + integration tests for all services (`go test ./...` must pass; at minimum happy path + error cases per handler)
- [ ] **[M]** Wire inbound shipping to returns flow — `order-returns-api` calls `shipping-api` `POST /shipping/labels/inbound` on RMA approval and returns the label URL to the customer
- [ ] **[M]** Wire inbound shipping to repair intake — `order-reservations-api` calls `shipping-api` `POST /shipping/labels/inbound` on repair booking confirmation; store `inbound_shipment_id` on the booking record
- [ ] **[M]** Wire outbound shipping to order fulfillment — `order-api` calls `shipping-api` `POST /shipping/labels/outbound` when order transitions to `shipped`; store tracking number and carrier on the order record
- [ ] **[M]** Guest + registered account identity model — email (not phone) is the cross-DB correlation key for both guest and registered accounts; each account still has a unique `account_id` as the primary key; email must be stored consistently across user-api, auth-api, and promotions-api to enable unsubscribe preference management for guest accounts without requiring registration; do not use phone as a linking key (privacy risk)
- [ ] **[M]** Reserve ports 7111–7113 in port allocation — 7111 `statistics-api` public, 7112 `statistics-api` private, 7113 `slm-api`; update CLAUDE.md port table when confirmed

> **SDK extractions moved** — items previously listed here (HTTP client base, health handler, circuit breaker with call-site context, shared hashing utility) now live in `komodo-forge-sdk-go/TODO.md`. Service-side code migration to the generated adapters is tracked there alongside the codegen pipeline.
