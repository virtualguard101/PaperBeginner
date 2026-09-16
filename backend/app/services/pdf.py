from io import BytesIO

from pypdf import PdfReader


def extract_text(data: bytes, max_chars: int = 12000) -> str:
    try:
        reader = PdfReader(BytesIO(data))
        parts: list[str] = []
        for page in reader.pages:
            text = page.extract_text() or ""
            parts.append(text)
        joined = "\n".join(parts).strip()
        if len(joined) > max_chars:
            return joined[:max_chars]
        return joined
    except Exception:
        return ""
