import re
from typing import NamedTuple

# ---------------------------------------------------------------------------
# PII patterns
# ---------------------------------------------------------------------------

class _Pattern(NamedTuple):
    flag: str
    placeholder: str
    regex: re.Pattern


_PATTERNS: list[_Pattern] = [
    _Pattern(
        flag="pii.email",
        placeholder="[EMAIL]",
        regex=re.compile(
            r"\b[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,}\b"
        ),
    ),
    _Pattern(
        flag="pii.phone",
        placeholder="[PHONE]",
        regex=re.compile(
            r"\b(\+?1[\s.\-]?)?(\(?\d{3}\)?[\s.\-]?)\d{3}[\s.\-]?\d{4}\b"
        ),
    ),
    _Pattern(
        flag="pii.ssn",
        placeholder="[SSN]",
        regex=re.compile(
            r"\b\d{3}[-\s]?\d{2}[-\s]?\d{4}\b"
        ),
    ),
    _Pattern(
        flag="pii.credit_card_full",
        placeholder="[CREDIT_CARD]",
        regex=re.compile(
            r"\b(?:\d[ \-]?){13,16}\b"
        ),
    ),
    _Pattern(
        flag="pii.credit_card_last4",
        placeholder="****-****-****-{last4}",
        regex=re.compile(
            r"\b\d{4}[-\s]?\d{4}[-\s]?\d{4}[-\s]?(\d{4})\b"
        ),
    ),
]


def redact_pii(text: str) -> tuple[str, list[str]]:
    """Replace PII patterns in text with safe placeholders.

    Returns:
        (redacted_text, detected_types) — detected_types is a list of flag
        strings (e.g. ["pii.email", "pii.phone"]) for every category that
        matched at least once.
    """
    detected: list[str] = []

    for pattern in _PATTERNS:
        if pattern.flag == "pii.credit_card_last4":
            def replacer(match):
                return pattern.placeholder.replace("{last4}", match.group(1))
            replaced, count = pattern.regex.subn(replacer, text)
        else:
            replaced, count = pattern.regex.subn(pattern.placeholder, text)
        if count > 0:
            text = replaced
            detected.append(pattern.flag)

    return text, detected
