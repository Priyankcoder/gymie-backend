
#!/bin/bash

echo "==================================="
echo "Testing Gymie Backend API"
echo "==================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

BASE_URL="http://localhost:8080"

echo -e "${BLUE}1. Testing Health Endpoint${NC}"
curl -s $BASE_URL/health | jq '.'
echo -e "\n"

echo -e "${BLUE}2. Registering a new user${NC}"
REGISTER_RESPONSE=$(curl -s -X POST $BASE_URL/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@gymie.com",
    "password": "password123",
    "name": "Test User"
  }')
echo "$REGISTER_RESPONSE" | jq '.'
TOKEN=$(echo "$REGISTER_RESPONSE" | jq -r '.data.token')
echo -e "${GREEN}Token saved: ${TOKEN:0:20}...${NC}"
echo -e "\n"

echo -e "${BLUE}3. Login with the user${NC}"
LOGIN_RESPONSE=$(curl -s -X POST $BASE_URL/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@gymie.com",
    "password": "password123"
  }')
echo "$LOGIN_RESPONSE" | jq '.'
echo -e "\n"

echo -e "${BLUE}4. Get user profile${NC}"
curl -s $BASE_URL/v1/users/profile \
  -H "Authorization: Bearer $TOKEN" | jq '.'
echo -e "\n"

echo -e "${BLUE}5. Update user profile${NC}"
curl -s -X PUT $BASE_URL/v1/users/profile \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "height": 175,
    "weight": 75,
    "age": 25,
    "gender": "male",
    "goal": "gain_muscle"
  }' | jq '.'
echo -e "\n"

echo -e "${BLUE}6. Create a workout${NC}"
WORKOUT_RESPONSE=$(curl -s -X POST $BASE_URL/v1/workouts \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Morning Workout",
    "date": "2024-12-18T08:00:00Z",
    "duration": 60,
    "notes": "Great workout!",
    "exercises": [
      {
        "name": "Bench Press",
        "muscle_group": "chest",
        "order": 1,
        "sets": [
          {
            "set_number": 1,
            "weight": 80,
            "reps": 10,
            "completed": true
          },
          {
            "set_number": 2,
            "weight": 80,
            "reps": 8,
            "completed": true
          }
        ]
      },
      {
        "name": "Squats",
        "muscle_group": "legs",
        "order": 2,
        "sets": [
          {
            "set_number": 1,
            "weight": 100,
            "reps": 12,
            "completed": true
          }
        ]
      }
    ]
  }')
echo "$WORKOUT_RESPONSE" | jq '.'
WORKOUT_ID=$(echo "$WORKOUT_RESPONSE" | jq -r '.data.id')
echo -e "${GREEN}Workout ID: $WORKOUT_ID${NC}"
echo -e "\n"

echo -e "${BLUE}7. Get all workouts${NC}"
curl -s "$BASE_URL/v1/workouts?page=1&page_size=10" \
  -H "Authorization: Bearer $TOKEN" | jq '.'
echo -e "\n"

echo -e "${BLUE}8. Get workout stats${NC}"
curl -s $BASE_URL/v1/workouts/stats \
  -H "Authorization: Bearer $TOKEN" | jq '.'
echo -e "\n"

echo -e "${BLUE}9. Create nutrition day${NC}"
NUTRITION_RESPONSE=$(curl -s -X POST $BASE_URL/v1/nutrition \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "date": "2024-12-18T00:00:00Z",
    "notes": "Healthy eating day",
    "meals": [
      {
        "name": "Breakfast",
        "foods": [
          {
            "name": "Eggs",
            "calories": 155,
            "protein": 13,
            "carbs": 1,
            "fat": 11,
            "quantity": 2,
            "unit": "pieces"
          },
          {
            "name": "Oatmeal",
            "calories": 150,
            "protein": 5,
            "carbs": 27,
            "fat": 3,
            "quantity": 1,
            "unit": "serving"
          }
        ]
      }
    ]
  }')
echo "$NUTRITION_RESPONSE" | jq '.'
echo -e "\n"

echo -e "${BLUE}10. Get nutrition stats${NC}"
curl -s $BASE_URL/v1/nutrition/stats \
  -H "Authorization: Bearer $TOKEN" | jq '.'
echo -e "\n"

echo -e "${BLUE}11. Create weight entry${NC}"
curl -s -X POST $BASE_URL/v1/progress/weight \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "date": "2024-12-18T00:00:00Z",
    "weight": 75,
    "notes": "Starting weight"
  }' | jq '.'
echo -e "\n"

echo -e "${BLUE}12. Get weight progress${NC}"
curl -s $BASE_URL/v1/progress/weight/stats \
  -H "Authorization: Bearer $TOKEN" | jq '.'
echo -e "\n"

echo -e "${GREEN}==================================="
echo "All API tests completed!"
echo "===================================${NC}"
