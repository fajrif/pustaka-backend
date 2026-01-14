package handlers

import (
	"pustaka-backend/config"
	"pustaka-backend/helpers"
	"pustaka-backend/models"

	"github.com/gofiber/fiber/v2"
)

// GetAllCurriculum godoc
// @Summary Get all curriculum
// @Description Retrieve all curriculum records ordered by creation date
// @Tags Curriculum
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param search query string false "Search by code, name, or description"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items per page (default: 20)"
// @Param all query bool false "Get all records without pagination"
// @Success 200 {object} map[string]interface{} "List of all curriculum with pagination"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/curriculum [get]
func GetAllCurriculums(c *fiber.Ctx) error {
	var curriculum []models.Curriculum

	// Get pagination parameters
	pagination := helpers.GetPaginationParams(c)

	query := config.DB.Order("created_at ASC")
	queryCount := config.DB.Model(&models.Curriculum{})

	// add params for not using pagination
	if c.Query("all") == "true" {
		pagination.Limit = -1 // No limit
		pagination.Offset = 0 // No offset
	}

	// Filter search
	if searchQuery := c.Query("search"); searchQuery != "" {
		// Wrap string search with wildcard SQL LIKE
		searchTerm := "%" + searchQuery + "%"
		cond := "curriculum.code ILIKE ? OR curriculum.name ILIKE ? OR curriculum.description ILIKE ?"
		args := []interface{}{searchTerm, searchTerm, searchTerm}

		query = query.Where(cond, args...)
		queryCount = queryCount.Where(cond, args...)
	}

	// Apply pagination and fetch data
	if err := query.Offset(pagination.Offset).Limit(pagination.Limit).Find(&curriculum).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch all curriculum",
		})
	}

	// Create pagination response
	response, err := helpers.CreatePaginationResponse(queryCount, curriculum, "curriculums", pagination.Page, pagination.Limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create pagination response",
		})
	}

	return c.JSON(response)
}

// GetCurriculum godoc
// @Summary Get a curriculum by ID
// @Description Retrieve a single curriculum by its ID
// @Tags Curriculum
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Curriculum ID (UUID)"
// @Success 200 {object} map[string]interface{} "Curriculum details"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Curriculum not found"
// @Router /api/curriculum/{id} [get]
func GetCurriculum(c *fiber.Ctx) error {
	id := c.Params("id")

	var curriculum models.Curriculum
	if err := config.DB.Where("id = ?", id).First(&curriculum).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Curriculum not found",
		})
	}

	return c.JSON(fiber.Map{
		"curriculum": curriculum,
	})
}

// CreateCurriculum godoc
// @Summary Create a new curriculum
// @Description Create a new curriculum entry
// @Tags Curriculum
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.Curriculum true "Curriculum details"
// @Success 201 {object} models.Curriculum "Created curriculum"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/curriculum [post]
func CreateCurriculum(c *fiber.Ctx) error {
	var curriculum models.Curriculum
	if err := c.BodyParser(&curriculum); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := config.DB.Create(&curriculum).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create curriculum",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(curriculum)
}

// UpdateCurriculum godoc
// @Summary Update a curriculum
// @Description Update an existing curriculum by ID
// @Tags Curriculum
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Curriculum ID (UUID)"
// @Param request body models.Curriculum true "Updated curriculum details"
// @Success 200 {object} models.Curriculum "Updated curriculum"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Curriculum not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/curriculum/{id} [put]
func UpdateCurriculum(c *fiber.Ctx) error {
	id := c.Params("id")

	var curriculum models.Curriculum
	if err := config.DB.Where("id = ?", id).First(&curriculum).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Curriculum not found",
		})
	}

	if err := c.BodyParser(&curriculum); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := config.DB.Model(&curriculum).Updates(curriculum).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update curriculum",
		})
	}

	return c.JSON(curriculum)
}

// DeleteCurriculum godoc
// @Summary Delete a curriculum
// @Description Delete a curriculum by ID
// @Tags Curriculum
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Curriculum ID (UUID)"
// @Success 200 {object} map[string]interface{} "Curriculum deleted successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Curriculum not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/curriculum/{id} [delete]
func DeleteCurriculum(c *fiber.Ctx) error {
	id := c.Params("id")

	result := config.DB.Delete(&models.Curriculum{}, "id = ?", id)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete curriculum",
		})
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Curriculum not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Curriculum deleted successfully",
	})
}
