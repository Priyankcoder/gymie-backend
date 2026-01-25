package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	brevo "github.com/getbrevo/brevo-go/lib"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("\n╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║          BREVO EMAIL TESTING SCRIPT                          ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝\n")

	// Load .env file if it exists (try multiple locations)
	err := godotenv.Load(".env")
	if err != nil {
		err = godotenv.Load("../.env")
	}
	if err != nil {
		log.Println("⚠️  No .env file found, using system environment variables")
	} else {
		log.Println("✅ Loaded .env file")
	}

	// Step 1: Load and validate configuration
	fmt.Println("📋 STEP 1: Loading Configuration")
	fmt.Println("─────────────────────────────────────")

	apiKey := os.Getenv("BREVO_API_KEY")
	fromEmail := os.Getenv("FROM_EMAIL")
	fromName := os.Getenv("FROM_NAME")
	testEmail := os.Getenv("TEST_EMAIL")

	if testEmail == "" {
		fmt.Println("⚠️  TEST_EMAIL not set, using FROM_EMAIL as test recipient")
		testEmail = fromEmail
	}

	// Validate configuration
	fmt.Printf("API Key: ")
	if apiKey == "" {
		fmt.Println("❌ NOT SET")
		fmt.Println("\n❌ ERROR: BREVO_API_KEY environment variable is required")
		fmt.Println("Set it in backend/.env or export it:")
		fmt.Println("  export BREVO_API_KEY='xkeysib-xxxxx'")
		os.Exit(1)
	}
	if len(apiKey) < 20 {
		fmt.Printf("❌ INVALID (too short: %d chars)\n", len(apiKey))
		fmt.Println("Expected format: xkeysib-xxxxx... (typically 64+ characters)")
		os.Exit(1)
	}
	fmt.Printf("✅ %s****%s (%d chars)\n", apiKey[:10], apiKey[len(apiKey)-4:], len(apiKey))

	fmt.Printf("From Email: ")
	if fromEmail == "" {
		fmt.Println("❌ NOT SET")
		os.Exit(1)
	}
	fmt.Printf("✅ %s\n", fromEmail)

	fmt.Printf("From Name: ")
	if fromName == "" {
		fmt.Println("⚠️  NOT SET (will use 'Gymie')")
		fromName = "Gymie"
	}
	fmt.Printf("✅ %s\n", fromName)

	fmt.Printf("Test Recipient: ")
	if testEmail == "" {
		fmt.Println("❌ NOT SET")
		os.Exit(1)
	}
	fmt.Printf("✅ %s\n", testEmail)

	// Check for common configuration mistakes
	fmt.Println("\n🔍 STEP 2: Configuration Validation")
	fmt.Println("─────────────────────────────────────")

	// Check if FROM_NAME contains email (common mistake)
	if strings.Contains(fromName, "<") || strings.Contains(fromName, "@") {
		fmt.Printf("⚠️  WARNING: FROM_NAME contains email address: '%s'\n", fromName)
		fmt.Println("   FROM_NAME should only contain the name, not the email")
		fmt.Println("   Example: 'Gymie' not 'Gymie <no-reply@example.com>'")
		fmt.Println("\n   Would you like to continue anyway? (y/n): ")
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("\n❌ Aborted. Please fix FROM_NAME and try again.")
			os.Exit(1)
		}
	}

	// Check if API key starts with xkeysib-
	if !strings.HasPrefix(apiKey, "xkeysib-") {
		fmt.Println("⚠️  WARNING: API key doesn't start with 'xkeysib-'")
		fmt.Println("   Are you sure this is a valid Brevo API key?")
	}

	fmt.Println("✅ Configuration looks good!")

	// Step 3: Initialize Brevo client
	fmt.Println("\n🔧 STEP 3: Initializing Brevo Client")
	fmt.Println("─────────────────────────────────────")

	cfg := brevo.NewConfiguration()
	cfg.AddDefaultHeader("api-key", apiKey)
	client := brevo.NewAPIClient(cfg)

	fmt.Println("✅ Brevo API client initialized")

	// Step 4: Create test email
	fmt.Println("\n📧 STEP 4: Creating Test Email")
	fmt.Println("─────────────────────────────────────")

	subject := "Brevo Test Email - Gymie"

	plainTextContent := fmt.Sprintf(`
Hello!

This is a test email from your Gymie backend to verify Brevo is working correctly.

Configuration Details:
- From Name: %s
- From Email: %s
- Brevo API Key: %s****%s
- Recipient: %s

If you received this email, Brevo is configured correctly!

---
Gymie Backend
Automated Test Email
	`, fromName, fromEmail, apiKey[:10], apiKey[len(apiKey)-4:], testEmail)

	htmlContent := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Brevo Test Email</title>
</head>
<body style="margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; background-color: #f5f5f5;">
    <table cellpadding="0" cellspacing="0" border="0" width="100%%" style="background-color: #f5f5f5; padding: 20px 0;">
        <tr>
            <td align="center">
                <table cellpadding="0" cellspacing="0" border="0" width="600" style="background-color: #ffffff; border-radius: 16px; overflow: hidden; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);">
                    <tr>
                        <td style="background: linear-gradient(135deg, #4F46E5 0%%, #7C3AED 100%%); padding: 40px 30px; text-align: center;">
                            <h1 style="margin: 0; color: #ffffff; font-size: 32px; font-weight: 700;">🏋️ Gymie</h1>
                            <p style="margin: 10px 0 0 0; color: #ffffff; font-size: 18px; opacity: 0.9;">Brevo Email Test</p>
                        </td>
                    </tr>
                    
                    <tr>
                        <td style="padding: 40px 30px;">
                            <h2 style="margin: 0 0 20px 0; color: #1f2937; font-size: 24px; font-weight: 600;">Test Email Successful! ✅</h2>
                            
                            <p style="margin: 0 0 20px 0; color: #4b5563; font-size: 16px; line-height: 1.6;">
                                This is a test email from your Gymie backend to verify Brevo is working correctly.
                            </p>
                            
                            <div style="background-color: #f0fdf4; border-left: 4px solid #10b981; padding: 16px; border-radius: 8px; margin: 0 0 20px 0;">
                                <p style="margin: 0; color: #065f46; font-size: 14px; font-weight: 500;">
                                    ✅ If you received this email, Brevo is configured correctly!
                                </p>
                            </div>
                            
                            <p style="margin: 0 0 10px 0; color: #6b7280; font-size: 14px; line-height: 1.6;">
                                <strong>Configuration Details:</strong>
                            </p>
                            
                            <ul style="margin: 0 0 20px 0; color: #6b7280; font-size: 14px; line-height: 1.8;">
                                <li>From Name: %s</li>
                                <li>From Email: %s</li>
                                <li>API Key: %s****%s</li>
                                <li>Recipient: %s</li>
                            </ul>
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
	`, fromName, fromEmail, apiKey[:10], apiKey[len(apiKey)-4:], testEmail)

	// Create the email structure
	email := brevo.SendSmtpEmail{
		Sender: &brevo.SendSmtpEmailSender{
			Email: fromEmail,
			Name:  fromName,
		},
		To: []brevo.SendSmtpEmailTo{
			{
				Email: testEmail,
				Name:  "Test Recipient",
			},
		},
		Subject:     subject,
		HtmlContent: htmlContent,
		TextContent: plainTextContent,
	}

	fmt.Printf("Subject: %s\n", subject)
	fmt.Printf("From: %s <%s>\n", fromName, fromEmail)
	fmt.Printf("To: %s\n", testEmail)
	fmt.Printf("Plain text: %d bytes\n", len(plainTextContent))
	fmt.Printf("HTML content: %d bytes\n", len(htmlContent))

	// Step 5: Send email
	fmt.Println("\n📤 STEP 5: Sending Email via Brevo")
	fmt.Println("─────────────────────────────────────")
	fmt.Println("Making API request to Brevo...")

	ctx := context.Background()
	result, response, err := client.TransactionalEmailsApi.SendTransacEmail(ctx, email)

	if err != nil {
		fmt.Println("\n❌ ERROR: Failed to send email")
		fmt.Printf("Error Type: %T\n", err)
		fmt.Printf("Error Message: %v\n", err)
		if response != nil {
			fmt.Printf("HTTP Status Code: %d\n", response.StatusCode)
		}
		fmt.Println("\nCommon causes:")
		fmt.Println("  • Network connectivity issues")
		fmt.Println("  • Invalid API key")
		fmt.Println("  • Brevo service down")
		fmt.Println("  • Firewall blocking port 443")
		fmt.Println("  • Sender email not verified")
		os.Exit(1)
	}

	// Step 6: Check response
	fmt.Println("\n📥 STEP 6: Processing Response")
	fmt.Println("─────────────────────────────────────")
	fmt.Printf("HTTP Status Code: %d\n", response.StatusCode)
	if result.MessageId != "" {
		fmt.Printf("Message ID: %s\n", result.MessageId)
	}

	// Interpret status code
	fmt.Println("\n📊 STEP 7: Result")
	fmt.Println("─────────────────────────────────────")

	switch {
	case response.StatusCode == 201:
		fmt.Println("✅ SUCCESS! Email sent successfully!")
		fmt.Println("\nWhat this means:")
		fmt.Println("  • Brevo accepted the email")
		fmt.Println("  • Email is being processed")
		fmt.Println("  • You should receive it shortly")
		fmt.Printf("  • Message ID: %s\n", result.MessageId)
		fmt.Println("\nNext steps:")
		fmt.Printf("  1. Check your inbox: %s\n", testEmail)
		fmt.Println("  2. Check spam folder if not in inbox")
		fmt.Println("  3. Check Brevo logs:")
		fmt.Println("     https://app.brevo.com/log")
		fmt.Println("\n✅ Brevo is configured correctly!")

	case response.StatusCode == 400:
		fmt.Println("❌ ERROR: Bad Request (400)")
		fmt.Println("\nCommon causes:")
		fmt.Println("  • Invalid email format")
		fmt.Println("  • Missing required fields")
		fmt.Println("  • Malformed request")
		os.Exit(1)

	case response.StatusCode == 401:
		fmt.Println("❌ ERROR: Unauthorized (401)")
		fmt.Println("\nThis means:")
		fmt.Println("  • Invalid API key")
		fmt.Println("\nTo fix:")
		fmt.Println("  1. Verify API key is correct")
		fmt.Println("  2. Check API key at:")
		fmt.Println("     https://app.brevo.com/settings/keys/api")
		os.Exit(1)

	case response.StatusCode == 402:
		fmt.Println("❌ ERROR: Payment Required (402)")
		fmt.Println("\nThis means:")
		fmt.Println("  • Account needs payment or credits")
		fmt.Println("  • Daily/monthly quota exceeded")
		fmt.Println("\nTo fix:")
		fmt.Println("  1. Check account credits at: https://app.brevo.com/account/plan")
		fmt.Println("  2. Add credits or upgrade plan")
		os.Exit(1)

	case response.StatusCode == 403:
		fmt.Println("❌ ERROR: Forbidden (403)")
		fmt.Println("\nCommon causes:")
		fmt.Println("  • Sender email not verified")
		fmt.Println("  • Account suspended")
		fmt.Println("  • API key doesn't have permission to send emails")
		fmt.Println("\nTo fix:")
		fmt.Println("  1. Verify sender at: https://app.brevo.com/senders")
		fmt.Println("  2. Check API key permissions")
		fmt.Println("  3. Check account status")
		os.Exit(1)

	case response.StatusCode == 429:
		fmt.Println("❌ ERROR: Rate Limit Exceeded (429)")
		fmt.Println("\nYou've sent too many emails. Wait and try again later.")
		os.Exit(1)

	case response.StatusCode >= 500:
		fmt.Println("❌ ERROR: Brevo Server Error")
		fmt.Printf("Status Code: %d\n", response.StatusCode)
		fmt.Println("\nBrevo is having issues. Try again later.")
		os.Exit(1)

	default:
		fmt.Printf("⚠️  Unexpected status code: %d\n", response.StatusCode)
		os.Exit(1)
	}

	fmt.Println("\n╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    TEST COMPLETED                            ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝\n")
}
