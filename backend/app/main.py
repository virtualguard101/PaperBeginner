from pathlib import Path

from fastapi import FastAPI, HTTPException, Request
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse

from app import models  # noqa: F401
from app.config import settings
from app.database import Base, SessionLocal, engine
from app.routers import auth, ccf, learning, papers, reviews, trending, users
from app.seed import seed_all


def _ensure_dirs() -> None:
    Path(settings.upload_dir).mkdir(parents=True, exist_ok=True)
    if settings.database_url.startswith("sqlite:///"):
        raw = settings.database_url.removeprefix("sqlite:///")
        Path(raw).parent.mkdir(parents=True, exist_ok=True)


def create_app() -> FastAPI:
    _ensure_dirs()
    Base.metadata.create_all(bind=engine)
    db = SessionLocal()
    try:
        seed_all(db)
    finally:
        db.close()

    app = FastAPI(title="PaperBeginner API", version="1.0.0")
    app.add_middleware(
        CORSMiddleware,
        allow_origins=["*"],
        allow_credentials=True,
        allow_methods=["*"],
        allow_headers=["*"],
    )

    @app.exception_handler(HTTPException)
    async def http_exc_handler(_request: Request, exc: HTTPException):
        code = "UNAUTHORIZED" if exc.status_code == 401 else "ERROR"
        detail = exc.detail if isinstance(exc.detail, str) else "Request failed"
        return JSONResponse(
            status_code=exc.status_code,
            content={"success": False, "error": {"code": code, "message": detail}},
        )

    @app.get("/health")
    def health():
        return {"status": "healthy", "version": "1.0.0"}

    prefix = "/api/v1"
    app.include_router(auth.router, prefix=prefix)
    app.include_router(users.router, prefix=prefix)
    app.include_router(papers.router, prefix=prefix)
    app.include_router(trending.router, prefix=prefix)
    app.include_router(learning.router, prefix=prefix)
    app.include_router(reviews.router, prefix=prefix)
    app.include_router(ccf.router, prefix=prefix)
    return app


app = create_app()
