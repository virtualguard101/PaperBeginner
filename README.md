<div align="center">

# PaperBeginner

![PaperBeginner Logo](https://img.shields.io/badge/PaperBeginner-AI%20Academic%20Guide-blue?style=for-the-badge)

**AI 驱动的学术研究引导平台，助力学术新人快速入门计算机学术研究**

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go)](https://go.dev/) [![React](https://img.shields.io/badge/React-18.3+-61DAFB?style=flat-square&logo=react)](https://react.dev/) [![License](https://img.shields.io/badge/License-AGPL%203.0-green?style=flat-square)](LICENSE)

[English](README.en.md) | [中文](#中文)

</div>

---

## 🎯 项目简介

PaperBeginner 是一个基于 AI Agent 的学术研究辅助平台，专为计算机科学领域的学术新人设计。通过智能化工具帮助用户：

- 📊 **追踪研究热点** - 实时监控 GitHub 热门项目和 CCF 顶会论文
- 📚 **规划学习路线** - AI 生成个性化学习路径，整合顶尖院校资源
- 📝 **智能论文分析** - 上传论文获取智能摘要、方法论解析
- 📖 **撰写论文综述** - AI 辅助生成高质量综述并提供评分

## ✨ 核心功能

### 1. 行业前沿热点嗅探

- 自动爬取 GitHub Trending 项目

- 整合 CCF 推荐国际学术刊物/会议论文

- AI 智能分类到 CCF 十大领域

- 生成周/月度热点报告

### 2. 个性化学习路线

- 根据研究方向生成学习路径

- 整合官方文档、MIT/Stanford 公开课、csdiy.wiki 资源

- 分阶段学习目标与进度追踪

### 3. 论文智能分析

- PDF 上传与文本提取

- 论文摘要、方法论、贡献点自动提取

- 向量化存储支持语义搜索

- 多维度深度分析

### 4. 综述撰写与评分

- 基于多篇论文生成论文综述

- 按学术标准进行评分

- 提供改进建议

## 🏗️ 技术架构

```
┌─────────────────────────────────────────────────────────────┐
│                      Frontend (React)                        
├─────────────────────────────────────────────────────────────┤
│                    API Gateway (Traefik)                     
├──────────────┬──────────────┬──────────────┬────────────────┤
│  User Service   Paper Service    Agent Service   Crawler Service  
├──────────────┴──────────────┴──────────────┴────────────────┤
│                     Message Queue (Redis Stream)             
├──────────────┬──────────────┬──────────────────────────────┤
│  PostgreSQL         Redis            MinIO (文件存储)            
└──────────────┴──────────────┴──────────────────────────────┘
```

## 🛠️ 技术栈

**后端:**

- Go 1.22+ / Gin / GORM

- PostgreSQL + pgvector

- Redis / Asynq

- LangChainGo

**前端:**

- React 18 / TypeScript

- Vite / TailwindCSS

- TanStack Query / Zustand

**基础设施:**

- Docker / Docker Compose

- MinIO (对象存储)

- Traefik (反向代理)

## 🚀 快速开始

### 前置要求

- Go 1.22+

- Node.js 20+

- Docker & Docker Compose

### 安装步骤

1. **克隆仓库**

```bash
git clone https://github.com/virtualguard/PaperBeginner.git
cd PaperBeginner
```

2. **启动基础设施**

```bash
make dev-infra
```

3. **配置环境变量**

```bash
cp deployments/config/config.example.yaml backend/config.yaml
# 编辑 config.yaml 配置 LLM API Keys
```

4. **运行后端**

```bash
cd backend
go mod download
go run ./cmd/api
```

5. **运行前端**

```bash
cd frontend
npm install
npm run dev
```

6. **访问应用**

- 前端: http://localhost:5173
- API: http://localhost:8080

## 📁 项目结构

```
PaperBeginner/
├── backend/                 # Go 后端
│   ├── cmd/                 # 入口点 (api, worker)
│   ├── internal/            # 内部包
│   │   ├── agent/           # AI Agent 实现
│   │   ├── config/          # 配置管理
│   │   ├── domain/          # 领域模型
│   │   ├── handler/         # HTTP 处理器
│   │   ├── llm/             # LLM 适配层
│   │   ├── repository/      # 数据访问层
│   │   ├── service/         # 业务逻辑
│   │   └── worker/          # 异步任务
│   ├── pkg/                 # 公共包
│   └── migrations/          # 数据库迁移
├── frontend/                # React 前端
│   └── src/
│       ├── components/      # UI 组件
│       ├── pages/           # 页面
│       ├── services/        # API 调用
│       └── stores/          # 状态管理
├── deployments/             # 部署配置
└── docker-compose.yml       # 容器编排
```

## 🔧 配置说明

### LLM 支持

- **OpenAI**: GPT-4o, GPT-4-turbo
- **Anthropic Claude**: Claude 3.5 Sonnet
- **DeepSeek**: DeepSeek Chat
- **Ollama**: 本地模型 (Llama, Mistral 等)

### CCF 领域分类

1. 计算机体系结构/并行与分布计算/存储系统
2. 计算机网络
3. 网络与信息安全
4. 软件工程/系统软件/程序设计语言
5. 数据库/数据挖掘/内容检索
6. 计算机科学理论
7. 计算机图形学与多媒体
8. 人工智能
9. 人机交互与普适计算
10. 交叉/综合/新兴

## 📜 许可证

本项目采用 [AGPL-3.0](LICENSE) 许可证。

---
