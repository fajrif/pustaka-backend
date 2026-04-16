package handlers

import (
	"pustaka-backend/config"
	"pustaka-backend/helpers"
	"pustaka-backend/models"

	"github.com/gofiber/fiber/v2"
)

// PurchasingReportSummary represents the summary for purchasing report
type PurchasingReportSummary struct {
	TotalTransactions int     `json:"total_transactions"`
	TotalAmount       float64 `json:"total_amount"`
	TotalItems        int     `json:"total_items"`
}

// SalesReportSummary represents the summary for sales report
type SalesReportSummary struct {
	TotalTransactions  int     `json:"total_transactions"`
	TotalAmount        float64 `json:"total_amount"`
	TotalItems         int     `json:"total_items"`
	CashTransactions   int     `json:"cash_transactions"`
	CashAmount         float64 `json:"cash_amount"`
	CreditTransactions int     `json:"credit_transactions"`
	CreditAmount       float64 `json:"credit_amount"`
}

// BooksStockSummary represents the summary for books stock report
type BooksStockSummary struct {
	TotalBooks    int `json:"total_books"`
	TotalStock    int `json:"total_stock"`
	LowStockCount int `json:"low_stock_count"`
}

// CreditReportSummary represents the summary for credits report
type CreditReportSummary struct {
	TotalOutstanding  float64 `json:"total_outstanding"`
	TotalTransactions int     `json:"total_transactions"`
	TotalItems        int     `json:"total_items"`
}

// PurchasingReportData represents a single purchase transaction with computed fields
type PurchasingReportData struct {
	models.PurchaseTransaction
	TotalItems int `json:"total_items"`
}

// SalesReportData represents a single sales transaction with computed fields
type SalesReportData struct {
	models.SalesTransaction
	TotalItems int `json:"total_items"`
}

// CreditReportItem represents a single item in the credits report
type CreditReportItem struct {
	Transaction     models.SalesTransaction `json:"transaction"`
	TotalPaid       float64                 `json:"total_paid"`
	RemainingAmount float64                 `json:"remaining_amount"`
	TotalItems      int                     `json:"total_items"`
}

// GetPurchasingReport godoc
// @Summary Get purchasing report
// @Description Get a report of all purchases from suppliers with filters
// @Tags Reports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items per page (default: 20)"
// @Param all query bool false "Get all records without pagination"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param supplier_id query string false "Filter by supplier ID"
// @Param status query int false "Filter by status (0=pending, 1=completed, 2=cancelled)"
// @Success 200 {object} map[string]interface{} "Purchasing report with summary and pagination"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/reports/purchases [get]
func GetPurchasingReport(c *fiber.Ctx) error {
	var transactions []models.PurchaseTransaction

	// Get pagination parameters
	pagination := helpers.GetPaginationParams(c)

	query := config.DB.Order("purchase_date DESC").
		Preload("Supplier").
		Preload("Items").
		Preload("Items.Book").
		Preload("Items.Book.MerkBuku")

	queryCount := config.DB.Model(&models.PurchaseTransaction{})

	// add params for not using pagination
	if c.Query("all") == "true" {
		pagination.Limit = -1 // No limit
		pagination.Offset = 0 // No offset
	}

	// Filter by date range
	if startDate := c.Query("start_date"); startDate != "" {
		query = query.Where("purchase_date >= ?", startDate)
		queryCount = queryCount.Where("purchase_date >= ?", startDate)
	}

	if endDate := c.Query("end_date"); endDate != "" {
		query = query.Where("purchase_date <= ?", endDate+" 23:59:59")
		queryCount = queryCount.Where("purchase_date <= ?", endDate+" 23:59:59")
	}

	// Filter by supplier
	if supplierID := c.Query("supplier_id"); supplierID != "" {
		query = query.Where("supplier_id = ?", supplierID)
		queryCount = queryCount.Where("supplier_id = ?", supplierID)
	}

	// Filter by status
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
		queryCount = queryCount.Where("status = ?", status)
	}

	// Apply pagination and fetch data
	if err := query.Offset(pagination.Offset).Limit(pagination.Limit).Find(&transactions).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch purchasing report",
		})
	}

	// Build report data with total_items for each transaction
	var reportData []PurchasingReportData
	var summary PurchasingReportSummary
	summary.TotalTransactions = len(transactions)

	for _, tx := range transactions {
		summary.TotalAmount += tx.TotalAmount

		// Calculate total items for this transaction
		totalItems := 0
		for _, item := range tx.Items {
			totalItems += item.Quantity
			summary.TotalItems += item.Quantity
		}

		reportData = append(reportData, PurchasingReportData{
			PurchaseTransaction: tx,
			TotalItems:          totalItems,
		})
	}

	// Create pagination response
	response, err := helpers.CreatePaginationResponse(queryCount, reportData, "data", pagination.Page, pagination.Limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create pagination response",
		})
	}

	// Add summary to response
	response["summary"] = summary

	return c.JSON(response)
}

// GetSalesReport godoc
// @Summary Get sales report
// @Description Get a report of all sales transactions with filters
// @Tags Reports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items per page (default: 20)"
// @Param all query bool false "Get all records without pagination"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param payment_type query string false "Filter by payment type (T=cash, K=credit, all=both)"
// @Param status query int false "Filter by status (0=booking, 1=paid-off, 2=installment)"
// @Param sales_associate_id query string false "Filter by sales associate ID"
// @Success 200 {object} map[string]interface{} "Sales report with summary and pagination"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/reports/sales [get]
func GetSalesReport(c *fiber.Ctx) error {
	var transactions []models.SalesTransaction

	// Get pagination parameters
	pagination := helpers.GetPaginationParams(c)

	query := config.DB.Order("transaction_date DESC").
		Preload("Biller").
		Preload("SalesAssociate").
		Preload("Items").
		Preload("Items.Book").
		Preload("Items.Book.MerkBuku").
		Preload("Payments").
		Preload("Shippings")

	queryCount := config.DB.Model(&models.SalesTransaction{})

	// add params for not using pagination
	if c.Query("all") == "true" {
		pagination.Limit = -1 // No limit
		pagination.Offset = 0 // No offset
	}

	// Filter by date range
	if startDate := c.Query("start_date"); startDate != "" {
		query = query.Where("transaction_date >= ?", startDate)
		queryCount = queryCount.Where("transaction_date >= ?", startDate)
	}

	if endDate := c.Query("end_date"); endDate != "" {
		query = query.Where("transaction_date <= ?", endDate+" 23:59:59")
		queryCount = queryCount.Where("transaction_date <= ?", endDate+" 23:59:59")
	}

	// Filter by payment type
	if paymentType := c.Query("payment_type"); paymentType != "" && paymentType != "all" {
		query = query.Where("payment_type = ?", paymentType)
		queryCount = queryCount.Where("payment_type = ?", paymentType)
	}

	// Filter by status
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
		queryCount = queryCount.Where("status = ?", status)
	}

	// Filter by sales associate
	if salesAssociateID := c.Query("sales_associate_id"); salesAssociateID != "" {
		query = query.Where("sales_associate_id = ?", salesAssociateID)
		queryCount = queryCount.Where("sales_associate_id = ?", salesAssociateID)
	}

	// Apply pagination and fetch data
	if err := query.Offset(pagination.Offset).Limit(pagination.Limit).Find(&transactions).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch sales report",
		})
	}

	// Build report data with total_items for each transaction
	var reportData []SalesReportData
	var summary SalesReportSummary
	summary.TotalTransactions = len(transactions)

	for _, tx := range transactions {
		summary.TotalAmount += tx.TotalAmount

		// Calculate total items for this transaction
		totalItems := 0
		for _, item := range tx.Items {
			totalItems += item.Quantity
			summary.TotalItems += item.Quantity
		}

		reportData = append(reportData, SalesReportData{
			SalesTransaction: tx,
			TotalItems:       totalItems,
		})

		if tx.PaymentType == "T" {
			summary.CashTransactions++
			summary.CashAmount += tx.TotalAmount
		} else {
			summary.CreditTransactions++
			summary.CreditAmount += tx.TotalAmount
		}
	}

	// Create pagination response
	response, err := helpers.CreatePaginationResponse(queryCount, reportData, "data", pagination.Page, pagination.Limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create pagination response",
		})
	}

	// Add summary to response
	response["summary"] = summary

	return c.JSON(response)
}

// GetBooksStockReport godoc
// @Summary Get books stock report
// @Description Get a report of all books with their current stock levels
// @Tags Reports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items per page (default: 20)"
// @Param all query bool false "Get all records without pagination"
// @Param jenis_buku_id query string false "Filter by jenis buku ID"
// @Param jenjang_studi_id query string false "Filter by jenjang studi ID"
// @Param curriculum_id query string false "Filter by curriculum ID"
// @Param kelas query string false "Filter by kelas code"
// @Param low_stock_threshold query int false "Low stock threshold (default: 10)"
// @Param sort_by query string false "Sort by field (stock, name, created_at)"
// @Param sort_order query string false "Sort order (asc, desc)"
// @Success 200 {object} map[string]interface{} "Books stock report with summary and pagination"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/reports/books-stock [get]
func GetBooksStockReport(c *fiber.Ctx) error {
	var books []models.Book

	// Get pagination parameters
	pagination := helpers.GetPaginationParams(c)

	// Default sort by stock ascending (low stock first)
	sortBy := c.Query("sort_by", "stock")
	sortOrder := c.Query("sort_order", "asc")
	orderClause := sortBy + " " + sortOrder

	query := config.DB.Order(orderClause).
		Preload("MerkBuku").
		Preload("JenisBuku").
		Preload("JenjangStudi").
		Preload("BidangStudi").
		Preload("Curriculum").
		Preload("Publisher")

	queryCount := config.DB.Model(&models.Book{})

	// add params for not using pagination
	if c.Query("all") == "true" {
		pagination.Limit = -1 // No limit
		pagination.Offset = 0 // No offset
	}

	// Filter by jenis buku
	if jenisBukuID := c.Query("jenis_buku_id"); jenisBukuID != "" {
		query = query.Where("jenis_buku_id = ?", jenisBukuID)
		queryCount = queryCount.Where("jenis_buku_id = ?", jenisBukuID)
	}

	// Filter by jenjang studi
	if jenjangStudiID := c.Query("jenjang_studi_id"); jenjangStudiID != "" {
		query = query.Where("jenjang_studi_id = ?", jenjangStudiID)
		queryCount = queryCount.Where("jenjang_studi_id = ?", jenjangStudiID)
	}

	// Filter by curriculum
	if curriculumID := c.Query("curriculum_id"); curriculumID != "" {
		query = query.Where("curriculum_id = ?", curriculumID)
		queryCount = queryCount.Where("curriculum_id = ?", curriculumID)
	}

	// Filter by kelas
	if kelas := c.Query("kelas"); kelas != "" {
		query = query.Where("kelas = ?", kelas)
		queryCount = queryCount.Where("kelas = ?", kelas)
	}

	// Apply pagination and fetch data
	if err := query.Offset(pagination.Offset).Limit(pagination.Limit).Find(&books).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch books stock report",
		})
	}

	// Calculate summary
	var summary BooksStockSummary
	summary.TotalBooks = len(books)

	// Default low stock threshold
	lowStockThreshold := 10
	if threshold := c.QueryInt("low_stock_threshold", 10); threshold > 0 {
		lowStockThreshold = threshold
	}

	for _, book := range books {
		summary.TotalStock += book.Stock
		if book.Stock <= lowStockThreshold {
			summary.LowStockCount++
		}
	}

	// Create pagination response
	response, err := helpers.CreatePaginationResponse(queryCount, books, "data", pagination.Page, pagination.Limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create pagination response",
		})
	}

	// Add summary to response
	response["summary"] = summary
	response["low_stock_threshold"] = lowStockThreshold

	return c.JSON(response)
}

// GetCreditsReport godoc
// @Summary Get credits (piutang) report
// @Description Get a report of all outstanding credit transactions (remaining balances)
// @Tags Reports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items per page (default: 20)"
// @Param all query bool false "Get all records without pagination"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param sales_associate_id query string false "Filter by sales associate ID"
// @Param overdue_only query bool false "Show only overdue transactions"
// @Success 200 {object} map[string]interface{} "Credits report with summary and pagination"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/reports/credits [get]
func GetCreditsReport(c *fiber.Ctx) error {
	var transactions []models.SalesTransaction

	// Get pagination parameters
	pagination := helpers.GetPaginationParams(c)

	// Only get credit transactions that are not fully paid
	query := config.DB.Order("transaction_date ASC").
		Where("payment_type = ?", "K").
		Where("status != ?", 1). // Not paid-off
		Preload("Biller").
		Preload("SalesAssociate").
		Preload("SalesAssociate.City").
		Preload("Items").
		Preload("Items.Book").
		Preload("Payments")

	queryCount := config.DB.Model(&models.SalesTransaction{}).
		Where("payment_type = ?", "K").
		Where("status != ?", 1) // Not paid-off

	// add params for not using pagination
	if c.Query("all") == "true" {
		pagination.Limit = -1 // No limit
		pagination.Offset = 0 // No offset
	}

	// Filter by date range (transaction date)
	if startDate := c.Query("start_date"); startDate != "" {
		query = query.Where("transaction_date >= ?", startDate)
		queryCount = queryCount.Where("transaction_date >= ?", startDate)
	}

	if endDate := c.Query("end_date"); endDate != "" {
		query = query.Where("transaction_date <= ?", endDate+" 23:59:59")
		queryCount = queryCount.Where("transaction_date <= ?", endDate+" 23:59:59")
	}

	// Filter by sales associate
	if salesAssociateID := c.Query("sales_associate_id"); salesAssociateID != "" {
		query = query.Where("sales_associate_id = ?", salesAssociateID)
		queryCount = queryCount.Where("sales_associate_id = ?", salesAssociateID)
	}

	// Apply pagination and fetch data
	if err := query.Offset(pagination.Offset).Limit(pagination.Limit).Find(&transactions).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch credits report",
		})
	}

	// Calculate summary and build report items
	var summary CreditReportSummary
	var reportItems []CreditReportItem

	for _, tx := range transactions {
		// Calculate total paid
		var totalPaid float64
		for _, payment := range tx.Payments {
			totalPaid += payment.Amount
		}

		remainingAmount := tx.TotalAmount - totalPaid
		if remainingAmount <= 0 {
			continue // Skip fully paid
		}

		// Calculate total items for this transaction
		totalItems := 0
		for _, item := range tx.Items {
			totalItems += item.Quantity
			summary.TotalItems += item.Quantity
		}

		item := CreditReportItem{
			Transaction:     tx,
			TotalPaid:       totalPaid,
			RemainingAmount: remainingAmount,
			TotalItems:      totalItems,
		}
		reportItems = append(reportItems, item)

		// Update summary
		summary.TotalTransactions++
		summary.TotalOutstanding += remainingAmount
	}

	// Create pagination response with the filtered report items
	// Note: We need to count the filtered items, not use the queryCount
	// because we filter out fully paid transactions in the loop
	paginationMeta := helpers.PaginationMeta{
		Page:       pagination.Page,
		Limit:      pagination.Limit,
		Total:      int64(len(reportItems)),
		TotalPages: 0,
	}

	// Calculate total pages if pagination is enabled
	if pagination.Limit > 0 {
		paginationMeta.TotalPages = (len(reportItems) + pagination.Limit - 1) / pagination.Limit
	}

	response := fiber.Map{
		"data":       reportItems,
		"pagination": paginationMeta,
		"summary":    summary,
	}

	return c.JSON(response)
}
