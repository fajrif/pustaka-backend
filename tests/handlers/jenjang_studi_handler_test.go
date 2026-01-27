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

func TestGetAllJenjangStudi(t *testing.T) {
	db, mock, err := testutil.SetupMockDB()
	assert.NoError(t, err)
	defer testutil.CloseMockDB(db)

	app := fiber.New()
	app.Get("/jenjang-studi", handlers.GetAllJenjangStudi)

	t.Run("Successfully get all jenjang studi", func(t *testing.T) {
		jenjangStudiID1 := uuid.New()
		jenjangStudiID2 := uuid.New()

		rows := sqlmock.NewRows([]string{"id", "code", "name", "description", "period", "created_at", "updated_at"}).
			AddRow(jenjangStudiID1, "JS001", "SD", "Test Description", "S", time.Now(), time.Now()).
			AddRow(jenjangStudiID2, "JS002", "SMP", "Test Description", "S", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "jenjang_studi" ORDER BY created_at ASC LIMIT 20`)).
			WillReturnRows(rows)

		countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "jenjang_studi"`)).
			WillReturnRows(countRows)

		req := httptest.NewRequest("GET", "/jenjang-studi", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["jenjang_studi"])
		assert.NotNil(t, response["pagination"])
	})

	t.Run("Search filter by code", func(t *testing.T) {
		db2, mock2, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db2)

		jenjangStudiID := uuid.New()
		description := "Test Description"

		jenjangStudiRows := sqlmock.NewRows([]string{"id", "code", "name", "description", "period", "created_at", "updated_at"}).
			AddRow(jenjangStudiID, "JS001", "Elementary", &description, "S", time.Now(), time.Now())

		mock2.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "jenjang_studi" WHERE jenjang_studi.code ILIKE $1 OR jenjang_studi.name ILIKE $2 OR jenjang_studi.description ILIKE $3 ORDER BY created_at ASC LIMIT 20`)).
			WithArgs("%JS001%", "%JS001%", "%JS001%").
			WillReturnRows(jenjangStudiRows)

		countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
		mock2.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "jenjang_studi"`)).
			WillReturnRows(countRows)

		req := httptest.NewRequest("GET", "/jenjang-studi?search=JS001", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["jenjang_studi"])
		assert.NotNil(t, response["pagination"])
	})

	t.Run("Search filter by name", func(t *testing.T) {
		db3, mock3, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db3)

		jenjangStudiID := uuid.New()

		jenjangStudiRows := sqlmock.NewRows([]string{"id", "code", "name", "description", "period", "created_at", "updated_at"}).
			AddRow(jenjangStudiID, "JS001", "Elementary", nil, "S", time.Now(), time.Now())

		mock3.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "jenjang_studi" WHERE jenjang_studi.code ILIKE $1 OR jenjang_studi.name ILIKE $2 OR jenjang_studi.description ILIKE $3 ORDER BY created_at ASC LIMIT 20`)).
			WithArgs("%Elementary%", "%Elementary%", "%Elementary%").
			WillReturnRows(jenjangStudiRows)

		countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
		mock3.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "jenjang_studi"`)).
			WillReturnRows(countRows)

		req := httptest.NewRequest("GET", "/jenjang-studi?search=Elementary", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["jenjang_studi"])
		assert.NotNil(t, response["pagination"])
	})

	t.Run("Search filter by description", func(t *testing.T) {
		db4, mock4, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db4)

		jenjangStudiID := uuid.New()

		jenjangStudiRows := sqlmock.NewRows([]string{"id", "code", "name", "description", "period", "created_at", "updated_at"}).
			AddRow(jenjangStudiID, "JS001", "Elementary", "Sekolah Dasar", "S", time.Now(), time.Now())

		mock4.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "jenjang_studi" WHERE jenjang_studi.code ILIKE $1 OR jenjang_studi.name ILIKE $2 OR jenjang_studi.description ILIKE $3 ORDER BY created_at ASC LIMIT 20`)).
			WithArgs("%Sekolah Dasar%", "%Sekolah Dasar%", "%Sekolah Dasar%").
			WillReturnRows(jenjangStudiRows)

		countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
		mock4.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "jenjang_studi"`)).
			WillReturnRows(countRows)

		req := httptest.NewRequest("GET", "/jenjang-studi?search=Sekolah%20Dasar", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["jenjang_studi"])
		assert.NotNil(t, response["pagination"])
	})
}

func TestGetJenjangStudi(t *testing.T) {
	app := fiber.New()
	app.Get("/jenjang-studi/:id", handlers.GetJenjangStudi)

	t.Run("Successfully get jenjang studi by ID", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		jenjangStudiID := uuid.New()

		rows := sqlmock.NewRows([]string{"id", "code", "name", "description", "period", "created_at", "updated_at"}).
			AddRow(jenjangStudiID, "JS001", "SD", "Test Description", "S", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "jenjang_studi" WHERE id = $1`)).
			WithArgs(jenjangStudiID.String()).
			WillReturnRows(rows)

		req := httptest.NewRequest("GET", "/jenjang-studi/"+jenjangStudiID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["jenjang_studi"])
	})

	t.Run("JenjangStudi not found", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		jenjangStudiID := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "jenjang_studi" WHERE id = $1`)).
			WithArgs(jenjangStudiID.String()).
			WillReturnError(gorm.ErrRecordNotFound)

		req := httptest.NewRequest("GET", "/jenjang-studi/"+jenjangStudiID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "JenjangStudi not found", response["error"])
	})
}

func TestCreateJenjangStudi(t *testing.T) {
	app := fiber.New()
	app.Post("/jenjang-studi", handlers.CreateJenjangStudi)

	t.Run("Successfully create jenjang studi", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		reqBody := models.JenjangStudi{
			Name: "SD",
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "jenjang_studi"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
				AddRow(uuid.New(), time.Now(), time.Now()))
		mock.ExpectCommit()

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/jenjang-studi", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	})

	t.Run("Invalid request body", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/jenjang-studi", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Invalid request body", response["error"])
	})
}

func TestUpdateJenjangStudi(t *testing.T) {
	app := fiber.New()
	app.Put("/jenjang-studi/:id", handlers.UpdateJenjangStudi)

	t.Run("Successfully update jenjang studi", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		jenjangStudiID := uuid.New()

		rows := sqlmock.NewRows([]string{"id", "code", "name", "description", "period", "created_at", "updated_at"}).
			AddRow(jenjangStudiID, "JS001", "SD", "Test Description", "S", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "jenjang_studi" WHERE id = $1`)).
			WithArgs(jenjangStudiID.String()).
			WillReturnRows(rows)

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "jenjang_studi" SET .+ WHERE "id" = .+`).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		reqBody := models.JenjangStudi{
			ID:   jenjangStudiID,
			Name: "SMP",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("PUT", "/jenjang-studi/"+jenjangStudiID.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("JenjangStudi not found", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		jenjangStudiID := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "jenjang_studi" WHERE id = $1`)).
			WithArgs(jenjangStudiID.String()).
			WillReturnError(gorm.ErrRecordNotFound)

		reqBody := models.JenjangStudi{
			Name: "SMP",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("PUT", "/jenjang-studi/"+jenjangStudiID.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "JenjangStudi not found", response["error"])
	})
}

func TestDeleteJenjangStudi(t *testing.T) {
	db, mock, err := testutil.SetupMockDB()
	assert.NoError(t, err)
	defer testutil.CloseMockDB(db)

	app := fiber.New()
	app.Delete("/jenjang-studi/:id", handlers.DeleteJenjangStudi)

	t.Run("Successfully delete jenjang studi", func(t *testing.T) {
		jenjangStudiID := uuid.New()

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "jenjang_studi" WHERE id = $1`)).
			WithArgs(jenjangStudiID.String()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		req := httptest.NewRequest("DELETE", "/jenjang-studi/"+jenjangStudiID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "JenjangStudi deleted successfully", response["message"])
	})

	t.Run("JenjangStudi not found", func(t *testing.T) {
		jenjangStudiID := uuid.New()

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "jenjang_studi" WHERE id = $1`)).
			WithArgs(jenjangStudiID.String()).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		req := httptest.NewRequest("DELETE", "/jenjang-studi/"+jenjangStudiID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "JenjangStudi not found", response["error"])
	})
}
