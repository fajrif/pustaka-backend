package handlers

import (
	"pustaka-backend/config"
	"pustaka-backend/helpers"
	"pustaka-backend/models"
	"github.com/gofiber/fiber/v2"
)

// GetAllJenjangStudi godoc
// @Summary Get all jenjang studi
// @Description Retrieve all education levels (jenjang studi) ordered by creation date
// @Tags JenjangStudi
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param search query string false "Search by code, name, or description"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items per page (default: 20)"
// @Success 200 {object} map[string]interface{} "List of all jenjang studi with pagination"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/jenjang-studi [get]
func GetAllJenjangStudi(c *fiber.Ctx) error {
	var jenjangStudi []models.JenjangStudi

	// Get pagination parameters
	pagination := helpers.GetPaginationParams(c)

	query := config.DB.Order("created_at ASC")
	queryCount := config.DB.Model(&models.JenjangStudi{})

	// add params for not using pagination
	if c.Query("all") == "true" {
		pagination.Limit = -1 // No limit
		pagination.Offset = 0 // No offset
	}

	// Filter search
	if searchQuery := c.Query("search"); searchQuery != "" {
		// Wrap string search with wildcard SQL LIKE
		searchTerm := "%" + searchQuery + "%"
		cond := "jenjang_studi.code ILIKE ? OR jenjang_studi.name ILIKE ? OR jenjang_studi.description ILIKE ?"
		args := []interface{}{searchTerm, searchTerm, searchTerm}

		query = query.Where(cond, args...)
		queryCount = queryCount.Where(cond, args...)
	}

	// Apply pagination and fetch data
	if err := query.Offset(pagination.Offset).Limit(pagination.Limit).Find(&jenjangStudi).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch all jenjang studi",
		})
	}

	// Create pagination response
	response, err := helpers.CreatePaginationResponse(queryCount, jenjangStudi, "jenjang_studi", pagination.Page, pagination.Limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create pagination response",
		})
	}

	return c.JSON(response)
}

// GetJenjangStudi godoc
// @Summary Get a jenjang studi by ID
// @Description Retrieve a single education level by its ID
// @Tags JenjangStudi
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "JenjangStudi ID (UUID)"
// @Success 200 {object} map[string]interface{} "JenjangStudi details"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "JenjangStudi not found"
// @Router /api/jenjang-studi/{id} [get]
func GetJenjangStudi(c *fiber.Ctx) error {
	id := c.Params("id")

	var jenjangStudi models.JenjangStudi
	if err := config.DB.Where("id = ?", id).First(&jenjangStudi).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "JenjangStudi not found",
		})
	}

	return c.JSON(fiber.Map{
		"jenjang_studi": jenjangStudi,
	})
}

// CreateJenjangStudi godoc
// @Summary Create a new jenjang studi
// @Description Create a new education level entry
// @Tags JenjangStudi
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.JenjangStudi true "JenjangStudi details"
// @Success 201 {object} models.JenjangStudi "Created jenjang studi"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/jenjang-studi [post]
func CreateJenjangStudi(c *fiber.Ctx) error {
	var jenjangStudi models.JenjangStudi
	if err := c.BodyParser(&jenjangStudi); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate period field - only 'S' or 'T' allowed
	if jenjangStudi.Period != "" && jenjangStudi.Period != "S" && jenjangStudi.Period != "T" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Period must be either 'S' or 'T'",
		})
	}

	if err := config.DB.Create(&jenjangStudi).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create jenjang studi",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(jenjangStudi)
}

// UpdateJenjangStudi godoc
// @Summary Update a jenjang studi
// @Description Update an existing education level by ID
// @Tags JenjangStudi
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "JenjangStudi ID (UUID)"
// @Param request body models.JenjangStudi true "Updated jenjang studi details"
// @Success 200 {object} models.JenjangStudi "Updated jenjang studi"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "JenjangStudi not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/jenjang-studi/{id} [put]
func UpdateJenjangStudi(c *fiber.Ctx) error {
	id := c.Params("id")

	var jenjangStudi models.JenjangStudi
	if err := config.DB.Where("id = ?", id).First(&jenjangStudi).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "JenjangStudi not found",
		})
	}

	if err := c.BodyParser(&jenjangStudi); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate period field - only 'S' or 'T' allowed
	if jenjangStudi.Period != "" && jenjangStudi.Period != "S" && jenjangStudi.Period != "T" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Period must be either 'S' or 'T'",
		})
	}

	if err := config.DB.Model(&jenjangStudi).Updates(jenjangStudi).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update jenjang studi",
		})
	}

	return c.JSON(jenjangStudi)
}

// DeleteJenjangStudi godoc
// @Summary Delete a jenjang studi
// @Description Delete an education level by ID
// @Tags JenjangStudi
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "JenjangStudi ID (UUID)"
// @Success 200 {object} map[string]interface{} "JenjangStudi deleted successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "JenjangStudi not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/jenjang-studi/{id} [delete]
func DeleteJenjangStudi(c *fiber.Ctx) error {
	id := c.Params("id")

	result := config.DB.Delete(&models.JenjangStudi{}, "id = ?", id)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete jenjang studi",
		})
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "JenjangStudi not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "JenjangStudi deleted successfully",
	})
}
