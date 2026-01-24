
# 📧 Email Debugging Guide

## 🎯 Comprehensive Logging Added

I've added detailed step-by-step debugging logs to track every phase of email sending. Here's what you'll see now:

---

## 🔍 What You'll See in Logs

### 1. Service Initialization (On Startup)
```
╔══════════════════════════════════════════════════════════════╗
║              EMAIL SERVICE INITIALIZATION                    ║
╚══════════════════════════════════════════════════════════════╝
✅ SendGrid API Key detected
   └─ Using SendGrid Email Service (Production)
```

or if SendGrid key is missing:
```
⚠️  No SendGrid API Key found
   └─ Falling back to Gmail SMTP Service
```

---

### 2. Email Preparation (When Sending)
```
╔══════════════════════════════════════════════════════════════╗
║     PREPARING VERIFICATION EMAIL - SendGrid Service         ║
╚══════════════════════════════════════════════════════════════╝
[PREPARE] Building verification email...
  ├─ To: user@example.com
  ├─ User Name: John Doe
  ├─ Token: abc123xyz...
  ├─ Frontend URL: https://gymie.fit
  └─ Verification Link: https://gymie.fit/verify-email?token=abc123xyz
```

---

### 3. SendGrid Sending Process (6 Steps)
```
╔══════════════════════════════════════════════════════════════╗
║          SENDGRID EMAIL SENDING - DEBUG LOG                  ║
╚══════════════════════════════════════════════════════════════╝

[STEP 1] Validating SendGrid configuration...
  ├─ API Key Length: 69 characters
  ├─ API Key Prefix: SG.u14zvkP****
  ├─ From Email: priyankrastogi14@gmail.com
  ├─ From Name: Gymie
  ├─ Frontend URL: https://gymie.fit
  └─ To Email: user@example.com

[STEP 2] Creating email message...
  ├─ From: Gymie <priyankrastogi14@gmail.com>
  ├─ To: user@example.com
  ├─ Subject: Verify Your Gymie Account
  ├─ Plain text length: 245 bytes
  └─ HTML content length: 4823 bytes

[STEP 3] Initializing SendGrid client...
  └─ Client created successfully

[STEP 4] Sending email via SendGrid API...
  ├─ Endpoint: https://api.sendgrid.com/v3/mail/send
  └─ Making API request...

[STEP 6] Processing SendGrid response...
  ├─ Status Code: 202
  ├─ Response Headers: map[...]

✅ [SUCCESS] Email sent via SendGrid!
  ├─ Status Code: 202
  ├─ Message ID: <message-id>
  └─ Email delivered to SendGrid successfully

═══════════════════════════════════════════════════════════════
```

---

### 4. Gmail SMTP Process (Alternative)
```
╔══════════════════════════════════════════════════════════════╗
║          GMAIL SMTP EMAIL SENDING - DEBUG LOG               ║
╚══════════════════════════════════════════════════════════════╝

[STEP 1] Validating SMTP configuration...
  ├─ SMTP Host: smtp.gmail.com
  ├─ SMTP Port: 587
  ├─ SMTP Username: your-email@gmail.com
  ├─ SMTP Password: rt****dpos
  ├─ From Email: priyankrastogi14@gmail.com
  ├─ From Name: Gymie
  └─ To Email: user@example.com

[STEP 2] Creating email message...
  ├─ From: Gymie <priyankrastogi14@gmail.com>
  ├─ To: user@example.com
  ├─ Subject: Verify Your Gymie Account
  └─ Body length: 4823 bytes

[STEP 3] Creating SMTP dialer...
  ├─ Dialer created for smtp.gmail.com:587
  └─ TLS enabled with ServerName: smtp.gmail.com

[STEP 4] Connecting to SMTP server...
  ├─ Attempting connection to smtp.gmail.com:587
  └─ Using STARTTLS...

✅ [SUCCESS] Email sent via Gmail SMTP!
  └─ Message delivered to SMTP server successfully

═══════════════════════════════════════════════════════════════
```

---

## 🐛 Common Errors & Solutions

### SendGrid Errors

#### ❌ Error: Empty API Key
```
❌ [ERROR] SendGrid API key is empty!
```
**Solution:** Add `SENDGRID_API_KEY` to environment variables

#### ❌ Error: Status 401 (Unauthorized)
```
❌ [ERROR] SendGrid returned error status
  ├─ Status Code: 401
  └─ Common causes:
      • Invalid API key
```
**Solution:** 
1. Check API key is correct
2. Regenerate key in SendGrid dashboard
3. Update environment variable

#### ❌ Error: Status 403 (Forbidden)
```
❌ [ERROR] SendGrid returned error status
  ├─ Status Code: 403
  └─ Common causes:
      • Forbidden (check sender verification)
```
**Solution:**
1. Verify sender email in SendGrid dashboard
2. Complete sender authentication
3. Check domain verification

#### ❌ Error: Status 400 (Bad Request)
```
❌ [ERROR] SendGrid returned error status
  ├─ Status Code: 400
  └─ Common causes:
      • Invalid request (check email format)
```
**Solution:**
1. Check recipient email format
2. Verify FROM_EMAIL is correct
3. Check content encoding

#### ❌ Error: Network/API Request Failed
```
❌ [STEP 5] SENDGRID API ERROR
  ├─ Error Message: connection refused
  └─ This usually means:
      • Network connectivity issues
      • Firewall blocking port 443
      • SendGrid service down
```
**Solution:**
1. Check internet connection
2. Test: `curl https://api.sendgrid.com/v3/mail/send`
3. Check firewall rules for port 443
4. Check SendGrid status: https://status.sendgrid.com/

---

### Gmail SMTP Errors

#### ❌ Error: Empty Credentials
```
❌ [ERROR] SMTP credentials are empty!
```
**Solution:** Add `SMTP_USERNAME` and `SMTP_PASSWORD`

#### ❌ Error: Authentication Failed
```
❌ [STEP 5] SMTP SEND FAILED
  ├─ Error Message: 535 authentication failed
  └─ Common causes:
      • Invalid credentials (check App Password)
```
**Solution:**
1. Generate new App Password at: https://myaccount.google.com/apppasswords
2. Use App Password (not regular password)
3. Update `SMTP_PASSWORD` with new App Password

#### ❌ Error: Connection Timeout
```
❌ [STEP 5] SMTP SEND FAILED
  ├─ Error Message: dial tcp: i/o timeout
  └─ Common causes:
      • Network/firewall blocking port 587
```
**Solution:**
1. Check firewall allows port 587
2. Test: `telnet smtp.gmail.com 587`
3. Use SendGrid instead (more reliable)

---

## 📊 Debugging Workflow

### Step 1: Check Service Initialization
Look for this at startup:
```
╔══════════════════════════════════════════════════════════════╗
║              EMAIL SERVICE INITIALIZATION                    ║
╚══════════════════════════════════════════════════════════════╝
```

**What to check:**
- ✅ SendGrid detected? Or ⚠️ fallback to Gmail?
- Is this what you expected?

### Step 2: Trigger Email Send
Register a new user or request password reset

### Step 3: Follow the Logs
Watch for:
1. **STEP 1**: Configuration validation
   - Is API key present?
   - Are credentials correct?

2. **STEP 2**: Message creation
   - Is recipient correct?
   - Is link properly formed?

3. **STEP 3-4**: API/SMTP connection
   - Any connection errors?
   - Network issues?

4. **STEP 5-6**: Response
   - Status code 202 (SendGrid) or success (SMTP)?
   - Any error messages?

### Step 4: Check Result
- ✅ Success? Check email inbox
- ❌ Failed? Look at error message and match with solutions above

---

## 🧪 Testing Commands

### Test SendGrid (Production)
```bash
# Set env vars
export SENDGRID_API_KEY="SG.your-key"
export FROM_EMAIL="priyankrastogi14@gmail.com"
export FROM_NAME="Gymie"
export FRONTEND_URL="https://gymie.fit"

# Run server
cd backend
go run cmd/api/main.go

# Watch logs for email sending
```

### Test Gmail SMTP (Fallback)
```bash
# Unset SendGrid key
unset SENDGRID_API_KEY

# Set SMTP vars
export SMTP_HOST="smtp.gmail.com"
export SMTP_PORT="587"
export SMTP_USERNAME="your-email@gmail.com"
export SMTP_PASSWORD="your-app-password"
export FROM_EMAIL="your-email@gmail.com"
export FROM_NAME="Gymie"
export FRONTEND_URL="https://gymie.fit"

# Run server
go run cmd/api/main.go
```

---

## 📝 What Changed

### Files Modified:
1. [`email_service_sendgrid.go`](backend/internal/services/email_service_sendgrid.go:1)
   - Added 6-step detailed logging
   - Shows API endpoint, status codes, headers
   - Provides specific error solutions

2. [`email_service.go`](backend/internal/services/email_service.go:1)
   - Added 5-step SMTP logging
   - Shows connection details
   - Diagnoses common SMTP issues

3. [`service.go`](backend/internal/service/service.go:1)
   - Shows which service is selected at startup
   - Clear indication of SendGrid vs SMTP

---

## 🎯 Quick Diagnosis

### Email Not Sending?

**Check logs for:**

1. **Service Selection**
   ```
   ✅ SendGrid API Key detected
   ```
   or
   ```
   ⚠️  No SendGrid API Key found
   ```

2. **Configuration**
   ```
   [STEP 1] Validating SendGrid configuration...
   ```
   - API key present?
   - FROM_EMAIL set?

3. **API Call**
   ```
   [STEP 4] Sending email via SendGrid API...
   ```
   - Did it reach this step?
   - Any error before?

4. **Response**
   ```
   ✅ [SUCCESS] Email sent via SendGrid!
   ```
   or
   ```
   ❌ [ERROR] SendGrid returned error status
   ```

5. **Status Code**
   - `202` = Success (SendGrid)
   - `401` = Bad API key
   - `403` = Sender not verified
   - `400` = Invalid request
   - `500+` = SendGrid server error

---

## 🚀 Production Checklist

Before deploying:
- [ ] Logs show "SendGrid API Key detected"
- [ ] FROM_EMAIL matches verified sender
- [ ] FROM_NAME is set correctly
- [ ] FRONTEND_URL is production URL
- [ ] Test registration sends email
- [ ] Email link works correctly
- [ ] Check spam folder if not in inbox
- [ ] Monitor SendGrid dashboard for delivery

---

## 📧 Need Help?

If email still fails after checking logs:

1. Copy the complete log output (all 6 steps)
2. Note which step failed
3. Check the error message
4. Match with solutions above
5. Verify environment variables in Render dashboard

The logs now show **exactly** where and why email sending fails! 🎯
