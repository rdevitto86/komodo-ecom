# komodo-search-api

Full-text and faceted search across the Komodo product and service catalog.

> ⚠️ **Migration needed:** `go.mod` currently depends on `gorilla/mux`. Must be migrated to `net/http` ServeMux before implementation begins.

| Key | Value |
|-----|-------|
| Port | 7042 |
| Domain | Commerce & Catalog |
| Status | Stub |
| Language | Go 1.26 |
| Router | gorilla/mux (⚠️ non-standard — migrate to net/http) |
| SDK | `komodo-forge-sdk-go` |

**Stub** — `go.mod` and `openapi.yaml` exist, no handler logic. When implemented, add `docs/README.md` and replace gorilla/mux with `net/http` ServeMux.
