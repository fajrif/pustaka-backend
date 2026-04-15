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

func TestGetAllDiscountRates(t *testing.T) {
	db, mock, err := testutil.SetupMockDB()
	assert.NoError(t, err)
	defer testutil.CloseMockDB(db)

	app := fiber.New()
	app.Get("/discount-rates", handlers.GetAllDiscountRates)

	t.Run("Successfully get all discount rates", func(t *testing.T) {
		discountRateID1 := uuid.New()
		discountRateID2 := uuid.New()

		discountRateRows := sqlmock.NewRows([]string{"id", "name", "discount", "description", "created_at", "updated_at"}).
			AddRow(discountRateID1, "Early Payment Discount", 8.00, "Discount for payments on time", time.Now(), time.Now()).
			AddRow(discountRateID2, "Secondary Payment Discount", 5.00, "Discount for secondary payments", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "discount_rates" ORDER BY created_at ASC LIMIT 20`)).
			WillReturnRows(discountRateRows)

		countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "discount_rates"`)).
			WillReturnRows(countRows)

		req := httptest.NewRequest("GET", "/discount-rates", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["discount_rates"])
		assert.NotNil(t, response["pagination"])
	})
}

func TestGetDiscountRate(t *testing.T) {
	app := fiber.New()
	app.Get("/discount-rates/:id", handlers.GetDiscountRate)

	t.Run("Successfully get discount rate by ID", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		discountRateID := uuid.New()

		discountRateRows := sqlmock.NewRows([]string{"id", "name", "discount", "description", "created_at", "updated_at"}).
			AddRow(discountRateID, "Early Payment Discount", 8.00, "Discount for payments on time", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "discount_rates" WHERE id = $1`)).
			WithArgs(discountRateID.String()).
			WillReturnRows(discountRateRows)

		req := httptest.NewRequest("GET", "/discount-rates/"+discountRateID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["discount_rate"])
	})

	t.Run("DiscountRate not found", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		discountRateID := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "discount_rates" WHERE id = $1`)).
			WithArgs(discountRateID.String()).
			WillReturnError(gorm.ErrRecordNotFound)

		req := httptest.NewRequest("GET", "/discount-rates/"+discountRateID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "DiscountRate not found", response["error"])
	})
}

func TestCreateDiscountRate(t *testing.T) {
	app := fiber.New()
	app.Post("/discount-rates", handlers.CreateDiscountRate)

	t.Run("Successfully create discount rate", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "discount_rates"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "discount", "description", "created_at", "updated_at"}).
				AddRow(uuid.New(), "New Discount", 10.00, nil, time.Now(), time.Now()))
		mock.ExpectCommit()

		reqBody := models.DiscountRate{
			Name:     "New Discount",
			Discount: 10.00,
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/discount-rates", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	})

	t.Run("Invalid request body", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/discount-rates", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Invalid request body", response["error"])
	})

	t.Run("Name is required", func(t *testing.T) {
		reqBody := models.DiscountRate{
			Name:     "",
			Discount: 8.00,
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/discount-rates", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "name is required", response["error"])
	})

	t.Run("Discount must be between 0 and 100", func(t *testing.T) {
		reqBody := models.DiscountRate{
			Name:     "Invalid Discount",
			Discount: 150.00,
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/discount-rates", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "discount must be between 0 and 100", response["error"])
	})
}

func TestUpdateDiscountRate(t *testing.T) {
	app := fiber.New()
	app.Put("/discount-rates/:id", handlers.UpdateDiscountRate)

	t.Run("Successfully update discount rate", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		discountRateID := uuid.New()

		discountRateRows := sqlmock.NewRows([]string{"id", "name", "discount", "description", "created_at", "updated_at"}).
			AddRow(discountRateID, "Early Payment Discount", 8.00, "Discount for payments on time", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "discount_rates" WHERE id = $1`)).
			WithArgs(discountRateID.String()).
			WillReturnRows(discountRateRows)

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "discount_rates" SET .+ WHERE "id" = .+`).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		reqBody := models.DiscountRate{
			ID:       discountRateID,
			Name:     "Updated Discount",
			Discount: 10.00,
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("PUT", "/discount-rates/"+discountRateID.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("DiscountRate not found", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		discountRateID := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "discount_rates" WHERE id = $1`)).
			WithArgs(discountRateID.String()).
			WillReturnError(gorm.ErrRecordNotFound)

		reqBody := models.DiscountRate{
			Name:     "Updated Discount",
			Discount: 10.00,
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("PUT", "/discount-rates/"+discountRateID.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "DiscountRate not found", response["error"])
	})
}

func TestDeleteDiscountRate(t *testing.T) {
	db, mock, err := testutil.SetupMockDB()
	assert.NoError(t, err)
	defer testutil.CloseMockDB(db)

	app := fiber.New()
	app.Delete("/discount-rates/:id", handlers.DeleteDiscountRate)

	t.Run("Successfully delete discount rate", func(t *testing.T) {
		discountRateID := uuid.New()

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "discount_rates" WHERE id = $1`)).
			WithArgs(discountRateID.String()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		req := httptest.NewRequest("DELETE", "/discount-rates/"+discountRateID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "DiscountRate deleted successfully", response["message"])
	})

	t.Run("DiscountRate not found", func(t *testing.T) {
		discountRateID := uuid.New()

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "discount_rates" WHERE id = $1`)).
			WithArgs(discountRateID.String()).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		req := httptest.NewRequest("DELETE", "/discount-rates/"+discountRateID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "DiscountRate not found", response["error"])
	})
}
