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

func TestGetAllMerkBuku(t *testing.T) {
	db, mock, err := testutil.SetupMockDB()
	assert.NoError(t, err)
	defer testutil.CloseMockDB(db)

	app := fiber.New()
	app.Get("/merk-buku", handlers.GetAllMerkBuku)

	t.Run("Successfully get all merk buku", func(t *testing.T) {
		merkBukuID1 := uuid.New()
		merkBukuID2 := uuid.New()

		merkBukuRows := sqlmock.NewRows([]string{"id", "code", "name", "description", "created_at", "updated_at"}).
			AddRow(merkBukuID1, "JB001", "Textbook", "Test Description", time.Now(), time.Now()).
			AddRow(merkBukuID2, "JB002", "Novel", "Test Description", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "merk_buku" ORDER BY created_at ASC LIMIT 20`)).
			WillReturnRows(merkBukuRows)

		countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "merk_buku"`)).
			WillReturnRows(countRows)

		req := httptest.NewRequest("GET", "/merk-buku", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["merk_buku"])
		assert.NotNil(t, response["pagination"])
	})

	t.Run("Search filter by code", func(t *testing.T) {
		db2, mock2, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db2)

		merkBukuID := uuid.New()
		description := "Test Description"

		merkBukuRows := sqlmock.NewRows([]string{"id", "code", "name", "description", "created_at", "updated_at"}).
			AddRow(merkBukuID, "JB001", "Textbook", &description, time.Now(), time.Now())

		mock2.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "merk_buku" WHERE merk_buku.code ILIKE $1 OR merk_buku.name ILIKE $2 OR merk_buku.description ILIKE $3 ORDER BY created_at ASC LIMIT 20`)).
			WithArgs("%JB001%", "%JB001%", "%JB001%").
			WillReturnRows(merkBukuRows)

		countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
		mock2.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "merk_buku"`)).
			WillReturnRows(countRows)

		req := httptest.NewRequest("GET", "/merk-buku?search=JB001", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["merk_buku"])
		assert.NotNil(t, response["pagination"])
	})

	t.Run("Search filter by name", func(t *testing.T) {
		db3, mock3, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db3)

		merkBukuID := uuid.New()

		merkBukuRows := sqlmock.NewRows([]string{"id", "code", "name", "description", "created_at", "updated_at"}).
			AddRow(merkBukuID, "JB001", "Textbook", nil, time.Now(), time.Now())

		mock3.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "merk_buku" WHERE merk_buku.code ILIKE $1 OR merk_buku.name ILIKE $2 OR merk_buku.description ILIKE $3 ORDER BY created_at ASC LIMIT 20`)).
			WithArgs("%Textbook%", "%Textbook%", "%Textbook%").
			WillReturnRows(merkBukuRows)

		countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
		mock3.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "merk_buku"`)).
			WillReturnRows(countRows)

		req := httptest.NewRequest("GET", "/merk-buku?search=Textbook", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["merk_buku"])
		assert.NotNil(t, response["pagination"])
	})
}

func TestGetMerkBuku(t *testing.T) {
	app := fiber.New()
	app.Get("/merk-buku/:id", handlers.GetMerkBuku)

	t.Run("Successfully get merk buku by ID", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		merkBukuID := uuid.New()

		merkBukuRows := sqlmock.NewRows([]string{"id", "code", "name", "description", "created_at", "updated_at"}).
			AddRow(merkBukuID, "JB001", "Textbook", "Test Description", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "merk_buku" WHERE id = $1`)).
			WithArgs(merkBukuID.String()).
			WillReturnRows(merkBukuRows)

		req := httptest.NewRequest("GET", "/merk-buku/"+merkBukuID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["merk_buku"])
	})

	t.Run("MerkBuku not found", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		merkBukuID := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "merk_buku" WHERE id = $1`)).
			WithArgs(merkBukuID.String()).
			WillReturnError(gorm.ErrRecordNotFound)

		req := httptest.NewRequest("GET", "/merk-buku/"+merkBukuID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "MerkBuku not found", response["error"])
	})
}

func TestCreateMerkBuku(t *testing.T) {
	app := fiber.New()
	app.Post("/merk-buku", handlers.CreateMerkBuku)

	t.Run("Successfully create merk buku", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		reqBody := models.MerkBuku{
			Name: "Textbook",
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "merk_buku"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
				AddRow(uuid.New(), time.Now(), time.Now()))
		mock.ExpectCommit()

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/merk-buku", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	})

	t.Run("Invalid request body", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/merk-buku", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Invalid request body", response["error"])
	})
}

func TestUpdateMerkBuku(t *testing.T) {
	app := fiber.New()
	app.Put("/merk-buku/:id", handlers.UpdateMerkBuku)

	t.Run("Successfully update merk buku", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		merkBukuID := uuid.New()

		merkBukuRows := sqlmock.NewRows([]string{"id", "code", "name", "description", "created_at", "updated_at"}).
			AddRow(merkBukuID, "JB001", "Textbook", "Test Description", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "merk_buku" WHERE id = $1`)).
			WithArgs(merkBukuID.String()).
			WillReturnRows(merkBukuRows)

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "merk_buku" SET .+ WHERE "id" = .+`).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		reqBody := models.MerkBuku{
			ID:   merkBukuID,
			Name: "Novel",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("PUT", "/merk-buku/"+merkBukuID.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("MerkBuku not found", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		merkBukuID := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "merk_buku" WHERE id = $1`)).
			WithArgs(merkBukuID.String()).
			WillReturnError(gorm.ErrRecordNotFound)

		reqBody := models.MerkBuku{
			Name: "Novel",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("PUT", "/merk-buku/"+merkBukuID.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "MerkBuku not found", response["error"])
	})
}

func TestDeleteMerkBuku(t *testing.T) {
	db, mock, err := testutil.SetupMockDB()
	assert.NoError(t, err)
	defer testutil.CloseMockDB(db)

	app := fiber.New()
	app.Delete("/merk-buku/:id", handlers.DeleteMerkBuku)

	t.Run("Successfully delete merk buku", func(t *testing.T) {
		merkBukuID := uuid.New()

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "merk_buku" WHERE id = $1`)).
			WithArgs(merkBukuID.String()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		req := httptest.NewRequest("DELETE", "/merk-buku/"+merkBukuID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "MerkBuku deleted successfully", response["message"])
	})

	t.Run("MerkBuku not found", func(t *testing.T) {
		merkBukuID := uuid.New()

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "merk_buku" WHERE id = $1`)).
			WithArgs(merkBukuID.String()).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		req := httptest.NewRequest("DELETE", "/merk-buku/"+merkBukuID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "MerkBuku not found", response["error"])
	})
}

