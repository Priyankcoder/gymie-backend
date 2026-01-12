
# Production Environment Setup

## ⚠️ IMPORTANT: Secure Your .env.production File

The [`backend/.env.production`](backend/.env.production:1) file contains sensitive credentials. Make sure it's properly secured:

### ✅ Already Protected
Your project's [`.gitignore`](.gitignore:1) already excludes `.env` files, so `.env.production` won't be committed to Git.

### 🔐 Before Deploying

**Update these values in [`backend/.env.production`](backend/.env.production:1):**

1. **DATABASE_URL** - Get from Render.com PostgreSQL database
   - Use the **INTERNAL** connection string (starts with `postgresql://`)
   - Format: `postgresql://user:password@internal-host:5432/database`

2. **REDIS_URL** - Get from Upstash.com
   - Format: `redis://default:password@host:port`

3. **JWT_SECRET** - Generate a secure random string
   ```bash
   # Generate secure secret (copy the output)
   openssl rand -base64 32
   ```
   - Must be at least 32 characters
   - Keep this secret! Anyone with this can forge tokens

### 📝 How to Use on Render.com

You **DON'T** need to upload `.env.production` to Render. Instead:

1. Go to Render Dashboard → Your Service → Environment
2. Add each variable manually:
   - `DATABASE_URL` = (from Render PostgreSQL)
   - `REDIS_URL` = (from Upstash)
   - `JWT_SECRET` = (your generated secret)
   - `PORT` = `8080`
   - `ENVIRONMENT` = `production`
   - `CORS_ALLOWED_ORIGINS` = `*`
   - `JWT_EXPIRATION` = `24`
   - `RATE_LIMIT_ENABLED` = `true`
   - `RATE_LIMIT_REQUESTS` = `100`
   - `RATE_LIMIT_WINDOW` = `60`

3. Render will automatically inject these into your container

### 🔄 Local Testing with Production Config

If you want to test locally with production settings:

```bash
cd backend

# Load production env (on macOS/Linux)
export $(cat .env.production | grep -v '^#' | xargs)

# Or use this command
set -a && source .env.production && set +a

# Run server
go run cmd/api/main.go
```

### 🚨 Security Checklist

- [ ] `.env.production` is NOT committed to git
- [ ] `JWT_SECRET` is at least 32 characters
- [ ] `JWT_SECRET` is randomly generated (not a dictionary word)
- [ ] `DATABASE_URL` uses INTERNAL host (for Render)
- [ ] Production values are only in Render dashboard, not in code
- [ ] Never share your `.env.production` file

### 🎯 Quick Setup

```bash
# 1. Copy the template
cp backend/.env.production.example backend/.env.production

# 2. Generate JWT secret
openssl rand -base64 32

# 3. Edit .env.production with your values
nano backend/.env.production

# 4. Add to Render (don't commit this file!)
```

### 📚 Related Files

- [`backend/.env.production`](backend/.env.production:1) - Your production config (DO NOT COMMIT)
- [`backend/.env.production.example`](backend/.env.production.example:1) - Template (safe to commit)
- [`backend/.env.example`](backend/.env.example:1) - Development template

---

**Remember**: The `.env.production` file is for your reference only. Always configure environment variables through the Render dashboard for actual deployment.
