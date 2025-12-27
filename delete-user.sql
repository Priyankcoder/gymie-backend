
-- Delete user and all related data for priyankrastogi14@gmail.com
-- Run this with: psql -U postgres -d gymie_dev -f delete-user.sql

BEGIN;

-- Delete user's data (cascading deletes should handle most of this)
DELETE FROM users WHERE email = 'priyankrastogi14@gmail.com';

COMMIT;

-- Verify deletion
SELECT email FROM users WHERE email = 'priyankrastogi14@gmail.com';
