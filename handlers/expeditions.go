package handlers

import (
	"pustaka-backend/config"
	"pustaka-backend/models"
	"github.com/gofiber/fiber/v2"
)

// GetAllExpeditions godoc
// @Summary Get all expeditions
// @Description Retrieve all expeditions with their related city information
// @Tags Expeditions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "List of all expeditions"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/expeditions [get]
func GetAllExpeditions(c *fiber.Ctx) error {
	var expeditions []models.Expedition
	query := config.DB.Order("created_at DESC")

	if err := query.Preload("City").Find(&expeditions).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch all expeditions",
		})
	}

	return c.JSON(fiber.Map{
		"expeditions": expeditions,
	})
}

// GetExpedition godoc
// @Summary Get an expedition by ID
// @Description Retrieve a single expedition by its ID with related city information
// @Tags Expeditions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Expedition ID (UUID)"
// @Success 200 {object} map[string]interface{} "Expedition details"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Expedition not found"
// @Router /api/expeditions/{id} [get]
func GetExpedition(c *fiber.Ctx) error {
	id := c.Params("id")

	var expedition models.Expedition
	if err := config.DB.Preload("City").Where("id = ?", id).First(&expedition).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Expedition not found",
		})
	}

	return c.JSON(fiber.Map{
		"expedition": expedition,
	})
}

// CreateExpedition godoc
// @Summary Create a new expedition
// @Description Create a new expedition entry
// @Tags Expeditions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.Expedition true "Expedition details"
// @Success 201 {object} models.Expedition "Created expedition"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/expeditions [post]
func CreateExpedition(c *fiber.Ctx) error {
	var expedition models.Expedition
	if err := c.BodyParser(&expedition); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := config.DB.Create(&expedition).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create expedition",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(expedition)
}

// UpdateExpedition godoc
// @Summary Update an expedition
// @Description Update an existing expedition by ID
// @Tags Expeditions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Expedition ID (UUID)"
// @Param request body models.Expedition true "Updated expedition details"
// @Success 200 {object} models.Expedition "Updated expedition"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Expedition not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/expeditions/{id} [put]
func UpdateExpedition(c *fiber.Ctx) error {
	id := c.Params("id")

	var expedition models.Expedition
	if err := config.DB.Where("id = ?", id).First(&expedition).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Expedition not found",
		})
	}

	if err := c.BodyParser(&expedition); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := config.DB.Model(&expedition).Updates(expedition).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update expedition",
		})
	}

	return c.JSON(expedition)
}

// DeleteExpedition godoc
// @Summary Delete an expedition
// @Description Delete an expedition by ID
// @Tags Expeditions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Expedition ID (UUID)"
// @Success 200 {object} map[string]interface{} "Expedition deleted successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Expedition not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/expeditions/{id} [delete]
func DeleteExpedition(c *fiber.Ctx) error {
	id := c.Params("id")

	result := config.DB.Delete(&models.Expedition{}, "id = ?", id)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete expedition",
		})
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Expedition not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Expedition deleted successfully",
	})
}
