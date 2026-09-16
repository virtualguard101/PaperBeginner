from typing import Any, Generic, TypeVar

from pydantic import BaseModel, Field

T = TypeVar("T")


class ErrorInfo(BaseModel):
    code: str
    message: str
    details: str | None = None


class Meta(BaseModel):
    page: int | None = None
    per_page: int | None = None
    total: int | None = None
    total_pages: int | None = None


class ApiResponse(BaseModel, Generic[T]):
    success: bool = True
    data: T | None = None
    error: ErrorInfo | None = None
    meta: Meta | None = None


class RegisterRequest(BaseModel):
    email: str
    password: str = Field(min_length=8)
    name: str = Field(min_length=1)


class LoginRequest(BaseModel):
    email: str
    password: str


class UpdateUserRequest(BaseModel):
    name: str | None = None
    avatar: str | None = None
    preferences: dict[str, Any] | None = None


class ChangePasswordRequest(BaseModel):
    old_password: str
    new_password: str = Field(min_length=8)


class AnalyzeRequest(BaseModel):
    types: list[str] | None = None


class LearningGenerateRequest(BaseModel):
    category_id: int
    difficulty: str = "beginner"
    prerequisites: list[str] | None = None
    focus_areas: list[str] | None = None


class ReviewGenerateRequest(BaseModel):
    title: str
    paper_ids: list[str]
    category_id: int | None = None
    style: str | None = "academic"


class ReviewScoreRequest(BaseModel):
    content: str | None = None
