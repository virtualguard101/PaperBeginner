package domain

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID           uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Email        string     `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash string     `json:"-" gorm:"not null"`
	Name         string     `json:"name" gorm:"not null"`
	Avatar       string     `json:"avatar,omitempty"`
	Role         UserRole   `json:"role" gorm:"default:user"`
	IsActive     bool       `json:"is_active" gorm:"default:true"`
	Preferences  *UserPrefs `json:"preferences,omitempty" gorm:"type:jsonb"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
}

// UserRole defines the role of a user
type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleAdmin UserRole = "admin"
)

// UserPrefs stores user preferences
type UserPrefs struct {
	PreferredLLM      string   `json:"preferred_llm,omitempty"`
	ResearchFields    []string `json:"research_fields,omitempty"`
	NotifyOnTrending  bool     `json:"notify_on_trending"`
	Language          string   `json:"language,omitempty"`
	OpenAIAPIKey      string   `json:"openai_api_key,omitempty"`
	AnthropicAPIKey   string   `json:"anthropic_api_key,omitempty"`
	DeepSeekAPIKey    string   `json:"deepseek_api_key,omitempty"`
}

// UserRepository defines the interface for user data access
type UserRepository interface {
	Create(user *User) error
	GetByID(id uuid.UUID) (*User, error)
	GetByEmail(email string) (*User, error)
	Update(user *User) error
	Delete(id uuid.UUID) error
	List(offset, limit int) ([]*User, int64, error)
	UpdateLastLogin(id uuid.UUID) error
}

// CreateUserRequest represents the request to create a user
type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required,min=2,max=100"`
}

// UpdateUserRequest represents the request to update a user
type UpdateUserRequest struct {
	Name        *string    `json:"name,omitempty" binding:"omitempty,min=2,max=100"`
	Avatar      *string    `json:"avatar,omitempty"`
	Preferences *UserPrefs `json:"preferences,omitempty"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
	User      *User  `json:"user"`
}

// UserResponse represents a user in API responses
type UserResponse struct {
	ID          uuid.UUID  `json:"id"`
	Email       string     `json:"email"`
	Name        string     `json:"name"`
	Avatar      string     `json:"avatar,omitempty"`
	Role        UserRole   `json:"role"`
	IsActive    bool       `json:"is_active"`
	Preferences *UserPrefs `json:"preferences,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
}

// ToResponse converts a User to UserResponse
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:          u.ID,
		Email:       u.Email,
		Name:        u.Name,
		Avatar:      u.Avatar,
		Role:        u.Role,
		IsActive:    u.IsActive,
		Preferences: u.Preferences,
		CreatedAt:   u.CreatedAt,
		LastLoginAt: u.LastLoginAt,
	}
}

// TableName returns the table name for GORM
func (User) TableName() string {
	return "users"
}

