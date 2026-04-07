# Komodo UI — TODO

Priority guide: **[H]** = blocking real backend simulation · **[M]** = important, not blocking · **[L]** = low priority (docs, testing)

Sections are ordered by dependency — auth must come before cart, cart before checkout, etc.

---

## Auth Flow
> All auth endpoints return 501. No JWT handling exists. Protected routes are unprotected.

- [ ] **[H]** Implement `POST /api/user/login` BFF route — exchange credentials with `komodo-auth-api`; store JWT in `httpOnly; Secure; SameSite=Strict` cookie
- [ ] **[H]** Implement `POST /api/user/logout` BFF route — call `komodo-auth-api` `/oauth/revoke`, clear session cookie
- [ ] **[H]** Implement `POST /api/user/register` BFF route — create user via `komodo-user-api` `POST /me/profile`
- [ ] **[H]** Add JWT validation to root `+layout.server.ts` — validate cookie on every server load, hydrate session into `event.locals`
- [ ] **[H]** Build `/login` page — login form wired to `/api/user/login`; redirect to previous route on success
- [ ] **[H]** Build `/signup` page — registration form wired to `/api/user/register`
- [ ] **[H]** Add `+layout.server.ts` guards to `/profile`, `/profile/settings`, and `/checkout` that redirect unauthenticated users to `/login`
- [ ] **[H]** Build `src/lib/server/auth.ts` — `AuthClient` wrapping `komodo-auth-api` (token exchange, revoke, introspect)
- [ ] **[M]** Wire `UserState` to server-loaded session — currently just sets a boolean flag with no data
- [ ] **[M]** Implement `GET /api/user/preferences` + `PUT /api/user/preferences` BFF routes — proxy to `komodo-user-api`
- [ ] **[M]** Build combined signup/login page — single email entry; check email existence via BFF, then render login form (password) or signup form (full registration) in-place; no separate `/login` and `/signup` routes needed

---

## Product Pages
> Server load works (shop-items-api wired). The page itself is empty.

- [ ] **[H]** Build `/shop/[id]/+page.svelte` — render product images, name, price, description, variant selector, and add-to-cart using data already returned by `+page.server.ts`
- [ ] **[H]** Implement `ShopItem.svelte` — reusable product card for listings, search results, and homepage
- [ ] **[H]** Implement `PricingLabel.svelte` — price with sale/original price and currency formatting
- [ ] **[H]** Implement `QuantitySelector.svelte` — increment/decrement with min=1, max=stock
- [ ] **[H]** Implement `AddToCart.svelte` — calls cart BFF and updates `CartState`; handles out-of-stock state
- [ ] **[M]** Build `/shop/service/[id]` page — same pattern as product; wire `GET /api/shop/service/[id]` to `komodo-shop-items-api`
- [ ] **[M]** Fix `Header.svelte` nav links — currently point to `/products` and `/services` (routes that don't exist)
- [ ] **[M]** Fix cart badge in `Header.svelte` — hardcoded to "0"; bind to `CartState.itemCount`
- [ ] **[M]** Display certification badges on product/service detail pages — render from `ShopItem` data
- [ ] **[M]** Display top 6 key features per product/service — structured feature highlight section on detail page
- [ ] **[M]** Display recommended products/tools — wire to `komodo-shop-items-api` `GET /suggestions` on detail + cart pages

---

## Home Page
> Renders a Hero with hardcoded Unsplash images. No real content.

- [ ] **[H]** Build featured products section — call `komodo-shop-items-api` or use curated static list to show real products using `ShopItem.svelte`
- [ ] **[M]** Replace hardcoded Unsplash images with real product/brand imagery
- [ ] **[M]** Build `Footer.svelte` — navigation links, legal, social
- [ ] **[M]** Build scrolling promotional banner component — configurable slides with CTA links, auto-play with pause-on-hover

---

## Cart
> All cart BFF routes return 501. `CartState` has structure but no methods.

- [ ] **[H]** Build `src/lib/server/cart.ts` — `CartClient` wrapping all `komodo-cart-api` guest + auth endpoints
- [ ] **[H]** Implement `GET /api/shop/cart` BFF route — fetch cart from `komodo-cart-api` (`GET /me/cart` for auth, `GET /cart/{cartId}` for guest)
- [ ] **[H]** Implement `POST /api/shop/cart` BFF route — add item; call `komodo-cart-api` `POST /me/cart/items`
- [ ] **[H]** Implement `PUT /api/shop/cart` BFF route — update quantity; call `komodo-cart-api` `PUT /me/cart/items/{itemId}`
- [ ] **[H]** Implement `DELETE /api/shop/cart` BFF route — remove item or clear; call `komodo-cart-api` `DELETE /me/cart/items/{itemId}`
- [ ] **[H]** Add `addItem()`, `updateQuantity()`, `removeItem()`, `clear()`, and `subtotal` derived to `CartState`
- [ ] **[H]** Persist guest cart ID in `localStorage`; pass `X-Session-ID` header on guest cart requests
- [ ] **[H]** Build `/shop/cart` page — line items, quantity controls, subtotal, proceed-to-checkout CTA
- [ ] **[M]** Call `POST /me/cart/merge` after login with stored `guest_cart_id` to merge guest → authenticated cart
- [ ] **[M]** Add optimistic UI for add/remove — update `CartState` immediately, reconcile with server response

---

## Checkout Flow
> All 3 checkout pages are empty. Checkout BFF returns 501.

- [ ] **[H]** Build `/checkout` page — 3-step form: (1) shipping address, (2) payment method, (3) review & confirm
- [ ] **[H]** Implement `POST /api/checkout` BFF route — call `komodo-cart-api` `POST /me/cart/checkout` for token, then `komodo-order-api` `POST /me/orders` to place the order
- [ ] **[H]** Build `/checkout/complete` page — order confirmation with order ID, summary, and estimated delivery
- [ ] **[H]** Implement `CheckoutForm.svelte` — multi-step form component with address, payment, and confirmation steps
- [ ] **[H]** Build `src/lib/server/checkout.ts` — orchestrates cart-api checkout token + order-api order placement
- [ ] **[M]** Pre-fill checkout address from user's saved addresses (`komodo-user-api` `GET /me/addresses`)
- [ ] **[M]** Pre-fill payment from user's saved methods (`komodo-user-api` `GET /me/payments`)
- [ ] **[M]** Wire `POST /api/payments/intent` BFF route — create payment intent via `komodo-payments-api`
- [ ] **[M]** Implement `PaymentProcessors.svelte` — credit card form (Stripe Elements or equivalent)

---

## User Profile
> Profile pages are empty stubs. All user BFF routes return 501.

- [ ] **[H]** Build `src/lib/server/user.ts` — `UserClient` wrapping `komodo-user-api` (profile, addresses, payments, preferences)
- [ ] **[H]** Implement `GET/PUT/DELETE /api/user/profile` BFF routes — proxy to `komodo-user-api`
- [ ] **[H]** Build `/profile` page — display name, email, order history summary, account actions
- [ ] **[M]** Build `/profile/settings` page — editable profile form, address book, saved payment methods, preferences toggles
- [ ] **[M]** Implement `GET/POST/PUT/DELETE /api/address` BFF routes — confirm correct target (`komodo-user-api` for saved addresses; `komodo-address-api` for validation only)
- [ ] **[M]** Implement payment method CRUD in profile settings — wire to `komodo-user-api` `GET/PUT/DELETE /me/payments`
- [ ] **[M]** Implement `ProfileSettings.svelte`, `PaymentMethods.svelte` components

---

## Orders
> Order BFF routes all return 501. No order pages exist.

- [ ] **[M]** Build `src/lib/server/order.ts` — `OrderClient` wrapping `komodo-order-api`
- [ ] **[M]** Implement `GET /api/orders/history/[userID]` BFF route — proxy to `komodo-order-api` `GET /me/orders`
- [ ] **[M]** Implement `GET /api/orders/[id]/details` BFF route — proxy to `komodo-order-api` `GET /me/orders/{orderId}`
- [ ] **[M]** Implement `POST /api/orders/[id]/returns` BFF route — proxy to `komodo-order-returns-api`
- [ ] **[M]** Build order detail page — items, status timeline, tracking info, return/cancel CTA
- [ ] **[M]** Implement `OrderHistory.svelte` — table of past orders with status badges

---

## Search
> `/shop/search` page is empty. Search BFF returns 501.

- [ ] **[M]** Build `src/lib/server/search.ts` — `SearchClient` wrapping `komodo-search-api`
- [ ] **[M]** Implement `GET /api/shop/search` BFF route — proxy to `komodo-search-api`
- [ ] **[M]** Build `/shop/search` page — results grid, filters sidebar, pagination
- [ ] **[M]** Implement `ItemFiltering.svelte` — category, price range, rating filter controls
- [ ] **[M]** Implement `Pagination.svelte` — page controls for search results
- [ ] **[M]** Wire search input in `Header.svelte` to navigate to `/shop/search?q=...`
- [ ] **[M]** Implement `Search.svelte` — debounced typeahead with instant results dropdown

---

## Core UI Components
> 69 of 70 components are empty stubs. These block every page above.

- [ ] **[H]** `Input.svelte` — text input with label, validation error state, helper text
- [ ] **[H]** `Select/DropdownSelector.svelte` — controlled select with options list
- [ ] **[H]** `CheckboxGroup.svelte` / `RadioGroup.svelte` — form group controls
- [ ] **[H]** `Spinner.svelte` / `Skeleton.svelte` — loading states for async content
- [ ] **[H]** `Toast.svelte` / `Alert.svelte` — error and confirmation feedback; wire to `NotificationsState`
- [ ] **[M]** `Modal.svelte` — overlay dialog for confirmations and inline forms
- [ ] **[M]** `Badge.svelte` — status labels (order status, stock status)
- [ ] **[M]** `Tooltip.svelte` — hover hints
- [ ] **[M]** `Breadcrumbs.svelte` — navigation trail; wire to `AppState.breadcrumbs` (restoration is a TODO in AppState)
- [ ] **[M]** `Sidebar.svelte` — collapsible side panel; wire to `AppState.sidebarOpen`
- [ ] **[M]** `Hamburger.svelte` + `DropdownMenu.svelte` — mobile nav and account menu in Header
- [ ] **[M]** `DatePicker.svelte` — for service scheduling flows
- [ ] **[L]** `Tabs.svelte`, `Accordion.svelte` — secondary layout components

---

## State Management
> Stores are declared but have no data flow end-to-end.

- [ ] **[H]** `UserState` — hydrate `#profile` from root server load (`event.locals.user`); wire `login()`/`logout()` to real session data
- [ ] **[H]** `CartState` — implement `addItem()`, `updateQuantity()`, `removeItem()`, `clear()`, `itemCount` and `subtotal` derived; sync with cart BFF on page load
- [ ] **[M]** `NotificationsState` — implement `add()`, `dismiss()`, `clear()`; drive `Toast.svelte`
- [ ] **[M]** `item.svelte.ts` — add selected variant, quantity, and loaded product state for product detail page
- [ ] **[M]** Fix `AppState` breadcrumb restoration — currently a `// TODO` comment

---

## Payments BFF
> All payment routes return 501.

- [ ] **[M]** Build `src/lib/server/payments.ts` — `PaymentsClient` wrapping `komodo-payments-api`
- [ ] **[M]** Implement `POST /api/payments/intent` — create payment intent; proxy to `komodo-payments-api`
- [ ] **[M]** Implement `GET/POST/DELETE /api/payments/methods` — saved card management; proxy to `komodo-payments-api`
- [ ] **[M]** Implement `GET /api/payments/transactions` — transaction history; proxy to `komodo-payments-api`

---

## Security
> Auth not yet implemented — these become relevant once sessions exist.

- [ ] **[H]** Set `httpOnly; Secure; SameSite=Strict` on session cookie when auth is wired
- [ ] **[M]** Add Content Security Policy headers in `hooks.server.ts` — restrict script/style/img sources
- [ ] **[M]** Add rate limiting on BFF write routes (`/api/user/login`, `/api/checkout`, `/api/payments/*`)
- [ ] **[M]** Add `.env.example` documenting all required env vars (`AUTH_API_URL`, `USER_API_URL`, `CART_API_URL`, `SHOP_ITEMS_API_URL`, etc.)
- [ ] **[M]** Implement BFF-layer field encryption/decryption for sensitive PII — encrypt fields in the SvelteKit server layer (never in the browser) using keys fetched from AWS Secrets Manager before passing to backend APIs; decrypt on read
- [ ] **[L]** Add Subresource Integrity (SRI) hashes for any external CDN assets
- [ ] **[L]** Replace hardcoded Unsplash URLs with a proxied image route to avoid leaking behavior to third-party

---

## Adapter & Build
> Currently `adapter-static`. Must switch to `adapter-node` for real server-side BFF calls to work in production.

- [ ] **[H]** Switch `svelte.config.js` from `adapter-static` → `adapter-node` when backend services are live (comment swap is already prepared in the file)
- [ ] **[M]** Add production `Dockerfile` optimized for `adapter-node` (current one is dev-only hot-reload)
- [ ] **[M]** Set `ORIGIN` env var for `adapter-node` CSRF validation

---

## Testing
- [ ] **[L]** Write unit tests for all BFF API routes (happy path + error cases) — Vitest
- [ ] **[L]** Write component tests for all implemented UI components — `@testing-library/svelte`
- [ ] **[L]** Write integration tests for BFF → API service flows — Vitest with MSW or real LocalStack endpoints
- [ ] **[L]** Write E2E tests for core flows: login → add to cart → checkout → order confirmation — Playwright
- [ ] **[L]** Add client-side error boundary logging to `/api/log` route
- [ ] **[L]** Implement log rehydration from `localStorage` in `handler.worker.ts` — replay any queued log entries on worker init (`lib/logger/common/handler.worker.ts:11`)
- [ ] **[L]** Add Core Web Vitals tracking (LCP, FID, CLS)

---

## Analytics & Telemetry
> No client-side event tracking exists. Required for understanding user behaviour and diagnosing failures.

- [ ] **[M]** Implement clickstream tracking — capture UI element interactions (clicks, navigation, meta) and batch-send to analytics endpoint
- [ ] **[M]** Implement client-side telemetry — collect device type on page load; report failures (unhandled errors, API timeouts) to telemetry endpoint; nothing beyond device + failures
- [ ] **[M]** Wire interaction events to backend (cart add, order start, order abandoned) — fire to event-bus-api on user actions

---

## Support & Chat
> Support chat UI not yet implemented.

- [ ] **[M]** Build chat widget — send/receive messages against `komodo-support-api` `/chat/messages`; guest and authenticated modes
- [ ] **[M]** Display full chat session history — fetch prior messages on widget open; paginate if session is long
- [ ] **[L]** Wire escalation CTA in chat widget — call `/chat/escalate` and confirm ticket creation to user

---

## Documentation
- [ ] **[L]** Add `.env.example` with all required env vars and descriptions
- [ ] **[L]** Document component API (props, slots, events) for shared components once implemented
- [ ] **[L]** Update `docs/architecture.md` once BFF wiring is complete

## SDK Extractions (komodo-forge-sdk-ts)

- [ ] **[M]** Extract `APIClient` base class into SDK — `src/lib/server/common/client.ts` is a generic HTTP client with RFC 7807 `ProblemDetail` error handling. Currently only used by the UI but needed by any TypeScript service (SSR engine, future BFF workers). Move to `komodo-forge-sdk-ts` so it's available without copying.
