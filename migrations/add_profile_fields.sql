
-- Migration: Add display_name, profile_picture, and bio fields to user_profiles table
-- Created: 2026-02-01

-- Add new columns to user_profiles table
ALTER TABLE user_profiles 
  ADD COLUMN IF NOT EXISTS display_name VARCHAR(255),
  ADD COLUMN IF NOT EXISTS profile_picture TEXT,
  ADD COLUMN IF NOT EXISTS bio TEXT;

-- Create index on display_name for faster searches
CREATE INDEX IF NOT EXISTS idx_user_profiles_display_name ON user_profiles(display_name);

-- Update existing profiles to set display_name from users.name if null
UPDATE user_profiles up
SET display_name = u.name
FROM users u
WHERE up.user_id = u.id AND up.display_name IS NULL;
