package handlers

import (
	"pustaka-backend/config"
	"pustaka-backend/helpers"
	"pustaka-backend/models"

	"github.com/gofiber/fiber/v2"
)

// CreateShippingRequest represents the request body for creating a shipping
type CreateShippingRequest struct {
	ExpeditionID string  `json:"expedition_id"`
	NoResi       *string `json:"no_resi"`
	TotalAmount  float64 `json:"total_amount"`
}

// UpdateShippingRequest represents the request body for updating a shipping
type UpdateShippingRequest struct {
	ExpeditionID *string  `json:"expedition_id"`
	NoResi       *string  `json:"no_resi"`
	TotalAmount  *float64 `json:"total_amount"`
}

// GetTransactionShippings godoc
// @Summary Get all shippings for a transaction
// @Description Retrieve all shippings for a specific sales transaction
// @Tags Shippings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param transaction_id path string true "Transaction ID (UUID)"
// @Success 200 {object} map[string]interface{} "List of shippings with total"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Transaction not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/sales-transactions/{transaction_id}/shippings [get]
func GetTransactionShippings(c *fiber.Ctx) error {
	transactionID := c.Params("transaction_id")

	// Verify transaction exists
	var transaction models.SalesTransaction
	if err := config.DB.Where("id = ?", transactionID).First(&transaction).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Transaction not found",
		})
	}

	var shippings []models.Shipping
	if err := config.DB.
		Preload("Expedition").
		Where("sales_transaction_id = ?", transactionID).
		Order("created_at ASC").
		Find(&shippings).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch shippings",
		})
	}

	// Calculate total shipping cost
	var totalShippingCost float64
	for _, shipping := range shippings {
		totalShippingCost += shipping.TotalAmount
	}

	return c.JSON(fiber.Map{
		"transaction_id":      transaction.ID,
		"total_shipping_cost": totalShippingCost,
		"shippings":           shippings,
	})
}

// CreateShipping godoc
// @Summary Create a new shipping for a transaction
// @Description Add a shipping entry to a sales transaction. Automatically updates transaction total.
// @Tags Shippings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param transaction_id path string true "Transaction ID (UUID)"
// @Param request body CreateShippingRequest true "Shipping details"
// @Success 201 {object} map[string]interface{} "Created shipping with updated transaction total"
// @Failure 400 {object} map[string]interface{} "Invalid request body or validation error"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Transaction or expedition not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/sales-transactions/{transaction_id}/shippings [post]
func CreateShipping(c *fiber.Ctx) error {
	transactionID := c.Params("transaction_id")

	// Verify transaction exists
	var transaction models.SalesTransaction
	if err := config.DB.Where("id = ?", transactionID).First(&transaction).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Transaction not found",
		})
	}

	var req CreateShippingRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate expedition_id
	if req.ExpeditionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "expedition_id is required",
		})
	}

	// Verify expedition exists
	var expedition models.Expedition
	if err := config.DB.Where("id = ?", req.ExpeditionID).First(&expedition).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Expedition not found",
		})
	}

	// Validate total_amount
	if req.TotalAmount < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "total_amount cannot be negative",
		})
	}

	// Start database transaction
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create shipping
	shipping := models.Shipping{
		SalesTransactionID: transaction.ID,
		ExpeditionID:       helpers.ParseUUID(req.ExpeditionID),
		NoResi:             req.NoResi,
		TotalAmount:        req.TotalAmount,
	}

	if err := tx.Create(&shipping).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create shipping",
		})
	}

	// Update transaction total_amount (add shipping cost)
	newTotal := transaction.TotalAmount + req.TotalAmount
	if err := tx.Model(&transaction).Update("total_amount", newTotal).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update transaction total",
		})
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to commit transaction",
		})
	}

	// Fetch the shipping with expedition details
	config.DB.Preload("Expedition").Where("id = ?", shipping.ID).First(&shipping)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":                   "Shipping created successfully",
		"shipping":                  shipping,
		"updated_transaction_total": newTotal,
	})
}

// UpdateShipping godoc
// @Summary Update a shipping
// @Description Update a shipping entry. Recalculates transaction total if amount changes.
// @Tags Shippings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param transaction_id path string true "Transaction ID (UUID)"
// @Param id path string true "Shipping ID (UUID)"
// @Param request body UpdateShippingRequest true "Updated shipping details"
// @Success 200 {object} map[string]interface{} "Updated shipping with transaction total"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Shipping or expedition not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/sales-transactions/{transaction_id}/shippings/{id} [put]
func UpdateShipping(c *fiber.Ctx) error {
	shippingID := c.Params("id")
	transactionID := c.Params("transaction_id")

	// Verify transaction exists
	var transaction models.SalesTransaction
	if err := config.DB.Where("id = ?", transactionID).First(&transaction).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Transaction not found",
		})
	}

	// Verify shipping exists and belongs to transaction
	var shipping models.Shipping
	if err := config.DB.
		Where("id = ? AND sales_transaction_id = ?", shippingID, transactionID).
		First(&shipping).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Shipping not found",
		})
	}

	var req UpdateShippingRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Start database transaction
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Build updates map
	updates := make(map[string]interface{})
	oldAmount := shipping.TotalAmount

	if req.ExpeditionID != nil && *req.ExpeditionID != "" {
		// Verify expedition exists
		var expedition models.Expedition
		if err := config.DB.Where("id = ?", *req.ExpeditionID).First(&expedition).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Expedition not found",
			})
		}
		updates["expedition_id"] = helpers.ParseUUID(*req.ExpeditionID)
	}

	if req.NoResi != nil {
		updates["no_resi"] = *req.NoResi
	}

	if req.TotalAmount != nil {
		if *req.TotalAmount < 0 {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "total_amount cannot be negative",
			})
		}
		updates["total_amount"] = *req.TotalAmount
	}

	// Update shipping
	if len(updates) > 0 {
		if err := tx.Model(&shipping).Updates(updates).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update shipping",
			})
		}
	}

	// If total_amount changed, update transaction total
	if req.TotalAmount != nil {
		amountDifference := *req.TotalAmount - oldAmount
		newTotal := transaction.TotalAmount + amountDifference

		if err := tx.Model(&transaction).Update("total_amount", newTotal).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update transaction total",
			})
		}
		transaction.TotalAmount = newTotal
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to commit transaction",
		})
	}

	// Fetch updated shipping with expedition
	config.DB.Preload("Expedition").Where("id = ?", shipping.ID).First(&shipping)

	return c.JSON(fiber.Map{
		"message":                   "Shipping updated successfully",
		"shipping":                  shipping,
		"updated_transaction_total": transaction.TotalAmount,
	})
}

// DeleteShipping godoc
// @Summary Delete a shipping
// @Description Delete a shipping entry from a sales transaction. Subtracts shipping cost from transaction total.
// @Tags Shippings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param transaction_id path string true "Transaction ID (UUID)"
// @Param id path string true "Shipping ID (UUID)"
// @Success 200 {object} map[string]interface{} "Shipping deleted successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Shipping not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/sales-transactions/{transaction_id}/shippings/{id} [delete]
func DeleteShipping(c *fiber.Ctx) error {
	shippingID := c.Params("id")
	transactionID := c.Params("transaction_id")

	// Verify transaction exists
	var transaction models.SalesTransaction
	if err := config.DB.Where("id = ?", transactionID).First(&transaction).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Transaction not found",
		})
	}

	// Verify shipping exists and get its amount
	var shipping models.Shipping
	if err := config.DB.
		Where("id = ? AND sales_transaction_id = ?", shippingID, transactionID).
		First(&shipping).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Shipping not found",
		})
	}

	// Start database transaction
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Delete the shipping
	if err := tx.Delete(&shipping).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete shipping",
		})
	}

	// Update transaction total_amount (subtract shipping cost)
	newTotal := transaction.TotalAmount - shipping.TotalAmount
	if err := tx.Model(&transaction).Update("total_amount", newTotal).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update transaction total",
		})
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to commit transaction",
		})
	}

	return c.JSON(fiber.Map{
		"message":                   "Shipping deleted successfully",
		"updated_transaction_total": newTotal,
	})
}
