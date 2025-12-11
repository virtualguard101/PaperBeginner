package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/virtualguard101/PaperBeginner/internal/config"
	"github.com/virtualguard101/PaperBeginner/internal/handler"
	"github.com/virtualguard101/PaperBeginner/internal/handler/middleware"
	"github.com/virtualguard101/PaperBeginner/internal/repository/postgres"
	"github.com/virtualguard101/PaperBeginner/internal/service"
	"github.com/virtualguard101/PaperBeginner/pkg/logger"
	"github.com/virtualguard101/PaperBeginner/pkg/validator"
	"go.uber.org/zap"
)

// @title PaperBeginner API
// @version 1.0
// @description AI-powered Academic Research Guide API
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
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

	logger.Info("Starting PaperBeginner API",
		zap.String("env", cfg.App.Env),
		zap.Int("port", cfg.App.Port),
	)

	// Initialize validator
	if err := validator.Init(); err != nil {
		logger.Fatal("Failed to initialize validator", zap.Error(err))
	}

	// Initialize database connection
	db, err := postgres.NewConnection(cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Run migrations
	if err := postgres.AutoMigrate(db); err != nil {
		logger.Fatal("Failed to run migrations", zap.Error(err))
	}

	// Initialize Redis client
	redisClient, err := postgres.NewRedisClient(cfg.Redis)
	if err != nil {
		logger.Fatal("Failed to connect to Redis", zap.Error(err))
	}
	defer redisClient.Close()

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	paperRepo := postgres.NewPaperRepository(db)
	trendingRepo := postgres.NewTrendingRepository(db)
	learningPathRepo := postgres.NewLearningPathRepository(db)
	reviewRepo := postgres.NewReviewRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo, cfg.JWT)
	paperService := service.NewPaperService(paperRepo, cfg)
	// agentService := service.NewAgentService(cfg.LLM)
	// crawlerService := service.NewCrawlerService(cfg.GitHub, trendingRepo)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)
	paperHandler := handler.NewPaperHandler(paperService)

	// Set Gin mode
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	router := gin.New()

	// Global middleware
	router.Use(gin.Recovery())
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())
	router.Use(middleware.RequestID())

	// Rate limiting (if configured)
	if cfg.App.RateLimitRPS > 0 {
		router.Use(middleware.RateLimit(cfg.App.RateLimitRPS, cfg.App.RateLimitBurst))
	}

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"version": "1.0.0",
		})
	})

	// API routes
	api := router.Group("/api/v1")
	{
		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
			auth.POST("/refresh", userHandler.RefreshToken)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.Auth(cfg.JWT.Secret))
		{
			// User routes
			users := protected.Group("/users")
			{
				users.GET("/me", userHandler.GetCurrentUser)
				users.PUT("/me", userHandler.UpdateCurrentUser)
				users.PUT("/me/password", userHandler.ChangePassword)
			}

			// Paper routes
			papers := protected.Group("/papers")
			{
				papers.POST("/upload", paperHandler.Upload)
				papers.GET("", paperHandler.List)
				papers.GET("/:id", paperHandler.Get)
				papers.DELETE("/:id", paperHandler.Delete)
				papers.POST("/:id/analyze", paperHandler.Analyze)
				papers.GET("/:id/analyses", paperHandler.GetAnalyses)
			}

			// Trending routes
			trending := protected.Group("/trending")
			{
				trending.GET("", handler.GetTrending)
				trending.GET("/reports", handler.GetTrendReports)
			}

			// Learning path routes
			learning := protected.Group("/learning")
			{
				learning.POST("/generate", handler.GenerateLearningPath)
				learning.GET("/paths", handler.GetLearningPaths)
				learning.GET("/paths/:id", handler.GetLearningPath)
			}

			// Review routes
			reviews := protected.Group("/reviews")
			{
				reviews.POST("/generate", handler.GenerateReview)
				reviews.GET("", handler.ListReviews)
				reviews.GET("/:id", handler.GetReview)
				reviews.POST("/:id/score", handler.ScoreReview)
				reviews.DELETE("/:id", handler.DeleteReview)
			}

			// CCF categories
			api.GET("/ccf/categories", paperHandler.GetCCFCategories)
		}
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.App.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Server starting", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited gracefully")

	// Suppress unused variable warnings (remove when implemented)
	_ = trendingRepo
	_ = learningPathRepo
	_ = reviewRepo
}
