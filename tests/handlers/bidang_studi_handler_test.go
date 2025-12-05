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

func TestGetAllBidangStudi(t *testing.T) {
	db, mock, err := testutil.SetupMockDB()
	assert.NoError(t, err)
	defer testutil.CloseMockDB(db)

	app := fiber.New()
	app.Get("/bidang-studi", handlers.GetAllBidangStudi)

	t.Run("Successfully get all bidang studi", func(t *testing.T) {
		bidangStudiID1 := uuid.New()
		bidangStudiID2 := uuid.New()

		rows := sqlmock.NewRows([]string{"id", "code", "name", "description", "created_at", "updated_at"}).
			AddRow(bidangStudiID1, "BS001", "Mathematics", "Test Description", time.Now(), time.Now()).
			AddRow(bidangStudiID2, "BS002", "Science", "Test Description", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "bidang_studi" ORDER BY created_at DESC`)).
			WillReturnRows(rows)

		countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "bidang_studi"`)).
			WillReturnRows(countRows)

		req := httptest.NewRequest("GET", "/bidang-studi", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["bidang_studi"])
		assert.NotNil(t, response["pagination"])
	})

	t.Run("Search filter by code", func(t *testing.T) {
		db2, mock2, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db2)

		bidangStudiID := uuid.New()
		description := "Test Description"

		bidangStudiRows := sqlmock.NewRows([]string{"id", "code", "name", "description", "created_at", "updated_at"}).
			AddRow(bidangStudiID, "BS001", "Science", &description, time.Now(), time.Now())

		mock2.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "bidang_studi" WHERE bidang_studi.code ILIKE $1 OR bidang_studi.name ILIKE $2 OR bidang_studi.description ILIKE $3 ORDER BY created_at DESC`)).
			WithArgs("%BS001%", "%BS001%", "%BS001%").
			WillReturnRows(bidangStudiRows)

		countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
		mock2.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "bidang_studi"`)).
			WillReturnRows(countRows)

		req := httptest.NewRequest("GET", "/bidang-studi?search=BS001", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["bidang_studi"])
		assert.NotNil(t, response["pagination"])
	})

	t.Run("Search filter by name", func(t *testing.T) {
		db3, mock3, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db3)

		bidangStudiID := uuid.New()

		bidangStudiRows := sqlmock.NewRows([]string{"id", "code", "name", "description", "created_at", "updated_at"}).
			AddRow(bidangStudiID, "BS001", "Science", nil, time.Now(), time.Now())

		mock3.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "bidang_studi" WHERE bidang_studi.code ILIKE $1 OR bidang_studi.name ILIKE $2 OR bidang_studi.description ILIKE $3 ORDER BY created_at DESC`)).
			WithArgs("%Science%", "%Science%", "%Science%").
			WillReturnRows(bidangStudiRows)

		countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
		mock3.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "bidang_studi"`)).
			WillReturnRows(countRows)

		req := httptest.NewRequest("GET", "/bidang-studi?search=Science", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["bidang_studi"])
		assert.NotNil(t, response["pagination"])
	})

	t.Run("Search filter by description", func(t *testing.T) {
		db4, mock4, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db4)

		bidangStudiID := uuid.New()

		bidangStudiRows := sqlmock.NewRows([]string{"id", "code", "name", "description", "created_at", "updated_at"}).
			AddRow(bidangStudiID, "BS001", "Science", "meta-science", time.Now(), time.Now())

		mock4.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "bidang_studi" WHERE bidang_studi.code ILIKE $1 OR bidang_studi.name ILIKE $2 OR bidang_studi.description ILIKE $3 ORDER BY created_at DESC`)).
			WithArgs("%meta-science%", "%meta-science%", "%meta-science%").
			WillReturnRows(bidangStudiRows)

		countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
		mock4.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "bidang_studi"`)).
			WillReturnRows(countRows)

		req := httptest.NewRequest("GET", "/bidang-studi?search=meta-science", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["bidang_studi"])
		assert.NotNil(t, response["pagination"])
	})
}

func TestGetBidangStudi(t *testing.T) {
	app := fiber.New()
	app.Get("/bidang-studi/:id", handlers.GetBidangStudi)

	t.Run("Successfully get bidang studi by ID", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		bidangStudiID := uuid.New()

		rows := sqlmock.NewRows([]string{"id", "code", "name", "description", "created_at", "updated_at"}).
			AddRow(bidangStudiID, "BS001", "Mathematics", "Test Description", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "bidang_studi" WHERE id = $1`)).
			WithArgs(bidangStudiID.String()).
			WillReturnRows(rows)

		req := httptest.NewRequest("GET", "/bidang-studi/"+bidangStudiID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["bidang_studi"])
	})

	t.Run("BidangStudi not found", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		bidangStudiID := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "bidang_studi" WHERE id = $1`)).
			WithArgs(bidangStudiID.String()).
			WillReturnError(gorm.ErrRecordNotFound)

		req := httptest.NewRequest("GET", "/bidang-studi/"+bidangStudiID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "BidangStudi not found", response["error"])
	})
}

func TestCreateBidangStudi(t *testing.T) {
	app := fiber.New()
	app.Post("/bidang-studi", handlers.CreateBidangStudi)

	t.Run("Successfully create bidang studi", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		reqBody := models.BidangStudi{
			Name: "Mathematics",
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "bidang_studi"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
				AddRow(uuid.New(), time.Now(), time.Now()))
		mock.ExpectCommit()

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/bidang-studi", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	})

	t.Run("Invalid request body", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/bidang-studi", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Invalid request body", response["error"])
	})
}

func TestUpdateBidangStudi(t *testing.T) {
	app := fiber.New()
	app.Put("/bidang-studi/:id", handlers.UpdateBidangStudi)

	t.Run("Successfully update bidang studi", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		bidangStudiID := uuid.New()

		rows := sqlmock.NewRows([]string{"id", "code", "name", "description", "created_at", "updated_at"}).
			AddRow(bidangStudiID, "BS001", "Mathematics", "Test Description", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "bidang_studi" WHERE id = $1`)).
			WithArgs(bidangStudiID.String()).
			WillReturnRows(rows)

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "bidang_studi" SET .+ WHERE "id" = .+`).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		reqBody := models.BidangStudi{
			ID:   bidangStudiID,
			Name: "Science",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("PUT", "/bidang-studi/"+bidangStudiID.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("BidangStudi not found", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		bidangStudiID := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "bidang_studi" WHERE id = $1`)).
			WithArgs(bidangStudiID.String()).
			WillReturnError(gorm.ErrRecordNotFound)

		reqBody := models.BidangStudi{
			Name: "Science",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("PUT", "/bidang-studi/"+bidangStudiID.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "BidangStudi not found", response["error"])
	})
}

func TestDeleteBidangStudi(t *testing.T) {
	db, mock, err := testutil.SetupMockDB()
	assert.NoError(t, err)
	defer testutil.CloseMockDB(db)

	app := fiber.New()
	app.Delete("/bidang-studi/:id", handlers.DeleteBidangStudi)

	t.Run("Successfully delete bidang studi", func(t *testing.T) {
		bidangStudiID := uuid.New()

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "bidang_studi" WHERE id = $1`)).
			WithArgs(bidangStudiID.String()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		req := httptest.NewRequest("DELETE", "/bidang-studi/"+bidangStudiID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "BidangStudi deleted successfully", response["message"])
	})

	t.Run("BidangStudi not found", func(t *testing.T) {
		bidangStudiID := uuid.New()

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "bidang_studi" WHERE id = $1`)).
			WithArgs(bidangStudiID.String()).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		req := httptest.NewRequest("DELETE", "/bidang-studi/"+bidangStudiID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "BidangStudi not found", response["error"])
	})
}
