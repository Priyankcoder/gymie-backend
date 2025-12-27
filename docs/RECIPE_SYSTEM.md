
# Recipe System Implementation Guide

## Overview

This system allows you to:
1. **Generate recipes using AI** (one-time cost)
2. **Store in database** (free forever after)
3. **Serve to users** (no ongoing API costs)

## Cost Analysis

### AI Generation (One-Time)
- **OpenAI GPT-4o-mini**: ~$0.0001 per recipe
- **Claude Haiku**: ~$0.0003 per recipe  
- **Generate 500 recipes**: ~$0.15 total (one-time cost!)

### Storage (Ongoing)
- **PostgreSQL**: Free (self-hosted) or ~$5/month (managed)
- **No ongoing API costs** ✅

## Database Schema

Already created in `backend/internal/models/recipe.go`:
- Title, description, image
- Nutrition: calories, protein, carbs, fat
- Instructions & ingredients
- Tags, category, difficulty
- AI generation flag

## Implementation Steps

### Step 1: Set up AI API Key

Add to your `.env`:
```bash
# Choose ONE:
OPENAI_API_KEY=sk-...
# OR
ANTHROPIC_API_KEY=sk-ant-...
```

### Step 2: Run Migration

Add to `backend/internal/config/database.go`:
```go
db.AutoMigrate(
    // ... existing models
    &models.Recipe{},
)
```

### Step 3: Generate Seed Recipes

Create `backend/cmd/seed-recipes/main.go`:
```go
package main

import (
    "context"
    "fmt"
    "log"
    
    // Your imports
)

func main() {
    // Connect to DB
    // Generate 100 recipes
    // Categories: breakfast (25), lunch (25), dinner (25), snacks (25)
    // Focus: high-protein, balanced, low-carb, etc.
}
```

### Step 4: Run Seed Script

```bash
cd backend
go run cmd/seed-recipes/main.go
```

This will:
- Generate 100 recipes using AI
- Cost: ~$0.01 total
- Save to database
- Never need to generate again!

### Step 5: Add API Endpoints

Routes to add:
- `GET /recipes` - List recipes (with filters)
- `GET /recipes/:id` - Get single recipe
- `GET /recipes/search` - Search recipes
- `POST /recipes/generate` - Generate new recipe (admin only)

## Recipe Categories to Generate

### Breakfast (25 recipes)
- High-protein options
- Quick prep (<15 min)
- Balanced macros
- Various cuisines

### Lunch (25 recipes)
- Meal prep friendly
- 400-600 calories
- High protein (>30g)
- Portable options

### Dinner (25 recipes)
- 500-700 calories
- Balanced nutrition
- Family-friendly
- Various difficulty levels

### Snacks (25 recipes)
- Under 300 calories
- High protein (>15g)
- Quick prep
- Portable

## Example AI Prompt for Generation

```
Generate a fitness-focused recipe with these requirements:

Category: {breakfast/lunch/dinner/snack}
Target Calories: {400-600}
Target Protein: {30-50g}
Diet Type: {balanced/high-protein/low-carb}
Difficulty: {easy/medium}
Max Prep Time: {15-30 minutes}

Return JSON with:
{
  "title": "Recipe Name",
  "description": "Brief description",
  "calories": 450,
  "protein": 35,
  "carbs": 40,
  "fat": 15,
  "prepTime": 10,
  "cookTime": 15,
  "servings": 1,
  "difficulty": "easy",
  "ingredients": ["ingredient 1", "ingredient 2"],
  "instructions": ["step 1", "step 2"],
  "tags": ["high-protein", "quick"]
}
```

## Future Enhancements

1. **User-Generated Recipes**: Allow users to contribute
2. **Favorites**: Let users save favorites
3. **Meal Plans**: Generate weekly meal plans
4. **Shopping Lists**: Auto-generate from recipes
5. **AI Meal Planner**: Match recipes to macro goals

## Cost Comparison

### Option 1: External API (Spoonacular)
- Free tier: 150 calls/day
- Paid: $49/month
- **Cost for 10,000 users**: $49-$199/month

### Option 2: AI + Database (This Approach)
- Initial: $0.10 (500 recipes)
- Ongoing: $0/month
- **Cost for 10,000 users**: $0/month ✅

## Next Steps

1. Run migration
2. Set up AI API key
3. Create seed script
4. Generate 100-500 recipes
5. Add API endpoints
6. Integrate in frontend

Want me to implement any of these steps?
