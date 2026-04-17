# komodo-ai-guardrails-api — Low-Level Design

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

## TODO

- [ ] Implement Ollama HTTP call in `routes/moderate.py` (`local` provider path)
- [ ] Implement AWS Bedrock `invoke_model` call in `routes/moderate.py` (`bedrock` provider path)
- [ ] Parse and map LLM JSON response (`{"safe": bool, "categories": list}`) to flags
- [ ] Add integration tests against a local Ollama stub
- [ ] Add auth middleware (internal-only service — mTLS or shared secret)
- [ ] Add structured logging with request ID propagation
