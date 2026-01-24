
# Database Migration Checklist

## ✅ Migration Status: READY TO DEPLOY

All migration files are correctly configured and in sync.

## Models Included in Migration (17 Total)

### ✅ Core User Models (2)
- [x] `User` - User accounts with email verification fields
  - EmailVerified (bool)
  - VerificationToken (string)
  - VerificationTokenExpiresAt (time.Time)
- [x] `UserProfile` - User profile details

### ✅ Workout Models (3)
- [x] `Workout` - Workout session records
- [x] `Exercise` - Exercises within a workout
- [x] `WorkoutSet` - Sets within exercises

### ✅ Nutrition Models (3)
- [x] `NutritionDay` - Daily nutrition logs
- [x] `Meal` - Meals within a day
- [x] `Food` - Food items within meals

### ✅ Progress Tracking Models (2)
- [x] `ProgressPhoto` - Progress photos
- [x] `WeightEntry` - Weight tracking

### ✅ Workout Planning Models (3)
- [x] `WorkoutPlan` - Workout plan templates
- [x] `WorkoutPlanDay` - Days within workout plans
- [x] `ScheduledWorkout` - Scheduled workout instances

### ✅ Offline Nutrition Models (4)
- [x] `DishMaster` - Offline dish database
- [x] `DishNutritionMaster` - Nutrition info for dishes
- [x] `UserCorrection` - User corrections to nutrition data
- [x] `ModelVersion` - ML model version tracking

## Files Verified

### ✅ Auto-Migration (Runs on API Startup)
**File:** `backend/internal/config/database.go`
- Function: `autoMigrate(db *gorm.DB)`
- Called by: `InitDB()` automatically on server start
- **Status:** ✅ Has all 17 models
- **Location:** Lines 70-113

### ✅ Manual Migration Script (Optional)
**File:** `backend/migrations/migrate.go`
- Command: `go run cmd/migrate/main.go`
- **Status:** ✅ Has all 17 models
- **Location:** Lines 26-45

### ✅ Both Files Match
- migrate.go: 17 models ✅
- database.go: 17 models ✅
- **Conclusion:** Perfect sync! 🎉

## Email Verification Fields on User Model

The `User` model (`internal/models/user.go`) includes:

```go
// Email verification
EmailVerified              bool      `gorm:"default:false" json:"emailVerified"`
VerificationToken          string    `gorm:"index:idx_users_verification_token;unique" json:"-"`
VerificationTokenExpiresAt time.Time `json:"-"`
```

**Status:** ✅ Correctly defined

## How Migration Runs

### Automatic (Recommended) ✅
Migration runs automatically when the API server starts:

```
API Server Starts
    ↓
config.Load() - Load environment
    ↓
config.InitDB() - Connect to database
    ↓
autoMigrate() - Run migration automatically
    ↓
Server Ready
```

**You don't need to do anything!** Just deploy and it will migrate automatically.

### Manual (If Needed)
If you need to run migration separately:

```bash
cd backend
go run migrations/migrate.go
```

## Deployment Steps

1. **Commit Changes**
   ```bash
   git add backend/internal/config/database.go
   git commit -m "Add all models to auto-migration"
   git push origin main
   ```

2. **Render Auto-Deploy**
   - Render detects the push
   - Builds and deploys automatically
   - Migration runs on startup

3. **Verify in Logs**
   Look for these lines in Render logs:
   ```
   Running auto-migration...
   ✅ Auto-migration completed successfully
   ✅ All tables synced including: users (with email_verified), workout_plans, scheduled_workouts, offline_nutrition tables
   ```

4. **Test Endpoints**
   ```bash
   # Test verification status
   curl https://gymie-api.onrender.com/v1/auth/verification-status/test@test.com
   
   # Should return JSON, not 404
   ```

## What Gets Created/Updated

### New Tables (if they don't exist)
- workout_plans
- workout_plan_days
- scheduled_workouts
- dish_masters
- dish_nutrition_masters
- user_corrections
- model_versions

### Updated Tables (if they exist)
- **users** table gets new columns:
  - `email_verified` (boolean, default false)
  - `verification_token` (string, indexed)
  - `verification_token_expires_at` (timestamp)

### Existing Data
- **Safe:** GORM AutoMigrate only ADDS columns, never deletes
- Existing users will have `email_verified = false` by default
- Existing data is preserved

## Verification Commands

### 1. Check Migration Logs (Render Dashboard)
```
Logs Tab → Look for:
"Running auto-migration..."
"✅ Auto-migration completed successfully"
```

### 2. Test Verification Endpoint
```bash
curl https://gymie-api.onrender.com/v1/auth/verification-status/test@test.com
```

**Expected:** JSON response (even if user doesn't exist)
**Bad:** 404 or "method not allowed"

### 3. Register New User
```bash
curl -X POST https://gymie-api.onrender.com/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "username": "testuser",
    "password": "Test123!@#"
  }'
```

**Expected:** Response includes `"emailVerified": false`

### 4. Check User Profile
```bash
curl https://gymie-api.onrender.com/v1/auth/me \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Expected:** User object includes `emailVerified` field

## Common Issues & Solutions

### Issue: "Column doesn't exist" Error
**Solution:** Migration hasn't run yet. Redeploy.

### Issue: No Migration Logs
**Solution:** Check if latest code is deployed. Force manual deploy.

### Issue: 404 on Verification Endpoint
**Solution:** Routes not updated. Check if routes include auth verification endpoints.

### Issue: "main redeclared" Build Error
**Solution:** Already fixed - scripts directory excluded from build.

## Next Steps After Migration

1. ✅ Verify migration ran (check logs)
2. ✅ Test all email verification endpoints
3. ✅ Add SMTP credentials to Render
4. ✅ Test email sending
5. ✅ Build production Android APK
6. ✅ Test with real devices

## Summary

**Status:** ✅ All migration files are correct and in sync
**Action Required:** Just commit, push, and deploy
**Migration Type:** Automatic on server startup
**Safety:** GORM only adds columns, never deletes data
**Total Models:** 17 (all present in both files)

---

**Ready to deploy!** 🚀

The migration will run automatically when Render deploys your backend. No manual intervention needed.
