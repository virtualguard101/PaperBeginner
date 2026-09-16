from pathlib import Path

from pydantic_settings import BaseSettings, SettingsConfigDict


ROOT_DIR = Path(__file__).resolve().parent.parent
PROJECT_ROOT = ROOT_DIR.parent


class Settings(BaseSettings):
    model_config = SettingsConfigDict(env_file=str(PROJECT_ROOT / ".env"), extra="ignore")

    app_env: str = "development"
    app_port: int = 8000
    jwt_secret: str = "paperbeginner-demo-secret-change-me"
    jwt_expire_hours: int = 72

    database_url: str = f"sqlite:///{(ROOT_DIR / 'data' / 'paperbeginner.db').as_posix()}"
    upload_dir: str = str(ROOT_DIR / "data" / "uploads")

    llm_api_key: str = ""
    llm_base_url: str = "https://api.deepseek.com/v1"
    llm_model: str = "deepseek-chat"
    llm_timeout_seconds: float = 90.0

    github_token: str = ""


settings = Settings()
