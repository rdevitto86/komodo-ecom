import re
import unicodedata

_MAX_BYTES = 8192

# Control characters to strip: 0x00–0x08, 0x0b–0x0c, 0x0e–0x1f, 0x7f
_CONTROL_CHARS_RE = re.compile(
    r"[\x00-\x08\x0b\x0c\x0e-\x1f\x7f]"
)

_WHITESPACE_RE = re.compile(r"[ \t\r\n]+")


def normalize(text: str) -> str:
    """Normalize input text before any guardrail check.

    Steps:
    1. NFC unicode normalization
    2. Strip dangerous control characters
    3. Collapse whitespace (tabs, newlines, multiple spaces → single space)
    4. Truncate to MAX_BYTES (8192) encoded as UTF-8
    """
    text = unicodedata.normalize("NFC", text)
    text = _CONTROL_CHARS_RE.sub("", text)
    text = _WHITESPACE_RE.sub(" ", text).strip()

    encoded = text.encode("utf-8")
    if len(encoded) > _MAX_BYTES:
        encoded = encoded[:_MAX_BYTES]
        # Decode safely, ignoring any partial multi-byte char at the boundary
        text = encoded.decode("utf-8", errors="ignore")

    return text
