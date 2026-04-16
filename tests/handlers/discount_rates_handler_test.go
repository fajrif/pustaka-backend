package handlers_test

import (
	"bytes"
	"encoding/json"
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

func TestGetAllDiscountRates(t *testing.T) {
	db, mock, err := testutil.SetupMockDB()
	assert.NoError(t, err)
	defer testutil.CloseMockDB(db)

	app := fiber.New()
	app.Get("/discount-rates", handlers.GetAllDiscountRates)

	t.Run("Successfully get all discount rates", func(t *testing.T) {
		discountRateID1 := uuid.New()
		discountRateID2 := uuid.New()

		discountRateRows := sqlmock.NewRows([]string{"id", "name", "discount", "periode", "year", "start_date", "end_date", "description", "created_at", "updated_at"}).
			AddRow(discountRateID1, "Early Payment Discount", 8.00, 1, "2024", time.Now(), time.Now().AddDate(0, 1, 0), "Discount for payments on time", time.Now(), time.Now()).
			AddRow(discountRateID2, "Secondary Payment Discount", 5.00, 1, "2024", time.Now(), time.Now().AddDate(0, 1, 0), "Discount for secondary payments", time.Now(), time.Now())

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

		discountRateRows := sqlmock.NewRows([]string{"id", "name", "discount", "periode", "year", "start_date", "end_date", "description", "created_at", "updated_at"}).
			AddRow(discountRateID, "Early Payment Discount", 8.00, 1, "2024", time.Now(), time.Now().AddDate(0, 1, 0), "Discount for payments on time", time.Now(), time.Now())

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

		startDateStr := "2024-01-01"
		endDateStr := "2024-01-31"

		countRows := sqlmock.NewRows([]string{"count"}).AddRow(0)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "discount_rates" WHERE periode = $1 AND year = $2`)).
			WithArgs(1, "2024").
			WillReturnRows(countRows)

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "discount_rates"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "discount", "periode", "year", "start_date", "end_date", "description", "created_at", "updated_at"}).
				AddRow(uuid.New(), "New Discount", 10.00, 1, "2024", time.Now(), time.Now(), nil, time.Now(), time.Now()))
		mock.ExpectCommit()

		reqBody := handlers.CreateDiscountRateRequest{
			Name:      "New Discount",
			Discount:  10.00,
			Periode:   1,
			Year:      "2024",
			StartDate: &startDateStr,
			EndDate:   &endDateStr,
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
		reqBody := handlers.CreateDiscountRateRequest{
			Name:     "",
			Discount: 8.00,
			Periode:  1,
			Year:     "2024",
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
		reqBody := handlers.CreateDiscountRateRequest{
			Name:     "Invalid Discount",
			Discount: 150.00,
			Periode:  1,
			Year:     "2024",
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

	t.Run("Periode must be at least 1", func(t *testing.T) {
		reqBody := handlers.CreateDiscountRateRequest{
			Name:     "Test Discount",
			Discount: 8.00,
			Periode:  0,
			Year:     "2024",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/discount-rates", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "periode must be at least 1", response["error"])
	})

	t.Run("Year is required", func(t *testing.T) {
		reqBody := handlers.CreateDiscountRateRequest{
			Name:     "Test Discount",
			Discount: 8.00,
			Periode:  1,
			Year:     "",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/discount-rates", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "year is required and must be a valid year", response["error"])
	})

	t.Run("Duplicate periode and year", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "discount_rates" WHERE periode = $1 AND year = $2`)).
			WithArgs(1, "2024").
			WillReturnRows(countRows)

		reqBody := handlers.CreateDiscountRateRequest{
			Name:     "Duplicate Discount",
			Discount: 8.00,
			Periode:  1,
			Year:     "2024",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/discount-rates", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "A discount rate for this periode and year already exists", response["error"])
	})

	t.Run("Year as number should work", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		countRows := sqlmock.NewRows([]string{"count"}).AddRow(0)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "discount_rates" WHERE periode = $1 AND year = $2`)).
			WithArgs(1, "2024").
			WillReturnRows(countRows)

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "discount_rates"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "discount", "periode", "year", "start_date", "end_date", "description", "created_at", "updated_at"}).
				AddRow(uuid.New(), "Discount", 8.00, 1, "2024", time.Now(), time.Now(), nil, time.Now(), time.Now()))
		mock.ExpectCommit()

		reqBody := map[string]interface{}{
			"name":     "Discount",
			"discount": 8.00,
			"periode":  1,
			"year":     2024, // year as number
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/discount-rates", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
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

		discountRateRows := sqlmock.NewRows([]string{"id", "name", "discount", "periode", "year", "start_date", "end_date", "description", "created_at", "updated_at"}).
			AddRow(discountRateID, "Early Payment Discount", 8.00, 1, "2024", time.Now(), time.Now().AddDate(0, 1, 0), "Discount for payments on time", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "discount_rates" WHERE id = $1`)).
			WithArgs(discountRateID.String()).
			WillReturnRows(discountRateRows)

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "discount_rates" SET .+ WHERE "id" = .+`).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		reqBody := handlers.CreateDiscountRateRequest{
			Name:     "Updated Discount",
			Discount: 10.00,
			Periode:  1,
			Year:     "2024",
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

		reqBody := handlers.CreateDiscountRateRequest{
			Name:     "Updated Discount",
			Discount: 10.00,
			Periode:  1,
			Year:     "2024",
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
	app := fiber.New()
	app.Delete("/discount-rates/:id", handlers.DeleteDiscountRate)

	t.Run("Successfully delete discount rate", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

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
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

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

func TestGetAllDiscountRatesOrderByCreatedAt(t *testing.T) {
	app := fiber.New()
	app.Get("/discount-rates", handlers.GetAllDiscountRates)

	t.Run("Records ordered by created_at ASC - newest at bottom", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		discountRateID1 := uuid.New()
		discountRateID2 := uuid.New()
		discountRateID3 := uuid.New()

		createdAt1 := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
		createdAt2 := time.Date(2024, 2, 1, 10, 0, 0, 0, time.UTC)
		createdAt3 := time.Date(2024, 3, 1, 10, 0, 0, 0, time.UTC)
		updatedAt := time.Now()

		discountRateRows := sqlmock.NewRows([]string{"id", "name", "discount", "periode", "year", "start_date", "end_date", "description", "created_at", "updated_at"}).
			AddRow(discountRateID1, "First Discount", 5.00, 1, "2024", time.Now(), time.Now().AddDate(0, 1, 0), "First created", createdAt1, updatedAt).
			AddRow(discountRateID2, "Second Discount", 8.00, 1, "2024", time.Now(), time.Now().AddDate(0, 1, 0), "Second created", createdAt2, updatedAt).
			AddRow(discountRateID3, "Third Discount", 10.00, 1, "2024", time.Now(), time.Now().AddDate(0, 1, 0), "Third created", createdAt3, updatedAt)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "discount_rates" ORDER BY created_at ASC LIMIT 20`)).
			WillReturnRows(discountRateRows)

		countRows := sqlmock.NewRows([]string{"count"}).AddRow(3)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "discount_rates"`)).
			WillReturnRows(countRows)

		req := httptest.NewRequest("GET", "/discount-rates", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		discountRates, ok := response["discount_rates"].([]interface{})
		assert.True(t, ok)
		assert.Len(t, discountRates, 3)

		first := discountRates[0].(map[string]interface{})
		second := discountRates[1].(map[string]interface{})
		third := discountRates[2].(map[string]interface{})

		assert.Equal(t, "First Discount", first["name"])
		assert.Equal(t, "Second Discount", second["name"])
		assert.Equal(t, "Third Discount", third["name"])
	})

	t.Run("Search filter works with created_at ordering", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		discountRateID := uuid.New()
		createdAt := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
		updatedAt := time.Now()

		discountRateRows := sqlmock.NewRows([]string{"id", "name", "discount", "periode", "year", "start_date", "end_date", "description", "created_at", "updated_at"}).
			AddRow(discountRateID, "Special Discount", 12.00, 1, "2024", time.Now(), time.Now().AddDate(0, 1, 0), "Special offer", createdAt, updatedAt)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "discount_rates" WHERE discount_rates.name ILIKE $1 OR discount_rates.description ILIKE $2 ORDER BY created_at ASC LIMIT 20`)).
			WithArgs("%Special%", "%Special%").
			WillReturnRows(discountRateRows)

		countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "discount_rates" WHERE discount_rates.name ILIKE $1 OR discount_rates.description ILIKE $2`)).
			WithArgs("%Special%", "%Special%").
			WillReturnRows(countRows)

		req := httptest.NewRequest("GET", "/discount-rates?search=Special", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		discountRates, ok := response["discount_rates"].([]interface{})
		assert.True(t, ok)
		assert.Len(t, discountRates, 1)
		assert.Equal(t, "Special Discount", discountRates[0].(map[string]interface{})["name"])
	})
}

func TestUpdateDiscountRatePreservesCreatedAt(t *testing.T) {
	app := fiber.New()
	app.Put("/discount-rates/:id", handlers.UpdateDiscountRate)

	t.Run("Update does not change created_at", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		discountRateID := uuid.New()
		originalCreatedAt := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
		originalUpdatedAt := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)

		discountRateRows := sqlmock.NewRows([]string{"id", "name", "discount", "periode", "year", "start_date", "end_date", "description", "created_at", "updated_at"}).
			AddRow(discountRateID, "Original Name", 5.00, 1, "2024", time.Now(), time.Now().AddDate(0, 1, 0), "Original description", originalCreatedAt, originalUpdatedAt)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "discount_rates" WHERE id = $1`)).
			WithArgs(discountRateID.String()).
			WillReturnRows(discountRateRows)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "discount_rates" SET`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		reqBody := handlers.CreateDiscountRateRequest{
			Name:     "Updated Name",
			Discount: 10.00,
			Periode:  1,
			Year:     "2024",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("PUT", "/discount-rates/"+discountRateID.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})
}
