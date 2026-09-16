from fastapi import APIRouter, Depends
from sqlalchemy.orm import Session

from app.database import get_db
from app.http import ok

router = APIRouter(prefix="/ccf", tags=["ccf"])


@router.get("/categories")
def categories(db: Session = Depends(get_db)):
    from app.models import CCFCategory

    rows = db.query(CCFCategory).order_by(CCFCategory.id).all()
    return ok([{"id": r.id, "name": r.name, "name_en": r.name_en} for r in rows])
