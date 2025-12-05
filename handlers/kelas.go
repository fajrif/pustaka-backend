package handlers

import (
	"pustaka-backend/config"
	"pustaka-backend/models"
	"github.com/gofiber/fiber/v2"
)

// GetAllKelas godoc
// @Summary Get all kelas
// @Description Retrieve all classes (kelas) ordered by creation date
// @Tags Kelas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param search query string false "Search by code, name, or description"
// @Success 200 {object} map[string]interface{} "List of all kelas"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/kelas [get]
func GetAllKelas(c *fiber.Ctx) error {
	var kelas []models.Kelas
	query := config.DB.Order("created_at DESC")

	// Filter search
	if searchQuery := c.Query("search"); searchQuery != "" {
		// Wrap string search with wildcard SQL LIKE
		searchTerm := "%" + searchQuery + "%"

		query = query.
			Where("kelas.code ILIKE ? OR kelas.name ILIKE ? OR kelas.description ILIKE ?", searchTerm, searchTerm, searchTerm)
	}

	if err := query.Find(&kelas).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch all kelas",
		})
	}

	return c.JSON(fiber.Map{
		"kelas": kelas,
	})
}

// GetKelas godoc
// @Summary Get a kelas by ID
// @Description Retrieve a single class by its ID
// @Tags Kelas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Kelas ID (UUID)"
// @Success 200 {object} map[string]interface{} "Kelas details"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Kelas not found"
// @Router /api/kelas/{id} [get]
func GetKelas(c *fiber.Ctx) error {
	id := c.Params("id")

	var kelas models.Kelas
	if err := config.DB.Where("id = ?", id).First(&kelas).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Kelas not found",
		})
	}

	return c.JSON(fiber.Map{
		"kelas": kelas,
	})
}

// CreateKelas godoc
// @Summary Create a new kelas
// @Description Create a new class entry
// @Tags Kelas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.Kelas true "Kelas details"
// @Success 201 {object} models.Kelas "Created kelas"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/kelas [post]
func CreateKelas(c *fiber.Ctx) error {
	var kelas models.Kelas
	if err := c.BodyParser(&kelas); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := config.DB.Create(&kelas).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create kelas",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(kelas)
}

// UpdateKelas godoc
// @Summary Update a kelas
// @Description Update an existing class by ID
// @Tags Kelas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Kelas ID (UUID)"
// @Param request body models.Kelas true "Updated kelas details"
// @Success 200 {object} models.Kelas "Updated kelas"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Kelas not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/kelas/{id} [put]
func UpdateKelas(c *fiber.Ctx) error {
	id := c.Params("id")

	var kelas models.Kelas
	if err := config.DB.Where("id = ?", id).First(&kelas).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Kelas not found",
		})
	}

	if err := c.BodyParser(&kelas); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := config.DB.Model(&kelas).Updates(kelas).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update kelas",
		})
	}

	return c.JSON(kelas)
}

// DeleteKelas godoc
// @Summary Delete a kelas
// @Description Delete a class by ID
// @Tags Kelas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Kelas ID (UUID)"
// @Success 200 {object} map[string]interface{} "Kelas deleted successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Kelas not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/kelas/{id} [delete]
func DeleteKelas(c *fiber.Ctx) error {
	id := c.Params("id")

	result := config.DB.Delete(&models.Kelas{}, "id = ?", id)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete kelas",
		})
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Kelas not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Kelas deleted successfully",
	})
}
