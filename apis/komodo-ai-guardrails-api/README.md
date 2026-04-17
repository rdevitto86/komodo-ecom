# komodo-ai-guardrails-api

## Purpose

A guardrails proxy that sits in front of LLM calls to filter and moderate content. All text passes through this service before reaching an LLM provider. It runs rule-based checks (PII redaction, prompt injection detection) locally with no external dependency, and routes obscenity/toxicity checks to a configurable LLM provider. The provider is swappable at runtime via env var — `local` targets an Ollama instance (dev default), `bedrock` targets AWS Bedrock.

## Port

`7113` — Analytics & Discovery block (7111–7120)

## Routes

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/health` | Health check — returns service name, version, active provider |
| `POST` | `/guardrails/moderate` | Run one or more guardrail checks on a text payload |

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `7113` | HTTP listen port |
| `GUARDRAILS_PROVIDER` | `local` | Active LLM provider: `local` (Ollama) or `bedrock` (AWS Bedrock) |
| `OLLAMA_BASE_URL` | `http://localhost:11434` | Ollama server base URL (used when `GUARDRAILS_PROVIDER=local`) |
| `OLLAMA_MODEL` | `llama-guard3` | Ollama model name for moderation |
| `AWS_REGION` | `us-east-1` | AWS region for Bedrock calls |
| `AWS_BEDROCK_MODEL_ID` | `meta.llama-guard-3-8b-v1:0` | Bedrock model ID (used when `GUARDRAILS_PROVIDER=bedrock`) |

## Run Locally

```bash
pipenv install
pipenv run uvicorn main:app --reload --port 7113
```

Or with Docker:

```bash
docker compose up
```

## Provider Config

**Local (Ollama) — default for dev:**
```
GUARDRAILS_PROVIDER=local
OLLAMA_BASE_URL=http://localhost:11434
OLLAMA_MODEL=llama-guard3
```
Requires an Ollama instance running with the `llama-guard3` model pulled:
```bash
ollama pull llama-guard3
```

**AWS Bedrock:**
```
GUARDRAILS_PROVIDER=bedrock
AWS_REGION=us-east-1
AWS_BEDROCK_MODEL_ID=meta.llama-guard-3-8b-v1:0
AWS_ACCESS_KEY_ID=...
AWS_SECRET_ACCESS_KEY=...
```

## Checks

| Check | Type | Description |
|-------|------|-------------|
| `pii` | Rule-based | Detects and redacts email, US phone, SSN, credit card numbers |
| `injection` | Rule-based | Pattern-matches common prompt injection phrases |
| `obscenity` | LLM-backed (stub) | Routes to active provider for obscenity classification |
| `toxicity` | LLM-backed (stub) | Routes to active provider for toxicity classification |

## TODO

- [ ] Implement Ollama HTTP call in `routes/moderate.py` (`local` provider path)
- [ ] Implement AWS Bedrock `invoke_model` call in `routes/moderate.py` (`bedrock` provider path)
- [ ] Parse and map LLM JSON response (`{"safe": bool, "categories": list}`) to flags
- [ ] Add integration tests against a local Ollama stub
- [ ] Add auth middleware (internal-only service — mTLS or shared secret)
- [ ] Add structured logging with request ID propagation
