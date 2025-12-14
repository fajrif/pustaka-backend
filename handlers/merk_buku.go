package handlers

import (
	"pustaka-backend/config"
	"pustaka-backend/helpers"
	"pustaka-backend/models"
	"github.com/gofiber/fiber/v2"
)

// GetAllMerkBuku godoc
// @Summary Get all merk buku
// @Description Retrieve all book types (merk buku) ordered by creation date
// @Tags MerkBuku
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param search query string false "Search by code, name, or description"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items per page (default: 20)"
// @Success 200 {object} map[string]interface{} "List of all merk buku with pagination"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/merk-buku [get]
func GetAllMerkBuku(c *fiber.Ctx) error {
	var merkBuku []models.MerkBuku

	// Get pagination parameters
	pagination := helpers.GetPaginationParams(c)

	query := config.DB.Order("created_at DESC")
	queryCount := config.DB.Model(&models.MerkBuku{})

	// add params for not using pagination
	if c.Query("all") == "true" {
		pagination.Limit = -1 // No limit
		pagination.Offset = 0 // No offset
	}

	// Filter search
	if searchQuery := c.Query("search"); searchQuery != "" {
		// Wrap string search with wildcard SQL LIKE
		searchTerm := "%" + searchQuery + "%"
		cond := "merk_buku.code ILIKE ? OR merk_buku.name ILIKE ? OR merk_buku.description ILIKE ?"
		args := []interface{}{searchTerm, searchTerm, searchTerm}

		query = query.Where(cond, args...)
		queryCount = queryCount.Where(cond, args...)
	}

	// Apply pagination and fetch data
	if err := query.Offset(pagination.Offset).Limit(pagination.Limit).Find(&merkBuku).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch all merk buku",
		})
	}

	// Create pagination response
	response, err := helpers.CreatePaginationResponse(queryCount, merkBuku, "merk_buku", pagination.Page, pagination.Limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create pagination response",
		})
	}

	return c.JSON(response)
}

// GetMerkBuku godoc
// @Summary Get a merk buku by ID
// @Description Retrieve a single book type by its ID
// @Tags MerkBuku
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "MerkBuku ID (UUID)"
// @Success 200 {object} map[string]interface{} "MerkBuku details"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "MerkBuku not found"
// @Router /api/merk-buku/{id} [get]
func GetMerkBuku(c *fiber.Ctx) error {
	id := c.Params("id")

	var merkBuku models.MerkBuku
	if err := config.DB.Where("id = ?", id).First(&merkBuku).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "MerkBuku not found",
		})
	}

	return c.JSON(fiber.Map{
		"merk_buku": merkBuku,
	})
}

// CreateMerkBuku godoc
// @Summary Create a new merk buku
// @Description Create a new book type entry
// @Tags MerkBuku
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.MerkBuku true "MerkBuku details"
// @Success 201 {object} models.MerkBuku "Created merk buku"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/merk-buku [post]
func CreateMerkBuku(c *fiber.Ctx) error {
	var merkBuku models.MerkBuku
	if err := c.BodyParser(&merkBuku); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := config.DB.Create(&merkBuku).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create merk buku",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(merkBuku)
}

// UpdateMerkBuku godoc
// @Summary Update a merk buku
// @Description Update an existing book type by ID
// @Tags MerkBuku
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "MerkBuku ID (UUID)"
// @Param request body models.MerkBuku true "Updated merk buku details"
// @Success 200 {object} models.MerkBuku "Updated merk buku"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "MerkBuku not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/merk-buku/{id} [put]
func UpdateMerkBuku(c *fiber.Ctx) error {
	id := c.Params("id")

	var merkBuku models.MerkBuku
	if err := config.DB.Where("id = ?", id).First(&merkBuku).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "MerkBuku not found",
		})
	}

	if err := c.BodyParser(&merkBuku); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := config.DB.Model(&merkBuku).Updates(merkBuku).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update merk buku",
		})
	}

	return c.JSON(merkBuku)
}

// DeleteMerkBuku godoc
// @Summary Delete a merk buku
// @Description Delete a book type by ID
// @Tags MerkBuku
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "MerkBuku ID (UUID)"
// @Success 200 {object} map[string]interface{} "MerkBuku deleted successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "MerkBuku not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/merk-buku/{id} [delete]
func DeleteMerkBuku(c *fiber.Ctx) error {
	id := c.Params("id")

	result := config.DB.Delete(&models.MerkBuku{}, "id = ?", id)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete merk buku",
		})
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "MerkBuku not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "MerkBuku deleted successfully",
	})
}

