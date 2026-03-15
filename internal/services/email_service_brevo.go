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
		log.Fatal("BREVO_API_KEY environment variable is required")
	}

	if len(cfg.BrevoAPIKey) < 20 {
		log.Fatalf("Invalid BREVO_API_KEY - key appears to be incomplete (%d chars)", len(cfg.BrevoAPIKey))
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

	log.Printf("Brevo email service initialized (from: %s <%s>)\n", service.fromName, service.fromEmail)

	return service
}

func (s *BrevoEmailService) SendVerificationEmail(toEmail, userName, verificationToken string) error {
	log.Printf("Sending verification email to %s\n", toEmail)

	verificationLink := fmt.Sprintf("%s/verify-email?token=%s", s.frontendURL, verificationToken)

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
	log.Printf("Sending password reset email to %s\n", toEmail)

	resetLink := fmt.Sprintf("%s/reset-password?token=%s", s.frontendURL, resetToken)

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
	if s.apiKey == "" {
		return fmt.Errorf("brevo API key is not configured")
	}

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

	ctx := context.Background()
	result, response, err := s.client.TransactionalEmailsApi.SendTransacEmail(ctx, email)

	if err != nil {
		statusCode := 0
		if response != nil {
			statusCode = response.StatusCode
		}
		log.Printf("Brevo API error sending to %s: %v (status %d)\n", to, err, statusCode)
		return fmt.Errorf("brevo API request failed: %w", err)
	}

	if response.StatusCode >= 400 {
		log.Printf("Brevo returned error status %d for email to %s\n", response.StatusCode, to)
		return fmt.Errorf("brevo error (status %d)", response.StatusCode)
	}

	log.Printf("Email sent via Brevo to %s (messageId: %s)\n", to, result.MessageId)
	return nil
}
