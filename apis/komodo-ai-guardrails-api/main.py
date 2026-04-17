import uvicorn
from fastapi import FastAPI

from config.consts import PORT, SERVICE_NAME, VERSION
from routes.health import router as health_router
from routes.moderate import router as moderate_router

app = FastAPI(
    title="Komodo AI Guardrails API",
    version=VERSION,
    description=(
        "Content moderation and guardrails proxy. "
        "Detects PII, prompt injection, and toxicity before forwarding to LLM providers."
    ),
)

app.include_router(health_router)
app.include_router(moderate_router)

if __name__ == "__main__":
    uvicorn.run("main:app", host="0.0.0.0", port=PORT, reload=False)
