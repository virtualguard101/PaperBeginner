package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/virtualguard/PaperBeginer/internal/config"
	"github.com/virtualguard/PaperBeginer/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserExists         = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// UserService handles user-related business logic
type UserService struct {
	userRepo  domain.UserRepository
	jwtConfig config.JWTConfig
}

// NewUserService creates a new UserService
func NewUserService(userRepo domain.UserRepository, jwtConfig config.JWTConfig) *UserService {
	return &UserService{
		userRepo:  userRepo,
		jwtConfig: jwtConfig,
	}
}

// Register creates a new user
func (s *UserService) Register(req *domain.CreateUserRequest) (*domain.User, error) {
	// Check if user already exists
	existing, _ := s.userRepo.GetByEmail(req.Email)
	if existing != nil {
		return nil, ErrUserExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Name:         req.Name,
		Role:         domain.RoleUser,
		IsActive:     true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login authenticates a user and returns a JWT token
func (s *UserService) Login(req *domain.LoginRequest) (*domain.LoginResponse, error) {
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, ErrInvalidCredentials
	}

	// Update last login
	_ = s.userRepo.UpdateLastLogin(user.ID)

	// Generate token
	token, expiresAt, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &domain.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      user,
	}, nil
}

// RefreshToken generates a new token for a user
func (s *UserService) RefreshToken(userID uuid.UUID) (*domain.LoginResponse, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	token, expiresAt, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &domain.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      user,
	}, nil
}

// GetByID retrieves a user by ID
func (s *UserService) GetByID(id uuid.UUID) (*domain.User, error) {
	return s.userRepo.GetByID(id)
}

// Update updates a user's profile
func (s *UserService) Update(id uuid.UUID, req *domain.UpdateUserRequest) (*domain.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Avatar != nil {
		user.Avatar = *req.Avatar
	}
	if req.Preferences != nil {
		user.Preferences = req.Preferences
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

// ChangePassword changes a user's password
func (s *UserService) ChangePassword(id uuid.UUID, oldPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return ErrInvalidCredentials
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashedPassword)
	return s.userRepo.Update(user)
}

// generateToken generates a JWT token for a user
func (s *UserService) generateToken(user *domain.User) (string, int64, error) {
	expiresAt := time.Now().Add(s.jwtConfig.Expiry).Unix()

	claims := jwt.MapClaims{
		"sub":   user.ID.String(),
		"email": user.Email,
		"role":  string(user.Role),
		"exp":   expiresAt,
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.jwtConfig.Secret))
	if err != nil {
		return "", 0, err
	}

	return signedToken, expiresAt, nil
}

