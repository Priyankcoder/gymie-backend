
# SendGrid Testing Script

A comprehensive script to test your SendGrid email configuration before deploying.

## 🎯 What It Does

This script will:
1. ✅ Load and validate your SendGrid configuration
2. ✅ Check for common configuration mistakes
3. ✅ Send a test email via SendGrid API
4. ✅ Show detailed results and troubleshooting steps
5. ✅ Verify sender email and API key validity

---

## 🚀 Quick Start

### Step 1: Set Environment Variables

Create or update `backend/.env`:

```bash
# Required
SENDGRID_API_KEY=SG.your_api_key_here
FROM_EMAIL=your-email@gmail.com
FROM_NAME=Gymie

# Optional (defaults to FROM_EMAIL)
TEST_EMAIL=recipient@example.com
```

**Important:** Make sure the .env file is at `backend/.env` (not in the root directory).

### Step 2: Run the Script

**From the backend directory:**
```bash
cd backend
go run scripts/test-sendgrid.go
```

**Or from project root:**
```bash
cd backend && go run scripts/test-sendgrid.go
```

The script will automatically find your `.env` file.

---

## 📋 Example Output

### ✅ Success (Status 202)

```
╔══════════════════════════════════════════════════════════════╗
║          SENDGRID EMAIL TESTING SCRIPT                       ║
╚══════════════════════════════════════════════════════════════╝

📋 STEP 1: Loading Configuration
─────────────────────────────────────
API Key: ✅ SG.u14zvkP****MJHl (69 chars)
From Email: ✅ test@gmail.com
From Name: ✅ Gymie
Test Recipient: ✅ test@gmail.com

🔍 STEP 2: Configuration Validation
─────────────────────────────────────
✅ Configuration looks good!

📧 STEP 3: Creating Test Email
─────────────────────────────────────
Subject: SendGrid Test Email - Gymie
From: Gymie <test@gmail.com>
To: test@gmail.com
Plain text: 455 bytes
HTML content: 120 bytes

📤 STEP 4: Sending Email via SendGrid
─────────────────────────────────────
Making API request to SendGrid...

📥 STEP 5: Processing Response
─────────────────────────────────────
Status Code: 202
Response Body: 

📊 STEP 6: Result
─────────────────────────────────────
✅ SUCCESS! Email sent successfully!

What this means:
  • SendGrid accepted the email
  • Email is being processed
  • You should receive it shortly

Next steps:
  1. Check your inbox: test@gmail.com
  2. Check spam folder if not in inbox
  3. Check SendGrid Activity Feed:
     https://app.sendgrid.com/email_activity

✅ SendGrid is configured correctly!

╔══════════════════════════════════════════════════════════════╗
║                    TEST COMPLETED                            ║
╚══════════════════════════════════════════════════════════════╝
```

---

### ❌ Error: Sender Not Verified (Status 401)

```
📊 STEP 6: Result
─────────────────────────────────────
❌ ERROR: Unauthorized (401)
Response: {"errors":[{"message":"Authenticated user is not authorized to send mail"}]}

This means:
  • Invalid API key, OR
  • Sender email not verified in SendGrid

To fix:
  1. Verify API key is correct
  2. Verify sender email at:
     https://app.sendgrid.com/settings/sender_auth/senders
  3. Make sure FROM_EMAIL matches verified sender
```

---

### ⚠️ Warning: Incorrect FROM_NAME Format

```
🔍 STEP 2: Configuration Validation
─────────────────────────────────────
⚠️  WARNING: FROM_NAME contains email address: 'Gymie <no-reply@gymie.fit>'
   FROM_NAME should only contain the name, not the email
   Example: 'Gymie' not 'Gymie <no-reply@example.com>'

   Would you like to continue anyway? (y/n):
```

---

## 🔧 Configuration Options

### Required Variables

| Variable | Example | Description |
|----------|---------|-------------|
| `SENDGRID_API_KEY` | `SG.xxxxx...` | Your SendGrid API key (69 chars) |
| `FROM_EMAIL` | `no-reply@example.com` | Verified sender email in SendGrid |
| `FROM_NAME` | `Gymie` | Display name (name only, no email) |

### Optional Variables

| Variable | Example | Description |
|----------|---------|-------------|
| `TEST_EMAIL` | `test@example.com` | Recipient for test email (defaults to FROM_EMAIL) |

---

## ✅ Common Fixes

### Issue 1: 401 Error - Sender Not Verified

**Fix:**
1. Go to: https://app.sendgrid.com/settings/sender_auth/senders
2. Click "Create New Sender"
3. Enter your FROM_EMAIL address
4. Click verification link in email
5. Confirm "Verified" status (green checkmark)

### Issue 2: FROM_NAME Contains Email

**Current (Wrong):**
```bash
FROM_NAME=Gymie <no-reply@gymie.fit>  ❌
```

**Correct:**
```bash
FROM_NAME=Gymie  ✅
FROM_EMAIL=no-reply@gymie.fit  ✅
```

### Issue 3: Invalid API Key

**Fix:**
1. Go to: https://app.sendgrid.com/settings/api_keys
2. Create new key with "Full Access" or "Mail Send" permission
3. Copy the key (starts with `SG.`)
4. Update `SENDGRID_API_KEY` in `.env`

---

## 🎯 What Each Status Code Means

| Status | Meaning | Action |
|--------|---------|--------|
| **202** | ✅ Success | Email sent! Check inbox |
| **400** | ❌ Bad Request | Check email format |
| **401** | ❌ Unauthorized | Verify API key or sender email |
| **403** | ❌ Forbidden | Verify sender or check account |
| **429** | ❌ Rate Limited | Wait and try again |
| **5xx** | ❌ Server Error | SendGrid issue, try later |

---

## 🔍 Troubleshooting

### "SENDGRID_API_KEY not set"

**Solution:**
```bash
# Add to backend/.env
SENDGRID_API_KEY=SG.your_actual_key_here

# Or export directly
export SENDGRID_API_KEY="SG.your_actual_key_here"
```

### "Sender email not verified"

**Solution:**
Your FROM_EMAIL must be verified in SendGrid:
1. https://app.sendgrid.com/settings/sender_auth/senders
2. Add your email
3. Verify it

### "Email not received"

**Check:**
1. ✅ Spam folder
2. ✅ SendGrid Activity Feed: https://app.sendgrid.com/email_activity
3. ✅ Email address is correct
4. ✅ Script showed "Status Code: 202"

---

## 🚨 Before Production

1. **Rotate exposed API key** (if it was in Git)
2. **Verify your production sender email**
3. **Test with this script first**
4. **Check SendGrid Activity Feed** after test
5. **Update Render environment variables**

---

## 📝 Quick Checklist

Before running the script:

- [ ] SendGrid account created
- [ ] API key generated (Full Access)
- [ ] Sender email verified in SendGrid
- [ ] `.env` file configured
- [ ] `FROM_NAME` is just the name (no email)
- [ ] `FROM_EMAIL` matches verified sender

After successful test:

- [ ] Email received in inbox
- [ ] Status code was 202
- [ ] SendGrid Activity Feed shows delivery
- [ ] Ready to deploy to production

---

## 💡 Tips

1. **Use personal email for testing** (easier to verify)
2. **Check spam folder** if email doesn't arrive
3. **Run test before each deployment** to catch config issues
4. **Monitor SendGrid dashboard** for delivery status

---

## 🔗 Useful Links

- **SendGrid Dashboard:** https://app.sendgrid.com/
- **Sender Verification:** https://app.sendgrid.com/settings/sender_auth/senders
- **API Keys:** https://app.sendgrid.com/settings/api_keys
- **Email Activity:** https://app.sendgrid.com/email_activity
- **SendGrid Docs:** https://docs.sendgrid.com/

---

**Time to run:** ~10 seconds  
**Difficulty:** Easy  
**Purpose:** Verify SendGrid before deployment 🚀
