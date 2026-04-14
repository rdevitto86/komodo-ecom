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
> **ABANDONED** — V1 is Rust. See `komodo-shop-inventory-api-rust`.

~~Go implementation items removed — all work happens in the Rust service.~~

---

## komodo-shop-inventory-api-rust
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
> **ABANDONED** — V1 is Rust. See `komodo-payments-api-rust`.

~~Go implementation items removed — all work happens in the Rust service.~~

---

## komodo-payments-api-rust
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
- [ ] **[L]** Emit CloudWatch metric (or fixed-key structured log) on unroutable CDC events (`cdc/handler.go:51`)
- [ ] **[L]** Add integration tests for event publish and consumer routing

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

## Cross-Cutting

- [ ] **[H]** Finalize `user-api` DynamoDB single-table schema — unblocks addresses, payments, and preferences across the UI
- [ ] **[H]** Add shared hashing utility (bcrypt or Argon2id) to forge SDK or cross-service layer — standardize password and token hashing across all services
- [ ] **[M]** Establish shared event type catalogue (`event-bus-api` has it defined but services aren't validating against it)
- [ ] **[M]** Wire `support-api` escalation → `communications-api` once communications-api is scaffolded
- [ ] **[M]** Confirm forge SDK `aws/dynamodb` package path and update `order-reservations-api` bootstrap
- [ ] **[L]** Add `docs/data-model.md` to every API that uses DynamoDB (currently missing for most)
- [ ] **[L]** Standardize `openapi.yaml` across all APIs (several stubs are missing or outdated)
- [ ] **[L]** Write unit + integration tests for all services (`go test ./...` must pass; at minimum happy path + error cases per handler)
- [ ] **[M]** Wire inbound shipping to returns flow — `order-returns-api` calls `shipping-api` `POST /shipping/labels/inbound` on RMA approval and returns the label URL to the customer
- [ ] **[M]** Wire inbound shipping to repair intake — `order-reservations-api` calls `shipping-api` `POST /shipping/labels/inbound` on repair booking confirmation; store `inbound_shipment_id` on the booking record
- [ ] **[M]** Wire outbound shipping to order fulfillment — `order-api` calls `shipping-api` `POST /shipping/labels/outbound` when order transitions to `shipped`; store tracking number and carrier on the order record
- [ ] **[M]** Add default version exports to all Go service `pkg/` packages — each `pkg/` root should re-export from the current stable versioned subpackage (e.g. `pkg/v1`) so consumers can import a single unversioned canonical path; older/newer versions remain importable via their versioned subpath

## SDK Extractions (komodo-forge-sdk-go)

- [ ] **[M]** Extract HTTP client base into `http/client` — `pkg/v1/client/client.go` in auth-api, user-api, and cart-api all implement the same `post()`/`get()` helpers (context, JSON marshal, bearer token, RFC 7807 error unwrap). Every service that calls another service will need this. Move to SDK so services only define their own endpoint methods.
- [ ] **[L]** Add health handler to `http/handlers` — 5+ services (`user-api`, `cart-api`, `shop-inventory-api`, `reviews-api`, `features-api`) implement an identical `{"status":"OK"}` health endpoint. Move to SDK as a one-liner registration.
- [ ] **[M]** **Enterprise Pattern: Circuit Breaker** — Extract circuit breaker to `resilience/circuitbreaker` in SDK. APIs with external service dependencies need circuit breaker protection when DB or downstream services are down:
  - **komodo-auth-api**: ElastiCache token revocation checks (`oauth_token_handler.go:173`)
  - **komodo-cart-api**: `shop-items-api` calls (product snapshots), `inventory-api` calls (stock holds)
  - **komodo-support-api**: Anthropic Haiku LLM calls (`anthropic.go:39`)
  - **komodo-address-api**: External address validation provider (SmartyStreets/Google) — currently stubs (`address.go:50,59,77`)
  - **komodo-search-api**: Typesense search queries
  - **komodo-communications-api**: SendGrid/SES email, Twilio/SNS SMS (future providers)
  - **komodo-shipping-api**: Carrier aggregator API (EasyPost/ShipStation)
  - **komodo-payments-api-rust**: Stripe API calls (`payment_intents`, `refunds`)
  - **komodo-event-bus-api**: SNS publish calls (CDC Lambda and relay publisher)
  - **Cross-service calls**: cart-api ↔ inventory-api, order-api ↔ payments-api, order-api ↔ shipping-api, returns-api ↔ payments-api/inventory-api
  - **Pattern requirements**: Configurable failure threshold, half-open state probe, exponential backoff, fallback to degraded mode (cache-only, async queue, or fail-fast with clear error codes)
