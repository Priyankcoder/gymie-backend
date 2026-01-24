package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// GoogleTokenInfo represents the validated Google token information
type GoogleTokenInfo struct {
	Sub           string      `json:"sub"` // Google user ID
	Email         string      `json:"email"`
	EmailVerified interface{} `json:"email_verified"` // Can be bool or string
	Name          string      `json:"name"`
	Picture       string      `json:"picture"`
	GivenName     string      `json:"given_name"`
	FamilyName    string      `json:"family_name"`
	Aud           string      `json:"aud"` // Client ID
	Iss           string      `json:"iss"` // Issuer
	Exp           interface{} `json:"exp"` // Expiration time (can be int64 or string)
}

// IsEmailVerified checks if email is verified (handles both bool and string)
func (g *GoogleTokenInfo) IsEmailVerified() bool {
	switch v := g.EmailVerified.(type) {
	case bool:
		return v
	case string:
		return v == "true"
	default:
		return false
	}
}

// GetExp returns the expiration time as int64 (handles both int64 and string)
func (g *GoogleTokenInfo) GetExp() int64 {
	switch v := g.Exp.(type) {
	case float64:
		return int64(v)
	case int64:
		return v
	case string:
		// Parse string to int64
		var exp int64
		fmt.Sscanf(v, "%d", &exp)
		return exp
	default:
		return 0
	}
}

// VerifyGoogleIDToken verifies a Google ID token and returns the token information
func VerifyGoogleIDToken(ctx context.Context, idToken string) (*GoogleTokenInfo, error) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Call Google's tokeninfo endpoint
	url := fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?id_token=%s", idToken)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token verification failed: %s", string(body))
	}

	// Parse response
	var tokenInfo GoogleTokenInfo
	if err := json.NewDecoder(resp.Body).Decode(&tokenInfo); err != nil {
		return nil, fmt.Errorf("failed to decode token info: %w", err)
	}

	// Validate token
	if !tokenInfo.IsEmailVerified() {
		return nil, fmt.Errorf("email not verified")
	}

	// Check if token is expired
	now := time.Now().Unix()
	if tokenInfo.GetExp() < now {
		return nil, fmt.Errorf("token expired")
	}

	// Validate issuer
	if tokenInfo.Iss != "https://accounts.google.com" && tokenInfo.Iss != "accounts.google.com" {
		return nil, fmt.Errorf("invalid token issuer")
	}

	return &tokenInfo, nil
}
