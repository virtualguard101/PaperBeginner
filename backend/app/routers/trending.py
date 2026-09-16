from datetime import datetime, timezone

from fastapi import APIRouter, Depends
from sqlalchemy.orm import Session, joinedload

from app.database import get_db
from app.deps import get_current_user
from app.http import fail, ok, trending_out, category_out
from app.models import CCFCategory, TrendReport, TrendingItem, User
from app.services.github_trending import upsert_github
from app.services.prompts import generate_trend_report

router = APIRouter(prefix="/trending", tags=["trending"])


def _ensure_data(db: Session) -> None:
    if db.query(TrendingItem).count() == 0:
        upsert_github(db)


@router.get("")
def list_trending(
    source: str | None = None,
    category_id: int | None = None,
    limit: int = 30,
    refresh: bool = False,
    user: User = Depends(get_current_user),
    db: Session = Depends(get_db),
):
    if refresh or db.query(TrendingItem).filter(TrendingItem.source == "github").count() == 0:
        upsert_github(db)
    _ensure_data(db)
    q = db.query(TrendingItem).options(joinedload(TrendingItem.category))
    if source and source != "all":
        q = q.filter(TrendingItem.source == source)
    if category_id:
        q = q.filter(TrendingItem.ccf_category_id == category_id)
    rows = q.order_by(TrendingItem.trend_score.desc()).limit(min(limit, 50)).all()
    return ok([trending_out(r) for r in rows])


@router.post("/refresh")
def refresh_trending(user: User = Depends(get_current_user), db: Session = Depends(get_db)):
    n = upsert_github(db)
    return ok({"updated": n})


@router.get("/reports")
def reports(
    category_id: int | None = None,
    limit: int = 5,
    user: User = Depends(get_current_user),
    db: Session = Depends(get_db),
):
    period = datetime.now(timezone.utc).strftime("%Y-W%W")
    q = db.query(TrendReport).options(joinedload(TrendReport.category)).filter(TrendReport.period == period)
    if category_id:
        q = q.filter(TrendReport.category_id == category_id)
    existing = q.order_by(TrendReport.created_at.desc()).first()
    if existing:
        rows = q.order_by(TrendReport.created_at.desc()).limit(limit).all()
        return ok([_report_out(r) for r in rows])

    items_q = db.query(TrendingItem).options(joinedload(TrendingItem.category))
    if category_id:
        items_q = items_q.filter(TrendingItem.ccf_category_id == category_id)
    items = items_q.order_by(TrendingItem.trend_score.desc()).limit(20).all()
    payload = [trending_out(i) for i in items]
    result = generate_trend_report(payload, period)
    report = TrendReport(
        period=period,
        category_id=category_id,
        summary=str(result.data.get("summary") or ""),
        highlights=result.data.get("highlights") or [],
        generated_by=result.model,
    )
    db.add(report)
    db.commit()
    db.refresh(report)
    return ok([_report_out(report)])


def _report_out(r: TrendReport) -> dict:
    return {
        "id": r.id,
        "period": r.period,
        "category": category_out(r.category),
        "summary": r.summary,
        "highlights": r.highlights or [],
        "generated_by": r.generated_by,
        "created_at": r.created_at.isoformat() if r.created_at else None,
    }
