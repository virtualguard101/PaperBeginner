package llm

import (
	"context"
	"errors"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/anthropic"
)

// ClaudeProvider implements Provider for Anthropic Claude
type ClaudeProvider struct {
	config ProviderConfig
	llm    *anthropic.LLM
}

// NewClaudeProvider creates a new Claude provider
func NewClaudeProvider(config ProviderConfig) (*ClaudeProvider, error) {
	if config.APIKey == "" {
		return &ClaudeProvider{config: config}, nil
	}

	opts := []anthropic.Option{
		anthropic.WithToken(config.APIKey),
	}

	if config.Model != "" {
		opts = append(opts, anthropic.WithModel(config.Model))
	}

	llm, err := anthropic.New(opts...)
	if err != nil {
		return nil, err
	}

	return &ClaudeProvider{
		config: config,
		llm:    llm,
	}, nil
}

// Name returns the provider name
func (p *ClaudeProvider) Name() string {
	return "anthropic"
}

// IsAvailable returns true if the provider is configured
func (p *ClaudeProvider) IsAvailable() bool {
	return p.config.APIKey != "" && p.llm != nil
}

// Complete generates a completion
func (p *ClaudeProvider) Complete(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error) {
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
func (p *ClaudeProvider) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
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
		return nil, errors.New("no response from Claude")
	}

	return &ChatResponse{
		Message: Message{
			Role:    "assistant",
			Content: resp.Choices[0].Content,
		},
		FinishReason: string(resp.Choices[0].StopReason),
		Usage: Usage{
			TotalTokens: 0, // Not available in langchaingo
		},
	}, nil
}

// Embed generates embeddings (Claude doesn't support embeddings natively)
func (p *ClaudeProvider) Embed(ctx context.Context, text string) ([]float32, error) {
	return nil, errors.New("Claude does not support embeddings, use OpenAI or another provider")
}

// GetMaxTokens returns the maximum tokens
func (p *ClaudeProvider) GetMaxTokens() int {
	switch p.config.Model {
	case "claude-3-5-sonnet-20241022", "claude-3-5-sonnet-latest":
		return 200000
	case "claude-3-opus-20240229":
		return 200000
	case "claude-3-sonnet-20240229":
		return 200000
	case "claude-3-haiku-20240307":
		return 200000
	default:
		return 100000
	}
}

// GetEmbeddingDimension returns the embedding dimension
func (p *ClaudeProvider) GetEmbeddingDimension() int {
	return 0 // Claude doesn't support embeddings
}

