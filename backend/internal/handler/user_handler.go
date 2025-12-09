package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/virtualguard/PaperBeginer/internal/domain"
	"github.com/virtualguard/PaperBeginer/internal/handler/middleware"
	"github.com/virtualguard/PaperBeginer/internal/service"
	"github.com/virtualguard/PaperBeginer/pkg/response"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Register handles user registration
// @Summary Register a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body domain.CreateUserRequest true "Registration request"
// @Success 201 {object} response.Response{data=domain.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 409 {object} response.Response
// @Router /auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req domain.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	user, err := h.userService.Register(&req)
	if err != nil {
		if err == service.ErrUserExists {
			response.Conflict(c, "User with this email already exists")
			return
		}
		response.InternalError(c, "Failed to create user")
		return
	}

	response.Created(c, user.ToResponse())
}

// Login handles user authentication
// @Summary Login
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body domain.LoginRequest true "Login request"
// @Success 200 {object} response.Response{data=domain.LoginResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	loginResp, err := h.userService.Login(&req)
	if err != nil {
		if err == service.ErrInvalidCredentials {
			response.Unauthorized(c, "Invalid email or password")
			return
		}
		response.InternalError(c, "Failed to authenticate")
		return
	}

	response.Success(c, loginResp)
}

// RefreshToken handles token refresh
// @Summary Refresh token
// @Description Refresh JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} response.Response{data=domain.LoginResponse}
// @Failure 401 {object} response.Response
// @Router /auth/refresh [post]
func (h *UserHandler) RefreshToken(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Unauthorized(c, "Invalid token")
		return
	}

	loginResp, err := h.userService.RefreshToken(userID)
	if err != nil {
		response.InternalError(c, "Failed to refresh token")
		return
	}

	response.Success(c, loginResp)
}

// GetCurrentUser gets the current authenticated user
// @Summary Get current user
// @Description Get the profile of the currently authenticated user
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=domain.UserResponse}
// @Failure 401 {object} response.Response
// @Router /users/me [get]
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Unauthorized(c, "Invalid token")
		return
	}

	user, err := h.userService.GetByID(userID)
	if err != nil {
		response.NotFound(c, "User not found")
		return
	}

	response.Success(c, user.ToResponse())
}

// UpdateCurrentUser updates the current user's profile
// @Summary Update current user
// @Description Update the profile of the currently authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.UpdateUserRequest true "Update request"
// @Success 200 {object} response.Response{data=domain.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /users/me [put]
func (h *UserHandler) UpdateCurrentUser(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Unauthorized(c, "Invalid token")
		return
	}

	var req domain.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	user, err := h.userService.Update(userID, &req)
	if err != nil {
		response.InternalError(c, "Failed to update user")
		return
	}

	response.Success(c, user.ToResponse())
}

// ChangePassword changes the user's password
// @Summary Change password
// @Description Change the password of the currently authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ChangePasswordRequest true "Change password request"
// @Success 204 "No Content"
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /users/me/password [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Unauthorized(c, "Invalid token")
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	if err := h.userService.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		if err == service.ErrInvalidCredentials {
			response.BadRequest(c, "Current password is incorrect")
			return
		}
		response.InternalError(c, "Failed to change password")
		return
	}

	response.NoContent(c)
}

// ChangePasswordRequest represents a password change request
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

