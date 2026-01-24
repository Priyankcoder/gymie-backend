
#!/bin/bash

# Script to completely delete a user and all related data from the database
# Usage: ./scripts/delete-user.sh <email>

set -e  # Exit on error

# Check if email is provided
if [ -z "$1" ]; then
    echo "❌ Error: Email address is required"
    echo "Usage: ./scripts/delete-user.sh <email>"
    exit 1
fi

EMAIL="$1"

echo "🗑️  Deleting user: $EMAIL"
echo "================================"

# Run the deletion with proper cascading
docker-compose exec -T postgres psql -U postgres -d gymie_dev << EOF
-- Start transaction
BEGIN;

-- Get user ID
DO \$\$
DECLARE
    v_user_id INTEGER;
BEGIN
    SELECT id INTO v_user_id FROM users WHERE email = '$EMAIL';
    
    IF v_user_id IS NULL THEN
        RAISE NOTICE 'User with email $EMAIL not found';
    ELSE
        RAISE NOTICE 'Found user ID: %', v_user_id;
        
        -- Delete workout-related data
        DELETE FROM workout_sets WHERE exercise_id IN (
            SELECT id FROM exercises WHERE workout_id IN (
                SELECT id FROM workouts WHERE user_id = v_user_id
            )
        );
        RAISE NOTICE 'Deleted workout_sets';
        
        DELETE FROM exercises WHERE workout_id IN (
            SELECT id FROM workouts WHERE user_id = v_user_id
        );
        RAISE NOTICE 'Deleted exercises';
        
        -- Delete nutrition-related data
        DELETE FROM foods WHERE meal_id IN (
            SELECT id FROM meals WHERE nutrition_day_id IN (
                SELECT id FROM nutrition_days WHERE user_id = v_user_id
            )
        );
        RAISE NOTICE 'Deleted foods';
        
        DELETE FROM meals WHERE nutrition_day_id IN (
            SELECT id FROM nutrition_days WHERE user_id = v_user_id
        );
        RAISE NOTICE 'Deleted meals';
        
        DELETE FROM nutrition_days WHERE user_id = v_user_id;
        RAISE NOTICE 'Deleted nutrition_days';
        
        -- Delete progress tracking data
        DELETE FROM progress_photos WHERE user_id = v_user_id;
        RAISE NOTICE 'Deleted progress_photos';
        
        DELETE FROM weight_entries WHERE user_id = v_user_id;
        RAISE NOTICE 'Deleted weight_entries';
        
        -- Delete workout plan data
        DELETE FROM scheduled_workouts WHERE user_id = v_user_id;
        RAISE NOTICE 'Deleted scheduled_workouts';
        
        DELETE FROM workout_plan_days WHERE plan_id IN (
            SELECT id FROM workout_plans WHERE user_id = v_user_id
        );
        RAISE NOTICE 'Deleted workout_plan_days';
        
        DELETE FROM workout_plans WHERE user_id = v_user_id;
        RAISE NOTICE 'Deleted workout_plans';
        
        DELETE FROM workouts WHERE user_id = v_user_id;
        RAISE NOTICE 'Deleted workouts';
        
        -- Delete offline nutrition corrections
        DELETE FROM user_corrections WHERE user_id = v_user_id;
        RAISE NOTICE 'Deleted user_corrections';
        
        -- Delete user profile
        DELETE FROM user_profiles WHERE user_id = v_user_id;
        RAISE NOTICE 'Deleted user_profile';
        
        -- Finally delete the user
        DELETE FROM users WHERE id = v_user_id;
        RAISE NOTICE 'Deleted user';
        
        RAISE NOTICE 'Successfully deleted user $EMAIL and all related data';
    END IF;
END \$\$;

-- Commit transaction
COMMIT;

SELECT 'User deletion completed successfully!' as result;
EOF

echo ""
echo "✅ User $EMAIL has been completely deleted"
