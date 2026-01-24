#!/bin/bash

# Email Testing Script Runner
echo "🧪 Gymie Email Testing Script"
echo "=============================="
echo ""

# Check if we're in the backend directory
if [ ! -f "go.mod" ]; then
    echo "❌ Error: Please run this script from the backend directory"
    echo "   cd backend && ./scripts/test-email.sh"
    exit 1
fi

# Check if .env exists
if [ ! -f ".env" ]; then
    echo "❌ Error: .env file not found"
    echo "   Please create .env file with SMTP credentials"
    exit 1
fi

# Run the test script
echo "Running email test..."
echo ""

cd scripts
go run test-email.go

echo ""
echo "Test complete!"
