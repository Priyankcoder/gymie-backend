
# SendGrid Integration Guide

## Why SendGrid Instead of Gmail SMTP?

### Problems with Gmail SMTP:
- ❌ Timeout issues from cloud providers (Render blocks port 587)
- ❌ Daily sending limits (500 emails/day for free Gmail)
- ❌ Requires App Passwords and 2FA setup
- ❌ Emails often go to spam
- ❌ Not reliable for production

### Benefits of SendGrid:
- ✅ 100 emails/day free tier
- ✅ No firewall/port blocking issues (uses HTTPS API)
- ✅ Better deliverability (emails don't go to spam)
- ✅ Email analytics and tracking
- ✅ Production-ready and reliable

## Setup Steps

### Step 1: Create SendGrid Account

1. Go to: https://sendgrid.com/
2. Click "Start for free"
3. Sign up with your email
4. Verify your email address

### Step 2: Verify Sender Identity

**Option A: Single Sender Verification (Quick)**
1. Go to: https://app.sendgrid.com/settings/sender_auth/senders
2. Click "Create New Sender"
3. Fill in:
   - From Name: `Gymie`
   - From Email: `noreply@gymie.com` (or your email)
   - Reply To: Your support email
   - Company Address: Your address
4. Click "Create"
5. **Check your email** and verify the sender

**Option B: Domain Authentication (Professional)**
1. Go to: https://app.sendgrid.com/settings/sender_auth
2. Click "Authenticate Your Domain"
3. Select your DNS provider
4. Add the DNS records they provide
5. Wait for verification (can take 24-48 hours)

### Step 3: Create API Key

1. Go to: https://app.sendgrid.com/settings/api_keys
2. Click "Create API Key"
3. Name: `Gymie Backend`
4. API Key Permissions: **Full Access** (or just "Mail Send")
5. Click "Create & View"
6. **Copy the API key** (you won't see it again!)
   - Format: `SG.xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx`

### Step 4: Add to Environment Variables

#### Local Development (backend/.env):
```env
# SendGrid Configuration
SENDGRID_API_KEY=SG.your-api-key-here
FRONTEND_URL=http://localhost:3000
```

#### Production (Render Dashboard):
1. Go to Render Dashboard → Your service → Environment
2. Add these variables:
```
SENDGRID_API_KEY = SG.your-api-key-here
FRONTEND_URL = https://gymie.fit
```
3. Click "Save Changes"

### Step 5: Install SendGrid Go Library

```bash
cd backend
go get github.com/sendgrid/sendgrid-go
```

### Step 6: Update Service Initialization

The service file [`email_service_sendgrid.go`](backend/internal/services/email_service_sendgrid.go) is already created.

Now update your service initialization in `internal/service/service.go`:

```go
package service

import (
	"github.com/yourusername/gymie-backend/internal/config"
	"github.com/yourusername/gymie-backend/internal/repository"
	"github.com/yourusername/gymie-backend/internal/services"
)

type Services struct {
	Auth            *AuthService
	User            *UserService
	Workout         *WorkoutService
	Nutrition       *NutritionService
	Progress        *ProgressService
	WorkoutPlan     *WorkoutPlanService
	OfflineNutrition *OfflineNutritionService
}

func NewServices(repos *repository.Repositories, cfg *config.Config) *Services {
	// Initialize email service based on configuration
	var emailService EmailServiceInterface
	
	if cfg.SendGridAPIKey != "" {
		// Use SendGrid if API key is provided (preferred)
		emailService = services.NewSendGridEmailService(cfg)
	} else {
		// Fall back to SMTP (Gmail)
		emailService = services.NewEmailService()
	}
	
	return &Services{
		Auth:            NewAuthService(repos.User, cfg, emailService),
		User:            NewUserService(repos.User),
		Workout:         NewWorkoutService(repos.Workout),
		Nutrition:       NewNutritionService(repos.Nutrition),
		Progress:        NewProgressService(repos.Progress),
		WorkoutPlan:     NewWorkoutPlanService(repos.WorkoutPlan),
		OfflineNutrition: NewOfflineNutritionService(repos.OfflineNutrition),
	}
}

// EmailServiceInterface defines the contract for email services
type EmailServiceInterface interface {
	SendVerificationEmail(toEmail, userName, verificationToken string) error
	SendPasswordResetEmail(toEmail, userName, resetToken string) error
}
```

### Step 7: Test SendGrid Integration

Create a test script:

```bash
cd backend
go run scripts/test-sendgrid.go
```

Create `backend/scripts/test-sendgrid.go`:

```go
package main

import (
	"fmt"
	"log"
	
	"github.com/joho/godotenv"
	"github.com/yourusername/gymie-backend/internal/config"
	"github.com/yourusername/gymie-backend/internal/services"
)

func main() {
	// Load environment
	_ = godotenv.Load(".env")
	
	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	
	if cfg.SendGridAPIKey == "" {
		log.Fatal("SENDGRID_API_KEY not set in environment")
	}
	
	// Create SendGrid email service
	emailService := services.NewSendGridEmailService(cfg)
	
	// Test email
	testEmail := "your-email@gmail.com"
	testName := "Test User"
	testToken := "test-token-123"
	
	fmt.Printf("Sending test email to: %s\n", testEmail)
	
	err = emailService.SendVerificationEmail(testEmail, testName, testToken)
	if err != nil {
		log.Fatalf("Failed to send email: %v", err)
	}
	
	fmt.Println("✅ Email sent successfully!")
	fmt.Println("Check your inbox (and spam folder)")
}
```

Run it:
```bash
SENDGRID_API_KEY=SG.your-key go run scripts/test-sendgrid.go
```

## Verification Checklist

- [ ] SendGrid account created
- [ ] Sender verified (noreply@gymie.com)
- [ ] API key created and copied
- [ ] SENDGRID_API_KEY added to local .env
- [ ] SENDGRID_API_KEY added to Render environment
- [ ] `go get github.com/sendgrid/sendgrid-go` executed
- [ ] Service initialization updated
- [ ] Test email sent successfully
- [ ] Email received (check spam folder)

## Troubleshooting

### Error: "API key does not start with 'SG.'"

**Solution:** Make sure you copied the complete API key including the `SG.` prefix.

### Error: "The from email does not match a verified Sender Identity"

**Solution:** 
1. Go to SendGrid → Settings → Sender Authentication
2. Verify that `noreply@gymie.com` is in the verified senders list
3. Check your email and click the verification link if pending

### Email not received

**Check:**
1. Spam/Junk folder
2. SendGrid Activity Feed: https://app.sendgrid.com/email_activity
3. Make sure sender is verified
4. Check Render logs for errors

### Error: "forbidden"

**Solution:** API key doesn't have Mail Send permission. Create a new API key with Full Access.

## SendGrid vs Gmail Comparison

| Feature | Gmail SMTP | SendGrid |
|---------|------------|----------|
| **Setup** | Medium (App Password + 2FA) | Easy (Just API key) |
| **Free Tier** | 500 emails/day | 100 emails/day |
| **Deliverability** | Often goes to spam | Professional delivery |
| **Cloud Compatible** | ❌ Blocked by Render | ✅ Works everywhere |
| **Analytics** | ❌ None | ✅ Full dashboard |
| **Reliability** | Low (timeouts) | High |
| **Production Ready** | ❌ Not recommended | ✅ Recommended |

## Migration Steps

If you're switching from Gmail SMTP to SendGrid:

1. ✅ Keep existing SMTP_* variables (for fallback)
2. ✅ Add SENDGRID_API_KEY
3. ✅ Service auto-detects and prefers SendGrid
4. ✅ Test thoroughly
5. ✅ Remove SMTP_* variables once confirmed working

## Cost

### Free Tier:
- 100 emails/day forever
- 2,000 contacts
- Email API access
- **Perfect for your app!**

### Essentials Plan ($19.95/month):
- 50,000 emails/month
- Email API & SMTP
- 24/7 Support

For Gymie, the **free tier is more than enough**!

## Next Steps

1. ✅ Complete setup above
2. ✅ Test email sending
3. ✅ Update `.env.production` 
4. ✅ Deploy to Render
5. ✅ Test production emails
6. ✅ Monitor SendGrid dashboard

## Support

- SendGrid Docs: https://docs.sendgrid.com/
- API Reference: https://docs.sendgrid.com/api-reference/mail-send/mail-send
- Support: https://support.sendgrid.com/

---

**Summary:** SendGrid is the professional, reliable choice for email delivery in production. The free tier is perfect for your needs, and it just works!
