package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/virtualguard/PaperBeginer/internal/config"
	"github.com/virtualguard/PaperBeginer/internal/domain"
	"github.com/virtualguard/PaperBeginer/internal/worker"
	"github.com/virtualguard/PaperBeginer/pkg/logger"
	"go.uber.org/zap"
)

var (
	ErrPaperNotFound = errors.New("paper not found")
	ErrUnauthorized  = errors.New("unauthorized access")
	ErrInvalidFile   = errors.New("invalid file")
)

// PaperService handles paper-related business logic
type PaperService struct {
	paperRepo   domain.PaperRepository
	minioClient *minio.Client
	asynqClient *asynq.Client
	config      *config.Config
}

// NewPaperService creates a new PaperService
func NewPaperService(paperRepo domain.PaperRepository, cfg *config.Config) *PaperService {
	// Initialize MinIO client
	minioClient, err := minio.New(cfg.MinIO.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIO.AccessKey, cfg.MinIO.SecretKey, ""),
		Secure: cfg.MinIO.UseSSL,
	})
	if err != nil {
		logger.Error("Failed to initialize MinIO client", zap.Error(err))
	}

	// Initialize Asynq client
	asynqClient := asynq.NewClient(asynq.RedisClientOpt{
		Addr: cfg.Redis.URL,
	})

	return &PaperService{
		paperRepo:   paperRepo,
		minioClient: minioClient,
		asynqClient: asynqClient,
		config:      cfg,
	}
}

// Upload uploads a paper and stores it
func (s *PaperService) Upload(
	userID uuid.UUID,
	file io.Reader,
	header *multipart.FileHeader,
	title, authors, abstract string,
) (*domain.Paper, error) {
	// Validate file type
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".pdf" {
		return nil, ErrInvalidFile
	}

	// Generate unique file path
	fileID := uuid.New().String()
	filePath := fmt.Sprintf("papers/%s/%s%s", userID.String(), fileID, ext)

	// Ensure bucket exists
	ctx := context.Background()
	exists, err := s.minioClient.BucketExists(ctx, s.config.MinIO.Bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket: %w", err)
	}
	if !exists {
		if err := s.minioClient.MakeBucket(ctx, s.config.MinIO.Bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	// Upload to MinIO
	_, err = s.minioClient.PutObject(ctx, s.config.MinIO.Bucket, filePath, file, header.Size, minio.PutObjectOptions{
		ContentType: "application/pdf",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	// Parse authors
	var authorList []string
	if authors != "" {
		authorList = strings.Split(authors, ",")
		for i := range authorList {
			authorList[i] = strings.TrimSpace(authorList[i])
		}
	}

	// Create paper record
	paper := &domain.Paper{
		ID:       uuid.New(),
		UserID:   userID,
		Title:    title,
		Authors:  authorList,
		Abstract: abstract,
		FilePath: filePath,
		FileSize: header.Size,
		Status:   domain.PaperStatusPending,
	}

	if err := s.paperRepo.Create(paper); err != nil {
		// Try to clean up uploaded file
		_ = s.minioClient.RemoveObject(ctx, s.config.MinIO.Bucket, filePath, minio.RemoveObjectOptions{})
		return nil, fmt.Errorf("failed to create paper record: %w", err)
	}

	// Queue paper processing task
	go func() {
		payload, _ := json.Marshal(worker.PaperAnalysisPayload{
			PaperID: paper.ID.String(),
			UserID:  userID.String(),
			Types:   []string{"summary"},
		})
		task := asynq.NewTask(worker.TypePaperAnalysis, payload)
		_, err := s.asynqClient.Enqueue(task, asynq.Queue("default"))
		if err != nil {
			logger.Error("Failed to enqueue paper analysis task", zap.Error(err))
		}
	}()

	return paper, nil
}

// GetByID retrieves a paper by ID
func (s *PaperService) GetByID(paperID, userID uuid.UUID) (*domain.Paper, error) {
	paper, err := s.paperRepo.GetByID(paperID)
	if err != nil {
		return nil, ErrPaperNotFound
	}

	// Check ownership
	if paper.UserID != userID {
		return nil, ErrUnauthorized
	}

	return paper, nil
}

// ListByUserID lists papers for a user
func (s *PaperService) ListByUserID(userID uuid.UUID, offset, limit int) ([]*domain.Paper, int64, error) {
	return s.paperRepo.GetByUserID(userID, offset, limit)
}

// Delete deletes a paper
func (s *PaperService) Delete(paperID, userID uuid.UUID) error {
	paper, err := s.paperRepo.GetByID(paperID)
	if err != nil {
		return ErrPaperNotFound
	}

	if paper.UserID != userID {
		return ErrUnauthorized
	}

	// Delete from MinIO
	ctx := context.Background()
	if paper.FilePath != "" {
		_ = s.minioClient.RemoveObject(ctx, s.config.MinIO.Bucket, paper.FilePath, minio.RemoveObjectOptions{})
	}

	return s.paperRepo.Delete(paperID)
}

// QueueAnalysis queues paper analysis tasks
func (s *PaperService) QueueAnalysis(paperID, userID uuid.UUID, types []string) error {
	// Verify paper exists and user owns it
	paper, err := s.paperRepo.GetByID(paperID)
	if err != nil {
		return ErrPaperNotFound
	}

	if paper.UserID != userID {
		return ErrUnauthorized
	}

	// Update status
	paper.Status = domain.PaperStatusProcessing
	_ = s.paperRepo.Update(paper)

	// Queue analysis task
	payload, _ := json.Marshal(worker.PaperAnalysisPayload{
		PaperID: paperID.String(),
		UserID:  userID.String(),
		Types:   types,
	})

	task := asynq.NewTask(worker.TypePaperAnalysis, payload)
	_, err = s.asynqClient.Enqueue(task, asynq.Queue("default"))
	return err
}

// GetAnalyses retrieves all analyses for a paper
func (s *PaperService) GetAnalyses(paperID, userID uuid.UUID) ([]domain.PaperAnalysis, error) {
	paper, err := s.paperRepo.GetByID(paperID)
	if err != nil {
		return nil, ErrPaperNotFound
	}

	if paper.UserID != userID {
		return nil, ErrUnauthorized
	}

	return s.paperRepo.GetAnalysesByPaperID(paperID)
}

// GetCCFCategories retrieves all CCF categories
func (s *PaperService) GetCCFCategories() ([]domain.CCFCategory, error) {
	return s.paperRepo.GetCCFCategories()
}

// Close closes the service connections
func (s *PaperService) Close() error {
	if s.asynqClient != nil {
		return s.asynqClient.Close()
	}
	return nil
}

