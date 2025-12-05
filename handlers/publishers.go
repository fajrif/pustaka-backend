package handlers

import (
	"pustaka-backend/config"
	"pustaka-backend/models"
	"github.com/gofiber/fiber/v2"
)

// GetAllPublishers godoc
// @Summary Get all publishers
// @Description Retrieve all publishers with their related city information
// @Tags Publishers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "List of all publishers"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/publishers [get]
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

// GetPublisher godoc
// @Summary Get a publisher by ID
// @Description Retrieve a single publisher by its ID with related city information
// @Tags Publishers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Publisher ID (UUID)"
// @Success 200 {object} map[string]interface{} "Publisher details"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Publisher not found"
// @Router /api/publishers/{id} [get]
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

// CreatePublisher godoc
// @Summary Create a new publisher
// @Description Create a new publisher entry
// @Tags Publishers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.Publisher true "Publisher details"
// @Success 201 {object} models.Publisher "Created publisher"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/publishers [post]
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

// UpdatePublisher godoc
// @Summary Update a publisher
// @Description Update an existing publisher by ID
// @Tags Publishers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Publisher ID (UUID)"
// @Param request body models.Publisher true "Updated publisher details"
// @Success 200 {object} models.Publisher "Updated publisher"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Publisher not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/publishers/{id} [put]
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

// DeletePublisher godoc
// @Summary Delete a publisher
// @Description Delete a publisher by ID
// @Tags Publishers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Publisher ID (UUID)"
// @Success 200 {object} map[string]interface{} "Publisher deleted successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Publisher not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/publishers/{id} [delete]
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
