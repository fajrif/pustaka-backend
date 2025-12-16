package handlers

import (
	"fmt"
	"strconv"
	"time"

	"pustaka-backend/config"
	"pustaka-backend/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// CreatePaymentRequest represents the request body for creating a payment
type CreatePaymentRequest struct {
	PaymentDate time.Time `json:"payment_date"`
	Amount      float64   `json:"amount"`
	Note        *string   `json:"note"`
}

// generatePaymentNumber generates sequential payment number: PMT + YYYYMMDD + 8-digit sequence
// Example: PMT2023120500000001
func generatePaymentNumber(db *gorm.DB) (string, error) {
	prefix := "PMT"
	dateStr := time.Now().Format("20060102") // YYYYMMDD
	pattern := prefix + dateStr + "%"

	var maxNumber string
	err := db.Model(&models.Payment{}).
		Where("no_payment LIKE ?", pattern).
		Select("COALESCE(MAX(no_payment), '')").
		Scan(&maxNumber).Error

	if err != nil {
		return "", err
	}

	nextSeq := 1
	if maxNumber != "" {
		// Extract sequence part (last 8 digits)
		seqStr := maxNumber[len(prefix)+8:] // Skip prefix (3) + date (8)
		if seq, err := strconv.Atoi(seqStr); err == nil {
			nextSeq = seq + 1
		}
	}

	return fmt.Sprintf("%s%s%08d", prefix, dateStr, nextSeq), nil
}

// GetTransactionPayments godoc
// @Summary Get all payments for a transaction
// @Description Retrieve all payments for a specific sales transaction
// @Tags Payments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param transaction_id path string true "Transaction ID (UUID)"
// @Success 200 {object} map[string]interface{} "List of payments with summary"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Transaction not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/sales-transactions/{transaction_id}/payments [get]
func GetTransactionPayments(c *fiber.Ctx) error {
	transactionID := c.Params("transaction_id")

	// Verify transaction exists
	var transaction models.SalesTransaction
	if err := config.DB.Where("id = ?", transactionID).First(&transaction).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Transaction not found",
		})
	}

	var payments []models.Payment
	if err := config.DB.
		Where("sales_transaction_id = ?", transactionID).
		Order("payment_date ASC").
		Find(&payments).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch payments",
		})
	}

	// Calculate total paid
	var totalPaid float64
	for _, payment := range payments {
		totalPaid += payment.Amount
	}

	return c.JSON(fiber.Map{
		"transaction_id":     transaction.ID,
		"total_amount":       transaction.TotalAmount,
		"total_paid":         totalPaid,
		"remaining_amount":   transaction.TotalAmount - totalPaid,
		"transaction_status": transaction.Status,
		"payments":           payments,
	})
}

// CreatePayment godoc
// @Summary Create a new payment for a transaction
// @Description Add a payment to a sales transaction. Validates that total payments don't exceed transaction total.
// @Tags Payments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param transaction_id path string true "Transaction ID (UUID)"
// @Param request body CreatePaymentRequest true "Payment details"
// @Success 201 {object} map[string]interface{} "Created payment with transaction summary"
// @Failure 400 {object} map[string]interface{} "Invalid request body or validation error"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Transaction not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/sales-transactions/{transaction_id}/payments [post]
func CreatePayment(c *fiber.Ctx) error {
	transactionID := c.Params("transaction_id")

	// Verify transaction exists
	var transaction models.SalesTransaction
	if err := config.DB.Where("id = ?", transactionID).First(&transaction).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Transaction not found",
		})
	}

	var req CreatePaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate amount
	if req.Amount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Amount must be greater than 0",
		})
	}

	// Calculate current total payments
	var currentTotalPaid float64
	config.DB.Model(&models.Payment{}).
		Where("sales_transaction_id = ?", transactionID).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&currentTotalPaid)

	// Validate that payment won't exceed total amount
	if currentTotalPaid+req.Amount > transaction.TotalAmount {
		remainingAmount := transaction.TotalAmount - currentTotalPaid
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":            "Payment amount exceeds remaining balance",
			"remaining_amount": remainingAmount,
			"requested_amount": req.Amount,
		})
	}

	// Generate payment number
	noPayment, err := generatePaymentNumber(config.DB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate payment number",
		})
	}

	payment := models.Payment{
		SalesTransactionID: transaction.ID,
		NoPayment:          noPayment,
		PaymentDate:        req.PaymentDate,
		Amount:             req.Amount,
		Note:               req.Note,
	}

	if err := config.DB.Create(&payment).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create payment",
		})
	}

	// Calculate new total paid
	newTotalPaid := currentTotalPaid + req.Amount

	// Update transaction status based on total payments
	var newStatus int
	if newTotalPaid >= transaction.TotalAmount {
		newStatus = 1 // Paid-off
	} else if newTotalPaid > 0 {
		newStatus = 2 // Installment (partial payment)
	} else {
		newStatus = 0 // Booking
	}

	config.DB.Model(&transaction).Update("status", newStatus)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":            "Payment created successfully",
		"payment":            payment,
		"transaction_status": newStatus,
		"total_paid":         newTotalPaid,
		"remaining_amount":   transaction.TotalAmount - newTotalPaid,
	})
}

// DeletePayment godoc
// @Summary Delete a payment
// @Description Delete a payment from a sales transaction. Updates transaction status accordingly.
// @Tags Payments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param transaction_id path string true "Transaction ID (UUID)"
// @Param id path string true "Payment ID (UUID)"
// @Success 200 {object} map[string]interface{} "Payment deleted successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Payment not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/sales-transactions/{transaction_id}/payments/{id} [delete]
func DeletePayment(c *fiber.Ctx) error {
	paymentID := c.Params("id")
	transactionID := c.Params("transaction_id")

	// Verify transaction exists
	var transaction models.SalesTransaction
	if err := config.DB.Where("id = ?", transactionID).First(&transaction).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Transaction not found",
		})
	}

	// Delete the payment
	result := config.DB.
		Where("id = ? AND sales_transaction_id = ?", paymentID, transactionID).
		Delete(&models.Payment{})

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete payment",
		})
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Payment not found",
		})
	}

	// Recalculate total paid and update status
	var totalPaid float64
	config.DB.Model(&models.Payment{}).
		Where("sales_transaction_id = ?", transactionID).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalPaid)

	// Update transaction status
	var newStatus int
	if totalPaid >= transaction.TotalAmount {
		newStatus = 1 // Paid-off
	} else if totalPaid > 0 {
		newStatus = 2 // Installment
	} else {
		newStatus = 0 // Booking
	}

	config.DB.Model(&transaction).Update("status", newStatus)

	return c.JSON(fiber.Map{
		"message":            "Payment deleted successfully",
		"transaction_status": newStatus,
		"total_paid":         totalPaid,
		"remaining_amount":   transaction.TotalAmount - totalPaid,
	})
}
