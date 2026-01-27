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
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

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

		transactionRows := sqlmock.NewRows([]string{
			"id", "biller_id", "sales_associate_id", "no_invoice", "payment_type",
			"transaction_date", "due_date", "total_amount", "status",
			"created_at", "updated_at",
		})

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
			TransactionDate: time.Now(),
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
			TransactionDate:  time.Now(),
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
			TransactionDate:  time.Now(),
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

func TestAddInstallment(t *testing.T) {
	db, mock, err := testutil.SetupMockDB()
	assert.NoError(t, err)
	defer testutil.CloseMockDB(db)

	app := fiber.New()
	app.Post("/sales-transactions/:transaction_id/installments", handlers.AddInstallment)

	t.Run("Transaction not found", func(t *testing.T) {
		transactionID := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "sales_transactions" WHERE id = $1`)).
			WithArgs(transactionID.String()).
			WillReturnError(gorm.ErrRecordNotFound)

		requestBody := handlers.CreateInstallmentRequest{
			InstallmentDate: time.Now(),
			Amount:          250000.0,
		}

		bodyBytes, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", fmt.Sprintf("/sales-transactions/%s/installments", transactionID.String()), bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Transaction not found", response["error"])
	})

	t.Run("Cannot add installment to cash transaction", func(t *testing.T) {
		transactionID := uuid.New()
		salesAssociateID := uuid.New()
		billerID := uuid.New()

		// Transaction with payment_type 'T' (cash)
		transactionRows := sqlmock.NewRows([]string{
			"id", "biller_id", "sales_associate_id", "no_invoice", "payment_type",
			"transaction_date", "due_date", "total_amount", "status",
			"created_at", "updated_at",
		}).
			AddRow(transactionID, billerID, salesAssociateID, "INV2024010100000001", "T",
				time.Now(), nil, 500000.00, 1, time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "sales_transactions" WHERE id = $1`)).
			WithArgs(transactionID.String()).
			WillReturnRows(transactionRows)

		requestBody := handlers.CreateInstallmentRequest{
			InstallmentDate: time.Now(),
			Amount:          250000.0,
		}

		bodyBytes, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", fmt.Sprintf("/sales-transactions/%s/installments", transactionID.String()), bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Installments can only be added to credit transactions", response["error"])
	})

	t.Run("Invalid request body", func(t *testing.T) {
		transactionID := uuid.New()
		salesAssociateID := uuid.New()
		billerID := uuid.New()

		transactionRows := sqlmock.NewRows([]string{
			"id", "biller_id", "sales_associate_id", "no_invoice", "payment_type",
			"transaction_date", "due_date", "total_amount", "status",
			"created_at", "updated_at",
		}).
			AddRow(transactionID, billerID, salesAssociateID, "INV2024010100000001", "K",
				time.Now(), time.Now().AddDate(0, 1, 0), 500000.00, 2, time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "sales_transactions" WHERE id = $1`)).
			WithArgs(transactionID.String()).
			WillReturnRows(transactionRows)

		req := httptest.NewRequest("POST", fmt.Sprintf("/sales-transactions/%s/installments", transactionID.String()), bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Invalid request body", response["error"])
	})
}

func TestGetTransactionInstallments(t *testing.T) {
	db, mock, err := testutil.SetupMockDB()
	assert.NoError(t, err)
	defer testutil.CloseMockDB(db)

	app := fiber.New()
	app.Get("/sales-transactions/:transaction_id/installments", handlers.GetTransactionInstallments)

	t.Run("Transaction not found", func(t *testing.T) {
		transactionID := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "sales_transactions" WHERE id = $1`)).
			WithArgs(transactionID.String()).
			WillReturnError(gorm.ErrRecordNotFound)

		req := httptest.NewRequest("GET", fmt.Sprintf("/sales-transactions/%s/installments", transactionID.String()), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Transaction not found", response["error"])
	})

	t.Run("Successfully get installments", func(t *testing.T) {
		transactionID := uuid.New()
		salesAssociateID := uuid.New()
		billerID := uuid.New()

		// Verify transaction exists
		transactionRows := sqlmock.NewRows([]string{
			"id", "biller_id", "sales_associate_id", "no_invoice", "payment_type",
			"transaction_date", "due_date", "total_amount", "status",
			"created_at", "updated_at",
		}).
			AddRow(transactionID, billerID, salesAssociateID, "INV2024010100000001", "K",
				time.Now(), time.Now().AddDate(0, 1, 0), 500000.00, 2, time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "sales_transactions" WHERE id = $1`)).
			WithArgs(transactionID.String()).
			WillReturnRows(transactionRows)

		// Fetch installments
		installmentID := uuid.New()
		installmentRows := sqlmock.NewRows([]string{"id", "transaction_id", "installment_date", "amount", "note", "created_at", "updated_at"}).
			AddRow(installmentID, transactionID, time.Now(), 250000.0, nil, time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "sales_transaction_installments" WHERE transaction_id = $1 ORDER BY installment_date ASC`)).
			WithArgs(transactionID.String()).
			WillReturnRows(installmentRows)

		req := httptest.NewRequest("GET", fmt.Sprintf("/sales-transactions/%s/installments", transactionID.String()), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["transaction_id"])
		assert.NotNil(t, response["total_amount"])
		assert.NotNil(t, response["total_paid"])
		assert.NotNil(t, response["remaining"])
		assert.NotNil(t, response["installments"])
	})
}

func TestDeleteSalesTransactionInstallment(t *testing.T) {
	db, mock, err := testutil.SetupMockDB()
	assert.NoError(t, err)
	defer testutil.CloseMockDB(db)

	app := fiber.New()
	app.Delete("/sales-transactions/:transaction_id/installments/:id", handlers.DeleteInstallment)

	t.Run("Successfully delete installment", func(t *testing.T) {
		transactionID := uuid.New()
		installmentID := uuid.New()

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "sales_transaction_installments" WHERE id = $1 AND sales_transaction_id = $2`)).
			WithArgs(installmentID.String(), transactionID.String()).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		req := httptest.NewRequest("DELETE", fmt.Sprintf("/sales-transactions/%s/installments/%s", transactionID.String(), installmentID.String()), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Installment deleted successfully", response["message"])
	})

	t.Run("Installment not found", func(t *testing.T) {
		transactionID := uuid.New()
		installmentID := uuid.New()

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "sales_transaction_installments" WHERE id = $1 AND sales_transaction_id = $2`)).
			WithArgs(installmentID.String(), transactionID.String()).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		req := httptest.NewRequest("DELETE", fmt.Sprintf("/sales-transactions/%s/installments/%s", transactionID.String(), installmentID.String()), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Installment not found", response["error"])
	})

	t.Run("Database error", func(t *testing.T) {
		transactionID := uuid.New()
		installmentID := uuid.New()

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "sales_transaction_installments" WHERE id = $1 AND sales_transaction_id = $2`)).
			WithArgs(installmentID.String(), transactionID.String()).
			WillReturnError(gorm.ErrInvalidDB)
		mock.ExpectRollback()

		req := httptest.NewRequest("DELETE", fmt.Sprintf("/sales-transactions/%s/installments/%s", transactionID.String(), installmentID.String()), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Failed to delete installment", response["error"])
	})
}
