package postgres

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/virtualguard/PaperBeginner/internal/config"
	"github.com/virtualguard/PaperBeginner/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewConnection creates a new database connection
func NewConnection(cfg config.DatabaseConfig) (*gorm.DB, error) {
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}

	db, err := gorm.Open(postgres.Open(cfg.URL), gormConfig)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	return db, nil
}

// NewRedisClient creates a new Redis client
func NewRedisClient(cfg config.RedisConfig) (*redis.Client, error) {
	opt, err := redis.ParseURL(cfg.URL)
	if err != nil {
		// Fall back to simple address parsing
		opt = &redis.Options{
			Addr:     cfg.URL,
			Password: cfg.Password,
			DB:       cfg.DB,
		}
	}

	client := redis.NewClient(opt)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return client, nil
}

// AutoMigrate runs database migrations
func AutoMigrate(db *gorm.DB) error {
	// Enable UUID extension
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	
	// Enable pgvector extension
	db.Exec("CREATE EXTENSION IF NOT EXISTS vector")

	return db.AutoMigrate(
		&domain.User{},
		&domain.Paper{},
		&domain.PaperAnalysis{},
		&domain.CCFCategory{},
		&domain.TrendingItem{},
		&domain.TrendReport{},
		&domain.LearningPath{},
		&domain.Review{},
	)
}

