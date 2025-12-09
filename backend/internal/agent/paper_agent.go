package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/virtualguard/PaperBeginer/internal/domain"
	"github.com/virtualguard/PaperBeginer/internal/llm"
)

// PaperAgent handles paper analysis tasks
type PaperAgent struct {
	router *llm.Router
}

// NewPaperAgent creates a new PaperAgent
func NewPaperAgent(router *llm.Router) *PaperAgent {
	return &PaperAgent{router: router}
}

// AnalyzePaper performs comprehensive paper analysis
func (a *PaperAgent) AnalyzePaper(ctx context.Context, paper *PaperInput, analysisType domain.AnalysisType) (*PaperAnalysisResult, error) {
	var systemPrompt, userPrompt string

	switch analysisType {
	case domain.AnalysisSummary:
		systemPrompt, userPrompt = a.buildSummaryPrompts(paper)
	case domain.AnalysisMethodology:
		systemPrompt, userPrompt = a.buildMethodologyPrompts(paper)
	case domain.AnalysisContributions:
		systemPrompt, userPrompt = a.buildContributionsPrompts(paper)
	case domain.AnalysisStrengths:
		systemPrompt, userPrompt = a.buildStrengthsPrompts(paper)
	case domain.AnalysisWeaknesses:
		systemPrompt, userPrompt = a.buildWeaknessesPrompts(paper)
	default:
		systemPrompt, userPrompt = a.buildSummaryPrompts(paper)
	}

	resp, err := a.router.Chat(ctx, llm.TaskPaperSummary, &llm.ChatRequest{
		Messages: []llm.Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Temperature: 0.5,
		MaxTokens:   2000,
	}, nil)

	if err != nil {
		return nil, err
	}

	return &PaperAnalysisResult{
		Type:    analysisType,
		Content: resp.Message.Content,
		Usage:   resp.Usage,
	}, nil
}

func (a *PaperAgent) buildSummaryPrompts(paper *PaperInput) (string, string) {
	systemPrompt := `You are an expert academic paper analyst. Your task is to provide a comprehensive summary of research papers.

Your summary should include:
1. Main research question/problem
2. Key methodology used
3. Main findings/results
4. Significance and impact
5. Target audience

Write in clear, academic language suitable for researchers new to the field.`

	userPrompt := fmt.Sprintf(`Please summarize the following paper:

Title: %s
Authors: %s
Abstract: %s

Full Text (excerpt):
%s

Provide a comprehensive summary covering the main aspects of this research.`, 
		paper.Title, paper.Authors, paper.Abstract, truncateText(paper.Content, 10000))

	return systemPrompt, userPrompt
}

func (a *PaperAgent) buildMethodologyPrompts(paper *PaperInput) (string, string) {
	systemPrompt := `You are an expert in research methodology. Your task is to analyze and explain the methodology used in academic papers.

Your analysis should cover:
1. Research design (experimental, observational, theoretical, etc.)
2. Data collection methods
3. Analysis techniques
4. Evaluation metrics
5. Limitations of the methodology

Be specific and technical where appropriate.`

	userPrompt := fmt.Sprintf(`Analyze the methodology of the following paper:

Title: %s
Abstract: %s

Full Text (excerpt):
%s

Provide a detailed analysis of the research methodology.`,
		paper.Title, paper.Abstract, truncateText(paper.Content, 10000))

	return systemPrompt, userPrompt
}

func (a *PaperAgent) buildContributionsPrompts(paper *PaperInput) (string, string) {
	systemPrompt := `You are an expert in evaluating academic contributions. Your task is to identify and explain the key contributions of research papers.

For each contribution, explain:
1. What the contribution is
2. Why it is significant
3. How it advances the field
4. Potential applications or implications

List contributions in order of significance.`

	userPrompt := fmt.Sprintf(`Identify the key contributions of the following paper:

Title: %s
Abstract: %s

Full Text (excerpt):
%s

List and explain the main contributions of this research.`,
		paper.Title, paper.Abstract, truncateText(paper.Content, 10000))

	return systemPrompt, userPrompt
}

func (a *PaperAgent) buildStrengthsPrompts(paper *PaperInput) (string, string) {
	systemPrompt := `You are an expert peer reviewer. Your task is to identify the strengths of academic papers.

Consider:
1. Novelty and originality
2. Technical soundness
3. Clarity of presentation
4. Experimental rigor
5. Practical relevance
6. Theoretical foundations

Be constructive and specific.`

	userPrompt := fmt.Sprintf(`Identify the strengths of the following paper:

Title: %s
Abstract: %s

Full Text (excerpt):
%s

List the main strengths of this research work.`,
		paper.Title, paper.Abstract, truncateText(paper.Content, 10000))

	return systemPrompt, userPrompt
}

func (a *PaperAgent) buildWeaknessesPrompts(paper *PaperInput) (string, string) {
	systemPrompt := `You are an expert peer reviewer. Your task is to identify potential weaknesses and areas for improvement in academic papers.

Consider:
1. Methodological limitations
2. Gaps in evaluation
3. Missing related work
4. Unclear explanations
5. Reproducibility concerns
6. Generalization limitations

Be constructive and suggest improvements where possible.`

	userPrompt := fmt.Sprintf(`Identify potential weaknesses of the following paper:

Title: %s
Abstract: %s

Full Text (excerpt):
%s

List the potential weaknesses and areas for improvement.`,
		paper.Title, paper.Abstract, truncateText(paper.Content, 10000))

	return systemPrompt, userPrompt
}

// ExtractKeywords extracts keywords from a paper
func (a *PaperAgent) ExtractKeywords(ctx context.Context, paper *PaperInput) ([]string, error) {
	systemPrompt := `You are an expert in academic keyword extraction. Extract relevant keywords and key phrases from the given paper that would help in categorization and search.

Return a JSON array of keywords, ordered by relevance. Include:
- Technical terms
- Methods/algorithms
- Application domains
- Key concepts

Return only the JSON array, no explanation.`

	userPrompt := fmt.Sprintf(`Extract keywords from:

Title: %s
Abstract: %s

Return as JSON array: ["keyword1", "keyword2", ...]`,
		paper.Title, paper.Abstract)

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

	var keywords []string
	if err := json.Unmarshal([]byte(resp.Message.Content), &keywords); err != nil {
		return []string{}, nil
	}

	return keywords, nil
}

// PaperInput represents input for paper analysis
type PaperInput struct {
	Title    string
	Authors  string
	Abstract string
	Content  string
}

// PaperAnalysisResult represents paper analysis result
type PaperAnalysisResult struct {
	Type    domain.AnalysisType
	Content string
	Usage   llm.Usage
}

// truncateText truncates text to a maximum length
func truncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen] + "..."
}

