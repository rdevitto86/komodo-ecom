# Business Rules

## Authentication & Authorization

- All checkout and profile routes require a valid JWT session.
- Sessions are validated server-side in `+layout.server.ts` — never trust client-side state for auth decisions.
- Unauthenticated users hitting protected routes are redirected to `/auth?next=<return_url>`.
- JWT expiry is handled by the auth-api token refresh flow — the UI should not manage token rotation directly.

## Cart

- Cart state is client-side only until checkout begins.
- Items in the cart are validated against live inventory at checkout initiation — cart state is not authoritative on stock.
- Guest cart is ephemeral (session storage or in-memory). Logged-in cart is persisted server-side via user-api.
- Merging guest cart into user cart on login is a defined use case.

## Product Catalog

- Product availability (in-stock / out-of-stock) is determined by shop-items-api at display time.
- Out-of-stock items may be displayed but cannot be added to cart.
- Prices displayed in the UI are informational — the authoritative price is confirmed at checkout by order-api.
- Product slugs are stable identifiers — do not construct URLs from item names or IDs alone.

## Checkout

- Checkout requires authentication. Guest checkout is not supported.
- The checkout flow must not allow double-submission (disable submit on first click, show loading state).
- Payment is handled by payments-api — the UI never processes raw card data.
- Order confirmation must be shown on success; the user must not be able to re-submit the same order via back navigation.

## Demo Mode

In `--mode mock`:
- No real API calls are made. All data is in-process mock data.
- Auth is simulated — mock sessions are always "logged in" with a test user.
- Cart actions work but do not persist.
- Checkout flow goes through all UI steps but never submits to order-api or payments-api.

## Data Validation

- All user-supplied input is validated client-side with `zod` before submission (for UX feedback).
- Server-side validation in API routes is always the source of truth — client validation is a convenience only.
- DOMPurify (`isomorphic-dompurify`) is used for any user-generated content rendered as HTML.
