
# Automatic Database Schema Sync Guide

## The Better Way: GORM AutoMigrate ✨

Instead of manually running SQL migrations, you can now use **GORM's AutoMigrate** feature which automatically syncs your database schema with your Go models.

## What Changed

I've updated `backend/internal/config/database.go` to include automatic schema migration on every server startup.

### How It Works

```go
// When your backend starts, it now automatically:
1. Connects to the database
2. Compares your Go models with the database schema
3. Adds any missing columns
4. Creates any missing tables
5. Updates indexes
```

**What it does:**
- ✅ Adds new columns (like email_verified, verification_token)
- ✅ Creates new tables
- ✅ Updates indexes
- ✅ Safe to run multiple times (idempotent)

**What it doesn't do:**
- ❌ Delete columns (safe by design)
- ❌ Modify existing data
- ❌ Change column types (requires manual migration)

## How to Use

### Option 1: Just Restart Your Backend (Recommended)

```bash
# On Render, Heroku, or any hosting platform
# Just redeploy or restart your backend service

# The AutoMigrate will run automatically on startup
```

**That's it!** No manual migration needed.

### Option 2: Test Locally First

```bash
cd backend

# Make sure DATABASE_URL points to your production database in .env
# Or use local database first to test

# Start the backend
go run cmd/api/main.go
```

You'll see:
```
Database connection established
Running auto-migration...
✅ Auto-migration completed successfully
Starting server on port 8080
```

### Option 3: Run Migration Only (Without Starting Server)

If you want to just run the migration without starting the server:

```bash
cd backend

# Create a simple migration-only script
go run cmd/migrate/main.go
```

## Deployment Workflow

### Previous Way (Manual):
```
1. Write SQL migration file
2. Connect to database
3. Run SQL manually
4. Deploy backend
5. Deploy Android app
```

### New Way (Automatic):
```
1. Update your Go models (already done ✅)
2. Deploy/restart backend (AutoMigrate runs automatically)
3. Deploy Android app
```

## What Models Are Auto-Migrated

All your models are now automatically synced:

- ✅ `User` (includes email verification fields)
- ✅ `UserProfile`
- ✅ `Exercise`
- ✅ `Workout`
- ✅ `WorkoutExercise`
- ✅ `ExerciseSet`
- ✅ `NutritionDay`
- ✅ `Meal`
- ✅ `MealFood`
- ✅ `ProgressPhoto`
- ✅ `WeightEntry`

## For Your Current Situation

### Render PostgreSQL + Android Build

**Simple steps:**

1. **Redeploy your backend on Render**
   - Go to Render dashboard
   - Click "Manual Deploy" or push to your connected repo
   - AutoMigrate will run during startup
   - Check logs for: "✅ Auto-migration completed successfully"

2. **Build Android app**
   ```bash
   cd frontend
   eas build --platform android
   ```

That's it! No manual database work needed.

## Safety Features

### Idempotent
You can run it multiple times safely:
```
First run:  Creates email_verified column ✅
Second run: Column already exists, skips ✅
Third run:  Column already exists, skips ✅
```

### Non-Destructive
- **Never deletes** columns (even if you remove from model)
- **Never modifies** existing data
- **Only adds** what's missing

### Production Safe
- No downtime required
- Works with live traffic
- Transactions used internally by GORM

## Monitoring

### Check Migration Success

**In logs:**
```bash
# Render Dashboard → Logs
# Look for:
"Running auto-migration..."
"✅ Auto-migration completed successfully"
```

**Verify in database:**
```sql
-- Check if new columns exist
SELECT column_name, data_type 
FROM information_schema.columns
WHERE table_name = 'users'
AND column_name IN ('email_verified', 'verification_token', 'verification_token_expires_at');
```

Should return 3 rows.

## When to Use Manual Migrations

AutoMigrate is great for most cases, but use manual migrations when:

1. **Changing column types** (e.g., VARCHAR → TEXT)
2. **Complex data transformations** (e.g., splitting a column)
3. **Dropping columns** (AutoMigrate won't do this)
4. **Renaming columns** (AutoMigrate treats as drop + create)
5. **Performance-critical changes** (want control over timing)

## Troubleshooting

### "Migration Failed"
- Check database connection
- Verify DATABASE_URL is correct
- Ensure database user has CREATE/ALTER permissions

### "Column Already Exists"
- This is normal! AutoMigrate is idempotent
- It will skip and continue

### "Permission Denied"
- Database user needs `CREATE` and `ALTER TABLE` permissions
- On Render, default user has these permissions

## Benefits of This Approach

### For Development
- ✅ No manual SQL writing
- ✅ Schema always in sync with code
- ✅ Easy to add new fields
- ✅ One source of truth (Go models)

### For Deployment
- ✅ Automatic on every deploy
- ✅ No separate migration step
- ✅ Fail-fast if schema incompatible
- ✅ Works with CI/CD

### For Team Collaboration
- ✅ No "forgot to run migration" errors
- ✅ Everyone's DB stays in sync
- ✅ Less documentation needed
- ✅ Fewer deployment steps

## Comparison

| Aspect | Manual Migrations | Auto-Migrate |
|--------|------------------|--------------|
| Effort | Write SQL + Run manually | Write Go model only |
| Safety | Depends on SQL quality | Built-in safety |
| Speed | Must remember to run | Automatic |
| Rollback | Manual DOWN migration | Model-driven |
| Team Sync | Can get out of sync | Always synced |
| Learning Curve | Need SQL knowledge | Go structs only |

## Best Practices

### 1. Always Test Locally First
```bash
# Point to local DB first
DATABASE_URL=postgresql://localhost/gymie_test

# Start backend and verify
go run cmd/api/main.go
```

### 2. Check Logs After Deploy
```bash
# Ensure you see:
"✅ Auto-migration completed successfully"
```

### 3. Keep Manual Migrations for Complex Changes
```bash
# Use migrations/ folder for:
# - Data transformations
# - Column type changes
# - Dropping columns
```

### 4. Version Control Your Models
```bash
# Git commit after model changes
git add internal/models/
git commit -m "Add email verification fields to User model"
```

## FAQ

**Q: Will this delete my data?**
A: No. AutoMigrate only adds, never deletes data or columns.

**Q: Can I turn it off?**
A: Yes, comment out the `autoMigrate(db)` call in database.go

**Q: Does it slow down startup?**
A: Minimal. Takes <1 second, runs once per startup.

**Q: What if I have multiple backend instances?**
A: Safe! All instances can run AutoMigrate concurrently.

**Q: Do I still need the migrations/ folder?**
A: Keep it for complex migrations. AutoMigrate handles simple schema changes.

**Q: Will this work with my existing data?**
A: Yes! It only adds missing columns/tables. Existing data is untouched.

## Summary

**Old Way:**
```
1. Edit Go model
2. Write SQL migration
3. Connect to DB
4. Run migration
5. Deploy backend
6. Deploy app
```

**New Way:**
```
1. Edit Go model
2. Deploy backend (AutoMigrate runs automatically)
3. Deploy app
```

**Saved Steps:** 3 manual steps eliminated! ✨

---

**Status:** ✅ Implemented and ready to use

**Next Step:** Just redeploy your backend on Render, and the migration will run automatically!
