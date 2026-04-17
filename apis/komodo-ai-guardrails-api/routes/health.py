from fastapi import APIRouter
from config.consts import SERVICE_NAME, VERSION, GUARDRAILS_PROVIDER

router = APIRouter()

@router.get("/health")
async def health() -> dict:
    return {
        "status": "ok",
        "service": SERVICE_NAME,
        "version": VERSION,
        "provider": GUARDRAILS_PROVIDER,
    }
