# Komodo UI — AI Context

## Project Purpose
SvelteKit 5 frontend for a portfolio-grade e-commerce platform. Personal project with a realistic path to a real small business. Cost-efficient today with a clear AWS scaling path.

---

## Active Mode

Agent profiles live at the monorepo root in `.agents/`. See root `CLAUDE.md` for the full mode table.

| Mode | Trigger | Role |
|------|---------|------|
| **ADVISOR** (default) | No prefix | Senior peer — challenge, guide, never implement |
| **SENIOR** | `[SWE]` | Implements with judgment — brief flags, then executes |

---

## Context Strategy
**Do not pre-load context.** Discover only what's relevant to the current task.

**Working inside `ui/`:**
1. Read `ui/docs/README.md` first — port, run commands, env vars.
2. Pull other `ui/docs/` files only if directly relevant (e.g. `architecture.md` for topology, `ux-design.md` for component/design questions).
3. Check `src/lib/components/` before creating new components.
4. Fall back to root `CLAUDE.md` only for monorepo-wide conventions.

---

## Tech Stack
- **Framework:** SvelteKit 5 + TypeScript (bun runtime)
- **Adapter:** `adapter-static` — SPA mode with `fallback: 'index.html'` for S3/CloudFront demo. Switch to `adapter-node` when wiring real backend.
- **Styling:** Tailwind CSS v4 (vite plugin), `tw-animate-css`, `class-variance-authority`, `tailwind-merge`
- **3D / Animation:** Threlte (Three.js), GSAP, Lenis (smooth scroll), Postprocessing
- **Auth:** JWT RS256 via `komodo-auth-api` (7011). Validated server-side in `+layout.server.ts`.
- **SDK:** `komodo-forge-sdk-ts` (types, API clients)
- **Tests:** `bun run test:unit` (Vitest + Testing Library), `bun run test:e2e` (Playwright)

## Deployment Strategy
> Current state: demo site only. Real backend not deployed.

| Target | Adapter | Script | Deploy |
|--------|---------|--------|--------|
| Demo (S3 + CloudFront) | `adapter-static` | `bun run build:demo` | Static files → S3 |
| Local dev (Docker) | — (dev server) | `bun run dev` | `just up api ui` |
| Real backend (EC2/ECS) | `adapter-node` | `bun run build` | EC2 docker-compose |

`build:demo` uses `--mode mock` → reads `.env.mock` → disables real API calls and uses mock data.

## Local Dev
Run from monorepo root:
```
just up api ui   # starts: localstack + redis + auth-api + enabled APIs + ui
just down        # stop all
```
Or standalone (requires `just up` first for `komodo-network`):
```
cd ui && bun run dev        # http://localhost:7001
cd ui && bun run dev:mock   # mock mode (no backend needed)
```

## Port
| Service | Port |
|---------|------|
| UI (this service) | 7001 |
| SSR Engine (fragment renderer) | 7003 |

## Conventions
- **Components:** `src/lib/components/` — extend existing before creating new. Check folder first.
- **State:** `src/lib/state/` for global stores
- **Types:** `src/lib/types/` — shared TypeScript types
- **Server logic:** `src/lib/server/` — never imported in client code
- **Errors:** Typed error objects. Surface at boundaries only.
- **Tests:** `*.test.ts` adjacent to source in `__tests__/`. E2E in `e2e/`.
- **Accessibility:** ARIA labels, keyboard navigation, contrast — not optional.

# currentDate
Today's date is 2026-03-08.
