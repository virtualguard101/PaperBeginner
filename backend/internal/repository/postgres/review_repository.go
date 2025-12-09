package postgres

import (
	"github.com/google/uuid"
	"github.com/virtualguard/PaperBeginer/internal/domain"
	"gorm.io/gorm"
)

// ReviewRepository implements domain.ReviewRepository
type ReviewRepository struct {
	db *gorm.DB
}

// NewReviewRepository creates a new ReviewRepository
func NewReviewRepository(db *gorm.DB) *ReviewRepository {
	return &ReviewRepository{db: db}
}

// Create creates a new review
func (r *ReviewRepository) Create(review *domain.Review) error {
	return r.db.Create(review).Error
}

// GetByID retrieves a review by ID
func (r *ReviewRepository) GetByID(id uuid.UUID) (*domain.Review, error) {
	var review domain.Review
	err := r.db.Preload("Category").Preload("Papers").Where("id = ?", id).First(&review).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

// GetByUserID retrieves reviews by user ID with pagination
func (r *ReviewRepository) GetByUserID(userID uuid.UUID, offset, limit int) ([]*domain.Review, int64, error) {
	var reviews []*domain.Review
	var total int64

	r.db.Model(&domain.Review{}).Where("user_id = ?", userID).Count(&total)

	err := r.db.Preload("Category").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&reviews).Error

	if err != nil {
		return nil, 0, err
	}

	return reviews, total, nil
}

// Update updates a review
func (r *ReviewRepository) Update(review *domain.Review) error {
	return r.db.Save(review).Error
}

// Delete deletes a review
func (r *ReviewRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&domain.Review{}, "id = ?", id).Error
}

