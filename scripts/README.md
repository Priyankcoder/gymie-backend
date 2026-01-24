
# Database Scripts

## Delete User Script

Completely removes a user and all their associated data from the database.

### Usage

#### Option 1: Using Make (Recommended)
```bash
cd backend
make delete-user EMAIL=user@example.com
```

#### Option 2: Direct Script
```bash
cd backend
./scripts/delete-user.sh user@example.com
```

### What Gets Deleted

The script performs a cascading delete of all user-related data:

1. **Workout Data**
   - Workout sets
   - Exercises
   - Workouts

2. **Nutrition Data**
   - Foods
   - Meals
   - Nutrition days

3. **Progress Tracking**
   - Progress photos
   - Weight entries

4. **Workout Plans**
   - Scheduled workouts
   - Workout plan days
   - Workout plans

5. **Offline Nutrition**
   - User corrections

6. **Profile & User**
   - User profile
   - User account

### Examples

Delete a test user:
```bash
make delete-user EMAIL=test@example.com
```

Delete multiple users (one at a time):
```bash
make delete-user EMAIL=user1@example.com
make delete-user EMAIL=user2@example.com
make delete-user EMAIL=user3@example.com
```

### Error Handling

If the user doesn't exist, the script will notify you and exit gracefully:
```
User with email user@example.com not found
```

### Transaction Safety

The script runs all deletions in a single database transaction:
- If any error occurs, all changes are rolled back
- Database remains in a consistent state

### Requirements

- Docker and docker-compose must be running
- PostgreSQL container must be accessible
- Database must be `gymie_dev`

### Notes

- This operation is **irreversible**
- All user data will be permanently deleted
- Use with caution in production environments
- Consider backing up data before deletion if needed
