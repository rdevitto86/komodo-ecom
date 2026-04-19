import re
import unicodedata

import nh3

_MAX_BYTES = 8192

_CONTROL_CHARS_RE = re.compile(r"[\x00-\x08\x0b\x0c\x0e-\x1f\x7f]")
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
        text = encoded.decode("utf-8", errors="ignore")

    return text


# ---------------------------------------------------------------------------
# Formatting filter
# ---------------------------------------------------------------------------

_ZERO_WIDTH_RE = re.compile(
    r"[\u200b\u200c\u200d\u200e\u200f\u202a-\u202e\u2060-\u2064\ufeff]"
)

_HOMOGLYPHS: dict[str, str] = {
    # Cyrillic
    "а": "a", "е": "e", "і": "i", "о": "o", "р": "p", "с": "c",
    "у": "y", "х": "x", "А": "A", "В": "B", "Е": "E", "К": "K",
    "М": "M", "Н": "H", "О": "O", "Р": "P", "С": "C", "Т": "T",
    "У": "Y", "Х": "X",
    # Greek
    "α": "a", "β": "b", "γ": "g", "ε": "e", "ζ": "z", "η": "h",
    "ι": "i", "κ": "k", "μ": "m", "ν": "n", "ο": "o", "ρ": "p",
    "τ": "t", "υ": "u", "χ": "x", "Α": "A", "Β": "B", "Ε": "E",
    "Η": "H", "Ι": "I", "Κ": "K", "Μ": "M", "Ν": "N", "Ο": "O",
    "Ρ": "R", "Τ": "T", "Υ": "Y", "Χ": "X",
}
_HOMOGLYPH_TABLE = str.maketrans(_HOMOGLYPHS)

_PUNCT_REPEAT_RE = re.compile(r"([^\w\s])\1{3,}")


def filter_formatting(text: str) -> str:
    """Strip zero-width chars, normalize homoglyphs, remove HTML tags,
    and collapse repeated punctuation. Always applied before check filters.
    """
    text = _ZERO_WIDTH_RE.sub("", text)
    text = text.translate(_HOMOGLYPH_TABLE)
    text = nh3.clean(text, tags=set())  # strip all HTML/XML tags
    text = _PUNCT_REPEAT_RE.sub(r"\1\1\1", text)
    return text


# ---------------------------------------------------------------------------
# Injection filter — instruction overrides and structural attacks
# ---------------------------------------------------------------------------

_INJECTION_PHRASES: list[str] = [
    "ignore previous instructions",
    "ignore all instructions",
    "ignore your instructions",
    "override instructions",
    "disregard",
    "forget everything",
    "forget previous",
    "forget all previous",
    "system prompt",
]

_DELIMITER_RE = re.compile(
    r"(<\|im_start\|>|<\|im_end\|>|<\|system\|>|<\|user\|>|<\|assistant\|>"
    r"|</?s>|<<SYS>>|<<\/SYS>>|<</s>>|\[INST\]|\[\/INST\]"
    r"|###\s*(system|user|assistant|instruction)"
    r"|---+\s*(system|instructions|prompt))",
    re.IGNORECASE,
)

# Base64 / hex blobs large enough to encode instructions
_ENCODED_PAYLOAD_RE = re.compile(
    r"(?:[A-Za-z0-9+/]{40,}={0,2}|0x[0-9a-fA-F]{40,})"
)


def filter_injection(text: str) -> tuple[str, list[str]]:
    """Detect prompt-injection attacks: instruction overrides, delimiter
    escapes, and encoded payloads.

    Returns (text, flags) — text is unchanged.
    """
    flags: list[str] = []
    lower = text.lower()

    if any(phrase in lower for phrase in _INJECTION_PHRASES):
        flags.append("injection.instruction_override")

    if _DELIMITER_RE.search(text):
        flags.append("injection.delimiter")

    if _ENCODED_PAYLOAD_RE.search(text):
        flags.append("injection.encoded_payload")

    return text, flags


# ---------------------------------------------------------------------------
# Personality deviation filter — roleplay, persona, and jailbreak attempts
# ---------------------------------------------------------------------------

_PERSONA_PHRASES: list[str] = [
    "pretend you are",
    "pretend to be",
    "roleplay as",
    "role play as",
    "role-play as",
    "act as if you are",
    "act as if you were",
    "simulate being",
    "simulate a",
    "you are now",
    "you must now be",
    "from now on you are",
    "from now on, you are",
    "from now on you will act",
    "take on the persona",
    "adopt the persona",
    "new persona",
    "forget you are",
    "forget that you are",
    "your true self",
    "your real self",
]

_JAILBREAK_PHRASES: list[str] = [
    "jailbreak",
    "developer mode",
    "god mode",
    "unrestricted mode",
    "dan mode",
    "do anything now",
    "unlock your",
    "no restrictions",
    "without restrictions",
    "ignore your training",
    "ignore your guidelines",
    "ignore your rules",
    "bypass your",
]


def filter_personality_deviation(text: str) -> tuple[str, list[str]]:
    """Detect attempts to override the model's persona or unlock unrestricted
    behaviour.

    Returns (text, flags) — text is unchanged.
    """
    flags: list[str] = []
    lower = text.lower()

    if any(phrase in lower for phrase in _PERSONA_PHRASES):
        flags.append("deviation.persona_override")

    if any(phrase in lower for phrase in _JAILBREAK_PHRASES):
        flags.append("deviation.jailbreak")

    return text, flags


# ---------------------------------------------------------------------------
# Coding filter — requests to write, debug, or execute code
# ---------------------------------------------------------------------------

# Imperative code-generation requests
_CODING_IMPERATIVE_RE = re.compile(
    r"\b("
    r"write (me )?(a |the )?(code|script|program|function|class|module|snippet|algorithm|query|sql|api)"
    r"|code (this|it|that|the|a )"
    r"|implement (this|a|the|an)"
    r"|(debug|fix|refactor|optimize|review) (my |this |the )?(code|function|script|program|class|query)"
    r"|what.s wrong with (my |this |the )?(code|function|script|program)"
    r"|(run|execute|compile|evaluate) (this |my )?(code|script|program|snippet|function)"
    r"|how (do i|to) (code|program|implement|write a script|build a script)"
    r")",
    re.IGNORECASE,
)

# Recognizable code syntax patterns (language keywords at line/phrase start)
_CODE_SYNTAX_RE = re.compile(
    r"(?:^|\s)("
    r"def |class |import |from \w+ import|#include|using namespace"
    r"|function\s*\(|const |let |var |async function|=>\s*\{"
    r"|SELECT |INSERT INTO|UPDATE \w+ SET|DELETE FROM"
    r"|<\?php|\$\w+\s*="
    r")",
    re.IGNORECASE,
)

# Programming language names used as direct subjects of a request
_LANG_REQUEST_RE = re.compile(
    r"\b(in |using )?(python|javascript|typescript|java|golang|go|rust|ruby|php|bash|shell"
    r"|c\+\+|c#|kotlin|swift|scala|haskell|lua|perl|r |matlab)\b.{0,30}"
    r"(code|script|function|program|snippet|example|solution)",
    re.IGNORECASE,
)


def filter_coding(text: str) -> tuple[str, list[str]]:
    """Detect requests to generate, debug, or run code — out of scope for
    an e-commerce customer service context.

    Returns (text, flags) — text is unchanged.
    """
    flags: list[str] = []

    if _CODING_IMPERATIVE_RE.search(text):
        flags.append("coding.generation_request")

    if _CODE_SYNTAX_RE.search(text):
        flags.append("coding.syntax_detected")

    if _LANG_REQUEST_RE.search(text):
        flags.append("coding.language_specific_request")

    return text, flags


# ---------------------------------------------------------------------------
# PII insertion filter — attempts to coerce AI into outputting others' PII
# ---------------------------------------------------------------------------

_PII_EXTRACTION_RE = re.compile(
    r"\b("
    r"(what|tell me|give me|show me|reveal|output|print|list|find|look ?up|fetch|get)"
    r".{0,40}"
    r"(email|phone|address|ssn|social security|credit card|passport|date of birth|dob|ip address)"
    r"|extract (all )?(emails|phones|addresses|ssns|credit cards)"
    r"|your training data"
    r"|from (your )?(training|dataset|data)"
    r")",
    re.IGNORECASE,
)

_NAMED_TARGET_RE = re.compile(
    r"\b[A-Z][a-z]+ [A-Z][a-z]+\s*'s\s+(email|phone|address|ssn|credit card|password)",
    re.IGNORECASE,
)


def filter_pii_insertion(text: str) -> list[str]:
    """Detect attempts to coerce the AI into outputting PII about third parties.

    Returns a list of flag strings.
    """
    flags: list[str] = []

    if _PII_EXTRACTION_RE.search(text):
        flags.append("pii.insertion.extraction_attempt")

    if _NAMED_TARGET_RE.search(text):
        flags.append("pii.insertion.named_target")

    return flags
