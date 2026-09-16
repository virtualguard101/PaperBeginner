# PaperBeginner

AI 驱动的学术研究引导平台（课设演示版）。

当前实现：**React 前端 + FastAPI 后端 + SQLite**。原 Go 后端已归档到 `backend-go/`，仅作历史参考。

## 演示功能

- 注册 / 登录（JWT）
- 热点：GitHub 搜索热门仓库 + CCF 种子条目，可生成周期短报
- 论文：上传 PDF，一次分析得到摘要 / 方法 / 贡献
- 学习路径：按 CCF 领域与难度生成分阶段路线
- 综述：基于已选论文生成 Markdown 并打分

未配置 `LLM_API_KEY` 时，LLM 相关接口返回离线模板，流程仍可走完。

## 本地启动（推荐演示）

需要 Python 3.11+ 与 Node.js 20+。

```bash
copy .env.example .env
# 可选：在 .env 中填写 LLM_API_KEY（DeepSeek 等 OpenAI 兼容接口）

cd backend
python -m pip install -r requirements.txt
python -m uvicorn app.main:app --reload --host 0.0.0.0 --port 8000
```

另开终端：

```bash
cd frontend
npm install
npm run dev
```

打开 http://localhost:5173

1. 注册账号（密码至少 8 位）
2. 热点页查看列表与报告
3. 上传 `frontend/public/sample-attention.pdf`（或任意短 PDF）并点击分析
4. 生成学习路线
5. 用该论文生成综述并评分

API：http://localhost:8000/health  
文档：http://localhost:8000/docs

## Docker 一键

```bash
docker compose up --build
```

浏览器打开 http://localhost:5173 （前端容器映射 80→5173，并反代 `/api`）。

## 环境变量

见 `.env.example`：

- `LLM_API_KEY` / `LLM_BASE_URL` / `LLM_MODEL`：默认 DeepSeek Chat
- `JWT_SECRET`
- `GITHUB_TOKEN`：可选，提高 GitHub API 限额；失败则用种子数据

## 架构（演示）

```
React (Vite :5173)
  -> /api/v1 代理到 FastAPI :8000
       -> SQLite + 本地 uploads/
       -> OpenAI 兼容 LLM
       -> GitHub API（失败回退种子）
```

两周内刻意不做：Postgres/pgvector、MinIO、Redis/Asynq、Traefik、多 LLM 路由、完整 CCF 爬虫。

## 许可证

AGPL-3.0
