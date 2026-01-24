
#!/bin/bash

# Restart Backend Server Script
# This script stops the current backend and restarts it with fresh environment variables

set -e  # Exit on error

echo "🏋️ Gymie Backend Restart"
echo "======================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Check if backend is running
if lsof -Pi :8080 -sTCP:LISTEN -t >/dev/null 2>&1; then
    PID=$(lsof -Pi :8080 -sTCP:LISTEN -t)
    echo -e "${YELLOW}📍 Backend is currently running (PID: $PID)${NC}"
    echo "Stopping backend server..."
    kill -15 $PID 2>/dev/null || true
    
    # Wait for process to stop
    sleep 2
    
    # Force kill if still running
    if lsof -Pi :8080 -sTCP:LISTEN -t >/dev/null 2>&1; then
        echo -e "${YELLOW}⚠️  Force stopping backend...${NC}"
        kill -9 $PID 2>/dev/null || true
        sleep 1
    fi
    
    echo -e "${GREEN}✓${NC} Backend stopped"
else
    echo -e "${BLUE}ℹ️  Backend is not currently running${NC}"
fi

echo ""
echo -e "${GREEN}🚀 Starting backend server...${NC}"
echo "================================"
echo ""
echo -e "${YELLOW}📝 Watch for these initialization logs:${NC}"
echo "  1. ${BLUE}=== EMAIL SERVICE INITIALIZED ===${NC}"
echo "  2. ${BLUE}SMTP Host: smtp.gmail.com${NC}"
echo "  3. ${BLUE}SMTP Username: priyankrastogi14@gmail.com${NC}"
echo ""
echo -e "${YELLOW}Press Ctrl+C to stop the server${NC}"
echo ""

# Start backend
cd backend
go run cmd/api/main.go
