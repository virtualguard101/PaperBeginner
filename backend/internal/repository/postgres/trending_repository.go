package postgres

import (
	"github.com/virtualguard/PaperBeginer/internal/domain"
	"gorm.io/gorm"
)

// TrendingRepository implements domain.TrendingRepository
type TrendingRepository struct {
	db *gorm.DB
}

// NewTrendingRepository creates a new TrendingRepository
func NewTrendingRepository(db *gorm.DB) *TrendingRepository {
	return &TrendingRepository{db: db}
}

// CreateTrendingItem creates a new trending item
func (r *TrendingRepository) CreateTrendingItem(item *domain.TrendingItem) error {
	return r.db.Create(item).Error
}

// GetTrendingItems retrieves trending items with filters
func (r *TrendingRepository) GetTrendingItems(source *domain.TrendSource, categoryID *int, limit int) ([]*domain.TrendingItem, error) {
	var items []*domain.TrendingItem

	db := r.db.Preload("CCFCategory")

	if source != nil {
		db = db.Where("source = ?", *source)
	}

	if categoryID != nil {
		db = db.Where("ccf_category_id = ?", *categoryID)
	}

	err := db.Order("trend_score DESC, crawled_at DESC").Limit(limit).Find(&items).Error
	return items, err
}

// CreateTrendReport creates a new trend report
func (r *TrendingRepository) CreateTrendReport(report *domain.TrendReport) error {
	return r.db.Create(report).Error
}

// GetTrendReports retrieves trend reports with filters
func (r *TrendingRepository) GetTrendReports(categoryID *int, limit int) ([]*domain.TrendReport, error) {
	var reports []*domain.TrendReport

	db := r.db.Preload("Category")

	if categoryID != nil {
		db = db.Where("category_id = ?", *categoryID)
	}

	err := db.Order("created_at DESC").Limit(limit).Find(&reports).Error
	return reports, err
}

// UpsertTrendingItem upserts a trending item based on source and source_id
func (r *TrendingRepository) UpsertTrendingItem(item *domain.TrendingItem) error {
	return r.db.Where("source = ? AND source_id = ?", item.Source, item.SourceID).
		Assign(item).
		FirstOrCreate(item).Error
}

