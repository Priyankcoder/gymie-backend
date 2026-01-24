package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yourusername/gymie-backend/internal/models"
	"gorm.io/gorm"
)

// ProgressRepository defines the interface for progress data operations
type ProgressRepository interface {
	// Progress Photos
	CreateProgressPhoto(ctx context.Context, photo *models.ProgressPhoto) error
	GetProgressPhotoByID(ctx context.Context, id uint, userID uint) (*models.ProgressPhoto, error)
	GetProgressPhotos(ctx context.Context, userID uint, query *models.ListQuery) ([]models.ProgressPhoto, int64, error)
	UpdateProgressPhoto(ctx context.Context, photo *models.ProgressPhoto) error
	DeleteProgressPhoto(ctx context.Context, id uint, userID uint) error

	// Weight Entries
	CreateWeightEntry(ctx context.Context, entry *models.WeightEntry) error
	GetWeightEntryByID(ctx context.Context, id uint, userID uint) (*models.WeightEntry, error)
	GetWeightEntries(ctx context.Context, userID uint, query *models.ListQuery) ([]models.WeightEntry, int64, error)
	GetWeightEntriesByDateRange(ctx context.Context, userID uint, startDate, endDate time.Time) ([]models.WeightEntry, error)
	UpdateWeightEntry(ctx context.Context, entry *models.WeightEntry) error
	DeleteWeightEntry(ctx context.Context, id uint, userID uint) error

	// Stats
	GetWeightProgress(ctx context.Context, userID uint) (*models.WeightProgressResponse, error)
}

type progressRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

// NewProgressRepository creates a new progress repository
func NewProgressRepository(db *gorm.DB, redis *redis.Client) ProgressRepository {
	return &progressRepository{
		db:    db,
		redis: redis,
	}
}

// CreateProgressPhoto creates a new progress photo
func (r *progressRepository) CreateProgressPhoto(ctx context.Context, photo *models.ProgressPhoto) error {
	return r.db.WithContext(ctx).Create(photo).Error
}

// GetProgressPhotoByID retrieves a progress photo by ID
func (r *progressRepository) GetProgressPhotoByID(ctx context.Context, id uint, userID uint) (*models.ProgressPhoto, error) {
	var photo models.ProgressPhoto
	if err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		First(&photo).Error; err != nil {
		return nil, err
	}
	return &photo, nil
}

// GetProgressPhotos retrieves progress photos for a user with pagination
func (r *progressRepository) GetProgressPhotos(ctx context.Context, userID uint, query *models.ListQuery) ([]models.ProgressPhoto, int64, error) {
	query.SetDefaults()

	var photos []models.ProgressPhoto
	var total int64

	// Count total
	if err := r.db.WithContext(ctx).
		Model(&models.ProgressPhoto{}).
		Where("user_id = ?", userID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("date DESC").
		Offset(query.GetOffset()).
		Limit(query.PageSize).
		Find(&photos).Error; err != nil {
		return nil, 0, err
	}

	return photos, total, nil
}

// UpdateProgressPhoto updates a progress photo
func (r *progressRepository) UpdateProgressPhoto(ctx context.Context, photo *models.ProgressPhoto) error {
	return r.db.WithContext(ctx).Save(photo).Error
}

// DeleteProgressPhoto soft deletes a progress photo
func (r *progressRepository) DeleteProgressPhoto(ctx context.Context, id uint, userID uint) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&models.ProgressPhoto{}).Error
}

// CreateWeightEntry creates a new weight entry
func (r *progressRepository) CreateWeightEntry(ctx context.Context, entry *models.WeightEntry) error {
	return r.db.WithContext(ctx).Create(entry).Error
}

// GetWeightEntryByID retrieves a weight entry by ID
func (r *progressRepository) GetWeightEntryByID(ctx context.Context, id uint, userID uint) (*models.WeightEntry, error) {
	var entry models.WeightEntry
	if err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		First(&entry).Error; err != nil {
		return nil, err
	}
	return &entry, nil
}

// GetWeightEntries retrieves weight entries for a user with pagination
func (r *progressRepository) GetWeightEntries(ctx context.Context, userID uint, query *models.ListQuery) ([]models.WeightEntry, int64, error) {
	query.SetDefaults()

	var entries []models.WeightEntry
	var total int64

	// Count total
	if err := r.db.WithContext(ctx).
		Model(&models.WeightEntry{}).
		Where("user_id = ?", userID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("date DESC").
		Offset(query.GetOffset()).
		Limit(query.PageSize).
		Find(&entries).Error; err != nil {
		return nil, 0, err
	}

	return entries, total, nil
}

// GetWeightEntriesByDateRange retrieves weight entries within a date range
func (r *progressRepository) GetWeightEntriesByDateRange(ctx context.Context, userID uint, startDate, endDate time.Time) ([]models.WeightEntry, error) {
	var entries []models.WeightEntry

	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND date BETWEEN ? AND ?", userID, startDate, endDate).
		Order("date ASC").
		Find(&entries).Error; err != nil {
		return nil, err
	}

	return entries, nil
}

// UpdateWeightEntry updates a weight entry
func (r *progressRepository) UpdateWeightEntry(ctx context.Context, entry *models.WeightEntry) error {
	return r.db.WithContext(ctx).Save(entry).Error
}

// DeleteWeightEntry soft deletes a weight entry
func (r *progressRepository) DeleteWeightEntry(ctx context.Context, id uint, userID uint) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&models.WeightEntry{}).Error
}

// GetWeightProgress retrieves weight progress statistics
func (r *progressRepository) GetWeightProgress(ctx context.Context, userID uint) (*models.WeightProgressResponse, error) {
	var entries []models.WeightEntry

	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("date ASC").
		Find(&entries).Error; err != nil {
		return nil, err
	}

	response := &models.WeightProgressResponse{
		Entries: entries,
	}

	if len(entries) == 0 {
		return response, nil
	}

	// Get start and current weight
	response.StartWeight = &entries[0].Weight
	response.CurrentWeight = &entries[len(entries)-1].Weight

	// Get goal weight from user profile
	var profile models.UserProfile
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		First(&profile).Error; err == nil && profile.Weight != nil {
		response.GoalWeight = profile.Weight
	}

	// Calculate total change
	response.TotalChange = *response.CurrentWeight - *response.StartWeight

	// Calculate average weekly change
	if len(entries) > 1 {
		firstDate := entries[0].Date
		lastDate := entries[len(entries)-1].Date
		weeks := lastDate.Sub(firstDate).Hours() / 24 / 7
		if weeks > 0 {
			response.AverageChange = response.TotalChange / weeks
		}
	}

	return response, nil
}
