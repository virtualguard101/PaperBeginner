package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
)

// Paper represents an academic paper
type Paper struct {
	ID           uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID       uuid.UUID       `json:"user_id" gorm:"type:uuid;index;not null"`
	Title        string          `json:"title" gorm:"not null"`
	Authors      []string        `json:"authors" gorm:"type:text[]"`
	Abstract     string          `json:"abstract"`
	Keywords     []string        `json:"keywords" gorm:"type:text[]"`
	PublishedAt  *time.Time      `json:"published_at,omitempty"`
	Venue        string          `json:"venue,omitempty"`
	DOI          string          `json:"doi,omitempty" gorm:"index"`
	SourceURL    string          `json:"source_url,omitempty"`
	FilePath     string          `json:"file_path,omitempty"`
	FileSize     int64           `json:"file_size,omitempty"`
	CCFCategory  *CCFCategory    `json:"ccf_category,omitempty" gorm:"foreignKey:CCFCategoryID"`
	CCFCategoryID *int           `json:"ccf_category_id,omitempty"`
	CCFRank      string          `json:"ccf_rank,omitempty"`
	Status       PaperStatus     `json:"status" gorm:"default:pending"`
	Embedding    pgvector.Vector `json:"-" gorm:"type:vector(1536)"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	
	// Relations
	User      *User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Analyses  []PaperAnalysis `json:"analyses,omitempty" gorm:"foreignKey:PaperID"`
}

// PaperStatus represents the processing status of a paper
type PaperStatus string

const (
	PaperStatusPending    PaperStatus = "pending"
	PaperStatusProcessing PaperStatus = "processing"
	PaperStatusCompleted  PaperStatus = "completed"
	PaperStatusFailed     PaperStatus = "failed"
)

// PaperAnalysis represents an AI analysis of a paper
type PaperAnalysis struct {
	ID            uuid.UUID    `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	PaperID       uuid.UUID    `json:"paper_id" gorm:"type:uuid;index;not null"`
	AnalysisType  AnalysisType `json:"analysis_type" gorm:"not null"`
	Content       string       `json:"content" gorm:"type:text"`
	StructuredData interface{} `json:"structured_data,omitempty" gorm:"type:jsonb"`
	LLMProvider   string       `json:"llm_provider"`
	LLMModel      string       `json:"llm_model"`
	TokensUsed    int          `json:"tokens_used"`
	CreatedAt     time.Time    `json:"created_at"`
}

// AnalysisType defines the type of paper analysis
type AnalysisType string

const (
	AnalysisSummary       AnalysisType = "summary"
	AnalysisMethodology   AnalysisType = "methodology"
	AnalysisContributions AnalysisType = "contributions"
	AnalysisCitations     AnalysisType = "citations"
	AnalysisComparison    AnalysisType = "comparison"
	AnalysisStrengths     AnalysisType = "strengths"
	AnalysisWeaknesses    AnalysisType = "weaknesses"
)

// CCFCategory represents a CCF academic category
type CCFCategory struct {
	ID     int    `json:"id" gorm:"primary_key"`
	Name   string `json:"name" gorm:"not null"`
	NameEN string `json:"name_en" gorm:"not null"`
}

// PaperRepository defines the interface for paper data access
type PaperRepository interface {
	Create(paper *Paper) error
	GetByID(id uuid.UUID) (*Paper, error)
	GetByUserID(userID uuid.UUID, offset, limit int) ([]*Paper, int64, error)
	Update(paper *Paper) error
	Delete(id uuid.UUID) error
	Search(query string, categoryID *int, limit int) ([]*Paper, error)
	FindSimilar(embedding pgvector.Vector, limit int) ([]*Paper, error)
	CreateAnalysis(analysis *PaperAnalysis) error
	GetAnalysesByPaperID(paperID uuid.UUID) ([]PaperAnalysis, error)
	GetCCFCategories() ([]CCFCategory, error)
}

// UploadPaperRequest represents a paper upload request
type UploadPaperRequest struct {
	Title    string `json:"title" binding:"required"`
	Authors  string `json:"authors,omitempty"`
	Abstract string `json:"abstract,omitempty"`
	DOI      string `json:"doi,omitempty"`
}

// PaperResponse represents a paper in API responses
type PaperResponse struct {
	ID          uuid.UUID       `json:"id"`
	Title       string          `json:"title"`
	Authors     []string        `json:"authors"`
	Abstract    string          `json:"abstract"`
	Keywords    []string        `json:"keywords"`
	PublishedAt *time.Time      `json:"published_at,omitempty"`
	Venue       string          `json:"venue,omitempty"`
	DOI         string          `json:"doi,omitempty"`
	CCFCategory *CCFCategory    `json:"ccf_category,omitempty"`
	CCFRank     string          `json:"ccf_rank,omitempty"`
	Status      PaperStatus     `json:"status"`
	Analyses    []PaperAnalysis `json:"analyses,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
}

// ToResponse converts a Paper to PaperResponse
func (p *Paper) ToResponse() *PaperResponse {
	return &PaperResponse{
		ID:          p.ID,
		Title:       p.Title,
		Authors:     p.Authors,
		Abstract:    p.Abstract,
		Keywords:    p.Keywords,
		PublishedAt: p.PublishedAt,
		Venue:       p.Venue,
		DOI:         p.DOI,
		CCFCategory: p.CCFCategory,
		CCFRank:     p.CCFRank,
		Status:      p.Status,
		Analyses:    p.Analyses,
		CreatedAt:   p.CreatedAt,
	}
}

// TableName returns the table name for GORM
func (Paper) TableName() string {
	return "papers"
}

func (PaperAnalysis) TableName() string {
	return "paper_analyses"
}

func (CCFCategory) TableName() string {
	return "ccf_categories"
}

