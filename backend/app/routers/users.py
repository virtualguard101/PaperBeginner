from fastapi import APIRouter, Depends
from sqlalchemy.orm import Session

from app.database import get_db
from app.deps import get_current_user, hash_password, verify_password
from app.http import fail, ok, user_out
from app.models import User
from app.schemas import ChangePasswordRequest, UpdateUserRequest

router = APIRouter(prefix="/users", tags=["users"])


@router.get("/me")
def me(user: User = Depends(get_current_user)):
    return ok(user_out(user))


@router.put("/me")
def update_me(body: UpdateUserRequest, user: User = Depends(get_current_user), db: Session = Depends(get_db)):
    if body.name:
        user.name = body.name
    if body.avatar is not None:
        user.avatar = body.avatar
    if body.preferences is not None:
        user.preferences = body.preferences
    db.commit()
    db.refresh(user)
    return ok(user_out(user))


@router.put("/me/password")
def change_password(body: ChangePasswordRequest, user: User = Depends(get_current_user), db: Session = Depends(get_db)):
    if not verify_password(body.old_password, user.password_hash):
        return fail(400, "BAD_REQUEST", "Current password is incorrect")
    user.password_hash = hash_password(body.new_password)
    db.commit()
    return ok({"ok": True})
