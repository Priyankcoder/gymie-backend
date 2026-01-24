
package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/yourusername/gymie-backend/internal/config"
	"github.com/yourusername/gymie-backend/internal/models"
	"github.com/yourusername/gymie-backend/internal/repository"
	"github.com/yourusername/gymie-backend/internal/services"
	"github.com/yourusername/gymie-backend/internal/utils"
	"gorm.io/gorm"
)

// AuthService defines the interface for authentication operations
type AuthService interface {
	Register(ctx context.Context, req *models.UserRegisterRequest) (*models.AuthResponse, error)
	Login(ctx context.Context, req *models.UserLoginRequest) (*models.AuthResponse, error)
	LoginWithGoogle(ctx context.Context, req *models.GoogleSignInRequest) (*models.AuthResponse, error)
	ValidateToken(tokenString string) (*utils.Claims, error)
}

type authService struct {
	userRepo     repository.UserRepository
	cfg          *config.Config
	emailService *services.EmailService
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo repository.UserRepository, cfg *config.Config, emailService *services.EmailService) AuthService {
	return &authService{
		userRepo:     userRepo,
		cfg:          cfg,
		emailService: emailService,
	}
}

// Register registers a new user
func (s *authService) Register(ctx context.Context, req *models.UserRegisterRequest) (*models.AuthResponse, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Generate verification token
	verificationToken, err := utils.GenerateVerificationToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate verification token: %w", err)
	}

	// Create user with email verification fields
	user := &models.User{
		Email:                      req.Email,
		Password:                   hashedPassword,
		Name:                       req.Name,
		EmailVerified:              false,
		VerificationToken:          verificationToken,
		VerificationTokenExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Create user profile
	profile := &models.UserProfile{
		UserID: user.ID,
	}
	if err := s.userRepo.UpdateProfile(ctx, profile); err != nil {
		return nil, fmt.Errorf("failed to create user profile: %w", err)
	}

	// Load profile
	user.Profile = profile

	// Send verification email
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:8081" // Default for development
	}
	verificationLink := fmt.Sprintf("%s/verify-email?token=%s", frontendURL, verificationToken)
	
	fmt.Printf("=== EMAIL DEBUG ===\n")
	fmt.Printf("Sending verification email to: %s\n", user.Email)
	fmt.Printf("User name: %s\n", user.Name)
	fmt.Printf("Verification link: %s\n", verificationLink)
	fmt.Printf("==================\n")
	
	if err := s.emailService.SendVerificationEmail(user.Email, user.Name, verificationLink); err != nil {
		// Log error but don't fail registration
		fmt.Printf("❌ FAILED to send verification email: %v\n", err)
	} else {
		fmt.Printf("✅ Verification email sent successfully to %s\n", user.Email)
	}

	// Return response WITHOUT token - user must verify email first
	return &models.AuthResponse{
		Token: "", // No token until email is verified
		User:  user.ToResponse(),
	}, nil
}

// Login authenticates a user and returns a token
func (s *authService) Login(ctx context.Context, req *models.UserLoginRequest) (*models.AuthResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("invalid email or password")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Compare password
	if err := utils.ComparePassword(user.Password, req.Password); err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Check if email is verified
	if !user.EmailVerified {
		return nil, fmt.Errorf("email_not_verified")
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID, user.Email, s.cfg.JWTSecret, s.cfg.JWTExpiration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &models.AuthResponse{
		Token: token,
		User:  user.ToResponse(),
	}, nil
}

// LoginWithGoogle authenticates a user via Google Sign-In
func (s *authService) LoginWithGoogle(ctx context.Context, req *models.GoogleSignInRequest) (*models.AuthResponse, error) {
	// Verify Google ID token
	tokenInfo, err := utils.VerifyGoogleIDToken(ctx, req.IDToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify Google token: %w", err)
	}

	// Use email from token (most reliable) or fallback to request
	email := tokenInfo.Email
	if email == "" && req.Email != nil {
		email = *req.Email
	}
	if email == "" {
		return nil, fmt.Errorf("email is required")
	}

	// Check if user exists
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}

	// If user doesn't exist, create new account
	if user == nil || errors.Is(err, gorm.ErrRecordNotFound) {
		// Use name from token or request
		name := tokenInfo.Name
		if name == "" && req.Name != nil {
			name = *req.Name
		}
		if name == "" {
			name = tokenInfo.GivenName + " " + tokenInfo.FamilyName
		}
		if name == "" {
			name = email // Fallback to email if no name available
		}

		// Create user with a random password (won't be used for Google sign-in)
		randomPassword, err := utils.GenerateRandomPassword(32)
		if err != nil {
			return nil, fmt.Errorf("failed to generate password: %w", err)
		}

		hashedPassword, err := utils.HashPassword(randomPassword)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}

		user = &models.User{
			Email:         email,
			Password:      hashedPassword,
			Name:          name,
			EmailVerified: true, // Google users are already verified
		}

		if err := s.userRepo.Create(ctx, user); err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}

		// Create user profile with Google profile picture if available
		profilePicture := tokenInfo.Picture
		profile := &models.UserProfile{
			UserID:         user.ID,
			ProfilePicture: &profilePicture,
		}
		if err := s.userRepo.UpdateProfile(ctx, profile); err != nil {
			return nil, fmt.Errorf("failed to create user profile: %w", err)
		}

		user.Profile = profile
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email, s.cfg.JWTSecret, s.cfg.JWTExpiration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &models.AuthResponse{
		Token: token,
		User:  user.ToResponse(),
	}, nil
}

// ValidateToken validates a JWT token
func (s *authService) ValidateToken(tokenString string) (*utils.Claims, error) {
	return utils.ValidateToken(tokenString, s.cfg.JWTSecret)
}
