package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"pustaka-backend/handlers"
	"pustaka-backend/models"
	"pustaka-backend/tests/testutil"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCreatePurchaseTransactionWithDateOnly(t *testing.T) {
	db, mock, err := testutil.SetupMockDB()
	assert.NoError(t, err)
	defer testutil.CloseMockDB(db)

	app := fiber.New()
	app.Post("/purchase-transactions", handlers.CreatePurchaseTransaction)

	t.Run("Create with date only (no time)", func(t *testing.T) {
		supplierID := uuid.New()
		bookID := uuid.New()

		// Mock supplier check
		supplierRows := sqlmock.NewRows([]string{"id", "code", "name", "address", "phone1"}).
			AddRow(supplierID, "SUP001", "Test Supplier", "Test Address", "08123456789")

		mock.ExpectQuery(`SELECT \* FROM "publishers" WHERE id = \$1`).
			WithArgs(supplierID.String()).
			WillReturnRows(supplierRows)

		// Mock transaction begin
		mock.ExpectBegin()

		// Mock book check
		bookRows := sqlmock.NewRows([]string{"id", "name", "year", "stock", "price"}).
			AddRow(bookID, "Test Book", "2026", 10, 50000.0)

		mock.ExpectQuery(`SELECT \* FROM "books" WHERE id = \$1`).
			WithArgs(bookID.String()).
			WillReturnRows(bookRows)

		// Mock invoice number generation
		mock.ExpectQuery(`SELECT COALESCE\(MAX\(no_invoice\), ''\) FROM "purchase_transactions"`).
			WillReturnRows(sqlmock.NewRows([]string{"max"}).AddRow(""))

		// Mock transaction insert
		mock.ExpectQuery(`INSERT INTO "purchase_transactions"`).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))

		// Mock items insert
		mock.ExpectQuery(`INSERT INTO "purchase_transaction_items"`).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))

		// Mock commit
		mock.ExpectCommit()

		// Mock final select with preloads
		transactionID := uuid.New()
		transactionRows := sqlmock.NewRows([]string{
			"id", "supplier_id", "no_invoice", "purchase_date", "total_amount", "status", "created_at", "updated_at",
		}).AddRow(transactionID, supplierID, "PRC2026013100000001", "2026-01-31", 225000.0, 0, time.Now(), time.Now())

		mock.ExpectQuery(`SELECT \* FROM "purchase_transactions" WHERE id = \$1`).
			WillReturnRows(transactionRows)

		mock.ExpectQuery(`SELECT \* FROM "publishers" WHERE`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "code", "name", "address", "phone1"}).
				AddRow(supplierID, "SUP001", "Test Supplier", "Test Address", "08123456789"))

		mock.ExpectQuery(`SELECT \* FROM "purchase_transaction_items" WHERE`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "purchase_transaction_id", "book_id", "quantity", "price", "subtotal"}).
				AddRow(uuid.New(), transactionID, bookID, 5, 45000.0, 225000.0))

		mock.ExpectQuery(`SELECT \* FROM "books" WHERE`).
			WillReturnRows(bookRows)

		requestBody := handlers.CreatePurchaseTransactionRequest{
			SupplierID:   supplierID.String(),
			PurchaseDate: models.Date{Time: time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC)},
			Items: []handlers.CreatePurchaseTransactionItemReq{
				{
					BookID:   bookID.String(),
					Quantity: 5,
					Price:    45000,
				},
			},
		}

		bodyBytes, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", "/purchase-transactions", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	})
}

func TestUpdatePurchaseTransactionWithDateOnly(t *testing.T) {
	db, mock, err := testutil.SetupMockDB()
	assert.NoError(t, err)
	defer testutil.CloseMockDB(db)

	app := fiber.New()
	app.Put("/purchase-transactions/:id", handlers.UpdatePurchaseTransaction)

	t.Run("Update purchase date", func(t *testing.T) {
		transactionID := uuid.New()
		supplierID := uuid.New()

		// Mock transaction check
		transactionRows := sqlmock.NewRows([]string{
			"id", "supplier_id", "no_invoice", "purchase_date", "total_amount", "status", "created_at", "updated_at",
		}).AddRow(transactionID, supplierID, "PRC2026013100000001", "2026-01-31", 100000.0, 0, time.Now(), time.Now())

		mock.ExpectQuery(`SELECT \* FROM "purchase_transactions" WHERE id = \$1`).
			WithArgs(transactionID.String()).
			WillReturnRows(transactionRows)

		// Mock transaction begin
		mock.ExpectBegin()

		// Mock update
		mock.ExpectExec(`UPDATE "purchase_transactions"`).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Mock commit
		mock.ExpectCommit()

		// Mock final select with preloads
		updatedRows := sqlmock.NewRows([]string{
			"id", "supplier_id", "no_invoice", "purchase_date", "total_amount", "status", "created_at", "updated_at",
		}).AddRow(transactionID, supplierID, "PRC2026013100000001", "2026-02-01", 100000.0, 0, time.Now(), time.Now())

		mock.ExpectQuery(`SELECT \* FROM "purchase_transactions" WHERE id = \$1`).
			WithArgs(transactionID.String()).
			WillReturnRows(updatedRows)

		mock.ExpectQuery(`SELECT \* FROM "publishers" WHERE`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "code", "name", "address", "phone1"}).
				AddRow(supplierID, "SUP001", "Test Supplier", "Test Address", "08123456789"))

		mock.ExpectQuery(`SELECT \* FROM "purchase_transaction_items" WHERE`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "purchase_transaction_id", "book_id", "quantity", "price", "subtotal"}))

		newDate := models.Date{Time: time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)}
		updateBody := handlers.UpdatePurchaseTransactionRequest{
			PurchaseDate: &newDate,
		}

		bodyBytes, _ := json.Marshal(updateBody)

		req := httptest.NewRequest("PUT", "/purchase-transactions/"+transactionID.String(), bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})
}

func TestDateOnlyJSONMarshaling(t *testing.T) {
	t.Run("Marshal date to JSON", func(t *testing.T) {
		date := models.Date{Time: time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC)}

		jsonBytes, err := json.Marshal(date)
		assert.NoError(t, err)
		assert.Equal(t, `"2026-01-31"`, string(jsonBytes))
	})

	t.Run("Unmarshal JSON to date", func(t *testing.T) {
		var date models.Date
		err := json.Unmarshal([]byte(`"2026-01-31"`), &date)
		assert.NoError(t, err)
		assert.Equal(t, "2026-01-31", date.String())
		assert.Equal(t, 2026, date.Year())
		assert.Equal(t, time.Month(1), date.Month())
		assert.Equal(t, 31, date.Day())
	})

	t.Run("Unmarshal null to date", func(t *testing.T) {
		var date models.Date
		err := json.Unmarshal([]byte(`null`), &date)
		assert.NoError(t, err)
		assert.True(t, date.IsZero())
	})

	t.Run("Unmarshal empty string to date", func(t *testing.T) {
		var date models.Date
		err := json.Unmarshal([]byte(`""`), &date)
		assert.NoError(t, err)
		assert.True(t, date.IsZero())
	})
}

func TestDateOnlyFiltering(t *testing.T) {
	db, mock, err := testutil.SetupMockDB()
	assert.NoError(t, err)
	defer testutil.CloseMockDB(db)

	app := fiber.New()
	app.Get("/purchase-transactions", handlers.GetAllPurchaseTransactions)

	t.Run("Filter by start_date", func(t *testing.T) {
		// Mock the query with date filter
		transactionRows := sqlmock.NewRows([]string{
			"id", "supplier_id", "no_invoice", "purchase_date", "total_amount", "status", "created_at", "updated_at",
		})

		mock.ExpectQuery(`SELECT \* FROM "purchase_transactions"`).
			WillReturnRows(transactionRows)

		countRows := sqlmock.NewRows([]string{"count"}).AddRow(0)
		mock.ExpectQuery(`SELECT count\(\*\) FROM "purchase_transactions"`).
			WillReturnRows(countRows)

		req := httptest.NewRequest("GET", "/purchase-transactions?start_date=2026-01-26", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("Filter by date range", func(t *testing.T) {
		transactionRows := sqlmock.NewRows([]string{
			"id", "supplier_id", "no_invoice", "purchase_date", "total_amount", "status", "created_at", "updated_at",
		})

		mock.ExpectQuery(`SELECT \* FROM "purchase_transactions"`).
			WillReturnRows(transactionRows)

		countRows := sqlmock.NewRows([]string{"count"}).AddRow(0)
		mock.ExpectQuery(`SELECT count\(\*\) FROM "purchase_transactions"`).
			WillReturnRows(countRows)

		req := httptest.NewRequest("GET", "/purchase-transactions?start_date=2026-01-25&end_date=2026-01-27", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})
}

func TestPurchaseTransactionNotFound(t *testing.T) {
	db, mock, err := testutil.SetupMockDB()
	assert.NoError(t, err)
	defer testutil.CloseMockDB(db)

	app := fiber.New()
	app.Get("/purchase-transactions/:id", handlers.GetPurchaseTransaction)

	t.Run("Transaction not found", func(t *testing.T) {
		transactionID := uuid.New()

		mock.ExpectQuery(`SELECT \* FROM "purchase_transactions" WHERE id = \$1`).
			WithArgs(transactionID.String()).
			WillReturnError(gorm.ErrRecordNotFound)

		req := httptest.NewRequest("GET", "/purchase-transactions/"+transactionID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})
}
