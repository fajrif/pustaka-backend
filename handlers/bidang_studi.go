package handlers

import (
	"pustaka-backend/config"
	"pustaka-backend/models"
	"github.com/gofiber/fiber/v2"
)

func GetAllBidangStudi(c *fiber.Ctx) error {
	var bidangStudi []models.BidangStudi
	query := config.DB.Order("created_at DESC")

	if err := query.Find(&bidangStudi).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch all bidang studi",
		})
	}

	return c.JSON(fiber.Map{
		"bidang_studi": bidangStudi,
	})
}

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
