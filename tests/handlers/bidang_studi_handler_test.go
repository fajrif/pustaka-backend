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

		rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
			AddRow(bidangStudiID1, "Mathematics", time.Now(), time.Now()).
			AddRow(bidangStudiID2, "Science", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "bidang_studi" ORDER BY created_at DESC`)).
			WillReturnRows(rows)

		req := httptest.NewRequest("GET", "/bidang-studi", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["bidang_studi"])
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

		rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
			AddRow(bidangStudiID, "Mathematics", time.Now(), time.Now())

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

		rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
			AddRow(bidangStudiID, "Mathematics", time.Now(), time.Now())

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
