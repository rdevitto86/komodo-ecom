# Skill: /new-page

Scaffold a new SvelteKit 5 page (route) in the UI following established conventions.

## Usage

```
/new-page <route-path> [--auth none|optional|required] [--api] [--layout <layout-group>]
```

- `<route-path>` — URL path, e.g. `profile/orders`, `shop/wishlist`, `admin/products`. Used to derive the filesystem path under `src/routes/`.
- `--auth` — Auth level (default: `none`):
  - `none` — Public page. No session check.
  - `optional` — Load user data if a valid session exists; continue if not.
  - `required` — Redirect to `/login` if no valid session. Guard runs in `+page.server.ts`.
- `--api` — Also generate a `+server.ts` BFF route at the same path for client-side API calls (e.g. form actions, SPA navigation). Generates `GET` and `POST` stubs by default.
- `--layout <layout-group>` — Wrap in a named layout group, e.g. `--layout auth` creates `src/routes/(auth)/<route-path>/`. Omit for default layout.

**Must be run from inside `ui/`.**

---

## Before generating anything

1. Read `ui/docs/README.md` — check the route table to see if the route already exists or conflicts with an existing path.
2. Read `src/routes/+layout.server.ts` — understand what data the root layout already provides (session, user, etc.) so the new page doesn't duplicate it.
3. If `--auth required` or `optional`, read `src/routes/(auth)/` to understand the existing auth guard pattern.
4. If `--layout` is specified, check whether that layout group and its `+layout.svelte` / `+layout.server.ts` already exist.

---

## Target state: `adapter-node`

All generated code assumes `adapter-node` (SSR + server-side loading). **Never generate static-only patterns:**
- Always generate `+page.server.ts` — even for simple public pages. A minimal `load` returning `{}` is correct for now; it marks the page as server-rendered and provides the extension point for later.
- Use `export const load: PageServerLoad` — not `PageLoad` (client-side load).
- Form actions go in `+page.server.ts` under `export const actions`.
- Do not use `export const prerender = true` or `export const ssr = false` unless explicitly requested.

---

## Files to generate

### `+page.server.ts`

**`--auth none` (public):**
```typescript
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async () => {
  // TODO: load page data
  return {};
};
```

**`--auth optional`:**
```typescript
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ locals }) => {
  // locals.user is populated by +layout.server.ts if a valid session exists.
  // Page should degrade gracefully when locals.user is undefined.
  return {
    // user: locals.user ?? null,
  };
};
```

**`--auth required`:**
```typescript
import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ locals }) => {
  // TODO: update once auth-api session validation is wired in +layout.server.ts
  if (!locals.user) {
    redirect(307, '/login');
  }
  return {
    // user: locals.user,
  };
};
```

---

### `+page.svelte`

Svelte 5 rules — identical to `/new-component` constraints:
- `$props()` for data from load, never `export let data`.
- `$state()` for local reactive state.
- `$derived()` for computed values.
- No `$:` reactive statements.

```svelte
<script lang="ts">
  import type { PageData } from './$types';

  let { data }: { data: PageData } = $props();
</script>

<svelte:head>
  <title>TODO — Komodo</title>
</svelte:head>

<!-- TODO: implement page -->
<main>
  <h1>TODO</h1>
</main>
```

---

### `+server.ts` (only with `--api`)

Generated at the same route path as the page. Returns RFC 7807 stubs until backed by a real service.

```typescript
import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

// TODO: wire to backend service when available.
const notImplemented = () =>
  json(
    { type: 'about:blank', title: 'Not Implemented', status: 501, detail: 'TODO: describe the planned endpoint.' },
    { status: 501 },
  );

export const GET: RequestHandler = notImplemented;
export const POST: RequestHandler = notImplemented;
```

Adjust exported methods to match the actual planned contract.

---

### Layout group files (only with `--layout <group>` if group doesn't yet exist)

If `src/routes/(<group>)/` does not exist, also generate:

**`+layout.server.ts`:**
```typescript
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async () => {
  // Shared layout data for the (<group>) group.
  return {};
};
```

**`+layout.svelte`:**
```svelte
<script lang="ts">
  let { children } = $props();
</script>

{@render children()}
```

---

## After generating

1. Update the route table in `ui/docs/README.md` — add a row for the new route with its path, auth level, and the backing service (or "static" / "BFF stub").
2. Print all files created and their paths.
3. Remind the developer to:
   - Run `bun run dev` and navigate to the route to confirm no build errors.
   - Fill in the `load` function once the backing service is ready.
   - Add E2E test in `e2e/` once the page is functional.
   - For `--auth required` pages: add a Playwright test that confirms unauthenticated users are redirected.
