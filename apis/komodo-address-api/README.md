# komodo-address-api

Address validation, normalization, and geocoding for the Komodo platform.

> ⚠️ **Migration needed:** This service currently uses the Gin framework. It must be migrated to `net/http` ServeMux (Go 1.26 standard) to align with monorepo conventions before being considered production-ready.

---

## Port

| Server | Port | Env Var |
|--------|------|---------|
| Public | 7031 | `PORT` |

---

## Routes

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/health` | None | Liveness check |
| `POST` | `/validate` | JWT | Validate address correctness and existence |
| `POST` | `/normalize` | JWT | Standardize address formatting |
| `POST` | `/geocode` | JWT | Convert address to lat/lng coordinates |

---

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `ENV` | Yes | Runtime environment (`local`, `dev`, `staging`, `prod`) |
| `PORT` | No | HTTP listen port (default: `7031`) |
| `AUTH_SERVICE_VALIDATE_URL` | Yes | URL to validate bearer tokens |

---

## Local Development

### Prerequisites

- LocalStack running: `cd ../localstack && docker compose up -d`

### Run

```bash
cd apis/komodo-address-api
go run .
```

### cURL Examples

```bash
# Health check
curl http://localhost:7031/health

# Validate an address
curl -X POST http://localhost:7031/validate \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"street":"123 Main St","city":"Chicago","state":"IL","postalCode":"60601"}'

# Normalize an address
curl -X POST http://localhost:7031/normalize \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"street":"123 main street","city":"chicago","state":"il","postalCode":"60601"}'

# Geocode an address
curl -X POST http://localhost:7031/geocode \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"street":"123 Main St","city":"Chicago","state":"IL","postalCode":"60601"}'
```

---

## Status

**Active** — routes implemented and running. Pending migration from Gin to `net/http` ServeMux.

| Key | Value |
|-----|-------|
| Language | Go 1.26 |
| Framework | Gin v1.11.0 (⚠️ non-standard — migrate to net/http) |
| Port | 7031 |
| Domain | Address & Geo |
| Docs | `docs/openapi.yaml` |
