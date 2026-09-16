from datetime import datetime, timezone

from fastapi import APIRouter, Depends
from sqlalchemy.orm import Session, joinedload

from app.database import get_db
from app.deps import get_current_user
from app.http import fail, ok, paginated, review_out
from app.models import Paper, PaperAnalysis, Review, User
from app.schemas import ReviewGenerateRequest, ReviewScoreRequest
from app.services.prompts import generate_review, score_review

router = APIRouter(prefix="/reviews", tags=["reviews"])


def _owned(db: Session, review_id: str, user_id: str) -> Review | None:
    return (
        db.query(Review)
        .options(joinedload(Review.category))
        .filter(Review.id == review_id, Review.user_id == user_id)
        .one_or_none()
    )


@router.post("/generate")
def generate(body: ReviewGenerateRequest, user: User = Depends(get_current_user), db: Session = Depends(get_db)):
    if not body.paper_ids:
        return fail(400, "BAD_REQUEST", "请至少选择一篇论文")
    papers = (
        db.query(Paper)
        .filter(Paper.user_id == user.id, Paper.id.in_(body.paper_ids))
        .all()
    )
    if not papers:
        return fail(400, "BAD_REQUEST", "未找到所选论文")
    blobs = []
    for p in papers:
        analyses = db.query(PaperAnalysis).filter(PaperAnalysis.paper_id == p.id).all()
        analysis_text = "\n".join(f"{a.analysis_type}: {a.content}" for a in analyses)
        blobs.append(f"## {p.title}\nAuthors: {p.authors}\nAbstract: {p.abstract}\n{analysis_text}")
    result = generate_review(body.title, body.style or "academic", "\n\n".join(blobs))
    review = Review(
        user_id=user.id,
        title=body.title,
        abstract=str(result.data.get("abstract") or ""),
        content=str(result.data.get("content") or ""),
        paper_ids=body.paper_ids,
        category_id=body.category_id,
        status="generated",
        generated_by=result.model,
    )
    db.add(review)
    db.commit()
    db.refresh(review)
    review = _owned(db, review.id, user.id)
    return ok(review_out(review), status_code=201)


@router.get("")
def list_reviews(page: int = 1, per_page: int = 20, user: User = Depends(get_current_user), db: Session = Depends(get_db)):
    page = max(page, 1)
    per_page = min(max(per_page, 1), 100)
    q = (
        db.query(Review)
        .options(joinedload(Review.category))
        .filter(Review.user_id == user.id)
        .order_by(Review.created_at.desc())
    )
    total = q.count()
    rows = q.offset((page - 1) * per_page).limit(per_page).all()
    return paginated([review_out(r) for r in rows], page, per_page, total)


@router.get("/{review_id}")
def get_review(review_id: str, user: User = Depends(get_current_user), db: Session = Depends(get_db)):
    review = _owned(db, review_id, user.id)
    if not review:
        return fail(404, "NOT_FOUND", "Review not found")
    return ok(review_out(review))


@router.post("/{review_id}/score")
def score(review_id: str, body: ReviewScoreRequest | None = None, user: User = Depends(get_current_user), db: Session = Depends(get_db)):
    review = _owned(db, review_id, user.id)
    if not review:
        return fail(404, "NOT_FOUND", "Review not found")
    content = (body.content if body and body.content else review.content) or ""
    result = score_review(review.title, content)
    payload = dict(result.data)
    payload["generated_at"] = datetime.now(timezone.utc).isoformat()
    review.score = payload
    review.status = "scored"
    db.commit()
    review = _owned(db, review_id, user.id)
    return ok(review_out(review))


@router.delete("/{review_id}")
def delete_review(review_id: str, user: User = Depends(get_current_user), db: Session = Depends(get_db)):
    review = db.query(Review).filter(Review.id == review_id, Review.user_id == user.id).one_or_none()
    if not review:
        return fail(404, "NOT_FOUND", "Review not found")
    db.delete(review)
    db.commit()
    return ok({"ok": True})
