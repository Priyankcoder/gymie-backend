package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	// Server
	Port        string
	Environment string

	// Database
	DatabaseURL string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string

	// Redis
	RedisURL  string
	RedisHost string
	RedisPort string

	// JWT
	JWTSecret     string
	JWTExpiration int // in hours

	// Storage (Cloudflare R2 / S3-compatible)
	StorageEndpoint      string
	StorageAccessKey     string
	StorageSecretKey     string
	StorageBucket        string
	StorageRegion        string
	StoragePublicURL     string
	StorageBasePath      string // Local file storage base path
	StoragePresignExpiry int    // in minutes

	// Rate Limiting
	RateLimitEnabled  bool
	RateLimitRequests int
	RateLimitWindow   int // in seconds

	// CORS
	CORSAllowedOrigins []string

	// SMTP Email Configuration (legacy - Gmail)
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	SMTPFrom     string
	
	// SendGrid Configuration (preferred)
	SendGridAPIKey string
	
	// Common Email Configuration
	FromEmail   string
	FromName    string
	FrontendURL string
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists (for local development)
	_ = godotenv.Load()

	cfg := &Config{
		// Server
		Port:        getEnv("PORT", "8080"),
		Environment: getEnv("ENVIRONMENT", "development"),

		// Database
		DatabaseURL: getEnv("DATABASE_URL", ""),
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBUser:      getEnv("DB_USER", "postgres"),
		DBPassword:  getEnv("DB_PASSWORD", "postgres"),
		DBName:      getEnv("DB_NAME", "gymie_dev"),

		// Redis
		RedisURL:  getEnv("REDIS_URL", ""),
		RedisHost: getEnv("REDIS_HOST", "localhost"),
		RedisPort: getEnv("REDIS_PORT", "6379"),

		// JWT
		JWTSecret:     getEnv("JWT_SECRET", "your-secret-key-change-this"),
		JWTExpiration: getEnvAsInt("JWT_EXPIRATION", 24), // 24 hours

		// Storage
		StorageEndpoint:      getEnv("STORAGE_ENDPOINT", ""),
		StorageAccessKey:     getEnv("STORAGE_ACCESS_KEY", ""),
		StorageSecretKey:     getEnv("STORAGE_SECRET_KEY", ""),
		StorageBucket:        getEnv("STORAGE_BUCKET", "gymie-dev"),
		StorageRegion:        getEnv("STORAGE_REGION", "auto"),
		StoragePublicURL:     getEnv("STORAGE_PUBLIC_URL", "http://localhost:8080/uploads"),
		StorageBasePath:      getEnv("STORAGE_BASE_PATH", "./uploads"),
		StoragePresignExpiry: getEnvAsInt("STORAGE_PRESIGN_EXPIRY", 15), // 15 minutes

		// Rate Limiting
		RateLimitEnabled:  getEnvAsBool("RATE_LIMIT_ENABLED", true),
		RateLimitRequests: getEnvAsInt("RATE_LIMIT_REQUESTS", 100),
		RateLimitWindow:   getEnvAsInt("RATE_LIMIT_WINDOW", 60), // 60 seconds

		// CORS
		CORSAllowedOrigins: []string{
			getEnv("CORS_ALLOWED_ORIGINS", "*"),
		},

		// SMTP Email (legacy - Gmail)
		SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:     getEnv("SMTP_PORT", "587"),
		SMTPUsername: getEnv("SMTP_USERNAME", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		SMTPFrom:     getEnv("SMTP_FROM", getEnv("FROM_EMAIL", "")), // Support both formats
		
		// SendGrid (preferred)
		SendGridAPIKey: getEnv("SENDGRID_API_KEY", ""),
		
		// Common Email Configuration
		FromEmail:   getEnv("FROM_EMAIL", "noreply@gymie.com"),
		FromName:    getEnv("FROM_NAME", "Gymie"),
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),
	}

	// Validate required fields
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// validate checks if required configuration values are present
func (c *Config) validate() error {
	if c.JWTSecret == "your-secret-key-change-this" && c.Environment == "production" {
		return fmt.Errorf("JWT_SECRET must be set in production")
	}

	return nil
}

// Helper functions

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}
