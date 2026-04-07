# UI System Design

## Route Structure

| Route | Type | Auth | Data Sources |
|-------|------|------|-------------|
| `/` | Page | No | Static + shop-items-api (featured) |
| `/shop` | Page | No | shop-items-api |
| `/shop/[slug]` | Page | No | shop-items-api |
| `/feed` | Page | No | Static content |
| `/about` | Page | No | Static |
| `/faq` | Page | No | Static |
| `/contact-us` | Page | No | — |
| `/profile` | Page | Yes | user-api |
| `/checkout` | Page | Yes | shop-items-api, order-api |
| `/auth` | Page | No | auth-api |
| `/api/*` | API routes | Varies | Proxied to backend services |

## Data Flow (Live Mode)

```
+page.server.ts (load fn)
  └─► forge SDK client (from @komodo-forge-sdk/typescript)
        └─► backend API (auth-api / user-api / shop-items-api / etc.)
              └─► DynamoDB / S3

+page.svelte
  └─► receives data from load fn (no direct API calls from client)
  └─► client interactions → API routes (/api/*) → backend
```

## Demo Mode (Mock)

When built with `--mode mock` (`.env.mock`), all backend calls are replaced by in-process mock data. No network requests leave the browser. This enables full demo deployment on S3/CloudFront with zero backend infrastructure.

## State Management

Global state lives in `src/lib/state/` using Svelte runes (`$state`, `$derived`):

| Store | File | Scope |
|-------|------|-------|
| App state (theme, nav) | `app.svelte.ts` | Global |
| Cart | `cart.svelte.ts` | Global |
| Current item | `item.svelte.ts` | Page-scoped |
| Notifications | `notifications.svelte.ts` | Global |
| User session | `user.svelte.ts` | Global |

## Component Architecture

Components live in `src/lib/components/` grouped by domain:

```
components/
├── animations/   — GSAP wrappers, motion primitives
├── commerce/     — Product cards, cart items, price display
├── data/         — Data tables, lists, pagination
├── display/      — Typography, layout, content sections
├── feedback/     — Toasts, alerts, error states, empty states
├── forms/        — Input, select, checkbox, validation wrappers
├── loading/      — Skeleton screens, spinners
├── media/        — Image, video with lazy loading
├── navigation/   — Navbar, breadcrumbs, tabs, sidebar
├── overlays/     — Modals, drawers, tooltips
├── primitives/   — Base unstyled building blocks (button, card, badge)
├── sections/     — Page-level composed sections
└── user/         — Avatar, profile card, auth gate
```

## Error Handling

- Server-side errors: caught in `+error.svelte`, return RFC 7807-style problem data
- Client-side errors: routed to notification store → toast/alert display
- Network failures in mock mode: never happen (all data is in-process)

## Testing Strategy

| Layer | Tool | Location |
|-------|------|----------|
| Unit (components, logic) | Vitest + Testing Library | `src/routes/__tests__/` |
| E2E (user flows) | Playwright | `e2e/` |
| Auth fixtures | Playwright fixtures | `e2e/fixtures/auth.json` |
