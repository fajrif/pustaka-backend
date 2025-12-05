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

func TestGetAllKelas(t *testing.T) {
	db, mock, err := testutil.SetupMockDB()
	assert.NoError(t, err)
	defer testutil.CloseMockDB(db)

	app := fiber.New()
	app.Get("/kelas", handlers.GetAllKelas)

	t.Run("Successfully get all kelas", func(t *testing.T) {
		kelasID1 := uuid.New()
		kelasID2 := uuid.New()

		rows := sqlmock.NewRows([]string{"id", "code", "name", "description", "created_at", "updated_at"}).
			AddRow(kelasID1, "K001", "Class 1", "Test Description", time.Now(), time.Now()).
			AddRow(kelasID2, "K002", "Class 2", "Test Description", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "kelas" ORDER BY created_at DESC`)).
			WillReturnRows(rows)

		req := httptest.NewRequest("GET", "/kelas", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["kelas"])
	})

	t.Run("Search filter by code", func(t *testing.T) {
		db2, mock2, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db2)

		kelasID := uuid.New()
		description := "Test Description"

		kelasRows := sqlmock.NewRows([]string{"id", "code", "name", "description", "created_at", "updated_at"}).
			AddRow(kelasID, "K001", "Grade 1", &description, time.Now(), time.Now())

		mock2.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "kelas" WHERE kelas.code ILIKE $1 OR kelas.name ILIKE $2 OR kelas.description ILIKE $3 ORDER BY created_at DESC`)).
			WithArgs("%K001%", "%K001%", "%K001%").
			WillReturnRows(kelasRows)

		req := httptest.NewRequest("GET", "/kelas?search=K001", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["kelas"])
	})

	t.Run("Search filter by name", func(t *testing.T) {
		db3, mock3, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db3)

		kelasID := uuid.New()

		kelasRows := sqlmock.NewRows([]string{"id", "code", "name", "description", "created_at", "updated_at"}).
			AddRow(kelasID, "K001", "Grade", nil, time.Now(), time.Now())

		mock3.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "kelas" WHERE kelas.code ILIKE $1 OR kelas.name ILIKE $2 OR kelas.description ILIKE $3 ORDER BY created_at DESC`)).
			WithArgs("%Grade%", "%Grade%", "%Grade%").
			WillReturnRows(kelasRows)

		req := httptest.NewRequest("GET", "/kelas?search=Grade", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["kelas"])
	})
}

func TestGetKelas(t *testing.T) {
	app := fiber.New()
	app.Get("/kelas/:id", handlers.GetKelas)

	t.Run("Successfully get kelas by ID", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		kelasID := uuid.New()

		rows := sqlmock.NewRows([]string{"id", "code", "name", "description", "created_at", "updated_at"}).
			AddRow(kelasID, "K001", "Class 1", "Test Description", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "kelas" WHERE id = $1`)).
			WithArgs(kelasID.String()).
			WillReturnRows(rows)

		req := httptest.NewRequest("GET", "/kelas/"+kelasID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["kelas"])
	})

	t.Run("Kelas not found", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		kelasID := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "kelas" WHERE id = $1`)).
			WithArgs(kelasID.String()).
			WillReturnError(gorm.ErrRecordNotFound)

		req := httptest.NewRequest("GET", "/kelas/"+kelasID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Kelas not found", response["error"])
	})
}

func TestCreateKelas(t *testing.T) {
	app := fiber.New()
	app.Post("/kelas", handlers.CreateKelas)

	t.Run("Successfully create kelas", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		reqBody := models.Kelas{
			Name: "Class 1",
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "kelas"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
				AddRow(uuid.New(), time.Now(), time.Now()))
		mock.ExpectCommit()

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/kelas", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	})

	t.Run("Invalid request body", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/kelas", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Invalid request body", response["error"])
	})
}

func TestUpdateKelas(t *testing.T) {
	app := fiber.New()
	app.Put("/kelas/:id", handlers.UpdateKelas)

	t.Run("Successfully update kelas", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		kelasID := uuid.New()

		rows := sqlmock.NewRows([]string{"id", "code", "name", "description", "created_at", "updated_at"}).
			AddRow(kelasID, "K001", "Class 1", "Test Description", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "kelas" WHERE id = $1`)).
			WithArgs(kelasID.String()).
			WillReturnRows(rows)

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "kelas" SET .+ WHERE "id" = .+`).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		reqBody := models.Kelas{
			ID:   kelasID,
			Name: "Class 2",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("PUT", "/kelas/"+kelasID.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("Kelas not found", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		kelasID := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "kelas" WHERE id = $1`)).
			WithArgs(kelasID.String()).
			WillReturnError(gorm.ErrRecordNotFound)

		reqBody := models.Kelas{
			Name: "Class 2",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("PUT", "/kelas/"+kelasID.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Kelas not found", response["error"])
	})
}

func TestDeleteKelas(t *testing.T) {
	db, mock, err := testutil.SetupMockDB()
	assert.NoError(t, err)
	defer testutil.CloseMockDB(db)

	app := fiber.New()
	app.Delete("/kelas/:id", handlers.DeleteKelas)

	t.Run("Successfully delete kelas", func(t *testing.T) {
		kelasID := uuid.New()

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "kelas" WHERE id = $1`)).
			WithArgs(kelasID.String()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		req := httptest.NewRequest("DELETE", "/kelas/"+kelasID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Kelas deleted successfully", response["message"])
	})

	t.Run("Kelas not found", func(t *testing.T) {
		kelasID := uuid.New()

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "kelas" WHERE id = $1`)).
			WithArgs(kelasID.String()).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		req := httptest.NewRequest("DELETE", "/kelas/"+kelasID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Kelas not found", response["error"])
	})
}
