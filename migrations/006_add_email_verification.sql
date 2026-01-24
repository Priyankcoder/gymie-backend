
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

-- Add comment
COMMENT ON COLUMN users.email_verified IS 'Whether the user has verified their email address';
COMMENT ON COLUMN users.verification_token IS 'Token for email verification (single use)';
COMMENT ON COLUMN users.verification_token_expires_at IS 'Expiration time for verification token';
