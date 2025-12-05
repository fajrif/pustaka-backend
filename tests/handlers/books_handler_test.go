package handlers_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"pustaka-backend/handlers"
	"pustaka-backend/models"
	"pustaka-backend/tests/testutil"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetAllBooks(t *testing.T) {
	db, mock, err := testutil.SetupMockDB()
	assert.NoError(t, err)
	defer testutil.CloseMockDB(db)

	app := fiber.New()
	app.Get("/books", handlers.GetAllBooks)

	t.Run("Successfully get all books", func(t *testing.T) {
		bookID1 := uuid.New()
		bookID2 := uuid.New()

		bookRows := sqlmock.NewRows([]string{"id", "name", "description", "year", "jenis_buku_id", "jenjang_studi_id", "bidang_studi_id", "kelas_id", "publisher_id", "price", "created_at", "updated_at"}).
			AddRow(bookID1, "Mathematics Grade 1", "Test book description", "2024", nil, nil, nil, nil, nil, 50000.00, time.Now(), time.Now()).
			AddRow(bookID2, "Science Grade 2", "Test book description", "2024", nil, nil, nil, nil, nil, 75000.00, time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "books" ORDER BY created_at DESC`)).
			WillReturnRows(bookRows)

		req := httptest.NewRequest("GET", "/books", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["books"])
	})

	t.Run("Empty list", func(t *testing.T) {
		db2, mock2, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db2)

		bookRows := sqlmock.NewRows([]string{"id", "name", "description", "year", "jenis_buku_id", "jenjang_studi_id", "bidang_studi_id", "kelas_id", "publisher_id", "price", "created_at", "updated_at"})

		mock2.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "books" ORDER BY created_at DESC`)).
			WillReturnRows(bookRows)

		req := httptest.NewRequest("GET", "/books", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["books"])
	})

	t.Run("Database error", func(t *testing.T) {
		db3, mock3, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db3)

		mock3.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "books" ORDER BY created_at DESC`)).
			WillReturnError(gorm.ErrInvalidDB)

		req := httptest.NewRequest("GET", "/books", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Failed to fetch all books", response["error"])
	})

	t.Run("Search filter by name", func(t *testing.T) {
		db3, mock3, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db3)

		bookID := uuid.New()
		description := "Test book description"

		bookRows := sqlmock.NewRows([]string{"id", "name", "description", "year", "jenis_buku_id", "jenjang_studi_id", "bidang_studi_id", "kelas_id", "publisher_id", "price", "created_at", "updated_at"}).
			AddRow(bookID, "Mathematics Grade 1", &description, "2024", nil, nil, nil, nil, nil, 50000.00, time.Now(), time.Now())

		mock3.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "books" WHERE books.name ILIKE $1 OR books.description ILIKE $2 ORDER BY created_at DESC`)).
			WithArgs("%Mathematics%", "%Mathematics%").
			WillReturnRows(bookRows)

		req := httptest.NewRequest("GET", "/books?search=Mathematics", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["books"])
	})

	t.Run("Search filter by description", func(t *testing.T) {
		db4, mock4, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db4)

		bookID := uuid.New()
		description := "Science book description"

		bookRows := sqlmock.NewRows([]string{"id", "name", "description", "year", "jenis_buku_id", "jenjang_studi_id", "bidang_studi_id", "kelas_id", "publisher_id", "price", "created_at", "updated_at"}).
			AddRow(bookID, "Science Grade 2", &description, "2024", nil, nil, nil, nil, nil, 75000.00, time.Now(), time.Now())

		mock4.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "books" WHERE books.name ILIKE $1 OR books.description ILIKE $2 ORDER BY created_at DESC`)).
			WithArgs("%Science%", "%Science%").
			WillReturnRows(bookRows)

		req := httptest.NewRequest("GET", "/books?search=Science", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["books"])
	})
}

func TestGetBook(t *testing.T) {
	app := fiber.New()
	app.Get("/books/:id", handlers.GetBook)

	t.Run("Successfully get book by ID", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		bookID := uuid.New()

		bookRows := sqlmock.NewRows([]string{"id", "name", "description", "year", "jenis_buku_id", "jenjang_studi_id", "bidang_studi_id", "kelas_id", "publisher_id", "price", "created_at", "updated_at"}).
			AddRow(bookID, "Mathematics Grade 1", "Test book description", "2024", nil, nil, nil, nil, nil, 50000.00, time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "books" WHERE id = $1`)).
			WithArgs(bookID.String()).
			WillReturnRows(bookRows)

		req := httptest.NewRequest("GET", "/books/"+bookID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["book"])
	})

	t.Run("Book not found", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		bookID := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "books" WHERE id = $1`)).
			WithArgs(bookID.String()).
			WillReturnError(gorm.ErrRecordNotFound)

		req := httptest.NewRequest("GET", "/books/"+bookID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Book not found", response["error"])
	})

	t.Run("Invalid UUID format", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "books" WHERE id = $1`)).
			WithArgs("invalid-uuid").
			WillReturnError(gorm.ErrRecordNotFound)

		req := httptest.NewRequest("GET", "/books/invalid-uuid", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})
}

func TestCreateBook(t *testing.T) {
	app := fiber.New()
	app.Post("/books", handlers.CreateBook)

	t.Run("Successfully create book", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		reqBody := models.Book{
			Name:  "Mathematics Grade 1",
			Price: 50000.00,
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "books"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
				AddRow(uuid.New(), time.Now(), time.Now()))
		mock.ExpectCommit()

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/books", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	})

	t.Run("Invalid request body", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/books", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Invalid request body", response["error"])
	})

	t.Run("Database error on create", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		reqBody := models.Book{
			Name:  "Science Grade 2",
			Price: 75000.00,
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "books"`)).
			WillReturnError(gorm.ErrInvalidDB)
		mock.ExpectRollback()

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/books", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Failed to create book", response["error"])
	})
}

func TestUpdateBook(t *testing.T) {
	app := fiber.New()
	app.Put("/books/:id", handlers.UpdateBook)

	t.Run("Successfully update book", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		bookID := uuid.New()

		bookRows := sqlmock.NewRows([]string{"id", "name", "description", "year", "jenis_buku_id", "jenjang_studi_id", "bidang_studi_id", "kelas_id", "publisher_id", "price", "created_at", "updated_at"}).
			AddRow(bookID, "Mathematics Grade 1", "Test book description", "2024", nil, nil, nil, nil, nil, 50000.00, time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "books" WHERE id = $1`)).
			WithArgs(bookID.String()).
			WillReturnRows(bookRows)

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "books" SET .+ WHERE "id" = .+`).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		reqBody := models.Book{
			ID:    bookID,
			Name:  "Mathematics Grade 1 Updated",
			Price: 55000.00,
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("PUT", "/books/"+bookID.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("Book not found", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		bookID := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "books" WHERE id = $1`)).
			WithArgs(bookID.String()).
			WillReturnError(gorm.ErrRecordNotFound)

		reqBody := models.Book{
			Name:  "Science Grade 2",
			Price: 75000.00,
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("PUT", "/books/"+bookID.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Book not found", response["error"])
	})

	t.Run("Invalid request body", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		bookID := uuid.New()

		bookRows := sqlmock.NewRows([]string{"id", "name", "description", "year", "jenis_buku_id", "jenjang_studi_id", "bidang_studi_id", "kelas_id", "publisher_id", "price", "created_at", "updated_at"}).
			AddRow(bookID, "Mathematics Grade 1", "Test book description", "2024", nil, nil, nil, nil, nil, 50000.00, time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "books" WHERE id = $1`)).
			WithArgs(bookID.String()).
			WillReturnRows(bookRows)

		req := httptest.NewRequest("PUT", "/books/"+bookID.String(), bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Invalid request body", response["error"])
	})

	t.Run("Database error on update", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		bookID := uuid.New()

		bookRows := sqlmock.NewRows([]string{"id", "name", "description", "year", "jenis_buku_id", "jenjang_studi_id", "bidang_studi_id", "kelas_id", "publisher_id", "price", "created_at", "updated_at"}).
			AddRow(bookID, "Mathematics Grade 1", "Test book description", "2024", nil, nil, nil, nil, nil, 50000.00, time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "books" WHERE id = $1`)).
			WithArgs(bookID.String()).
			WillReturnRows(bookRows)

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "books" SET .+ WHERE "id" = .+`).
			WillReturnError(gorm.ErrInvalidDB)
		mock.ExpectRollback()

		reqBody := models.Book{
			ID:    bookID,
			Name:  "Mathematics Grade 1 Updated",
			Price: 55000.00,
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("PUT", "/books/"+bookID.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Failed to update book", response["error"])
	})
}

func TestDeleteBook(t *testing.T) {
	db, mock, err := testutil.SetupMockDB()
	assert.NoError(t, err)
	defer testutil.CloseMockDB(db)

	app := fiber.New()
	app.Delete("/books/:id", handlers.DeleteBook)

	t.Run("Successfully delete book", func(t *testing.T) {
		bookID := uuid.New()

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "books" WHERE id = $1`)).
			WithArgs(bookID.String()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		req := httptest.NewRequest("DELETE", "/books/"+bookID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Book deleted successfully", response["message"])
	})

	t.Run("Book not found", func(t *testing.T) {
		bookID := uuid.New()

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "books" WHERE id = $1`)).
			WithArgs(bookID.String()).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		req := httptest.NewRequest("DELETE", "/books/"+bookID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Book not found", response["error"])
	})

	t.Run("Database error on delete", func(t *testing.T) {
		bookID := uuid.New()

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "books" WHERE id = $1`)).
			WithArgs(bookID.String()).
			WillReturnError(gorm.ErrInvalidDB)
		mock.ExpectRollback()

		req := httptest.NewRequest("DELETE", "/books/"+bookID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Failed to delete book", response["error"])
	})
}
