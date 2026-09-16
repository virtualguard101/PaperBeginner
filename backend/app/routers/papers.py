from pathlib import Path

from fastapi import APIRouter, Depends, File, Form, UploadFile
from sqlalchemy.orm import Session, joinedload

from app.config import settings
from app.database import get_db
from app.deps import get_current_user
from app.http import analysis_out, fail, ok, paginated, paper_out
from app.models import Paper, PaperAnalysis, User
from app.schemas import AnalyzeRequest
from app.services.pdf import extract_text
from app.services.prompts import analyze_paper

router = APIRouter(prefix="/papers", tags=["papers"])


def _get_owned(db: Session, paper_id: str, user_id: str) -> Paper | None:
    return (
        db.query(Paper)
        .options(joinedload(Paper.category), joinedload(Paper.analyses))
        .filter(Paper.id == paper_id, Paper.user_id == user_id)
        .one_or_none()
    )


@router.post("/upload")
async def upload(
    file: UploadFile = File(...),
    title: str | None = Form(None),
    authors: str | None = Form(None),
    abstract: str | None = Form(None),
    user: User = Depends(get_current_user),
    db: Session = Depends(get_db),
):
    data = await file.read()
    if not data:
        return fail(400, "BAD_REQUEST", "File is required")
    name = file.filename or "paper.pdf"
    paper_title = title or Path(name).stem
    author_list = [a.strip() for a in (authors or "").split(",") if a.strip()]
    text = extract_text(data)
    upload_root = Path(settings.upload_dir)
    upload_root.mkdir(parents=True, exist_ok=True)
    paper = Paper(
        user_id=user.id,
        title=paper_title,
        authors=author_list,
        abstract=abstract or (text[:400] if text else ""),
        extracted_text=text,
        status="pending",
    )
    db.add(paper)
    db.commit()
    db.refresh(paper)
    dest = upload_root / f"{paper.id}.pdf"
    dest.write_bytes(data)
    paper.file_path = str(dest)
    db.commit()
    paper = _get_owned(db, paper.id, user.id)
    return ok(paper_out(paper), status_code=201)


@router.get("")
def list_papers(page: int = 1, per_page: int = 20, user: User = Depends(get_current_user), db: Session = Depends(get_db)):
    page = max(page, 1)
    per_page = min(max(per_page, 1), 100)
    q = (
        db.query(Paper)
        .options(joinedload(Paper.category))
        .filter(Paper.user_id == user.id)
        .order_by(Paper.created_at.desc())
    )
    total = q.count()
    rows = q.offset((page - 1) * per_page).limit(per_page).all()
    return paginated([paper_out(p) for p in rows], page, per_page, total)


@router.get("/{paper_id}")
def get_paper(paper_id: str, user: User = Depends(get_current_user), db: Session = Depends(get_db)):
    paper = _get_owned(db, paper_id, user.id)
    if not paper:
        return fail(404, "NOT_FOUND", "Paper not found")
    return ok(paper_out(paper, include_analyses=True))


@router.delete("/{paper_id}")
def delete_paper(paper_id: str, user: User = Depends(get_current_user), db: Session = Depends(get_db)):
    paper = db.query(Paper).filter(Paper.id == paper_id, Paper.user_id == user.id).one_or_none()
    if not paper:
        return fail(404, "NOT_FOUND", "Paper not found")
    db.delete(paper)
    db.commit()
    return ok({"ok": True})


@router.post("/{paper_id}/analyze")
def analyze(paper_id: str, body: AnalyzeRequest | None = None, user: User = Depends(get_current_user), db: Session = Depends(get_db)):
    paper = _get_owned(db, paper_id, user.id)
    if not paper:
        return fail(404, "NOT_FOUND", "Paper not found")
    paper.status = "processing"
    db.commit()
    result = analyze_paper(paper.title, paper.extracted_text or paper.abstract)
    if result.data.get("authors") and not paper.authors:
        paper.authors = result.data["authors"]
    if result.data.get("abstract"):
        paper.abstract = result.data["abstract"]
    db.query(PaperAnalysis).filter(PaperAnalysis.paper_id == paper.id).delete()
    for key in ("summary", "methodology", "contributions"):
        db.add(
            PaperAnalysis(
                paper_id=paper.id,
                analysis_type=key,
                content=str(result.data.get(key) or ""),
                llm_provider=result.provider,
                llm_model=result.model,
            )
        )
    paper.status = "completed"
    db.commit()
    paper = _get_owned(db, paper_id, user.id)
    return ok(paper_out(paper, include_analyses=True))


@router.get("/{paper_id}/analyses")
def get_analyses(paper_id: str, user: User = Depends(get_current_user), db: Session = Depends(get_db)):
    paper = _get_owned(db, paper_id, user.id)
    if not paper:
        return fail(404, "NOT_FOUND", "Paper not found")
    return ok([analysis_out(a) for a in paper.analyses])
