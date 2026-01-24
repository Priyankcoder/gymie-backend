
# Checking Database Migration Status

## Quick Check: Is Your Database Updated?

### 1. Check User Schema (Email Verification)

Test if the `email_verified` and related columns exist:

```bash
# Get your user profile (requires valid JWT token)
curl -X GET https://gymie-api.onrender.com/v1/auth/me \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json"
```

**Expected Response (if migrated):**
```json
{
  "id": 1,
  "email": "user@example.com",
  "username": "user123",
  "email_verified": false,  // ✅ This field should exist
  "created_at": "2024-01-01T00:00:00Z"
}
```

**If NOT migrated:**
- Field `email_verified` will be missing
- Or endpoint returns 500 error about missing column

### 2. Check Email Verification Endpoints

```bash
# Test verification status endpoint
curl -X GET https://gymie-api.onrender.com/v1/auth/verification-status/test@example.com \
  -H "Content-Type: application/json"
```

**Expected Response (if deployed):**
```json
{
  "email": "test@example.com",
  "verified": false,
  "user_exists": false
}
```

**If NOT deployed:**
- Returns 404 Not Found
- Or "method not allowed"

### 3. Check Health/Status Endpoint

```bash
# Check if API is responding
curl https://gymie-api.onrender.com/v1/health
```

### 4. Check Server Logs on Render

1. Go to your Render dashboard
2. Click on your backend service
3. Go to **Logs** tab
4. Look for these lines:

**If migration succeeded:**
```
Running auto-migration...
✅ Auto-migration completed successfully
✅ All tables synced including: users (with email_verified), workout_plans, scheduled_workouts, offline_nutrition tables
Database connection established
Starting server on port 8080
```

**If migration NOT run:**
- Won't see the auto-migration log lines
- May see errors about missing columns when trying to use endpoints

## How to Force Migration

### Option 1: Redeploy on Render (Recommended)

After you've pushed the updated `database.go` file:

1. Go to Render Dashboard → Your Backend Service
2. Click **Manual Deploy** → **Deploy latest commit**
3. Wait for deployment (2-3 minutes)
4. Check logs for migration messages

### Option 2: Run Migration Command Manually

If you have SSH access or can run commands:

```bash
# From your backend directory
go run cmd/migrate/main.go
```

This runs the standalone migration script that uses the same models.

### Option 3: Trigger Render Build via Git

```bash
# Make any small change (like adding a comment)
echo "# Migration update" >> backend/README.md
git add .
git commit -m "Trigger migration"
git push origin main
```

Render will auto-deploy and run the migration.

## Verify Migration Worked

### Step 1: Check Logs

Look for these exact lines in Render logs:
```
Running auto-migration...
✅ Auto-migration completed successfully
✅ All tables synced including: users (with email_verified), workout_plans, scheduled_workouts, offline_nutrition tables
```

### Step 2: Test Registration with Email Verification

```bash
# Register a new user
curl -X POST https://gymie-api.onrender.com/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newuser@example.com",
    "username": "newuser",
    "password": "Test123!@#"
  }'
```

**Expected Response (if working):**
```json
{
  "user": {
    "id": 123,
    "email": "newuser@example.com",
    "username": "newuser",
    "email_verified": false
  },
  "token": "eyJhbGc...",
  "message": "Registration successful. Please check your email to verify your account."
}
```

### Step 3: Check Verification Status

```bash
curl -X GET https://gymie-api.onrender.com/v1/auth/verification-status/newuser@example.com
```

**Expected:**
```json
{
  "email": "newuser@example.com",
  "verified": false,
  "user_exists": true
}
```

### Step 4: Test Email Resend (Optional)

```bash
curl -X POST https://gymie-api.onrender.com/v1/auth/resend-verification \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newuser@example.com"
  }'
```

## Common Issues

### Issue 1: Migration Doesn't Run on Deploy

**Cause:** The `InitDB()` function calls `autoMigrate()` which should run automatically on every API server start.

**Fix:**
1. Check if `database.go` has the updated `autoMigrate()` function
2. Verify the file is committed and pushed
3. Force a new deployment on Render

### Issue 2: "Column doesn't exist" Errors

**Symptoms:**
```
ERROR: column "email_verified" does not exist
```

**This means:** Migration hasn't run yet.

**Fix:**
1. Redeploy the service
2. Check logs to confirm migration ran
3. If migration logs appear but error persists, your code might be using old field names

### Issue 3: SMTP Errors in Logs

**Symptoms:**
```
Failed to send email: dial tcp: lookup smtp.gmail.com: no such host
```

**This is OK** - it means the migration worked and the code is trying to send emails, but SMTP isn't configured yet.

**Fix:** Add SMTP environment variables to Render (see PRODUCTION_ENV_CHECKLIST.md)

### Issue 4: No Migration Logs Appear

**Possible Causes:**
1. The server started but migration was skipped
2. Database connection failed before migration
3. Old code is still deployed

**Fix:**
1. Check Render deployment history - is latest commit deployed?
2. Check database credentials in Render environment variables
3. Force a new deployment

## Database Direct Check (Advanced)

If you have database access, you can check directly:

```sql
-- Check if email_verified column exists
SELECT column_name, data_type 
FROM information_schema.columns 
WHERE table_name = 'users';

-- Should see:
-- email_verified | boolean
-- email_verified_at | timestamp
-- verification_token | varchar
-- verification_token_expires | timestamp

-- Check all tables exist
SELECT table_name 
FROM information_schema.tables 
WHERE table_schema = 'public'
ORDER BY table_name;

-- Should see tables like:
-- users
-- user_profiles  
-- workouts
-- exercises
-- nutrition_days
-- progress_photos
-- workout_plans
-- workout_plan_days
-- scheduled_workouts
-- dish_masters
-- dish_nutrition_masters
-- user_corrections
-- model_versions
```

## What Tables Should Exist After Migration

After successful migration, you should have these tables:

**Core Tables:**
- `users` - User accounts (with email_verified column)
- `user_profiles` - User profile details

**Workout Tables:**
- `workouts` - Workout records
- `exercises` - Exercise details
- `workout_sets` - Sets within exercises
- `workout_plans` - Workout plan templates
- `workout_plan_days` - Days within plans
- `scheduled_workouts` - Scheduled workout instances

**Nutrition Tables:**
- `nutrition_days` - Daily nutrition logs
- `meals` - Meals within a day
- `foods` - Foods within a meal

**Progress Tables:**
- `progress_photos` - Progress photo uploads
- `weight_entries` - Weight tracking entries

**Offline Nutrition Tables:**
- `dish_masters` - Dish database
- `dish_nutrition_masters` - Nutrition info for dishes
- `user_corrections` - User corrections to nutrition data
- `model_versions` - ML model version tracking

## Next Steps After Confirming Migration

1. ✅ Verify all endpoints work
2. ✅ Add SMTP credentials to Render (for email sending)
3. ✅ Test email verification flow end-to-end
4. ✅ Build and deploy Android APK
5. ✅ Test production app with real devices

---

**Quick Summary Command:**

```bash
# One command to check everything
curl -s https://gymie-api.onrender.com/v1/auth/verification-status/test@test.com | python -m json.tool
```

If you get a valid JSON response (even with `user_exists: false`), your migration is working! 🎉
