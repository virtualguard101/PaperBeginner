from fastapi import APIRouter, Depends
from sqlalchemy.orm import Session

from app.database import get_db
from app.deps import create_token, get_current_user, hash_password, verify_password
from app.http import fail, ok, user_out
from app.models import User, utcnow
from app.schemas import LoginRequest, RegisterRequest

router = APIRouter(prefix="/auth", tags=["auth"])


@router.post("/register")
def register(body: RegisterRequest, db: Session = Depends(get_db)):
    exists = db.query(User).filter(User.email == body.email.lower()).one_or_none()
    if exists:
        return fail(409, "CONFLICT", "User with this email already exists")
    user = User(
        email=body.email.lower(),
        password_hash=hash_password(body.password),
        name=body.name,
    )
    db.add(user)
    db.commit()
    db.refresh(user)
    return ok(user_out(user), status_code=201)


@router.post("/login")
def login(body: LoginRequest, db: Session = Depends(get_db)):
    user = db.query(User).filter(User.email == body.email.lower()).one_or_none()
    if user is None or not verify_password(body.password, user.password_hash):
        return fail(401, "UNAUTHORIZED", "Invalid email or password")
    user.last_login_at = utcnow()
    db.commit()
    token, exp = create_token(user.id)
    return ok({"token": token, "expires_at": exp, "user": user_out(user)})


@router.post("/refresh")
def refresh(user: User = Depends(get_current_user)):
    token, exp = create_token(user.id)
    return ok({"token": token, "expires_at": exp, "user": user_out(user)})
