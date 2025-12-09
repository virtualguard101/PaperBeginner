package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/virtualguard/PaperBeginer/internal/domain"
	"github.com/virtualguard/PaperBeginer/internal/llm"
)

// LearningAgent handles learning path generation
type LearningAgent struct {
	router *llm.Router
}

// NewLearningAgent creates a new LearningAgent
func NewLearningAgent(router *llm.Router) *LearningAgent {
	return &LearningAgent{router: router}
}

// GenerateLearningPath generates a personalized learning path
func (a *LearningAgent) GenerateLearningPath(ctx context.Context, input *LearningPathInput) (*domain.LearningPath, error) {
	systemPrompt := `You are an expert computer science educator specializing in guiding new researchers.
Your task is to create comprehensive learning paths for students entering academic research.

When creating a learning path:
1. Start with foundational concepts and gradually increase complexity
2. Include both theoretical knowledge and practical skills
3. Reference high-quality resources:
   - Official documentation
   - Top university courses (MIT, Stanford, CMU, Berkeley)
   - CS自学指南 (csdiy.wiki) resources when applicable
   - Classic textbooks
   - Important research papers
4. Estimate realistic time requirements
5. Include hands-on projects or exercises

Respond in JSON format:
{
  "title": "<learning path title>",
  "description": "<overview of the learning path>",
  "estimated_time": "<total estimated time, e.g., '3-6 months'>",
  "prerequisites": ["<prerequisite 1>", "<prerequisite 2>"],
  "stages": [
    {
      "order": 1,
      "title": "<stage title>",
      "description": "<stage description>",
      "duration": "<estimated duration>",
      "skills": ["<skill 1>", "<skill 2>"],
      "resources": [
        {
          "type": "<documentation|course|book|video|tutorial>",
          "title": "<resource title>",
          "url": "<resource URL>",
          "provider": "<provider name>",
          "language": "<en|zh>",
          "is_free": true,
          "description": "<brief description>"
        }
      ]
    }
  ]
}`

	categoryInfo := ""
	if input.Category != nil {
		categoryInfo = fmt.Sprintf("Field: %s (%s)", input.Category.Name, input.Category.NameEN)
	}

	prerequisitesInfo := ""
	if len(input.Prerequisites) > 0 {
		prerequisitesInfo = fmt.Sprintf("Student's background: %v", input.Prerequisites)
	}

	focusInfo := ""
	if len(input.FocusAreas) > 0 {
		focusInfo = fmt.Sprintf("Areas of interest: %v", input.FocusAreas)
	}

	userPrompt := fmt.Sprintf(`Create a learning path for a student wanting to enter academic research:

%s
Difficulty level: %s
%s
%s

Please create a comprehensive, stage-by-stage learning path with high-quality resources.
Include resources from official documentation, top universities (MIT OCW, Stanford Online, etc.), and csdiy.wiki where relevant.`,
		categoryInfo, input.Difficulty, prerequisitesInfo, focusInfo)

	resp, err := a.router.Chat(ctx, llm.TaskLearningPath, &llm.ChatRequest{
		Messages: []llm.Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Temperature: 0.7,
		MaxTokens:   4000,
	}, nil)

	if err != nil {
		return nil, err
	}

	var result LearningPathResult
	if err := json.Unmarshal([]byte(resp.Message.Content), &result); err != nil {
		return nil, fmt.Errorf("failed to parse learning path response: %w", err)
	}

	// Convert to domain model
	stages := make([]domain.LearningStage, len(result.Stages))
	for i, stage := range result.Stages {
		resources := make([]domain.LearningResource, len(stage.Resources))
		for j, res := range stage.Resources {
			resources[j] = domain.LearningResource{
				Type:        res.Type,
				Title:       res.Title,
				URL:         res.URL,
				Provider:    res.Provider,
				Language:    res.Language,
				IsFree:      res.IsFree,
				Description: res.Description,
			}
		}
		stages[i] = domain.LearningStage{
			Order:       stage.Order,
			Title:       stage.Title,
			Description: stage.Description,
			Duration:    stage.Duration,
			Skills:      stage.Skills,
			Resources:   resources,
		}
	}

	learningPath := &domain.LearningPath{
		Title:         result.Title,
		Description:   result.Description,
		Stages:        stages,
		Prerequisites: result.Prerequisites,
		EstimatedTime: result.EstimatedTime,
		Difficulty:    input.Difficulty,
	}

	if input.Category != nil {
		learningPath.CategoryID = input.Category.ID
	}

	return learningPath, nil
}

// LearningPathInput represents input for learning path generation
type LearningPathInput struct {
	Category      *domain.CCFCategory
	Difficulty    string   // beginner, intermediate, advanced
	Prerequisites []string
	FocusAreas    []string
}

// LearningPathResult represents the parsed learning path from LLM
type LearningPathResult struct {
	Title         string                   `json:"title"`
	Description   string                   `json:"description"`
	EstimatedTime string                   `json:"estimated_time"`
	Prerequisites []string                 `json:"prerequisites"`
	Stages        []LearningStageResult    `json:"stages"`
}

// LearningStageResult represents a stage in the learning path
type LearningStageResult struct {
	Order       int                      `json:"order"`
	Title       string                   `json:"title"`
	Description string                   `json:"description"`
	Duration    string                   `json:"duration"`
	Skills      []string                 `json:"skills"`
	Resources   []LearningResourceResult `json:"resources"`
}

// LearningResourceResult represents a learning resource
type LearningResourceResult struct {
	Type        string `json:"type"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	Provider    string `json:"provider"`
	Language    string `json:"language"`
	IsFree      bool   `json:"is_free"`
	Description string `json:"description"`
}

