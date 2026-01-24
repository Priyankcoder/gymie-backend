
# Why Email Works Locally But Not on Render

## TL;DR

**The problem:** Render doesn't automatically read your `.env.production` file!

You need to **manually set environment variables** in Render's dashboard.

---

## How Environment Variables Work

### Local Development ✅

When you run locally:
```bash
go run cmd/api/main.go
```

The code reads from:
1. `backend/.env` file (if exists)
2. Your system environment variables

**Your local `.env` has:**
```env
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=priyankrastogi14@gmail.com
SMTP_PASSWORD=rtcakrynvmgvdpos  # Your real App Password
SMTP_FROM=Gymie <priyankrastogi14@gmail.com>
```

**Result:** Email works! ✅

### Render Production ❌

When Render deploys your code:

**What Render DOES:**
1. ✅ Clones your git repository
2. ✅ Builds your Go code
3. ✅ Runs your API server

**What Render DOES NOT:**
1. ❌ Read `.env` files (they're in .gitignore)
2. ❌ Read `.env.production` files
3. ❌ Automatically set environment variables

**Where Render gets environment variables:**
- **Only from Render Dashboard → Environment tab**
- You must set them manually!

### Why `.env.production` Doesn't Work

```
Your Repo:
├── backend/
│   ├── .env                  ← In .gitignore, not pushed
│   ├── .env.production       ← Pushed, but Render IGNORES it
│   └── cmd/api/main.go       ← Render runs this
```

**The `.env.production` file is:**
- ✅ In your git repo (for reference)
- ✅ A template for what variables are needed
- ❌ NOT read by Render automatically
- ❌ NOT used at runtime

**To use it, you must:**
1. Open Render Dashboard
2. Go to Environment tab
3. Manually copy each variable from `.env.production`
4. Paste into Render's UI

---

## Check Your Render Environment Variables

### Step 1: Go to Render Dashboard

1. Visit: https://dashboard.render.com
2. Click on your backend service (e.g., "gymie-backend")
3. Click **"Environment"** in the left sidebar

### Step 2: Check What Variables Exist

Look for these variables:

**Check if these exist:**
- [ ] `SMTP_HOST`
- [ ] `SMTP_PORT`
- [ ] `SMTP_USERNAME`
- [ ] `SMTP_PASSWORD`
- [ ] `SMTP_FROM`
- [ ] `FRONTEND_URL`

**Check if OLD variables exist (wrong names):**
- [ ] `FROM_EMAIL` ← Delete this if exists
- [ ] `FROM_NAME` ← Delete this if exists

### Step 3: Compare Values

**Check SMTP_PASSWORD:**
```
Is it: rtcakrynvmgvdpos  ✅ Correct (no spaces)
Or: rtca kryn vmgv dpos  ❌ Wrong (has spaces)
```

**Check SMTP_FROM:**
```
Is it: Gymie <priyankrastogi14@gmail.com>  ✅ Correct
Or: priyankrastogi14@gmail.com            ❌ Wrong (missing name)
```

---

## How to Fix

### Option 1: Set Variables Manually in Render

1. **Go to Environment tab**
2. **Click "Add Environment Variable"**
3. **Add each one:**

```
Key: SMTP_HOST
Value: smtp.gmail.com

Key: SMTP_PORT
Value: 587

Key: SMTP_USERNAME
Value: priyankrastogi14@gmail.com

Key: SMTP_PASSWORD
Value: rtcakrynvmgvdpos

Key: SMTP_FROM
Value: Gymie <priyankrastogi14@gmail.com>

Key: FRONTEND_URL
Value: https://gymie.fit
```

4. **Delete old variables** (if they exist):
   - Delete `FROM_EMAIL`
   - Delete `FROM_NAME`

5. **Click "Save Changes"**

Render will automatically redeploy with new variables.

### Option 2: Use render.yaml (Automated)

Create `render.yaml` in your repo root:

```yaml
services:
  - type: web
    name: gymie-backend
    env: go
    buildCommand: go build -o bin/api cmd/api/main.go
    startCommand: ./bin/api
    envVars:
      - key: SMTP_HOST
        value: smtp.gmail.com
      - key: SMTP_PORT
        value: 587
      - key: SMTP_USERNAME
        value: priyankrastogi14@gmail.com
      - key: SMTP_PASSWORD
        sync: false  # Will prompt you to set manually (for security)
      - key: SMTP_FROM
        value: Gymie <priyankrastogi14@gmail.com>
      - key: FRONTEND_URL
        value: https://gymie.fit
      - key: DATABASE_URL
        fromDatabase:
          name: gymie-db
          property: connectionString
```

**But you still need to set `SMTP_PASSWORD` manually in dashboard for security!**

---

## How to Verify What Render Is Using

### Add Debug Logging

Temporarily add this to your `cmd/api/main.go`:

```go
func main() {
    log.Println("=== ENVIRONMENT VARIABLES ===")
    log.Println("SMTP_HOST:", os.Getenv("SMTP_HOST"))
    log.Println("SMTP_PORT:", os.Getenv("SMTP_PORT"))
    log.Println("SMTP_USERNAME:", os.Getenv("SMTP_USERNAME"))
    
    smtpPass := os.Getenv("SMTP_PASSWORD")
    if smtpPass == "" {
        log.Println("SMTP_PASSWORD: (empty!)")
    } else {
        log.Printf("SMTP_PASSWORD: %s****%s (%d chars)\n", 
            smtpPass[:2], smtpPass[len(smtpPass)-2:], len(smtpPass))
    }
    
    log.Println("SMTP_FROM:", os.Getenv("SMTP_FROM"))
    log.Println("FROM_EMAIL:", os.Getenv("FROM_EMAIL"))
    log.Println("FROM_NAME:", os.Getenv("FROM_NAME"))
    log.Println("=============================")
    
    // ... rest of your code
}
```

### Check Render Logs

After deploying, go to Render → Logs:

**If variables are set correctly:**
```
=== ENVIRONMENT VARIABLES ===
SMTP_HOST: smtp.gmail.com
SMTP_PORT: 587
SMTP_USERNAME: priyankrastogi14@gmail.com
SMTP_PASSWORD: rt****os (16 chars)
SMTP_FROM: Gymie <priyankrastogi14@gmail.com>
FROM_EMAIL: (empty)
FROM_NAME: (empty)
=============================
```

**If variables are NOT set:**
```
=== ENVIRONMENT VARIABLES ===
SMTP_HOST: smtp.gmail.com
SMTP_PORT: 587
SMTP_USERNAME: priyankrastogi14@gmail.com
SMTP_PASSWORD: (empty!)  ← Problem!
SMTP_FROM: 
FROM_EMAIL: priyankrastogi14@gmail.com  ← Using old name!
FROM_NAME: Gymie  ← Using old name!
=============================
```

---

## Why Your Code Changes ARE Needed

Even though the issue is Render environment variables, the code changes are needed because:

### 1. Support New Variable Format

**Old code:**
```go
fromEmail: getEnv("FROM_EMAIL", "...")  // Only reads FROM_EMAIL
fromName:  getEnv("FROM_NAME", "...")   // Only reads FROM_NAME
```

**New code:**
```go
// Reads SMTP_FROM first, falls back to FROM_EMAIL/FROM_NAME
smtpFrom := getEnv("SMTP_FROM", "")
```

### 2. Parse Combined Format

**Old format (2 variables):**
```env
FROM_EMAIL=priyankrastogi14@gmail.com
FROM_NAME=Gymie
```

**New format (1 variable):**
```env
SMTP_FROM=Gymie <priyankrastogi14@gmail.com>
```

Code needed to parse the new format!

### 3. Config Struct Was Missing SMTP Fields

**Before:**
```go
type Config struct {
    Port string
    DBHost string
    // No SMTP fields!
}
```

**After:**
```go
type Config struct {
    Port string
    DBHost string
    SMTPHost string  // Added
    SMTPPort string  // Added
    SMTPFrom string  // Added
    // ...
}
```

Without these fields, the config package couldn't even store SMTP settings!

---

## Summary

### Why It Works Locally:
✅ Reads from `.env` file
✅ Has correct SMTP credentials
✅ Your network allows SMTP connections

### Why It Doesn't Work on Render:
❌ Render doesn't read `.env` or `.env.production` files
❌ Environment variables not set in Render Dashboard
❌ Or they're set with old names (`FROM_EMAIL`/`FROM_NAME`)
❌ Or they have spaces in the password

### What You Need to Do:
1. ✅ Push code changes (already done)
2. ✅ Set environment variables in Render Dashboard
3. ✅ Use correct variable names (SMTP_FROM, not FROM_EMAIL)
4. ✅ Remove spaces from SMTP_PASSWORD
5. ✅ Wait for Render to redeploy
6. ✅ Check logs for success

### What Will Happen After Fix:
```
=== SENDING EMAIL ===
To: priyankrastogi145@gmail.com
Subject: Verify Your Gymie Account
From: Gymie <priyankrastogi14@gmail.com>
SMTP: smtp.gmail.com:587
==================
Attempting to send email...
✅ Email sent successfully!
```

---

## Quick Checklist

Before asking for help again, verify:

- [ ] Committed and pushed code changes
- [ ] Opened Render Dashboard → Environment tab
- [ ] Added SMTP_HOST = smtp.gmail.com
- [ ] Added SMTP_PORT = 587
- [ ] Added SMTP_USERNAME = priyankrastogi14@gmail.com
- [ ] Added SMTP_PASSWORD = rtcakrynvmgvdpos (NO SPACES!)
- [ ] Added SMTP_FROM = Gymie <priyankrastogi14@gmail.com>
- [ ] Deleted FROM_EMAIL (if it existed)
- [ ] Deleted FROM_NAME (if it existed)
- [ ] Clicked "Save Changes"
- [ ] Waited for redeploy to complete
- [ ] Checked logs for "✅ Email sent successfully"

If ALL checkboxes are ✅ and it still doesn't work, then it might be Render blocking SMTP.
