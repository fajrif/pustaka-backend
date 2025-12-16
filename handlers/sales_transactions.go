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

// CreateTransactionRequest represents the request body for creating a transaction
type CreateTransactionRequest struct {
	SalesAssociateID string                         `json:"sales_associate_id"`
	PaymentType      string                         `json:"payment_type"` // 'T' or 'K'
	TransactionDate  time.Time                      `json:"transaction_date"`
	DueDate          *time.Time                     `json:"due_date"`
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

// generateInvoiceNumber generates sequential invoice number: INV + YYYYMMDD + 8-digit sequence
// Example: INV2023120500000001
func generateInvoiceNumber(db *gorm.DB) (string, error) {
	prefix := "INV"
	dateStr := time.Now().Format("20060102") // YYYYMMDD
	pattern := prefix + dateStr + "%"

	var maxNumber string
	err := db.Model(&models.SalesTransaction{}).
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

// generateInstallmentNumber generates sequential installment number: PKR + YYYYMMDD + 8-digit sequence
// Example: PKR2023120500000001
func generateInstallmentNumber(db *gorm.DB) (string, error) {
	prefix := "PKR"
	dateStr := time.Now().Format("20060102") // YYYYMMDD
	pattern := prefix + dateStr + "%"

	var maxNumber string
	err := db.Model(&models.SalesTransactionInstallment{}).
		Where("no_installment LIKE ?", pattern).
		Select("COALESCE(MAX(no_installment), '')").
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
		Preload("Biller").
		Preload("SalesAssociate").
		Preload("SalesAssociate.City").
		Preload("Items").
		Preload("Items.Book").
		Preload("Payments").
		Preload("Shippings").
		Preload("Shippings.Expedition").
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
		Preload("Biller").
		Preload("SalesAssociate").
		Preload("SalesAssociate.City").
		Preload("Items").
		Preload("Items.Book").
		Preload("Items.Book.Publisher").
		Preload("Items.Book.JenisBuku").
		Preload("Payments").
		Preload("Shippings").
		Preload("Shippings.Expedition").
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
	// Get Default Biller ID
	var defaultBiller models.Biller
	if err := config.DB.Select("id").First(&defaultBiller).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get default biller",
		})
	}

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

	// Calculate total amount from items and validate stock
	var totalItemsPrice float64
	var transactionItems []models.SalesTransactionItem
	var booksToUpdate []models.Book

	for _, item := range req.Items {
		// Fetch book to get current price and stock
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

		// Validate stock availability
		if book.Stock < item.Quantity {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":           fmt.Sprintf("Insufficient stock for book: %s", book.Name),
				"available_stock": book.Stock,
				"requested":       item.Quantity,
			})
		}

		// Calculate subtotal
		subtotal := book.Price * float64(item.Quantity)
		totalItemsPrice += subtotal

		// Reduce stock
		book.Stock -= item.Quantity
		booksToUpdate = append(booksToUpdate, book)

		// Create transaction item (we'll save this after creating the transaction)
		transactionItems = append(transactionItems, models.SalesTransactionItem{
			BookID:   book.ID,
			Quantity: item.Quantity,
			Price:    book.Price,
			Subtotal: subtotal,
		})
	}

	// Update book stocks
	for _, book := range booksToUpdate {
		if err := tx.Model(&book).Update("stock", book.Stock).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update book stock",
			})
		}
	}

	// Calculate total amount (items only, shipping added separately)
	totalAmount := totalItemsPrice

	// Generate invoice number
	noInvoice, err := generateInvoiceNumber(tx)
	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate invoice number",
		})
	}

	// Create the transaction
	transaction := models.SalesTransaction{
		BillerID:         &defaultBiller.ID,
		SalesAssociateID: helpers.ParseUUID(req.SalesAssociateID),
		NoInvoice:        noInvoice,
		PaymentType:      req.PaymentType,
		TransactionDate:  req.TransactionDate,
		DueDate:          req.DueDate,
		TotalAmount:      totalAmount,
		Status:           0, // Default to booking
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
		Preload("Biller").
		Preload("SalesAssociate").
		Preload("Items").
		Preload("Items.Book").
		Preload("Payments").
		Preload("Shippings").
		Preload("Shippings.Expedition").
		Where("id = ?", transaction.ID).First(&createdTransaction)

	return c.Status(fiber.StatusCreated).JSON(createdTransaction)
}

// UpdateTransactionRequest represents the request body for updating a transaction
type UpdateTransactionRequest struct {
	SalesAssociateID *string                        `json:"sales_associate_id"`
	PaymentType      *string                        `json:"payment_type"`
	TransactionDate  *time.Time                     `json:"transaction_date"`
	DueDate          *time.Time                     `json:"due_date"`
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

	// check existing transaction biller id is nil
	// if nil, then get default biller id from biller table (first record)
	// set transaction biller id to default biller id
	if transaction.BillerID == nil {
		var defaultBiller models.Biller
		if err := config.DB.Select("id").First(&defaultBiller).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get default biller",
			})
		}
		transaction.BillerID = &defaultBiller.ID
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

	if req.PaymentType != nil {
		// Validate payment type
		if *req.PaymentType != "T" && *req.PaymentType != "K" {
			tx.Rollback()
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

	if req.Status != nil {
		updates["status"] = *req.Status
	}

	// Handle items updates with stock management
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

			// Fetch book to get current price and stock
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
				// Calculate stock adjustment (difference between old and new quantity)
				quantityDiff := itemReq.Quantity - existingItem.Quantity

				if quantityDiff > 0 {
					// Need more stock - check availability
					if book.Stock < quantityDiff {
						tx.Rollback()
						return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
							"error":           fmt.Sprintf("Insufficient stock for book: %s", book.Name),
							"available_stock": book.Stock,
							"additional_requested": quantityDiff,
						})
					}
					// Reduce stock
					if err := tx.Model(&book).Update("stock", book.Stock-quantityDiff).Error; err != nil {
						tx.Rollback()
						return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
							"error": "Failed to update book stock",
						})
					}
				} else if quantityDiff < 0 {
					// Returning stock
					if err := tx.Model(&book).Update("stock", book.Stock-quantityDiff).Error; err != nil {
						tx.Rollback()
						return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
							"error": "Failed to update book stock",
						})
					}
				}

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
				// New item - check stock availability
				if book.Stock < itemReq.Quantity {
					tx.Rollback()
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"error":           fmt.Sprintf("Insufficient stock for book: %s", book.Name),
						"available_stock": book.Stock,
						"requested":       itemReq.Quantity,
					})
				}

				// Reduce stock
				if err := tx.Model(&book).Update("stock", book.Stock-itemReq.Quantity).Error; err != nil {
					tx.Rollback()
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"error": "Failed to update book stock",
					})
				}

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

		// Delete items that are no longer in the request and restore stock
		for bookID, existingItem := range existingItemsMap {
			if !requestedBookIDs[bookID] {
				// Restore stock for removed item
				var book models.Book
				if err := tx.Where("id = ?", existingItem.BookID).First(&book).Error; err == nil {
					if err := tx.Model(&book).Update("stock", book.Stock+existingItem.Quantity).Error; err != nil {
						tx.Rollback()
						return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
							"error": "Failed to restore book stock",
						})
					}
				}

				if err := tx.Delete(&existingItem).Error; err != nil {
					tx.Rollback()
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"error": "Failed to delete transaction item",
					})
				}
			}
		}

		// Recalculate total amount (items + existing shipping costs)
		var totalShippingCost float64
		tx.Model(&models.Shipping{}).
			Where("sales_transaction_id = ?", transaction.ID).
			Select("COALESCE(SUM(total_amount), 0)").
			Scan(&totalShippingCost)

		updates["total_amount"] = totalItemsPrice + totalShippingCost
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
		Preload("Biller").
		Preload("SalesAssociate").
		Preload("Items").
		Preload("Items.Book").
		Preload("Payments").
		Preload("Shippings").
		Preload("Shippings.Expedition").
		Where("id = ?", id).First(&transaction)

	return c.JSON(transaction)
}

// DeleteSalesTransaction godoc
// @Summary Delete a sales transaction
// @Description Delete a sales transaction by ID (this will cascade delete items, payments, and shippings). Stock is restored for all items.
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

	// Verify transaction exists and get items for stock restoration
	var transaction models.SalesTransaction
	if err := config.DB.Preload("Items").Where("id = ?", id).First(&transaction).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Transaction not found",
		})
	}

	// Start database transaction
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Restore stock for all items
	for _, item := range transaction.Items {
		var book models.Book
		if err := tx.Where("id = ?", item.BookID).First(&book).Error; err == nil {
			if err := tx.Model(&book).Update("stock", book.Stock+item.Quantity).Error; err != nil {
				tx.Rollback()
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to restore book stock",
				})
			}
		}
	}

	// Delete the transaction (cascade will delete items, payments, shippings)
	if err := tx.Delete(&transaction).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete transaction",
		})
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to commit transaction",
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
// @Param transaction_id path string true "Transaction ID (UUID)"
// @Param request body CreateInstallmentRequest true "Installment details"
// @Success 201 {object} models.SalesTransactionInstallment "Created installment"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Transaction not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/sales-transactions/{transaction_id}/installments [post]
func AddInstallment(c *fiber.Ctx) error {
	id := c.Params("transaction_id")

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

	// Generate installment number
	noInstallment, err := generateInstallmentNumber(config.DB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate installment number",
		})
	}

	installment := models.SalesTransactionInstallment{
		TransactionID:   transaction.ID,
		NoInstallment:   noInstallment,
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
// @Param transaction_id path string true "Transaction ID (UUID)"
// @Success 200 {object} map[string]interface{} "List of installments"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Transaction not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/sales-transactions/{transaction_id}/installments [get]
func GetTransactionInstallments(c *fiber.Ctx) error {
	id := c.Params("transaction_id")

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

// DeleteInstallment godoc
// @Summary Delete a sales transaction installment
// @Description Delete a sales transaction installment by ID
// @Tags Sales Transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param transaction_id path string true "Transaction ID (UUID)"
// @Success 200 {object} map[string]interface{} "Transaction deleted successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Transaction not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/sales-transactions/{transaction_id}/installments/{id} [delete]
func DeleteInstallment(c *fiber.Ctx) error {
	id := c.Params("id")
	transactionID := c.Params("transaction_id")

	result := config.DB.
		Where("id = ? AND sales_transaction_id = ?", id, transactionID).
		Delete(&models.SalesTransactionInstallment{})
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete installment",
		})
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Installment not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Installment deleted successfully",
	})
}
