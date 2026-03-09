# UI Service — Reference

## Port
| Service | Port |
|---------|------|
| UI (dev server) | 7001 |
| SSR Engine (fragments) | 7003 |

## Run Commands

```bash
# Development
bun run dev            # hot-reload dev server at :7001
bun run dev:mock       # dev server with mock API (no backend needed)

# Build
bun run build          # production build (adapter-static → build/)
bun run build:demo     # demo build (--mode mock, for S3/CloudFront)
bun run preview        # preview production build locally

# Via monorepo root (preferred for full-stack work)
make up-ui             # starts infra + all backend services + ui
make down              # stop everything
```

## Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `ENV` | No | `development` | Runtime environment label |
| `PORT` | No | `7001` | Exposed port (Docker only) |
| `LOG_LEVEL` | No | `debug` | Log verbosity |
| `VITE_AUTH_API_URL` | Demo: No, Live: Yes | — | Public auth API base URL |
| `VITE_USER_API_URL` | Demo: No, Live: Yes | — | Public user API base URL |
| `VITE_SHOP_ITEMS_API_URL` | Demo: No, Live: Yes | — | Shop items API base URL |

In demo mode (`--mode mock`), all `VITE_*_API_URL` vars are ignored — the app uses in-process mock data.

## Routes

| Route | Auth Required | Description |
|-------|-------------|-------------|
| `/` | No | Landing / home page |
| `/shop` | No | Product catalog |
| `/shop/[slug]` | No | Product detail page |
| `/feed` | No | Content/editorial feed |
| `/about` | No | About page |
| `/faq` | No | FAQ page |
| `/contact-us` | No | Contact form |
| `/profile` | Yes | User profile |
| `/checkout` | Yes | Checkout flow |
| `/auth` | No | Login / OAuth callback |

## Project Structure

```
src/
├── lib/
│   ├── components/     # Shared component library
│   │   ├── animations/ # GSAP / motion components
│   │   ├── commerce/   # Cart, product cards, etc.
│   │   ├── display/    # Layout, typography primitives
│   │   ├── feedback/   # Toasts, alerts, modals
│   │   ├── forms/      # Input, select, checkbox, etc.
│   │   ├── loading/    # Skeletons, spinners
│   │   ├── media/      # Image, video wrappers
│   │   ├── navigation/ # Nav bar, breadcrumbs, tabs
│   │   ├── overlays/   # Drawers, dialogs
│   │   ├── primitives/ # Base unstyled building blocks
│   │   ├── sections/   # Page-level section components
│   │   └── user/       # Auth / user-specific components
│   ├── server/         # Server-only code (never imported client-side)
│   ├── state/          # Svelte runes-based global state
│   ├── styles/         # Global CSS, design tokens
│   └── types/          # Shared TypeScript types
├── routes/             # SvelteKit file-based routing
└── app.html            # HTML shell
```

## Dependencies (notable)

| Package | Purpose |
|---------|---------|
| `@sveltejs/kit` | Framework |
| `tailwindcss` v4 | Styling (vite plugin) |
| `@threlte/core` + `extras` | Three.js integration |
| `gsap` | Animation |
| `lenis` | Smooth scroll |
| `zod` | Schema validation |
| `lucide-svelte` | Icons |
