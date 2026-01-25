package services

import (
	"context"
	"fmt"
	"log"

	brevo "github.com/getbrevo/brevo-go/lib"
	"github.com/yourusername/gymie-backend/internal/config"
)

type BrevoEmailService struct {
	client      *brevo.APIClient
	apiKey      string
	fromEmail   string
	fromName    string
	frontendURL string
}

func NewBrevoEmailService(cfg *config.Config) *BrevoEmailService {
	// Validate API key
	if cfg.BrevoAPIKey == "" {
		log.Println("❌ CRITICAL ERROR: Brevo API key is empty!")
		panic("BREVO_API_KEY environment variable is required")
	}

	if len(cfg.BrevoAPIKey) < 20 {
		log.Printf("❌ CRITICAL ERROR: Brevo API key is too short (%d chars)\n", len(cfg.BrevoAPIKey))
		log.Println("   Expected format: xkeysib-xxxxx... (typically 64+ characters)")
		panic("Invalid BREVO_API_KEY - key appears to be incomplete")
	}

	// Initialize Brevo client
	brevoConfig := brevo.NewConfiguration()
	brevoConfig.AddDefaultHeader("api-key", cfg.BrevoAPIKey)
	client := brevo.NewAPIClient(brevoConfig)

	service := &BrevoEmailService{
		client:      client,
		apiKey:      cfg.BrevoAPIKey,
		fromEmail:   cfg.FromEmail,
		fromName:    cfg.FromName,
		frontendURL: cfg.FrontendURL,
	}

	log.Println("=== BREVO EMAIL SERVICE INITIALIZED ===")
	// Safe logging - only show first 10 and last 4 chars
	maskedKey := cfg.BrevoAPIKey[:min(10, len(cfg.BrevoAPIKey))] + "****"
	if len(cfg.BrevoAPIKey) > 4 {
		maskedKey += cfg.BrevoAPIKey[len(cfg.BrevoAPIKey)-4:]
	}
	log.Printf("API Key: %s (%d chars)\n", maskedKey, len(cfg.BrevoAPIKey))
	log.Printf("From: %s <%s>\n", service.fromName, service.fromEmail)
	log.Printf("Frontend URL: %s\n", service.frontendURL)
	log.Println("==========================================")

	return service
}

func (s *BrevoEmailService) SendVerificationEmail(toEmail, userName, verificationToken string) error {
	log.Printf("\n╔══════════════════════════════════════════════════════════════╗\n")
	log.Printf("║     PREPARING VERIFICATION EMAIL - Brevo Service            ║\n")
	log.Printf("╚══════════════════════════════════════════════════════════════╝\n")
	log.Printf("[PREPARE] Building verification email...\n")
	log.Printf("  ├─ To: %s\n", toEmail)
	log.Printf("  ├─ User Name: %s\n", userName)
	log.Printf("  ├─ Token: %s\n", verificationToken)
	log.Printf("  ├─ Frontend URL: %s\n", s.frontendURL)

	verificationLink := fmt.Sprintf("%s/verify-email?token=%s", s.frontendURL, verificationToken)
	log.Printf("  └─ Verification Link: %s\n", verificationLink)

	subject := "Verify Your Gymie Account"

	htmlContent := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Verify Your Email</title>
</head>
<body style="margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; background-color: #f5f5f5;">
    <table cellpadding="0" cellspacing="0" border="0" width="100%%" style="background-color: #f5f5f5; padding: 20px 0;">
        <tr>
            <td align="center">
                <table cellpadding="0" cellspacing="0" border="0" width="600" style="background-color: #ffffff; border-radius: 16px; overflow: hidden; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);">
                    <tr>
                        <td style="background: linear-gradient(135deg, #4F46E5 0%%, #7C3AED 100%%); padding: 40px 30px; text-align: center;">
                            <h1 style="margin: 0; color: #ffffff; font-size: 32px; font-weight: 700;">🏋️ Gymie</h1>
                            <p style="margin: 10px 0 0 0; color: #ffffff; font-size: 18px; opacity: 0.9;">Your Fitness Journey Starts Here</p>
                        </td>
                    </tr>
                    <tr>
                        <td style="padding: 40px 30px;">
                            <h2 style="margin: 0 0 20px 0; color: #1f2937; font-size: 24px; font-weight: 600;">Welcome to Gymie, %s! 👋</h2>
                            <p style="margin: 0 0 20px 0; color: #4b5563; font-size: 16px; line-height: 1.6;">
                                Thank you for signing up! We're excited to help you on your fitness journey.
                            </p>
                            <p style="margin: 0 0 30px 0; color: #4b5563; font-size: 16px; line-height: 1.6;">
                                To get started, please verify your email address by clicking the button below:
                            </p>
                            <table cellpadding="0" cellspacing="0" border="0" width="100%%">
                                <tr>
                                    <td align="center" style="padding: 0 0 30px 0;">
                                        <a href="%s" style="display: inline-block; background: linear-gradient(135deg, #4F46E5 0%%, #7C3AED 100%%); color: #ffffff; text-decoration: none; padding: 16px 40px; border-radius: 12px; font-size: 16px; font-weight: 600; box-shadow: 0 4px 6px rgba(79, 70, 229, 0.3);">
                                            Verify Email Address
                                        </a>
                                    </td>
                                </tr>
                            </table>
                            <div style="background-color: #fef3c7; border-left: 4px solid #f59e0b; padding: 16px; border-radius: 8px; margin: 0 0 20px 0;">
                                <p style="margin: 0; color: #92400e; font-size: 14px; font-weight: 500;">
                                    ⏰ This verification link will expire in 24 hours.
                                </p>
                            </div>
                            <p style="margin: 0; color: #6b7280; font-size: 14px; line-height: 1.6;">
                                If you didn't create an account with Gymie, you can safely ignore this email.
                            </p>
                        </td>
                    </tr>
                    <tr>
                        <td style="background-color: #f9fafb; padding: 30px; text-align: center; border-top: 1px solid #e5e7eb;">
                            <p style="margin: 0 0 10px 0; color: #6b7280; font-size: 14px;">
                                © 2024 Gymie. All rights reserved.
                            </p>
                            <p style="margin: 0; color: #9ca3af; font-size: 12px;">
                                Track your fitness journey, one workout at a time.
                            </p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>
	`, userName, verificationLink)

	plainTextContent := fmt.Sprintf(`
Welcome to Gymie, %s!

Thank you for signing up! We're excited to help you on your fitness journey.

To get started, please verify your email address by visiting this link:
%s

This verification link will expire in 24 hours.

If you didn't create an account with Gymie, you can safely ignore this email.

© 2024 Gymie. All rights reserved.
	`, userName, verificationLink)

	return s.sendEmail(toEmail, subject, plainTextContent, htmlContent)
}

func (s *BrevoEmailService) SendPasswordResetEmail(toEmail, userName, resetToken string) error {
	log.Printf("\n╔══════════════════════════════════════════════════════════════╗\n")
	log.Printf("║     PREPARING PASSWORD RESET EMAIL - Brevo Service          ║\n")
	log.Printf("╚══════════════════════════════════════════════════════════════╝\n")
	log.Printf("[PREPARE] Building password reset email...\n")
	log.Printf("  ├─ To: %s\n", toEmail)
	log.Printf("  ├─ User Name: %s\n", userName)
	log.Printf("  ├─ Token: %s\n", resetToken)
	log.Printf("  ├─ Frontend URL: %s\n", s.frontendURL)

	resetLink := fmt.Sprintf("%s/reset-password?token=%s", s.frontendURL, resetToken)
	log.Printf("  └─ Reset Link: %s\n", resetLink)

	subject := "Reset Your Gymie Password"

	htmlContent := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Reset Your Password</title>
</head>
<body style="margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; background-color: #f5f5f5;">
    <table cellpadding="0" cellspacing="0" border="0" width="100%%" style="background-color: #f5f5f5; padding: 20px 0;">
        <tr>
            <td align="center">
                <table cellpadding="0" cellspacing="0" border="0" width="600" style="background-color: #ffffff; border-radius: 16px; overflow: hidden; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);">
                    <tr>
                        <td style="background: linear-gradient(135deg, #EF4444 0%%, #DC2626 100%%); padding: 40px 30px; text-align: center;">
                            <h1 style="margin: 0; color: #ffffff; font-size: 32px; font-weight: 700;">🔒 Password Reset</h1>
                        </td>
                    </tr>
                    <tr>
                        <td style="padding: 40px 30px;">
                            <h2 style="margin: 0 0 20px 0; color: #1f2937; font-size: 24px; font-weight: 600;">Hi %s,</h2>
                            <p style="margin: 0 0 20px 0; color: #4b5563; font-size: 16px; line-height: 1.6;">
                                We received a request to reset your Gymie password. Click the button below to create a new password:
                            </p>
                            <table cellpadding="0" cellspacing="0" border="0" width="100%%">
                                <tr>
                                    <td align="center" style="padding: 0 0 30px 0;">
                                        <a href="%s" style="display: inline-block; background: linear-gradient(135deg, #EF4444 0%%, #DC2626 100%%); color: #ffffff; text-decoration: none; padding: 16px 40px; border-radius: 12px; font-size: 16px; font-weight: 600;">
                                            Reset Password
                                        </a>
                                    </td>
                                </tr>
                            </table>
                            <div style="background-color: #fef3c7; border-left: 4px solid #f59e0b; padding: 16px; border-radius: 8px; margin: 0 0 20px 0;">
                                <p style="margin: 0; color: #92400e; font-size: 14px; font-weight: 500;">
                                    ⏰ This link will expire in 1 hour.
                                </p>
                            </div>
                            <p style="margin: 0; color: #6b7280; font-size: 14px; line-height: 1.6;">
                                If you didn't request a password reset, please ignore this email or contact support if you have concerns.
                            </p>
                        </td>
                    </tr>
                    <tr>
                        <td style="background-color: #f9fafb; padding: 30px; text-align: center; border-top: 1px solid #e5e7eb;">
                            <p style="margin: 0; color: #6b7280; font-size: 14px;">
                                © 2024 Gymie. All rights reserved.
                            </p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>
	`, userName, resetLink)

	plainTextContent := fmt.Sprintf(`
Hi %s,

We received a request to reset your Gymie password.

To reset your password, visit this link:
%s

This link will expire in 1 hour.

If you didn't request a password reset, please ignore this email or contact support if you have concerns.

© 2024 Gymie. All rights reserved.
	`, userName, resetLink)

	return s.sendEmail(toEmail, subject, plainTextContent, htmlContent)
}

func (s *BrevoEmailService) sendEmail(to, subject, plainTextContent, htmlContent string) error {
	log.Printf("\n╔══════════════════════════════════════════════════════════════╗\n")
	log.Printf("║          BREVO EMAIL SENDING - DEBUG LOG                     ║\n")
	log.Printf("╚══════════════════════════════════════════════════════════════╝\n")

	log.Printf("[STEP 1] Validating Brevo configuration...\n")
	log.Printf("  ├─ API Key Length: %d characters\n", len(s.apiKey))
	log.Printf("  ├─ API Key Prefix: %s****\n", s.apiKey[:min(10, len(s.apiKey))])
	log.Printf("  ├─ From Email: %s\n", s.fromEmail)
	log.Printf("  ├─ From Name: %s\n", s.fromName)
	log.Printf("  ├─ Frontend URL: %s\n", s.frontendURL)
	log.Printf("  └─ To Email: %s\n", to)

	if s.apiKey == "" {
		log.Printf("❌ [ERROR] Brevo API key is empty!\n")
		return fmt.Errorf("brevo API key is not configured")
	}

	log.Printf("\n[STEP 2] Creating email message...\n")
	email := brevo.SendSmtpEmail{
		Sender: &brevo.SendSmtpEmailSender{
			Email: s.fromEmail,
			Name:  s.fromName,
		},
		To: []brevo.SendSmtpEmailTo{
			{
				Email: to,
				Name:  "",
			},
		},
		Subject:     subject,
		HtmlContent: htmlContent,
		TextContent: plainTextContent,
	}
	log.Printf("  ├─ From: %s <%s>\n", s.fromName, s.fromEmail)
	log.Printf("  ├─ To: %s\n", to)
	log.Printf("  ├─ Subject: %s\n", subject)
	log.Printf("  ├─ Plain text length: %d bytes\n", len(plainTextContent))
	log.Printf("  └─ HTML content length: %d bytes\n", len(htmlContent))

	log.Printf("\n[STEP 3] Sending email via Brevo API...\n")
	log.Printf("  ├─ Endpoint: https://api.brevo.com/v3/smtp/email\n")
	log.Printf("  └─ Making API request...\n")

	ctx := context.Background()
	result, response, err := s.client.TransactionalEmailsApi.SendTransacEmail(ctx, email)

	if err != nil {
		log.Printf("\n❌ [STEP 4] BREVO API ERROR\n")
		log.Printf("  ├─ Error Type: %T\n", err)
		log.Printf("  ├─ Error Message: %v\n", err)
		if response != nil {
			log.Printf("  ├─ HTTP Status: %d\n", response.StatusCode)
		}
		log.Printf("  └─ This usually means:\n")
		log.Printf("      • Network connectivity issues\n")
		log.Printf("      • Invalid API key\n")
		log.Printf("      • Brevo service down\n")
		log.Printf("      • Sender email not verified\n")
		return fmt.Errorf("brevo API request failed: %w", err)
	}

	log.Printf("\n[STEP 5] Processing Brevo response...\n")
	log.Printf("  ├─ Status Code: %d\n", response.StatusCode)
	if result.MessageId != "" {
		log.Printf("  ├─ Message ID: %s\n", result.MessageId)
	}

	if response.StatusCode >= 400 {
		log.Printf("\n❌ [ERROR] Brevo returned error status\n")
		log.Printf("  ├─ Status Code: %d\n", response.StatusCode)
		log.Printf("  └─ Common causes:\n")

		switch response.StatusCode {
		case 400:
			log.Printf("      • Invalid request (check email format)\n")
		case 401:
			log.Printf("      • Invalid API key\n")
		case 402:
			log.Printf("      • Account needs payment or credits\n")
		case 403:
			log.Printf("      • Forbidden (check sender verification)\n")
		case 429:
			log.Printf("      • Rate limit exceeded\n")
		case 500, 502, 503:
			log.Printf("      • Brevo server error (try again later)\n")
		}

		return fmt.Errorf("brevo error (status %d)", response.StatusCode)
	}

	log.Printf("\n✅ [SUCCESS] Email sent via Brevo!\n")
	log.Printf("  ├─ Status Code: %d\n", response.StatusCode)
	log.Printf("  ├─ Message ID: %s\n", result.MessageId)
	log.Printf("  └─ Email delivered to Brevo successfully\n")
	log.Printf("\n═══════════════════════════════════════════════════════════════\n\n")

	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
