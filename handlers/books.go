package handlers

import (
	"strings"
	"pustaka-backend/config"
	"pustaka-backend/helpers"
	"pustaka-backend/models"
	"github.com/gofiber/fiber/v2"
)

// GetAllBooks godoc
// @Summary Get all books
// @Description Retrieve all books with their related entities
// @Tags Books
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param search query string false "Search by name or description"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items per page (default: 20)"
// @Param code query string false "Filter by book code"
// @Param bidang_studi_id query string false "Filter by bidang studi ID"
// @Param jenis_buku_id query string false "Filter by jenis buku ID"
// @Param jenjang_studi_id query string false "Filter by jenjang studi ID"
// @Param curriculum_id query string false "Filter by curriculum ID"
// @Param publisher_id query string false "Filter by publisher ID"
// @Param merk_buku_id query string false "Filter by merk buku ID"
// @Param periode query int false "Filter by periode"
// @Param year query string false "Filter by year"
// @Param kelas query string false "Filter by kelas code (e.g., 1, 2, A, B, ALL)"
// @Param price_min query number false "Minimum price filter"
// @Param price_max query number false "Maximum price filter"
// @Success 200 {object} map[string]interface{} "List of all books with pagination"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/books [get]
func GetAllBooks(c *fiber.Ctx) error {
	var books []models.Book

	// Get pagination parameters
	pagination := helpers.GetPaginationParams(c)

	conds := []string{}
	args := []interface{}{}

	query := config.DB.Order("created_at DESC")
	queryCount := config.DB.Model(&models.Book{})

	// add params for not using pagination
	if c.Query("all") == "true" {
		pagination.Limit = -1 // No limit
		pagination.Offset = 0 // No offset
	}

	// Filter search (name or description)
	if searchQuery := c.Query("search"); searchQuery != "" {
		searchTerm := "%" + searchQuery + "%"
		conds = append(conds, "(books.name ILIKE ? OR books.description ILIKE ?)")
		args = append(args, searchTerm, searchTerm)
	}

	// Filter by code (if books have a code field - using name for now as fallback)
	if code := c.Query("code"); code != "" {
		searchTerm := "%" + code + "%"
		conds = append(conds, "books.name ILIKE ?")
		args = append(args, searchTerm)
	}

	// Filter bidang_studi_id
	if bidangStudiId := c.Query("bidang_studi_id"); bidangStudiId != "" {
		conds = append(conds, "books.bidang_studi_id = ?")
		args = append(args, bidangStudiId)
	}

	// Filter jenis_buku_id
	if jenisBukuId := c.Query("jenis_buku_id"); jenisBukuId != "" {
		conds = append(conds, "books.jenis_buku_id = ?")
		args = append(args, jenisBukuId)
	}

	// Filter jenjang_studi_id
	if jenjangStudiId := c.Query("jenjang_studi_id"); jenjangStudiId != "" {
		conds = append(conds, "books.jenjang_studi_id = ?")
		args = append(args, jenjangStudiId)
	}

	// Filter curriculum_id
	if curriculumId := c.Query("curriculum_id"); curriculumId != "" {
		conds = append(conds, "books.curriculum_id = ?")
		args = append(args, curriculumId)
	}

	// Filter publisher_id
	if publisherId := c.Query("publisher_id"); publisherId != "" {
		conds = append(conds, "books.publisher_id = ?")
		args = append(args, publisherId)
	}

	// Filter merk_buku_id
	if merkBukuId := c.Query("merk_buku_id"); merkBukuId != "" {
		conds = append(conds, "books.merk_buku_id = ?")
		args = append(args, merkBukuId)
	}

	// Filter periode
	if periode := c.Query("periode"); periode != "" {
		conds = append(conds, "books.periode = ?")
		args = append(args, periode)
	}

	// Filter year
	if year := c.Query("year"); year != "" {
		conds = append(conds, "books.year = ?")
		args = append(args, year)
	}

	// Filter kelas (CHAR(5) code)
	if kelas := c.Query("kelas"); kelas != "" {
		conds = append(conds, "books.kelas = ?")
		args = append(args, kelas)
	}

	// Filter price range
	if priceMin := c.Query("price_min"); priceMin != "" {
		conds = append(conds, "books.price >= ?")
		args = append(args, priceMin)
	}

	if priceMax := c.Query("price_max"); priceMax != "" {
		conds = append(conds, "books.price <= ?")
		args = append(args, priceMax)
	}

	// Apply all conditions
	if len(conds) > 0 {
		whereClause := strings.Join(conds, " AND ")
		query = query.Where(whereClause, args...)
		queryCount = queryCount.Where(whereClause, args...)
	}

	// Apply pagination and fetch data
	if err := query.
		Offset(pagination.Offset).Limit(pagination.Limit).
		Preload("MerkBuku").
		Preload("JenisBuku").
		Preload("JenjangStudi").
		Preload("BidangStudi").
		Preload("Curriculum").
		Preload("Publisher").
		Preload("Publisher.City").
		Find(&books).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch all books",
		})
	}

	// Create pagination response
	response, err := helpers.CreatePaginationResponse(queryCount, books, "books", pagination.Page, pagination.Limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create pagination response",
		})
	}

	return c.JSON(response)
}

// GetBook godoc
// @Summary Get a book by ID
// @Description Retrieve a single book by its ID with all related entities
// @Tags Books
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Book ID (UUID)"
// @Success 200 {object} map[string]interface{} "Book details"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Book not found"
// @Router /api/books/{id} [get]
func GetBook(c *fiber.Ctx) error {
	id := c.Params("id")

	var book models.Book
	if err := config.DB.
		Preload("MerkBuku").
		Preload("JenisBuku").
		Preload("JenjangStudi").
		Preload("BidangStudi").
		Preload("Curriculum").
		Preload("Publisher").
		Preload("Publisher.City").
		Where("id = ?", id).First(&book).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Book not found",
		})
	}

	return c.JSON(fiber.Map{
		"book": book,
	})
}

// CreateBook godoc
// @Summary Create a new book
// @Description Create a new book entry. If name is empty and bidang_studi_id is provided, name will be auto-populated from bidang_studi.name
// @Tags Books
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.Book true "Book details"
// @Success 201 {object} models.Book "Created book"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/books [post]
func CreateBook(c *fiber.Ctx) error {
	var book models.Book
	if err := c.BodyParser(&book); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Auto-populate name from BidangStudi if name is empty and BidangStudiID is provided
	if book.Name == "" && book.BidangStudiID != nil {
		var bidangStudi models.BidangStudi
		if err := config.DB.Where("id = ?", book.BidangStudiID).First(&bidangStudi).Error; err == nil {
			book.Name = bidangStudi.Name
		}
	}

	// Validate that name is not empty
	if book.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Book name is required (provide name or bidang_studi_id)",
		})
	}

	if err := config.DB.Create(&book).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create book",
		})
	}

	// Fetch created book with relations
	config.DB.
		Preload("MerkBuku").
		Preload("JenisBuku").
		Preload("JenjangStudi").
		Preload("BidangStudi").
		Preload("Curriculum").
		Preload("Publisher").
		Where("id = ?", book.ID).First(&book)

	return c.Status(fiber.StatusCreated).JSON(book)
}

// UpdateBook godoc
// @Summary Update a book
// @Description Update an existing book by ID
// @Tags Books
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Book ID (UUID)"
// @Param request body models.Book true "Updated book details"
// @Success 200 {object} models.Book "Updated book"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Book not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/books/{id} [put]
func UpdateBook(c *fiber.Ctx) error {
	id := c.Params("id")

	var book models.Book
	if err := config.DB.Where("id = ?", id).First(&book).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Book not found",
		})
	}

	if err := c.BodyParser(&book); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := config.DB.Model(&book).Updates(book).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update book",
		})
	}

	return c.JSON(book)
}

// DeleteBook godoc
// @Summary Delete a book
// @Description Delete a book by ID
// @Tags Books
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Book ID (UUID)"
// @Success 200 {object} map[string]interface{} "Book deleted successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Book not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/books/{id} [delete]
func DeleteBook(c *fiber.Ctx) error {
	id := c.Params("id")

	result := config.DB.Delete(&models.Book{}, "id = ?", id)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete book",
		})
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Book not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Book deleted successfully",
	})
}
