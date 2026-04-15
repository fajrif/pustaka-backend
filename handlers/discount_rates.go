package handlers

import (
	"pustaka-backend/config"
	"pustaka-backend/helpers"
	"pustaka-backend/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

// GetAllDiscountRates godoc
// @Summary Get all discount rates
// @Description Retrieve all discount rates with pagination
// @Tags DiscountRates
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param all query bool false "Get all records without pagination"
// @Success 200 {object} map[string]interface{} "List of discount rates with pagination"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/discount-rates [get]
func GetAllDiscountRates(c *fiber.Ctx) error {
	var discountRates []models.DiscountRate

	pagination := helpers.GetPaginationParams(c)

	query := config.DB.Model(&models.DiscountRate{}).Order("created_at ASC")
	queryCount := config.DB.Model(&models.DiscountRate{})

	if c.Query("all") == "true" {
		pagination.Limit = -1
		pagination.Offset = 0
	}

	if err := query.Offset(pagination.Offset).Limit(pagination.Limit).Find(&discountRates).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch discount rates",
		})
	}

	response, err := helpers.CreatePaginationResponse(queryCount, discountRates, "discount_rates", pagination.Page, pagination.Limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create pagination response",
		})
	}

	return c.JSON(response)
}

// GetDiscountRate godoc
// @Summary Get a discount rate by ID
// @Description Retrieve a single discount rate by its ID
// @Tags DiscountRates
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "DiscountRate ID (UUID)"
// @Success 200 {object} map[string]interface{} "DiscountRate details"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "DiscountRate not found"
// @Router /api/discount-rates/{id} [get]
func GetDiscountRate(c *fiber.Ctx) error {
	id := c.Params("id")

	var discountRate models.DiscountRate
	if err := config.DB.Where("id = ?", id).First(&discountRate).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "DiscountRate not found",
		})
	}

	return c.JSON(fiber.Map{
		"discount_rate": discountRate,
	})
}

// CreateDiscountRate godoc
// @Summary Create a new discount rate
// @Description Create a new discount rate entry
// @Tags DiscountRates
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.DiscountRate true "DiscountRate details"
// @Success 201 {object} models.DiscountRate "Created discount rate"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/discount-rates [post]
func CreateDiscountRate(c *fiber.Ctx) error {
	var discountRate models.DiscountRate
	if err := c.BodyParser(&discountRate); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if discountRate.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "name is required",
		})
	}

	if discountRate.Discount < 0 || discountRate.Discount > 100 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "discount must be between 0 and 100",
		})
	}

	if err := config.DB.Create(&discountRate).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create discount rate",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(discountRate)
}

// UpdateDiscountRate godoc
// @Summary Update a discount rate
// @Description Update an existing discount rate by ID
// @Tags DiscountRates
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "DiscountRate ID (UUID)"
// @Param request body models.DiscountRate true "Updated discount rate details"
// @Success 200 {object} models.DiscountRate "Updated discount rate"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "DiscountRate not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/discount-rates/{id} [put]
func UpdateDiscountRate(c *fiber.Ctx) error {
	id := c.Params("id")

	var discountRate models.DiscountRate
	if err := config.DB.Where("id = ?", id).First(&discountRate).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "DiscountRate not found",
		})
	}

	if err := c.BodyParser(&discountRate); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if discountRate.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "name is required",
		})
	}

	if discountRate.Discount < 0 || discountRate.Discount > 100 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "discount must be between 0 and 100",
		})
	}

	if err := config.DB.Model(&discountRate).Updates(discountRate).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update discount rate",
		})
	}

	return c.JSON(discountRate)
}

// DeleteDiscountRate godoc
// @Summary Delete a discount rate
// @Description Delete a discount rate by ID
// @Tags DiscountRates
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "DiscountRate ID (UUID)"
// @Success 200 {object} map[string]interface{} "DiscountRate deleted successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "DiscountRate not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/discount-rates/{id} [delete]
func DeleteDiscountRate(c *fiber.Ctx) error {
	id := c.Params("id")

	result := config.DB.Delete(&models.DiscountRate{}, "id = ?", id)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete discount rate",
		})
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "DiscountRate not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "DiscountRate deleted successfully",
	})
}

// GetApplicableDiscount godoc
// @Summary Preview applicable discount
// @Description Preview the discount percentage and name based on installment date vs due dates
// @Tags DiscountRates
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param installment_date query string true "Installment date (ISO 8601: YYYY-MM-DD or RFC3339)"
// @Param due_date query string false "Due date (ISO 8601: YYYY-MM-DD or RFC3339)"
// @Param secondary_due_date query string false "Secondary due date (ISO 8601: YYYY-MM-DD or RFC3339)"
// @Success 200 {object} map[string]interface{} "Applicable discount details"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /api/discount-rates/applicable [get]
func GetApplicableDiscount(c *fiber.Ctx) error {
	installmentDateStr := c.Query("installment_date")
	dueDateStr := c.Query("due_date")
	secondaryDueDateStr := c.Query("secondary_due_date")

	if installmentDateStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "installment_date query parameter is required",
		})
	}

	installmentDate, err := time.Parse(time.RFC3339, installmentDateStr)
	if err != nil {
		installmentDate, err = time.Parse("2006-01-02", installmentDateStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid installment_date format. Use ISO 8601 (YYYY-MM-DD or RFC3339)",
			})
		}
	}

	var dueDate, secondaryDueDate *time.Time
	if dueDateStr != "" {
		parsed, err := time.Parse(time.RFC3339, dueDateStr)
		if err != nil {
			parsed, err = time.Parse("2006-01-02", dueDateStr)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Invalid due_date format. Use ISO 8601 (YYYY-MM-DD or RFC3339)",
				})
			}
		}
		dueDate = &parsed
	}

	if secondaryDueDateStr != "" {
		parsed, err := time.Parse(time.RFC3339, secondaryDueDateStr)
		if err != nil {
			parsed, err = time.Parse("2006-01-02", secondaryDueDateStr)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Invalid secondary_due_date format. Use ISO 8601 (YYYY-MM-DD or RFC3339)",
				})
			}
		}
		secondaryDueDate = &parsed
	}

	var discountPercentage float64
	var discountName string

	if dueDate != nil && !installmentDate.After(*dueDate) {
		var earlyDiscount models.DiscountRate
		if err := config.DB.Where("name ILIKE ?", "%early%").Order("created_at ASC").First(&earlyDiscount).Error; err == nil {
			discountPercentage = earlyDiscount.Discount
			discountName = earlyDiscount.Name
		} else {
			discountPercentage = 8
			discountName = "Early Payment Discount"
		}
	} else if secondaryDueDate != nil && !installmentDate.After(*secondaryDueDate) {
		var secondaryDiscount models.DiscountRate
		if err := config.DB.Where("name ILIKE ?", "%secondary%").Order("created_at ASC").First(&secondaryDiscount).Error; err == nil {
			discountPercentage = secondaryDiscount.Discount
			discountName = secondaryDiscount.Name
		} else {
			discountPercentage = 5
			discountName = "Secondary Payment Discount"
		}
	}

	return c.JSON(fiber.Map{
		"installment_date":    installmentDate,
		"due_date":            dueDate,
		"secondary_due_date":  secondaryDueDate,
		"discount_percentage": discountPercentage,
		"discount_name":       discountName,
	})
}
