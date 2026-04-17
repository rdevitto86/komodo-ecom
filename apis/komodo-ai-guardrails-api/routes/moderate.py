import time

from fastapi import APIRouter

from config.consts import GUARDRAILS_PROVIDER
from models.models import CheckType, ModerateRequest, ModerateResponse
from utils.sanitization import normalize
from utils.redaction import redact_pii


router = APIRouter()

_INJECTION_PHRASES = [
    "ignore previous instructions",
    "ignore all instructions",
    "you are now",
    "disregard",
    "forget everything",
    "system prompt",
    "jailbreak",
]


@router.post("/guardrails/moderate", response_model=ModerateResponse)
async def moderate(req: ModerateRequest) -> ModerateResponse:
    start = time.monotonic()
    flags: list[str] = []

    # 1. Normalize input
    text = normalize(req.text)
    redacted_text = text

    # 2. PII detection and redaction
    if CheckType.PII in req.checks:
        redacted_text, pii_flags = redact_pii(text)
        flags.extend(pii_flags)

    # 3. Prompt injection detection (pattern-based)
    if CheckType.INJECTION in req.checks:
        lower = text.lower()
        if any(phrase in lower for phrase in _INJECTION_PHRASES):
            flags.append("injection")

    # 4. Obscenity / toxicity — LLM-backed (stubbed)
    if CheckType.OBSCENITY in req.checks or CheckType.TOXICITY in req.checks:
        # TODO: implement — route text to GUARDRAILS_PROVIDER
        # local  → POST {OLLAMA_BASE_URL}/api/chat with OLLAMA_MODEL
        # bedrock → boto3 bedrock-runtime invoke_model with AWS_BEDROCK_MODEL_ID
        # Parse JSON response: {"safe": bool, "categories": list[str]}
        # Append returned categories to flags if not safe
        pass

    latency_ms = int((time.monotonic() - start) * 1000)

    return ModerateResponse(
        passed=len(flags) == 0,
        flags=flags,
        redacted_text=redacted_text,
        provider=GUARDRAILS_PROVIDER,
        latency_ms=latency_ms,
    )
