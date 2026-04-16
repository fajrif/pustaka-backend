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

var paymentColumns = []string{
	"id", "sales_transaction_id", "no_payment", "payment_date", "amount",
	"discount_percentage", "discount_amount", "note", "created_at", "updated_at",
}

func TestGetTransactionPayments(t *testing.T) {
	db, mock, err := testutil.SetupMockDB()
	assert.NoError(t, err)
	defer testutil.CloseMockDB(db)

	app := fiber.New()
	app.Get("/sales-transactions/:transaction_id/payments", handlers.GetTransactionPayments)

	t.Run("Transaction not found", func(t *testing.T) {
		transactionID := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "sales_transactions" WHERE id = $1`)).
			WithArgs(transactionID.String()).
			WillReturnError(gorm.ErrRecordNotFound)

		req := httptest.NewRequest("GET", fmt.Sprintf("/sales-transactions/%s/payments", transactionID.String()), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Transaction not found", response["error"])
	})

	t.Run("Successfully get payments", func(t *testing.T) {
		transactionID := uuid.New()
		salesAssociateID := uuid.New()
		billerID := uuid.New()

		transactionRows := sqlmock.NewRows([]string{
			"id", "biller_id", "sales_associate_id", "no_invoice", "payment_type",
			"transaction_date", "total_amount", "status",
			"periode", "year", "curriculum_id", "merk_buku_id", "jenjang_studi_id",
			"created_at", "updated_at",
		}).AddRow(
			transactionID, billerID, salesAssociateID, "INV2024010100000001", "K",
			time.Now(), 500000.00, 2, 1, "2024", nil, nil, nil, time.Now(), time.Now(),
		)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "sales_transactions" WHERE id = $1`)).
			WithArgs(transactionID.String()).
			WillReturnRows(transactionRows)

		paymentID := uuid.New()
		paymentRows := sqlmock.NewRows(paymentColumns).
			AddRow(paymentID, transactionID, "PMT2024010100000001", time.Now(), 250000.0, 8.0, 20000.0, nil, time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "payments" WHERE sales_transaction_id = $1 ORDER BY payment_date ASC`)).
			WithArgs(transactionID.String()).
			WillReturnRows(paymentRows)

		req := httptest.NewRequest("GET", fmt.Sprintf("/sales-transactions/%s/payments", transactionID.String()), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["transaction_id"])
		assert.NotNil(t, response["total_amount"])
		assert.NotNil(t, response["total_paid"])
		assert.NotNil(t, response["total_discount"])
		assert.NotNil(t, response["remaining_amount"])
		assert.NotNil(t, response["payments"])
	})

	t.Run("Empty payments list", func(t *testing.T) {
		transactionID := uuid.New()
		salesAssociateID := uuid.New()
		billerID := uuid.New()

		transactionRows := sqlmock.NewRows([]string{
			"id", "biller_id", "sales_associate_id", "no_invoice", "payment_type",
			"transaction_date", "due_date", "secondary_due_date", "total_amount", "status",
			"periode", "year", "curriculum_id", "merk_buku_id", "jenjang_studi_id",
			"created_at", "updated_at",
		}).AddRow(
			transactionID, billerID, salesAssociateID, "INV2024010100000001", "K",
			time.Now(), time.Now().AddDate(0, 1, 0), nil, 500000.00, 0, 1, "2024", nil, nil, nil, time.Now(), time.Now(),
		)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "sales_transactions" WHERE id = $1`)).
			WithArgs(transactionID.String()).
			WillReturnRows(transactionRows)

		paymentRows := sqlmock.NewRows(paymentColumns)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "payments" WHERE sales_transaction_id = $1 ORDER BY payment_date ASC`)).
			WithArgs(transactionID.String()).
			WillReturnRows(paymentRows)

		req := httptest.NewRequest("GET", fmt.Sprintf("/sales-transactions/%s/payments", transactionID.String()), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, float64(0), response["total_paid"])
		assert.Equal(t, float64(0), response["total_discount"])
	})
}

func TestCreatePayment(t *testing.T) {
	app := fiber.New()
	app.Post("/sales-transactions/:transaction_id/payments", handlers.CreatePayment)

	t.Run("Transaction not found", func(t *testing.T) {
		_, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)

		transactionID := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "sales_transactions" WHERE id = $1`)).
			WithArgs(transactionID.String()).
			WillReturnError(gorm.ErrRecordNotFound)

		requestBody := handlers.CreatePaymentRequest{
			PaymentDate: testutil.StringPtr("2024-01-15"),
			Amount:      500000.0,
		}

		bodyBytes, _ := json.Marshal(requestBody)
		req := httptest.NewRequest("POST", fmt.Sprintf("/sales-transactions/%s/payments", transactionID.String()), bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Transaction not found", response["error"])
	})

	t.Run("Invalid request body", func(t *testing.T) {
		transactionID := uuid.New()
		salesAssociateID := uuid.New()
		billerID := uuid.New()

		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		transactionRows := sqlmock.NewRows([]string{
			"id", "biller_id", "sales_associate_id", "no_invoice", "payment_type",
			"transaction_date", "due_date", "secondary_due_date", "total_amount", "status",
			"periode", "year", "curriculum_id", "merk_buku_id", "jenjang_studi_id",
			"created_at", "updated_at",
		}).AddRow(
			transactionID, billerID, salesAssociateID, "INV2024010100000001", "K",
			time.Now(), time.Now().AddDate(0, 1, 0), nil, 500000.00, 0, 1, "2024", nil, nil, nil, time.Now(), time.Now(),
		)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "sales_transactions" WHERE id = $1`)).
			WithArgs(transactionID.String()).
			WillReturnRows(transactionRows)

		req := httptest.NewRequest("POST", fmt.Sprintf("/sales-transactions/%s/payments", transactionID.String()), bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Invalid request body", response["error"])
	})

	t.Run("Invalid amount - zero", func(t *testing.T) {
		transactionID := uuid.New()
		salesAssociateID := uuid.New()
		billerID := uuid.New()

		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		transactionRows := sqlmock.NewRows([]string{
			"id", "biller_id", "sales_associate_id", "no_invoice", "payment_type",
			"transaction_date", "due_date", "secondary_due_date", "total_amount", "status",
			"periode", "year", "curriculum_id", "merk_buku_id", "jenjang_studi_id",
			"created_at", "updated_at",
		}).AddRow(
			transactionID, billerID, salesAssociateID, "INV2024010100000001", "K",
			time.Now(), time.Now().AddDate(0, 1, 0), nil, 500000.00, 0, 1, "2024", nil, nil, nil, time.Now(), time.Now(),
		)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "sales_transactions" WHERE id = $1`)).
			WithArgs(transactionID.String()).
			WillReturnRows(transactionRows)

		requestBody := handlers.CreatePaymentRequest{
			PaymentDate: testutil.StringPtr("2024-01-15"),
			Amount:      0,
		}

		bodyBytes, _ := json.Marshal(requestBody)
		req := httptest.NewRequest("POST", fmt.Sprintf("/sales-transactions/%s/payments", transactionID.String()), bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Amount must be greater than 0", response["error"])
	})

	t.Run("Invalid amount - negative", func(t *testing.T) {
		transactionID := uuid.New()
		salesAssociateID := uuid.New()
		billerID := uuid.New()

		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		transactionRows := sqlmock.NewRows([]string{
			"id", "biller_id", "sales_associate_id", "no_invoice", "payment_type",
			"transaction_date", "due_date", "secondary_due_date", "total_amount", "status",
			"periode", "year", "curriculum_id", "merk_buku_id", "jenjang_studi_id",
			"created_at", "updated_at",
		}).AddRow(
			transactionID, billerID, salesAssociateID, "INV2024010100000001", "K",
			time.Now(), time.Now().AddDate(0, 1, 0), nil, 500000.00, 0, 1, "2024", nil, nil, nil, time.Now(), time.Now(),
		)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "sales_transactions" WHERE id = $1`)).
			WithArgs(transactionID.String()).
			WillReturnRows(transactionRows)

		requestBody := handlers.CreatePaymentRequest{
			PaymentDate: testutil.StringPtr("2024-01-15"),
			Amount:      -100.0,
		}

		bodyBytes, _ := json.Marshal(requestBody)
		req := httptest.NewRequest("POST", fmt.Sprintf("/sales-transactions/%s/payments", transactionID.String()), bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Amount must be greater than 0", response["error"])
	})

	t.Run("Payment exceeds remaining balance", func(t *testing.T) {
		transactionID := uuid.New()
		salesAssociateID := uuid.New()
		billerID := uuid.New()

		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		transactionRows := sqlmock.NewRows([]string{
			"id", "biller_id", "sales_associate_id", "no_invoice", "payment_type",
			"transaction_date", "due_date", "secondary_due_date", "total_amount", "status",
			"periode", "year", "curriculum_id", "merk_buku_id", "jenjang_studi_id",
			"created_at", "updated_at",
		}).AddRow(
			transactionID, billerID, salesAssociateID, "INV2024010100000001", "K",
			time.Now(), time.Now().AddDate(0, 1, 0), nil, 500000.00, 0, 1, "2024", nil, nil, nil, time.Now(), time.Now(),
		)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "sales_transactions" WHERE id = $1`)).
			WithArgs(transactionID.String()).
			WillReturnRows(transactionRows)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT COALESCE`)).WillReturnError(gorm.ErrInvalidDB)

		requestBody := handlers.CreatePaymentRequest{
			PaymentDate: testutil.StringPtr("2024-01-15"),
			Amount:      500000.0,
		}

		bodyBytes, _ := json.Marshal(requestBody)
		req := httptest.NewRequest("POST", fmt.Sprintf("/sales-transactions/%s/payments", transactionID.String()), bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})
}

func TestDeletePayment(t *testing.T) {
	app := fiber.New()
	app.Delete("/sales-transactions/:transaction_id/payments/:id", handlers.DeletePayment)

	t.Run("Transaction not found", func(t *testing.T) {
		_, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)

		transactionID := uuid.New()
		paymentID := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "sales_transactions" WHERE id = $1`)).
			WithArgs(transactionID.String()).
			WillReturnError(gorm.ErrRecordNotFound)

		req := httptest.NewRequest("DELETE", fmt.Sprintf("/sales-transactions/%s/payments/%s", transactionID.String(), paymentID.String()), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Transaction not found", response["error"])
	})

	t.Run("Payment not found", func(t *testing.T) {
		transactionID := uuid.New()
		paymentID := uuid.New()
		salesAssociateID := uuid.New()
		billerID := uuid.New()

		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		transactionRows := sqlmock.NewRows([]string{
			"id", "biller_id", "sales_associate_id", "no_invoice", "payment_type",
			"transaction_date", "due_date", "secondary_due_date", "total_amount", "status",
			"periode", "year", "curriculum_id", "merk_buku_id", "jenjang_studi_id",
			"created_at", "updated_at",
		}).AddRow(
			transactionID, billerID, salesAssociateID, "INV2024010100000001", "K",
			time.Now(), time.Now().AddDate(0, 1, 0), nil, 500000.00, 0, 1, "2024", nil, nil, nil, time.Now(), time.Now(),
		)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "sales_transactions" WHERE id = $1`)).
			WithArgs(transactionID.String()).
			WillReturnRows(transactionRows)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "payments" WHERE id = $1 AND sales_transaction_id = $2`)).
			WithArgs(paymentID.String(), transactionID.String()).
			WillReturnError(gorm.ErrRecordNotFound)

		req := httptest.NewRequest("DELETE", fmt.Sprintf("/sales-transactions/%s/payments/%s", transactionID.String(), paymentID.String()), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Payment not found", response["error"])
	})

	t.Run("Successfully delete payment", func(t *testing.T) {
		transactionID := uuid.New()
		paymentID := uuid.New()
		salesAssociateID := uuid.New()
		billerID := uuid.New()

		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		transactionRows := sqlmock.NewRows([]string{
			"id", "biller_id", "sales_associate_id", "no_invoice", "payment_type",
			"transaction_date", "total_amount", "status",
			"periode", "year", "curriculum_id", "merk_buku_id", "jenjang_studi_id",
			"created_at", "updated_at",
		}).AddRow(
			transactionID, billerID, salesAssociateID, "INV2024010100000001", "K",
			time.Now(), 500000.00, 2, 1, "2024", nil, nil, nil, time.Now(), time.Now(),
		)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "sales_transactions" WHERE id = $1`)).
			WithArgs(transactionID.String()).
			WillReturnRows(transactionRows)

		paymentRows := sqlmock.NewRows(paymentColumns).
			AddRow(paymentID, transactionID, "PMT2024010100000001", time.Now(), 250000.0, 8.0, 20000.0, nil, time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "payments" WHERE id = $1 AND sales_transaction_id = $2`)).
			WithArgs(paymentID.String(), transactionID.String()).
			WillReturnRows(paymentRows)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "payments" WHERE "id" = $1`)).
			WithArgs(paymentID.String()).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT COALESCE`)).WillReturnError(gorm.ErrInvalidDB)

		req := httptest.NewRequest("DELETE", fmt.Sprintf("/sales-transactions/%s/payments/%s", transactionID.String(), paymentID.String()), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})
}
