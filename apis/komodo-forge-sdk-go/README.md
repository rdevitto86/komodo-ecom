# komodo-forge-sdk-go

Internal Go SDK for all Komodo services. Provides AWS clients, HTTP middleware, JWT/OAuth crypto, structured logging, concurrency utilities, and a universal server entry point.

Module: `komodo-forge-sdk-go`

All services reference this SDK via a local `replace` directive in their `go.mod`:
```
replace komodo-forge-sdk-go => ../komodo-forge-sdk-go
```

---

## Packages

### `aws/secrets-manager`
Bootstraps AWS Secrets Manager at service startup. Resolves secrets into in-memory config.

```go
awsSM.Bootstrap(cfg)
awsSM.GetSecret(key, prefix)
```

### `aws/s3`
S3 client with typed get/put/delete.

```go
awsS3.Init(cfg)
awsS3.GetObject(ctx, bucket, key)           // returns []byte
awsS3.GetObjectAs(ctx, bucket, key, &out)   // unmarshals JSON into out
awsS3.PutObject(ctx, bucket, key, data, contentType)
awsS3.DeleteObject(ctx, bucket, key)
```

### `aws/dynamodb`
DynamoDB client with typed CRUD and query/scan helpers.

```go
dynamodb.Init(cfg)
dynamodb.GetItem(ctx, table, key, &out)
dynamodb.WriteItem(ctx, table, item)
dynamodb.UpdateItem(ctx, table, key, updates)
dynamodb.DeleteItem(ctx, table, key)
dynamodb.Query(ctx, input)
dynamodb.QueryAll(ctx, input)
dynamodb.Scan(ctx, input)
```

### `aws/elasticache`
Redis client via `go-redis`. Includes distributed token-bucket rate limiter.

```go
elasticache.Init(cfg)
elasticache.Get(key)
elasticache.Set(key, value, ttlSec)
elasticache.Delete(key)
elasticache.AllowDistributed(ctx, key, rate, burst, ttlSec) // returns (allowed bool, retryAfter time.Duration, err)
```

### `aws/aurora`
Aurora-compatible SQL client (planned).

### `config`
In-memory config store. Checks local store first, falls back to `os.Getenv`. Used by the Secrets Manager bootstrap to inject resolved secrets.

```go
config.GetConfigValue(key)
config.SetConfigValue(key, val)
config.DeleteConfigValue(key)
config.GetAllKeys()
```

### `crypto/jwt`
RS256 JWT signing and validation.

```go
jwt.InitializeKeys()
jwt.SignToken(issuer, subject, audience, ttl, scopes)
jwt.ValidateToken(tokenString)
jwt.ParseClaims(tokenString)
jwt.ExtractTokenFromRequest(req)
```

### `crypto/oauth`
OAuth 2.0 scope and grant type validation.

```go
oauth.IsValidScope(scope)
oauth.IsValidGrantType(grantType)
```

### `http/middleware`
Full middleware stack, exported from `http/middleware/exports.go`.

| Middleware | Description |
|-----------|-------------|
| `AuthMiddleware` | Validates RS256 JWT from Authorization header |
| `ClientTypeMiddleware` | Identifies client type (M2M vs user) |
| `CORSMiddleware` | CORS preflight and headers |
| `CSRFMiddleware` | CSRF token validation |
| `IdempotencyMiddleware` | Idempotency key deduplication via Redis |
| `IPAccessMiddleware` | IP whitelist/blacklist enforcement |
| `NormalizationMiddleware` | Request path and header normalization |
| `RateLimiterMiddleware` | Token bucket rate limiting via Redis |
| `RedactionMiddleware` | Redacts sensitive fields from logs |
| `RequestIDMiddleware` | Injects `X-Request-ID` |
| `RuleValidationMiddleware` | YAML-driven field validation rules |
| `SanitizationMiddleware` | Input sanitization (XSS, injection) |
| `ScopeMiddleware` | Requires `svc:` prefixed JWT scope (internal routes) |
| `SecurityHeadersMiddleware` | Sets security response headers |
| `TelemetryMiddleware` | Request/response telemetry logging |
| `Chain` | Composes a slice of middleware onto an `http.Handler` |

### `http/server`
Universal service entry point. Auto-detects Lambda vs Fargate/local.

```go
// On AWS Lambda (AWS_LAMBDA_FUNCTION_NAME set): starts lambda.Start with httpadapter
// Otherwise: ListenAndServe with graceful SIGINT/SIGTERM shutdown
server.Run(srv, port, gracefulTimeout)
```

### `http/errors`
RFC 7807 Problem+JSON error responses.

```go
errors.SendError(w, r, errCode, overrides...)
errors.SendCustomError(w, r, status, message, detail, code)
errors.WithMessage(msg)
errors.WithDetail(detail)
errors.WithStatus(status)
```

### `http/request`
Request parsing helpers.

```go
request.GetPathParams(req)
request.GetQueryParams(req)
request.GetClientType(req)
request.GetRequestID(req)
request.GenerateRequestId()
request.ExtractTokenFromRequest(req)
```

### `http/response`
Response writer wrapper with status tracking.

```go
response.Bind(res, &target)
response.IsSuccess(status)
response.IsError(status)
```

### `logging/runtime`
`slog`-based structured logger. `tint` handler for local/dev, JSON handler for staging/prod. Includes automatic field redaction.

```go
logger.Init(appName, logLevel, env)
logger.Debug(msg, args...)
logger.Info(msg, args...)
logger.Warn(msg, args...)
logger.Error(msg, args...)
logger.Fatal(msg, err)
```

### `concurrency/worker`
Bounded goroutine worker pool.

```go
pool, _ := worker.NewWorkerPool(cfg)
pool.Submit(ctx, job)
pool.SubmitAsync(ctx, job)  // returns <-chan error
pool.Shutdown(ctx)
```

### `concurrency/semaphore`
Counting semaphore for concurrency control.

---

## Testing Utilities

### `testing/moxtox`
HTTP mock server builder for unit testing middleware and handlers without network calls.

### `testing/chaos`
Chaos injection helpers — latency simulation and error injection.

### `testing/performance`
Latency measurement and benchmarking utilities.

---

## Development

```bash
cd apis/komodo-forge-sdk-go
make test     # run all tests with race detector
make lint     # run golangci-lint
make audit    # security audit
make coverage # generate coverage report
```
