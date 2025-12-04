package handlers

import (
	"pustaka-backend/config"
	"pustaka-backend/models"
	"github.com/gofiber/fiber/v2"
)

func GetAllPublishers(c *fiber.Ctx) error {
	var publishers []models.Publisher
	query := config.DB.Order("created_at DESC")

	if err := query.Preload("City").Find(&publishers).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch all publishers",
		})
	}

	return c.JSON(fiber.Map{
		"publishers": publishers,
	})
}

func GetPublisher(c *fiber.Ctx) error {
	id := c.Params("id")

	var publisher models.Publisher
	if err := config.DB.Preload("City").Where("id = ?", id).First(&publisher).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Publisher not found",
		})
	}

	return c.JSON(fiber.Map{
		"publisher": publisher,
	})
}

func CreatePublisher(c *fiber.Ctx) error {
	var publisher models.Publisher
	if err := c.BodyParser(&publisher); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := config.DB.Create(&publisher).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create publisher",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(publisher)
}

func UpdatePublisher(c *fiber.Ctx) error {
	id := c.Params("id")

	var publisher models.Publisher
	if err := config.DB.Where("id = ?", id).First(&publisher).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Publisher not found",
		})
	}

	if err := c.BodyParser(&publisher); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := config.DB.Model(&publisher).Updates(publisher).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update publisher",
		})
	}

	return c.JSON(publisher)
}

func DeletePublisher(c *fiber.Ctx) error {
	id := c.Params("id")

	result := config.DB.Delete(&models.Publisher{}, "id = ?", id)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete publisher",
		})
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Publisher not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Publisher deleted successfully",
	})
}
