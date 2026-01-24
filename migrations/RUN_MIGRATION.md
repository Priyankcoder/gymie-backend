
# Running Email Verification Migration

## Option 1: Using Go Migration Runner (Recommended)

The easiest way is to use the provided Go migration runner:

```bash
cd backend
go run cmd/migrate/main.go
```

**What it does:**
- ✅ Reads DATABASE_URL from your .env file
- ✅ Connects to the database securely
- ✅ Runs the migration SQL
- ✅ Creates necessary indexes
- ✅ Sets existing users as verified (backward compatibility)
- ✅ Shows clear success/error messages

**Output on success:**
```
📡 Connecting to database...
✅ Connected to database
🔄 Running email verification migration...
✅ Migration completed successfully!

📋 New columns added:
   - email_verified (boolean)
   - verification_token (varchar)
   - verification_token_expires_at (timestamp)
```

## Option 2: Using Docker

If you're running PostgreSQL in Docker:

```bash
cd backend

# Find your postgres container name
docker ps | grep postgres

# Run migration
docker exec -i YOUR_CONTAINER_NAME psql -U YOUR_USERNAME -d YOUR_DATABASE < migrations/006_add_email_verification.sql

# Example:
docker exec -i gymie-db psql -U postgres -d gymie < migrations/006_add_email_verification.sql
```

## Option 3: Using GORM Auto Migrate

Add to your main.go initialization:

```go
// In your database initialization code
db.AutoMigrate(&models.User{})

// Or manually:
db.Exec(`
	ALTER TABLE users 
	ADD COLUMN IF NOT EXISTS email_verified BOOLEAN DEFAULT FALSE,
	ADD COLUMN IF NOT EXISTS verification_token VARCHAR(255) UNIQUE,
	ADD COLUMN IF NOT EXISTS verification_token_expires_at TIMESTAMP;
`)
```

## Option 4: Using Database GUI Tools

### DBeaver / TablePlus / pgAdmin
1. Connect to your database
2. Open SQL editor
3. Copy contents of `migrations/006_add_email_verification.sql`
4. Execute

### Railway / Render Dashboard
1. Go to your database dashboard
2. Find "Query" or "SQL Editor"
3. Paste migration SQL
4. Execute

## Option 5: Using psql via Homebrew (Mac)

```bash
# Install PostgreSQL client
brew install postgresql

# Then run migration
psql $DATABASE_URL -f migrations/006_add_email_verification.sql
```

## Verify Migration Worked

Run this SQL to check:

```sql
SELECT 
    column_name, 
    data_type, 
    is_nullable, 
    column_default
FROM information_schema.columns
WHERE table_name = 'users'
AND column_name IN ('email_verified', 'verification_token', 'verification_token_expires_at');
```

Should return 3 rows showing the new columns.

## Rollback (If Needed)

```sql
ALTER TABLE users 
DROP COLUMN IF EXISTS email_verified,
DROP COLUMN IF EXISTS verification_token,
DROP COLUMN IF EXISTS verification_token_expires_at;

DROP INDEX IF EXISTS idx_users_verification_token;
DROP INDEX IF EXISTS idx_users_email_verified;
```
