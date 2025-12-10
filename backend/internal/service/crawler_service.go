package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/virtualguard/PaperBeginner/internal/agent"
	"github.com/virtualguard/PaperBeginner/internal/config"
	"github.com/virtualguard/PaperBeginner/internal/domain"
	"github.com/virtualguard/PaperBeginner/internal/llm"
	"github.com/virtualguard/PaperBeginner/pkg/logger"
	"go.uber.org/zap"
)

// CrawlerService handles crawling GitHub and CCF data
type CrawlerService struct {
	githubConfig config.GitHubConfig
	trendingRepo domain.TrendingRepository
	trendAgent   *agent.TrendAgent
	httpClient   *http.Client
}

// NewCrawlerService creates a new CrawlerService
func NewCrawlerService(
	githubConfig config.GitHubConfig,
	trendingRepo domain.TrendingRepository,
	llmRouter *llm.Router,
) *CrawlerService {
	return &CrawlerService{
		githubConfig: githubConfig,
		trendingRepo: trendingRepo,
		trendAgent:   agent.NewTrendAgent(llmRouter),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// CrawlGitHubTrending crawls GitHub trending repositories
func (s *CrawlerService) CrawlGitHubTrending(ctx context.Context, language string, since string) error {
	// Build URL for GitHub trending
	baseURL := "https://api.github.com/search/repositories"
	
	// Calculate date range
	var dateFrom time.Time
	switch since {
	case "daily":
		dateFrom = time.Now().AddDate(0, 0, -1)
	case "weekly":
		dateFrom = time.Now().AddDate(0, 0, -7)
	case "monthly":
		dateFrom = time.Now().AddDate(0, -1, 0)
	default:
		dateFrom = time.Now().AddDate(0, 0, -7)
	}

	// Build search query
	query := fmt.Sprintf("created:>%s", dateFrom.Format("2006-01-02"))
	if language != "" {
		query += fmt.Sprintf(" language:%s", language)
	}

	params := url.Values{}
	params.Set("q", query)
	params.Set("sort", "stars")
	params.Set("order", "desc")
	params.Set("per_page", "50")

	reqURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	if s.githubConfig.Token != "" {
		req.Header.Set("Authorization", "Bearer "+s.githubConfig.Token)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("GitHub API error: %s - %s", resp.Status, string(body))
	}

	var result GitHubSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	// Process each repository
	for _, repo := range result.Items {
		// Create trending item
		item := &domain.TrendingItem{
			ID:          uuid.New(),
			Source:      domain.TrendSourceGitHub,
			SourceID:    fmt.Sprintf("%d", repo.ID),
			Title:       repo.FullName,
			Description: repo.Description,
			URL:         repo.HTMLURL,
			Stars:       repo.Stars,
			Forks:       repo.Forks,
			Language:    repo.Language,
			Topics:      repo.Topics,
			TrendScore:  s.calculateTrendScore(repo),
			CrawledAt:   time.Now(),
		}

		// Classify into CCF category using AI
		classifyResult, err := s.trendAgent.ClassifyCCFCategory(ctx, &agent.ClassifyInput{
			Title:       repo.FullName,
			Description: repo.Description,
			Topics:      repo.Topics,
			Language:    repo.Language,
		})
		if err == nil && classifyResult.Confidence > 0.5 {
			item.CCFCategoryID = &classifyResult.CategoryID
		}

		// Save to database
		if err := s.trendingRepo.CreateTrendingItem(item); err != nil {
			logger.Error("Failed to save trending item", zap.Error(err), zap.String("repo", repo.FullName))
		}
	}

	logger.Info("GitHub trending crawl completed",
		zap.Int("count", len(result.Items)),
		zap.String("language", language),
		zap.String("since", since),
	)

	return nil
}

// calculateTrendScore calculates a trend score for a repository
func (s *CrawlerService) calculateTrendScore(repo GitHubRepo) float64 {
	// Simple scoring algorithm based on stars, forks, and recency
	score := float64(repo.Stars) * 0.6
	score += float64(repo.Forks) * 0.3
	
	// Boost for recent activity
	if repo.PushedAt != "" {
		pushedAt, err := time.Parse(time.RFC3339, repo.PushedAt)
		if err == nil {
			daysSincePush := time.Since(pushedAt).Hours() / 24
			if daysSincePush < 7 {
				score *= 1.5
			} else if daysSincePush < 30 {
				score *= 1.2
			}
		}
	}

	return score
}

// CrawlCCFVenues crawls CCF conference/journal information
func (s *CrawlerService) CrawlCCFVenues(ctx context.Context, categoryID int, rank string) error {
	// Note: CCF website requires web scraping, which may need special handling
	// For now, we'll use a simplified approach with predefined venue data
	
	// This is a placeholder - in production, you would:
	// 1. Scrape CCF website: https://www.ccf.org.cn/Academic_Evaluation/By_category/
	// 2. Or use DBLP API for publication data
	// 3. Or integrate with Semantic Scholar API
	
	logger.Info("CCF venue crawl started",
		zap.Int("category_id", categoryID),
		zap.String("rank", rank),
	)
	
	// TODO: Implement CCF crawling
	// Options:
	// 1. Use colly or goquery for web scraping
	// 2. Integrate with DBLP API: https://dblp.org/faq/How+to+use+the+dblp+search+API.html
	// 3. Use Semantic Scholar API: https://api.semanticscholar.org/
	
	return nil
}

// FetchDPLBPublications fetches recent publications from DBLP
func (s *CrawlerService) FetchDPLBPublications(ctx context.Context, venue string, year int) ([]DPLBPublication, error) {
	// DBLP API endpoint
	apiURL := fmt.Sprintf("https://dblp.org/search/publ/api?q=venue:%s+year:%d&format=json&h=100", 
		url.QueryEscape(venue), year)

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result DPLBSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var publications []DPLBPublication
	if result.Result != nil && result.Result.Hits != nil {
		for _, hit := range result.Result.Hits.Hit {
			pub := DPLBPublication{
				Title:   hit.Info.Title,
				Authors: strings.Split(hit.Info.Authors.Author, ", "),
				Venue:   hit.Info.Venue,
				Year:    hit.Info.Year,
				URL:     hit.Info.URL,
			}
			publications = append(publications, pub)
		}
	}

	return publications, nil
}

// GitHubSearchResult represents GitHub search API response
type GitHubSearchResult struct {
	TotalCount int          `json:"total_count"`
	Items      []GitHubRepo `json:"items"`
}

// GitHubRepo represents a GitHub repository
type GitHubRepo struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	FullName    string   `json:"full_name"`
	Description string   `json:"description"`
	HTMLURL     string   `json:"html_url"`
	Stars       int      `json:"stargazers_count"`
	Forks       int      `json:"forks_count"`
	Language    string   `json:"language"`
	Topics      []string `json:"topics"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
	PushedAt    string   `json:"pushed_at"`
}

// DPLBSearchResult represents DBLP search result
type DPLBSearchResult struct {
	Result *struct {
		Hits *struct {
			Hit []struct {
				Info struct {
					Title   string `json:"title"`
					Authors struct {
						Author string `json:"author"`
					} `json:"authors"`
					Venue string `json:"venue"`
					Year  string `json:"year"`
					URL   string `json:"url"`
				} `json:"info"`
			} `json:"hit"`
		} `json:"hits"`
	} `json:"result"`
}

// DPLBPublication represents a DBLP publication
type DPLBPublication struct {
	Title   string
	Authors []string
	Venue   string
	Year    string
	URL     string
}

