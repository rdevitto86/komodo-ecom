# Skill: /new-service

Scaffold a new Go microservice in this monorepo following established conventions.

## Usage

```
/new-service <name> <domain> [--lambda|--fargate]
```

- `<name>` — short name without prefix/suffix (e.g. `cart`, `loyalty`, `returns`). Full name becomes `komodo-<name>-api`.
- `<domain>` — domain group for port allocation. Must match one of the domain blocks in the Port Allocation table in `CLAUDE.md`.
- `--lambda` or `--fargate` — compute target. Defaults to `--fargate` if omitted.

## Before generating anything

1. Read `CLAUDE.md` for the Port Allocation table. Scan existing `docker-compose.yaml` files in sibling service directories to find which ports in the domain block are already taken. Pick the next available anchor port.
2. Read `apis/komodo-forge-sdk-go/http/server/` to get the exact `srv.Run` signature.
3. Read `apis/komodo-forge-sdk-go/http/middleware/exports.go` for the middleware list.
4. Read `apis/komodo-forge-sdk-go/http/errors/` for `httpErr.SendError` signature.

## Files to generate

Generate all files under `apis/komodo-<name>-api/`.

---

### `go.mod`

```
module komodo-<name>-api

go 1.26

require komodo-forge-sdk-go v0.1.0

// TODO: replace with github.com/komodo-hq/forge-sdk-go when extracted to its own repo
replace komodo-forge-sdk-go => ../komodo-forge-sdk-go
```

Run `go mod tidy` from the service directory after generating.

---

### `main.go`

**Fargate target:**

```go
package main

import (
	awsSM "komodo-forge-sdk-go/aws/secrets-manager"
	"komodo-forge-sdk-go/config"
	mw "komodo-forge-sdk-go/http/middleware"
	srv "komodo-forge-sdk-go/http/server"
	logger "komodo-forge-sdk-go/logging/runtime"
	"komodo-<name>-api/internal/handlers"
	"net/http"
	"os"
	"time"
)

func init() {
	logger.Init(
		config.GetConfigValue("APP_NAME"),
		config.GetConfigValue("LOG_LEVEL"),
		config.GetConfigValue("ENV"),
	)

	smCfg := awsSM.Config{
		Region:   config.GetConfigValue("AWS_REGION"),
		Endpoint: config.GetConfigValue("AWS_ENDPOINT"),
		Prefix:   config.GetConfigValue("AWS_SECRET_PREFIX"),
		Batch:    config.GetConfigValue("AWS_SECRET_BATCH"),
		Keys: []string{
			// TODO: add service-specific secret keys
			"IP_WHITELIST",
			"IP_BLACKLIST",
			"MAX_CONTENT_LENGTH",
			"IDEMPOTENCY_TTL_SEC",
			"RATE_LIMIT_RPS",
			"RATE_LIMIT_BURST",
			"BUCKET_TTL_SECOND",
		},
	}
	if err := awsSM.Bootstrap(smCfg); err != nil {
		logger.Fatal("failed to initialize secrets manager", err)
		os.Exit(1)
	}

	logger.Info("komodo-<name>-api: bootstrap complete")
}

func main() {
	readMW := []func(http.Handler) http.Handler{
		mw.RequestIDMiddleware,
		mw.TelemetryMiddleware,
		mw.RateLimiterMiddleware,
		mw.CORSMiddleware,
		mw.SecurityHeadersMiddleware,
		mw.AuthMiddleware,
		mw.NormalizationMiddleware,
		mw.RuleValidationMiddleware,
		mw.SanitizationMiddleware,
	}

	writeMW := append(readMW, mw.IdempotencyMiddleware)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handlers.HealthHandler)

	// TODO: register routes here
	// mux.Handle("GET /example", mw.Chain(handlers.GetExample, readMW...))
	// mux.Handle("POST /example", mw.Chain(handlers.CreateExample, writeMW...))
	_ = writeMW

	server := &http.Server{
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	srv.Run(server, config.GetConfigValue("PORT"), 30*time.Second)
}
```

**Lambda target** — move bootstrap from `init()` into `main()` and use `srv.Run` which auto-detects `AWS_LAMBDA_FUNCTION_NAME`. Structure is otherwise identical.

---

### `internal/handlers/health.go`

```go
package handlers

import (
	"encoding/json"
	"net/http"
)

func HealthHandler(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")
	wtr.WriteHeader(http.StatusOK)
	json.NewEncoder(wtr).Encode(map[string]string{"status": "ok"})
}
```

---

### `pkg/v1/models/errors.go`

```go
package models

import httpErr "komodo-forge-sdk-go/http/errors"

var Err = struct {
	// TODO: define service-specific errors
	// Example:
	// NotFound httpErr.ErrorCode
}{
	// NotFound: httpErr.ErrorCode{...},
}

// Ensure the forge SDK error types are available for handler use.
var _ = httpErr.Global
```

---

### `pkg/v1/models/<name>.go`

Skeleton with a placeholder domain model:

```go
package models

// TODO: define domain models for <name>-api
```

---

### `pkg/v1/adapter/adapter.go`

```go
package adapter

import "net/http"

// Client is the typed HTTP client for calling komodo-<name>-api.
// Consuming services inject <NAME>_API_INTERNAL_URL.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

// TODO: add typed methods matching the service's routes.
```

---

### `pkg/v1/exports.go`

```go
package v1

import "komodo-<name>-api/pkg/v1/adapter"

// Adapter is the typed HTTP client for calling komodo-<name>-api.
type Adapter = adapter.Client

var NewAdapter = adapter.NewClient
```

---

### `pkg/v1/mocks/exports.go`

```go
package mocks

// Mocks holds JSON fixture loaders for komodo-<name>-api responses.
// Add JSON files alongside this file and expose them here.
var Mocks = struct{}{} // TODO: add mock fixtures
```

---

### `Dockerfile`

```dockerfile
FROM golang:1.26 AS build

COPY komodo-forge-sdk-go /komodo-forge-sdk-go

WORKDIR /app

COPY komodo-<name>-api/go.mod komodo-<name>-api/go.sum ./
RUN go mod download

COPY komodo-<name>-api ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/komodo ./...

FROM gcr.io/distroless/base-debian12
COPY --from=build /bin/komodo /komodo
EXPOSE <PORT>
ENTRYPOINT ["/komodo"]
```

Build context is `..` (the `apis/` directory) so the forge SDK can be `COPY`'d. Always run `docker compose` from inside `apis/komodo-<name>-api/`.

---

### `docker-compose.yaml`

```yaml
name: komodo-<name>-api

services:
  <name>-api:
    image: ${APP_NAME}:${VERSION:-latest}
    build:
      context: ..
      dockerfile: komodo-<name>-api/Dockerfile
      args:
        ENV: ${ENV:-local}
    restart: ${RESTART_POLICY:-no}
    deploy:
      resources:
        limits:
          memory: ${MEM_LIMIT:-512M}
    environment:
      APP_NAME: ${APP_NAME}
      ENV: ${ENV}
      PORT: <PORT>
      VERSION: ${VERSION}
      AWS_REGION: ${AWS_REGION}
      AWS_ENDPOINT: ${AWS_ENDPOINT}
      AWS_SECRET_PREFIX: ${AWS_SECRET_PREFIX}
      AWS_SECRET_BATCH: ${AWS_SECRET_BATCH}
    ports:
      - "<PORT>:<PORT>"
    extra_hosts:
      - "host.docker.internal:host-gateway"
    healthcheck:
      test: ["CMD", "wget", "-q", "-O", "-", "http://localhost:<PORT>/health"]
      interval: 30s
      timeout: 5s
      retries: 3
      start_period: 40s
    networks:
      - komodo-network
networks:
  komodo-network:
    external: true
```

---

### `Makefile`

Copy the pattern from any sibling service, substituting `APP_NAME := komodo-<name>-api` and the correct `AWS_SECRET_PREFIX`.

---

### `.golangci.yaml`

Copy verbatim from any sibling service. The `wrapcheck.ignore-package-globs` entry must include the forge SDK glob — update it when the SDK moves to its GitHub path.

---

### `docs/` skeleton

Create all five files with placeholder content:

| File | Minimum content |
|------|----------------|
| `README.md` | Port, run commands, env vars table (copy structure from a sibling service) |
| `openapi.yaml` | Minimal OpenAPI 3.1 spec with `/health` GET only |
| `architecture.md` | H1 heading + one-line description |
| `design-decisions.md` | H1 heading only |
| `data-model.md` | H1 heading only |

---

## After generating

1. Remind the developer to:
   - Add the service to `infra/local/services.jsonc` under the correct profile group (`api`, `ui`, or `support`)
   - Add a profile entry in `infra/local/docker-compose.yml` if it isn't auto-included
   - Update the port allocation table in root `CLAUDE.md` and `MEMORY.md`
   - Run `go mod tidy` from inside the service directory
2. Print the allocated port and full service name.
