# PaperBeginner 技术栈学习指南

> 本文档面向初学者，帮助你了解继续完善本项目所需的技术栈知识。

---

## 🎯 技术栈总览与学习路线

### 一、后端开发（核心）

#### 1. Go 语言基础 ⭐⭐⭐⭐⭐ (最优先)

项目使用 Go 1.22+，你需要掌握：

| 知识点 | 重要程度 | 学习资源 |
|--------|---------|---------|
| 基本语法、数据类型 | 必须 | [Go Tour](https://go.dev/tour) |
| Goroutine 和 Channel | 必须 | [Go by Example](https://gobyexample.com/) |
| 接口和组合 | 必须 | [Effective Go](https://go.dev/doc/effective_go) |
| 错误处理 | 必须 | 项目中大量使用 |
| Context 包 | 必须 | 用于超时控制和取消 |
| 项目结构规范 | 重要 | [Go 标准项目布局](https://github.com/golang-standards/project-layout) |

#### 2. Gin Web 框架 ⭐⭐⭐⭐⭐

```go
// 项目使用的核心模式
router := gin.New()
router.Use(middleware.Logger())  // 中间件
router.GET("/api/v1/papers", paperHandler.List)  // 路由
```

需要掌握：
- 路由定义与分组
- 中间件机制（认证、日志、CORS、限流）
- 请求参数绑定与验证
- JSON 响应处理

**学习资源：** [Gin 官方文档](https://gin-gonic.com/docs/)

#### 3. GORM (ORM) ⭐⭐⭐⭐

```go
// 项目中的数据库操作模式
type Paper struct {
    ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
    Title     string    `gorm:"not null"`
    UserID    uuid.UUID `gorm:"type:uuid;index"`
}
```

需要掌握：
- 模型定义与关联关系（一对多、多对多）
- CRUD 操作
- 数据库迁移
- 查询构建器与预加载

**学习资源：** [GORM 官方文档](https://gorm.io/docs/)

#### 4. PostgreSQL + pgvector ⭐⭐⭐⭐

```sql
-- 项目使用向量存储支持语义搜索
CREATE EXTENSION IF NOT EXISTS vector;
```

需要掌握：
- SQL 基础（查询、索引、事务）
- pgvector 向量搜索（用于论文语义检索）
- 数据库设计范式

**学习资源：**
- [PostgreSQL 教程](https://www.postgresqltutorial.com/)
- [pgvector 文档](https://github.com/pgvector/pgvector)

#### 5. Redis ⭐⭐⭐

用于缓存和任务队列：
- 基本数据结构（String、Hash、List、Set）
- 缓存策略（TTL、缓存穿透/击穿/雪崩）
- 与 Asynq 任务队列配合

**学习资源：** [Redis 官方教程](https://redis.io/docs/getting-started/)

#### 6. Asynq 异步任务 ⭐⭐⭐

```go
// Worker 处理后台任务
type TaskPayload struct {
    PaperID string `json:"paper_id"`
}
```

需要掌握：
- 任务定义与调度
- Worker 模式
- 任务重试与错误处理
- 优先级队列

**学习资源：** [Asynq 官方文档](https://github.com/hibiken/asynq)

---

### 二、前端开发

#### 1. React 18 + TypeScript ⭐⭐⭐⭐⭐

```typescript
// 项目使用函数式组件 + Hooks
function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const { isAuthenticated } = useAuthStore()
  if (!isAuthenticated) {
    return <Navigate to="/login" replace />
  }
  return <>{children}</>
}
```

需要掌握：
- React Hooks（useState, useEffect, useContext, useMemo, useCallback）
- TypeScript 类型系统（接口、泛型、类型推断）
- 组件设计模式（组合 vs 继承）
- React Router v6 路由

**学习资源：**
- [React 官方文档](https://react.dev/)
- [TypeScript 官方文档](https://www.typescriptlang.org/docs/)

#### 2. 状态管理 ⭐⭐⭐⭐

| 库 | 用途 | 学习资源 |
|---|------|---------|
| Zustand | 全局状态（如用户认证状态） | [Zustand 文档](https://github.com/pmndrs/zustand) |
| TanStack Query | 服务端状态（API 数据缓存） | [TanStack Query 文档](https://tanstack.com/query/latest) |

```typescript
// Zustand 示例
const useAuthStore = create((set) => ({
  isAuthenticated: false,
  login: () => set({ isAuthenticated: true }),
  logout: () => set({ isAuthenticated: false }),
}))
```

#### 3. TailwindCSS ⭐⭐⭐

```html
<!-- 实用优先的 CSS 框架 -->
<button className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600">
  提交
</button>
```

需要掌握：
- 实用类命名规则
- 响应式设计（sm, md, lg, xl）
- 自定义主题配置
- 组件抽象

**学习资源：** [TailwindCSS 官方文档](https://tailwindcss.com/docs)

#### 4. 其他前端库 ⭐⭐

| 库 | 用途 |
|---|------|
| `react-dropzone` | 文件拖拽上传 |
| `react-markdown` | Markdown 渲染 |
| `recharts` | 数据可视化图表 |
| `framer-motion` | 动画效果 |
| `lucide-react` | 图标库 |
| `react-hot-toast` | 通知提示 |

---

### 三、AI/LLM 集成 ⭐⭐⭐⭐

这是项目的核心亮点，需要理解：

#### 1. LLM Provider 抽象

```go
// 项目定义的 Provider 接口（见 internal/llm/provider.go）
type Provider interface {
    // 文本补全
    Complete(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error)
    // 对话
    Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error)
    // 向量嵌入
    Embed(ctx context.Context, text string) ([]float32, error)
}
```

#### 2. 支持的 LLM 提供商

| 提供商 | 适用场景 | API 文档 |
|--------|---------|---------|
| OpenAI (GPT-4) | 通用任务 | [OpenAI API](https://platform.openai.com/docs) |
| Claude (Anthropic) | 长文本分析 | [Anthropic API](https://docs.anthropic.com/) |
| DeepSeek | 中文优化 | [DeepSeek API](https://platform.deepseek.com/) |
| Ollama | 本地部署 | [Ollama](https://ollama.ai/) |

#### 3. LangChainGo

项目使用 `github.com/tmc/langchaingo`，需要了解：
- Prompt 模板设计
- Chain 组合（顺序链、路由链）
- Agent 模式（ReAct、工具调用）
- 向量存储与检索

**学习资源：** [LangChainGo 文档](https://tmc.github.io/langchaingo/docs/)

#### 4. 项目中的 AI Agent

```
internal/agent/
├── learning_agent.go   # 学习路径生成
├── paper_agent.go      # 论文分析
├── review_agent.go     # 综述撰写
└── trend_agent.go      # 热点分析
```

---

### 四、DevOps 与基础设施

#### 1. Docker & Docker Compose ⭐⭐⭐⭐

```yaml
# 项目使用多容器编排（见 docker-compose.yml）
services:
  postgres:
    image: pgvector/pgvector:pg16
  redis:
    image: redis:7-alpine
  minio:
    image: minio/minio:latest
```

需要掌握：
- Dockerfile 编写（多阶段构建）
- 容器网络与数据卷
- 多服务编排
- 健康检查

**学习资源：** [Docker 官方文档](https://docs.docker.com/)

#### 2. MinIO ⭐⭐⭐

对象存储，用于存储上传的 PDF 论文文件：
- S3 兼容 API
- Bucket 管理
- 访问策略

**学习资源：** [MinIO 文档](https://min.io/docs/minio/linux/index.html)

#### 3. Traefik ⭐⭐

反向代理和负载均衡（生产环境）：
- 自动服务发现
- HTTPS 配置
- 路由规则

**学习资源：** [Traefik 文档](https://doc.traefik.io/traefik/)

---

### 五、API 设计与文档

#### 1. OpenAPI/Swagger ⭐⭐⭐

项目有 `api/openapi.yaml` 定义 API 规范：
- RESTful API 设计原则
- API 文档编写
- 请求/响应模式定义

**学习资源：** [OpenAPI 规范](https://swagger.io/specification/)

#### 2. JWT 认证 ⭐⭐⭐

```go
// 项目使用 JWT 进行用户认证
protected.Use(middleware.Auth(cfg.JWT.Secret))
```

需要了解：
- JWT 结构（Header、Payload、Signature）
- Access Token 与 Refresh Token
- Token 过期与刷新策略

---

## 📚 推荐学习路线（按时间顺序）

### 第一阶段：基础入门（2-4 周）

```
1. Go 语言基础 → Go Tour + Go by Example
2. SQL 基础 → PostgreSQL 教程
3. HTTP 和 REST API 概念
4. Git 版本控制
```

**目标：** 能够阅读和理解项目中的 Go 代码

### 第二阶段：后端开发（4-6 周）

```
1. Gin 框架 → 官方文档 + 实践
2. GORM 数据库操作
3. 中间件（认证、日志、错误处理）
4. Redis 缓存基础
```

**目标：** 能够独立添加新的 API 端点

### 第三阶段：前端开发（3-4 周）

```
1. TypeScript 基础
2. React 18 核心概念
3. TailwindCSS 样式
4. 状态管理（Zustand + TanStack Query）
```

**目标：** 能够独立开发前端页面

### 第四阶段：进阶技能（4-6 周）

```
1. Docker 容器化
2. LLM API 调用（OpenAI/Claude）
3. 向量数据库和语义搜索
4. 异步任务处理（Asynq）
```

**目标：** 能够完善 AI Agent 功能

---

## 🔧 适合初学者的待完成任务

根据当前项目状态，以下是一些适合初学者的任务：

### 简单难度 ⭐⭐

| 任务 | 描述 | 涉及技术 |
|------|------|---------|
| 完善 Settings 页面 | 添加用户设置表单 | React, TailwindCSS |
| 优化 Dashboard 布局 | 改进仪表板 UI | React, 组件设计 |
| 添加 Loading 状态 | 给 API 请求添加加载提示 | React, TanStack Query |

### 中等难度 ⭐⭐⭐

| 任务 | 描述 | 涉及技术 |
|------|------|---------|
| 实现 GitHub Trending 爬虫 | 爬取热门项目 | Go, HTTP Client |
| 添加分页功能 | 论文列表分页 | Go, React |
| 实现密码重置 | 邮箱验证流程 | Go, 邮件服务 |

### 较难难度 ⭐⭐⭐⭐

| 任务 | 描述 | 涉及技术 |
|------|------|---------|
| 完善 LLM Agent | 实现论文分析 Agent | Go, LangChainGo |
| 添加向量搜索 | 论文语义搜索 | pgvector, 嵌入模型 |
| 实现 SSE 流式输出 | AI 生成内容流式显示 | Go, React |

---

## 📁 项目结构详解

```
PaperBeginner/
├── backend/                      # Go 后端
│   ├── cmd/
│   │   ├── api/main.go          # API 服务入口
│   │   └── worker/main.go       # Worker 服务入口
│   ├── internal/
│   │   ├── agent/               # AI Agent 实现
│   │   ├── config/              # 配置管理（Viper）
│   │   ├── domain/              # 领域模型（实体定义）
│   │   ├── handler/             # HTTP 处理器（控制器）
│   │   │   └── middleware/      # 中间件
│   │   ├── llm/                 # LLM 适配层
│   │   ├── repository/          # 数据访问层（DAO）
│   │   ├── service/             # 业务逻辑层
│   │   └── worker/              # 异步任务处理
│   ├── pkg/                     # 公共工具包
│   │   ├── logger/              # 日志（Zap）
│   │   ├── pdf/                 # PDF 解析
│   │   ├── response/            # 统一响应格式
│   │   └── validator/           # 请求验证
│   ├── data/                    # 数据库迁移文件
│   └── api/openapi.yaml         # API 规范
├── frontend/                     # React 前端
│   └── src/
│       ├── components/          # 可复用组件
│       ├── pages/               # 页面组件
│       ├── services/            # API 调用封装
│       ├── stores/              # Zustand 状态
│       └── types/               # TypeScript 类型
├── deployments/                  # 部署配置
│   ├── config/                  # 配置文件模板
│   └── docker/                  # Dockerfile
├── docs/                         # 文档
├── scripts/                      # 脚本工具
├── docker-compose.yml           # 容器编排
└── Makefile                     # 常用命令
```

---

## 💡 学习建议

1. **从阅读代码开始**：先理解现有代码的结构和风格
2. **边做边学**：选择一个小任务，遇到不懂的概念就去查文档
3. **使用 AI 辅助**：善用 Claude/GPT 解释代码和概念
4. **写测试验证**：通过写测试来验证你对代码的理解
5. **提交小 PR**：保持提交粒度小，便于 Review

---

## 🔗 推荐资源汇总

### Go 语言
- [Go Tour](https://go.dev/tour) - 官方入门教程
- [Go by Example](https://gobyexample.com/) - 示例驱动学习
- [Effective Go](https://go.dev/doc/effective_go) - 最佳实践

### 前端
- [React 官方文档](https://react.dev/)
- [TypeScript 官方文档](https://www.typescriptlang.org/docs/)
- [TailwindCSS 文档](https://tailwindcss.com/docs)

### 数据库
- [PostgreSQL 教程](https://www.postgresqltutorial.com/)
- [GORM 文档](https://gorm.io/docs/)
- [Redis 文档](https://redis.io/docs/)

### AI/LLM
- [OpenAI API 文档](https://platform.openai.com/docs)
- [LangChain 概念](https://python.langchain.com/docs/concepts/)

### DevOps
- [Docker 入门](https://docs.docker.com/get-started/)
- [Docker Compose 教程](https://docs.docker.com/compose/gettingstarted/)

---

*最后更新：2024年12月*

