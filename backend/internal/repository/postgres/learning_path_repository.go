package postgres

import (
	"github.com/google/uuid"
	"github.com/virtualguard/PaperBeginer/internal/domain"
	"gorm.io/gorm"
)

// LearningPathRepository implements domain.LearningPathRepository
type LearningPathRepository struct {
	db *gorm.DB
}

// NewLearningPathRepository creates a new LearningPathRepository
func NewLearningPathRepository(db *gorm.DB) *LearningPathRepository {
	return &LearningPathRepository{db: db}
}

// Create creates a new learning path
func (r *LearningPathRepository) Create(path *domain.LearningPath) error {
	return r.db.Create(path).Error
}

// GetByID retrieves a learning path by ID
func (r *LearningPathRepository) GetByID(id uuid.UUID) (*domain.LearningPath, error) {
	var path domain.LearningPath
	err := r.db.Preload("Category").Where("id = ?", id).First(&path).Error
	if err != nil {
		return nil, err
	}
	return &path, nil
}

// GetByUserID retrieves learning paths by user ID
func (r *LearningPathRepository) GetByUserID(userID uuid.UUID) ([]*domain.LearningPath, error) {
	var paths []*domain.LearningPath
	err := r.db.Preload("Category").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&paths).Error
	return paths, err
}

// GetByCategoryID retrieves learning paths by category ID
func (r *LearningPathRepository) GetByCategoryID(categoryID int) ([]*domain.LearningPath, error) {
	var paths []*domain.LearningPath
	err := r.db.Preload("Category").
		Where("category_id = ?", categoryID).
		Order("created_at DESC").
		Find(&paths).Error
	return paths, err
}

// Update updates a learning path
func (r *LearningPathRepository) Update(path *domain.LearningPath) error {
	return r.db.Save(path).Error
}

// Delete deletes a learning path
func (r *LearningPathRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&domain.LearningPath{}, "id = ?", id).Error
}

