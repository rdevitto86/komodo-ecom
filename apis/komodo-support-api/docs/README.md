# komodo-support-api

AI-powered customer support chatbot with human escalation path.

| Key | Value |
|-----|-------|
| Port | 7101 |
| Domain | Support & CX |
| Status | Implemented (in-memory store — wire DynamoDB before prod) |
| Language | Go 1.26 |
| Router | `net/http` ServeMux |
| SDK | `komodo-forge-sdk-go` |
| LLM | Anthropic Haiku 4.5 (swappable — see `service.LLMProvider`) |

---

## Routes

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/health` | None | Health check |
| `POST` | `/chat/session` | None | Create anonymous session (sets `komodo_chat_sid` cookie) |
| `GET` | `/chat/session` | None | Validate session cookie |
| `POST` | `/chat/message` | Cookie or JWT | Send message, receive AI response |
| `GET` | `/chat/history` | Cookie or JWT | Fetch conversation history |
| `DELETE` | `/chat/history` | Cookie or JWT | Clear conversation history |
| `POST` | `/chat/escalate` | Cookie or JWT | Explicitly escalate to human agent |
| `GET` | `/me/chat/history` | JWT required | Persistent history for authenticated users |
| `DELETE` | `/me/chat/history` | JWT required | Clear persistent history for authenticated users |

---

## Session Model

**Anonymous:** `POST /chat/session` → generates UUID `session_id` → sets `komodo_chat_sid` cookie (HTTPOnly, Secure, SameSite=Strict). TTL configurable via `CHAT_SESSION_TTL_DAYS` (default: 30 days).

**Authenticated:** Session ID is derived from the JWT `user_id` claim. No cookie needed.

**Merge path (TODO):** Call `repo.MergeSession(sessionID, userID)` at login to associate an anonymous history with the user's account. Not yet wired to the auth flow.

---

## Escalation

Two signals trigger escalation:

1. **Model signal** — model prefixes response with `[ESCALATE]` (instructed via system prompt for: human requests, legal threats, fraud claims, repeated failures)
2. **Keyword scan** — client-side check on user input (phrases like "speak to human", "chargeback", "lawyer", etc.)

Escalated flag is persisted on the session. The `POST /chat/escalate` endpoint handles explicit user-initiated escalation.

**TODO:** Wire escalation to `komodo-communications-api` for async ticket creation (human handoff model not yet decided — see design decision pending).

---

## LLM Provider

`service.LLMProvider` is an interface — swap implementations in `main.go`:

```go
// Current default
llm := service.NewAnthropicProvider(config.GetConfigValue("ANTHROPIC_API_KEY"))

// Future: OpenAI, Bedrock, Ollama, etc.
// llm := service.NewOpenAIProvider(...)
```

Provider implementations live in `internal/service/`. Each one handles system prompt formatting and message history mapping for its API.

---

## Environment Variables

| Variable | Source | Description |
|----------|--------|-------------|
| `PORT` | Env | `7101` |
| `ENV` | Env | `local` / `staging` / `prod` |
| `APP_NAME` | Env | Service name for logging |
| `AWS_REGION` | Env | AWS region |
| `AWS_ENDPOINT` | Env | LocalStack endpoint (local only) |
| `AWS_SECRET_PREFIX` | Env | Secrets Manager key prefix |
| `AWS_SECRET_BATCH` | Env | Secrets Manager batch name |
| `ANTHROPIC_API_KEY` | Secrets Manager | Anthropic API key |
| `CHAT_SESSION_TTL_DAYS` | Secrets Manager | Anonymous session TTL (default: 30) |
| `CHAT_MAX_HISTORY` | Secrets Manager | Max history turns in LLM context (default: 20) |
| `RATE_LIMIT_RPS` | Secrets Manager | Rate limiter requests/sec |
| `RATE_LIMIT_BURST` | Secrets Manager | Rate limiter burst |
| `IP_WHITELIST` | Secrets Manager | Allowed IPs (optional) |
| `IP_BLACKLIST` | Secrets Manager | Blocked IPs (optional) |

---

## Storage (TODO: DynamoDB)

Currently using in-memory store (`repository.InMemoryChatRepository`). Replace with DynamoDB before production.

Proposed table design:

| Attribute | Type | Notes |
|-----------|------|-------|
| `PK` | `session_id` | Partition key for all records |
| `SK` | `"META"` or `message_id` | Sort key — META for session record, UUID for messages |
| `user_id` | GSI PK | Used to fetch all sessions for a logged-in user |
| `created_at` | GSI SK | Sort by time within user's sessions |
| `expires_at` | TTL | DynamoDB TTL auto-deletes anonymous sessions |

---

## Run Commands

```bash
# Local (from service root)
go run .

# Docker (from apis/ directory)
docker compose -f komodo-support-api/docker-compose.yaml up
```
