package repository

import (
	"context"

	"github.com/yourusername/gymie-backend/internal/models"
	"gorm.io/gorm"
)

// RecipeRepository defines the interface for recipe data operations
type RecipeRepository interface {
	Create(ctx context.Context, recipe *models.Recipe) error
	GetByID(ctx context.Context, id uint) (*models.Recipe, error)
	List(ctx context.Context, query *models.RecipeSearchRequest) ([]models.Recipe, int64, error)
	Update(ctx context.Context, recipe *models.Recipe) error
	Delete(ctx context.Context, id uint) error
	Search(ctx context.Context, query *models.RecipeSearchRequest) ([]models.Recipe, int64, error)
}

type recipeRepository struct {
	db *gorm.DB
}

// NewRecipeRepository creates a new recipe repository
func NewRecipeRepository(db *gorm.DB) RecipeRepository {
	return &recipeRepository{db: db}
}

// Create creates a new recipe
func (r *recipeRepository) Create(ctx context.Context, recipe *models.Recipe) error {
	return r.db.WithContext(ctx).Create(recipe).Error
}

// GetByID retrieves a recipe by ID
func (r *recipeRepository) GetByID(ctx context.Context, id uint) (*models.Recipe, error) {
	var recipe models.Recipe
	if err := r.db.WithContext(ctx).First(&recipe, id).Error; err != nil {
		return nil, err
	}
	return &recipe, nil
}

// List retrieves recipes with pagination
func (r *recipeRepository) List(ctx context.Context, query *models.RecipeSearchRequest) ([]models.Recipe, int64, error) {
	query.SetDefaults()

	var recipes []models.Recipe
	var total int64

	db := r.db.WithContext(ctx).Model(&models.Recipe{})

	// Apply filters
	if query.Category != "" {
		db = db.Where("category = ?", query.Category)
	}
	if query.MaxCalories > 0 {
		db = db.Where("calories <= ?", query.MaxCalories)
	}
	if query.MinProtein > 0 {
		db = db.Where("protein >= ?", query.MinProtein)
	}
	if query.MaxPrepTime > 0 {
		db = db.Where("prep_time + cook_time <= ?", query.MaxPrepTime)
	}
	if query.Difficulty != "" {
		db = db.Where("difficulty = ?", query.Difficulty)
	}

	// Count total
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := db.Order("created_at DESC").
		Offset(query.GetOffset()).
		Limit(query.PageSize).
		Find(&recipes).Error; err != nil {
		return nil, 0, err
	}

	return recipes, total, nil
}

// Search performs text search on recipes
func (r *recipeRepository) Search(ctx context.Context, query *models.RecipeSearchRequest) ([]models.Recipe, int64, error) {
	query.SetDefaults()

	var recipes []models.Recipe
	var total int64

	db := r.db.WithContext(ctx).Model(&models.Recipe{})

	// Text search
	if query.Query != "" {
		searchPattern := "%" + query.Query + "%"
		db = db.Where("title LIKE ? OR description LIKE ?", searchPattern, searchPattern)
	}

	// Apply other filters
	if query.Category != "" {
		db = db.Where("category = ?", query.Category)
	}
	if query.MaxCalories > 0 {
		db = db.Where("calories <= ?", query.MaxCalories)
	}
	if query.MinProtein > 0 {
		db = db.Where("protein >= ?", query.MinProtein)
	}

	// Count total
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get results
	if err := db.Order("created_at DESC").
		Offset(query.GetOffset()).
		Limit(query.PageSize).
		Find(&recipes).Error; err != nil {
		return nil, 0, err
	}

	return recipes, total, nil
}

// Update updates a recipe
func (r *recipeRepository) Update(ctx context.Context, recipe *models.Recipe) error {
	return r.db.WithContext(ctx).Save(recipe).Error
}

// Delete soft deletes a recipe
func (r *recipeRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Recipe{}, id).Error
}
