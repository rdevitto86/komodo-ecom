# komodo-ssr-engine-svelte

Backend SSR service for the Komodo platform. Pre-renders and caches SvelteKit component trees server-side, delivering HTML fragments to the `ui` frontend shell as injectable slots. Handles CloudFront cache invalidation and S3-backed content storage.

---

## Port

| Server | Port | Env Var |
|--------|------|---------|
| Public | 7003 | `PORT` |

---

## Routes

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/api/health` | None | Liveness check (includes S3 connectivity) |
| `GET` | `/api/v1/landing` | None | Landing page seed data |
| `GET` | `/api/v1/products/{id}` | None | Pre-rendered product page content from S3 |
| `GET` | `/api/v1/services/{id}` | None | Pre-rendered service page content from S3 |
| `GET` | `/api/v1/orders/{id}` | None | Pre-rendered order detail content |
| `GET` | `/api/v1/marketing/content/{id}` | None | Marketing content by ID |
| `GET` | `/api/v1/marketing/user/{id}` | None | User-targeted marketing content |
| `POST` | `/api/v1/admin/manage/content/upsert` | (TODO: admin JWT) | Upsert cached content and invalidate CloudFront |
| `POST` | `/api/v1/admin/manage/content/invalidate` | (TODO: admin JWT) | Invalidate a CloudFront cache key by pageKey |

---

## Architecture

```
ui (7001) ──fetch──► komodo-ssr-engine-svelte (7003)
                              │
                    ┌─────────┴──────────┐
                    S3 (content cache)   CloudFront (CDN)
```

The `ui` shell requests pre-rendered fragments from this service and injects them as Svelte slots for performance-critical routes (product pages, landing, marketing). Content is read from S3 and CloudFront invalidations are triggered on updates.

---

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `PORT` | No | HTTP listen port (default: `7003`) |
| `NODE_ENV` | Yes | Runtime environment (`development`, `production`) |
| `AWS_REGION` | Yes | AWS region (e.g. `us-east-1`) |
| `AWS_ENDPOINT` | No | LocalStack endpoint URL (local dev only) |
| `S3_CONTENT_BUCKET` | Yes | S3 bucket name for pre-rendered content |
| `CLOUDFRONT_DISTRIBUTION_ID` | Yes | CloudFront distribution ID for invalidations |

---

## Local Development

### Prerequisites

- Bun 1.2+
- LocalStack running: `just up` from repo root

### Run

```bash
cd apis/komodo-ssr-engine-svelte
bun install
bun dev          # dev server with HMR on :7003
bun run build    # production build
bun start        # serve production build
```

### Docker

```bash
cd apis/komodo-ssr-engine-svelte
docker compose up --build
```

### Health Check

```bash
curl http://localhost:7003/api/health
# {"status":"ok","timestamp":"...","s3":"ok"}
```

---

## Status

**Active**

| Key | Value |
|-----|-------|
| Runtime | Bun 1.2+ |
| Framework | SvelteKit 5 (adapter-node) |
| Port | 7003 |
| Domain | Frontend & Infrastructure |
| Storage | AWS S3, SQLite (`better-sqlite3`) |
| CDN | AWS CloudFront |
| SDK | [`@komodo-forge-sdk/typescript`](https://github.com/rdevitto86/komodo-forge-sdk-ts) |
