package domain

import (
	"time"

	"github.com/google/uuid"
)

// TrendingItem represents a trending topic or project
type TrendingItem struct {
	ID            uuid.UUID    `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Source        TrendSource  `json:"source" gorm:"not null;index"`
	SourceID      string       `json:"source_id" gorm:"index"`
	Title         string       `json:"title" gorm:"not null"`
	Description   string       `json:"description" gorm:"type:text"`
	URL           string       `json:"url"`
	Stars         int          `json:"stars,omitempty"`
	Forks         int          `json:"forks,omitempty"`
	Language      string       `json:"language,omitempty"`
	Topics        []string     `json:"topics" gorm:"type:text[]"`
	CCFCategoryID *int         `json:"ccf_category_id,omitempty"`
	CCFCategory   *CCFCategory `json:"ccf_category,omitempty" gorm:"foreignKey:CCFCategoryID"`
	TrendScore    float64      `json:"trend_score"`
	CrawledAt     time.Time    `json:"crawled_at"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

// TrendSource defines the source of a trending item
type TrendSource string

const (
	TrendSourceGitHub TrendSource = "github"
	TrendSourceCCF    TrendSource = "ccf"
	TrendSourceArxiv  TrendSource = "arxiv"
)

// TrendReport represents an aggregated trend report
type TrendReport struct {
	ID           uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Period       string          `json:"period" gorm:"not null"` // e.g., "2024-W01", "2024-01"
	CategoryID   int             `json:"category_id" gorm:"index"`
	Category     *CCFCategory    `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	Summary      string          `json:"summary" gorm:"type:text"`
	Highlights   []TrendHighlight `json:"highlights" gorm:"type:jsonb"`
	GeneratedBy  string          `json:"generated_by"` // LLM model used
	CreatedAt    time.Time       `json:"created_at"`
}

// TrendHighlight represents a key highlight in a trend report
type TrendHighlight struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Sources     []string `json:"sources"`
	Importance  int      `json:"importance"` // 1-5
}

// LearningPath represents a recommended learning path
type LearningPath struct {
	ID          uuid.UUID          `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID      uuid.UUID          `json:"user_id" gorm:"type:uuid;index"`
	CategoryID  int                `json:"category_id" gorm:"index"`
	Category    *CCFCategory       `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	Title       string             `json:"title" gorm:"not null"`
	Description string             `json:"description" gorm:"type:text"`
	Stages      []LearningStage    `json:"stages" gorm:"type:jsonb"`
	Prerequisites []string         `json:"prerequisites" gorm:"type:text[]"`
	EstimatedTime string           `json:"estimated_time"`
	Difficulty  string             `json:"difficulty"` // beginner, intermediate, advanced
	GeneratedBy string             `json:"generated_by"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

// LearningStage represents a stage in a learning path
type LearningStage struct {
	Order       int                `json:"order"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Resources   []LearningResource `json:"resources"`
	Duration    string             `json:"duration"`
	Skills      []string           `json:"skills"`
}

// LearningResource represents a learning resource
type LearningResource struct {
	Type        string `json:"type"` // documentation, course, book, video, tutorial
	Title       string `json:"title"`
	URL         string `json:"url"`
	Provider    string `json:"provider"` // e.g., MIT, Stanford, official docs
	Language    string `json:"language"` // en, zh, etc.
	IsFree      bool   `json:"is_free"`
	Description string `json:"description,omitempty"`
}

// Review represents a literature review
type Review struct {
	ID           uuid.UUID     `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID       uuid.UUID     `json:"user_id" gorm:"type:uuid;index;not null"`
	Title        string        `json:"title" gorm:"not null"`
	Abstract     string        `json:"abstract" gorm:"type:text"`
	Content      string        `json:"content" gorm:"type:text"`
	PaperIDs     []uuid.UUID   `json:"paper_ids" gorm:"type:uuid[]"`
	CategoryID   *int          `json:"category_id,omitempty"`
	Category     *CCFCategory  `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	Status       ReviewStatus  `json:"status" gorm:"default:draft"`
	Score        *ReviewScore  `json:"score,omitempty" gorm:"type:jsonb"`
	GeneratedBy  string        `json:"generated_by"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	
	// Relations
	User   *User    `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Papers []*Paper `json:"papers,omitempty" gorm:"many2many:review_papers"`
}

// ReviewStatus defines the status of a review
type ReviewStatus string

const (
	ReviewStatusDraft     ReviewStatus = "draft"
	ReviewStatusGenerated ReviewStatus = "generated"
	ReviewStatusScored    ReviewStatus = "scored"
)

// ReviewScore represents the AI-generated score for a review
type ReviewScore struct {
	OverallScore    float64           `json:"overall_score"`    // 0-100
	Criteria        []ScoreCriterion  `json:"criteria"`
	Strengths       []string          `json:"strengths"`
	Weaknesses      []string          `json:"weaknesses"`
	Suggestions     []string          `json:"suggestions"`
	GeneratedAt     time.Time         `json:"generated_at"`
}

// ScoreCriterion represents a scoring criterion
type ScoreCriterion struct {
	Name        string  `json:"name"`
	Score       float64 `json:"score"`       // 0-100
	Weight      float64 `json:"weight"`      // 0-1
	Description string  `json:"description"`
	Feedback    string  `json:"feedback"`
}

// TrendingRepository defines the interface for trending data access
type TrendingRepository interface {
	CreateTrendingItem(item *TrendingItem) error
	GetTrendingItems(source *TrendSource, categoryID *int, limit int) ([]*TrendingItem, error)
	CreateTrendReport(report *TrendReport) error
	GetTrendReports(categoryID *int, limit int) ([]*TrendReport, error)
}

// LearningPathRepository defines the interface for learning path data access
type LearningPathRepository interface {
	Create(path *LearningPath) error
	GetByID(id uuid.UUID) (*LearningPath, error)
	GetByUserID(userID uuid.UUID) ([]*LearningPath, error)
	GetByCategoryID(categoryID int) ([]*LearningPath, error)
	Update(path *LearningPath) error
	Delete(id uuid.UUID) error
}

// ReviewRepository defines the interface for review data access
type ReviewRepository interface {
	Create(review *Review) error
	GetByID(id uuid.UUID) (*Review, error)
	GetByUserID(userID uuid.UUID, offset, limit int) ([]*Review, int64, error)
	Update(review *Review) error
	Delete(id uuid.UUID) error
}

// TableNames
func (TrendingItem) TableName() string {
	return "trending_items"
}

func (TrendReport) TableName() string {
	return "trend_reports"
}

func (LearningPath) TableName() string {
	return "learning_paths"
}

func (Review) TableName() string {
	return "reviews"
}

