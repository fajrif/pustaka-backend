package handlers

import (
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
// @Success 200 {object} map[string]interface{} "List of all books with pagination"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/books [get]
func GetAllBooks(c *fiber.Ctx) error {
	var books []models.Book

	// Get pagination parameters
	pagination := helpers.GetPaginationParams(c)

	query := config.DB.Order("created_at DESC")
	queryCount := config.DB.Model(&models.Book{})

	// Filter search
	if searchQuery := c.Query("search"); searchQuery != "" {
		// Wrap string search with wildcard SQL LIKE
		searchTerm := "%" + searchQuery + "%"
		cond := "books.name ILIKE ? OR books.description ILIKE ?"
		args := []interface{}{searchTerm, searchTerm}

		query = query.Where(cond, args...)
		queryCount = queryCount.Where(cond, args...)
	}

	// Apply pagination and fetch data
	if err := query.
		Offset(pagination.Offset).Limit(pagination.Limit).
		Preload("JenisBuku").
		Preload("JenjangStudi").
		Preload("BidangStudi").
		Preload("Kelas").
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
		Preload("JenisBuku").
		Preload("JenjangStudi").
		Preload("BidangStudi").
		Preload("Kelas").
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
// @Description Create a new book entry
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

	if err := config.DB.Create(&book).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create book",
		})
	}

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
