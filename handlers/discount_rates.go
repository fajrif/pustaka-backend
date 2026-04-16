package handlers

import (
	"fmt"
	"pustaka-backend/config"
	"pustaka-backend/helpers"
	"pustaka-backend/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

type CreateDiscountRateRequest struct {
	Name        string  `json:"name"`
	Discount    float64 `json:"discount"`
	Periode     int     `json:"periode"`
	Year        any     `json:"year"`
	StartDate   *string `json:"start_date"`
	EndDate     *string `json:"end_date"`
	Description *string `json:"description"`
}

func parseYearValue(year any) (string, error) {
	switch v := year.(type) {
	case string:
		return v, nil
	case float64:
		return fmt.Sprintf("%.0f", v), nil
	case int:
		return fmt.Sprintf("%d", v), nil
	default:
		return "", fmt.Errorf("year must be a string or number")
	}
}

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

	query := config.DB.Order("created_at ASC")
	queryCount := config.DB.Model(&models.DiscountRate{})

	if c.Query("all") == "true" {
		pagination.Limit = -1
		pagination.Offset = 0
	}

	if searchQuery := c.Query("search"); searchQuery != "" {
		searchTerm := "%" + searchQuery + "%"
		cond := "discount_rates.name ILIKE ? OR discount_rates.description ILIKE ?"
		args := []interface{}{searchTerm, searchTerm}

		query = query.Where(cond, args...)
		queryCount = queryCount.Where(cond, args...)
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
// @Param request body CreateDiscountRateRequest true "DiscountRate details"
// @Success 201 {object} models.DiscountRate "Created discount rate"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/discount-rates [post]
func CreateDiscountRate(c *fiber.Ctx) error {
	var req CreateDiscountRateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "name is required",
		})
	}

	if req.Discount < 0 || req.Discount > 100 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "discount must be between 0 and 100",
		})
	}

	if req.Periode < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "periode must be at least 1",
		})
	}

	yearStr, err := parseYearValue(req.Year)
	if err != nil || yearStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "year is required and must be a valid year",
		})
	}

	startDate, err := helpers.ParseDateString(req.StartDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid start_date format. Use YYYY-MM-DD",
		})
	}

	endDate, err := helpers.ParseDateString(req.EndDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid end_date format. Use YYYY-MM-DD",
		})
	}

	var existingCount int64
	config.DB.Model(&models.DiscountRate{}).
		Where("periode = ? AND year = ?", req.Periode, yearStr).
		Count(&existingCount)
	if existingCount > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "A discount rate for this periode and year already exists",
		})
	}

	if startDate != nil && endDate != nil {
		if endDate.Before(*startDate) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "end_date must be after start_date",
			})
		}
	}

	discountRate := models.DiscountRate{
		Name:        req.Name,
		Discount:    req.Discount,
		Periode:     req.Periode,
		Year:        yearStr,
		StartDate:   startDate,
		EndDate:     endDate,
		Description: req.Description,
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
// @Param request body CreateDiscountRateRequest true "Updated discount rate details"
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

	var req CreateDiscountRateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "name is required",
		})
	}

	if req.Discount < 0 || req.Discount > 100 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "discount must be between 0 and 100",
		})
	}

	if req.Periode < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "periode must be at least 1",
		})
	}

	yearStr, err := parseYearValue(req.Year)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "year is required and must be a valid year",
		})
	}

	if req.Periode != discountRate.Periode || yearStr != discountRate.Year {
		var existingCount int64
		config.DB.Model(&models.DiscountRate{}).
			Where("periode = ? AND year = ? AND id != ?", req.Periode, yearStr, id).
			Count(&existingCount)
		if existingCount > 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "A discount rate for this periode and year already exists",
			})
		}
	}

	startDate, _ := helpers.ParseDateString(req.StartDate)
	endDate, _ := helpers.ParseDateString(req.EndDate)
	if startDate != nil && endDate != nil {
		if endDate.Before(*startDate) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "end_date must be after start_date",
			})
		}
	}

	discountRate.Name = req.Name
	discountRate.Discount = req.Discount
	discountRate.Periode = req.Periode
	discountRate.Year = yearStr
	discountRate.StartDate = startDate
	discountRate.EndDate = endDate
	discountRate.Description = req.Description

	if err := config.DB.Model(&discountRate).Select("name", "discount", "periode", "year", "start_date", "end_date", "description", "updated_at").Updates(discountRate).Error; err != nil {
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

// GetSalesTransactionDiscountValue godoc
// @Summary Get discount value for a sales transaction
// @Description Calculate discount percentage and amount based on transaction's periode, year, and payment date
// @Tags DiscountRates
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param sales_transaction_id path string true "Sales Transaction ID (UUID)"
// @Param date query string true "Payment transaction date (ISO 8601: YYYY-MM-DD or RFC3339)"
// @Success 200 {object} map[string]interface{} "Discount calculation details"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Transaction or discount rate not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/sales-transactions/{sales_transaction_id}/discount-value [get]
func GetSalesTransactionDiscountValue(c *fiber.Ctx) error {
	salesTransactionID := c.Params("sales_transaction_id")
	dateStr := c.Query("date")

	if dateStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "date query parameter is required",
		})
	}

	paymentDate, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		paymentDate, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid date format. Use ISO 8601 (YYYY-MM-DD or RFC3339)",
			})
		}
	}

	var transaction models.SalesTransaction
	if err := config.DB.Where("id = ?", salesTransactionID).First(&transaction).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Transaction not found",
		})
	}

	if transaction.PaymentType != "K" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Discount calculation is only applicable for credit transactions",
		})
	}

	var discountRate models.DiscountRate
	err = config.DB.Where("periode = ? AND year = ?", transaction.Periode, transaction.Year).
		Where("start_date IS NOT NULL AND end_date IS NOT NULL").
		Where("? BETWEEN start_date AND end_date", paymentDate).
		First(&discountRate).Error

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":        "No applicable discount rate found for this transaction's periode, year, and payment date",
			"periode":      transaction.Periode,
			"year":         transaction.Year,
			"payment_date": paymentDate,
		})
	}

	discountAmount := transaction.TotalAmount * (discountRate.Discount / 100)
	amountAfterDiscount := transaction.TotalAmount - discountAmount

	return c.JSON(fiber.Map{
		"sales_transaction_id": transaction.ID,
		"total_amount":         transaction.TotalAmount,
		"payment_date":         paymentDate,
		"periode":              transaction.Periode,
		"year":                 transaction.Year,
		"discount_rate": fiber.Map{
			"id":                  discountRate.ID,
			"name":                discountRate.Name,
			"discount_percentage": discountRate.Discount,
			"start_date":          discountRate.StartDate,
			"end_date":            discountRate.EndDate,
		},
		"discount_percentage":   discountRate.Discount,
		"discount_amount":       discountAmount,
		"amount_after_discount": amountAfterDiscount,
	})
}
