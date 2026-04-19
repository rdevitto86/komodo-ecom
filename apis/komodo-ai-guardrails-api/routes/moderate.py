import time

from fastapi import APIRouter

from config.consts import GUARDRAILS_PROVIDER
from models.models import CheckType, ModerateRequest, ModerateResponse
from utils.logger import get_logger
from utils.redaction import redact_pii
from utils.sanitization import (
    filter_coding,
    filter_formatting,
    filter_injection,
    filter_personality_deviation,
    filter_pii_insertion,
    normalize,
)

log = get_logger("routes.moderate")
router = APIRouter()


@router.post("/guardrails/moderate", response_model=ModerateResponse)
async def moderate(req: ModerateRequest) -> ModerateResponse:
    start = time.monotonic()
    flags: list[str] = []

    try:
        # 1. Normalize and format-clean input
        text = normalize(req.text)
        text = filter_formatting(text)
        redacted_text = text

        # 2. PII — redact leaks in user input, detect coercion attempts
        if CheckType.PII in req.checks:
            redacted_text, pii_flags = redact_pii(text)
            flags.extend(pii_flags)
            flags.extend(filter_pii_insertion(text))

        # 3. Prompt injection — instruction overrides, delimiter escapes, encoded payloads
        if CheckType.INJECTION in req.checks:
            _, injection_flags = filter_injection(text)
            flags.extend(injection_flags)

        # 4. Personality deviation — persona override and jailbreak attempts
        if CheckType.DEVIATION in req.checks:
            _, deviation_flags = filter_personality_deviation(text)
            flags.extend(deviation_flags)

        # 5. Coding — out-of-scope code generation / execution requests
        if CheckType.CODING in req.checks:
            _, coding_flags = filter_coding(text)
            flags.extend(coding_flags)

        # 6. Obscenity / toxicity — LLM-backed (stubbed)
        if CheckType.OBSCENITY in req.checks or CheckType.TOXICITY in req.checks:
            # TODO: implement — route text to GUARDRAILS_PROVIDER
            # local  → POST {OLLAMA_BASE_URL}/api/chat with OLLAMA_MODEL
            # bedrock → boto3 bedrock-runtime invoke_model with AWS_BEDROCK_MODEL_ID
            # Parse JSON response: {"safe": bool, "categories": list[str]}
            # Append returned categories to flags if not safe
            pass

    except Exception as exc:
        latency_ms = int((time.monotonic() - start) * 1000)
        log.error("moderate check failed", error=str(exc), latency_ms=latency_ms)
        raise

    latency_ms = int((time.monotonic() - start) * 1000)
    passed = len(flags) == 0

    log.info(
        "request moderated",
        passed=passed,
        flags=flags or None,
        checks=req.checks,
        latency_ms=latency_ms,
    )

    return ModerateResponse(
        passed=passed,
        flags=flags,
        redacted_text=redacted_text,
        provider=GUARDRAILS_PROVIDER,
        latency_ms=latency_ms,
    )
