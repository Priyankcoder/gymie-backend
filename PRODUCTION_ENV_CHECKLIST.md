
# Production Environment Variables Checklist

## ⚠️ CRITICAL ISSUES FOUND

Your production deployment is **MISSING** several required environment variables!

## 🔴 Missing from Render Configuration

### 1. **Email/SMTP Configuration** (CRITICAL for OTP)
```yaml
# Currently MISSING from render.yaml:
- SMTP_HOST
- SMTP_PORT
- SMTP_USERNAME
- SMTP_PASSWORD
- FROM_EMAIL
- FROM_NAME
```

**Impact:** Email verification won't work! Users can't verify their accounts.

### 2. **Frontend URL** (CRITICAL for email links)
```yaml
# Currently MISSING from render.yaml:
- FRONTEND_URL
```

**Impact:** Email verification links will be broken (default to localhost).

### 3. **Storage Configuration** (for progress photos)
```yaml
# Currently MISSING from render.yaml:
- STORAGE_ENDPOINT
- STORAGE_ACCESS_KEY
- STORAGE_SECRET_KEY
- STORAGE_BUCKET
- STORAGE_REGION
- STORAGE_PUBLIC_URL
- STORAGE_PRESIGN_EXPIRY
```

**Impact:** Progress photos won't upload.

### 4. **Frontend API URL**
```yaml
# Frontend needs EXPO_PUBLIC_API_URL
```

**Impact:** App will try to connect to localhost instead of production API.

---

## ✅ Currently Configured (Good)

- ✅ `PORT` - 8080
- ✅ `ENVIRONMENT` - production
- ✅ `DATABASE_URL` - Auto from Render DB
- ✅ `JWT_SECRET` - Auto-generated
- ✅ `JWT_EXPIRATION` - 24 hours
- ✅ `CORS_ALLOWED_ORIGINS` - *
- ✅ `RATE_LIMIT_*` - All set
- ⚠️ `REDIS_URL` - Needs manual setup from Upstash

---

## 🔧 HOW TO FIX

### Step 1: Update Render Environment Variables

Go to your Render service dashboard and add these:

#### **Email Configuration (Choose One)**

**Option A: Gmail (Quick Setup)**
```
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-gmail@gmail.com
SMTP_PASSWORD=your-app-password (not regular password!)
FROM_EMAIL=your-gmail@gmail.com
FROM_NAME=Gymie
```

**How to get Gmail App Password:**
1. Go to Google Account settings
2. Security → 2-Step Verification → App passwords
3. Generate app password
4. Use that password (not your regular Gmail password)

**Option B: SendGrid (Recommended for Production)**
```
SMTP_HOST=smtp.sendgrid.net
SMTP_PORT=587
SMTP_USERNAME=apikey
SMTP_PASSWORD=your-sendgrid-api-key
FROM_EMAIL=noreply@gymie.fit
FROM_NAME=Gymie
```

**How to get SendGrid:**
1. Sign up at sendgrid.com (free tier: 100 emails/day)
2. Create API Key
3. Use "apikey" as username (literal text)
4. Use API key as password

#### **Frontend URL**
```
FRONTEND_URL=https://gymie.fit
```

Or whatever your actual frontend domain is. This is used for:
- Email verification links
- Password reset links (if you add that feature)

#### **Storage Configuration (Cloudflare R2)**

If you're using Cloudflare R2 for progress photos:
```
STORAGE_ENDPOINT=https://YOUR-ACCOUNT-ID.r2.cloudflarestorage.com
STORAGE_ACCESS_KEY=your-r2-access-key
STORAGE_SECRET_KEY=your-r2-secret-key
STORAGE_BUCKET=gymie-production
STORAGE_REGION=auto
STORAGE_PUBLIC_URL=https://YOUR-BUCKET.YOUR-DOMAIN.com
STORAGE_PRESIGN_EXPIRY=15
```

**If not using Cloudflare R2 yet:**
- Can skip for now
- Progress photos will use local storage (works but not ideal)
- Set up R2 later when needed

#### **Redis URL (Upstash)**
```
REDIS_URL=redis://default:YOUR-PASSWORD@YOUR-HOST.upstash.io:PORT
```

Get this from your Upstash dashboard.

---

### Step 2: Update Frontend API URL

In `frontend/.env` or app.json:
```
EXPO_PUBLIC_API_URL=https://your-render-app.onrender.com/v1
```

Replace `your-render-app` with your actual Render service URL.

---

### Step 3: Updated render.yaml

I'll create an updated render.yaml with all required variables...

---

## 🚨 DEPLOYMENT BLOCKER

**You CANNOT deploy Android app until these are fixed:**

1. ❌ Email verification won't work (users stuck at registration)
2. ❌ Verification links will point to localhost
3. ❌ Frontend can't connect to backend API

**Fix Order:**
1. Add SMTP configuration to Render
2. Add FRONTEND_URL to Render
3. Update REDIS_URL in Render
4. Update EXPO_PUBLIC_API_URL in frontend
5. Redeploy backend (AutoMigrate will run)
6. Build Android app

---

## 📋 Complete Environment Variables List

### Backend (Render)

**Required (MUST have):**
- [ ] `PORT` - ✅ Set (8080)
- [ ] `ENVIRONMENT` - ✅ Set (production)
- [ ] `DATABASE_URL` - ✅ Auto-configured
- [ ] `REDIS_URL` - ⚠️ Need to add from Upstash
- [ ] `JWT_SECRET` - ✅ Auto-generated
- [ ] `JWT_EXPIRATION` - ✅ Set (24)
- [ ] `SMTP_HOST` - ❌ MISSING
- [ ] `SMTP_PORT` - ❌ MISSING
- [ ] `SMTP_USERNAME` - ❌ MISSING
- [ ] `SMTP_PASSWORD` - ❌ MISSING
- [ ] `FROM_EMAIL` - ❌ MISSING
- [ ] `FROM_NAME` - ❌ MISSING
- [ ] `FRONTEND_URL` - ❌ MISSING

**Optional (Nice to have):**
- [ ] `CORS_ALLOWED_ORIGINS` - ✅ Set (*)
- [ ] `RATE_LIMIT_ENABLED` - ✅ Set
- [ ] `RATE_LIMIT_REQUESTS` - ✅ Set
- [ ] `RATE_LIMIT_WINDOW` - ✅ Set
- [ ] `STORAGE_ENDPOINT` - ⚠️ For progress photos
- [ ] `STORAGE_ACCESS_KEY` - ⚠️ For progress photos
- [ ] `STORAGE_SECRET_KEY` - ⚠️ For progress photos
- [ ] `STORAGE_BUCKET` - ⚠️ For progress photos
- [ ] `STORAGE_REGION` - ⚠️ For progress photos
- [ ] `STORAGE_PUBLIC_URL` - ⚠️ For progress photos

### Frontend (Expo)

**Required:**
- [ ] `EXPO_PUBLIC_API_URL` - ❌ MISSING (still using localhost)

---

## 🧪 How to Test

### Test Email Configuration
```bash
# After deploying with SMTP vars
curl -X POST https://your-app.onrender.com/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "testpass123",
    "name": "Test User"
  }'

# Check Render logs for:
"=== EMAIL DEBUG ==="
"Sending verification email to: test@example.com"
```

### Test Frontend Connection
```bash
# In your app
console.log(API_CONFIG.BASE_URL)
# Should show: https://your-app.onrender.com/v1
# NOT: http://localhost:8080/v1
```

---

## 📝 Quick Setup Commands

### Add to Render via CLI (if using Render CLI)
```bash
render env set SMTP_HOST smtp.gmail.com
render env set SMTP_PORT 587
render env set SMTP_USERNAME your-email@gmail.com
render env set SMTP_PASSWORD your-app-password
render env set FROM_EMAIL your-email@gmail.com
render env set FROM_NAME Gymie
render env set FRONTEND_URL https://gymie.fit
```

### Or use Render Dashboard:
1. Go to your Render service
2. Environment → Add Environment Variable
3. Add each variable one by one
4. Click "Save Changes"
5. Render will auto-redeploy

---

## ⚡ Priority Order

**Do these IN ORDER:**

1. **HIGH PRIORITY (Deploy Blocker):**
   - [ ] Add SMTP configuration (emails won't work)
   - [ ] Add FRONTEND_URL (verification links broken)
   - [ ] Add REDIS_URL (sessions won't work)
   - [ ] Update frontend API URL (app can't connect)

2. **MEDIUM PRIORITY:**
   - [ ] Add storage configuration (for progress photos)

3. **LOW PRIORITY:**
   - [ ] Fine-tune rate limiting
   - [ ] Update CORS if needed

---

## 🎯 Summary

**Status:** ❌ NOT READY FOR PRODUCTION

**Missing:** 7 critical environment variables

**Action Required:** Add environment variables to Render before deploying Android app

**Time to Fix:** ~15 minutes

**Next Step:** I'll create an updated render.yaml file with all variables marked...
