package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/virtualguard/PaperBeginner/internal/domain"
	"github.com/virtualguard/PaperBeginner/internal/llm"
)

// TrendAgent handles trend analysis and categorization
type TrendAgent struct {
	router *llm.Router
}

// NewTrendAgent creates a new TrendAgent
func NewTrendAgent(router *llm.Router) *TrendAgent {
	return &TrendAgent{router: router}
}

// ClassifyCCFCategory classifies an item into CCF categories
func (a *TrendAgent) ClassifyCCFCategory(ctx context.Context, item *ClassifyInput) (*ClassifyResult, error) {
	systemPrompt := `You are an expert in computer science research classification.
Your task is to classify research topics and projects into CCF (China Computer Federation) categories.

Available CCF Categories:
1. 计算机体系结构/并行与分布计算/存储系统 (Computer Architecture/Parallel and Distributed Computing/Storage Systems)
2. 计算机网络 (Computer Networks)
3. 网络与信息安全 (Network and Information Security)
4. 软件工程/系统软件/程序设计语言 (Software Engineering/System Software/Programming Languages)
5. 数据库/数据挖掘/内容检索 (Database/Data Mining/Content Retrieval)
6. 计算机科学理论 (Computer Science Theory)
7. 计算机图形学与多媒体 (Computer Graphics and Multimedia)
8. 人工智能 (Artificial Intelligence)
9. 人机交互与普适计算 (Human-Computer Interaction and Ubiquitous Computing)
10. 交叉/综合/新兴 (Interdisciplinary/Comprehensive/Emerging)

Respond in JSON format with the following structure:
{
  "category_id": <number 1-10>,
  "category_name": "<category name in Chinese>",
  "confidence": <0.0-1.0>,
  "reasoning": "<brief explanation>"
}`

	userPrompt := fmt.Sprintf(`Classify the following item:

Title: %s
Description: %s
Topics/Keywords: %v
Language/Framework: %s

Provide classification in JSON format.`, item.Title, item.Description, item.Topics, item.Language)

	resp, err := a.router.Chat(ctx, llm.TaskCategoryClassify, &llm.ChatRequest{
		Messages: []llm.Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Temperature: 0.3,
		MaxTokens:   500,
	}, nil)

	if err != nil {
		return nil, err
	}

	var result ClassifyResult
	if err := json.Unmarshal([]byte(resp.Message.Content), &result); err != nil {
		// Try to extract JSON from response
		result = ClassifyResult{
			CategoryID:   10, // Default to Interdisciplinary
			CategoryName: "交叉/综合/新兴",
			Confidence:   0.5,
			Reasoning:    "Unable to parse classification response",
		}
	}

	return &result, nil
}

// GenerateTrendReport generates a trend report for a category
func (a *TrendAgent) GenerateTrendReport(ctx context.Context, items []*domain.TrendingItem, category *domain.CCFCategory, period string) (*TrendReportResult, error) {
	// Build items summary
	var itemsSummary string
	for i, item := range items {
		if i >= 20 { // Limit to top 20 items
			break
		}
		itemsSummary += fmt.Sprintf("- %s: %s (Score: %.2f)\n", item.Title, item.Description, item.TrendScore)
	}

	systemPrompt := `You are an expert in computer science research trends.
Your task is to analyze trending topics and generate insightful trend reports for academic researchers.

The report should:
1. Identify key themes and emerging technologies
2. Highlight significant projects and papers
3. Analyze potential research directions
4. Be written in a professional academic style
5. Be helpful for researchers entering the field

Respond in JSON format:
{
  "summary": "<2-3 paragraph executive summary>",
  "highlights": [
    {
      "title": "<highlight title>",
      "description": "<detailed description>",
      "importance": <1-5>,
      "sources": ["<source references>"]
    }
  ],
  "emerging_topics": ["<topic1>", "<topic2>"],
  "research_opportunities": ["<opportunity1>", "<opportunity2>"]
}`

	categoryName := "General"
	if category != nil {
		categoryName = category.Name
	}

	userPrompt := fmt.Sprintf(`Generate a trend report for the following:

Category: %s
Period: %s

Trending Items:
%s

Please analyze these items and generate a comprehensive trend report.`, categoryName, period, itemsSummary)

	resp, err := a.router.Chat(ctx, llm.TaskTrendAnalysis, &llm.ChatRequest{
		Messages: []llm.Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Temperature: 0.7,
		MaxTokens:   2000,
	}, nil)

	if err != nil {
		return nil, err
	}

	var result TrendReportResult
	if err := json.Unmarshal([]byte(resp.Message.Content), &result); err != nil {
		return &TrendReportResult{
			Summary: resp.Message.Content,
		}, nil
	}

	return &result, nil
}

// ClassifyInput represents input for classification
type ClassifyInput struct {
	Title       string
	Description string
	Topics      []string
	Language    string
}

// ClassifyResult represents classification result
type ClassifyResult struct {
	CategoryID   int     `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Confidence   float64 `json:"confidence"`
	Reasoning    string  `json:"reasoning"`
}

// TrendReportResult represents a generated trend report
type TrendReportResult struct {
	Summary               string            `json:"summary"`
	Highlights            []TrendHighlight  `json:"highlights"`
	EmergingTopics        []string          `json:"emerging_topics"`
	ResearchOpportunities []string          `json:"research_opportunities"`
}

// TrendHighlight represents a highlight in the report
type TrendHighlight struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Importance  int      `json:"importance"`
	Sources     []string `json:"sources"`
}

