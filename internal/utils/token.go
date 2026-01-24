package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// GenerateSecureToken generates a cryptographically secure random token
func GenerateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate secure token: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateVerificationToken generates a 32-byte (64 character) verification token
func GenerateVerificationToken() (string, error) {
	return GenerateSecureToken(32)
}

// GeneratePasswordResetToken generates a 32-byte (64 character) password reset token
func GeneratePasswordResetToken() (string, error) {
	return GenerateSecureToken(32)
}
