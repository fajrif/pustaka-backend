package handlers

import (
	"pustaka-backend/config"
	"pustaka-backend/models"
	"github.com/gofiber/fiber/v2"
)

func GetAllSalesAssociates(c *fiber.Ctx) error {
	var salesAssociates []models.SalesAssociate
	query := config.DB.Order("created_at DESC")

	if err := query.Preload("City").Find(&salesAssociates).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch all sales associates",
		})
	}

	return c.JSON(fiber.Map{
		"sales_associates": salesAssociates,
	})
}

func GetSalesAssociate(c *fiber.Ctx) error {
	id := c.Params("id")

	var salesAssociate models.SalesAssociate
	if err := config.DB.Preload("City").Where("id = ?", id).First(&salesAssociate).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "SalesAssociate not found",
		})
	}

	return c.JSON(fiber.Map{
		"sales_associate": salesAssociate,
	})
}

func CreateSalesAssociate(c *fiber.Ctx) error {
	var salesAssociate models.SalesAssociate
	if err := c.BodyParser(&salesAssociate); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := config.DB.Create(&salesAssociate).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create sales associate",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(salesAssociate)
}

func UpdateSalesAssociate(c *fiber.Ctx) error {
	id := c.Params("id")

	var salesAssociate models.SalesAssociate
	if err := config.DB.Where("id = ?", id).First(&salesAssociate).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "SalesAssociate not found",
		})
	}

	if err := c.BodyParser(&salesAssociate); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := config.DB.Model(&salesAssociate).Updates(salesAssociate).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update sales associate",
		})
	}

	return c.JSON(salesAssociate)
}

func DeleteSalesAssociate(c *fiber.Ctx) error {
	id := c.Params("id")

	result := config.DB.Delete(&models.SalesAssociate{}, "id = ?", id)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete sales associate",
		})
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "SalesAssociate not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "SalesAssociate deleted successfully",
	})
}
