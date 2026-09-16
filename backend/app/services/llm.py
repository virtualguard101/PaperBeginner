import json
import re
from typing import Any

import httpx

from app.config import settings


class LLMResult:
    def __init__(self, data: dict[str, Any], provider: str, model: str, fallback: bool = False, raw: str = ""):
        self.data = data
        self.provider = provider
        self.model = model
        self.fallback = fallback
        self.raw = raw


def _extract_json(text: str) -> dict[str, Any] | None:
    text = text.strip()
    try:
        return json.loads(text)
    except json.JSONDecodeError:
        pass
    match = re.search(r"\{[\s\S]*\}", text)
    if match:
        try:
            return json.loads(match.group(0))
        except json.JSONDecodeError:
            return None
    return None


def chat_json(system: str, user: str, fallback: dict[str, Any]) -> LLMResult:
    if not settings.llm_api_key:
        return LLMResult(fallback, "fallback", "template", fallback=True)

    url = settings.llm_base_url.rstrip("/") + "/chat/completions"
    payload = {
        "model": settings.llm_model,
        "messages": [
            {"role": "system", "content": system},
            {"role": "user", "content": user},
        ],
        "temperature": 0.4,
    }
    headers = {
        "Authorization": f"Bearer {settings.llm_api_key}",
        "Content-Type": "application/json",
    }
    try:
        with httpx.Client(timeout=settings.llm_timeout_seconds) as client:
            resp = client.post(url, json=payload, headers=headers)
            resp.raise_for_status()
            body = resp.json()
            content = body["choices"][0]["message"]["content"]
            parsed = _extract_json(content)
            if not parsed:
                return LLMResult(fallback, "fallback", "parse-error", fallback=True, raw=content)
            return LLMResult(parsed, "openai-compatible", settings.llm_model, raw=content)
    except Exception:
        return LLMResult(fallback, "fallback", "error", fallback=True)
