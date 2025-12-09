package llm

import (
	"context"
	"errors"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

// OllamaProvider implements Provider for local Ollama models
type OllamaProvider struct {
	config ProviderConfig
	llm    *ollama.LLM
}

// NewOllamaProvider creates a new Ollama provider
func NewOllamaProvider(config ProviderConfig) (*OllamaProvider, error) {
	host := config.BaseURL
	if host == "" {
		host = "http://localhost:11434"
	}

	model := config.Model
	if model == "" {
		model = "llama3.1"
	}

	llm, err := ollama.New(
		ollama.WithServerURL(host),
		ollama.WithModel(model),
	)
	if err != nil {
		return &OllamaProvider{config: config}, nil // Return without LLM if connection fails
	}

	return &OllamaProvider{
		config: config,
		llm:    llm,
	}, nil
}

// Name returns the provider name
func (p *OllamaProvider) Name() string {
	return "ollama"
}

// IsAvailable returns true if Ollama is available
func (p *OllamaProvider) IsAvailable() bool {
	return p.llm != nil
}

// Complete generates a completion
func (p *OllamaProvider) Complete(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error) {
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
func (p *OllamaProvider) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
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
		return nil, errors.New("no response from Ollama")
	}

	return &ChatResponse{
		Message: Message{
			Role:    "assistant",
			Content: resp.Choices[0].Content,
		},
		FinishReason: string(resp.Choices[0].StopReason),
		Usage:        Usage{},
	}, nil
}

// Embed generates embeddings using Ollama
func (p *OllamaProvider) Embed(ctx context.Context, text string) ([]float32, error) {
	if !p.IsAvailable() {
		return nil, ErrProviderNotConfigured
	}

	// Use nomic-embed-text or mxbai-embed-large for embeddings
	embedder, err := ollama.New(
		ollama.WithServerURL(p.config.BaseURL),
		ollama.WithModel("nomic-embed-text"),
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
func (p *OllamaProvider) GetMaxTokens() int {
	// Most Ollama models support 4096-8192 context
	return 8192
}

// GetEmbeddingDimension returns the embedding dimension
func (p *OllamaProvider) GetEmbeddingDimension() int {
	return 768 // nomic-embed-text default dimension
}

