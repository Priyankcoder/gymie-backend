# 🚀 Brevo Email Service Setup Guide

## Quick Start

### 1. Install Brevo Go SDK

```bash
cd backend
go get github.com/getbrevo/brevo-go/lib
```

### 2. Get Your API Key

1. Go to [Brevo Dashboard](https://app.brevo.com)
2. Navigate to **Settings** → **SMTP & API** → **API Keys**
3. Click **Generate a new API key**
4. Give it a name (e.g., "Gymie Backend")
5. Copy the API key (starts with `xkeysib-`)

### 3. Add to Environment Variables

Add to `backend/.env`:

```env
# Brevo Configuration
BREVO_API_KEY=xkeysib-your-actual-api-key-here
FROM_EMAIL=your-verified-sender@example.com
FROM_NAME=Gymie
TEST_EMAIL=your-email-to-receive-test@example.com
```

### 4. Verify Sender Email

**IMPORTANT:** Before sending emails, you must verify your sender email:

1. Go to [Brevo Senders](https://app.brevo.com/senders)
2. Click **Add a new sender**
3. Enter your email address (e.g., `no-reply@yourdomain.com`)
4. Brevo will send a verification email
5. Click the verification link
6. Wait for approval (usually instant for personal emails)

### 5. Run Test Script

```bash
cd backend
go run scripts/test-brevo.go
```

## Expected Output

### ✅ Success (HTTP 201)
```
✅ SUCCESS! Email sent successfully!

What this means:
  • Brevo accepted the email
  • Email is being processed
  • You should receive it shortly
  • Message ID: <unique-id>

✅ Brevo is configured correctly!
```

### ❌ Common Errors

#### 401 Unauthorized
- **Cause:** Invalid API key
- **Fix:** Check your API key at https://app.brevo.com/settings/keys/api

#### 403 Forbidden
- **Cause:** Sender email not verified
- **Fix:** Verify your sender at https://app.brevo.com/senders

#### 402 Payment Required
- **Cause:** No credits or quota exceeded
- **Fix:** Check your account at https://app.brevo.com/account/plan

## Brevo Free Tier Limits

- **300 emails/day** for free accounts
- Unlimited contacts
- Email API access
- Transactional email templates

## Next Steps

Once the test succeeds:

1. ✅ Create production email service implementation
2. ✅ Update backend configuration
3. ✅ Deploy to Render with Brevo credentials
4. ✅ Test verification and password reset flows

## Brevo vs SendGrid

| Feature | Brevo (Free) | SendGrid (Free) |
|---------|--------------|-----------------|
| Daily Limit | 300 emails | 100 emails |
| Setup | Simpler | More complex |
| Sender Verification | Easy | Domain Auth preferred |
| API | Simple REST | More features |
| Free Tier | Permanent | Trial-based |

## Support Resources

- [Brevo Documentation](https://developers.brevo.com/)
- [Go SDK GitHub](https://github.com/getbrevo/brevo-go)
- [Brevo API Reference](https://developers.brevo.com/reference)
- [Email Logs](https://app.brevo.com/log)
