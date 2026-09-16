from datetime import datetime, timezone
from math import ceil
from typing import Any

from fastapi.responses import JSONResponse


def iso(dt: datetime | None) -> str | None:
    if dt is None:
        return None
    if dt.tzinfo is None:
        dt = dt.replace(tzinfo=timezone.utc)
    return dt.isoformat()


def ok(data: Any = None, status_code: int = 200, meta: dict | None = None) -> JSONResponse:
    body: dict[str, Any] = {"success": True, "data": data}
    if meta:
        body["meta"] = meta
    return JSONResponse(content=body, status_code=status_code)


def fail(status_code: int, code: str, message: str, details: str | None = None) -> JSONResponse:
    err: dict[str, Any] = {"code": code, "message": message}
    if details:
        err["details"] = details
    return JSONResponse(content={"success": False, "error": err}, status_code=status_code)


def paginated(data: Any, page: int, per_page: int, total: int) -> JSONResponse:
    total_pages = ceil(total / per_page) if per_page else 0
    return ok(
        data,
        meta={"page": page, "per_page": per_page, "total": total, "total_pages": total_pages},
    )


def category_out(cat) -> dict | None:
    if cat is None:
        return None
    return {"id": cat.id, "name": cat.name, "name_en": cat.name_en}


def user_out(user) -> dict:
    return {
        "id": user.id,
        "email": user.email,
        "name": user.name,
        "avatar": user.avatar or "",
        "role": user.role,
        "is_active": user.is_active,
        "preferences": user.preferences,
        "created_at": iso(user.created_at),
        "last_login_at": iso(user.last_login_at),
    }


def paper_out(paper, include_analyses: bool = False) -> dict:
    payload = {
        "id": paper.id,
        "title": paper.title,
        "authors": paper.authors or [],
        "abstract": paper.abstract or "",
        "keywords": paper.keywords or [],
        "ccf_category": category_out(paper.category),
        "status": paper.status,
        "created_at": iso(paper.created_at),
    }
    if include_analyses:
        payload["analyses"] = [analysis_out(a) for a in (paper.analyses or [])]
    return payload


def analysis_out(a) -> dict:
    return {
        "id": a.id,
        "paper_id": a.paper_id,
        "analysis_type": a.analysis_type,
        "content": a.content,
        "structured_data": a.structured_data,
        "llm_provider": a.llm_provider,
        "llm_model": a.llm_model,
        "tokens_used": a.tokens_used,
        "created_at": iso(a.created_at),
    }


def trending_out(item) -> dict:
    return {
        "id": item.id,
        "source": item.source,
        "source_id": item.source_id,
        "title": item.title,
        "description": item.description,
        "url": item.url,
        "stars": item.stars,
        "forks": item.forks,
        "language": item.language,
        "topics": item.topics or [],
        "ccf_category": category_out(item.category),
        "trend_score": item.trend_score,
        "crawled_at": iso(item.crawled_at),
    }


def path_out(path) -> dict:
    return {
        "id": path.id,
        "title": path.title,
        "description": path.description,
        "category": category_out(path.category),
        "stages": path.stages or [],
        "prerequisites": path.prerequisites or [],
        "estimated_time": path.estimated_time,
        "difficulty": path.difficulty,
        "created_at": iso(path.created_at),
    }


def review_out(review) -> dict:
    score = review.score or None
    return {
        "id": review.id,
        "title": review.title,
        "abstract": review.abstract,
        "content": review.content,
        "paper_ids": review.paper_ids or [],
        "category": category_out(review.category),
        "status": review.status,
        "score": score,
        "generated_by": review.generated_by,
        "created_at": iso(review.created_at),
        "updated_at": iso(review.updated_at),
    }
