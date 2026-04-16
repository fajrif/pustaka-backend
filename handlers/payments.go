package handlers

import (
	"fmt"
	"strconv"
	"time"

	"pustaka-backend/config"
	"pustaka-backend/helpers"
	"pustaka-backend/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CreatePaymentRequest struct {
	PaymentDate *string `json:"payment_date" example:"2024-01-15"`
	Amount      float64 `json:"amount" example:"500000.00"`
	Note        *string `json:"note" example:"Payment for invoice INV2024010100000001"`
}

func generatePaymentNumber(db *gorm.DB) (string, error) {
	prefix := "PMT"
	dateStr := time.Now().Format("20060102")
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
		seqStr := maxNumber[len(prefix)+8:]
		if seq, err := strconv.Atoi(seqStr); err == nil {
			nextSeq = seq + 1
		}
	}

	return fmt.Sprintf("%s%s%08d", prefix, dateStr, nextSeq), nil
}

func calculateDiscount(db *gorm.DB, transaction *models.SalesTransaction, paymentDate time.Time) (float64, float64, error) {
	var discountRate models.DiscountRate
	err := db.Where("periode = ? AND year = ?", transaction.Periode, transaction.Year).
		Where("start_date IS NOT NULL AND end_date IS NOT NULL").
		Where("? BETWEEN start_date AND end_date", paymentDate).
		First(&discountRate).Error

	if err != nil {
		return 0, 0, err
	}

	discountPercentage := discountRate.Discount
	discountAmount := transaction.TotalAmount * (discountPercentage / 100)

	return discountPercentage, discountAmount, nil
}

func GetTransactionPayments(c *fiber.Ctx) error {
	transactionID := c.Params("transaction_id")

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

	var totalPaid float64
	var totalDiscount float64
	for _, payment := range payments {
		totalPaid += payment.Amount
		totalDiscount += payment.DiscountAmount
	}

	return c.JSON(fiber.Map{
		"transaction_id":     transaction.ID,
		"total_amount":       transaction.TotalAmount,
		"total_paid":         totalPaid,
		"total_discount":     totalDiscount,
		"remaining_amount":   transaction.TotalAmount - totalPaid - totalDiscount,
		"transaction_status": transaction.Status,
		"payments":           payments,
	})
}

func CreatePayment(c *fiber.Ctx) error {
	transactionID := c.Params("transaction_id")

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

	if req.Amount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Amount must be greater than 0",
		})
	}

	paymentDate, err := helpers.ParseDateString(req.PaymentDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if paymentDate == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "payment_date is required",
		})
	}

	var discountPercentage float64
	var discountAmount float64

	if transaction.PaymentType == "K" {
		dp, da, err := calculateDiscount(config.DB, &transaction, *paymentDate)
		if err == nil {
			discountPercentage = dp
			discountAmount = da
		}
	}

	effectiveAmount := req.Amount - discountAmount
	if effectiveAmount < 0 {
		effectiveAmount = 0
	}

	var totals struct {
		TotalPaid     float64
		TotalDiscount float64
	}
	config.DB.Model(&models.Payment{}).
		Where("sales_transaction_id = ?", transactionID).
		Select("COALESCE(SUM(amount), 0) as total_paid, COALESCE(SUM(discount_amount), 0) as total_discount").
		Scan(&totals)

	currentTotalPaid := totals.TotalPaid
	currentTotalDiscount := totals.TotalDiscount

	maxAllowed := transaction.TotalAmount - currentTotalDiscount

	if req.Amount > maxAllowed {
		remainingAmount := maxAllowed - currentTotalPaid
		if remainingAmount < 0 {
			remainingAmount = 0
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":            "Payment amount exceeds remaining balance",
			"remaining_amount": remainingAmount,
			"requested_amount": req.Amount,
		})
	}

	noPayment, err := generatePaymentNumber(config.DB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate payment number",
		})
	}

	payment := models.Payment{
		SalesTransactionID: transaction.ID,
		NoPayment:          noPayment,
		PaymentDate:        *paymentDate,
		Amount:             req.Amount,
		DiscountPercentage: discountPercentage,
		DiscountAmount:     discountAmount,
		Note:               req.Note,
	}

	if err := config.DB.Create(&payment).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create payment",
		})
	}

	newTotalPaid := currentTotalPaid + req.Amount
	newTotalDiscount := currentTotalDiscount + discountAmount
	newTotalEffective := newTotalPaid - newTotalDiscount

	var newStatus int
	totalCoverage := newTotalEffective
	if transaction.PaymentType == "K" {
		totalCoverage = newTotalPaid + newTotalDiscount
	}

	if totalCoverage >= transaction.TotalAmount {
		newStatus = 1
	} else if newTotalEffective > 0 {
		newStatus = 2
	} else {
		newStatus = 0
	}

	config.DB.Model(&transaction).Update("status", newStatus)

	remainingAmount := transaction.TotalAmount - totalCoverage
	if remainingAmount < 0 {
		remainingAmount = 0
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":               "Payment created successfully",
		"payment":               payment,
		"transaction_status":    newStatus,
		"total_paid":            newTotalPaid,
		"total_discount":        newTotalDiscount,
		"amount_after_discount": newTotalEffective,
		"remaining_amount":      remainingAmount,
	})
}

func DeletePayment(c *fiber.Ctx) error {
	paymentID := c.Params("id")
	transactionID := c.Params("transaction_id")

	var transaction models.SalesTransaction
	if err := config.DB.Where("id = ?", transactionID).First(&transaction).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Transaction not found",
		})
	}

	var payment models.Payment
	if err := config.DB.Where("id = ? AND sales_transaction_id = ?", paymentID, transactionID).First(&payment).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Payment not found",
		})
	}

	if err := config.DB.Delete(&payment).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete payment",
		})
	}

	var deleteTotals struct {
		TotalPaid     float64
		TotalDiscount float64
	}
	config.DB.Model(&models.Payment{}).
		Where("sales_transaction_id = ?", transactionID).
		Select("COALESCE(SUM(amount), 0) as total_paid, COALESCE(SUM(discount_amount), 0) as total_discount").
		Scan(&deleteTotals)

	totalPaid := deleteTotals.TotalPaid
	totalDiscount := deleteTotals.TotalDiscount
	totalEffective := totalPaid - totalDiscount

	var newStatus int
	totalCoverage := totalEffective
	if transaction.PaymentType == "K" {
		totalCoverage = totalPaid + totalDiscount
	}

	if totalCoverage >= transaction.TotalAmount {
		newStatus = 1
	} else if totalEffective > 0 {
		newStatus = 2
	} else {
		newStatus = 0
	}

	config.DB.Model(&transaction).Update("status", newStatus)

	remainingAmount := transaction.TotalAmount - totalCoverage
	if remainingAmount < 0 {
		remainingAmount = 0
	}

	return c.JSON(fiber.Map{
		"message":            "Payment deleted successfully",
		"transaction_status": newStatus,
		"total_paid":         totalPaid,
		"total_discount":     totalDiscount,
		"remaining_amount":   remainingAmount,
	})
}
