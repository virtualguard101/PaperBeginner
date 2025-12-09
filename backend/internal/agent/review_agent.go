package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/virtualguard/PaperBeginer/internal/domain"
	"github.com/virtualguard/PaperBeginer/internal/llm"
)

// ReviewAgent handles literature review generation and scoring
type ReviewAgent struct {
	router *llm.Router
}

// NewReviewAgent creates a new ReviewAgent
func NewReviewAgent(router *llm.Router) *ReviewAgent {
	return &ReviewAgent{router: router}
}

// GenerateReview generates a literature review from multiple papers
func (a *ReviewAgent) GenerateReview(ctx context.Context, input *ReviewInput) (*ReviewResult, error) {
	systemPrompt := `You are an expert academic writer specializing in literature reviews.
Your task is to synthesize multiple research papers into a comprehensive literature review.

A good literature review should:
1. Provide context and background
2. Organize papers thematically, not just chronologically
3. Identify patterns, trends, and gaps in the research
4. Compare and contrast different approaches
5. Synthesize findings to create new insights
6. Be well-structured with clear sections
7. Use proper academic writing style

Structure your review with:
- Abstract (150-200 words)
- Introduction (background, scope, objectives)
- Main body (organized by themes/topics)
- Discussion (synthesis, gaps, future directions)
- Conclusion

Write in formal academic English.`

	// Build papers summary
	var papersSummary strings.Builder
	for i, paper := range input.Papers {
		papersSummary.WriteString(fmt.Sprintf("\n--- Paper %d ---\n", i+1))
		papersSummary.WriteString(fmt.Sprintf("Title: %s\n", paper.Title))
		papersSummary.WriteString(fmt.Sprintf("Authors: %s\n", strings.Join(paper.Authors, ", ")))
		if paper.Year != "" {
			papersSummary.WriteString(fmt.Sprintf("Year: %s\n", paper.Year))
		}
		papersSummary.WriteString(fmt.Sprintf("Abstract: %s\n", paper.Abstract))
		if paper.Summary != "" {
			papersSummary.WriteString(fmt.Sprintf("Summary: %s\n", paper.Summary))
		}
	}

	styleHint := ""
	switch input.Style {
	case "academic":
		styleHint = "Write in formal academic style suitable for journal publication."
	case "summary":
		styleHint = "Write a concise summary-style review focusing on key findings."
	case "comprehensive":
		styleHint = "Write a comprehensive review covering all aspects in depth."
	}

	userPrompt := fmt.Sprintf(`Generate a literature review on the topic: %s

Papers to review:
%s

Additional context: %s

%s

Generate a complete literature review with abstract, introduction, main body, discussion, and conclusion.`,
		input.Title, papersSummary.String(), input.Context, styleHint)

	resp, err := a.router.Chat(ctx, llm.TaskReviewWriting, &llm.ChatRequest{
		Messages: []llm.Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Temperature: 0.7,
		MaxTokens:   8000,
	}, nil)

	if err != nil {
		return nil, err
	}

	// Extract abstract from the generated content
	abstract := extractAbstract(resp.Message.Content)

	return &ReviewResult{
		Title:    input.Title,
		Abstract: abstract,
		Content:  resp.Message.Content,
	}, nil
}

// ScoreReview evaluates a literature review against academic standards
func (a *ReviewAgent) ScoreReview(ctx context.Context, review *ScoreReviewInput) (*domain.ReviewScore, error) {
	systemPrompt := `You are an expert peer reviewer and academic writing evaluator.
Your task is to score a literature review against standard academic criteria.

Evaluation Criteria (each scored 0-100):
1. Comprehensiveness (20%): Coverage of relevant literature, completeness
2. Organization (15%): Structure, flow, logical organization
3. Critical Analysis (20%): Depth of analysis, synthesis, comparison
4. Writing Quality (15%): Clarity, grammar, academic style
5. Relevance (15%): Focus on topic, appropriate scope
6. Contribution (15%): New insights, identification of gaps, future directions

Respond in JSON format:
{
  "overall_score": <0-100>,
  "criteria": [
    {
      "name": "<criterion name>",
      "score": <0-100>,
      "weight": <0-1>,
      "description": "<what this criterion measures>",
      "feedback": "<specific feedback for this criterion>"
    }
  ],
  "strengths": ["<strength 1>", "<strength 2>"],
  "weaknesses": ["<weakness 1>", "<weakness 2>"],
  "suggestions": ["<suggestion 1>", "<suggestion 2>"]
}`

	userPrompt := fmt.Sprintf(`Please evaluate the following literature review:

Title: %s

Content:
%s

Provide detailed scoring and feedback.`, review.Title, review.Content)

	resp, err := a.router.Chat(ctx, llm.TaskReviewScoring, &llm.ChatRequest{
		Messages: []llm.Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Temperature: 0.3,
		MaxTokens:   2000,
	}, nil)

	if err != nil {
		return nil, err
	}

	var scoreResult ScoreResult
	if err := json.Unmarshal([]byte(resp.Message.Content), &scoreResult); err != nil {
		return nil, fmt.Errorf("failed to parse score response: %w", err)
	}

	// Convert to domain model
	criteria := make([]domain.ScoreCriterion, len(scoreResult.Criteria))
	for i, c := range scoreResult.Criteria {
		criteria[i] = domain.ScoreCriterion{
			Name:        c.Name,
			Score:       c.Score,
			Weight:      c.Weight,
			Description: c.Description,
			Feedback:    c.Feedback,
		}
	}

	return &domain.ReviewScore{
		OverallScore: scoreResult.OverallScore,
		Criteria:     criteria,
		Strengths:    scoreResult.Strengths,
		Weaknesses:   scoreResult.Weaknesses,
		Suggestions:  scoreResult.Suggestions,
	}, nil
}

// extractAbstract extracts the abstract from a literature review
func extractAbstract(content string) string {
	// Simple extraction - look for "Abstract" section
	lines := strings.Split(content, "\n")
	inAbstract := false
	var abstract strings.Builder

	for _, line := range lines {
		lowerLine := strings.ToLower(strings.TrimSpace(line))
		
		if strings.Contains(lowerLine, "abstract") && len(line) < 50 {
			inAbstract = true
			continue
		}
		
		if inAbstract {
			if strings.Contains(lowerLine, "introduction") || 
			   strings.Contains(lowerLine, "1.") ||
			   line == "" && abstract.Len() > 100 {
				break
			}
			if strings.TrimSpace(line) != "" {
				abstract.WriteString(line)
				abstract.WriteString(" ")
			}
		}
	}

	result := strings.TrimSpace(abstract.String())
	if result == "" {
		// Fallback: use first 200 words
		words := strings.Fields(content)
		if len(words) > 200 {
			words = words[:200]
		}
		result = strings.Join(words, " ") + "..."
	}

	return result
}

// ReviewInput represents input for review generation
type ReviewInput struct {
	Title   string
	Papers  []PaperSummary
	Context string
	Style   string // academic, summary, comprehensive
}

// PaperSummary represents a summarized paper for review
type PaperSummary struct {
	Title    string
	Authors  []string
	Year     string
	Abstract string
	Summary  string
}

// ReviewResult represents a generated review
type ReviewResult struct {
	Title    string
	Abstract string
	Content  string
}

// ScoreReviewInput represents input for review scoring
type ScoreReviewInput struct {
	Title   string
	Content string
}

// ScoreResult represents the score response from LLM
type ScoreResult struct {
	OverallScore float64              `json:"overall_score"`
	Criteria     []CriterionResult    `json:"criteria"`
	Strengths    []string             `json:"strengths"`
	Weaknesses   []string             `json:"weaknesses"`
	Suggestions  []string             `json:"suggestions"`
}

// CriterionResult represents a scoring criterion
type CriterionResult struct {
	Name        string  `json:"name"`
	Score       float64 `json:"score"`
	Weight      float64 `json:"weight"`
	Description string  `json:"description"`
	Feedback    string  `json:"feedback"`
}

