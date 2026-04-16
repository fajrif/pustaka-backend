package handlers_test

import (
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

func TestGetSalesTransactionDiscountValue(t *testing.T) {
	app := fiber.New()
	app.Get("/sales-transactions/:sales_transaction_id/discount-value", handlers.GetSalesTransactionDiscountValue)

	t.Run("Transaction not found", func(t *testing.T) {
		_, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)

		transactionID := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "sales_transactions" WHERE id = $1`)).
			WithArgs(transactionID.String()).
			WillReturnError(gorm.ErrRecordNotFound)

		req := httptest.NewRequest("GET", fmt.Sprintf("/sales-transactions/%s/discount-value?date=2024-01-15", transactionID.String()), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Transaction not found", response["error"])
	})

	t.Run("Missing date parameter", func(t *testing.T) {
		transactionID := uuid.New()

		req := httptest.NewRequest("GET", fmt.Sprintf("/sales-transactions/%s/discount-value", transactionID.String()), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "date query parameter is required", response["error"])
	})

	t.Run("Invalid date format", func(t *testing.T) {
		transactionID := uuid.New()

		req := httptest.NewRequest("GET", fmt.Sprintf("/sales-transactions/%s/discount-value?date=invalid-date", transactionID.String()), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Invalid date format. Use ISO 8601 (YYYY-MM-DD or RFC3339)", response["error"])
	})

	t.Run("Discount only applicable for credit transactions", func(t *testing.T) {
		transactionID := uuid.New()
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
			transactionID, billerID, salesAssociateID, "INV2024010100000001", "T", // Cash transaction
			time.Now(), 500000.00, 0, 1, "2024", nil, nil, nil, time.Now(), time.Now(),
		)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "sales_transactions" WHERE id = $1`)).
			WithArgs(transactionID.String()).
			WillReturnRows(transactionRows)

		req := httptest.NewRequest("GET", fmt.Sprintf("/sales-transactions/%s/discount-value?date=2024-01-15", transactionID.String()), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Discount calculation is only applicable for credit transactions", response["error"])
	})

	t.Run("No applicable discount rate found", func(t *testing.T) {
		transactionID := uuid.New()
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
			time.Now(), 500000.00, 0, 1, "2024", nil, nil, nil, time.Now(), time.Now(),
		)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "sales_transactions" WHERE id = $1`)).
			WithArgs(transactionID.String()).
			WillReturnRows(transactionRows)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "discount_rates" WHERE periode = $1 AND year = $2`)).
			WithArgs(1, "2024").
			WillReturnError(gorm.ErrRecordNotFound)

		req := httptest.NewRequest("GET", fmt.Sprintf("/sales-transactions/%s/discount-value?date=2024-01-15", transactionID.String()), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "No applicable discount rate found for this transaction's periode, year, and payment date", response["error"])
	})

	t.Run("Successfully get discount value", func(t *testing.T) {
		transactionID := uuid.New()
		discountRateID := uuid.New()
		salesAssociateID := uuid.New()
		billerID := uuid.New()

		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)

		transactionRows := sqlmock.NewRows([]string{
			"id", "biller_id", "sales_associate_id", "no_invoice", "payment_type",
			"transaction_date", "total_amount", "status",
			"periode", "year", "curriculum_id", "merk_buku_id", "jenjang_studi_id",
			"created_at", "updated_at",
		}).AddRow(
			transactionID, billerID, salesAssociateID, "INV2024010100000001", "K",
			time.Now(), 500000.00, 0, 1, "2024", nil, nil, nil, time.Now(), time.Now(),
		)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "sales_transactions" WHERE id = $1`)).
			WithArgs(transactionID.String()).
			WillReturnRows(transactionRows)

		discountRateRows := sqlmock.NewRows([]string{
			"id", "name", "discount", "periode", "year", "start_date", "end_date", "description", "created_at", "updated_at",
		}).AddRow(
			discountRateID, "Early Payment Discount", 8.00, 1, "2024", startDate, endDate, "Discount for early payment", time.Now(), time.Now(),
		)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "discount_rates" WHERE`)).
			WithArgs(1, "2024", sqlmock.AnyArg()).
			WillReturnRows(discountRateRows)

		paymentDate := "2024-01-15"
		req := httptest.NewRequest("GET", fmt.Sprintf("/sales-transactions/%s/discount-value?date=%s", transactionID.String(), paymentDate), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, transactionID.String(), response["sales_transaction_id"])
		assert.Equal(t, float64(500000.00), response["total_amount"])
		assert.Equal(t, float64(8.0), response["discount_percentage"])
		assert.Equal(t, float64(40000.0), response["discount_amount"])        // 500000 * 8%
		assert.Equal(t, float64(460000.0), response["amount_after_discount"]) // 500000 - 40000

		discountRate := response["discount_rate"].(map[string]interface{})
		assert.Equal(t, discountRateID.String(), discountRate["id"])
		assert.Equal(t, "Early Payment Discount", discountRate["name"])
		assert.Equal(t, float64(8.0), discountRate["discount_percentage"])
	})

	t.Run("Discount calculation with RFC3339 date format", func(t *testing.T) {
		transactionID := uuid.New()
		discountRateID := uuid.New()
		salesAssociateID := uuid.New()
		billerID := uuid.New()

		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)

		transactionRows := sqlmock.NewRows([]string{
			"id", "biller_id", "sales_associate_id", "no_invoice", "payment_type",
			"transaction_date", "total_amount", "status",
			"periode", "year", "curriculum_id", "merk_buku_id", "jenjang_studi_id",
			"created_at", "updated_at",
		}).AddRow(
			transactionID, billerID, salesAssociateID, "INV2024010100000001", "K",
			time.Now(), 1000000.00, 0, 2, "2024", nil, nil, nil, time.Now(), time.Now(),
		)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "sales_transactions" WHERE id = $1`)).
			WithArgs(transactionID.String()).
			WillReturnRows(transactionRows)

		discountRateRows := sqlmock.NewRows([]string{
			"id", "name", "discount", "periode", "year", "start_date", "end_date", "description", "created_at", "updated_at",
		}).AddRow(
			discountRateID, "Full Year Discount", 10.00, 2, "2024", startDate, endDate, "Discount for full year", time.Now(), time.Now(),
		)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "discount_rates" WHERE`)).
			WithArgs(2, "2024", sqlmock.AnyArg()).
			WillReturnRows(discountRateRows)

		paymentDate := "2024-06-15T10:30:00Z"
		req := httptest.NewRequest("GET", fmt.Sprintf("/sales-transactions/%s/discount-value?date=%s", transactionID.String(), paymentDate), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, float64(10.0), response["discount_percentage"])
		assert.Equal(t, float64(100000.0), response["discount_amount"])       // 1000000 * 10%
		assert.Equal(t, float64(900000.0), response["amount_after_discount"]) // 1000000 - 100000
		assert.Equal(t, 2, int(response["periode"].(float64)))
		assert.Equal(t, "2024", response["year"])
	})
}
