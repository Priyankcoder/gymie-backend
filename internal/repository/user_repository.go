package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yourusername/gymie-backend/internal/models"
	"gorm.io/gorm"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uint) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uint) error
	UpdateProfile(ctx context.Context, profile *models.UserProfile) error
}

type userRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB, redis *redis.Client) UserRepository {
	return &userRepository{
		db:    db,
		redis: redis,
	}
}

// Create creates a new user
func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// GetByID retrieves a user by ID with caching
func (r *userRepository) GetByID(ctx context.Context, id uint) (*models.User, error) {
	cacheKey := fmt.Sprintf("user:%d", id)

	// Try cache first
	cached, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var user models.User
		if err := json.Unmarshal([]byte(cached), &user); err == nil {
			return &user, nil
		}
	}

	// Fetch from database
	var user models.User
	if err := r.db.WithContext(ctx).
		Preload("Profile").
		First(&user, id).Error; err != nil {
		return nil, err
	}

	// Cache the result
	if data, err := json.Marshal(user); err == nil {
		r.redis.Set(ctx, cacheKey, data, 30*time.Minute)
	}

	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).
		Preload("Profile").
		Where("email = ?", email).
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Update updates a user and invalidates cache
func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return err
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("user:%d", user.ID)
	r.redis.Del(ctx, cacheKey)

	return nil
}

// Delete soft deletes a user and invalidates cache
func (r *userRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&models.User{}, id).Error; err != nil {
		return err
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("user:%d", id)
	r.redis.Del(ctx, cacheKey)

	return nil
}

// UpdateProfile updates or creates a user profile
func (r *userRepository) UpdateProfile(ctx context.Context, profile *models.UserProfile) error {
	if err := r.db.WithContext(ctx).Save(profile).Error; err != nil {
		return err
	}

	// Invalidate user cache
	cacheKey := fmt.Sprintf("user:%d", profile.UserID)
	r.redis.Del(ctx, cacheKey)

	return nil
}
