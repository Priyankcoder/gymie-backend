
# What GORM AutoMigrate Does (Explained Simply)

## The Simple Answer

`db.AutoMigrate(&models.User{})` looks at your Go struct definition and makes sure your database table matches it.

## Example: User Model to Database Table

### Your Go Code (models/user.go)
```go
type User struct {
    ID        uint      `gorm:"primaryKey"`
    Email     string    `gorm:"uniqueIndex;not null"`
    Password  string    `gorm:"not null"`
    Name      string    `gorm:"not null"`
    
    // NEW: Email verification fields
    EmailVerified              bool      `gorm:"default:false"`
    VerificationToken          string    `gorm:"index"`
    VerificationTokenExpiresAt time.Time
}
```

### What AutoMigrate Does

When you run AutoMigrate, GORM:

1. **Checks if table exists**
   - If NO → Creates the entire `users` table
   - If YES → Goes to step 2

2. **Checks if columns exist**
   - Compares Go struct fields with database columns
   - If a field exists in Go but NOT in database → **ADDS the column**
   - If a column exists in database but NOT in Go → **DOES NOTHING** (keeps it)

3. **Checks indexes**
   - Creates missing indexes
   - Example: `gorm:"uniqueIndex"` → Creates unique index on email

4. **Checks constraints**
   - Adds NOT NULL constraints
   - Adds unique constraints
   - Adds foreign keys

## What AutoMigrate DOES ✅

### 1. Creates New Tables
If `users` table doesn't exist, creates:
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    email_verified BOOLEAN DEFAULT false,
    verification_token VARCHAR(255),
    verification_token_expires_at TIMESTAMP,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);
```

### 2. Adds New Columns to Existing Tables
If `users` table exists but missing `email_verified` column:
```sql
ALTER TABLE users 
ADD COLUMN email_verified BOOLEAN DEFAULT false;

ALTER TABLE users 
ADD COLUMN verification_token VARCHAR(255);

ALTER TABLE users 
ADD COLUMN verification_token_expires_at TIMESTAMP;
```

### 3. Creates Indexes
```sql
CREATE UNIQUE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_verification_token ON users(verification_token);
```

### 4. Creates Foreign Keys
Example: WorkoutPlanDay has a foreign key to WorkoutPlan
```sql
ALTER TABLE workout_plan_days
ADD CONSTRAINT fk_workout_plans_days 
FOREIGN KEY (workout_plan_id) 
REFERENCES workout_plans(id);
```

## What AutoMigrate DOES NOT Do ❌

### 1. ❌ Does NOT Delete Columns
If you remove a field from your Go struct, the database column stays:
```go
// OLD: Had Username field
type User struct {
    Username string // Removed this
}

// Database still has 'username' column
// AutoMigrate will NOT delete it
```

### 2. ❌ Does NOT Delete Tables
If you remove a model from migration, table stays in database.

### 3. ❌ Does NOT Modify Column Types
If you change a field type:
```go
// OLD:
Age int

// NEW:
Age string  // AutoMigrate won't change int to string
```
**You need a manual migration for type changes.**

### 4. ❌ Does NOT Rename Columns
If you rename a field:
```go
// OLD:
Email string

// NEW: 
EmailAddress string  // AutoMigrate creates NEW column, keeps old one
```

### 5. ❌ Does NOT Delete or Modify Data
- Never deletes existing rows
- Never modifies existing data
- Only changes table structure

## Real Example: What Happens When You Deploy

### Before Deployment
Your production `users` table:
```
users table:
├── id (int)
├── email (string)
├── password (string)
├── name (string)
├── created_at (timestamp)
└── updated_at (timestamp)
```

Existing users:
```
id | email           | password | name
1  | john@email.com  | hash123  | John
2  | jane@email.com  | hash456  | Jane
```

### After AutoMigrate Runs

1. **Table structure updated:**
```
users table:
├── id (int)
├── email (string)
├── password (string)
├── name (string)
├── email_verified (bool) ← NEW!
├── verification_token (string) ← NEW!
├── verification_token_expires_at (timestamp) ← NEW!
├── created_at (timestamp)
└── updated_at (timestamp)
```

2. **Existing data preserved:**
```
id | email           | password | name | email_verified
1  | john@email.com  | hash123  | John | false (default)
2  | jane@email.com  | hash456  | Jane | false (default)
```

3. **New users work correctly:**
```sql
INSERT INTO users (email, password, name, email_verified)
VALUES ('new@email.com', 'hash789', 'New User', false);
```

## All 17 Models We're Migrating

### What Happens for Each Model:

1. **User** → Adds email_verified, verification_token columns
2. **UserProfile** → Creates table if doesn't exist
3. **Workout** → Creates table if doesn't exist
4. **Exercise** → Creates table if doesn't exist
5. **WorkoutSet** → Creates table if doesn't exist
6. **NutritionDay** → Creates table if doesn't exist
7. **Meal** → Creates table if doesn't exist
8. **Food** → Creates table if doesn't exist
9. **ProgressPhoto** → Creates table if doesn't exist
10. **WeightEntry** → Creates table if doesn't exist
11. **WorkoutPlan** → Creates NEW table (workout planning feature)
12. **WorkoutPlanDay** → Creates NEW table
13. **ScheduledWorkout** → Creates NEW table
14. **DishMaster** → Creates NEW table (offline nutrition)
15. **DishNutritionMaster** → Creates NEW table
16. **UserCorrection** → Creates NEW table
17. **ModelVersion** → Creates NEW table

## Safety Guarantees

### ✅ Safe Operations (AutoMigrate Does)
- Creating new tables
- Adding new columns
- Creating indexes
- Adding foreign keys
- Setting default values on new columns

### ⚠️ Requires Manual Migration
- Dropping columns
- Renaming columns
- Changing column types
- Complex data transformations
- Dropping tables

## Why It's Safe for Production

1. **Non-Destructive**
   - Never deletes data
   - Only adds structure
   - Existing users keep working

2. **Automatic**
   - Runs on every server start
   - No manual SQL needed
   - Catches up automatically

3. **Idempotent**
   - Can run multiple times safely
   - Only makes changes if needed
   - Won't break if run twice

## What You'll See in Logs

### Successful Migration:
```
Running auto-migration...
✅ Auto-migration completed successfully
✅ All tables synced including: users (with email_verified), workout_plans, scheduled_workouts, offline_nutrition tables
Database connection established
```

### What GORM Does Behind the Scenes:
```sql
-- Check if users table exists
SELECT * FROM information_schema.tables WHERE table_name = 'users';

-- Check existing columns
SELECT column_name FROM information_schema.columns WHERE table_name = 'users';

-- Add missing columns
ALTER TABLE users ADD COLUMN email_verified BOOLEAN DEFAULT false;
ALTER TABLE users ADD COLUMN verification_token VARCHAR(255);
ALTER TABLE users ADD COLUMN verification_token_expires_at TIMESTAMP;

-- Create new tables
CREATE TABLE workout_plans (...);
CREATE TABLE scheduled_workouts (...);
CREATE TABLE dish_masters (...);

-- Create indexes
CREATE INDEX idx_users_verification_token ON users(verification_token);
```

## Common Questions

### Q: Will my existing users lose their data?
**A:** No! All existing user data is preserved. They just get new columns with default values.

### Q: What if I run AutoMigrate twice?
**A:** It's safe! GORM checks what exists and only makes necessary changes.

### Q: Can I roll back if something goes wrong?
**A:** AutoMigrate doesn't delete anything, so rolling back just means:
- Old code ignores new columns
- New columns remain empty

### Q: Do I need to stop the server to run migration?
**A:** No! Migration runs on server startup. Zero downtime.

### Q: What if migration fails?
**A:** Server won't start. Error is logged. Fix the issue and restart.

## Summary

**AutoMigrate = "Make my database tables match my Go code"**

It's like having a smart assistant that:
- ✅ Creates missing tables
- ✅ Adds missing columns
- ✅ Creates missing indexes
- ❌ Never deletes anything
- ❌ Never breaks existing data

**Perfect for:**
- Adding new features
- Adding new fields to existing models
- Development and production
- Continuous deployment

**Not for:**
- Complex schema changes
- Renaming columns
- Changing data types
- Data migrations

For those, you need manual SQL migrations.
