# PaperBeginer

## 🎯 Introduction

PaperBeginer is an AI Agent-powered academic research guidance platform designed for newcomers in computer science. It provides intelligent tools to help users:

- 📊 **Track Research Trends** - Monitor GitHub trending projects and CCF top-tier conference papers

- 📚 **Plan Learning Paths** - AI-generated personalized learning paths with top university resources

- 📝 **Analyze Papers** - Upload papers for intelligent summaries and methodology analysis

- 📖 **Write Literature Reviews** - AI-assisted high-quality review generation with scoring

## 🚀 Quick Start

```bash
# Clone the repository
git clone https://github.com/virtualguard/PaperBeginer.git
cd PaperBeginer

# Start infrastructure
make dev-infra

# Configure
cp deployments/config/config.example.yaml backend/config.yaml

# Run backend
cd backend && go run ./cmd/api

# Run frontend (in another terminal)
cd frontend && npm install && npm run dev
```

Visit http://localhost:5173 to access the application.

## 📄 License

This project is licensed under [AGPL-3.0](LICENSE).

---

<div align="center">

**Made with ❤️ for Academic Researchers**

</div>

