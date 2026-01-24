
# Fix Email Sending Issues

## Current Error

```
❌ Email send failed: dial tcp 172.253.118.108:587: i/o timeout
```

**This means:** Your server can't connect to Gmail's SMTP server (smtp.gmail.com:587)

## Why This Happens

Gmail blocks SMTP connections by default for security. You need to:
1. Enable 2-Factor Authentication on your Gmail account
2. Create an "App Password" specifically for this application

## Fix: Set Up Gmail App Password

### Step 1: Enable 2FA on Gmail

1. Go to https://myaccount.google.com/security
2. Find "2-Step Verification"
3. Click "Get Started" and follow the setup
4. Verify with your phone

### Step 2: Create App Password

1. Go to https://myaccount.google.com/apppasswords
   - Or: Google Account → Security → 2-Step Verification → App passwords
2. Select app: **Mail**
3. Select device: **Other (Custom name)**
4. Enter name: **Gymie Backend**
5. Click **Generate**
6. **Copy the 16-character password** (looks like: `abcd efgh ijkl mnop`)

### Step 3: Update Backend Environment Variables

#### For Local Development (backend/.env):
```env
# Email Configuration
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=priyankrastogi14@gmail.com
SMTP_PASSWORD=abcd efgh ijkl mnop  # Replace with your App Password (remove spaces)
SMTP_FROM=Gymie <priyankrastogi14@gmail.com>
```

**IMPORTANT:** Remove spaces from the App Password!
```
Wrong: abcd efgh ijkl mnop
Right: abcdefghijklmnop
```

#### For Production (Render Dashboard):

1. Go to your Render service
2. Click **Environment** tab
3. Add these variables:

```
SMTP_HOST = smtp.gmail.com
SMTP_PORT = 587
SMTP_USERNAME = priyankrastogi14@gmail.com
SMTP_PASSWORD = abcdefghijklmnop  (your App Password, no spaces)
SMTP_FROM = Gymie <priyankrastogi14@gmail.com>
```

4. Click **Save Changes**
5. Service will auto-deploy

### Step 4: Update Frontend URL

In Render, also add:
```
FRONTEND_URL = https://your-frontend-url.com
```

This is used in verification email links.

### Step 5: Restart Backend

#### Local:
```bash
# Stop the server (Ctrl+C)
# Start again
go run cmd/api/main.go
```

#### Production (Render):
- Render will restart automatically after saving env vars
- Or click "Manual Deploy"

## Test Email Sending

### Test 1: From Terminal (Local)
```bash
curl -X POST http://localhost:8080/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "your-test-email@gmail.com",
    "username": "testuser",
    "password": "Test123!@#"
  }'
```

**Check server logs for:**
```
Attempting to send email...
=== SENDING EMAIL ===
To: your-test-email@gmail.com
Subject: Verify Your Gymie Account
From: Gymie <priyankrastogi14@gmail.com>
SMTP: smtp.gmail.com:587
==================
✅ Email sent successfully
```

### Test 2: From App

1. Open app in emulator
2. Try to register with a real email
3. Check your inbox for verification email

### Test 3: Resend Verification
```bash
curl -X POST http://localhost:8080/v1/auth/resend-verification \
  -H "Content-Type: application/json" \
  -d '{"email": "your-email@gmail.com"}'
```

## Alternative SMTP Providers

If Gmail doesn't work, try these:

### Option 1: SendGrid (Free Tier: 100 emails/day)
1. Sign up at https://sendgrid.com
2. Create an API key
3. Update env vars:
```env
SMTP_HOST=smtp.sendgrid.net
SMTP_PORT=587
SMTP_USERNAME=apikey
SMTP_PASSWORD=your-sendgrid-api-key
SMTP_FROM=your-verified-sender@domain.com
```

### Option 2: Mailgun (Free Tier: 5,000 emails/month)
1. Sign up at https://mailgun.com
2. Get SMTP credentials
3. Update env vars:
```env
SMTP_HOST=smtp.mailgun.org
SMTP_PORT=587
SMTP_USERNAME=postmaster@your-domain.mailgun.org
SMTP_PASSWORD=your-mailgun-password
SMTP_FROM=Gymie <noreply@your-domain.mailgun.org>
```

### Option 3: AWS SES (Free Tier: 62,000 emails/month)
1. Set up AWS SES
2. Verify your email/domain
3. Get SMTP credentials
4. Update env vars accordingly

## Common Issues & Solutions

### Issue 1: Still Getting Timeout

**Causes:**
- App Password not generated correctly
- Spaces in password
- 2FA not enabled
- Wrong Gmail account

**Solution:**
1. Delete old App Password
2. Create new one
3. Copy it without spaces
4. Update .env file
5. Restart server

### Issue 2: "Username and Password not accepted"

**Causes:**
- Wrong email/password
- Spaces in App Password
- Using account password instead of App Password

**Solution:**
1. Verify SMTP_USERNAME matches your Gmail
2. Verify you're using App Password, not account password
3. Remove all spaces from password

### Issue 3: "TLS Handshake Failed"

**Solution:**
Change port from 587 to 465:
```env
SMTP_PORT=465
```

Or try without TLS:
```env
SMTP_PORT=25  # Less secure
```

### Issue 4: Emails Going to Spam

**Solution:**
1. Add SPF record to your domain
2. Use a custom domain instead of Gmail
3. Warm up your email sender
4. Don't send too many emails too quickly

### Issue 5: Rate Limiting

**Gmail limits:**
- 500 emails/day for free accounts
- 2000 emails/day for Google Workspace

**Solution:**
- Use SendGrid or Mailgun for production
- Implement rate limiting in your app

## Verify Configuration

### Check Environment Variables Loaded
Add this to your backend temporarily:

```go
// In cmd/api/main.go
log.Println("SMTP Host:", cfg.SMTPHost)
log.Println("SMTP Port:", cfg.SMTPPort)
log.Println("SMTP User:", cfg.SMTPUsername)
log.Println("SMTP From:", cfg.SMTPFrom)
log.Println("SMTP Pass:", len(cfg.SMTPPassword), "characters")
```

Should output:
```
SMTP Host: smtp.gmail.com
SMTP Port: 587
SMTP User: priyankrastogi14@gmail.com
SMTP From: Gymie <priyankrastogi14@gmail.com>
SMTP Pass: 16 characters  (or your password length)
```

### Test SMTP Connection Directly

Create a test script:

```go
// backend/scripts/test-smtp.go
package main

import (
    "fmt"
    "net/smtp"
)

func main() {
    auth := smtp.PlainAuth(
        "",
        "priyankrastogi14@gmail.com",
        "your-app-password",  // Replace with your App Password
        "smtp.gmail.com",
    )

    err := smtp.SendMail(
        "smtp.gmail.com:587",
        auth,
        "priyankrastogi14@gmail.com",
        []string{"priyankrastogi145@gmail.com"},
        []byte("Subject: Test\r\n\r\nThis is a test email!"),
    )

    if err != nil {
        fmt.Println("❌ Failed:", err)
        return
    }

    fmt.Println("✅ Email sent successfully!")
}
```

Run:
```bash
cd backend
go run scripts/test-smtp.go
```

## Production Deployment Checklist

Before deploying to production:

- [ ] Generated Gmail App Password
- [ ] Added SMTP_HOST to Render env vars
- [ ] Added SMTP_PORT to Render env vars  
- [ ] Added SMTP_USERNAME to Render env vars
- [ ] Added SMTP_PASSWORD (App Password, no spaces) to Render env vars
- [ ] Added SMTP_FROM to Render env vars
- [ ] Added FRONTEND_URL to Render env vars
- [ ] Redeployed service on Render
- [ ] Tested registration with real email
- [ ] Received verification email
- [ ] Clicked verification link works

## Security Notes

1. **Never commit App Passwords to git**
   - Always use environment variables
   - Add .env to .gitignore

2. **Use App Passwords, not account password**
   - More secure
   - Can be revoked without changing account password

3. **Rotate passwords regularly**
   - Generate new App Password every few months
   - Delete old ones

4. **Consider using a dedicated email**
   - Create noreply@yourdomain.com
   - Don't use personal Gmail for production

## Summary

**Quick Fix:**
1. Enable 2FA on Gmail: https://myaccount.google.com/security
2. Create App Password: https://myaccount.google.com/apppasswords
3. Add to .env (no spaces):
   ```env
   SMTP_PASSWORD=abcdefghijklmnop
   ```
4. Restart server
5. Test registration

The timeout will be replaced with:
```
✅ Email sent successfully
```

And you'll receive the verification email in your inbox!
