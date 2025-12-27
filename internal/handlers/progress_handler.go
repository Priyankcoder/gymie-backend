
package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/gymie-backend/internal/middleware"
	"github.com/yourusername/gymie-backend/internal/models"
	"github.com/yourusername/gymie-backend/internal/service"
)

type ProgressHandler struct {
	progressService service.ProgressService
}

func NewProgressHandler(progressService service.ProgressService) *ProgressHandler {
	return &ProgressHandler{
		progressService: progressService,
	}
}

// UploadProgressPhoto uploads a new progress photo
// @Summary Upload progress photo
// @Description Upload a new progress photo with image file
// @Tags progress
// @Accept multipart/form-data
// @Produce json
// @Security Bearer
// @Param photo formData file true "Progress photo image"
// @Param date formData string true "Date (ISO 8601 format)"
// @Param notes formData string false "Optional notes"
// @Param weight formData number false "Optional weight in kg"
// @Success 201 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /progress/photos [post]
func (h *ProgressHandler) UploadProgressPhoto(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	// Get the uploaded file
	file, err := c.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"validation_error",
			"Photo file is required",
			err.Error(),
		))
		return
	}

	// Validate file type
	if !isValidImageType(file.Header.Get("Content-Type")) {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"validation_error",
			"Invalid file type. Only JPEG, PNG, and WebP images are allowed",
			nil,
		))
		return
	}

	// Validate file size (max 10MB)
	const maxSize = 10 << 20 // 10MB
	if file.Size > maxSize {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"validation_error",
			"File size exceeds maximum limit of 10MB",
			nil,
		))
		return
	}

	// Parse form data
	dateStr := c.PostForm("date")
	notes := c.PostForm("notes")
	
	photoDate, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"validation_error",
			"Invalid date format. Use ISO 8601 format (e.g., 2024-12-26T09:00:00Z)",
			err.Error(),
		))
		return
	}

	// Optional weight
	var weight *float64
	if weightStr := c.PostForm("weight"); weightStr != "" {
		w, err := strconv.ParseFloat(weightStr, 64)
		if err == nil {
			weight = &w
		}
	}

	// Upload file and create photo record
	photo, err := h.progressService.UploadProgressPhoto(c.Request.Context(), userID, file, photoDate, weight, notes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"upload_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusCreated, models.NewSuccessResponse(
		"Progress photo uploaded successfully",
		photo,
	))
}

// CreateProgressPhoto creates a new progress photo (JSON API - for metadata only)
// @Summary Create progress photo
// @Description Create a new progress photo entry
// @Tags progress
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body models.ProgressPhotoCreateRequest true "Progress photo details"
// @Success 201 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /progress/photos [post]
func (h *ProgressHandler) CreateProgressPhoto(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var req models.ProgressPhotoCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"validation_error",
			"Invalid request body",
			err.Error(),
		))
		return
	}

	photo, err := h.progressService.CreateProgressPhoto(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"create_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusCreated, models.NewSuccessResponse(
		"Progress photo created successfully",
		photo,
	))
}

// Helper function to validate image type
func isValidImageType(contentType string) bool {
	validTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/webp": true,
	}
	return validTypes[contentType]
}

// GetProgressPhotos lists progress photos
// @Summary List progress photos
// @Description Get all progress photos for the authenticated user
// @Tags progress
// @Produce json
// @Security Bearer
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} models.PaginatedResponse
// @Router /progress/photos [get]
func (h *ProgressHandler) GetProgressPhotos(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var query models.ListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"validation_error",
			"Invalid query parameters",
			err.Error(),
		))
		return
	}

	photos, total, err := h.progressService.GetProgressPhotos(c.Request.Context(), userID, &query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"fetch_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, models.NewPaginatedResponse(photos, total, query.Page, query.PageSize))
}

// UpdateProgressPhoto updates a progress photo
// @Summary Update progress photo
// @Description Update an existing progress photo
// @Tags progress
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Progress Photo ID"
// @Param request body models.ProgressPhotoUpdateRequest true "Update details"
// @Success 200 {object} models.SuccessResponse
// @Router /progress/photos/{id} [put]
func (h *ProgressHandler) UpdateProgressPhoto(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req models.ProgressPhotoUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"validation_error",
			"Invalid request body",
			err.Error(),
		))
		return
	}

	photo, err := h.progressService.UpdateProgressPhoto(c.Request.Context(), uint(id), userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"update_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(
		"Progress photo updated successfully",
		photo,
	))
}

// DeleteProgressPhoto deletes a progress photo
// @Summary Delete progress photo
// @Description Delete a progress photo
// @Tags progress
// @Produce json
// @Security Bearer
// @Param id path int true "Progress Photo ID"
// @Success 200 {object} models.SuccessResponse
// @Router /progress/photos/{id} [delete]
func (h *ProgressHandler) DeleteProgressPhoto(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	if err := h.progressService.DeleteProgressPhoto(c.Request.Context(), uint(id), userID); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"delete_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(
		"Progress photo deleted successfully",
		nil,
	))
}

// CreateWeightEntry creates a new weight entry
// @Summary Create weight entry
// @Description Create a new weight tracking entry
// @Tags progress
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body models.WeightEntryCreateRequest true "Weight entry details"
// @Success 201 {object} models.SuccessResponse
// @Router /progress/weight [post]
func (h *ProgressHandler) CreateWeightEntry(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var req models.WeightEntryCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"validation_error",
			"Invalid request body",
			err.Error(),
		))
		return
	}

	entry, err := h.progressService.CreateWeightEntry(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"create_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusCreated, models.NewSuccessResponse(
		"Weight entry created successfully",
		entry,
	))
}

// GetWeightEntries lists weight entries
// @Summary List weight entries
// @Description Get all weight entries for the authenticated user
// @Tags progress
// @Produce json
// @Security Bearer
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} models.PaginatedResponse
// @Router /progress/weight [get]
func (h *ProgressHandler) GetWeightEntries(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var query models.ListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"validation_error",
			"Invalid query parameters",
			err.Error(),
		))
		return
	}

	var dateQuery models.DateRangeQuery
	c.ShouldBindQuery(&dateQuery)

	// Check if date range is provided
	if dateQuery.StartDate != "" && dateQuery.EndDate != "" {
		startDate, err1 := time.Parse("2006-01-02", dateQuery.StartDate)
		endDate, err2 := time.Parse("2006-01-02", dateQuery.EndDate)
		
		if err1 == nil && err2 == nil {
			entries, err := h.progressService.GetWeightEntriesByDateRange(c.Request.Context(), userID, startDate, endDate)
			if err != nil {
				c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
					"fetch_failed",
					err.Error(),
					nil,
				))
				return
			}

			c.JSON(http.StatusOK, models.NewSuccessResponse(
				"Weight entries retrieved successfully",
				entries,
			))
			return
		}
	}

	entries, total, err := h.progressService.GetWeightEntries(c.Request.Context(), userID, &query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"fetch_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, models.NewPaginatedResponse(entries, total, query.Page, query.PageSize))
}

// UpdateWeightEntry updates a weight entry
// @Summary Update weight entry
// @Description Update an existing weight entry
// @Tags progress
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Weight Entry ID"
// @Param request body models.WeightEntryUpdateRequest true "Update details"
// @Success 200 {object} models.SuccessResponse
// @Router /progress/weight/{id} [put]
func (h *ProgressHandler) UpdateWeightEntry(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req models.WeightEntryUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"validation_error",
			"Invalid request body",
			err.Error(),
		))
		return
	}

	entry, err := h.progressService.UpdateWeightEntry(c.Request.Context(), uint(id), userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"update_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(
		"Weight entry updated successfully",
		entry,
	))
}

// DeleteWeightEntry deletes a weight entry
// @Summary Delete weight entry
// @Description Delete a weight entry
// @Tags progress
// @Produce json
// @Security Bearer
// @Param id path int true "Weight Entry ID"
// @Success 200 {object} models.SuccessResponse
// @Router /progress/weight/{id} [delete]
func (h *ProgressHandler) DeleteWeightEntry(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	if err := h.progressService.DeleteWeightEntry(c.Request.Context(), uint(id), userID); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"delete_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(
		"Weight entry deleted successfully",
		nil,
	))
}

// GetWeightProgress retrieves weight progress statistics
// @Summary Get weight progress
// @Description Get weight progress statistics and trend
// @Tags progress
// @Produce json
// @Security Bearer
// @Success 200 {object} models.SuccessResponse
// @Router /progress/weight/stats [get]
func (h *ProgressHandler) GetWeightProgress(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	progress, err := h.progressService.GetWeightProgress(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"fetch_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(
		"Weight progress retrieved successfully",
		progress,
	))
}

// GenerateUploadURL generates a presigned URL for image upload
// @Summary Generate upload URL
// @Description Generate a presigned URL for uploading progress photos
// @Tags progress
// @Produce json
// @Security Bearer
// @Param filename query string true "Filename"
// @Success 200 {object} models.SuccessResponse
// @Router /progress/upload-url [get]
func (h *ProgressHandler) GenerateUploadURL(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	filename := c.Query("filename")

	if filename == "" {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"validation_error",
			"Filename is required",
			nil,
		))
		return
	}

	uploadURL, err := h.progressService.GenerateUploadURL(c.Request.Context(), userID, filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"generation_failed",
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(
		"Upload URL generated successfully",
		uploadURL,
	))
}
