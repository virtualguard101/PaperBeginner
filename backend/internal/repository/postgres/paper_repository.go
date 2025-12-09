package postgres

import (
	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
	"github.com/virtualguard/PaperBeginer/internal/domain"
	"gorm.io/gorm"
)

// PaperRepository implements domain.PaperRepository
type PaperRepository struct {
	db *gorm.DB
}

// NewPaperRepository creates a new PaperRepository
func NewPaperRepository(db *gorm.DB) *PaperRepository {
	return &PaperRepository{db: db}
}

// Create creates a new paper
func (r *PaperRepository) Create(paper *domain.Paper) error {
	return r.db.Create(paper).Error
}

// GetByID retrieves a paper by ID
func (r *PaperRepository) GetByID(id uuid.UUID) (*domain.Paper, error) {
	var paper domain.Paper
	err := r.db.Preload("CCFCategory").Preload("Analyses").Where("id = ?", id).First(&paper).Error
	if err != nil {
		return nil, err
	}
	return &paper, nil
}

// GetByUserID retrieves papers by user ID with pagination
func (r *PaperRepository) GetByUserID(userID uuid.UUID, offset, limit int) ([]*domain.Paper, int64, error) {
	var papers []*domain.Paper
	var total int64

	r.db.Model(&domain.Paper{}).Where("user_id = ?", userID).Count(&total)

	err := r.db.Preload("CCFCategory").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&papers).Error

	if err != nil {
		return nil, 0, err
	}

	return papers, total, nil
}

// Update updates a paper
func (r *PaperRepository) Update(paper *domain.Paper) error {
	return r.db.Save(paper).Error
}

// Delete deletes a paper
func (r *PaperRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&domain.Paper{}, "id = ?", id).Error
}

// Search searches papers by query
func (r *PaperRepository) Search(query string, categoryID *int, limit int) ([]*domain.Paper, error) {
	var papers []*domain.Paper

	db := r.db.Preload("CCFCategory")

	if query != "" {
		db = db.Where("title ILIKE ? OR abstract ILIKE ?", "%"+query+"%", "%"+query+"%")
	}

	if categoryID != nil {
		db = db.Where("ccf_category_id = ?", *categoryID)
	}

	err := db.Limit(limit).Find(&papers).Error
	return papers, err
}

// FindSimilar finds similar papers using vector similarity
func (r *PaperRepository) FindSimilar(embedding pgvector.Vector, limit int) ([]*domain.Paper, error) {
	var papers []*domain.Paper

	err := r.db.Preload("CCFCategory").
		Order(gorm.Expr("embedding <-> ?", embedding)).
		Limit(limit).
		Find(&papers).Error

	return papers, err
}

// CreateAnalysis creates a new paper analysis
func (r *PaperRepository) CreateAnalysis(analysis *domain.PaperAnalysis) error {
	return r.db.Create(analysis).Error
}

// GetAnalysesByPaperID retrieves all analyses for a paper
func (r *PaperRepository) GetAnalysesByPaperID(paperID uuid.UUID) ([]domain.PaperAnalysis, error) {
	var analyses []domain.PaperAnalysis
	err := r.db.Where("paper_id = ?", paperID).Order("created_at DESC").Find(&analyses).Error
	return analyses, err
}

// GetCCFCategories retrieves all CCF categories
func (r *PaperRepository) GetCCFCategories() ([]domain.CCFCategory, error) {
	var categories []domain.CCFCategory
	err := r.db.Order("id").Find(&categories).Error
	return categories, err
}

