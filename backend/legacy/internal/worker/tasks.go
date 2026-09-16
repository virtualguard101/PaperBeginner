package worker

// Task type constants
const (
	TypePaperAnalysis    = "paper:analysis"
	TypeGitHubCrawl      = "crawler:github"
	TypeCCFCrawl         = "crawler:ccf"
	TypeTrendReport      = "trend:report"
	TypeLearningPath     = "learning:path"
	TypeReviewGeneration = "review:generate"
	TypeReviewScoring    = "review:score"
)

// PaperAnalysisPayload represents the payload for paper analysis tasks
type PaperAnalysisPayload struct {
	PaperID string   `json:"paper_id"`
	UserID  string   `json:"user_id"`
	Types   []string `json:"types"` // summary, methodology, contributions, etc.
}

// GitHubCrawlPayload represents the payload for GitHub crawl tasks
type GitHubCrawlPayload struct {
	Language string `json:"language,omitempty"`
	Since    string `json:"since,omitempty"` // daily, weekly, monthly
}

// CCFCrawlPayload represents the payload for CCF crawl tasks
type CCFCrawlPayload struct {
	CategoryID int    `json:"category_id,omitempty"`
	Rank       string `json:"rank,omitempty"` // A, B, C
}

// TrendReportPayload represents the payload for trend report generation
type TrendReportPayload struct {
	CategoryID int    `json:"category_id"`
	Period     string `json:"period"` // e.g., "2024-W01"
}

// LearningPathPayload represents the payload for learning path generation
type LearningPathPayload struct {
	UserID        string   `json:"user_id"`
	CategoryID    int      `json:"category_id"`
	Difficulty    string   `json:"difficulty"`
	Prerequisites []string `json:"prerequisites,omitempty"`
	FocusAreas    []string `json:"focus_areas,omitempty"`
}

// ReviewGenerationPayload represents the payload for review generation
type ReviewGenerationPayload struct {
	UserID     string   `json:"user_id"`
	ReviewID   string   `json:"review_id"`
	PaperIDs   []string `json:"paper_ids"`
	Title      string   `json:"title"`
	CategoryID *int     `json:"category_id,omitempty"`
	Style      string   `json:"style,omitempty"`
}

// ReviewScoringPayload represents the payload for review scoring
type ReviewScoringPayload struct {
	UserID   string `json:"user_id"`
	ReviewID string `json:"review_id"`
	Content  string `json:"content,omitempty"`
}

