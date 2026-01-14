package handlers

import (
	"pustaka-backend/config"
	"pustaka-backend/models"
	"time"

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
	OverdueCount      int     `json:"overdue_count"`
	OverdueAmount     float64 `json:"overdue_amount"`
}

// CreditReportItem represents a single item in the credits report
type CreditReportItem struct {
	Transaction     models.SalesTransaction `json:"transaction"`
	TotalPaid       float64                 `json:"total_paid"`
	RemainingAmount float64                 `json:"remaining_amount"`
	IsOverdue       bool                    `json:"is_overdue"`
}

// GetPurchasingReport godoc
// @Summary Get purchasing report
// @Description Get a report of all purchases from suppliers with filters
// @Tags Reports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param supplier_id query string false "Filter by supplier ID"
// @Param status query int false "Filter by status (0=pending, 1=completed, 2=cancelled)"
// @Success 200 {object} map[string]interface{} "Purchasing report with summary"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/reports/purchases [get]
func GetPurchasingReport(c *fiber.Ctx) error {
	var transactions []models.PurchaseTransaction

	query := config.DB.Order("purchase_date DESC").
		Preload("Supplier").
		Preload("Items").
		Preload("Items.Book").
		Preload("Items.Book.MerkBuku")

	// Filter by date range
	if startDate := c.Query("start_date"); startDate != "" {
		query = query.Where("purchase_date >= ?", startDate)
	}

	if endDate := c.Query("end_date"); endDate != "" {
		query = query.Where("purchase_date <= ?", endDate+" 23:59:59")
	}

	// Filter by supplier
	if supplierID := c.Query("supplier_id"); supplierID != "" {
		query = query.Where("supplier_id = ?", supplierID)
	}

	// Filter by status
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Find(&transactions).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch purchasing report",
		})
	}

	// Calculate summary
	var summary PurchasingReportSummary
	summary.TotalTransactions = len(transactions)

	for _, tx := range transactions {
		summary.TotalAmount += tx.TotalAmount
		for range tx.Items {
			summary.TotalItems++
		}
	}

	return c.JSON(fiber.Map{
		"data":    transactions,
		"summary": summary,
	})
}

// GetSalesReport godoc
// @Summary Get sales report
// @Description Get a report of all sales transactions with filters
// @Tags Reports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param payment_type query string false "Filter by payment type (T=cash, K=credit, all=both)"
// @Param status query int false "Filter by status (0=booking, 1=paid-off, 2=installment)"
// @Param sales_associate_id query string false "Filter by sales associate ID"
// @Success 200 {object} map[string]interface{} "Sales report with summary"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/reports/sales [get]
func GetSalesReport(c *fiber.Ctx) error {
	var transactions []models.SalesTransaction

	query := config.DB.Order("transaction_date DESC").
		Preload("Biller").
		Preload("SalesAssociate").
		Preload("Items").
		Preload("Items.Book").
		Preload("Items.Book.MerkBuku").
		Preload("Payments").
		Preload("Shippings")

	// Filter by date range
	if startDate := c.Query("start_date"); startDate != "" {
		query = query.Where("transaction_date >= ?", startDate)
	}

	if endDate := c.Query("end_date"); endDate != "" {
		query = query.Where("transaction_date <= ?", endDate+" 23:59:59")
	}

	// Filter by payment type
	if paymentType := c.Query("payment_type"); paymentType != "" && paymentType != "all" {
		query = query.Where("payment_type = ?", paymentType)
	}

	// Filter by status
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// Filter by sales associate
	if salesAssociateID := c.Query("sales_associate_id"); salesAssociateID != "" {
		query = query.Where("sales_associate_id = ?", salesAssociateID)
	}

	if err := query.Find(&transactions).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch sales report",
		})
	}

	// Calculate summary
	var summary SalesReportSummary
	summary.TotalTransactions = len(transactions)

	for _, tx := range transactions {
		summary.TotalAmount += tx.TotalAmount

		if tx.PaymentType == "T" {
			summary.CashTransactions++
			summary.CashAmount += tx.TotalAmount
		} else {
			summary.CreditTransactions++
			summary.CreditAmount += tx.TotalAmount
		}
	}

	return c.JSON(fiber.Map{
		"data":    transactions,
		"summary": summary,
	})
}

// GetBooksStockReport godoc
// @Summary Get books stock report
// @Description Get a report of all books with their current stock levels
// @Tags Reports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param jenis_buku_id query string false "Filter by jenis buku ID"
// @Param jenjang_studi_id query string false "Filter by jenjang studi ID"
// @Param curriculum_id query string false "Filter by curriculum ID"
// @Param kelas query string false "Filter by kelas code"
// @Param low_stock_threshold query int false "Low stock threshold (default: 10)"
// @Param sort_by query string false "Sort by field (stock, name, created_at)"
// @Param sort_order query string false "Sort order (asc, desc)"
// @Success 200 {object} map[string]interface{} "Books stock report with summary"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/reports/books-stock [get]
func GetBooksStockReport(c *fiber.Ctx) error {
	var books []models.Book

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

	// Filter by jenis buku
	if jenisBukuID := c.Query("jenis_buku_id"); jenisBukuID != "" {
		query = query.Where("jenis_buku_id = ?", jenisBukuID)
	}

	// Filter by jenjang studi
	if jenjangStudiID := c.Query("jenjang_studi_id"); jenjangStudiID != "" {
		query = query.Where("jenjang_studi_id = ?", jenjangStudiID)
	}

	// Filter by curriculum
	if curriculumID := c.Query("curriculum_id"); curriculumID != "" {
		query = query.Where("curriculum_id = ?", curriculumID)
	}

	// Filter by kelas
	if kelas := c.Query("kelas"); kelas != "" {
		query = query.Where("kelas = ?", kelas)
	}

	if err := query.Find(&books).Error; err != nil {
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

	return c.JSON(fiber.Map{
		"data":                 books,
		"summary":              summary,
		"low_stock_threshold":  lowStockThreshold,
	})
}

// GetCreditsReport godoc
// @Summary Get credits (piutang) report
// @Description Get a report of all outstanding credit transactions (remaining balances)
// @Tags Reports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param sales_associate_id query string false "Filter by sales associate ID"
// @Param overdue_only query bool false "Show only overdue transactions"
// @Success 200 {object} map[string]interface{} "Credits report with summary"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/reports/credits [get]
func GetCreditsReport(c *fiber.Ctx) error {
	var transactions []models.SalesTransaction

	// Only get credit transactions that are not fully paid
	query := config.DB.Order("due_date ASC").
		Where("payment_type = ?", "K").
		Where("status != ?", 1). // Not paid-off
		Preload("Biller").
		Preload("SalesAssociate").
		Preload("SalesAssociate.City").
		Preload("Items").
		Preload("Items.Book").
		Preload("Payments")

	// Filter by date range (transaction date)
	if startDate := c.Query("start_date"); startDate != "" {
		query = query.Where("transaction_date >= ?", startDate)
	}

	if endDate := c.Query("end_date"); endDate != "" {
		query = query.Where("transaction_date <= ?", endDate+" 23:59:59")
	}

	// Filter by sales associate
	if salesAssociateID := c.Query("sales_associate_id"); salesAssociateID != "" {
		query = query.Where("sales_associate_id = ?", salesAssociateID)
	}

	if err := query.Find(&transactions).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch credits report",
		})
	}

	// Calculate summary and build report items
	var summary CreditReportSummary
	var reportItems []CreditReportItem
	now := time.Now()
	overdueOnly := c.Query("overdue_only") == "true"

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

		// Check if overdue
		isOverdue := false
		if tx.DueDate != nil && tx.DueDate.Before(now) {
			isOverdue = true
		}

		// Filter by overdue only if requested
		if overdueOnly && !isOverdue {
			continue
		}

		item := CreditReportItem{
			Transaction:     tx,
			TotalPaid:       totalPaid,
			RemainingAmount: remainingAmount,
			IsOverdue:       isOverdue,
		}
		reportItems = append(reportItems, item)

		// Update summary
		summary.TotalTransactions++
		summary.TotalOutstanding += remainingAmount

		if isOverdue {
			summary.OverdueCount++
			summary.OverdueAmount += remainingAmount
		}
	}

	return c.JSON(fiber.Map{
		"data":    reportItems,
		"summary": summary,
	})
}
