package llm

import (
	"context"
	"errors"
)

// Common errors
var (
	ErrProviderNotConfigured = errors.New("LLM provider not configured")
	ErrNoAvailableProvider   = errors.New("no available LLM provider")
	ErrRateLimited           = errors.New("rate limited")
	ErrContextTooLong        = errors.New("context too long")
)

// Provider defines the interface for LLM providers
type Provider interface {
	// Name returns the provider name
	Name() string

	// IsAvailable returns true if the provider is configured and available
	IsAvailable() bool

	// Complete generates a completion for the given prompt
	Complete(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error)

	// Chat generates a chat completion
	Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error)

	// Embed generates embeddings for the given text
	Embed(ctx context.Context, text string) ([]float32, error)

	// GetMaxTokens returns the maximum tokens supported by the model
	GetMaxTokens() int

	// GetEmbeddingDimension returns the embedding dimension
	GetEmbeddingDimension() int
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`    // system, user, assistant
	Content string `json:"content"`
}

// CompletionRequest represents a completion request
type CompletionRequest struct {
	Prompt      string   `json:"prompt"`
	MaxTokens   int      `json:"max_tokens,omitempty"`
	Temperature float64  `json:"temperature,omitempty"`
	TopP        float64  `json:"top_p,omitempty"`
	Stop        []string `json:"stop,omitempty"`
}

// CompletionResponse represents a completion response
type CompletionResponse struct {
	Text         string `json:"text"`
	FinishReason string `json:"finish_reason"`
	TokensUsed   int    `json:"tokens_used"`
}

// ChatRequest represents a chat completion request
type ChatRequest struct {
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
	TopP        float64   `json:"top_p,omitempty"`
	Stop        []string  `json:"stop,omitempty"`
	Stream      bool      `json:"stream,omitempty"`
}

// ChatResponse represents a chat completion response
type ChatResponse struct {
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
	Usage        Usage   `json:"usage"`
}

// Usage represents token usage
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// TaskType defines the type of AI task
type TaskType string

const (
	TaskTrendAnalysis   TaskType = "trend_analysis"
	TaskPaperParsing    TaskType = "paper_parsing"
	TaskPaperSummary    TaskType = "paper_summary"
	TaskLearningPath    TaskType = "learning_path"
	TaskReviewWriting   TaskType = "review_writing"
	TaskReviewScoring   TaskType = "review_scoring"
	TaskCategoryClassify TaskType = "category_classify"
	TaskEmbedding       TaskType = "embedding"
)

// ProviderPreference defines user's provider preferences
type ProviderPreference struct {
	PreferredProvider string            `json:"preferred_provider,omitempty"`
	APIKeys           map[string]string `json:"api_keys,omitempty"`
	FallbackEnabled   bool              `json:"fallback_enabled"`
}

// ProviderConfig contains common provider configuration
type ProviderConfig struct {
	APIKey     string
	Model      string
	BaseURL    string
	MaxRetries int
	Timeout    int // seconds
}

