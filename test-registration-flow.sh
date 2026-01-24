
#!/bin/bash

# Test Registration Flow Script
# This script helps test the complete registration and email verification flow

set -e  # Exit on error

echo "🏋️ Gymie Registration Flow Test"
echo "================================"
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if backend server is running
if ! lsof -Pi :8080 -sTCP:LISTEN -t >/dev/null 2>&1; then
    echo -e "${RED}❌ Backend server is not running on port 8080${NC}"
    echo ""
    echo "Please start the backend server first:"
    echo "  cd backend"
    echo "  go run cmd/api/main.go"
    echo ""
    exit 1
fi

echo -e "${GREEN}✓${NC} Backend server is running"
echo ""

# Ask for test email
read -p "Enter test email to delete (e.g., test@example.com): " TEST_EMAIL

if [ -z "$TEST_EMAIL" ]; then
    echo -e "${RED}❌ No email provided${NC}"
    exit 1
fi

echo ""
echo "🗑️  Deleting user with email: $TEST_EMAIL"
echo "================================"

# Delete the user
cd backend
make delete-user EMAIL="$TEST_EMAIL"

echo ""
echo -e "${GREEN}✓${NC} User deleted successfully"
echo ""
echo -e "${YELLOW}📝 Now register with the app using email: $TEST_EMAIL${NC}"
echo ""
echo "What to look for in the backend logs:"
echo "  1. === EMAIL SERVICE INITIALIZED === (on server start)"
echo "  2. === EMAIL DEBUG === (during registration)"
echo "  3. === SENDING EMAIL === (when email is sent)"
echo "  4. ✅ Verification email sent successfully"
echo ""
echo "If you see ❌ instead, check the error message in the logs."
