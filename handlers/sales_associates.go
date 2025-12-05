package handlers

import (
	"pustaka-backend/config"
	"pustaka-backend/helpers"
	"pustaka-backend/models"
	"github.com/gofiber/fiber/v2"
)

// GetAllSalesAssociates godoc
// @Summary Get all sales associates
// @Description Retrieve all sales associates with their related city information
// @Tags SalesAssociates
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param search query string false "Search by code, name, or description"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items per page (default: 20)"
// @Success 200 {object} map[string]interface{} "List of all sales associates with pagination"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/sales-associates [get]
func GetAllSalesAssociates(c *fiber.Ctx) error {
	var salesAssociates []models.SalesAssociate

	// Get pagination parameters
	pagination := helpers.GetPaginationParams(c)

	query := config.DB.Order("created_at DESC")
	queryCount := config.DB.Model(&models.SalesAssociate{})

	// Filter search
	if searchQuery := c.Query("search"); searchQuery != "" {
		// Wrap string search with wildcard SQL LIKE
		searchTerm := "%" + searchQuery + "%"
		cond := "sales_associates.code ILIKE ? OR sales_associates.name ILIKE ? OR sales_associates.description ILIKE ?"
		args := []interface{}{searchTerm, searchTerm, searchTerm}

		query = query.Where(cond, args...)
		queryCount = queryCount.Where(cond, args...)
	}

	// Apply pagination and fetch data
	if err := query.Offset(pagination.Offset).Limit(pagination.Limit).Preload("City").Find(&salesAssociates).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch all sales associates",
		})
	}

	// Create pagination response
	response, err := helpers.CreatePaginationResponse(queryCount, salesAssociates, "sales_associates", pagination.Page, pagination.Limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create pagination response",
		})
	}

	return c.JSON(response)
}

// GetSalesAssociate godoc
// @Summary Get a sales associate by ID
// @Description Retrieve a single sales associate by its ID with related city information
// @Tags SalesAssociates
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "SalesAssociate ID (UUID)"
// @Success 200 {object} map[string]interface{} "SalesAssociate details"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "SalesAssociate not found"
// @Router /api/sales-associates/{id} [get]
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

// CreateSalesAssociate godoc
// @Summary Create a new sales associate
// @Description Create a new sales associate entry
// @Tags SalesAssociates
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.SalesAssociate true "SalesAssociate details"
// @Success 201 {object} models.SalesAssociate "Created sales associate"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/sales-associates [post]
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

// UpdateSalesAssociate godoc
// @Summary Update a sales associate
// @Description Update an existing sales associate by ID
// @Tags SalesAssociates
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "SalesAssociate ID (UUID)"
// @Param request body models.SalesAssociate true "Updated sales associate details"
// @Success 200 {object} models.SalesAssociate "Updated sales associate"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "SalesAssociate not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/sales-associates/{id} [put]
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

// DeleteSalesAssociate godoc
// @Summary Delete a sales associate
// @Description Delete a sales associate by ID
// @Tags SalesAssociates
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "SalesAssociate ID (UUID)"
// @Success 200 {object} map[string]interface{} "SalesAssociate deleted successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "SalesAssociate not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/sales-associates/{id} [delete]
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
