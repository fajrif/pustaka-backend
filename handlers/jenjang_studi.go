package handlers

import (
	"pustaka-backend/config"
	"pustaka-backend/models"
	"github.com/gofiber/fiber/v2"
)

func GetAllJenjangStudi(c *fiber.Ctx) error {
	var jenjangStudi []models.JenjangStudi
	query := config.DB.Order("created_at DESC")

	if err := query.Find(&jenjangStudi).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch all jenjang studi",
		})
	}

	return c.JSON(fiber.Map{
		"jenjang_studi": jenjangStudi,
	})
}

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

func CreateJenjangStudi(c *fiber.Ctx) error {
	var jenjangStudi models.JenjangStudi
	if err := c.BodyParser(&jenjangStudi); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := config.DB.Create(&jenjangStudi).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create jenjang studi",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(jenjangStudi)
}

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

	if err := config.DB.Model(&jenjangStudi).Updates(jenjangStudi).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update jenjang studi",
		})
	}

	return c.JSON(jenjangStudi)
}

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
