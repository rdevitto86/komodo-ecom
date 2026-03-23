# Skill: /add-route

Add a new route (handler + test stub + `main.go` wiring + OpenAPI stub) to an existing Go microservice.

## Usage

```
/add-route <METHOD> <path> [--open|--protected|--write] [--file <handler-file>]
```

- `<METHOD>` — HTTP method: `GET`, `POST`, `PUT`, `PATCH`, `DELETE`
- `<path>` — Route path in Go ServeMux syntax, e.g. `/order/{id}`, `/order/{id}/status`
- `--open` — No auth. Use for public read endpoints (default if omitted and method is GET).
- `--protected` — Auth + CSRF + normalization + sanitization + rule validation. Use for authenticated reads and state changes that don't need idempotency.
- `--write` — Protected stack + idempotency middleware. Use for POST/PUT/PATCH/DELETE on resources that must be safe to retry.
- `--file <handler-file>` — Which file under `internal/handlers/` to append the handler to. Inferred from the resource name in the path if omitted (e.g. `/order/{id}` → `order.go`).

**Must be run from inside the service directory** (`apis/komodo-<name>-api/`).

---

## Before generating anything

1. Read `main.go` — understand the middleware stack names and which `chain`/`mw.Chain` style is used.
2. Read `internal/handlers/<file>.go` (the target file, or closest resource match) — follow naming conventions already established.
3. Read `apis/komodo-forge-sdk-go/http/middleware/exports.go` — confirm middleware names.
4. Read `apis/komodo-forge-sdk-go/http/errors/` — confirm `httpErr.SendError` signature and available error sets.
5. Read `pkg/v1/models/errors.go` — see what error codes exist; add new ones only if the handler needs them.

---

## Handler stub

Append to `internal/handlers/<file>.go`. Follow the exact style already in that file.

```go
// <METHOD> <path> — TODO: describe what this handler does.
func <HandlerName>(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")

	// TODO: extract path/query params, validate, call service layer
	httpErr.SendError(wtr, req, httpErr.Global.NotImplemented, httpErr.WithDetail("<HandlerName> not yet implemented"))
}
```

**Naming convention:** verb + resource + qualifier, e.g. `GetOrder`, `CreateOrder`, `UpdateOrderStatus`, `DeleteOrder`. Match the casing and style of existing handlers in the file.

---

## Test stub

Create (or append to) `internal/handlers/<file>_test.go` adjacent to the handler file.

```go
package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test<HandlerName>(t *testing.T) {
	t.Run("returns 501 when not implemented", func(t *testing.T) {
		req := httptest.NewRequest(http.Method<METHOD>, "<path>", nil)
		wtr := httptest.NewRecorder()

		<HandlerName>(wtr, req)

		if wtr.Code != http.StatusNotImplemented {
			t.Errorf("expected 501, got %d", wtr.Code)
		}
	})

	// TODO: add test cases for success, validation errors, not-found, auth failures
}
```

If the test file already exists, append the new `Test<HandlerName>` function — do not replace existing tests.

---

## `main.go` wiring

Read the existing route block in `main.go`. Locate the appropriate middleware stack variable for the requested protection level:

| Flag | Middleware stack to use |
|------|------------------------|
| `--open` | The base/read stack (no auth, no CSRF) |
| `--protected` | The auth/protected stack |
| `--write` | The auth/protected stack + idempotency |

**If the service uses `mw.Chain` (new-style):**
```go
mux.Handle("<METHOD> <path>", mw.Chain(http.HandlerFunc(handlers.<HandlerName>), <stack>...))
```

**If the service uses a local `chain` function (existing style):**
```go
mux.Handle("<METHOD> <path>", chain(http.HandlerFunc(handlers.<HandlerName>), <stack>...))
```

Do not change the style — follow what `main.go` already uses. Insert the new route near other routes for the same resource group (not arbitrarily at the end).

---

## OpenAPI stub

Append a minimal path entry to `docs/openapi.yaml`. If the file is empty or has only the minimal `/health` spec, add the new path. If the path already exists (e.g. adding a second method), add a new operation under the existing path key.

```yaml
  <path>:
    <method_lower>:
      summary: TODO
      operationId: <handlerName camelCase>
      tags:
        - <resource>
      # TODO: add parameters, requestBody, responses
      responses:
        '200':
          description: Success
        '501':
          description: Not implemented
```

---

## Error codes

If the new handler will need domain-specific error codes that don't exist yet, append them to `pkg/v1/models/errors.go`:

```go
var Err = struct {
    // existing fields ...
    <NewError> httpErr.ErrorCode
}{
    // existing values ...
    // <NewError>: httpErr.ErrorCode{ ... }, // TODO: fill in when implementing
}
```

Only add error stubs for errors that are clearly needed by this route's contract (e.g. not-found, invalid param). Do not invent errors speculatively.

---

## After generating

1. Print the full route as registered: `<METHOD> <path>` with the middleware level.
2. Remind the developer to:
   - Run `go build ./...` from the service directory to confirm the handler compiles.
   - Run `go test ./internal/handlers/...` to confirm the test stub passes.
   - Fill in the handler body and replace the 501 stub before shipping.
   - Update `docs/openapi.yaml` responses with the real schema once the handler is implemented.
