# 🔄 Migration Guide: SendGrid → Brevo

This document outlines the complete migration from SendGrid to Brevo for email services.

## ✅ What Was Changed

### 1. **Email Service Implementation**
- ❌ Removed: `email_service_sendgrid.go`
- ✅ Added: `email_service_brevo.go`
- Same interface, different provider

### 2. **Configuration**
- ❌ Old: `SENDGRID_API_KEY`
- ✅ New: `BREVO_API_KEY`
- All other env vars remain the same (`FROM_EMAIL`, `FROM_NAME`, `FRONTEND_URL`)

### 3. **Dependencies**
- ❌ Removed: `github.com/sendgrid/sendgrid-go`
- ✅ Added: `github.com/getbrevo/brevo-go/lib`

## 📋 Migration Steps

### Step 1: Install Brevo SDK

```bash
cd backend
go get github.com/getbrevo/brevo-go/lib
go mod tidy
```

### Step 2: Update Local Environment

Update `backend/.env`:

```env
# Replace this:
# SENDGRID_API_KEY=SG.xxxxx...

# With this:
BREVO_API_KEY=xkeysib-xxxxx...
FROM_EMAIL=your-verified-email@example.com
FROM_NAME=Gymie
FRONTEND_URL=http://localhost:3000
```

### Step 3: Verify Sender in Brevo

1. Go to https://app.brevo.com/senders
2. Click **Add a new sender**
3. Enter your email address
4. Verify via the confirmation email
5. Wait for approval (usually instant)

### Step 4: Test Locally

```bash
cd backend
go run scripts/test-brevo.go
```

Expected output:
```
✅ SUCCESS! Email sent successfully!
Message ID: <unique-id>
✅ Brevo is configured correctly!
```

### Step 5: Update Render Environment Variables

1. Go to https://dashboard.render.com
2. Select your `gymie-api` service
3. Go to **Environment** tab
4. **Delete** the old variable:
   - `SENDGRID_API_KEY`
5. **Add** new variable:
   - Key: `BREVO_API_KEY`
   - Value: `xkeysib-xxxxx...` (your Brevo API key)
6. **Update** `FROM_EMAIL` if needed:
   - Must match the verified sender in Brevo

### Step 6: Deploy

The next deployment will automatically use Brevo:

```bash
git add .
git commit -m "Migrate from SendGrid to Brevo for email service"
git push origin main
```

Render will auto-deploy and restart with Brevo.

## 🔍 Verification

### Local Testing
```bash
cd backend
go run scripts/test-brevo.go
```

### Production Testing
1. Sign up for a new account on your app
2. Check that verification email arrives
3. Check Brevo logs: https://app.brevo.com/log

## ⚠️ Important Notes

### API Key Format
- **SendGrid:** `SG.xxxxx...` (69 chars)
- **Brevo:** `xkeysib-xxxxx...` (64+ chars)

### Sender Verification
- **SendGrid:** Domain authentication preferred
- **Brevo:** Individual email verification is easier

### Rate Limits
- **SendGrid Free:** 100 emails/day
- **Brevo Free:** 300 emails/day ✅

### Success Response
- **SendGrid:** HTTP 202 (Accepted)
- **Brevo:** HTTP 201 (Created)

## 🐛 Troubleshooting

### Error: "Invalid API key"
- Check that `BREVO_API_KEY` starts with `xkeysib-`
- Verify key at: https://app.brevo.com/settings/keys/api

### Error: "Sender not verified"
- Verify sender at: https://app.brevo.com/senders
- Make sure `FROM_EMAIL` matches verified sender

### Error: "Account needs payment"
- Check daily quota at: https://app.brevo.com/account/plan
- Free tier: 300 emails/day

## 📊 Comparison

| Feature | SendGrid (Old) | Brevo (New) |
|---------|----------------|-------------|
| Free Emails/Day | 100 | 300 ✅ |
| Setup Complexity | High | Low ✅ |
| Sender Verification | Domain Auth | Email Verification ✅ |
| API Simplicity | Moderate | Simple ✅ |
| Free Tier Stability | Trial-based | Permanent ✅ |

## 🎯 Why Brevo?

1. **Higher free tier** - 300 emails/day vs 100
2. **Simpler setup** - Just verify email, no domain auth needed
3. **Better stability** - Permanent free tier
4. **Easier verification** - No complex DNS records
5. **Better documentation** - Clearer Go SDK examples

## ✅ Migration Checklist

- [ ] Install Brevo SDK (`go get github.com/getbrevo/brevo-go/lib`)
- [ ] Get Brevo API key from https://app.brevo.com
- [ ] Verify sender email in Brevo
- [ ] Update local `.env` with `BREVO_API_KEY`
- [ ] Run test script successfully
- [ ] Update Render environment variables
- [ ] Remove old `SENDGRID_API_KEY` from Render
- [ ] Deploy to production
- [ ] Test signup/verification flow
- [ ] Monitor Brevo logs for 24 hours

## 🆘 Rollback Plan

If you need to rollback to SendGrid:

1. Revert the code changes:
   ```bash
   git revert HEAD
   git push origin main
   ```

2. Update Render environment variables:
   - Remove: `BREVO_API_KEY`
   - Add back: `SENDGRID_API_KEY`

3. Verify SendGrid sender authentication

## 📚 Resources

- [Brevo Documentation](https://developers.brevo.com/)
- [Brevo Go SDK](https://github.com/getbrevo/brevo-go)
- [Brevo Dashboard](https://app.brevo.com)
- [Sender Verification](https://app.brevo.com/senders)
- [API Keys](https://app.brevo.com/settings/keys/api)
- [Email Logs](https://app.brevo.com/log)

---

**Migration completed successfully! 🎉**

Your backend now uses Brevo for all transactional emails with a higher free tier and simpler management.
