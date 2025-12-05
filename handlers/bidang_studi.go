package handlers

import (
	"pustaka-backend/config"
	"pustaka-backend/models"
	"github.com/gofiber/fiber/v2"
)

// GetAllBidangStudi godoc
// @Summary Get all bidang studi
// @Description Retrieve all fields of study (bidang studi) ordered by creation date
// @Tags BidangStudi
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param search query string false "Search by code, name, or description"
// @Success 200 {object} map[string]interface{} "List of all bidang studi"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/bidang-studi [get]
func GetAllBidangStudi(c *fiber.Ctx) error {
	var bidangStudi []models.BidangStudi
	query := config.DB.Order("created_at DESC")

	// Filter search
	if searchQuery := c.Query("search"); searchQuery != "" {
		// Wrap string search with wildcard SQL LIKE
		searchTerm := "%" + searchQuery + "%"

		query = query.
			Where("bidang_studi.code ILIKE ? OR bidang_studi.name ILIKE ? OR bidang_studi.description ILIKE ?", searchTerm, searchTerm, searchTerm)
	}

	if err := query.Find(&bidangStudi).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch all bidang studi",
		})
	}

	return c.JSON(fiber.Map{
		"bidang_studi": bidangStudi,
	})
}

// GetBidangStudi godoc
// @Summary Get a bidang studi by ID
// @Description Retrieve a single field of study by its ID
// @Tags BidangStudi
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "BidangStudi ID (UUID)"
// @Success 200 {object} map[string]interface{} "BidangStudi details"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "BidangStudi not found"
// @Router /api/bidang-studi/{id} [get]
func GetBidangStudi(c *fiber.Ctx) error {
	id := c.Params("id")

	var bidangStudi models.BidangStudi
	if err := config.DB.Where("id = ?", id).First(&bidangStudi).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "BidangStudi not found",
		})
	}

	return c.JSON(fiber.Map{
		"bidang_studi": bidangStudi,
	})
}

// CreateBidangStudi godoc
// @Summary Create a new bidang studi
// @Description Create a new field of study entry
// @Tags BidangStudi
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.BidangStudi true "BidangStudi details"
// @Success 201 {object} models.BidangStudi "Created bidang studi"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/bidang-studi [post]
func CreateBidangStudi(c *fiber.Ctx) error {
	var bidangStudi models.BidangStudi
	if err := c.BodyParser(&bidangStudi); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := config.DB.Create(&bidangStudi).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create bidang studi",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(bidangStudi)
}

// UpdateBidangStudi godoc
// @Summary Update a bidang studi
// @Description Update an existing field of study by ID
// @Tags BidangStudi
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "BidangStudi ID (UUID)"
// @Param request body models.BidangStudi true "Updated bidang studi details"
// @Success 200 {object} models.BidangStudi "Updated bidang studi"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "BidangStudi not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/bidang-studi/{id} [put]
func UpdateBidangStudi(c *fiber.Ctx) error {
	id := c.Params("id")

	var bidangStudi models.BidangStudi
	if err := config.DB.Where("id = ?", id).First(&bidangStudi).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "BidangStudi not found",
		})
	}

	if err := c.BodyParser(&bidangStudi); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := config.DB.Model(&bidangStudi).Updates(bidangStudi).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update bidang studi",
		})
	}

	return c.JSON(bidangStudi)
}

// DeleteBidangStudi godoc
// @Summary Delete a bidang studi
// @Description Delete a field of study by ID
// @Tags BidangStudi
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "BidangStudi ID (UUID)"
// @Success 200 {object} map[string]interface{} "BidangStudi deleted successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "BidangStudi not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/bidang-studi/{id} [delete]
func DeleteBidangStudi(c *fiber.Ctx) error {
	id := c.Params("id")

	result := config.DB.Delete(&models.BidangStudi{}, "id = ?", id)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete bidang studi",
		})
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "BidangStudi not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "BidangStudi deleted successfully",
	})
}
