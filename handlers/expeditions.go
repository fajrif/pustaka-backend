package handlers

import (
	"pustaka-backend/config"
	"pustaka-backend/models"
	"github.com/gofiber/fiber/v2"
)

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
