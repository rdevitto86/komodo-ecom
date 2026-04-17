from enum import Enum
from pydantic import BaseModel, Field


class CheckType(str, Enum):
    PII = "pii"
    OBSCENITY = "obscenity"
    INJECTION = "injection"
    TOXICITY = "toxicity"


_ALL_CHECKS = list(CheckType)


class ModerateRequest(BaseModel):
    text: str = Field(..., max_length=8192)
    checks: list[CheckType] = Field(default_factory=lambda: list(_ALL_CHECKS))


class ModerateResponse(BaseModel):
    passed: bool
    flags: list[str]
    redacted_text: str
    provider: str
    latency_ms: int
