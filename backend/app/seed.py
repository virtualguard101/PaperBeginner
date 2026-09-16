import json
from pathlib import Path

from sqlalchemy.orm import Session

from app.models import CCFCategory, TrendingItem, utcnow
from app.services.classify import classify_text

SEED_DIR = Path(__file__).resolve().parent / "seed"


def seed_all(db: Session) -> None:
    cats = json.loads((SEED_DIR / "ccf_categories.json").read_text(encoding="utf-8"))
    for row in cats:
        if db.get(CCFCategory, row["id"]) is None:
            db.add(CCFCategory(id=row["id"], name=row["name"], name_en=row["name_en"]))
    db.commit()

    if db.query(TrendingItem).filter(TrendingItem.source == "ccf").count() == 0:
        items = json.loads((SEED_DIR / "ccf_items.json").read_text(encoding="utf-8"))
        now = utcnow()
        for row in items:
            db.add(
                TrendingItem(
                    source="ccf",
                    source_id=row["source_id"],
                    title=row["title"],
                    description=row.get("description") or "",
                    url=row.get("url") or "",
                    topics=row.get("topics") or [],
                    ccf_category_id=row.get("ccf_category_id") or classify_text(row["title"], row.get("description", "")),
                    trend_score=row.get("trend_score") or 80,
                    crawled_at=now,
                )
            )
        db.commit()
