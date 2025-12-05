package handlers

import (
	"pustaka-backend/config"
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
// @Success 200 {object} map[string]interface{} "List of all jenis buku"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/jenis-buku [get]
func GetAllJenisBuku(c *fiber.Ctx) error {
	var jenisBuku []models.JenisBuku
	query := config.DB.Order("created_at DESC")

	if err := query.Find(&jenisBuku).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch all jenis buku",
		})
	}

	return c.JSON(fiber.Map{
		"jenis_buku": jenisBuku,
	})
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
