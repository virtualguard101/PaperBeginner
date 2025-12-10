package llm

import (
	"context"
	"sync"

	"github.com/virtualguard/PaperBeginner/internal/config"
	"github.com/virtualguard/PaperBeginner/pkg/logger"
	"go.uber.org/zap"
)

// Router manages multiple LLM providers and routes requests
type Router struct {
	providers   map[string]Provider
	config      config.LLMConfig
	mu          sync.RWMutex
	taskMapping map[TaskType]string // Default provider for each task type
}

// NewRouter creates a new LLM router
func NewRouter(cfg config.LLMConfig) (*Router, error) {
	r := &Router{
		providers:   make(map[string]Provider),
		config:      cfg,
		taskMapping: defaultTaskMapping(),
	}

	// Initialize OpenAI provider
	if cfg.OpenAI.APIKey != "" {
		openaiProvider, err := NewOpenAIProvider(ProviderConfig{
			APIKey:  cfg.OpenAI.APIKey,
			Model:   cfg.OpenAI.Model,
			BaseURL: cfg.OpenAI.BaseURL,
		})
		if err != nil {
			logger.Warn("Failed to initialize OpenAI provider", zap.Error(err))
		} else {
			r.providers["openai"] = openaiProvider
		}
	}

	// Initialize Claude provider
	if cfg.Anthropic.APIKey != "" {
		claudeProvider, err := NewClaudeProvider(ProviderConfig{
			APIKey: cfg.Anthropic.APIKey,
			Model:  cfg.Anthropic.Model,
		})
		if err != nil {
			logger.Warn("Failed to initialize Claude provider", zap.Error(err))
		} else {
			r.providers["anthropic"] = claudeProvider
		}
	}

	// Initialize DeepSeek provider
	if cfg.DeepSeek.APIKey != "" {
		deepseekProvider, err := NewDeepSeekProvider(ProviderConfig{
			APIKey:  cfg.DeepSeek.APIKey,
			Model:   cfg.DeepSeek.Model,
			BaseURL: cfg.DeepSeek.BaseURL,
		})
		if err != nil {
			logger.Warn("Failed to initialize DeepSeek provider", zap.Error(err))
		} else {
			r.providers["deepseek"] = deepseekProvider
		}
	}

	// Initialize Ollama provider
	ollamaProvider, err := NewOllamaProvider(ProviderConfig{
		BaseURL: cfg.Ollama.Host,
		Model:   cfg.Ollama.Model,
	})
	if err == nil && ollamaProvider.IsAvailable() {
		r.providers["ollama"] = ollamaProvider
	}

	return r, nil
}

// defaultTaskMapping returns the default provider mapping for each task
func defaultTaskMapping() map[TaskType]string {
	return map[TaskType]string{
		TaskTrendAnalysis:    "anthropic", // Claude is good for analysis
		TaskPaperParsing:     "openai",    // GPT-4o supports multimodal
		TaskPaperSummary:     "deepseek",  // Cost-effective for summaries
		TaskLearningPath:     "deepseek",  // Cost-effective
		TaskReviewWriting:    "anthropic", // Claude excels at long-form writing
		TaskReviewScoring:    "openai",    // GPT-4 for evaluation
		TaskCategoryClassify: "deepseek",  // Simple classification
		TaskEmbedding:        "openai",    // OpenAI embeddings
	}
}

// GetProvider returns a specific provider
func (r *Router) GetProvider(name string) (Provider, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	provider, ok := r.providers[name]
	return provider, ok
}

// SelectProvider selects the best provider for a task
func (r *Router) SelectProvider(task TaskType, pref *ProviderPreference) Provider {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Check user preference first
	if pref != nil && pref.PreferredProvider != "" {
		if provider, ok := r.providers[pref.PreferredProvider]; ok && provider.IsAvailable() {
			// Check if user has custom API key
			if apiKey, hasKey := pref.APIKeys[pref.PreferredProvider]; hasKey && apiKey != "" {
				// Create a new provider instance with user's API key
				return r.createProviderWithKey(pref.PreferredProvider, apiKey)
			}
			return provider
		}
	}

	// Use default mapping
	if defaultProvider, ok := r.taskMapping[task]; ok {
		if provider, ok := r.providers[defaultProvider]; ok && provider.IsAvailable() {
			return provider
		}
	}

	// Fallback to any available provider
	for _, provider := range r.providers {
		if provider.IsAvailable() {
			return provider
		}
	}

	return nil
}

// createProviderWithKey creates a provider with a custom API key
func (r *Router) createProviderWithKey(name, apiKey string) Provider {
	switch name {
	case "openai":
		provider, _ := NewOpenAIProvider(ProviderConfig{
			APIKey: apiKey,
			Model:  r.config.OpenAI.Model,
		})
		return provider
	case "anthropic":
		provider, _ := NewClaudeProvider(ProviderConfig{
			APIKey: apiKey,
			Model:  r.config.Anthropic.Model,
		})
		return provider
	case "deepseek":
		provider, _ := NewDeepSeekProvider(ProviderConfig{
			APIKey:  apiKey,
			Model:   r.config.DeepSeek.Model,
			BaseURL: r.config.DeepSeek.BaseURL,
		})
		return provider
	}
	return nil
}

// Chat sends a chat request using the appropriate provider
func (r *Router) Chat(ctx context.Context, task TaskType, req *ChatRequest, pref *ProviderPreference) (*ChatResponse, error) {
	provider := r.SelectProvider(task, pref)
	if provider == nil {
		return nil, ErrNoAvailableProvider
	}

	resp, err := provider.Chat(ctx, req)
	if err != nil {
		// Try fallback if enabled
		if pref != nil && pref.FallbackEnabled {
			for _, p := range r.providers {
				if p.IsAvailable() && p.Name() != provider.Name() {
					if fallbackResp, fallbackErr := p.Chat(ctx, req); fallbackErr == nil {
						return fallbackResp, nil
					}
				}
			}
		}
		return nil, err
	}

	return resp, nil
}

// Embed generates embeddings using the appropriate provider
func (r *Router) Embed(ctx context.Context, text string, pref *ProviderPreference) ([]float32, error) {
	provider := r.SelectProvider(TaskEmbedding, pref)
	if provider == nil {
		return nil, ErrNoAvailableProvider
	}

	return provider.Embed(ctx, text)
}

// GetAvailableProviders returns list of available providers
func (r *Router) GetAvailableProviders() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var available []string
	for name, provider := range r.providers {
		if provider.IsAvailable() {
			available = append(available, name)
		}
	}
	return available
}

// HasEmbeddingProvider returns true if an embedding provider is available
func (r *Router) HasEmbeddingProvider() bool {
	for _, provider := range r.providers {
		if provider.IsAvailable() && provider.GetEmbeddingDimension() > 0 {
			return true
		}
	}
	return false
}

