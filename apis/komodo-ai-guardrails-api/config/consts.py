import os

PORT: int = int(os.getenv("PORT", "7113"))

SERVICE_NAME: str = "komodo-ai-guardrails-api"
VERSION: str = "0.1.0"

GUARDRAILS_PROVIDER: str = os.getenv("GUARDRAILS_PROVIDER", "local")

OLLAMA_BASE_URL: str = os.getenv("OLLAMA_BASE_URL", "http://localhost:11434")
OLLAMA_MODEL: str = os.getenv("OLLAMA_MODEL", "llama-guard3")

AWS_REGION: str = os.getenv("AWS_REGION", "us-east-1")
AWS_BEDROCK_MODEL_ID: str = os.getenv(
    "AWS_BEDROCK_MODEL_ID", "meta.llama-guard-3-8b-v1:0"
)
