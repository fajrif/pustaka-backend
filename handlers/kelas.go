package handlers

import (
	"pustaka-backend/config"
	"pustaka-backend/models"
	"github.com/gofiber/fiber/v2"
)

func GetAllKelas(c *fiber.Ctx) error {
	var kelas []models.Kelas
	query := config.DB.Order("created_at DESC")

	if err := query.Find(&kelas).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch all kelas",
		})
	}

	return c.JSON(fiber.Map{
		"kelas": kelas,
	})
}

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
