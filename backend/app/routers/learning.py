from fastapi import APIRouter, Depends
from sqlalchemy.orm import Session, joinedload

from app.database import get_db
from app.deps import get_current_user
from app.http import fail, ok, path_out
from app.models import CCFCategory, LearningPath, User
from app.schemas import LearningGenerateRequest
from app.services.prompts import generate_learning_path

router = APIRouter(prefix="/learning", tags=["learning"])


@router.post("/generate")
def generate(body: LearningGenerateRequest, user: User = Depends(get_current_user), db: Session = Depends(get_db)):
    cat = db.get(CCFCategory, body.category_id)
    if not cat:
        return fail(400, "BAD_REQUEST", "Unknown category")
    extras = ""
    if body.prerequisites:
        extras += "prerequisites: " + ", ".join(body.prerequisites)
    if body.focus_areas:
        extras += " focus: " + ", ".join(body.focus_areas)
    result = generate_learning_path(cat.name, body.difficulty, extras)
    path = LearningPath(
        user_id=user.id,
        category_id=cat.id,
        title=str(result.data.get("title") or f"{cat.name} 学习路线"),
        description=str(result.data.get("description") or ""),
        stages=result.data.get("stages") or [],
        prerequisites=result.data.get("prerequisites") or (body.prerequisites or []),
        estimated_time=str(result.data.get("estimated_time") or ""),
        difficulty=body.difficulty,
        generated_by=result.model,
    )
    db.add(path)
    db.commit()
    db.refresh(path)
    path = (
        db.query(LearningPath)
        .options(joinedload(LearningPath.category))
        .filter(LearningPath.id == path.id)
        .one()
    )
    return ok(path_out(path), status_code=201)


@router.get("/paths")
def list_paths(user: User = Depends(get_current_user), db: Session = Depends(get_db)):
    rows = (
        db.query(LearningPath)
        .options(joinedload(LearningPath.category))
        .filter(LearningPath.user_id == user.id)
        .order_by(LearningPath.created_at.desc())
        .all()
    )
    return ok([path_out(p) for p in rows])


@router.get("/paths/{path_id}")
def get_path(path_id: str, user: User = Depends(get_current_user), db: Session = Depends(get_db)):
    path = (
        db.query(LearningPath)
        .options(joinedload(LearningPath.category))
        .filter(LearningPath.id == path_id, LearningPath.user_id == user.id)
        .one_or_none()
    )
    if not path:
        return fail(404, "NOT_FOUND", "Learning path not found")
    return ok(path_out(path))
