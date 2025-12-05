package handlers

import (
	"pustaka-backend/config"
	"pustaka-backend/models"
	"github.com/gofiber/fiber/v2"
	//"github.com/google/uuid"
)

// GetAllCities godoc
// @Summary Get all cities
// @Description Retrieve all cities ordered by creation date
// @Tags Cities
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param search query string false "Search by code or name"
// @Success 200 {object} map[string]interface{} "List of all cities"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/cities [get]
func GetAllCities(c *fiber.Ctx) error {
	var cities []models.City
	query := config.DB.Order("created_at DESC")

	// Filter search
	if searchQuery := c.Query("search"); searchQuery != "" {
		// Wrap string search with wildcard SQL LIKE
		searchTerm := "%" + searchQuery + "%"

		query = query.
			Where("cities.code ILIKE ? OR cities.name ILIKE ?", searchTerm, searchTerm)
	}

	if err := query.Find(&cities).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch all cities",
		})
	}

	return c.JSON(fiber.Map{
		"cities": cities,
	})
}

// GetCity godoc
// @Summary Get a city by ID
// @Description Retrieve a single city by its ID
// @Tags Cities
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "City ID (UUID)"
// @Success 200 {object} map[string]interface{} "City details"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "City not found"
// @Router /api/cities/{id} [get]
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

// CreateCity godoc
// @Summary Create a new city
// @Description Create a new city entry
// @Tags Cities
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.City true "City details"
// @Success 201 {object} models.City "Created city"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/cities [post]
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

// UpdateCity godoc
// @Summary Update a city
// @Description Update an existing city by ID
// @Tags Cities
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "City ID (UUID)"
// @Param request body models.City true "Updated city details"
// @Success 200 {object} models.City "Updated city"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "City not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/cities/{id} [put]
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

// DeleteCity godoc
// @Summary Delete a city
// @Description Delete a city by ID
// @Tags Cities
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "City ID (UUID)"
// @Success 200 {object} map[string]interface{} "City deleted successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "City not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/cities/{id} [delete]
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
