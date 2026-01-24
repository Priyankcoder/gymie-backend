
# Fix Render Environment Variables for Email

## What Was Wrong

Your production email wasn't working because of **3 issues**:

### Issue 1: Config Struct Missing SMTP Fields ❌
The `Config` struct in `config.go` didn't have SMTP fields, so the email service couldn't read them.

**Fixed:** Added SMTP fields to config struct ✅

### Issue 2: Wrong Environment Variable Names ❌
You were using:
```
FROM_EMAIL=priyankrastogi14@gmail.com
FROM_NAME=Gymie
```

But the code expects:
```
SMTP_FROM=Gymie <priyankrastogi14@gmail.com>
```

**Fixed:** Updated .env.production to use `SMTP_FROM` ✅

### Issue 3: App Password Has Spaces ❌
Gmail App Password shown as: `rtca kryn vmgv dpos`
But it must be used without spaces: `rtcakrynvmgvdpos`

**Fixed:** Removed spaces from password ✅

## Update Render Environment Variables

Go to your Render dashboard and **update these variables**:

### Delete These (Wrong Names):
- ❌ `FROM_EMAIL`
- ❌ `FROM_NAME`

### Add/Update These (Correct):

```
SMTP_HOST = smtp.gmail.com
SMTP_PORT = 587
SMTP_USERNAME = priyankrastogi14@gmail.com
SMTP_PASSWORD = rtcakrynvmgvdpos
SMTP_FROM = Gymie <priyankrastogi14@gmail.com>
FRONTEND_URL = https://gymie.fit
```

**IMPORTANT:** 
1. Remove ALL spaces from `SMTP_PASSWORD`
2. Gmail shows it as `rtca kryn vmgv dpos` but use `rtcakrynvmgvdpos`
3. Include angle brackets in `SMTP_FROM`: `Gymie <email@gmail.com>`

## Steps to Update

### 1. Go to Render Dashboard
https://dashboard.render.com

### 2. Select Your Backend Service
Click on your Gymie backend service

### 3. Go to Environment Tab
Click "Environment" in the left sidebar

### 4. Delete Old Variables
Find and delete:
- `FROM_EMAIL`
- `FROM_NAME`

### 5. Add/Update Variables
Click "Add Environment Variable" for each:

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

### 6. Save Changes
Click "Save Changes" at the bottom

### 7. Render Will Auto-Deploy
The service will automatically redeploy with new variables

## Verify It's Working

### Check Deployment Logs
After deployment completes, check logs for:

```
Database connection established
Running auto-migration...
✅ Auto-migration completed successfully
Starting server on port 8080
```

### Test Email Sending
Register a new account or use resend verification:

```bash
curl -X POST https://gymie-api.onrender.com/v1/auth/resend-verification \
  -H "Content-Type: application/json" \
  -d '{"email":"priyankrastogi145@gmail.com"}'
```

**Check logs for:**
```
Attempting to send email...
=== SENDING EMAIL ===
To: priyankrastogi145@gmail.com
Subject: Verify Your Gymie Account
From: Gymie <priyankrastogi14@gmail.com>
SMTP: smtp.gmail.com:587
==================
✅ Email sent successfully  # Should see this!
```

## About Gmail App Passwords

### What is an App Password?
A 16-character password that lets apps access your Gmail without your main password.

### How Gmail Shows It:
```
rtca kryn vmgv dpos  (with spaces for readability)
```

### How You Must Use It:
```
rtcakrynvmgvdpos  (no spaces!)
```

### Why Spaces Cause Issues:
Environment variables include the spaces as part of the password, making it invalid.

### Where to Get Your App Password:
1. Go to: https://myaccount.google.com/apppasswords
2. You'll see your existing "Gymie Backend" password
3. If you need to see it again, delete and create a new one
4. Copy it WITHOUT the spaces

## Alternative: Create New App Password

If the current one doesn't work:

### 1. Delete Old App Password
- Go to https://myaccount.google.com/apppasswords
- Find "Gymie Backend"
- Click trash icon to delete

### 2. Create New One
- Click "Select app" → Mail
- Click "Select device" → Other (Custom name)
- Enter: "Gymie Backend Production"
- Click "Generate"
- Copy the password (remove spaces!)

### 3. Update Render
- Go to Render → Environment
- Update `SMTP_PASSWORD` with new password (no spaces)
- Save

## About Whitelisting/Less Secure Apps

**Good News:** You don't need to whitelist domains or enable "Less Secure Apps"!

### Why It Works:
1. ✅ You're using App Password (not account password)
2. ✅ App Passwords bypass "Less Secure Apps" setting
3. ✅ No domain whitelisting needed for Gmail SMTP
4. ✅ Works from any IP address/server

### What You DO Need:
- ✅ 2-Factor Authentication enabled
- ✅ Valid App Password
- ✅ Correct password (no spaces)
- ✅ Correct environment variable names

## Production vs Development

### Development (.env):
```env
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=priyankrastogi14@gmail.com
SMTP_PASSWORD=rtcakrynvmgvdpos
SMTP_FROM=Gymie <priyankrastogi14@gmail.com>
FRONTEND_URL=http://localhost:3000
```

### Production (Render Environment):
```env
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=priyankrastogi14@gmail.com
SMTP_PASSWORD=rtcakrynvmgvdpos
SMTP_FROM=Gymie <priyankrastogi14@gmail.com>
FRONTEND_URL=https://gymie.fit
```

**Same SMTP settings, different FRONTEND_URL!**

## Troubleshooting

### Still Getting Timeout?
1. Verify 2FA is enabled on Gmail
2. Delete and regenerate App Password
3. Double-check no spaces in password
4. Restart Render service manually

### "Username and Password not accepted"?
1. Make sure using App Password, not account password
2. Remove all spaces from password
3. Verify email address is correct

### Emails Going to Spam?
This is normal for Gmail accounts. Solutions:
- Ask users to check spam folder
- Use SendGrid or Mailgun instead
- Set up custom domain with SPF records

## After Fixing

Once deployed with correct variables:

1. ✅ Registration emails will be sent
2. ✅ Verification links will work
3. ✅ "Resend verification" will work
4. ✅ No more timeout errors

## Summary

**Changes Made:**
1. ✅ Added SMTP fields to Config struct
2. ✅ Fixed environment variable names
3. ✅ Removed spaces from App Password

**Action Required:**
1. Commit and push the code changes
2. Update Render environment variables (see above)
3. Wait for auto-deployment
4. Test email sending

**Expected Result:**
```
✅ Email sent successfully
```

No more timeout errors! 🎉
