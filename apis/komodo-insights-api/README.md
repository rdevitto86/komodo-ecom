# komodo-insights-api

LLM-powered customer intelligence — item summaries, sentiment analysis, and trending signals.

| Key | Value |
|-----|-------|
| Port | 7111 |
| Domain | Analytics & Discovery |
| Status | Scaffolded |
| Language | Go 1.26 |
| Router | `net/http` ServeMux |
| SDK | `komodo-forge-sdk-go` |

## Routes

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/health` | — | Health check |
| `GET` | `/insights/items/{itemId}/summary` | JWT | Natural-language summary of an item's reviews and product details |
| `GET` | `/insights/items/{itemId}/sentiment` | JWT | Sentiment breakdown and top themes for an item |
| `GET` | `/insights/trending` | JWT | Ranked trending items across the catalog |

## LLM Provider

The `service.SummaryProvider` interface abstracts the backend. Swap the concrete
implementation in `cmd/public/main.go` without touching handler or service code:

```go
// Anthropic
service.NewAnthropicProvider(config.GetConfigValue("LLM_API_KEY"))

// AWS Bedrock
service.NewBedrockProvider(cfg)

// On-prem / OpenAI-compatible
service.NewOpenAICompatProvider(config.GetConfigValue("LLM_PROVIDER_URL"), ...)
```

## Environment Variables

| Variable | Source | Description |
|----------|--------|-------------|
| `APP_NAME` | env | Service name |
| `PORT` | env | Listen address (e.g. `:7111`) |
| `ENV` | env | `local` \| `dev` \| `staging` \| `prod` |
| `LOG_LEVEL` | env | `debug` \| `info` \| `error` |
| `LLM_API_KEY` | Secrets Manager | API key for the configured LLM provider |
| `LLM_PROVIDER_URL` | Secrets Manager | Override endpoint (on-prem / Bedrock); empty for hosted APIs |
| `JWT_PUBLIC_KEY` | Secrets Manager | RSA public key for token validation |
| `JWT_PRIVATE_KEY` | Secrets Manager | Required by `jwt.InitializeKeys()` — not used for signing |
| `JWT_AUDIENCE` | Secrets Manager | Expected JWT audience claim |
| `JWT_ISSUER` | Secrets Manager | Expected JWT issuer claim |
| `JWT_KID` | Secrets Manager | Key ID |
| `RATE_LIMIT_RPS` | Secrets Manager | Requests per second per client |
| `RATE_LIMIT_BURST` | Secrets Manager | Burst allowance |

## Running Locally

```bash
make bootstrap          # build + start (local)
make bootstrap ENV=dev  # against LocalStack dev
make stop               # tear down
make test_e2e           # e2e suite (requires just up api first)
```

**After cloning / first run:**

```bash
cd apis/komodo-insights-api
go mod tidy   # populate go.sum
```

## Status

**Scaffolded** — directory structure, middleware stack, and route stubs are in place.
All handlers return 404/500 until the LLM provider is wired and upstream clients
(reviews-api, shop-items-api) are integrated.

Next steps before implementing:
- Add `httpErr.RangeInsights` to `komodo-forge-sdk-go/http/errors/ranges.go`
- Choose LLM backend and add concrete `SummaryProvider` implementation
- Decide data-fetching strategy: on-demand vs. pre-computed + cached
