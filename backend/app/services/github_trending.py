from __future__ import annotations

import json
from pathlib import Path

import httpx
from sqlalchemy.orm import Session

from app.config import settings
from app.models import TrendingItem, utcnow
from app.services.classify import classify_text

SEED_DIR = Path(__file__).resolve().parent.parent / "seed"


def _fallback_github() -> list[dict]:
    path = SEED_DIR / "github_fallback.json"
    if path.exists():
        return json.loads(path.read_text(encoding="utf-8"))
    return []


def fetch_github_trending() -> list[dict]:
    headers = {"Accept": "application/vnd.github+json", "User-Agent": "PaperBeginner"}
    if settings.github_token:
        headers["Authorization"] = f"Bearer {settings.github_token}"
    try:
        with httpx.Client(timeout=20.0) as client:
            resp = client.get(
                "https://api.github.com/search/repositories",
                params={"q": "stars:>500", "sort": "updated", "order": "desc", "per_page": 15},
                headers=headers,
            )
            resp.raise_for_status()
            items = resp.json().get("items", [])
            out = []
            for repo in items:
                out.append(
                    {
                        "source_id": str(repo.get("id")),
                        "title": repo.get("full_name") or repo.get("name"),
                        "description": repo.get("description") or "",
                        "url": repo.get("html_url") or "",
                        "stars": repo.get("stargazers_count") or 0,
                        "forks": repo.get("forks_count") or 0,
                        "language": repo.get("language"),
                        "topics": repo.get("topics") or [],
                    }
                )
            if out:
                return out
    except Exception:
        pass
    return _fallback_github()


def upsert_github(db: Session) -> int:
    rows = fetch_github_trending()
    count = 0
    now = utcnow()
    for row in rows:
        existing = (
            db.query(TrendingItem)
            .filter(TrendingItem.source == "github", TrendingItem.source_id == row["source_id"])
            .one_or_none()
        )
        cat_id = classify_text(row.get("title", ""), row.get("description", ""), " ".join(row.get("topics") or []))
        stars = int(row.get("stars") or 0)
        score = min(99.9, round(stars / 1000 + 50, 1))
        if existing:
            existing.title = row["title"]
            existing.description = row.get("description") or ""
            existing.url = row.get("url") or ""
            existing.stars = stars
            existing.forks = int(row.get("forks") or 0)
            existing.language = row.get("language")
            existing.topics = row.get("topics") or []
            existing.ccf_category_id = cat_id
            existing.trend_score = score
            existing.crawled_at = now
        else:
            db.add(
                TrendingItem(
                    source="github",
                    source_id=row["source_id"],
                    title=row["title"],
                    description=row.get("description") or "",
                    url=row.get("url") or "",
                    stars=stars,
                    forks=int(row.get("forks") or 0),
                    language=row.get("language"),
                    topics=row.get("topics") or [],
                    ccf_category_id=cat_id,
                    trend_score=score,
                    crawled_at=now,
                )
            )
        count += 1
    db.commit()
    return count
