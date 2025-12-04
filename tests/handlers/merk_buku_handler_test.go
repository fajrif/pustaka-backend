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
		kodeMerk1 := "MB001"
		namaMerk1 := "Merk 1"
		kodeMerk2 := "MB002"
		namaMerk2 := "Merk 2"

		merkID1 := uuid.New()
		merkID2 := uuid.New()
		userID1 := uuid.New()
		userID2 := uuid.New()

		// Mock merk_buku query
		merkRows := sqlmock.NewRows([]string{"id", "kode_merk", "nama_merk", "bantuan_promosi", "user_id", "tstamp"}).
			AddRow(merkID1, kodeMerk1, namaMerk1, 1000, userID1, time.Now()).
			AddRow(merkID2, kodeMerk2, namaMerk2, 2000, userID2, time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "merk_buku" ORDER BY created_at DESC`)).
			WillReturnRows(merkRows)

		// Mock user preload - GORM uses IN query for multiple records
		userRows := sqlmock.NewRows([]string{"id", "email", "password_hash", "full_name", "role", "created_at", "updated_at"}).
			AddRow(userID1, "user1@example.com", "hashedpass", "User 1", "user", time.Now(), time.Now()).
			AddRow(userID2, "user2@example.com", "hashedpass", "User 2", "user", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" IN`)).
			WithArgs(userID1, userID2).
			WillReturnRows(userRows)

		req := httptest.NewRequest("GET", "/merk-buku", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["merk_buku"])
	})

	t.Run("Empty list", func(t *testing.T) {
		db2, mock2, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db2)

		merkRows := sqlmock.NewRows([]string{"id", "kode_merk", "nama_merk", "bantuan_promosi", "user_id", "tstamp"})

		mock2.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "merk_buku" ORDER BY created_at DESC`)).
			WillReturnRows(merkRows)

		req := httptest.NewRequest("GET", "/merk-buku", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["merk_buku"])
	})
}

func TestGetMerkBuku(t *testing.T) {
	app := fiber.New()
	app.Get("/merk-buku/:id", handlers.GetMerkBuku)

	t.Run("Successfully get merk buku by ID", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		kodeMerk := "MB001"
		namaMerk := "Test Merk"
		merkID := uuid.New()
		userID := uuid.New()

		merkRows := sqlmock.NewRows([]string{"id", "kode_merk", "nama_merk", "bantuan_promosi", "user_id", "tstamp"}).
			AddRow(merkID, kodeMerk, namaMerk, 1000, userID, time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "merk_buku" WHERE id = $1`)).
			WithArgs(merkID.String()).
			WillReturnRows(merkRows)

		userRows := sqlmock.NewRows([]string{"id", "email", "password_hash", "full_name", "role", "created_at", "updated_at"}).
			AddRow(userID, "user@example.com", "hashedpass", "Test User", "user", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1`)).
			WithArgs(userID).
			WillReturnRows(userRows)

		req := httptest.NewRequest("GET", "/merk-buku/"+merkID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["merk_buku"])
	})

	t.Run("Merk buku not found", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		merkID := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "merk_buku" WHERE id = $1`)).
			WithArgs(merkID.String()).
			WillReturnError(gorm.ErrRecordNotFound)

		req := httptest.NewRequest("GET", "/merk-buku/"+merkID.String(), nil)
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

	// Add middleware to set userID in locals
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"))
		return c.Next()
	})

	app.Post("/merk-buku", handlers.CreateMerkBuku)

	t.Run("Successfully create merk buku", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		kodeMerk := "MB001"
		namaMerk := "Test Merk"
		bantuanPromosi := 1000

		reqBody := models.MerkBuku{
			KodeMerk:       &kodeMerk,
			NamaMerk:       &namaMerk,
			BantuanPromosi: &bantuanPromosi,
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "merk_buku"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "tstamp"}).
				AddRow(uuid.New(), time.Now()))
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

	// Add middleware to set userID in locals
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"))
		return c.Next()
	})

	app.Put("/merk-buku/:id", handlers.UpdateMerkBuku)

	t.Run("Successfully update merk buku", func(t *testing.T) {
		db3, mock3, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db3)

		merkID := uuid.New()
		userID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
		kodeMerk := "MB001"
		namaMerk := "Test Merk"

		// Mock: Find existing merk buku
		merkRows := sqlmock.NewRows([]string{"id", "kode_merk", "nama_merk", "bantuan_promosi", "user_id", "tstamp"}).
			AddRow(merkID, kodeMerk, namaMerk, 1000, userID, time.Now())

		mock3.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "merk_buku" WHERE id = $1`)).
			WithArgs(merkID.String()).
			WillReturnRows(merkRows)

		// Mock: Update merk buku
		mock3.ExpectBegin()
		mock3.ExpectExec(`UPDATE "merk_buku" SET .+ WHERE "id" = .+`).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock3.ExpectCommit()

		updatedNamaMerk := "Updated Merk"
		reqBody := models.MerkBuku{
			ID:             merkID, // Include ID to preserve it after body parsing
			KodeMerk:       &kodeMerk,
			NamaMerk:       &updatedNamaMerk,
			BantuanPromosi: new(int),
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("PUT", "/merk-buku/"+merkID.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("Merk buku not found", func(t *testing.T) {
		db5, mock5, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db5)

		merkID := uuid.New()

		mock5.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "merk_buku" WHERE id = $1`)).
			WithArgs(merkID.String()).
			WillReturnError(gorm.ErrRecordNotFound)

		kodeMerk := "MB001"
		namaMerk := "Test Merk"
		reqBody := models.MerkBuku{
			KodeMerk: &kodeMerk,
			NamaMerk: &namaMerk,
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("PUT", "/merk-buku/"+merkID.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "MerkBuku not found", response["error"])
	})

	t.Run("Invalid request body", func(t *testing.T) {
		db4, mock4, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db4)

		merkID := uuid.New()
		userID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
		kodeMerk := "MB001"
		namaMerk := "Test Merk"

		// Mock: Find existing merk buku (this happens before body parsing)
		merkRows := sqlmock.NewRows([]string{"id", "kode_merk", "nama_merk", "bantuan_promosi", "user_id", "tstamp"}).
			AddRow(merkID, kodeMerk, namaMerk, 1000, userID, time.Now())

		mock4.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "merk_buku" WHERE id = $1`)).
			WithArgs(merkID.String()).
			WillReturnRows(merkRows)

		req := httptest.NewRequest("PUT", "/merk-buku/"+merkID.String(), bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})
}

func TestDeleteMerkBuku(t *testing.T) {
	db, mock, err := testutil.SetupMockDB()
	assert.NoError(t, err)
	defer testutil.CloseMockDB(db)

	app := fiber.New()
	app.Delete("/merk-buku/:id", handlers.DeleteMerkBuku)

	t.Run("Successfully delete merk buku", func(t *testing.T) {
		merkID := uuid.New()

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "merk_buku" WHERE id = $1`)).
			WithArgs(merkID.String()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		req := httptest.NewRequest("DELETE", "/merk-buku/"+merkID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "MerkBuku deleted successfully", response["message"])
	})

	t.Run("Merk buku not found", func(t *testing.T) {
		merkID := uuid.New()

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "merk_buku" WHERE id = $1`)).
			WithArgs(merkID.String()).
			WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected
		mock.ExpectCommit()

		req := httptest.NewRequest("DELETE", "/merk-buku/"+merkID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "MerkBuku not found", response["error"])
	})
}
