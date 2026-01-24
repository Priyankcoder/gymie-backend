
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Get database URL
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("❌ DATABASE_URL environment variable not set")
	}

	// Connect to database
	fmt.Println("📡 Connecting to database...")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

	fmt.Println("✅ Connected to database")
	fmt.Println("🔄 Running email verification migration...")

	// Run migration SQL
	sql := `
	-- Add email verification fields to users table
	ALTER TABLE users 
	ADD COLUMN IF NOT EXISTS email_verified BOOLEAN DEFAULT FALSE,
	ADD COLUMN IF NOT EXISTS verification_token VARCHAR(255) UNIQUE,
	ADD COLUMN IF NOT EXISTS verification_token_expires_at TIMESTAMP;

	-- Create indexes for faster lookups
	CREATE INDEX IF NOT EXISTS idx_users_verification_token ON users(verification_token);
	CREATE INDEX IF NOT EXISTS idx_users_email_verified ON users(email_verified);

	-- Update existing users to be verified (backward compatibility)
	UPDATE users SET email_verified = TRUE WHERE email_verified IS NULL;
	`

	if err := db.Exec(sql).Error; err != nil {
		log.Fatal("❌ Migration failed:", err)
	}

	fmt.Println("✅ Migration completed successfully!")
	fmt.Println("\n📋 New columns added:")
	fmt.Println("   - email_verified (boolean)")
	fmt.Println("   - verification_token (varchar)")
	fmt.Println("   - verification_token_expires_at (timestamp)")
	fmt.Println("\n🎯 Next steps:")
	fmt.Println("   1. Install email package: go get gopkg.in/mail.v2")
	fmt.Println("   2. Configure SMTP in .env")
	fmt.Println("   3. Update main.go and auth_service.go")
	fmt.Println("   4. Test registration flow")
}
