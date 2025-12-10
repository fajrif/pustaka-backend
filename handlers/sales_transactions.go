package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"pustaka-backend/config"
	"pustaka-backend/helpers"
	"pustaka-backend/models"

	"github.com/gofiber/fiber/v2"
	// "gorm.io/gorm"
)

// CreateTransactionRequest represents the request body for creating a transaction
type CreateTransactionRequest struct {
	SalesAssociateID string                        `json:"sales_associate_id"`
	ExpeditionID     *string                       `json:"expedition_id"`
	PaymentType      string                        `json:"payment_type"` // 'T' or 'K'
	TransactionDate  time.Time                     `json:"transaction_date"`
	DueDate          *time.Time                    `json:"due_date"`
	ExpeditionPrice  float64                       `json:"expedition_price"`
	Items            []CreateTransactionItemRequest `json:"items"`
}

// CreateTransactionItemRequest represents an item in the transaction
type CreateTransactionItemRequest struct {
	BookID   string `json:"book_id"`
	Quantity int    `json:"quantity"`
}

// CreateInstallmentRequest represents an installment payment
type CreateInstallmentRequest struct {
	InstallmentDate time.Time `json:"installment_date"`
	Amount          float64   `json:"amount"`
	Note            *string   `json:"note"`
}

// generateInvoiceNumber generates a unique invoice number in format: JL#timestamp#random_hex
func generateInvoiceNumber() string {
	timestamp := time.Now().Unix()
	randomBytes := make([]byte, 4)
	rand.Read(randomBytes)
	randomHex := hex.EncodeToString(randomBytes)
	return fmt.Sprintf("JL%d%s", timestamp, randomHex)
}

// GetAllSalesTransactions godoc
// @Summary Get all sales transactions
// @Description Retrieve all sales transactions with their related entities
// @Tags Sales Transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param search query string false "Search by invoice number or sales associate name"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items per page (default: 20)"
// @Param status query int false "Filter by status (0=booking, 1=paid-off, 2=installment)"
// @Param payment_type query string false "Filter by payment type (T=cash, K=credit)"
// @Success 200 {object} map[string]interface{} "List of all sales transactions with pagination"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/sales-transactions [get]
func GetAllSalesTransactions(c *fiber.Ctx) error {
	var transactions []models.SalesTransaction

	// Get pagination parameters
	pagination := helpers.GetPaginationParams(c)

	query := config.DB.Order("created_at DESC")
	queryCount := config.DB.Model(&models.SalesTransaction{})

	// add params for not using pagination
	if c.Query("all") == "true" {
		pagination.Limit = -1 // No limit
		pagination.Offset = 0 // No offset
	}

	// Filter search
	if searchQuery := c.Query("search"); searchQuery != "" {
		searchTerm := "%" + searchQuery + "%"
		cond := "sales_transactions.no_invoice ILIKE ? OR sales_associates.name ILIKE ?"
		args := []interface{}{searchTerm, searchTerm}

		query = query.Joins("LEFT JOIN sales_associates ON sales_transactions.sales_associate_id = sales_associates.id").
			Where(cond, args...)
		queryCount = queryCount.Joins("LEFT JOIN sales_associates ON sales_transactions.sales_associate_id = sales_associates.id").
			Where(cond, args...)
	}

	// Filter by status
	if status := c.Query("status"); status != "" {
		query = query.Where("sales_transactions.status = ?", status)
		queryCount = queryCount.Where("sales_transactions.status = ?", status)
	}

	// Filter by payment type
	if paymentType := c.Query("payment_type"); paymentType != "" {
		query = query.Where("sales_transactions.payment_type = ?", paymentType)
		queryCount = queryCount.Where("sales_transactions.payment_type = ?", paymentType)
	}

	// Apply pagination and fetch data
	if err := query.
		Offset(pagination.Offset).Limit(pagination.Limit).
		Preload("SalesAssociate").
		Preload("SalesAssociate.City").
		Preload("Expedition").
		Preload("Items").
		Preload("Items.Book").
		Preload("Installments").
		Find(&transactions).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch sales transactions",
		})
	}

	// Create pagination response
	response, err := helpers.CreatePaginationResponse(queryCount, transactions, "sales_transactions", pagination.Page, pagination.Limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create pagination response",
		})
	}

	return c.JSON(response)
}

// GetSalesTransaction godoc
// @Summary Get a sales transaction by ID
// @Description Retrieve a single sales transaction by its ID with all related entities
// @Tags Sales Transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Transaction ID (UUID)"
// @Success 200 {object} map[string]interface{} "Transaction details"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Transaction not found"
// @Router /api/sales-transactions/{id} [get]
func GetSalesTransaction(c *fiber.Ctx) error {
	id := c.Params("id")

	var transaction models.SalesTransaction
	if err := config.DB.
		Preload("SalesAssociate").
		Preload("SalesAssociate.City").
		Preload("Expedition").
		Preload("Items").
		Preload("Items.Book").
		Preload("Items.Book.Publisher").
		Preload("Installments").
		Where("id = ?", id).First(&transaction).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Transaction not found",
		})
	}

	return c.JSON(fiber.Map{
		"transaction": transaction,
	})
}

// CreateSalesTransaction godoc
// @Summary Create a new sales transaction
// @Description Create a new sales transaction with items and optional installments
// @Tags Sales Transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateTransactionRequest true "Transaction details"
// @Success 201 {object} models.SalesTransaction "Created transaction"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/sales-transactions [post]
func CreateSalesTransaction(c *fiber.Ctx) error {
	var req CreateTransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if req.SalesAssociateID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "sales_associate_id is required",
		})
	}

	if len(req.Items) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "At least one item is required",
		})
	}

	// Validate payment type
	if req.PaymentType != "T" && req.PaymentType != "K" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "payment_type must be either 'T' (cash) or 'K' (credit)",
		})
	}

	// If payment type is credit, due_date is required
	if req.PaymentType == "K" && req.DueDate == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "due_date is required for credit payment",
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
	var totalItemsPrice float64
	var transactionItems []models.SalesTransactionItem

	for _, item := range req.Items {
		// Fetch book to get current price
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

		// Calculate subtotal
		subtotal := book.Price * float64(item.Quantity)
		totalItemsPrice += subtotal

		// Create transaction item (we'll save this after creating the transaction)
		transactionItems = append(transactionItems, models.SalesTransactionItem{
			BookID:   book.ID,
			Quantity: item.Quantity,
			Price:    book.Price,
			Subtotal: subtotal,
		})
	}

	// Calculate total amount (items + expedition)
	totalAmount := totalItemsPrice + req.ExpeditionPrice

	// Create the transaction
	transaction := models.SalesTransaction{
		SalesAssociateID: helpers.ParseUUID(req.SalesAssociateID),
		NoInvoice:        generateInvoiceNumber(),
		PaymentType:      req.PaymentType,
		TransactionDate:  req.TransactionDate,
		DueDate:          req.DueDate,
		ExpeditionPrice:  req.ExpeditionPrice,
		TotalAmount:      totalAmount,
		Status:           0, // Default to booking
	}

	// Set expedition if provided
	if req.ExpeditionID != nil && *req.ExpeditionID != "" {
		expeditionUUID := helpers.ParseUUID(*req.ExpeditionID)
		transaction.ExpeditionID = &expeditionUUID
	}

	// Save the transaction
	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create transaction",
		})
	}

	// Save transaction items
	for i := range transactionItems {
		transactionItems[i].TransactionID = transaction.ID
	}
	if err := tx.Create(&transactionItems).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create transaction items",
		})
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to commit transaction",
		})
	}

	// Fetch the complete transaction with all relations
	var createdTransaction models.SalesTransaction
	config.DB.
		Preload("SalesAssociate").
		Preload("Expedition").
		Preload("Items").
		Preload("Items.Book").
		Preload("Installments").
		Where("id = ?", transaction.ID).First(&createdTransaction)

	return c.Status(fiber.StatusCreated).JSON(createdTransaction)
}

// UpdateTransactionRequest represents the request body for updating a transaction
type UpdateTransactionRequest struct {
	SalesAssociateID *string                        `json:"sales_associate_id"`
	ExpeditionID     *string                        `json:"expedition_id"`
	PaymentType      *string                        `json:"payment_type"`
	TransactionDate  *time.Time                     `json:"transaction_date"`
	DueDate          *time.Time                     `json:"due_date"`
	ExpeditionPrice  *float64                       `json:"expedition_price"`
	TotalAmount      *float64                       `json:"total_amount"`
	Status           *int                           `json:"status"`
	Items            []CreateTransactionItemRequest `json:"items,omitempty"`
}

// UpdateSalesTransaction godoc
// @Summary Update a sales transaction
// @Description Update an existing sales transaction by ID
// @Tags Sales Transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Transaction ID (UUID)"
// @Param request body UpdateTransactionRequest true "Updated transaction details"
// @Success 200 {object} models.SalesTransaction "Updated transaction"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Transaction not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/sales-transactions/{id} [put]
func UpdateSalesTransaction(c *fiber.Ctx) error {
	id := c.Params("id")

	var transaction models.SalesTransaction
	if err := config.DB.Where("id = ?", id).First(&transaction).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Transaction not found",
		})
	}

	var req UpdateTransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Start a database transaction for atomic updates
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Build updates map to handle zero values and nil properly
	updates := make(map[string]interface{})

	if req.SalesAssociateID != nil {
		updates["sales_associate_id"] = helpers.ParseUUID(*req.SalesAssociateID)
	}

	// Handle expedition_id - can be set to nil
	if req.ExpeditionID != nil {
		if *req.ExpeditionID == "" {
			updates["expedition_id"] = nil
		} else {
			expeditionUUID := helpers.ParseUUID(*req.ExpeditionID)
			updates["expedition_id"] = expeditionUUID
		}
	}

	if req.PaymentType != nil {
		// Validate payment type
		if *req.PaymentType != "T" && *req.PaymentType != "K" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "payment_type must be either 'T' (cash) or 'K' (credit)",
			})
		}
		updates["payment_type"] = *req.PaymentType
	}

	if req.TransactionDate != nil {
		updates["transaction_date"] = *req.TransactionDate
	}

	// Handle due_date - can be set to nil
	if req.DueDate != nil {
		updates["due_date"] = *req.DueDate
	}

	// Handle expedition_price - can be set to 0
	if req.ExpeditionPrice != nil {
		updates["expedition_price"] = *req.ExpeditionPrice
	}

	if req.Status != nil {
		updates["status"] = *req.Status
	}

	// Handle items updates
	if req.Items != nil && len(req.Items) > 0 {
		// Get existing items for this transaction
		var existingItems []models.SalesTransactionItem
		if err := tx.Where("transaction_id = ?", transaction.ID).Find(&existingItems).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch existing items",
			})
		}

		// Create a map of existing items by book_id for quick lookup
		existingItemsMap := make(map[string]models.SalesTransactionItem)
		for _, item := range existingItems {
			existingItemsMap[item.BookID.String()] = item
		}

		// Track which book IDs are in the update request
		requestedBookIDs := make(map[string]bool)
		var totalItemsPrice float64

		// Process each item in the request
		for _, itemReq := range req.Items {
			// Check for duplicate book_id in the request
			if requestedBookIDs[itemReq.BookID] {
				tx.Rollback()
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": fmt.Sprintf("Duplicate book_id in request: %s", itemReq.BookID),
				})
			}
			requestedBookIDs[itemReq.BookID] = true

			// Fetch book to get current price
			var book models.Book
			if err := tx.Where("id = ?", itemReq.BookID).First(&book).Error; err != nil {
				tx.Rollback()
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": fmt.Sprintf("Book with ID %s not found", itemReq.BookID),
				})
			}

			// Validate quantity
			if itemReq.Quantity <= 0 {
				tx.Rollback()
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Quantity must be greater than 0",
				})
			}

			// Calculate subtotal
			subtotal := book.Price * float64(itemReq.Quantity)
			totalItemsPrice += subtotal

			// Check if this book_id already exists in the transaction
			if existingItem, exists := existingItemsMap[itemReq.BookID]; exists {
				// Update existing item
				if err := tx.Model(&existingItem).Updates(map[string]interface{}{
					"quantity": itemReq.Quantity,
					"price":    book.Price,
					"subtotal": subtotal,
				}).Error; err != nil {
					tx.Rollback()
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"error": "Failed to update transaction item",
					})
				}
			} else {
				// Create new item
				newItem := models.SalesTransactionItem{
					TransactionID: transaction.ID,
					BookID:        book.ID,
					Quantity:      itemReq.Quantity,
					Price:         book.Price,
					Subtotal:      subtotal,
				}
				if err := tx.Create(&newItem).Error; err != nil {
					tx.Rollback()
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"error": "Failed to create transaction item",
					})
				}
			}
		}

		// Delete items that are no longer in the request
		for bookID, existingItem := range existingItemsMap {
			if !requestedBookIDs[bookID] {
				if err := tx.Delete(&existingItem).Error; err != nil {
					tx.Rollback()
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"error": "Failed to delete transaction item",
					})
				}
			}
		}

		// Recalculate total amount
		expeditionPrice := transaction.ExpeditionPrice
		if req.ExpeditionPrice != nil {
			expeditionPrice = *req.ExpeditionPrice
		}
		updates["total_amount"] = totalItemsPrice + expeditionPrice
	} else if req.ExpeditionPrice != nil {
		// If only expedition price changed (no items update), recalculate total
		// Fetch current items total
		var currentItemsTotal float64
		tx.Model(&models.SalesTransactionItem{}).
			Where("transaction_id = ?", transaction.ID).
			Select("COALESCE(SUM(subtotal), 0)").
			Scan(&currentItemsTotal)

		updates["total_amount"] = currentItemsTotal + *req.ExpeditionPrice
	}

	// Update using map to handle zero values
	if len(updates) > 0 {
		if err := tx.Model(&transaction).Updates(updates).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update transaction",
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
		Preload("SalesAssociate").
		Preload("Expedition").
		Preload("Items").
		Preload("Items.Book").
		Preload("Installments").
		Where("id = ?", id).First(&transaction)

	return c.JSON(transaction)
}

// DeleteSalesTransaction godoc
// @Summary Delete a sales transaction
// @Description Delete a sales transaction by ID (this will cascade delete items and installments)
// @Tags Sales Transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Transaction ID (UUID)"
// @Success 200 {object} map[string]interface{} "Transaction deleted successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Transaction not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/sales-transactions/{id} [delete]
func DeleteSalesTransaction(c *fiber.Ctx) error {
	id := c.Params("id")

	result := config.DB.Delete(&models.SalesTransaction{}, "id = ?", id)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete transaction",
		})
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Transaction not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Transaction deleted successfully",
	})
}

// AddInstallment godoc
// @Summary Add an installment to an existing transaction
// @Description Add a new installment payment to a credit transaction
// @Tags Sales Transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Transaction ID (UUID)"
// @Param request body CreateInstallmentRequest true "Installment details"
// @Success 201 {object} models.SalesTransactionInstallment "Created installment"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Transaction not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/sales-transactions/{id}/installments [post]
func AddInstallment(c *fiber.Ctx) error {
	id := c.Params("id")

	// Verify transaction exists and is a credit transaction
	var transaction models.SalesTransaction
	if err := config.DB.Where("id = ?", id).First(&transaction).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Transaction not found",
		})
	}

	if transaction.PaymentType != "K" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Installments can only be added to credit transactions",
		})
	}

	var req CreateInstallmentRequest
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

	installment := models.SalesTransactionInstallment{
		TransactionID:   transaction.ID,
		InstallmentDate: req.InstallmentDate,
		Amount:          req.Amount,
		Note:            req.Note,
	}

	if err := config.DB.Create(&installment).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create installment",
		})
	}

	// Check if total installments equal total amount, update status if paid off
	var totalInstallments float64
	config.DB.Model(&models.SalesTransactionInstallment{}).
		Where("transaction_id = ?", transaction.ID).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalInstallments)

	if totalInstallments >= transaction.TotalAmount {
		config.DB.Model(&transaction).Update("status", 1) // Paid-off
	} else {
		config.DB.Model(&transaction).Update("status", 2) // Installment
	}

	return c.Status(fiber.StatusCreated).JSON(installment)
}

// GetTransactionInstallments godoc
// @Summary Get all installments for a transaction
// @Description Retrieve all installment payments for a specific transaction
// @Tags Sales Transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Transaction ID (UUID)"
// @Success 200 {object} map[string]interface{} "List of installments"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Transaction not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/sales-transactions/{id}/installments [get]
func GetTransactionInstallments(c *fiber.Ctx) error {
	id := c.Params("id")

	// Verify transaction exists
	var transaction models.SalesTransaction
	if err := config.DB.Where("id = ?", id).First(&transaction).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Transaction not found",
		})
	}

	var installments []models.SalesTransactionInstallment
	if err := config.DB.
		Where("transaction_id = ?", id).
		Order("installment_date ASC").
		Find(&installments).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch installments",
		})
	}

	// Calculate total paid
	var totalPaid float64
	for _, inst := range installments {
		totalPaid += inst.Amount
	}

	return c.JSON(fiber.Map{
		"transaction_id": transaction.ID,
		"total_amount":   transaction.TotalAmount,
		"total_paid":     totalPaid,
		"remaining":      transaction.TotalAmount - totalPaid,
		"installments":   installments,
	})
}
