package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("🧪 Email Testing Script")
	fmt.Println("========================\n")

	// Load .env file
	fmt.Println("📁 Loading .env file...")
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("❌ Error loading .env file: %v", err)
	}
	fmt.Println("✅ .env file loaded\n")

	// Get SMTP configuration
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	fromEmail := os.Getenv("FROM_EMAIL")
	fromName := os.Getenv("FROM_NAME")

	fmt.Println("📧 SMTP Configuration:")
	fmt.Printf("   Host: %s\n", host)
	fmt.Printf("   Port: %s\n", port)
	fmt.Printf("   Username: %s\n", username)
	fmt.Printf("   Password: %s (length: %d)\n", strings.Repeat("*", len(password)), len(password))
	fmt.Printf("   From Email: %s\n", fromEmail)
	fmt.Printf("   From Name: %s\n\n", fromName)

	// Validate configuration
	if host == "" || port == "" || username == "" || password == "" || fromEmail == "" {
		log.Fatalf("❌ Missing required SMTP configuration in .env file")
	}

	// Test recipient
	toEmail := "priyankrastogi145@gmail.com"
	fmt.Printf("📬 Test recipient: %s\n\n", toEmail)

	// Test connection
	fmt.Printf("🔌 Testing connection to %s:%s...\n", host, port)
	addr := fmt.Sprintf("%s:%s", host, port)

	// Create TLS config
	tlsConfig := &tls.Config{
		ServerName: host,
	}

	// Setup authentication
	auth := smtp.PlainAuth("", username, password, host)
	fmt.Println("✅ Authentication setup complete\n")

	// Compose email
	subject := "Test Email from Gymie"
	body := "This is a test email to verify SMTP configuration."

	message := []byte(fmt.Sprintf(
		"From: %s <%s>\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/plain; charset=UTF-8\r\n"+
			"\r\n"+
			"%s\r\n",
		fromName, fromEmail, toEmail, subject, body,
	))

	fmt.Println("📤 Sending email...")
	fmt.Printf("   From: %s <%s>\n", fromName, fromEmail)
	fmt.Printf("   To: %s\n", toEmail)
	fmt.Printf("   Subject: %s\n", subject)
	fmt.Printf("   Message size: %d bytes\n\n", len(message))

	// Send email with detailed error handling
	err = sendMailWithDebug(addr, auth, fromEmail, []string{toEmail}, message, tlsConfig)
	if err != nil {
		fmt.Printf("\n❌ Failed to send email: %v\n", err)
		fmt.Println("\n🔍 Troubleshooting:")
		fmt.Println("   1. Check if SMTP_PASSWORD is an App Password (not regular password)")
		fmt.Println("   2. Verify 2FA is enabled on Gmail account")
		fmt.Println("   3. Generate new App Password at: https://myaccount.google.com/apppasswords")
		fmt.Println("   4. Check if 'Less secure app access' is enabled (if using regular password)")
		fmt.Println("   5. Try port 465 with SSL instead of 587 with TLS")
		os.Exit(1)
	}

	fmt.Println("\n✅ Email sent successfully!")
	fmt.Println("📨 Check your inbox (and spam folder) at:", toEmail)
	fmt.Println("\n💡 If email doesn't arrive in 2-3 minutes, check Gmail account settings")
}

func sendMailWithDebug(addr string, auth smtp.Auth, from string, to []string, msg []byte, tlsConfig *tls.Config) error {
	fmt.Println("🔐 Connecting to SMTP server...")

	// Connect to server
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("dial failed: %w", err)
	}
	defer client.Close()
	fmt.Println("✅ Connected to SMTP server")

	// Send HELLO
	fmt.Println("👋 Sending EHLO...")
	if err = client.Hello("localhost"); err != nil {
		return fmt.Errorf("hello failed: %w", err)
	}
	fmt.Println("✅ EHLO sent")

	// Start TLS
	fmt.Println("🔒 Starting TLS...")
	if err = client.StartTLS(tlsConfig); err != nil {
		return fmt.Errorf("starttls failed: %w", err)
	}
	fmt.Println("✅ TLS connection established")

	// Authenticate
	fmt.Println("🔑 Authenticating...")
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("authentication failed: %w (check if using App Password)", err)
	}
	fmt.Println("✅ Authentication successful")

	// Set sender
	fmt.Println("📧 Setting sender...")
	if err = client.Mail(from); err != nil {
		return fmt.Errorf("mail failed: %w", err)
	}
	fmt.Println("✅ Sender set")

	// Set recipients
	fmt.Println("👥 Setting recipients...")
	for _, addr := range to {
		if err = client.Rcpt(addr); err != nil {
			return fmt.Errorf("rcpt failed for %s: %w", addr, err)
		}
		fmt.Printf("✅ Recipient added: %s\n", addr)
	}

	// Send message
	fmt.Println("💌 Sending message data...")
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("data failed: %w", err)
	}

	_, err = w.Write(msg)
	if err != nil {
		return fmt.Errorf("write failed: %w", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("close failed: %w", err)
	}
	fmt.Println("✅ Message data sent")

	// Quit
	fmt.Println("👋 Sending QUIT...")
	err = client.Quit()
	if err != nil {
		return fmt.Errorf("quit failed: %w", err)
	}
	fmt.Println("✅ Connection closed gracefully")

	return nil
}
