package postgres

import (
	"time"

	"github.com/google/uuid"
	"github.com/virtualguard/PaperBeginer/internal/domain"
	"gorm.io/gorm"
)

// UserRepository implements domain.UserRepository
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id uuid.UUID) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update updates a user
func (r *UserRepository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

// Delete deletes a user (soft delete by deactivating)
func (r *UserRepository) Delete(id uuid.UUID) error {
	return r.db.Model(&domain.User{}).Where("id = ?", id).Update("is_active", false).Error
}

// List retrieves a paginated list of users
func (r *UserRepository) List(offset, limit int) ([]*domain.User, int64, error) {
	var users []*domain.User
	var total int64

	r.db.Model(&domain.User{}).Count(&total)

	err := r.db.Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// UpdateLastLogin updates the last login timestamp
func (r *UserRepository) UpdateLastLogin(id uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&domain.User{}).Where("id = ?", id).Update("last_login_at", now).Error
}

