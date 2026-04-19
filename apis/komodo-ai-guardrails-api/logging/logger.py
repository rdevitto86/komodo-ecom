import logging
import sys
from datetime import datetime, timezone

from config.consts import APP_ENV

class _KomodoFormatter(logging.Formatter):
    """Formats records as: timestamp [LEVEL] - message | key=value ..."""

    def format(self, record: logging.LogRecord) -> str:
        ts = (
            datetime.fromtimestamp(record.created, tz=timezone.utc)
            .strftime("%Y-%m-%dT%H:%M:%S.%f")[:-3] + "Z"
        )
        msg = record.getMessage()
        base = f"{ts} [{record.levelname}] - {msg}"

        details: dict | None = getattr(record, "details", None)
        if details:
            detail_str = " ".join(f"{k}={v}" for k, v in details.items())
            return f"{base} | {detail_str}"
        return base


def _make_level() -> int:
    return logging.INFO if APP_ENV in ("dev", "qa") else logging.ERROR

_handler = logging.StreamHandler(sys.stdout)
_handler.setFormatter(_KomodoFormatter())

_root = logging.getLogger("komodo") 
_root.setLevel(_make_level())
_root.addHandler(_handler)
_root.propagate = False


class KomodoLogger:
    """Thin wrapper that forwards **kwargs as structured detail fields."""

    def __init__(self, log: logging.Logger) -> None:
        self._log = log

    def info(self, msg: str, **details) -> None:
        if self._log.isEnabledFor(logging.INFO):
            self._log.info(msg, extra={"details": details} if details else {})

    def warn(self, msg: str, **details) -> None:
        if self._log.isEnabledFor(logging.WARNING):
            self._log.warning(msg, extra={"details": details} if details else {})

    def error(self, msg: str, **details) -> None:
        self._log.error(msg, extra={"details": details} if details else {})


def get_logger(name: str) -> KomodoLogger:
    """Return a KomodoLogger scoped to *name* under the komodo root logger."""
    return KomodoLogger(_root.getChild(name))
