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

func TestGetAllCities(t *testing.T) {
	db, mock, err := testutil.SetupMockDB()
	assert.NoError(t, err)
	defer testutil.CloseMockDB(db)

	app := fiber.New()
	app.Get("/cities", handlers.GetAllCities)

	t.Run("Successfully get all cities", func(t *testing.T) {
		cityID1 := uuid.New()
		cityID2 := uuid.New()

		cityRows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
			AddRow(cityID1, "Jakarta", time.Now(), time.Now()).
			AddRow(cityID2, "Bandung", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cities" ORDER BY created_at DESC`)).
			WillReturnRows(cityRows)

		req := httptest.NewRequest("GET", "/cities", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["cities"])
	})

	t.Run("Empty list", func(t *testing.T) {
		db2, mock2, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db2)

		cityRows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"})

		mock2.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cities" ORDER BY created_at DESC`)).
			WillReturnRows(cityRows)

		req := httptest.NewRequest("GET", "/cities", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["cities"])
	})
}

func TestGetCity(t *testing.T) {
	app := fiber.New()
	app.Get("/cities/:id", handlers.GetCity)

	t.Run("Successfully get city by ID", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		cityID := uuid.New()

		cityRows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
			AddRow(cityID, "Jakarta", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cities" WHERE id = $1`)).
			WithArgs(cityID.String()).
			WillReturnRows(cityRows)

		req := httptest.NewRequest("GET", "/cities/"+cityID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["city"])
	})

	t.Run("City not found", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		cityID := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cities" WHERE id = $1`)).
			WithArgs(cityID.String()).
			WillReturnError(gorm.ErrRecordNotFound)

		req := httptest.NewRequest("GET", "/cities/"+cityID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "City not found", response["error"])
	})
}

func TestCreateCity(t *testing.T) {
	app := fiber.New()
	app.Post("/cities", handlers.CreateCity)

	t.Run("Successfully create city", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		reqBody := models.City{
			Name: "Jakarta",
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "cities"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
				AddRow(uuid.New(), time.Now(), time.Now()))
		mock.ExpectCommit()

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/cities", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	})

	t.Run("Invalid request body", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/cities", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Invalid request body", response["error"])
	})
}

func TestUpdateCity(t *testing.T) {
	app := fiber.New()
	app.Put("/cities/:id", handlers.UpdateCity)

	t.Run("Successfully update city", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		cityID := uuid.New()

		cityRows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
			AddRow(cityID, "Jakarta", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cities" WHERE id = $1`)).
			WithArgs(cityID.String()).
			WillReturnRows(cityRows)

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "cities" SET .+ WHERE "id" = .+`).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		reqBody := models.City{
			ID:   cityID,
			Name: "Bandung",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("PUT", "/cities/"+cityID.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("City not found", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		cityID := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cities" WHERE id = $1`)).
			WithArgs(cityID.String()).
			WillReturnError(gorm.ErrRecordNotFound)

		reqBody := models.City{
			Name: "Bandung",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("PUT", "/cities/"+cityID.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "City not found", response["error"])
	})

	t.Run("Invalid request body", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		cityID := uuid.New()

		cityRows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
			AddRow(cityID, "Jakarta", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cities" WHERE id = $1`)).
			WithArgs(cityID.String()).
			WillReturnRows(cityRows)

		req := httptest.NewRequest("PUT", "/cities/"+cityID.String(), bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})
}

func TestDeleteCity(t *testing.T) {
	db, mock, err := testutil.SetupMockDB()
	assert.NoError(t, err)
	defer testutil.CloseMockDB(db)

	app := fiber.New()
	app.Delete("/cities/:id", handlers.DeleteCity)

	t.Run("Successfully delete city", func(t *testing.T) {
		cityID := uuid.New()

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "cities" WHERE id = $1`)).
			WithArgs(cityID.String()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		req := httptest.NewRequest("DELETE", "/cities/"+cityID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "City deleted successfully", response["message"])
	})

	t.Run("City not found", func(t *testing.T) {
		cityID := uuid.New()

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "cities" WHERE id = $1`)).
			WithArgs(cityID.String()).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		req := httptest.NewRequest("DELETE", "/cities/"+cityID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "City not found", response["error"])
	})
}
