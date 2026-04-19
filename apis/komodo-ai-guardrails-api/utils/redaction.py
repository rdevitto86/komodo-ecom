from presidio_analyzer import AnalyzerEngine
from presidio_analyzer.nlp_engine import NlpEngineProvider
from presidio_anonymizer import AnonymizerEngine
from presidio_anonymizer.entities import OperatorConfig

# Entity types → (flag, redaction placeholder)
_ENTITY_CONFIG: dict[str, tuple[str, str]] = {
    "EMAIL_ADDRESS": ("pii.email",        "[EMAIL]"),
    "PHONE_NUMBER":  ("pii.phone",        "[PHONE]"),
    "US_SSN":        ("pii.ssn",          "[SSN]"),
    "CREDIT_CARD":   ("pii.credit_card",  "[CREDIT_CARD]"),
}

_ENTITIES = list(_ENTITY_CONFIG.keys())

_OPERATORS = {
    entity: OperatorConfig("replace", {"new_value": placeholder})
    for entity, (_, placeholder) in _ENTITY_CONFIG.items()
}

# Initialized once at startup. en_core_web_sm is ~12 MB and is a statistical
# NLP model (not an LLM) — EMAIL/PHONE/SSN/CREDIT_CARD are all pattern-based
# with checksum validation; no NER inference runs for these entity types.
_nlp_engine = NlpEngineProvider(nlp_configuration={
    "nlp_engine_name": "spacy",
    "models": [{"lang_code": "en", "model_name": "en_core_web_sm"}],
}).create_engine()

_analyzer = AnalyzerEngine(nlp_engine=_nlp_engine, supported_languages=["en"])
_anonymizer = AnonymizerEngine()


def redact_pii(text: str) -> tuple[str, list[str]]:
    """Detect and redact PII using Presidio pattern recognizers.

    Credit cards are validated via Luhn checksum. Returns
    (redacted_text, detected_flags).
    """
    results = _analyzer.analyze(text=text, entities=_ENTITIES, language="en")
    if not results:
        return text, []

    detected = list({_ENTITY_CONFIG[r.entity_type][0] for r in results})
    anonymized = _anonymizer.anonymize(
        text=text, analyzer_results=results, operators=_OPERATORS
    )
    return anonymized.text, detected
