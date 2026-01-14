package handlers

import (
	"pustaka-backend/config"
	"pustaka-backend/helpers"
	"pustaka-backend/models"
	"github.com/gofiber/fiber/v2"
)

// GetAllJenisBuku godoc
// @Summary Get all jenis buku
// @Description Retrieve all book types (jenis buku) ordered by creation date
// @Tags JenisBuku
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param search query string false "Search by code, name, or description"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items per page (default: 20)"
// @Success 200 {object} map[string]interface{} "List of all jenis buku with pagination"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/jenis-buku [get]
func GetAllJenisBuku(c *fiber.Ctx) error {
	var jenisBuku []models.JenisBuku

	// Get pagination parameters
	pagination := helpers.GetPaginationParams(c)

	query := config.DB.Order("created_at ASC")
	queryCount := config.DB.Model(&models.JenisBuku{})

	// add params for not using pagination
	if c.Query("all") == "true" {
		pagination.Limit = -1 // No limit
		pagination.Offset = 0 // No offset
	}

	// Filter search
	if searchQuery := c.Query("search"); searchQuery != "" {
		// Wrap string search with wildcard SQL LIKE
		searchTerm := "%" + searchQuery + "%"
		cond := "jenis_buku.code ILIKE ? OR jenis_buku.name ILIKE ? OR jenis_buku.description ILIKE ?"
		args := []interface{}{searchTerm, searchTerm, searchTerm}

		query = query.Where(cond, args...)
		queryCount = queryCount.Where(cond, args...)
	}

	// Apply pagination and fetch data
	if err := query.Offset(pagination.Offset).Limit(pagination.Limit).Find(&jenisBuku).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch all jenis buku",
		})
	}

	// Create pagination response
	response, err := helpers.CreatePaginationResponse(queryCount, jenisBuku, "jenis_buku", pagination.Page, pagination.Limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create pagination response",
		})
	}

	return c.JSON(response)
}

// GetJenisBuku godoc
// @Summary Get a jenis buku by ID
// @Description Retrieve a single book type by its ID
// @Tags JenisBuku
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "JenisBuku ID (UUID)"
// @Success 200 {object} map[string]interface{} "JenisBuku details"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "JenisBuku not found"
// @Router /api/jenis-buku/{id} [get]
func GetJenisBuku(c *fiber.Ctx) error {
	id := c.Params("id")

	var jenisBuku models.JenisBuku
	if err := config.DB.Where("id = ?", id).First(&jenisBuku).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "JenisBuku not found",
		})
	}

	return c.JSON(fiber.Map{
		"jenis_buku": jenisBuku,
	})
}

// CreateJenisBuku godoc
// @Summary Create a new jenis buku
// @Description Create a new book type entry
// @Tags JenisBuku
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.JenisBuku true "JenisBuku details"
// @Success 201 {object} models.JenisBuku "Created jenis buku"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/jenis-buku [post]
func CreateJenisBuku(c *fiber.Ctx) error {
	var jenisBuku models.JenisBuku
	if err := c.BodyParser(&jenisBuku); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := config.DB.Create(&jenisBuku).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create jenis buku",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(jenisBuku)
}

// UpdateJenisBuku godoc
// @Summary Update a jenis buku
// @Description Update an existing book type by ID
// @Tags JenisBuku
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "JenisBuku ID (UUID)"
// @Param request body models.JenisBuku true "Updated jenis buku details"
// @Success 200 {object} models.JenisBuku "Updated jenis buku"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "JenisBuku not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/jenis-buku/{id} [put]
func UpdateJenisBuku(c *fiber.Ctx) error {
	id := c.Params("id")

	var jenisBuku models.JenisBuku
	if err := config.DB.Where("id = ?", id).First(&jenisBuku).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "JenisBuku not found",
		})
	}

	if err := c.BodyParser(&jenisBuku); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := config.DB.Model(&jenisBuku).Updates(jenisBuku).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update jenis buku",
		})
	}

	return c.JSON(jenisBuku)
}

// DeleteJenisBuku godoc
// @Summary Delete a jenis buku
// @Description Delete a book type by ID
// @Tags JenisBuku
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "JenisBuku ID (UUID)"
// @Success 200 {object} map[string]interface{} "JenisBuku deleted successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "JenisBuku not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/jenis-buku/{id} [delete]
func DeleteJenisBuku(c *fiber.Ctx) error {
	id := c.Params("id")

	result := config.DB.Delete(&models.JenisBuku{}, "id = ?", id)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete jenis buku",
		})
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "JenisBuku not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "JenisBuku deleted successfully",
	})
}
