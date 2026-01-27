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

func TestGetAllJenisBuku(t *testing.T) {
	db, mock, err := testutil.SetupMockDB()
	assert.NoError(t, err)
	defer testutil.CloseMockDB(db)

	app := fiber.New()
	app.Get("/jenis-buku", handlers.GetAllJenisBuku)

	t.Run("Successfully get all jenis buku", func(t *testing.T) {
		jenisBukuID1 := uuid.New()
		jenisBukuID2 := uuid.New()

		jenisBukuRows := sqlmock.NewRows([]string{"id", "code", "name", "description", "created_at", "updated_at"}).
			AddRow(jenisBukuID1, "JB001", "Textbook", "Test Description", time.Now(), time.Now()).
			AddRow(jenisBukuID2, "JB002", "Novel", "Test Description", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "jenis_buku" ORDER BY created_at ASC LIMIT 20`)).
			WillReturnRows(jenisBukuRows)

		countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "jenis_buku"`)).
			WillReturnRows(countRows)

		req := httptest.NewRequest("GET", "/jenis-buku", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["jenis_buku"])
		assert.NotNil(t, response["pagination"])
	})

	t.Run("Search filter by code", func(t *testing.T) {
		db2, mock2, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db2)

		jenisBukuID := uuid.New()
		description := "Test Description"

		jenisBukuRows := sqlmock.NewRows([]string{"id", "code", "name", "description", "created_at", "updated_at"}).
			AddRow(jenisBukuID, "JB001", "Textbook", &description, time.Now(), time.Now())

		mock2.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "jenis_buku" WHERE jenis_buku.code ILIKE $1 OR jenis_buku.name ILIKE $2 OR jenis_buku.description ILIKE $3 ORDER BY created_at ASC LIMIT 20`)).
			WithArgs("%JB001%", "%JB001%", "%JB001%").
			WillReturnRows(jenisBukuRows)

		countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
		mock2.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "jenis_buku"`)).
			WillReturnRows(countRows)

		req := httptest.NewRequest("GET", "/jenis-buku?search=JB001", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["jenis_buku"])
		assert.NotNil(t, response["pagination"])
	})

	t.Run("Search filter by name", func(t *testing.T) {
		db3, mock3, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db3)

		jenisBukuID := uuid.New()

		jenisBukuRows := sqlmock.NewRows([]string{"id", "code", "name", "description", "created_at", "updated_at"}).
			AddRow(jenisBukuID, "JB001", "Textbook", nil, time.Now(), time.Now())

		mock3.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "jenis_buku" WHERE jenis_buku.code ILIKE $1 OR jenis_buku.name ILIKE $2 OR jenis_buku.description ILIKE $3 ORDER BY created_at ASC LIMIT 20`)).
			WithArgs("%Textbook%", "%Textbook%", "%Textbook%").
			WillReturnRows(jenisBukuRows)

		countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
		mock3.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "jenis_buku"`)).
			WillReturnRows(countRows)

		req := httptest.NewRequest("GET", "/jenis-buku?search=Textbook", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["jenis_buku"])
		assert.NotNil(t, response["pagination"])
	})
}

func TestGetJenisBuku(t *testing.T) {
	app := fiber.New()
	app.Get("/jenis-buku/:id", handlers.GetJenisBuku)

	t.Run("Successfully get jenis buku by ID", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		jenisBukuID := uuid.New()

		jenisBukuRows := sqlmock.NewRows([]string{"id", "code", "name", "description", "created_at", "updated_at"}).
			AddRow(jenisBukuID, "JB001", "Textbook", "Test Description", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "jenis_buku" WHERE id = $1`)).
			WithArgs(jenisBukuID.String()).
			WillReturnRows(jenisBukuRows)

		req := httptest.NewRequest("GET", "/jenis-buku/"+jenisBukuID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["jenis_buku"])
	})

	t.Run("JenisBuku not found", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		jenisBukuID := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "jenis_buku" WHERE id = $1`)).
			WithArgs(jenisBukuID.String()).
			WillReturnError(gorm.ErrRecordNotFound)

		req := httptest.NewRequest("GET", "/jenis-buku/"+jenisBukuID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "JenisBuku not found", response["error"])
	})
}

func TestCreateJenisBuku(t *testing.T) {
	app := fiber.New()
	app.Post("/jenis-buku", handlers.CreateJenisBuku)

	t.Run("Successfully create jenis buku", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		reqBody := models.JenisBuku{
			Name: "Textbook",
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "jenis_buku"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
				AddRow(uuid.New(), time.Now(), time.Now()))
		mock.ExpectCommit()

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/jenis-buku", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	})

	t.Run("Invalid request body", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/jenis-buku", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Invalid request body", response["error"])
	})
}

func TestUpdateJenisBuku(t *testing.T) {
	app := fiber.New()
	app.Put("/jenis-buku/:id", handlers.UpdateJenisBuku)

	t.Run("Successfully update jenis buku", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		jenisBukuID := uuid.New()

		jenisBukuRows := sqlmock.NewRows([]string{"id", "code", "name", "description", "created_at", "updated_at"}).
			AddRow(jenisBukuID, "JB001", "Textbook", "Test Description", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "jenis_buku" WHERE id = $1`)).
			WithArgs(jenisBukuID.String()).
			WillReturnRows(jenisBukuRows)

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "jenis_buku" SET .+ WHERE "id" = .+`).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		reqBody := models.JenisBuku{
			ID:   jenisBukuID,
			Name: "Novel",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("PUT", "/jenis-buku/"+jenisBukuID.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("JenisBuku not found", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		jenisBukuID := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "jenis_buku" WHERE id = $1`)).
			WithArgs(jenisBukuID.String()).
			WillReturnError(gorm.ErrRecordNotFound)

		reqBody := models.JenisBuku{
			Name: "Novel",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("PUT", "/jenis-buku/"+jenisBukuID.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "JenisBuku not found", response["error"])
	})
}

func TestDeleteJenisBuku(t *testing.T) {
	db, mock, err := testutil.SetupMockDB()
	assert.NoError(t, err)
	defer testutil.CloseMockDB(db)

	app := fiber.New()
	app.Delete("/jenis-buku/:id", handlers.DeleteJenisBuku)

	t.Run("Successfully delete jenis buku", func(t *testing.T) {
		jenisBukuID := uuid.New()

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "jenis_buku" WHERE id = $1`)).
			WithArgs(jenisBukuID.String()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		req := httptest.NewRequest("DELETE", "/jenis-buku/"+jenisBukuID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "JenisBuku deleted successfully", response["message"])
	})

	t.Run("JenisBuku not found", func(t *testing.T) {
		jenisBukuID := uuid.New()

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "jenis_buku" WHERE id = $1`)).
			WithArgs(jenisBukuID.String()).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		req := httptest.NewRequest("DELETE", "/jenis-buku/"+jenisBukuID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "JenisBuku not found", response["error"])
	})
}
