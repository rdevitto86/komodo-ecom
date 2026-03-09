# UI Architecture

## Service Topology

```
Browser
  │
  ▼
ui (SvelteKit 5 — :7001)
  │
  ├─► auth-api (:7011)         — login, OAuth, token refresh
  ├─► user-api (:7051)         — profile, addresses
  ├─► shop-items-api (:7041)   — product catalog
  ├─► order-api (:7061)        — order history, checkout
  └─► ssr-engine-svelte (:7003) — HTML fragment pre-rendering (not active yet)
```

## Two-Layer SSR Design

The UI has a two-layer rendering architecture (second layer not yet active):

1. **SvelteKit shell** (`ui/`, port 7001) — The primary app. Handles routing, auth, user interaction, client-side state. Currently runs as a static SPA (`adapter-static`) for demo deployment.

2. **SSR fragment engine** (`apis/komodo-ssr-engine-svelte`, port 7003) — Bun-based service that pre-renders and caches SvelteKit component subtrees. Delivers HTML fragments to the shell for performance-critical and SEO-sensitive sections (product detail pages, catalog grids). Not wired until real backend is live.

## Adapter Strategy

| Mode | Adapter | Trigger | Target |
|------|---------|---------|--------|
| Demo | `adapter-static` | `bun run build:demo` | S3 + CloudFront |
| Dev | dev server (no adapter) | `bun run dev` | localhost:7001 |
| Live | `adapter-node` (comment swap in svelte.config.js) | `bun run build` | EC2 / ECS |

## Auth Flow

1. Unauthenticated requests to protected routes are caught in `+layout.server.ts`
2. Server redirects to `/auth` with return URL
3. OAuth exchange happens via `komodo-auth-api` (port 7011)
4. JWT (RS256) is stored in an httpOnly cookie
5. All protected `+page.server.ts` loaders validate the JWT server-side via forge SDK middleware

## Docker (local dev)

The root `docker-compose.yml` manages the UI container under the `ui-backend` and `full` profiles:
- Build context: `./ui` with `target: dev`
- Hot-reload via volume mount (`./ui:/app`)
- Depends on `auth-api-public` health check before starting

## Dependencies on Monorepo Services

| Service | Port | Used For |
|---------|------|----------|
| `komodo-auth-api` | 7011 | Auth tokens, session management |
| `komodo-user-api` | 7051 | Profile data, preferences |
| `komodo-shop-items-api` | 7041 | Product catalog, inventory |
| `komodo-order-api` | 7061 | Order history, checkout |
| `komodo-ssr-engine-svelte` | 7003 | HTML fragment pre-rendering (future) |
