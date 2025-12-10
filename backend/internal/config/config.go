package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	MinIO    MinIOConfig    `mapstructure:"minio"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	LLM      LLMConfig      `mapstructure:"llm"`
	GitHub   GitHubConfig   `mapstructure:"github"`
}

type AppConfig struct {
	Env            string `mapstructure:"env"`
	Port           int    `mapstructure:"port"`
	Secret         string `mapstructure:"secret"`
	RateLimitRPS   int    `mapstructure:"rate_limit_rps"`
	RateLimitBurst int    `mapstructure:"rate_limit_burst"`
}

type DatabaseConfig struct {
	URL             string        `mapstructure:"url"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

type RedisConfig struct {
	URL      string `mapstructure:"url"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type MinIOConfig struct {
	Endpoint  string `mapstructure:"endpoint"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	UseSSL    bool   `mapstructure:"use_ssl"`
	Bucket    string `mapstructure:"bucket"`
}

type JWTConfig struct {
	Secret string        `mapstructure:"secret"`
	Expiry time.Duration `mapstructure:"expiry"`
}

type LLMConfig struct {
	OpenAI    OpenAIConfig    `mapstructure:"openai"`
	Anthropic AnthropicConfig `mapstructure:"anthropic"`
	DeepSeek  DeepSeekConfig  `mapstructure:"deepseek"`
	Ollama    OllamaConfig    `mapstructure:"ollama"`
	Default   string          `mapstructure:"default"`
}

type OpenAIConfig struct {
	APIKey  string `mapstructure:"api_key"`
	Model   string `mapstructure:"model"`
	BaseURL string `mapstructure:"base_url"`
}

type AnthropicConfig struct {
	APIKey string `mapstructure:"api_key"`
	Model  string `mapstructure:"model"`
}

type DeepSeekConfig struct {
	APIKey  string `mapstructure:"api_key"`
	Model   string `mapstructure:"model"`
	BaseURL string `mapstructure:"base_url"`
}

type OllamaConfig struct {
	Host  string `mapstructure:"host"`
	Model string `mapstructure:"model"`
}

type GitHubConfig struct {
	Token string `mapstructure:"token"`
}

// Load reads configuration from environment variables and config files
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/etc/paperbeginner")

	// Set defaults
	setDefaults()

	// Read from environment variables
	viper.AutomaticEnv()

	// Map environment variables to config keys
	bindEnvVars()

	// Try to read config file (optional)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found is okay, we'll use env vars
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &cfg, nil
}

func setDefaults() {
	// App defaults
	viper.SetDefault("app.env", "development")
	viper.SetDefault("app.port", 8080)
	viper.SetDefault("app.rate_limit_rps", 10)
	viper.SetDefault("app.rate_limit_burst", 20)

	// Database defaults
	viper.SetDefault("database.max_open_conns", 25)
	viper.SetDefault("database.max_idle_conns", 5)
	viper.SetDefault("database.conn_max_lifetime", "5m")

	// Redis defaults
	viper.SetDefault("redis.db", 0)

	// MinIO defaults
	viper.SetDefault("minio.use_ssl", false)
	viper.SetDefault("minio.bucket", "paperbeginner")

	// JWT defaults
	viper.SetDefault("jwt.expiry", "24h")

	// LLM defaults
	viper.SetDefault("llm.default", "openai")
	viper.SetDefault("llm.openai.model", "gpt-4o")
	viper.SetDefault("llm.anthropic.model", "claude-3-5-sonnet-20241022")
	viper.SetDefault("llm.deepseek.model", "deepseek-chat")
	viper.SetDefault("llm.deepseek.base_url", "https://api.deepseek.com/v1")
	viper.SetDefault("llm.ollama.host", "http://localhost:11434")
	viper.SetDefault("llm.ollama.model", "llama3.1")
}

func bindEnvVars() {
	// App
	viper.BindEnv("app.env", "APP_ENV")
	viper.BindEnv("app.port", "APP_PORT")
	viper.BindEnv("app.secret", "APP_SECRET")
	viper.BindEnv("app.rate_limit_rps", "RATE_LIMIT_RPS")
	viper.BindEnv("app.rate_limit_burst", "RATE_LIMIT_BURST")

	// Database
	viper.BindEnv("database.url", "DATABASE_URL")

	// Redis
	viper.BindEnv("redis.url", "REDIS_URL")
	viper.BindEnv("redis.password", "REDIS_PASSWORD")

	// MinIO
	viper.BindEnv("minio.endpoint", "MINIO_ENDPOINT")
	viper.BindEnv("minio.access_key", "MINIO_ACCESS_KEY")
	viper.BindEnv("minio.secret_key", "MINIO_SECRET_KEY")
	viper.BindEnv("minio.use_ssl", "MINIO_USE_SSL")
	viper.BindEnv("minio.bucket", "MINIO_BUCKET")

	// JWT
	viper.BindEnv("jwt.secret", "JWT_SECRET")
	viper.BindEnv("jwt.expiry", "JWT_EXPIRY")

	// LLM
	viper.BindEnv("llm.openai.api_key", "OPENAI_API_KEY")
	viper.BindEnv("llm.anthropic.api_key", "ANTHROPIC_API_KEY")
	viper.BindEnv("llm.deepseek.api_key", "DEEPSEEK_API_KEY")
	viper.BindEnv("llm.ollama.host", "OLLAMA_HOST")

	// GitHub
	viper.BindEnv("github.token", "GITHUB_TOKEN")
}

// IsDevelopment returns true if the app is running in development mode
func (c *Config) IsDevelopment() bool {
	return c.App.Env == "development"
}

// IsProduction returns true if the app is running in production mode
func (c *Config) IsProduction() bool {
	return c.App.Env == "production"
}

