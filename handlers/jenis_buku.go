package handlers

import (
	"pustaka-backend/config"
	"pustaka-backend/models"
	"github.com/gofiber/fiber/v2"
)

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
