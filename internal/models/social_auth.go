package models

// GoogleSignInRequest represents the Google sign-in request
type GoogleSignInRequest struct {
	IDToken      string  `json:"id_token" binding:"required"`
	Email        *string `json:"email"`
	Name         *string `json:"name"`
	ProfileImage *string `json:"profile_image"`
}

// AppleSignInRequest represents the Apple sign-in request
type AppleSignInRequest struct {
	IDToken string  `json:"id_token" binding:"required"`
	Email   *string `json:"email"`
	Name    *string `json:"name"`
}

// GoogleTokenInfo represents the validated Google token information
type GoogleTokenInfo struct {
	Sub           string `json:"sub"` // Google user ID
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
}
