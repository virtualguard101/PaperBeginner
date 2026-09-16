from datetime import datetime, timezone
from uuid import uuid4

from sqlalchemy import Boolean, DateTime, Float, ForeignKey, Integer, String, Text
from sqlalchemy.orm import Mapped, mapped_column, relationship
from sqlalchemy.types import JSON

from app.database import Base


def utcnow() -> datetime:
    return datetime.now(timezone.utc)


def new_id() -> str:
    return str(uuid4())


class User(Base):
    __tablename__ = "users"

    id: Mapped[str] = mapped_column(String, primary_key=True, default=new_id)
    email: Mapped[str] = mapped_column(String, unique=True, index=True)
    password_hash: Mapped[str] = mapped_column(String)
    name: Mapped[str] = mapped_column(String)
    avatar: Mapped[str | None] = mapped_column(String, nullable=True)
    role: Mapped[str] = mapped_column(String, default="user")
    is_active: Mapped[bool] = mapped_column(Boolean, default=True)
    preferences: Mapped[dict | None] = mapped_column(JSON, nullable=True)
    created_at: Mapped[datetime] = mapped_column(DateTime, default=utcnow)
    updated_at: Mapped[datetime] = mapped_column(DateTime, default=utcnow, onupdate=utcnow)
    last_login_at: Mapped[datetime | None] = mapped_column(DateTime, nullable=True)


class CCFCategory(Base):
    __tablename__ = "ccf_categories"

    id: Mapped[int] = mapped_column(Integer, primary_key=True)
    name: Mapped[str] = mapped_column(String)
    name_en: Mapped[str] = mapped_column(String)


class Paper(Base):
    __tablename__ = "papers"

    id: Mapped[str] = mapped_column(String, primary_key=True, default=new_id)
    user_id: Mapped[str] = mapped_column(String, ForeignKey("users.id"), index=True)
    title: Mapped[str] = mapped_column(String)
    authors: Mapped[list] = mapped_column(JSON, default=list)
    abstract: Mapped[str] = mapped_column(Text, default="")
    keywords: Mapped[list] = mapped_column(JSON, default=list)
    file_path: Mapped[str] = mapped_column(String, default="")
    extracted_text: Mapped[str] = mapped_column(Text, default="")
    status: Mapped[str] = mapped_column(String, default="pending")
    ccf_category_id: Mapped[int | None] = mapped_column(Integer, ForeignKey("ccf_categories.id"), nullable=True)
    created_at: Mapped[datetime] = mapped_column(DateTime, default=utcnow)
    updated_at: Mapped[datetime] = mapped_column(DateTime, default=utcnow, onupdate=utcnow)

    category = relationship("CCFCategory")
    analyses = relationship("PaperAnalysis", back_populates="paper", cascade="all, delete-orphan")


class PaperAnalysis(Base):
    __tablename__ = "paper_analyses"

    id: Mapped[str] = mapped_column(String, primary_key=True, default=new_id)
    paper_id: Mapped[str] = mapped_column(String, ForeignKey("papers.id"), index=True)
    analysis_type: Mapped[str] = mapped_column(String)
    content: Mapped[str] = mapped_column(Text, default="")
    structured_data: Mapped[dict | None] = mapped_column(JSON, nullable=True)
    llm_provider: Mapped[str] = mapped_column(String, default="")
    llm_model: Mapped[str] = mapped_column(String, default="")
    tokens_used: Mapped[int] = mapped_column(Integer, default=0)
    created_at: Mapped[datetime] = mapped_column(DateTime, default=utcnow)

    paper = relationship("Paper", back_populates="analyses")


class TrendingItem(Base):
    __tablename__ = "trending_items"

    id: Mapped[str] = mapped_column(String, primary_key=True, default=new_id)
    source: Mapped[str] = mapped_column(String, index=True)
    source_id: Mapped[str] = mapped_column(String, index=True)
    title: Mapped[str] = mapped_column(String)
    description: Mapped[str] = mapped_column(Text, default="")
    url: Mapped[str] = mapped_column(String, default="")
    stars: Mapped[int | None] = mapped_column(Integer, nullable=True)
    forks: Mapped[int | None] = mapped_column(Integer, nullable=True)
    language: Mapped[str | None] = mapped_column(String, nullable=True)
    topics: Mapped[list] = mapped_column(JSON, default=list)
    ccf_category_id: Mapped[int | None] = mapped_column(Integer, ForeignKey("ccf_categories.id"), nullable=True)
    trend_score: Mapped[float] = mapped_column(Float, default=0)
    crawled_at: Mapped[datetime] = mapped_column(DateTime, default=utcnow)

    category = relationship("CCFCategory")


class TrendReport(Base):
    __tablename__ = "trend_reports"

    id: Mapped[str] = mapped_column(String, primary_key=True, default=new_id)
    period: Mapped[str] = mapped_column(String)
    category_id: Mapped[int | None] = mapped_column(Integer, ForeignKey("ccf_categories.id"), nullable=True)
    summary: Mapped[str] = mapped_column(Text, default="")
    highlights: Mapped[list] = mapped_column(JSON, default=list)
    generated_by: Mapped[str] = mapped_column(String, default="")
    created_at: Mapped[datetime] = mapped_column(DateTime, default=utcnow)

    category = relationship("CCFCategory")


class LearningPath(Base):
    __tablename__ = "learning_paths"

    id: Mapped[str] = mapped_column(String, primary_key=True, default=new_id)
    user_id: Mapped[str] = mapped_column(String, ForeignKey("users.id"), index=True)
    category_id: Mapped[int | None] = mapped_column(Integer, ForeignKey("ccf_categories.id"), nullable=True)
    title: Mapped[str] = mapped_column(String)
    description: Mapped[str] = mapped_column(Text, default="")
    stages: Mapped[list] = mapped_column(JSON, default=list)
    prerequisites: Mapped[list] = mapped_column(JSON, default=list)
    estimated_time: Mapped[str] = mapped_column(String, default="")
    difficulty: Mapped[str] = mapped_column(String, default="beginner")
    generated_by: Mapped[str] = mapped_column(String, default="")
    created_at: Mapped[datetime] = mapped_column(DateTime, default=utcnow)
    updated_at: Mapped[datetime] = mapped_column(DateTime, default=utcnow, onupdate=utcnow)

    category = relationship("CCFCategory")


class Review(Base):
    __tablename__ = "reviews"

    id: Mapped[str] = mapped_column(String, primary_key=True, default=new_id)
    user_id: Mapped[str] = mapped_column(String, ForeignKey("users.id"), index=True)
    title: Mapped[str] = mapped_column(String)
    abstract: Mapped[str] = mapped_column(Text, default="")
    content: Mapped[str] = mapped_column(Text, default="")
    paper_ids: Mapped[list] = mapped_column(JSON, default=list)
    category_id: Mapped[int | None] = mapped_column(Integer, ForeignKey("ccf_categories.id"), nullable=True)
    status: Mapped[str] = mapped_column(String, default="generated")
    score: Mapped[dict | None] = mapped_column(JSON, nullable=True)
    generated_by: Mapped[str] = mapped_column(String, default="")
    created_at: Mapped[datetime] = mapped_column(DateTime, default=utcnow)
    updated_at: Mapped[datetime] = mapped_column(DateTime, default=utcnow, onupdate=utcnow)

    category = relationship("CCFCategory")
