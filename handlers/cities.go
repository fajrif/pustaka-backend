package handlers

import (
	"pustaka-backend/config"
	"pustaka-backend/models"
	"github.com/gofiber/fiber/v2"
	//"github.com/google/uuid"
)

func GetAllCities(c *fiber.Ctx) error {
	var cities []models.City
	query := config.DB.Order("created_at DESC")

	if err := query.Find(&cities).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch all cities",
		})
	}

	return c.JSON(fiber.Map{
		"cities": cities,
	})
}

func GetCity(c *fiber.Ctx) error {
	id := c.Params("id")

	var city models.City
	if err := config.DB.Where("id = ?", id).First(&city).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "City not found",
		})
	}

	return c.JSON(fiber.Map{
		"city": city,
	})
}

func CreateCity(c *fiber.Ctx) error {
	var city models.City
	if err := c.BodyParser(&city); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := config.DB.Create(&city).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create city",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(city)
}

func UpdateCity(c *fiber.Ctx) error {
	id := c.Params("id")

	var city models.City
	if err := config.DB.Where("id = ?", id).First(&city).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "City not found",
		})
	}

	if err := c.BodyParser(&city); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := config.DB.Model(&city).Updates(city).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update city",
		})
	}

	return c.JSON(city)
}

func DeleteCity(c *fiber.Ctx) error {
	id := c.Params("id")

	result := config.DB.Delete(&models.City{}, "id = ?", id)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete city",
		})
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "City not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "City deleted successfully",
	})
}
