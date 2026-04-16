package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"pustaka-backend/handlers"
	"pustaka-backend/tests/testutil"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var salesTransactionColumns = []string{
	"id", "biller_id", "sales_associate_id", "no_invoice", "payment_type",
	"transaction_date", "total_amount", "status",
	"periode", "year", "curriculum_id", "merk_buku_id", "jenjang_studi_id",
	"created_at", "updated_at",
}

func TestGetAllSalesTransactions(t *testing.T) {
	db, _, err := testutil.SetupMockDB()
	assert.NoError(t, err)
	defer testutil.CloseMockDB(db)

	app := fiber.New()
	app.Get("/sales-transactions", handlers.GetAllSalesTransactions)

	t.Run("Database error", func(t *testing.T) {
		db3, mock3, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db3)

		mock3.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "sales_transactions" ORDER BY sales_transactions.created_at desc LIMIT 20`)).
			WillReturnError(gorm.ErrInvalidDB)

		req := httptest.NewRequest("GET", "/sales-transactions", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Failed to fetch sales transactions", response["error"])
	})

	// Use mock with MatchExpectationsInOrder(false) for success cases
	t.Run("Empty list", func(t *testing.T) {
		db2, mock2, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db2)

		transactionRows := sqlmock.NewRows(salesTransactionColumns)

		mock2.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "sales_transactions" ORDER BY sales_transactions.created_at desc LIMIT 20`)).
			WillReturnRows(transactionRows)

		countRows := sqlmock.NewRows([]string{"count"}).AddRow(0)
		mock2.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "sales_transactions"`)).
			WillReturnRows(countRows)

		req := httptest.NewRequest("GET", "/sales-transactions", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["sales_transactions"])
		assert.NotNil(t, response["pagination"])
	})
}

func TestGetSalesTransaction(t *testing.T) {
	db, mock, err := testutil.SetupMockDB()
	assert.NoError(t, err)
	defer testutil.CloseMockDB(db)

	app := fiber.New()
	app.Get("/sales-transactions/:id", handlers.GetSalesTransaction)

	t.Run("Transaction not found", func(t *testing.T) {
		transactionID := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "sales_transactions" WHERE id = $1`)).
			WithArgs(transactionID.String()).
			WillReturnError(gorm.ErrRecordNotFound)

		req := httptest.NewRequest("GET", fmt.Sprintf("/sales-transactions/%s", transactionID.String()), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Transaction not found", response["error"])
	})
}

func TestCreateSalesTransaction(t *testing.T) {
	db, mock, err := testutil.SetupMockDB()
	assert.NoError(t, err)
	defer testutil.CloseMockDB(db)

	app := fiber.New()
	app.Post("/sales-transactions", handlers.CreateSalesTransaction)

	t.Run("Invalid request body", func(t *testing.T) {
		defaultBillerID := uuid.New()

		defaultBillerRow := sqlmock.NewRows([]string{"id"}).
			AddRow(defaultBillerID)

		mock.ExpectQuery(`SELECT "id" FROM "billers"`).
			WillReturnRows(defaultBillerRow)

		req := httptest.NewRequest("POST", "/sales-transactions", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Invalid request body", response["error"])
	})

	t.Run("Missing sales_associate_id", func(t *testing.T) {
		defaultBillerID := uuid.New()

		defaultBillerRow := sqlmock.NewRows([]string{"id"}).
			AddRow(defaultBillerID)

		mock.ExpectQuery(`SELECT "id" FROM "billers"`).
			WillReturnRows(defaultBillerRow)

		requestBody := handlers.CreateTransactionRequest{
			PaymentType:     "T",
			TransactionDate: testutil.StringPtr("2024-01-15"),
			Items: []handlers.CreateTransactionItemRequest{
				{
					BookID:   uuid.New().String(),
					Quantity: 10,
				},
			},
		}

		bodyBytes, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", "/sales-transactions", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "sales_associate_id is required", response["error"])
	})

	t.Run("No items provided", func(t *testing.T) {
		defaultBillerID := uuid.New()

		defaultBillerRow := sqlmock.NewRows([]string{"id"}).
			AddRow(defaultBillerID)

		mock.ExpectQuery(`SELECT "id" FROM "billers"`).
			WillReturnRows(defaultBillerRow)

		requestBody := handlers.CreateTransactionRequest{
			SalesAssociateID: uuid.New().String(),
			PaymentType:      "T",
			TransactionDate:  testutil.StringPtr("2024-01-15"),
			Periode:          1,
			Year:             "2024",
			Items:            []handlers.CreateTransactionItemRequest{},
		}

		bodyBytes, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", "/sales-transactions", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "At least one item is required", response["error"])
	})

	t.Run("Invalid payment type", func(t *testing.T) {
		defaultBillerID := uuid.New()

		defaultBillerRow := sqlmock.NewRows([]string{"id"}).
			AddRow(defaultBillerID)

		mock.ExpectQuery(`SELECT "id" FROM "billers"`).
			WillReturnRows(defaultBillerRow)

		requestBody := handlers.CreateTransactionRequest{
			SalesAssociateID: uuid.New().String(),
			PaymentType:      "X",
			TransactionDate:  testutil.StringPtr("2024-01-15"),
			Periode:          1,
			Year:             "2024",
			Items: []handlers.CreateTransactionItemRequest{
				{
					BookID:   uuid.New().String(),
					Quantity: 10,
				},
			},
		}

		bodyBytes, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", "/sales-transactions", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "payment_type must be either 'T' (cash) or 'K' (credit)", response["error"])
	})
}

func TestDeleteSalesTransaction(t *testing.T) {
	app := fiber.New()
	app.Delete("/sales-transactions/:id", handlers.DeleteSalesTransaction)

	t.Run("Transaction not found", func(t *testing.T) {
		_, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)

		transactionID := uuid.New()

		// First, handler tries to find the transaction with preloaded items
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "sales_transactions" WHERE id = $1`)).
			WithArgs(transactionID.String()).
			WillReturnError(gorm.ErrRecordNotFound)

		req := httptest.NewRequest("DELETE", fmt.Sprintf("/sales-transactions/%s", transactionID.String()), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Transaction not found", response["error"])
	})
}
