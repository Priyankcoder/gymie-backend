package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/gymie-backend/internal/config"
	"github.com/yourusername/gymie-backend/internal/models"
	"github.com/yourusername/gymie-backend/internal/repository"
)

// ProgressService defines the interface for progress operations
type ProgressService interface {
	// Progress Photos
	CreateProgressPhoto(ctx context.Context, userID uint, req *models.ProgressPhotoCreateRequest) (*models.ProgressPhoto, error)
	UploadProgressPhoto(ctx context.Context, userID uint, file *multipart.FileHeader, date time.Time, weight *float64, notes string) (*models.ProgressPhoto, error)
	GetProgressPhotoByID(ctx context.Context, id uint, userID uint) (*models.ProgressPhoto, error)
	GetProgressPhotos(ctx context.Context, userID uint, query *models.ListQuery) ([]models.ProgressPhoto, int64, error)
	UpdateProgressPhoto(ctx context.Context, id uint, userID uint, req *models.ProgressPhotoUpdateRequest) (*models.ProgressPhoto, error)
	DeleteProgressPhoto(ctx context.Context, id uint, userID uint) error
	GenerateUploadURL(ctx context.Context, userID uint, filename string) (*models.UploadURLResponse, error)

	// Weight Entries
	CreateWeightEntry(ctx context.Context, userID uint, req *models.WeightEntryCreateRequest) (*models.WeightEntry, error)
	GetWeightEntryByID(ctx context.Context, id uint, userID uint) (*models.WeightEntry, error)
	GetWeightEntries(ctx context.Context, userID uint, query *models.ListQuery) ([]models.WeightEntry, int64, error)
	GetWeightEntriesByDateRange(ctx context.Context, userID uint, startDate, endDate time.Time) ([]models.WeightEntry, error)
	UpdateWeightEntry(ctx context.Context, id uint, userID uint, req *models.WeightEntryUpdateRequest) (*models.WeightEntry, error)
	DeleteWeightEntry(ctx context.Context, id uint, userID uint) error
	GetWeightProgress(ctx context.Context, userID uint) (*models.WeightProgressResponse, error)
}

type progressService struct {
	progressRepo repository.ProgressRepository
	cfg          *config.Config
}

// NewProgressService creates a new progress service
func NewProgressService(progressRepo repository.ProgressRepository, cfg *config.Config) ProgressService {
	return &progressService{
		progressRepo: progressRepo,
		cfg:          cfg,
	}
}

// UploadProgressPhoto uploads and creates a new progress photo
func (s *progressService) UploadProgressPhoto(ctx context.Context, userID uint, file *multipart.FileHeader, date time.Time, weight *float64, notes string) (*models.ProgressPhoto, error) {
	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	// Create upload directory if it doesn't exist
	uploadDir := filepath.Join(s.cfg.StorageBasePath, "progress_photos", fmt.Sprintf("%d", userID))
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	// Full file path
	filePath := filepath.Join(uploadDir, filename)

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	// Copy file content
	if _, err := io.Copy(dst, src); err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// Generate public URL
	imageURL := fmt.Sprintf("%s/progress_photos/%d/%s", s.cfg.StoragePublicURL, userID, filename)

	// Create database record
	photo := &models.ProgressPhoto{
		UserID:   userID,
		Date:     date,
		ImageURL: imageURL,
		Weight:   weight,
		Notes:    notes,
	}

	if err := s.progressRepo.CreateProgressPhoto(ctx, photo); err != nil {
		// Clean up file if database insert fails
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to create progress photo record: %w", err)
	}

	return photo, nil
}

// CreateProgressPhoto creates a new progress photo
func (s *progressService) CreateProgressPhoto(ctx context.Context, userID uint, req *models.ProgressPhotoCreateRequest) (*models.ProgressPhoto, error) {
	photo := &models.ProgressPhoto{
		UserID:   userID,
		Date:     req.Date,
		Weight:   req.Weight,
		Notes:    req.Notes,
		ImageURL: "", // Will be set after upload
	}

	if err := s.progressRepo.CreateProgressPhoto(ctx, photo); err != nil {
		return nil, fmt.Errorf("failed to create progress photo: %w", err)
	}

	return photo, nil
}

// GetProgressPhotoByID retrieves a progress photo by ID
func (s *progressService) GetProgressPhotoByID(ctx context.Context, id uint, userID uint) (*models.ProgressPhoto, error) {
	photo, err := s.progressRepo.GetProgressPhotoByID(ctx, id, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get progress photo: %w", err)
	}
	return photo, nil
}

// GetProgressPhotos retrieves progress photos for a user
func (s *progressService) GetProgressPhotos(ctx context.Context, userID uint, query *models.ListQuery) ([]models.ProgressPhoto, int64, error) {
	photos, total, err := s.progressRepo.GetProgressPhotos(ctx, userID, query)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get progress photos: %w", err)
	}
	return photos, total, nil
}

// UpdateProgressPhoto updates a progress photo
func (s *progressService) UpdateProgressPhoto(ctx context.Context, id uint, userID uint, req *models.ProgressPhotoUpdateRequest) (*models.ProgressPhoto, error) {
	photo, err := s.progressRepo.GetProgressPhotoByID(ctx, id, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get progress photo: %w", err)
	}

	if req.Date != nil {
		photo.Date = *req.Date
	}
	if req.Weight != nil {
		photo.Weight = req.Weight
	}
	if req.Notes != nil {
		photo.Notes = *req.Notes
	}

	if err := s.progressRepo.UpdateProgressPhoto(ctx, photo); err != nil {
		return nil, fmt.Errorf("failed to update progress photo: %w", err)
	}

	return photo, nil
}

// DeleteProgressPhoto deletes a progress photo
func (s *progressService) DeleteProgressPhoto(ctx context.Context, id uint, userID uint) error {
	// Get photo to retrieve file path
	photo, err := s.progressRepo.GetProgressPhotoByID(ctx, id, userID)
	if err != nil {
		return fmt.Errorf("failed to get progress photo: %w", err)
	}

	// Delete from database
	if err := s.progressRepo.DeleteProgressPhoto(ctx, id, userID); err != nil {
		return fmt.Errorf("failed to delete progress photo: %w", err)
	}

	// Delete file from storage
	if photo.ImageURL != "" {
		// Extract filename from URL
		parts := strings.Split(photo.ImageURL, "/")
		if len(parts) >= 3 {
			filename := parts[len(parts)-1]
			filePath := filepath.Join(s.cfg.StorageBasePath, "progress_photos", fmt.Sprintf("%d", userID), filename)
			// Ignore error if file doesn't exist
			os.Remove(filePath)
		}
	}

	return nil
}

// GenerateUploadURL generates a presigned URL for uploading images
func (s *progressService) GenerateUploadURL(ctx context.Context, userID uint, filename string) (*models.UploadURLResponse, error) {
	// TODO: Implement S3/R2 presigned URL generation
	// This is a placeholder implementation
	uploadURL := fmt.Sprintf("%s/uploads/%d/%s", s.cfg.StoragePublicURL, userID, filename)
	fileURL := uploadURL
	expiresAt := time.Now().Add(time.Duration(s.cfg.StoragePresignExpiry) * time.Minute)

	return &models.UploadURLResponse{
		UploadURL: uploadURL,
		FileURL:   fileURL,
		ExpiresAt: expiresAt,
	}, nil
}

// CreateWeightEntry creates a new weight entry
func (s *progressService) CreateWeightEntry(ctx context.Context, userID uint, req *models.WeightEntryCreateRequest) (*models.WeightEntry, error) {
	entry := &models.WeightEntry{
		UserID: userID,
		Date:   req.Date,
		Weight: req.Weight,
		Notes:  req.Notes,
	}

	if err := s.progressRepo.CreateWeightEntry(ctx, entry); err != nil {
		return nil, fmt.Errorf("failed to create weight entry: %w", err)
	}

	return entry, nil
}

// GetWeightEntryByID retrieves a weight entry by ID
func (s *progressService) GetWeightEntryByID(ctx context.Context, id uint, userID uint) (*models.WeightEntry, error) {
	entry, err := s.progressRepo.GetWeightEntryByID(ctx, id, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get weight entry: %w", err)
	}
	return entry, nil
}

// GetWeightEntries retrieves weight entries for a user
func (s *progressService) GetWeightEntries(ctx context.Context, userID uint, query *models.ListQuery) ([]models.WeightEntry, int64, error) {
	entries, total, err := s.progressRepo.GetWeightEntries(ctx, userID, query)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get weight entries: %w", err)
	}
	return entries, total, nil
}

// GetWeightEntriesByDateRange retrieves weight entries within a date range
func (s *progressService) GetWeightEntriesByDateRange(ctx context.Context, userID uint, startDate, endDate time.Time) ([]models.WeightEntry, error) {
	entries, err := s.progressRepo.GetWeightEntriesByDateRange(ctx, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get weight entries by date range: %w", err)
	}
	return entries, nil
}

// UpdateWeightEntry updates a weight entry
func (s *progressService) UpdateWeightEntry(ctx context.Context, id uint, userID uint, req *models.WeightEntryUpdateRequest) (*models.WeightEntry, error) {
	entry, err := s.progressRepo.GetWeightEntryByID(ctx, id, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get weight entry: %w", err)
	}

	if req.Date != nil {
		entry.Date = *req.Date
	}
	if req.Weight != nil {
		entry.Weight = *req.Weight
	}
	if req.Notes != nil {
		entry.Notes = *req.Notes
	}

	if err := s.progressRepo.UpdateWeightEntry(ctx, entry); err != nil {
		return nil, fmt.Errorf("failed to update weight entry: %w", err)
	}

	return entry, nil
}

// DeleteWeightEntry deletes a weight entry
func (s *progressService) DeleteWeightEntry(ctx context.Context, id uint, userID uint) error {
	if err := s.progressRepo.DeleteWeightEntry(ctx, id, userID); err != nil {
		return fmt.Errorf("failed to delete weight entry: %w", err)
	}
	return nil
}

// GetWeightProgress retrieves weight progress statistics
func (s *progressService) GetWeightProgress(ctx context.Context, userID uint) (*models.WeightProgressResponse, error) {
	progress, err := s.progressRepo.GetWeightProgress(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get weight progress: %w", err)
	}
	return progress, nil
}
