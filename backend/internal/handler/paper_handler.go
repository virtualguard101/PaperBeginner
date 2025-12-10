package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/virtualguard/PaperBeginner/internal/handler/middleware"
	"github.com/virtualguard/PaperBeginner/internal/service"
	"github.com/virtualguard/PaperBeginner/pkg/response"
)

// PaperHandler handles paper-related HTTP requests
type PaperHandler struct {
	paperService *service.PaperService
}

// NewPaperHandler creates a new PaperHandler
func NewPaperHandler(paperService *service.PaperService) *PaperHandler {
	return &PaperHandler{
		paperService: paperService,
	}
}

// Upload handles paper upload
// @Summary Upload a paper
// @Description Upload a PDF paper for analysis
// @Tags papers
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formance file true "PDF file"
// @Param title formance string true "Paper title"
// @Param authors formance string false "Paper authors (comma-separated)"
// @Success 201 {object} response.Response{data=domain.PaperResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /papers/upload [post]
func (h *PaperHandler) Upload(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Unauthorized(c, "Invalid token")
		return
	}

	// Get file from request
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.BadRequest(c, "File is required")
		return
	}
	defer file.Close()

	// Get metadata from form
	title := c.PostForm("title")
	if title == "" {
		response.BadRequest(c, "Title is required")
		return
	}
	authors := c.PostForm("authors")
	abstract := c.PostForm("abstract")

	paper, err := h.paperService.Upload(userID, file, header, title, authors, abstract)
	if err != nil {
		response.InternalError(c, "Failed to upload paper")
		return
	}

	response.Created(c, paper.ToResponse())
}

// List lists user's papers
// @Summary List papers
// @Description Get a paginated list of user's papers
// @Tags papers
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(20)
// @Success 200 {object} response.Response{data=[]domain.PaperResponse}
// @Failure 401 {object} response.Response
// @Router /papers [get]
func (h *PaperHandler) List(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Unauthorized(c, "Invalid token")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	offset := (page - 1) * perPage

	papers, total, err := h.paperService.ListByUserID(userID, offset, perPage)
	if err != nil {
		response.InternalError(c, "Failed to list papers")
		return
	}

	paperResponses := make([]*interface{}, len(papers))
	for i, p := range papers {
		resp := p.ToResponse()
		paperResponses[i] = new(interface{})
		*paperResponses[i] = resp
	}

	response.Paginated(c, papers, page, perPage, total)
}

// Get gets a single paper
// @Summary Get paper
// @Description Get a paper by ID
// @Tags papers
// @Produce json
// @Security BearerAuth
// @Param id path string true "Paper ID"
// @Success 200 {object} response.Response{data=domain.PaperResponse}
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /papers/{id} [get]
func (h *PaperHandler) Get(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Unauthorized(c, "Invalid token")
		return
	}

	paperID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid paper ID")
		return
	}

	paper, err := h.paperService.GetByID(paperID, userID)
	if err != nil {
		if err == service.ErrPaperNotFound {
			response.NotFound(c, "Paper not found")
			return
		}
		response.InternalError(c, "Failed to get paper")
		return
	}

	response.Success(c, paper.ToResponse())
}

// Delete deletes a paper
// @Summary Delete paper
// @Description Delete a paper by ID
// @Tags papers
// @Produce json
// @Security BearerAuth
// @Param id path string true "Paper ID"
// @Success 204 "No Content"
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /papers/{id} [delete]
func (h *PaperHandler) Delete(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Unauthorized(c, "Invalid token")
		return
	}

	paperID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid paper ID")
		return
	}

	if err := h.paperService.Delete(paperID, userID); err != nil {
		if err == service.ErrPaperNotFound {
			response.NotFound(c, "Paper not found")
			return
		}
		response.InternalError(c, "Failed to delete paper")
		return
	}

	response.NoContent(c)
}

// Analyze triggers analysis for a paper
// @Summary Analyze paper
// @Description Trigger AI analysis for a paper
// @Tags papers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Paper ID"
// @Param request body AnalyzeRequest true "Analysis request"
// @Success 202 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /papers/{id}/analyze [post]
func (h *PaperHandler) Analyze(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Unauthorized(c, "Invalid token")
		return
	}

	paperID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid paper ID")
		return
	}

	var req AnalyzeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Default to all analysis types if not specified
		req.Types = []string{"summary", "methodology", "contributions"}
	}

	if err := h.paperService.QueueAnalysis(paperID, userID, req.Types); err != nil {
		if err == service.ErrPaperNotFound {
			response.NotFound(c, "Paper not found")
			return
		}
		response.InternalError(c, "Failed to queue analysis")
		return
	}

	response.Success(c, gin.H{
		"message": "Analysis queued successfully",
		"types":   req.Types,
	})
}

// GetAnalyses gets all analyses for a paper
// @Summary Get paper analyses
// @Description Get all AI analyses for a paper
// @Tags papers
// @Produce json
// @Security BearerAuth
// @Param id path string true "Paper ID"
// @Success 200 {object} response.Response{data=[]domain.PaperAnalysis}
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /papers/{id}/analyses [get]
func (h *PaperHandler) GetAnalyses(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Unauthorized(c, "Invalid token")
		return
	}

	paperID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid paper ID")
		return
	}

	analyses, err := h.paperService.GetAnalyses(paperID, userID)
	if err != nil {
		if err == service.ErrPaperNotFound {
			response.NotFound(c, "Paper not found")
			return
		}
		response.InternalError(c, "Failed to get analyses")
		return
	}

	response.Success(c, analyses)
}

// GetCCFCategories gets all CCF categories
// @Summary Get CCF categories
// @Description Get all CCF research categories
// @Tags papers
// @Produce json
// @Success 200 {object} response.Response{data=[]domain.CCFCategory}
// @Router /ccf/categories [get]
func (h *PaperHandler) GetCCFCategories(c *gin.Context) {
	categories, err := h.paperService.GetCCFCategories()
	if err != nil {
		response.InternalError(c, "Failed to get categories")
		return
	}

	response.Success(c, categories)
}

// AnalyzeRequest represents an analysis request
type AnalyzeRequest struct {
	Types []string `json:"types"` // summary, methodology, contributions, etc.
}
