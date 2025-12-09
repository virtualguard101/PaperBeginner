package llm

import (
	"context"
	"errors"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

// OpenAIProvider implements Provider for OpenAI
type OpenAIProvider struct {
	config ProviderConfig
	llm    *openai.LLM
}

// NewOpenAIProvider creates a new OpenAI provider
func NewOpenAIProvider(config ProviderConfig) (*OpenAIProvider, error) {
	if config.APIKey == "" {
		return &OpenAIProvider{config: config}, nil
	}

	opts := []openai.Option{
		openai.WithToken(config.APIKey),
	}

	if config.Model != "" {
		opts = append(opts, openai.WithModel(config.Model))
	}

	if config.BaseURL != "" {
		opts = append(opts, openai.WithBaseURL(config.BaseURL))
	}

	llm, err := openai.New(opts...)
	if err != nil {
		return nil, err
	}

	return &OpenAIProvider{
		config: config,
		llm:    llm,
	}, nil
}

// Name returns the provider name
func (p *OpenAIProvider) Name() string {
	return "openai"
}

// IsAvailable returns true if the provider is configured
func (p *OpenAIProvider) IsAvailable() bool {
	return p.config.APIKey != "" && p.llm != nil
}

// Complete generates a completion
func (p *OpenAIProvider) Complete(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error) {
	if !p.IsAvailable() {
		return nil, ErrProviderNotConfigured
	}

	opts := []llms.CallOption{}
	if req.MaxTokens > 0 {
		opts = append(opts, llms.WithMaxTokens(req.MaxTokens))
	}
	if req.Temperature > 0 {
		opts = append(opts, llms.WithTemperature(req.Temperature))
	}
	if len(req.Stop) > 0 {
		opts = append(opts, llms.WithStopWords(req.Stop))
	}

	result, err := llms.GenerateFromSinglePrompt(ctx, p.llm, req.Prompt, opts...)
	if err != nil {
		return nil, err
	}

	return &CompletionResponse{
		Text:         result,
		FinishReason: "stop",
	}, nil
}

// Chat generates a chat completion
func (p *OpenAIProvider) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	if !p.IsAvailable() {
		return nil, ErrProviderNotConfigured
	}

	messages := make([]llms.MessageContent, len(req.Messages))
	for i, msg := range req.Messages {
		role := llms.ChatMessageTypeHuman
		switch msg.Role {
		case "system":
			role = llms.ChatMessageTypeSystem
		case "assistant":
			role = llms.ChatMessageTypeAI
		case "user":
			role = llms.ChatMessageTypeHuman
		}
		messages[i] = llms.MessageContent{
			Role:  role,
			Parts: []llms.ContentPart{llms.TextContent{Text: msg.Content}},
		}
	}

	opts := []llms.CallOption{}
	if req.MaxTokens > 0 {
		opts = append(opts, llms.WithMaxTokens(req.MaxTokens))
	}
	if req.Temperature > 0 {
		opts = append(opts, llms.WithTemperature(req.Temperature))
	}

	resp, err := p.llm.GenerateContent(ctx, messages, opts...)
	if err != nil {
		return nil, err
	}

	if len(resp.Choices) == 0 {
		return nil, errors.New("no response from OpenAI")
	}

	return &ChatResponse{
		Message: Message{
			Role:    "assistant",
			Content: resp.Choices[0].Content,
		},
		FinishReason: string(resp.Choices[0].StopReason),
		Usage: Usage{
			PromptTokens:     0, // Not available in langchaingo
			CompletionTokens: 0,
			TotalTokens:      0,
		},
	}, nil
}

// Embed generates embeddings
func (p *OpenAIProvider) Embed(ctx context.Context, text string) ([]float32, error) {
	if !p.IsAvailable() {
		return nil, ErrProviderNotConfigured
	}

	// Create embeddings using OpenAI embeddings model
	embedder, err := openai.New(
		openai.WithToken(p.config.APIKey),
		openai.WithEmbeddingModel("text-embedding-3-small"),
	)
	if err != nil {
		return nil, err
	}

	embeddings, err := embedder.CreateEmbedding(ctx, []string{text})
	if err != nil {
		return nil, err
	}

	if len(embeddings) == 0 {
		return nil, errors.New("no embeddings returned")
	}

	return embeddings[0], nil
}

// GetMaxTokens returns the maximum tokens
func (p *OpenAIProvider) GetMaxTokens() int {
	switch p.config.Model {
	case "gpt-4o", "gpt-4o-mini":
		return 128000
	case "gpt-4-turbo":
		return 128000
	case "gpt-4":
		return 8192
	case "gpt-3.5-turbo":
		return 16385
	default:
		return 8192
	}
}

// GetEmbeddingDimension returns the embedding dimension
func (p *OpenAIProvider) GetEmbeddingDimension() int {
	return 1536 // text-embedding-3-small default
}

