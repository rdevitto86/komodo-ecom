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
- [ ] **[L]** Add unit tests for S3 fetch and item parsing

---

## komodo-cart-api
> Status: Fully implemented. Guest + authenticated carts, checkout token generation, and stock hold coordination all work.

- [ ] **[M]** Design and implement "Save for Later" feature (separate DynamoDB entity, no TTL) — see README TODO
- [ ] **[L]** Add integration tests for guest cart TTL, merge flow, and checkout token lifecycle

---

## komodo-shop-inventory-api
> Status: Complete stub — `main.go` is a single TODO comment. Required for cart checkout holds.

- [ ] **[H]** Scaffold service: bootstrap (logger, secrets, DynamoDB), middleware stack, ServeMux routes
- [ ] **[H]** Implement `POST /stock/{sku}/reserve` — create hold record in DynamoDB with TTL (`HOLD_TTL_SEC`)
- [ ] **[H]** Implement `DELETE /stock/{sku}/holds/{holdId}` — release a hold
- [ ] **[H]** Implement `POST /stock/{sku}/decrement` — confirm hold → permanent stock decrement (called by order-api)
- [ ] **[H]** Implement `GET /stock/{sku}` — return current available quantity
- [ ] **[H]** Implement TTL-based auto-release: DynamoDB TTL attribute on hold records so holds expire without a cron
- [ ] **[M]** Implement `POST /stock/bulk` — bulk availability check for cart validation
- [ ] **[M]** Add internal route `GET /internal/stock/{sku}/holds` for ops/debug inspection
- [ ] **[L]** Add integration tests for hold/release/decrement lifecycle

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
> Status: Complete stub — `cmd/public/main.go` is an empty package declaration. Required for order checkout.

- [ ] **[H]** Scaffold service: bootstrap (logger, secrets, DynamoDB), middleware stacks, dual-server ServeMux
- [ ] **[H]** Implement `POST /payments/charge` — charge a saved payment method; integrate with Stripe (or stub with configurable test mode)
- [ ] **[H]** Implement `POST /payments/refund` — refund by order/charge ID; called by order-api on cancellation
- [ ] **[H]** Implement internal `GET /internal/payments/{orderId}` — payment status lookup for order-api and returns-api
- [ ] **[M]** Implement `POST /payments/methods` + `DELETE /payments/methods/{id}` for tokenized card storage (delegates to Stripe; stores token ref in DynamoDB)
- [ ] **[M]** Add webhook handler for async Stripe events (`payment_intent.succeeded`, `charge.refunded`, etc.)
- [ ] **[M]** Publish `payment.succeeded` / `payment.failed` / `payment.refunded` events to event-bus-api
- [ ] **[M]** Implement `GET /payments/card-identify` — BIN lookup to identify card network (Visa, Mastercard, Amex, etc.) for UI display before charge
- [ ] **[M]** Implement subscription billing flow — recurring charge schedule, webhook handling for renewal/failure/cancellation events
- [ ] **[L]** Add integration tests with Stripe test mode keys

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
- [ ] **[L]** Add integration tests for booking lifecycle

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

## SDK Extractions (komodo-forge-sdk-go)

- [ ] **[M]** Extract HTTP client base into `http/client` — `pkg/v1/client/client.go` in auth-api, user-api, and cart-api all implement the same `post()`/`get()` helpers (context, JSON marshal, bearer token, RFC 7807 error unwrap). Every service that calls another service will need this. Move to SDK so services only define their own endpoint methods.
- [ ] **[L]** Add health handler to `http/handlers` — 5+ services (`user-api`, `cart-api`, `shop-inventory-api`, `reviews-api`, `features-api`) implement an identical `{"status":"OK"}` health endpoint. Move to SDK as a one-liner registration.
