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

// CreatePurchaseTransactionRequest represents the request body for creating a purchase transaction
type CreatePurchaseTransactionRequest struct {
	SupplierID   string                             `json:"supplier_id"`
	PurchaseDate time.Time                          `json:"purchase_date"`
	Note         *string                            `json:"note"`
	Items        []CreatePurchaseTransactionItemReq `json:"items"`
}

// CreatePurchaseTransactionItemReq represents an item in the purchase transaction
type CreatePurchaseTransactionItemReq struct {
	BookID   string  `json:"book_id"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

// UpdatePurchaseTransactionRequest represents the request body for updating a purchase transaction
type UpdatePurchaseTransactionRequest struct {
	SupplierID   *string                            `json:"supplier_id"`
	PurchaseDate *time.Time                         `json:"purchase_date"`
	Note         *string                            `json:"note"`
	Items        []CreatePurchaseTransactionItemReq `json:"items,omitempty"`
}

// generatePurchaseInvoiceNumber generates sequential purchase invoice number: PRC + YYYYMMDD + 8-digit sequence
// Example: PRC2023120500000001
func generatePurchaseInvoiceNumber(db *gorm.DB) (string, error) {
	prefix := "PRC"
	dateStr := time.Now().Format("20060102") // YYYYMMDD
	pattern := prefix + dateStr + "%"

	var maxNumber string
	err := db.Model(&models.PurchaseTransaction{}).
		Where("no_invoice LIKE ?", pattern).
		Select("COALESCE(MAX(no_invoice), '')").
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

// GetAllPurchaseTransactions godoc
// @Summary Get all purchase transactions
// @Description Retrieve all purchase transactions with their related entities
// @Tags Purchase Transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param search query string false "Search by invoice number or supplier name"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items per page (default: 20)"
// @Param status query int false "Filter by status (0=pending, 1=completed, 2=cancelled)"
// @Param supplier_id query string false "Filter by supplier ID"
// @Param start_date query string false "Filter by start date (YYYY-MM-DD)"
// @Param end_date query string false "Filter by end date (YYYY-MM-DD)"
// @Success 200 {object} map[string]interface{} "List of all purchase transactions with pagination"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/purchase-transactions [get]
func GetAllPurchaseTransactions(c *fiber.Ctx) error {
	var transactions []models.PurchaseTransaction

	// Get pagination parameters
	pagination := helpers.GetPaginationParams(c)

	query := config.DB.Order("created_at DESC")
	queryCount := config.DB.Model(&models.PurchaseTransaction{})

	// add params for not using pagination
	if c.Query("all") == "true" {
		pagination.Limit = -1 // No limit
		pagination.Offset = 0 // No offset
	}

	// Filter search
	if searchQuery := c.Query("search"); searchQuery != "" {
		searchTerm := "%" + searchQuery + "%"
		cond := "purchase_transactions.no_invoice ILIKE ? OR publishers.name ILIKE ?"
		args := []interface{}{searchTerm, searchTerm}

		query = query.Joins("LEFT JOIN publishers ON purchase_transactions.supplier_id = publishers.id").
			Where(cond, args...)
		queryCount = queryCount.Joins("LEFT JOIN publishers ON purchase_transactions.supplier_id = publishers.id").
			Where(cond, args...)
	}

	// Filter by status
	if status := c.Query("status"); status != "" {
		query = query.Where("purchase_transactions.status = ?", status)
		queryCount = queryCount.Where("purchase_transactions.status = ?", status)
	}

	// Filter by supplier
	if supplierID := c.Query("supplier_id"); supplierID != "" {
		query = query.Where("purchase_transactions.supplier_id = ?", supplierID)
		queryCount = queryCount.Where("purchase_transactions.supplier_id = ?", supplierID)
	}

	// Filter by date range
	if startDate := c.Query("start_date"); startDate != "" {
		query = query.Where("purchase_transactions.purchase_date >= ?", startDate)
		queryCount = queryCount.Where("purchase_transactions.purchase_date >= ?", startDate)
	}

	if endDate := c.Query("end_date"); endDate != "" {
		query = query.Where("purchase_transactions.purchase_date <= ?", endDate+" 23:59:59")
		queryCount = queryCount.Where("purchase_transactions.purchase_date <= ?", endDate+" 23:59:59")
	}

	// Apply pagination and fetch data
	if err := query.
		Offset(pagination.Offset).Limit(pagination.Limit).
		Preload("Supplier").
		Preload("Items").
		Preload("Items.Book").
		Preload("Items.Book.MerkBuku").
		Find(&transactions).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch purchase transactions",
		})
	}

	// Create pagination response
	response, err := helpers.CreatePaginationResponse(queryCount, transactions, "purchase_transactions", pagination.Page, pagination.Limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create pagination response",
		})
	}

	return c.JSON(response)
}

// GetPurchaseTransaction godoc
// @Summary Get a purchase transaction by ID
// @Description Retrieve a single purchase transaction by its ID with all related entities
// @Tags Purchase Transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Transaction ID (UUID)"
// @Success 200 {object} map[string]interface{} "Transaction details"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Transaction not found"
// @Router /api/purchase-transactions/{id} [get]
func GetPurchaseTransaction(c *fiber.Ctx) error {
	id := c.Params("id")

	var transaction models.PurchaseTransaction
	if err := config.DB.
		Preload("Supplier").
		Preload("Items").
		Preload("Items.Book").
		Preload("Items.Book.MerkBuku").
		Preload("Items.Book.JenisBuku").
		Preload("Items.Book.Publisher").
		Where("id = ?", id).First(&transaction).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Purchase transaction not found",
		})
	}

	return c.JSON(fiber.Map{
		"purchase_transaction": transaction,
	})
}

// CreatePurchaseTransaction godoc
// @Summary Create a new purchase transaction
// @Description Create a new purchase transaction with items. Stock is NOT increased until marked as completed.
// @Tags Purchase Transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreatePurchaseTransactionRequest true "Transaction details"
// @Success 201 {object} models.PurchaseTransaction "Created transaction"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/purchase-transactions [post]
func CreatePurchaseTransaction(c *fiber.Ctx) error {
	var req CreatePurchaseTransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if req.SupplierID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "supplier_id is required",
		})
	}

	if len(req.Items) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "At least one item is required",
		})
	}

	// Verify supplier exists
	var supplier models.Publisher
	if err := config.DB.Where("id = ?", req.SupplierID).First(&supplier).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Supplier not found",
		})
	}

	// Start a database transaction
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Calculate total amount from items
	var totalAmount float64
	var transactionItems []models.PurchaseTransactionItem

	for _, item := range req.Items {
		// Fetch book to verify it exists
		var book models.Book
		if err := tx.Where("id = ?", item.BookID).First(&book).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Book with ID %s not found", item.BookID),
			})
		}

		// Validate quantity
		if item.Quantity <= 0 {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Quantity must be greater than 0",
			})
		}

		// Validate price
		if item.Price < 0 {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Price cannot be negative",
			})
		}

		// Calculate subtotal
		subtotal := item.Price * float64(item.Quantity)
		totalAmount += subtotal

		// Create transaction item (we'll save this after creating the transaction)
		transactionItems = append(transactionItems, models.PurchaseTransactionItem{
			BookID:   book.ID,
			Quantity: item.Quantity,
			Price:    item.Price,
			Subtotal: subtotal,
		})
	}

	// Generate invoice number
	noInvoice, err := generatePurchaseInvoiceNumber(tx)
	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate invoice number",
		})
	}

	// Create the transaction with status = 0 (pending)
	// Stock is NOT increased here - only when completed
	transaction := models.PurchaseTransaction{
		SupplierID:   helpers.ParseUUID(req.SupplierID),
		NoInvoice:    noInvoice,
		PurchaseDate: req.PurchaseDate,
		TotalAmount:  totalAmount,
		Status:       models.PurchaseStatusPending, // 0 = pending
		Note:         req.Note,
	}

	// Save the transaction
	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create purchase transaction",
		})
	}

	// Save transaction items
	for i := range transactionItems {
		transactionItems[i].PurchaseTransactionID = transaction.ID
	}
	if err := tx.Create(&transactionItems).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create purchase transaction items",
		})
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to commit transaction",
		})
	}

	// Fetch the complete transaction with all relations
	var createdTransaction models.PurchaseTransaction
	config.DB.
		Preload("Supplier").
		Preload("Items").
		Preload("Items.Book").
		Preload("Items.Book.MerkBuku").
		Where("id = ?", transaction.ID).First(&createdTransaction)

	return c.Status(fiber.StatusCreated).JSON(createdTransaction)
}

// UpdatePurchaseTransaction godoc
// @Summary Update a purchase transaction
// @Description Update an existing purchase transaction by ID. Cannot update completed or cancelled transactions.
// @Tags Purchase Transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Transaction ID (UUID)"
// @Param request body UpdatePurchaseTransactionRequest true "Updated transaction details"
// @Success 200 {object} models.PurchaseTransaction "Updated transaction"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Transaction not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/purchase-transactions/{id} [put]
func UpdatePurchaseTransaction(c *fiber.Ctx) error {
	id := c.Params("id")

	var transaction models.PurchaseTransaction
	if err := config.DB.Where("id = ?", id).First(&transaction).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Purchase transaction not found",
		})
	}

	// Cannot update completed or cancelled transactions
	if transaction.Status != models.PurchaseStatusPending {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot update completed or cancelled purchase transactions",
		})
	}

	var req UpdatePurchaseTransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Start a database transaction
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Build updates map
	updates := make(map[string]interface{})

	if req.SupplierID != nil {
		// Verify supplier exists
		var supplier models.Publisher
		if err := tx.Where("id = ?", *req.SupplierID).First(&supplier).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Supplier not found",
			})
		}
		updates["supplier_id"] = helpers.ParseUUID(*req.SupplierID)
	}

	if req.PurchaseDate != nil {
		updates["purchase_date"] = *req.PurchaseDate
	}

	if req.Note != nil {
		updates["note"] = *req.Note
	}

	// Handle items updates
	if req.Items != nil && len(req.Items) > 0 {
		// Delete existing items
		if err := tx.Where("purchase_transaction_id = ?", transaction.ID).
			Delete(&models.PurchaseTransactionItem{}).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to delete existing items",
			})
		}

		// Create new items
		var totalAmount float64
		var newItems []models.PurchaseTransactionItem

		for _, item := range req.Items {
			// Verify book exists
			var book models.Book
			if err := tx.Where("id = ?", item.BookID).First(&book).Error; err != nil {
				tx.Rollback()
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": fmt.Sprintf("Book with ID %s not found", item.BookID),
				})
			}

			// Validate quantity and price
			if item.Quantity <= 0 {
				tx.Rollback()
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Quantity must be greater than 0",
				})
			}

			if item.Price < 0 {
				tx.Rollback()
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Price cannot be negative",
				})
			}

			subtotal := item.Price * float64(item.Quantity)
			totalAmount += subtotal

			newItems = append(newItems, models.PurchaseTransactionItem{
				PurchaseTransactionID: transaction.ID,
				BookID:                book.ID,
				Quantity:              item.Quantity,
				Price:                 item.Price,
				Subtotal:              subtotal,
			})
		}

		// Create new items
		if err := tx.Create(&newItems).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create new items",
			})
		}

		updates["total_amount"] = totalAmount
	}

	// Apply updates
	if len(updates) > 0 {
		if err := tx.Model(&transaction).Updates(updates).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update purchase transaction",
			})
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to commit transaction",
		})
	}

	// Fetch the updated transaction with all relations
	config.DB.
		Preload("Supplier").
		Preload("Items").
		Preload("Items.Book").
		Preload("Items.Book.MerkBuku").
		Where("id = ?", id).First(&transaction)

	return c.JSON(transaction)
}

// DeletePurchaseTransaction godoc
// @Summary Delete a purchase transaction
// @Description Delete a purchase transaction by ID. If status is completed, stock will be restored.
// @Tags Purchase Transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Transaction ID (UUID)"
// @Success 200 {object} map[string]interface{} "Transaction deleted successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Transaction not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/purchase-transactions/{id} [delete]
func DeletePurchaseTransaction(c *fiber.Ctx) error {
	id := c.Params("id")

	// Verify transaction exists and get items
	var transaction models.PurchaseTransaction
	if err := config.DB.Preload("Items").Where("id = ?", id).First(&transaction).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Purchase transaction not found",
		})
	}

	// Start database transaction
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// If transaction was completed, restore stock (decrease it)
	if transaction.Status == models.PurchaseStatusCompleted {
		for _, item := range transaction.Items {
			var book models.Book
			if err := tx.Where("id = ?", item.BookID).First(&book).Error; err == nil {
				newStock := book.Stock - item.Quantity
				if newStock < 0 {
					newStock = 0 // Prevent negative stock
				}
				if err := tx.Model(&book).Update("stock", newStock).Error; err != nil {
					tx.Rollback()
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"error": "Failed to restore book stock",
					})
				}
			}
		}
	}

	// Delete the transaction (cascade will delete items)
	if err := tx.Delete(&transaction).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete purchase transaction",
		})
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to commit transaction",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Purchase transaction deleted successfully",
	})
}

// CompletePurchaseTransaction godoc
// @Summary Mark a purchase transaction as completed
// @Description Mark a pending purchase transaction as completed. This will increase book stock for all items.
// @Tags Purchase Transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Transaction ID (UUID)"
// @Success 200 {object} models.PurchaseTransaction "Completed transaction"
// @Failure 400 {object} map[string]interface{} "Transaction cannot be completed"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Transaction not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/purchase-transactions/{id}/complete [post]
func CompletePurchaseTransaction(c *fiber.Ctx) error {
	id := c.Params("id")

	// Verify transaction exists and get items
	var transaction models.PurchaseTransaction
	if err := config.DB.Preload("Items").Where("id = ?", id).First(&transaction).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Purchase transaction not found",
		})
	}

	// Check if transaction is pending
	if transaction.Status != models.PurchaseStatusPending {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Only pending transactions can be completed",
		})
	}

	// Start database transaction
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Increase stock for all items
	for _, item := range transaction.Items {
		var book models.Book
		if err := tx.Where("id = ?", item.BookID).First(&book).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Sprintf("Book with ID %s not found", item.BookID.String()),
			})
		}

		// Increase stock
		if err := tx.Model(&book).Update("stock", book.Stock+item.Quantity).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update book stock",
			})
		}
	}

	// Update transaction status to completed
	if err := tx.Model(&transaction).Update("status", models.PurchaseStatusCompleted).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update transaction status",
		})
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to commit transaction",
		})
	}

	// Fetch the updated transaction with all relations
	config.DB.
		Preload("Supplier").
		Preload("Items").
		Preload("Items.Book").
		Preload("Items.Book.MerkBuku").
		Where("id = ?", id).First(&transaction)

	return c.JSON(fiber.Map{
		"message":              "Purchase transaction completed successfully. Stock has been increased.",
		"purchase_transaction": transaction,
	})
}

// CancelPurchaseTransaction godoc
// @Summary Cancel a purchase transaction
// @Description Cancel a pending purchase transaction. Stock is not affected since pending transactions don't affect stock.
// @Tags Purchase Transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Transaction ID (UUID)"
// @Success 200 {object} map[string]interface{} "Transaction cancelled successfully"
// @Failure 400 {object} map[string]interface{} "Transaction cannot be cancelled"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Transaction not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/purchase-transactions/{id}/cancel [post]
func CancelPurchaseTransaction(c *fiber.Ctx) error {
	id := c.Params("id")

	var transaction models.PurchaseTransaction
	if err := config.DB.Where("id = ?", id).First(&transaction).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Purchase transaction not found",
		})
	}

	// Check if transaction is pending
	if transaction.Status != models.PurchaseStatusPending {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Only pending transactions can be cancelled",
		})
	}

	// Update transaction status to cancelled
	if err := config.DB.Model(&transaction).Update("status", models.PurchaseStatusCancelled).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to cancel transaction",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Purchase transaction cancelled successfully",
	})
}

// UploadPurchaseReceipt godoc
// @Summary Upload receipt image for a purchase transaction
// @Description Upload a receipt image for an existing purchase transaction
// @Tags Purchase Transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Transaction ID (UUID)"
// @Param receipt_image_url body string true "Receipt image URL"
// @Success 200 {object} map[string]interface{} "Receipt uploaded successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Transaction not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/purchase-transactions/{id}/receipt [put]
func UploadPurchaseReceipt(c *fiber.Ctx) error {
	id := c.Params("id")

	var transaction models.PurchaseTransaction
	if err := config.DB.Where("id = ?", id).First(&transaction).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Purchase transaction not found",
		})
	}

	var req struct {
		ReceiptImageUrl string `json:"receipt_image_url"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := config.DB.Model(&transaction).Update("receipt_image_url", req.ReceiptImageUrl).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update receipt image",
		})
	}

	return c.JSON(fiber.Map{
		"message":           "Receipt image uploaded successfully",
		"receipt_image_url": req.ReceiptImageUrl,
	})
}
