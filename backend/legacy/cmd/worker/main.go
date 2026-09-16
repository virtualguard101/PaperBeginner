package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/hibiken/asynq"
	"github.com/virtualguard/PaperBeginner/internal/config"
	"github.com/virtualguard/PaperBeginner/internal/repository/postgres"
	"github.com/virtualguard/PaperBeginner/internal/worker"
	"github.com/virtualguard/PaperBeginner/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	if err := logger.Init(cfg.App.Env); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("Starting PaperBeginner Worker",
		zap.String("env", cfg.App.Env),
	)

	// Initialize database connection
	db, err := postgres.NewConnection(cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Create Asynq server
	redisOpt := asynq.RedisClientOpt{
		Addr: cfg.Redis.URL,
	}

	srv := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				logger.Error("Task failed",
					zap.String("type", task.Type()),
					zap.Error(err),
				)
			}),
		},
	)

	// Create worker handler
	workerHandler := worker.NewHandler(db, cfg)

	// Register task handlers
	mux := asynq.NewServeMux()
	mux.HandleFunc(worker.TypePaperAnalysis, workerHandler.HandlePaperAnalysis)
	mux.HandleFunc(worker.TypeGitHubCrawl, workerHandler.HandleGitHubCrawl)
	mux.HandleFunc(worker.TypeCCFCrawl, workerHandler.HandleCCFCrawl)
	mux.HandleFunc(worker.TypeTrendReport, workerHandler.HandleTrendReport)
	mux.HandleFunc(worker.TypeLearningPath, workerHandler.HandleLearningPath)
	mux.HandleFunc(worker.TypeReviewGeneration, workerHandler.HandleReviewGeneration)
	mux.HandleFunc(worker.TypeReviewScoring, workerHandler.HandleReviewScoring)

	// Handle shutdown gracefully
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		logger.Info("Shutting down worker...")
		cancel()
		srv.Shutdown()
	}()

	// Start the worker server
	logger.Info("Worker starting...")
	if err := srv.Run(mux); err != nil {
		logger.Fatal("Failed to run worker", zap.Error(err))
	}

	<-ctx.Done()
	logger.Info("Worker exited gracefully")
}

