package handlers

import (
	"pustaka-backend/config"
	"pustaka-backend/helpers"
	"pustaka-backend/models"

	"github.com/gofiber/fiber/v2"
)

// GetAllBillers godoc
// @Summary Get all billers
// @Description Retrieve all billers with their related city information
// @Tags Billers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param search query string false "Search by code, name, or description"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items per page (default: 20)"
// @Success 200 {object} map[string]interface{} "List of all billers with pagination"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/billers [get]
func GetAllBillers(c *fiber.Ctx) error {
	var billers []models.Biller

	// Get pagination parameters
	pagination := helpers.GetPaginationParams(c)

	query := config.DB.Order("created_at DESC")
	queryCount := config.DB.Model(&models.Biller{})

	// add params for not using pagination
	if c.Query("all") == "true" {
		pagination.Limit = -1 // No limit
		pagination.Offset = 0 // No offset
	}

	// Filter search
	if searchQuery := c.Query("search"); searchQuery != "" {
		// Wrap string search with wildcard SQL LIKE
		searchTerm := "%" + searchQuery + "%"
		cond := "billers.code ILIKE ? OR billers.name ILIKE ? OR billers.description ILIKE ?"
		args := []interface{}{searchTerm, searchTerm, searchTerm}

		query = query.Where(cond, args...)
		queryCount = queryCount.Where(cond, args...)
	}

	// Apply pagination and fetch data
	if err := query.Offset(pagination.Offset).Limit(pagination.Limit).Preload("City").Find(&billers).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch all billers",
		})
	}

	// Create pagination response
	response, err := helpers.CreatePaginationResponse(queryCount, billers, "billers", pagination.Page, pagination.Limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create pagination response",
		})
	}

	return c.JSON(response)
}

// GetBiller godoc
// @Summary Get a biller by ID
// @Description Retrieve a single biller by its ID with related city information
// @Tags Billers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Biller ID (UUID)"
// @Success 200 {object} map[string]interface{} "Biller details"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Biller not found"
// @Router /api/billers/{id} [get]
func GetBiller(c *fiber.Ctx) error {
	id := c.Params("id")

	var biller models.Biller
	if err := config.DB.Preload("City").Where("id = ?", id).First(&biller).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Biller not found",
		})
	}

	return c.JSON(fiber.Map{
		"biller": biller,
	})
}

// CreateBiller godoc
// @Summary Create a new biller
// @Description Create a new biller entry
// @Tags Billers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.Biller true "Biller details"
// @Success 201 {object} models.Biller "Created biller"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/billers [post]
func CreateBiller(c *fiber.Ctx) error {
	var biller models.Biller
	if err := c.BodyParser(&biller); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := config.DB.Create(&biller).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create biller",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(biller)
}

// UpdateBiller godoc
// @Summary Update a biller
// @Description Update an existing biller by ID
// @Tags Billers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Biller ID (UUID)"
// @Param request body models.Biller true "Updated biller details"
// @Success 200 {object} models.Biller "Updated biller"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Biller not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/billers/{id} [put]
func UpdateBiller(c *fiber.Ctx) error {
	id := c.Params("id")

	var biller models.Biller
	if err := config.DB.Where("id = ?", id).First(&biller).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Biller not found",
		})
	}

	if err := c.BodyParser(&biller); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := config.DB.Model(&biller).Updates(biller).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update biller",
		})
	}

	return c.JSON(biller)
}

// DeleteBiller godoc
// @Summary Delete a biller
// @Description Delete a biller by ID
// @Tags Billers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Biller ID (UUID)"
// @Success 200 {object} map[string]interface{} "Biller deleted successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Biller not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/billers/{id} [delete]
func DeleteBiller(c *fiber.Ctx) error {
	id := c.Params("id")

	result := config.DB.Delete(&models.Biller{}, "id = ?", id)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete biller",
		})
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Biller not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Biller deleted successfully",
	})
}
