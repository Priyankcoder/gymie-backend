
# Google Sign-In Backend Implementation

## ✅ Implementation Complete

This document describes the backend implementation for Google Sign-In authentication in the Gymie fitness app.

## Files Created/Modified

### New Files:
1. **`backend/internal/models/social_auth.go`** - Request/response models for social authentication
2. **`backend/internal/utils/google_auth.go`** - Google token verification utility
3. **`backend/internal/utils/password.go`** - Random password generation utility

### Modified Files:
1. **`backend/internal/service/auth_service.go`** - Added `LoginWithGoogle()` method
2. **`backend/internal/handlers/auth_handler.go`** - Added `LoginWithGoogle()` handler
3. **`backend/internal/routes/routes.go`** - Added `/auth/google` route

## API Endpoint

### POST `/v1/auth/google`

Authenticates a user via Google Sign-In by verifying the Google ID token.

**Request Body:**
```json
{
  "id_token": "eyJhbGciOiJSUzI1NiIs...",
  "email": "user@gmail.com",
  "name": "John Doe",
  "profile_image": "https://..."
}
```

**Response (Success - 200 OK):**
```json
{
  "success": true,
  "message": "Google sign-in successful",
  "data": {
    "token": "jwt-token-here",
    "user": {
      "id": 1,
      "email": "user@gmail.com",
      "name": "John Doe",
      "profile": {
        "profilePicture": "https://..."
      },
      "createdAt": "2024-01-13T16:00:00Z",
      "updatedAt": "2024-01-13T16:00:00Z"
    }
  }
}
```

**Response (Error - 401 Unauthorized):**
```json
{
  "success": false,
  "code": "google_signin_failed",
  "message": "failed to verify Google token: token expired",
  "details": null
}
```

## Implementation Details

### Token Verification Flow

1. **Receive Request**: Frontend sends Google ID token to `/v1/auth/google`
2. **Verify Token**: Backend calls Google's tokeninfo API to verify the token
3. **Validate Claims**: Check email verification, expiration, and issuer
4. **Check User**: Look up user by email in database
5. **Create/Login**:
   - If user exists: Generate JWT and return
   - If new user: Create account with Google profile data, then generate JWT
6. **Return Response**: Send JWT token and user data to frontend

### Security Features

1. **Token Verification**: All Google ID tokens are verified with Google's API
2. **Email Verification**: Only verified emails are accepted
3. **Expiration Check**: Tokens are validated for expiration
4. **Issuer Validation**: Only tokens from `accounts.google.com` are accepted
5. **Secure Password**: New Google users get a random 32-character password (unused for Google login)

### Database Changes

No schema changes required. The implementation uses existing `users` and `user_profiles` tables:

- **users table**: Stores email, name, and hashed random password
- **user_profiles table**: Stores Google profile picture URL

### Auto-Registration

When a user signs in with Google for the first time:
1. A new user account is created automatically
2. Profile information from Google is saved (name, picture)
3. A random password is generated (not used for Google sign-in)
4. User is immediately authenticated and returned a JWT

## Testing

### Using cURL

```bash
# First, get a Google ID token from frontend or use a test token
# Then test the endpoint:

curl -X POST http://localhost:8080/v1/auth/google \
  -H "Content-Type: application/json" \
  -d '{
    "id_token": "YOUR_GOOGLE_ID_TOKEN_HERE",
    "email": "user@gmail.com",
    "name": "John Doe"
  }'
```

### Integration Testing

1. Start backend server: `cd backend && make run`
2. Frontend makes Google Sign-In request
3. Frontend receives ID token from Google
4. Frontend sends ID token to backend
5. Backend verifies and returns JWT
6. Frontend stores JWT and navigates to app

## Error Handling

The implementation handles these error scenarios:

1. **Invalid Token**: Returns 401 with "failed to verify Google token"
2. **Expired Token**: Returns 401 with "token expired"
3. **Unverified Email**: Returns 401 with "email not verified"
4. **Invalid Issuer**: Returns 401 with "invalid token issuer"
5. **Database Errors**: Returns appropriate error messages
6. **Network Errors**: Returns error from Google API call

## Configuration

No additional configuration required. The implementation:
- Uses Google's public tokeninfo API (no API keys needed)
- Works with any valid Google OAuth client ID configured in frontend
- Reuses existing JWT configuration from `.env`

## Dependencies

No new Go dependencies required. Uses standard library:
- `net/http` - HTTP client for Google API calls
- `encoding/json` - JSON parsing
- `context` - Request context handling

## Frontend Integration

The frontend is already configured to use this endpoint:

**File**: `frontend/src/services/realApi.ts`
```typescript
loginWithGoogle: async (idToken, email, name, profileImage) => {
  const response = await apiClient.post('/auth/google', {
    id_token: idToken,
    email,
    name,
    profile_image: profileImage
  });
  return response.data.data;
}
```

## Production Considerations

### Security
1. ✅ Token verification with Google's API
2. ✅ Email verification required
3. ✅ Issuer validation
4. ✅ Expiration check
5. ✅ Secure password generation for new users

### Performance
- Google token verification adds ~100-200ms latency
- Consider caching verified tokens (with expiration) for repeat requests
- Use connection pooling for HTTP client

### Monitoring
Log these events:
- Successful Google sign-ins
- Failed verification attempts
- New user registrations via Google
- Token verification errors

### Rate Limiting
Consider adding rate limiting for:
- Failed verification attempts per IP
- Sign-in requests per user
- Overall endpoint requests

## Next Steps

1. **Install Frontend Dependencies**:
   ```bash
   cd frontend
   npm install
   ```

2. **Build and Test**:
   ```bash
   cd backend
   go build ./...
   make run
   ```

3. **Test Integration**:
   - Run frontend app
   - Try Google Sign-In
   - Verify JWT token is stored
   - Check user is logged in

4. **Production Deployment**:
   - Deploy backend with Google Sign-In support
   - Update frontend API URL
   - Test on production environment

## Troubleshooting

### "failed to verify Google token: EOF"
- Network connectivity issue
- Google API is down
- Invalid token format

### "token expired"
- Token was generated more than 1 hour ago
- Request client to get a fresh token from Google

### "email not verified"
- User's Google account email is not verified
- Request user to verify email with Google

### "user with email X already exists"
- User registered with email/password before
- This shouldn't happen as the code handles existing users
- Check database for duplicate email issues

## Support

For questions or issues:
- Check frontend implementation in `SOCIAL_AUTH_FRONTEND_IMPLEMENTATION.md`
- Review backend logs for detailed error messages
- Verify Google OAuth configuration in frontend
- Test token verification with Google's tokeninfo API directly
