
package service

// EmailServiceInterface defines the contract for email services
type EmailServiceInterface interface {
	SendVerificationEmail(toEmail, userName, verificationToken string) error
	SendPasswordResetEmail(toEmail, userName, resetToken string) error
}
