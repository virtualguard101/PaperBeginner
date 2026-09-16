from app.services.llm import LLMResult, chat_json


def analyze_paper(title: str, text: str) -> LLMResult:
    excerpt = text[:12000] if text else "(no extractable text)"
    system = (
        "You are an academic paper analyst. Reply with JSON only, keys: "
        "title, authors (array of strings), abstract, summary, methodology, contributions. "
        "Write in Chinese unless the paper is clearly English-only; summaries may be Chinese."
    )
    user = f"Paper title hint: {title}\n\nPaper text:\n{excerpt}"
    fallback = {
        "title": title or "未命名论文",
        "authors": [],
        "abstract": excerpt[:400] or "未能从 PDF 提取摘要，以下为占位分析。",
        "summary": "（离线模板）本文讨论了相关方法与实验。请配置 LLM_API_KEY 以获得真实分析。",
        "methodology": "（离线模板）方法论部分未能由模型生成。请检查 PDF 文本提取与 API Key。",
        "contributions": "（离线模板）贡献点占位：提出问题、给出方法、进行实验验证。",
    }
    return chat_json(system, user, fallback)


def generate_trend_report(items: list[dict], period: str) -> LLMResult:
    lines = "\n".join(f"- [{i.get('source')}] {i.get('title')}: {i.get('description', '')[:120]}" for i in items[:20])
    system = (
        "You write weekly CS research trend briefs in Chinese. "
        "JSON keys: summary (string), highlights (array of {title, description, sources, importance 1-5})."
    )
    user = f"Period: {period}\nItems:\n{lines or '(empty)'}"
    fallback = {
        "summary": f"{period} 热点速览（离线模板）：当前列表包含 GitHub 热门项目与 CCF 代表工作，建议结合人工智能与系统方向跟进。",
        "highlights": [
            {
                "title": "开源项目热度集中",
                "description": "GitHub Trending 反映工程落地兴趣。",
                "sources": ["github"],
                "importance": 4,
            },
            {
                "title": "顶会方向可对照学习",
                "description": "CCF 种子条目可用于入门选题对照。",
                "sources": ["ccf"],
                "importance": 3,
            },
        ],
    }
    return chat_json(system, user, fallback)


def generate_learning_path(category_name: str, difficulty: str, extras: str) -> LLMResult:
    system = (
        "You design CS learning paths for beginners in Chinese. JSON keys: "
        "title, description, estimated_time, prerequisites (string array), "
        "stages (array of {order, title, description, duration, skills (string array), "
        "resources (array of {type, title, url, provider, language, is_free, description})}). "
        "Use 3 or 4 stages. Prefer real URLs: MIT OCW, Stanford, csdiy.wiki, official docs."
    )
    user = f"Category: {category_name}\nDifficulty: {difficulty}\nNotes: {extras}"
    fallback = {
        "title": f"{category_name} 入门学习路线",
        "description": f"面向 {difficulty} 学习者的分阶段路径（离线模板）。配置 LLM_API_KEY 后可生成更贴合的资源。",
        "estimated_time": "8-12 周",
        "prerequisites": ["离散数学基础", "一门编程语言"],
        "stages": [
            {
                "order": 1,
                "title": "基础概念",
                "description": "建立该方向的核心术语与问题意识。",
                "duration": "2-3 周",
                "skills": ["文献检索", "基础概念"],
                "resources": [
                    {
                        "type": "course",
                        "title": "MIT OpenCourseWare",
                        "url": "https://ocw.mit.edu/",
                        "provider": "MIT",
                        "language": "en",
                        "is_free": True,
                        "description": "公开课入口",
                    }
                ],
            },
            {
                "order": 2,
                "title": "核心课程",
                "description": "跟一门体系化公开课。",
                "duration": "4-6 周",
                "skills": ["动手实现"],
                "resources": [
                    {
                        "type": "tutorial",
                        "title": "CS DIY",
                        "url": "https://csdiy.wiki/",
                        "provider": "csdiy.wiki",
                        "language": "zh",
                        "is_free": True,
                        "description": "自学路线索引",
                    }
                ],
            },
            {
                "order": 3,
                "title": "读论文与复现",
                "description": "选 2-3 篇领域代表论文精读。",
                "duration": "3-4 周",
                "skills": ["论文精读"],
                "resources": [
                    {
                        "type": "documentation",
                        "title": "Papers With Code",
                        "url": "https://paperswithcode.com/",
                        "provider": "PapersWithCode",
                        "language": "en",
                        "is_free": True,
                    }
                ],
            },
        ],
    }
    return chat_json(system, user, fallback)


def generate_review(title: str, style: str, paper_blobs: str) -> LLMResult:
    system = (
        "You write a short Chinese literature review in Markdown. JSON keys: "
        "abstract, content. Style hint: " + style
    )
    user = f"Review title: {title}\nSource papers:\n{paper_blobs[:14000]}"
    fallback = {
        "abstract": f"本文围绕「{title}」对所选论文做简要综述（离线模板）。",
        "content": f"# {title}\n\n## 引言\n本文基于用户选择的论文生成占位综述。\n\n## 主要工作\n{paper_blobs[:800]}\n\n## 小结\n请配置 LLM_API_KEY 以获得完整综述。\n",
    }
    return chat_json(system, user, fallback)


def score_review(title: str, content: str) -> LLMResult:
    system = (
        "You score a Chinese academic literature review. JSON keys: "
        "overall_score (0-100), criteria (array of {name, score, weight, description, feedback}), "
        "strengths (string array), weaknesses (string array), suggestions (string array)."
    )
    user = f"Title: {title}\nContent:\n{content[:8000]}"
    fallback = {
        "overall_score": 72,
        "criteria": [
            {"name": "全面性", "score": 70, "weight": 0.25, "description": "覆盖广度", "feedback": "覆盖有限，可增加对照工作。"},
            {"name": "组织结构", "score": 78, "weight": 0.25, "description": "结构清晰度", "feedback": "章节清楚。"},
            {"name": "批判性分析", "score": 68, "weight": 0.25, "description": "评价深度", "feedback": "批判不足。"},
            {"name": "写作质量", "score": 74, "weight": 0.25, "description": "语言表达", "feedback": "可读性尚可。"},
        ],
        "strengths": ["结构完整", "主题集中"],
        "weaknesses": ["缺少定量对比"],
        "suggestions": ["补充未来工作与开放问题"],
    }
    return chat_json(system, user, fallback)
