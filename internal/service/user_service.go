
package service

import (
	"context"
	"fmt"

	"github.com/yourusername/gymie-backend/internal/models"
	"github.com/yourusername/gymie-backend/internal/repository"
)

// UserService defines the interface for user operations
type UserService interface {
	GetByID(ctx context.Context, id uint) (*models.User, error)
	Update(ctx context.Context, userID uint, req *models.UserUpdateRequest) (*models.User, error)
	Delete(ctx context.Context, id uint) error
}

type userService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// GetByID retrieves a user by ID
func (s *userService) GetByID(ctx context.Context, id uint) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// Update updates a user's information
func (s *userService) Update(ctx context.Context, userID uint, req *models.UserUpdateRequest) (*models.User, error) {
	// Get existing user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Update user fields
	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Email != nil {
		user.Email = *req.Email
	}

	// Save user
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Update profile if provided
	if req.Height != nil || req.Weight != nil || req.Age != nil || req.Gender != nil || req.Goal != nil {
		profile := user.Profile
		if profile == nil {
			profile = &models.UserProfile{
				UserID: userID,
			}
		}

		if req.Height != nil {
			profile.Height = req.Height
		}
		if req.Weight != nil {
			profile.Weight = req.Weight
		}
		if req.Age != nil {
			profile.Age = req.Age
		}
		if req.Gender != nil {
			profile.Gender = *req.Gender
		}
		if req.Goal != nil {
			profile.Goal = *req.Goal
		}

		if err := s.userRepo.UpdateProfile(ctx, profile); err != nil {
			return nil, fmt.Errorf("failed to update profile: %w", err)
		}

		user.Profile = profile
	}

	return user, nil
}

// Delete deletes a user
func (s *userService) Delete(ctx context.Context, id uint) error {
	if err := s.userRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
