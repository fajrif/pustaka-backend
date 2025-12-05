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

		cityRows := sqlmock.NewRows([]string{"id", "code", "name", "created_at", "updated_at"}).
			AddRow(cityID1, "JKT", "Jakarta", time.Now(), time.Now()).
			AddRow(cityID2, "BDG", "Bandung", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cities" ORDER BY created_at DESC`)).
			WillReturnRows(cityRows)

		countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "cities"`)).
			WillReturnRows(countRows)

		req := httptest.NewRequest("GET", "/cities", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["cities"])
		assert.NotNil(t, response["pagination"])
	})

	t.Run("Empty list", func(t *testing.T) {
		db2, mock2, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db2)

		cityRows := sqlmock.NewRows([]string{"id", "code", "name", "created_at", "updated_at"})

		mock2.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cities" ORDER BY created_at DESC LIMIT 20`)).
			WillReturnRows(cityRows)

		countRows := sqlmock.NewRows([]string{"count"}).AddRow(0)
		mock2.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "cities"`)).
			WillReturnRows(countRows)

		req := httptest.NewRequest("GET", "/cities", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["cities"])
		assert.NotNil(t, response["pagination"])
	})

	t.Run("Search filter by code", func(t *testing.T) {
		db3, mock3, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db3)

		cityID := uuid.New()

		cityRows := sqlmock.NewRows([]string{"id", "code", "name", "created_at", "updated_at"}).
			AddRow(cityID, "JKT", "Jakarta", time.Now(), time.Now())

		mock3.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cities" WHERE cities.code ILIKE $1 OR cities.name ILIKE $2 ORDER BY created_at DESC LIMIT 20`)).
			WithArgs("%JKT%", "%JKT%").
			WillReturnRows(cityRows)

		countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
		mock3.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "cities" WHERE cities.code ILIKE $1 OR cities.name ILIKE $2`)).
			WithArgs("%JKT%", "%JKT%").
			WillReturnRows(countRows)

		req := httptest.NewRequest("GET", "/cities?search=JKT", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["cities"])
		assert.NotNil(t, response["pagination"])
	})

	t.Run("Search filter by name", func(t *testing.T) {
		db4, mock4, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db4)

		cityID := uuid.New()

		cityRows := sqlmock.NewRows([]string{"id", "code", "name", "created_at", "updated_at"}).
			AddRow(cityID, "JKT", "Jakarta", time.Now(), time.Now())

		mock4.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cities" WHERE cities.code ILIKE $1 OR cities.name ILIKE $2 ORDER BY created_at DESC LIMIT 20`)).
			WithArgs("%Jakarta%", "%Jakarta%").
			WillReturnRows(cityRows)

		countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
		mock4.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "cities" WHERE cities.code ILIKE $1 OR cities.name ILIKE $2`)).
			WithArgs("%Jakarta%", "%Jakarta%").
			WillReturnRows(countRows)

		req := httptest.NewRequest("GET", "/cities?search=Jakarta", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["cities"])
		assert.NotNil(t, response["pagination"])
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

		cityRows := sqlmock.NewRows([]string{"id", "code", "name", "created_at", "updated_at"}).
			AddRow(cityID, "JKT", "Jakarta", time.Now(), time.Now())

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
			Code: "JKT",
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

		cityRows := sqlmock.NewRows([]string{"id", "code", "name", "created_at", "updated_at"}).
			AddRow(cityID, "JKT", "Jakarta", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cities" WHERE id = $1`)).
			WithArgs(cityID.String()).
			WillReturnRows(cityRows)

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "cities" SET .+ WHERE "id" = .+`).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		reqBody := models.City{
			ID:   cityID,
			Code: "BDG",
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
			Code: "BDG",
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

		cityRows := sqlmock.NewRows([]string{"id", "code", "name", "created_at", "updated_at"}).
			AddRow(cityID, "JKT", "Jakarta", time.Now(), time.Now())

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
