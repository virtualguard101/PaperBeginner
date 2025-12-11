# PaperBeginner 架构设计文档

---

## 项目概述

### 项目定位

PaperBeginner 是一个 **AI 驱动的学术研究引导平台**，专为计算机科学领域的学术新人设计。平台通过多个 AI Agent 协同工作，帮助用户追踪研究热点、规划学习路线、分析论文以及撰写论文综述。

### 核心设计理念

- **Agent 驱动**: 多个专业化 AI Agent 分工协作，各司其职

- **多 LLM 支持**: 支持 OpenAI、Claude、DeepSeek、Ollama 等多种 LLM 提供商，可根据任务类型智能路由

- **异步处理**: 采用任务队列处理耗时操作，提升用户体验

- **向量检索**: 利用 pgvector 实现论文语义搜索

- **前后端分离**: Go 后端 + React 前端，清晰的职责划分

---

## 技术栈

### 后端技术

| 技术 | 用途 | 版本 |
|------|------|------|
| Go | 主要开发语言 | 1.22+ |
| Gin | Web 框架 | - |
| GORM | ORM 框架 | - |
| PostgreSQL | 主数据库 | 16 |
| pgvector | 向量检索扩展 | - |
| Redis | 缓存 & 任务队列 | 7 |
| Asynq | 异步任务处理 | - |
| MinIO | 对象存储 | - |

### 前端技术

| 技术 | 用途 | 版本 |
|------|------|------|
| React | UI 框架 | 18.3+ |
| TypeScript | 类型安全 | - |
| Vite | 构建工具 | - |
| TailwindCSS | 样式框架 | - |
| Zustand | 状态管理 | - |
| TanStack Query | 数据请求 | - |

### 基础设施

| 技术 | 用途 |
|------|------|
| Docker | 容器化 |
| Docker Compose | 容器编排 |
| Traefik | 反向代理 / API 网关 |

---

## 系统架构图

### 整体架构

```
┌──────────────────────────────────────────────────────────────────────────┐
│                            客户端 (React SPA)                              │
├──────────────────────────────────────────────────────────────────────────┤
│                          API 网关 (Traefik)                               │
├────────────────┬─────────────────────────────────────────────────────────┤
│                │                     后端服务                              │
│                ├───────────────────┬───────────────────┬─────────────────┤
│   API Server   │   Worker Service  │    AI Agents      │   LLM Router    │
│    (Gin)       │     (Asynq)       │                   │                 │
├────────────────┴───────────────────┴───────────────────┴─────────────────┤
│                            消息队列 (Redis Stream)                         │
├────────────────┬───────────────────┬───────────────────┬─────────────────┤
│   PostgreSQL   │      Redis        │      MinIO        │   LLM APIs      │
│   + pgvector   │     (缓存)         │    (文件存储)      │  (外部服务)      │
└────────────────┴───────────────────┴───────────────────┴─────────────────┘
```

### 模块依赖图

```mermaid
graph TB
    subgraph "入口层"
        API[cmd/api/main.go<br/>API 服务入口]
        Worker[cmd/worker/main.go<br/>Worker 服务入口]
    end
    
    subgraph "处理层"
        Handler[handler/<br/>HTTP 处理器]
        Middleware[middleware/<br/>中间件]
        WorkerHandler[worker/<br/>任务处理器]
    end
    
    subgraph "业务层"
        Service[service/<br/>业务逻辑]
        Agent[agent/<br/>AI Agents]
    end
    
    subgraph "基础设施层"
        LLM[llm/<br/>LLM 提供商]
        Repository[repository/<br/>数据访问层]
        Config[config/<br/>配置管理]
    end
    
    subgraph "公共包"
        Logger[pkg/logger]
        PDF[pkg/pdf]
        Response[pkg/response]
        Validator[pkg/validator]
    end
    
    API --> Handler
    API --> Middleware
    API --> Config
    
    Worker --> WorkerHandler
    Worker --> Config
    
    Handler --> Service
    WorkerHandler --> Service
    WorkerHandler --> Agent
    
    Service --> Repository
    Service --> Agent
    
    Agent --> LLM
    
    LLM --> Config
    Repository --> Config
```

---

## 项目结构

```
PaperBeginner/
├── backend/                     # Go 后端
│   ├── cmd/                     # 应用入口
│   │   ├── api/                 # API 服务
│   │   │   └── main.go
│   │   └── worker/              # Worker 服务
│   │       └── main.go
│   │
│   ├── internal/                # 内部包（不对外暴露）
│   │   ├── agent/               # AI Agent 实现
│   │   │   ├── learning_agent.go    # 学习路径 Agent
│   │   │   ├── paper_agent.go       # 论文分析 Agent
│   │   │   ├── review_agent.go      # 综述撰写 Agent
│   │   │   └── trend_agent.go       # 热点分析 Agent
│   │   │
│   │   ├── config/              # 配置管理
│   │   │   └── config.go
│   │   │
│   │   ├── domain/              # 领域模型
│   │   │   ├── analysis.go
│   │   │   ├── paper.go
│   │   │   └── user.go
│   │   │
│   │   ├── handler/             # HTTP 处理器
│   │   │   ├── middleware/      # 中间件
│   │   │   ├── learning_handler.go
│   │   │   ├── paper_handler.go
│   │   │   ├── review_handler.go
│   │   │   ├── trending_handler.go
│   │   │   └── user_handler.go
│   │   │
│   │   ├── llm/                 # LLM 适配层
│   │   │   ├── claude.go            # Claude 提供商
│   │   │   ├── deepseek.go          # DeepSeek 提供商
│   │   │   ├── ollama.go            # Ollama 提供商
│   │   │   ├── openai.go            # OpenAI 提供商
│   │   │   ├── provider.go          # 提供商接口定义
│   │   │   └── router.go            # LLM 路由器
│   │   │
│   │   ├── repository/          # 数据访问层
│   │   │   └── postgres/
│   │   │       ├── connection.go
│   │   │       ├── learning_path_repository.go
│   │   │       ├── paper_repository.go
│   │   │       ├── review_repository.go
│   │   │       ├── trending_repository.go
│   │   │       └── user_repository.go
│   │   │
│   │   ├── service/             # 业务逻辑层
│   │   │   ├── crawler_service.go
│   │   │   ├── paper_service.go
│   │   │   └── user_service.go
│   │   │
│   │   └── worker/              # 异步任务
│   │       ├── handler.go
│   │       └── tasks.go
│   │
│   ├── pkg/                     # 公共包（可对外）
│   │   ├── logger/              # 日志
│   │   ├── pdf/                 # PDF 解析
│   │   ├── response/            # 响应封装
│   │   └── validator/           # 参数校验
│   │
│   ├── data/              # 数据库迁移
│   ├── api/                     # OpenAPI 规范
│   ├── config.yaml              # 配置文件
│   ├── go.mod
│   └── go.sum
│
├── frontend/                    # React 前端
│   ├── src/
│   │   ├── components/          # UI 组件
│   │   │   └── Layout.tsx
│   │   ├── pages/               # 页面组件
│   │   │   ├── DashboardPage.tsx
│   │   │   ├── HomePage.tsx
│   │   │   ├── LearningPage.tsx
│   │   │   ├── LoginPage.tsx
│   │   │   ├── PapersPage.tsx
│   │   │   ├── RegisterPage.tsx
│   │   │   ├── ReviewsPage.tsx
│   │   │   ├── SettingsPage.tsx
│   │   │   └── TrendingPage.tsx
│   │   ├── services/            # API 调用
│   │   │   └── api.ts
│   │   ├── stores/              # 状态管理
│   │   │   └── authStore.ts
│   │   ├── types/               # 类型定义
│   │   ├── hooks/               # 自定义 Hooks
│   │   ├── App.tsx
│   │   ├── main.tsx
│   │   └── index.css
│   │
│   ├── package.json
│   ├── vite.config.ts
│   ├── tailwind.config.js
│   └── tsconfig.json
│
├── deployments/                 # 部署配置
│   ├── config/
│   │   ├── config.example.yaml
│   │   └── nginx.conf
│   └── docker/
│       ├── backend.Dockerfile
│       ├── backend.dev.Dockerfile
│       ├── frontend.Dockerfile
│       ├── frontend.dev.Dockerfile
│       ├── worker.Dockerfile
│       ├── docker-compose.monitoring.yml
│       ├── prometheus.yml
│       └── promtail-config.yml
│
├── scripts/                     # 脚本
│   ├── init-db.sh
│   └── seed-data.sh
│
├── docs/                        # 文档
├── docker-compose.yml           # 生产环境编排
├── docker-compose.dev.yml       # 开发环境编排
├── Makefile                     # 构建命令
└── README.md
```

---

## 核心模块详解

### 1. AI Agent 模块 (`internal/agent/`)

**职责**: 封装各类 AI 任务的业务逻辑，通过 LLM 提供智能分析能力

#### Agent 类型

| Agent | 文件 | 职责 |
|-------|------|------|
| PaperAgent | `paper_agent.go` | 论文分析（摘要、方法论、贡献点、优缺点） |
| TrendAgent | `trend_agent.go` | 热点趋势分析与分类 |
| LearningAgent | `learning_agent.go` | 学习路径生成 |
| ReviewAgent | `review_agent.go` | 综述撰写与评分 |

#### PaperAgent 设计

```go
type PaperAgent struct {
    router *llm.Router  // LLM 路由器
}

// 分析类型
const (
    AnalysisSummary       // 摘要
    AnalysisMethodology   // 方法论
    AnalysisContributions // 贡献点
    AnalysisStrengths     // 优点
    AnalysisWeaknesses    // 缺点
)

// 核心方法
func (a *PaperAgent) AnalyzePaper(ctx, paper, analysisType) (*PaperAnalysisResult, error)
func (a *PaperAgent) ExtractKeywords(ctx, paper) ([]string, error)
```

### 2. LLM 适配层 (`internal/llm/`)

**职责**: 统一管理多个 LLM 提供商，实现智能路由与降级

#### 提供商接口

```go
type Provider interface {
    Name() string
    IsAvailable() bool
    Complete(ctx, req) (*CompletionResponse, error)
    Chat(ctx, req) (*ChatResponse, error)
    Embed(ctx, text) ([]float32, error)
    GetMaxTokens() int
    GetEmbeddingDimension() int
}
```

#### LLM Router

```go
type Router struct {
    providers   map[string]Provider      // 已注册的提供商
    taskMapping map[TaskType]string      // 任务类型 -> 默认提供商
}

// 任务类型与默认提供商映射
TaskTrendAnalysis   -> "anthropic"    // Claude 擅长分析
TaskPaperParsing    -> "openai"       // GPT-4o 支持多模态
TaskPaperSummary    -> "deepseek"     // 成本效益高
TaskLearningPath    -> "deepseek"     // 成本效益高
TaskReviewWriting   -> "anthropic"    // Claude 擅长长文写作
TaskReviewScoring   -> "openai"       // GPT-4 评估准确
TaskEmbedding       -> "openai"       // OpenAI Embedding
```

#### 路由策略

```mermaid
graph TD
    Request[请求] --> CheckPref{用户偏好?}
    CheckPref -->|是| UsePreferred[使用用户指定提供商]
    CheckPref -->|否| UseDefault[使用任务默认提供商]
    
    UsePreferred --> CheckAvail{可用?}
    UseDefault --> CheckAvail
    
    CheckAvail -->|是| Execute[执行请求]
    CheckAvail -->|否| Fallback[降级到其他可用提供商]
    
    Execute --> CheckError{成功?}
    CheckError -->|是| Return[返回结果]
    CheckError -->|否| TryFallback{启用降级?}
    TryFallback -->|是| Fallback
    TryFallback -->|否| Error[返回错误]
    
    Fallback --> Execute
```

### 3. 数据访问层 (`internal/repository/`)

**职责**: 封装数据库操作，提供统一的数据访问接口

#### Repository 接口设计

```go
// PaperRepository 论文数据访问接口
type PaperRepository interface {
    Create(paper *Paper) error
    GetByID(id uuid.UUID) (*Paper, error)
    GetByUserID(userID uuid.UUID, offset, limit int) ([]*Paper, int64, error)
    Update(paper *Paper) error
    Delete(id uuid.UUID) error
    Search(query string, categoryID *int, limit int) ([]*Paper, error)
    FindSimilar(embedding pgvector.Vector, limit int) ([]*Paper, error)
    CreateAnalysis(analysis *PaperAnalysis) error
    GetAnalysesByPaperID(paperID uuid.UUID) ([]PaperAnalysis, error)
    GetCCFCategories() ([]CCFCategory, error)
}
```

### 4. Worker 模块 (`internal/worker/`)

**职责**: 处理异步任务，如论文分析、热点爬取等耗时操作

#### 任务类型

```go
const (
    TypePaperAnalysis    = "paper:analysis"      // 论文分析
    TypeGitHubCrawl      = "crawler:github"      // GitHub 爬取
    TypeCCFCrawl         = "crawler:ccf"         // CCF 论文爬取
    TypeTrendReport      = "trend:report"        // 趋势报告生成
    TypeLearningPath     = "learning:path"       // 学习路径生成
    TypeReviewGeneration = "review:generate"     // 综述生成
    TypeReviewScoring    = "review:score"        // 综述评分
)
```

---

## 数据模型

### ER 图

```mermaid
erDiagram
    users ||--o{ papers : owns
    users ||--o{ learning_paths : creates
    users ||--o{ reviews : writes
    
    papers ||--o{ paper_analyses : has
    papers }o--|| ccf_categories : belongs_to
    
    trending_items }o--|| ccf_categories : belongs_to
    trend_reports }o--|| ccf_categories : belongs_to
    learning_paths }o--|| ccf_categories : belongs_to
    reviews }o--|| ccf_categories : belongs_to
    
    reviews ||--o{ review_papers : contains
    papers ||--o{ review_papers : included_in
    
    users {
        uuid id PK
        string email UK
        string password_hash
        string name
        string avatar
        string role
        boolean is_active
        jsonb preferences
        timestamp created_at
        timestamp updated_at
        timestamp last_login_at
    }
    
    papers {
        uuid id PK
        uuid user_id FK
        string title
        text[] authors
        text abstract
        text[] keywords
        timestamp published_at
        string venue
        string doi UK
        string file_path
        int ccf_category_id FK
        string ccf_rank
        string status
        vector embedding
        timestamp created_at
        timestamp updated_at
    }
    
    paper_analyses {
        uuid id PK
        uuid paper_id FK
        string analysis_type
        text content
        jsonb structured_data
        string llm_provider
        string llm_model
        int tokens_used
        timestamp created_at
    }
    
    ccf_categories {
        int id PK
        string name
        string name_en
    }
    
    trending_items {
        uuid id PK
        string source
        string source_id
        string title
        text description
        string url
        int stars
        int forks
        string language
        text[] topics
        int ccf_category_id FK
        float trend_score
        timestamp crawled_at
    }
    
    trend_reports {
        uuid id PK
        string period
        int category_id FK
        text summary
        jsonb highlights
        string generated_by
        timestamp created_at
    }
    
    learning_paths {
        uuid id PK
        uuid user_id FK
        int category_id FK
        string title
        text description
        jsonb stages
        text[] prerequisites
        string estimated_time
        string difficulty
        timestamp created_at
    }
    
    reviews {
        uuid id PK
        uuid user_id FK
        string title
        text abstract
        text content
        uuid[] paper_ids
        int category_id FK
        string status
        jsonb score
        timestamp created_at
    }
```

### CCF 分类

| ID | 中文名称 | 英文名称 |
|----|---------|---------|
| 1 | 计算机体系结构/并行与分布计算/存储系统 | Computer Architecture/Parallel and Distributed Computing/Storage Systems |
| 2 | 计算机网络 | Computer Networks |
| 3 | 网络与信息安全 | Network and Information Security |
| 4 | 软件工程/系统软件/程序设计语言 | Software Engineering/System Software/Programming Languages |
| 5 | 数据库/数据挖掘/内容检索 | Database/Data Mining/Content Retrieval |
| 6 | 计算机科学理论 | Computer Science Theory |
| 7 | 计算机图形学与多媒体 | Computer Graphics and Multimedia |
| 8 | 人工智能 | Artificial Intelligence |
| 9 | 人机交互与普适计算 | Human-Computer Interaction and Ubiquitous Computing |
| 10 | 交叉/综合/新兴 | Interdisciplinary/Comprehensive/Emerging |

---

## 核心工作流

### 论文分析流程

```mermaid
sequenceDiagram
    participant User as 用户
    participant API as API Server
    participant Queue as Redis Queue
    participant Worker as Worker
    participant Agent as PaperAgent
    participant LLM as LLM Provider
    participant DB as PostgreSQL
    
    User->>API: POST /papers/upload
    API->>DB: 保存论文元数据
    API-->>User: 返回 paper_id
    
    User->>API: POST /papers/:id/analyze
    API->>Queue: 入队分析任务
    API-->>User: 返回 "处理中"
    
    Queue->>Worker: 分发任务
    activate Worker
    
    Worker->>DB: 获取论文信息
    Worker->>Agent: 调用 AnalyzePaper()
    
    Agent->>Agent: 构建 Prompt
    Agent->>LLM: Chat 请求
    LLM-->>Agent: 分析结果
    
    Agent-->>Worker: PaperAnalysisResult
    Worker->>DB: 保存分析结果
    deactivate Worker
    
    User->>API: GET /papers/:id/analyses
    API->>DB: 查询分析结果
    API-->>User: 返回分析列表
```

### 学习路径生成流程

```mermaid
sequenceDiagram
    participant User as 用户
    participant API as API Server
    participant Queue as Redis Queue
    participant Worker as Worker
    participant Agent as LearningAgent
    participant LLM as LLM Provider
    participant DB as PostgreSQL
    
    User->>API: POST /learning/generate
    Note right of User: 包含: 研究方向、难度、前置知识
    
    API->>Queue: 入队学习路径生成任务
    API-->>User: 返回 path_id
    
    Queue->>Worker: 分发任务
    activate Worker
    
    Worker->>Agent: GenerateLearningPath()
    Agent->>Agent: 整合资源(MIT/Stanford/csdiy.wiki)
    Agent->>LLM: 生成个性化路径
    LLM-->>Agent: 学习路径 JSON
    
    Agent-->>Worker: LearningPathResult
    Worker->>DB: 保存学习路径
    deactivate Worker
    
    User->>API: GET /learning/paths/:id
    API->>DB: 查询学习路径
    API-->>User: 返回完整路径
```

### 综述生成与评分流程

```mermaid
sequenceDiagram
    participant User as 用户
    participant API as API Server
    participant Queue as Redis Queue
    participant Worker as Worker
    participant Agent as ReviewAgent
    participant LLM as LLM Provider
    participant DB as PostgreSQL
    
    User->>API: POST /reviews/generate
    Note right of User: 包含: 论文 ID 列表、主题
    
    API->>DB: 创建综述记录 (status=draft)
    API->>Queue: 入队综述生成任务
    API-->>User: 返回 review_id
    
    Queue->>Worker: 分发任务
    activate Worker
    
    Worker->>DB: 获取相关论文
    Worker->>Agent: GenerateReview()
    Agent->>LLM: 生成综述内容
    LLM-->>Agent: 综述文本
    
    Agent-->>Worker: ReviewResult
    Worker->>DB: 更新综述 (status=completed)
    deactivate Worker
    
    User->>API: POST /reviews/:id/score
    API->>Queue: 入队评分任务
    
    Queue->>Worker: 分发任务
    activate Worker
    Worker->>Agent: ScoreReview()
    Agent->>LLM: 按学术标准评分
    LLM-->>Agent: 评分结果
    
    Agent-->>Worker: ScoreResult
    Worker->>DB: 保存评分
    deactivate Worker
    
    User->>API: GET /reviews/:id
    API->>DB: 查询综述及评分
    API-->>User: 返回完整综述
```

---

## API 设计

### 路由结构

```
/api/v1
├── /auth                    # 认证相关（公开）
│   ├── POST /register       # 用户注册
│   ├── POST /login          # 用户登录
│   └── POST /refresh        # 刷新 Token
│
├── /users                   # 用户相关（需认证）
│   ├── GET /me              # 获取当前用户
│   ├── PUT /me              # 更新当前用户
│   └── PUT /me/password     # 修改密码
│
├── /papers                  # 论文管理（需认证）
│   ├── POST /upload         # 上传论文
│   ├── GET /                # 论文列表
│   ├── GET /:id             # 论文详情
│   ├── DELETE /:id          # 删除论文
│   ├── POST /:id/analyze    # 分析论文
│   └── GET /:id/analyses    # 获取分析结果
│
├── /trending                # 热点趋势（需认证）
│   ├── GET /                # 热点列表
│   └── GET /reports         # 趋势报告
│
├── /learning                # 学习路径（需认证）
│   ├── POST /generate       # 生成学习路径
│   ├── GET /paths           # 路径列表
│   └── GET /paths/:id       # 路径详情
│
├── /reviews                 # 综述管理（需认证）
│   ├── POST /generate       # 生成综述
│   ├── GET /                # 综述列表
│   ├── GET /:id             # 综述详情
│   ├── POST /:id/score      # 综述评分
│   └── DELETE /:id          # 删除综述
│
└── /ccf
    └── GET /categories      # CCF 分类列表
```

### 中间件链

```
Request → RequestID → Logger → CORS → RateLimit → Auth → Handler → Response
```

---

## 配置管理

### 配置结构

```go
type Config struct {
    App      AppConfig      // 应用配置
    Database DatabaseConfig // 数据库配置
    Redis    RedisConfig    // Redis 配置
    MinIO    MinIOConfig    // 对象存储配置
    JWT      JWTConfig      // JWT 配置
    LLM      LLMConfig      // LLM 配置
    GitHub   GitHubConfig   // GitHub API 配置
}
```

### 配置加载优先级

1. **环境变量** (最高优先级)
2. **配置文件** (`config.yaml`)
3. **默认值** (代码内置)

### 示例配置

```yaml
app:
  env: development
  port: 8080
  rate_limit_rps: 10
  rate_limit_burst: 20

database:
  url: postgres://user:pass@localhost:5432/paperbeginner?sslmode=disable
  max_open_conns: 25
  max_idle_conns: 5
  conn_max_lifetime: 5m

redis:
  url: redis://localhost:6379
  db: 0

minio:
  endpoint: localhost:9000
  access_key: minioadmin
  secret_key: minioadmin
  bucket: paperbeginner

jwt:
  secret: your-secret-key
  expiry: 24h

llm:
  default: openai
  openai:
    api_key: sk-xxx
    model: gpt-4o
  anthropic:
    api_key: sk-xxx
    model: claude-3-5-sonnet-20241022
  deepseek:
    api_key: sk-xxx
    model: deepseek-chat
  ollama:
    host: http://localhost:11434
    model: llama3.1

github:
  token: ghp_xxx
```

---

## 部署架构

### Docker Compose 服务

```yaml
services:
  postgres:    # PostgreSQL + pgvector
  redis:       # Redis 缓存 & 任务队列
  minio:       # MinIO 对象存储
  api:         # 后端 API 服务
  worker:      # 异步任务处理服务
  frontend:    # 前端静态资源
  traefik:     # 反向代理 (生产环境)
```

### 服务依赖关系

```mermaid
graph LR
    Frontend[Frontend] --> Traefik[Traefik]
    Traefik --> API[API Server]
    
    API --> PostgreSQL[(PostgreSQL)]
    API --> Redis[(Redis)]
    API --> MinIO[(MinIO)]
    
    Worker[Worker] --> PostgreSQL
    Worker --> Redis
    Worker --> MinIO
    Worker --> LLM[LLM APIs]
```

### 健康检查

| 服务 | 检查方式 | 间隔 |
|------|---------|------|
| PostgreSQL | `pg_isready` | 10s |
| Redis | `redis-cli ping` | 10s |
| MinIO | `mc ready local` | 30s |
| API | `GET /health` | - |

---

## 设计原则

### 1. 分层架构

**原则**: 严格遵循分层架构，各层职责清晰

```
Handler → Service → Repository → Database
    ↓         ↓
  Agent → LLM Provider
```

- **Handler**: 处理 HTTP 请求/响应，参数校验
- **Service**: 业务逻辑编排
- **Agent**: AI 任务封装
- **Repository**: 数据访问抽象
- **LLM Provider**: LLM 调用封装

### 2. 依赖注入

**原则**: 通过接口注入依赖，便于测试和替换

```go
// 接口定义
type PaperRepository interface {
    Create(paper *Paper) error
    GetByID(id uuid.UUID) (*Paper, error)
    // ...
}

// 依赖注入
func NewPaperService(repo PaperRepository, cfg *config.Config) *PaperService {
    return &PaperService{repo: repo, cfg: cfg}
}
```

### 3. 异步优先

**原则**: 耗时操作通过任务队列异步处理

- ✅ 论文分析、综述生成等 AI 任务
- ✅ GitHub/CCF 数据爬取
- ✅ 报告生成

### 4. 优雅降级

**原则**: LLM 调用支持多提供商降级

```go
// 主提供商失败时自动降级
if pref.FallbackEnabled {
    for _, p := range r.providers {
        if p.IsAvailable() && p.Name() != provider.Name() {
            return p.Chat(ctx, req)
        }
    }
}
```

### 5. 安全性

**原则**: 多层安全防护

- JWT Token 认证
- 请求频率限制
- CORS 配置
- 密码哈希存储
- SQL 注入防护 (GORM)

---

## 扩展指南

### 添加新的 LLM 提供商

1. 实现 `Provider` 接口:

```go
// internal/llm/new_provider.go
type NewProvider struct {
    config ProviderConfig
}

func (p *NewProvider) Name() string { return "new_provider" }
func (p *NewProvider) IsAvailable() bool { /* ... */ }
func (p *NewProvider) Chat(ctx, req) (*ChatResponse, error) { /* ... */ }
func (p *NewProvider) Embed(ctx, text) ([]float32, error) { /* ... */ }
// ...
```

2. 在 Router 中注册:

```go
// internal/llm/router.go
func NewRouter(cfg config.LLMConfig) (*Router, error) {
    // ...
    if cfg.NewProvider.APIKey != "" {
        provider, _ := NewNewProvider(...)
        r.providers["new_provider"] = provider
    }
}
```

### 添加新的 Agent

1. 创建 Agent 文件:

```go
// internal/agent/new_agent.go
type NewAgent struct {
    router *llm.Router
}

func NewNewAgent(router *llm.Router) *NewAgent {
    return &NewAgent{router: router}
}

func (a *NewAgent) Process(ctx, input) (*Result, error) {
    // 构建 Prompt
    // 调用 LLM
    // 解析结果
}
```

2. 在 Worker 中集成:

```go
// internal/worker/handler.go
func (h *Handler) HandleNewTask(ctx, task) error {
    agent := agent.NewNewAgent(h.router)
    return agent.Process(ctx, payload)
}
```

### 添加新的 API 端点

1. 在 `handler/` 创建 Handler
2. 在 `main.go` 注册路由
3. 更新 OpenAPI 规范

---

## 参考资料

### 项目相关

- [README](../README.md) - 项目介绍与快速开始
- [OpenAPI 规范](../backend/api/openapi.yaml) - API 文档
- [数据库迁移](../backend/migrations/) - 数据库 Schema

### 技术文档

- [Gin Web Framework](https://gin-gonic.com/docs/)
- [GORM](https://gorm.io/docs/)
- [Asynq](https://github.com/hibiken/asynq) - 异步任务队列
- [pgvector](https://github.com/pgvector/pgvector) - 向量相似度搜索

### 设计参考

- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [12-Factor App](https://12factor.net/)

---

*文档版本: 1.0.0*  
*最后更新: 2025-12-10*

