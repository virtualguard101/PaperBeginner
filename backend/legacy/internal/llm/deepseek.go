package llm

import (
	"context"
	"errors"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

// DeepSeekProvider implements Provider for DeepSeek (OpenAI-compatible API)
type DeepSeekProvider struct {
	config ProviderConfig
	llm    *openai.LLM
}

// NewDeepSeekProvider creates a new DeepSeek provider
func NewDeepSeekProvider(config ProviderConfig) (*DeepSeekProvider, error) {
	if config.APIKey == "" {
		return &DeepSeekProvider{config: config}, nil
	}

	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = "https://api.deepseek.com/v1"
	}

	model := config.Model
	if model == "" {
		model = "deepseek-chat"
	}

	llm, err := openai.New(
		openai.WithToken(config.APIKey),
		openai.WithBaseURL(baseURL),
		openai.WithModel(model),
	)
	if err != nil {
		return nil, err
	}

	return &DeepSeekProvider{
		config: config,
		llm:    llm,
	}, nil
}

// Name returns the provider name
func (p *DeepSeekProvider) Name() string {
	return "deepseek"
}

// IsAvailable returns true if the provider is configured
func (p *DeepSeekProvider) IsAvailable() bool {
	return p.config.APIKey != "" && p.llm != nil
}

// Complete generates a completion
func (p *DeepSeekProvider) Complete(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error) {
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
func (p *DeepSeekProvider) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
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
		return nil, errors.New("no response from DeepSeek")
	}

	return &ChatResponse{
		Message: Message{
			Role:    "assistant",
			Content: resp.Choices[0].Content,
		},
		FinishReason: string(resp.Choices[0].StopReason),
		Usage: Usage{
			TotalTokens: 0,
		},
	}, nil
}

// Embed generates embeddings (DeepSeek doesn't support embeddings yet)
func (p *DeepSeekProvider) Embed(ctx context.Context, text string) ([]float32, error) {
	return nil, errors.New("DeepSeek does not support embeddings, use OpenAI or another provider")
}

// GetMaxTokens returns the maximum tokens
func (p *DeepSeekProvider) GetMaxTokens() int {
	switch p.config.Model {
	case "deepseek-chat":
		return 64000
	case "deepseek-coder":
		return 64000
	default:
		return 32000
	}
}

// GetEmbeddingDimension returns the embedding dimension
func (p *DeepSeekProvider) GetEmbeddingDimension() int {
	return 0 // DeepSeek doesn't support embeddings
}

