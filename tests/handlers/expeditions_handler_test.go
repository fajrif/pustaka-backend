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

func TestGetAllExpeditions(t *testing.T) {
	db, mock, err := testutil.SetupMockDB()
	assert.NoError(t, err)
	defer testutil.CloseMockDB(db)

	app := fiber.New()
	app.Get("/expeditions", handlers.GetAllExpeditions)

	t.Run("Successfully get all expeditions", func(t *testing.T) {
		expeditionID1 := uuid.New()
		expeditionID2 := uuid.New()
		cityID1 := uuid.New()
		cityID2 := uuid.New()

		expeditionRows := sqlmock.NewRows([]string{"id", "name", "address", "city_id", "area", "phone1", "phone2", "email", "website", "created_at", "updated_at"}).
			AddRow(expeditionID1, "JNE", "Jl. Test 1", cityID1, "Area 1", "021-111", "021-112", "jne@example.com", "www.jne.com", time.Now(), time.Now()).
			AddRow(expeditionID2, "TIKI", "Jl. Test 2", cityID2, "Area 2", "021-222", nil, nil, nil, time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "expeditions" ORDER BY created_at DESC`)).
			WillReturnRows(expeditionRows)

		cityRows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
			AddRow(cityID1, "Jakarta", time.Now(), time.Now()).
			AddRow(cityID2, "Bandung", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cities" WHERE "cities"."id" IN`)).
			WithArgs(cityID1, cityID2).
			WillReturnRows(cityRows)

		req := httptest.NewRequest("GET", "/expeditions", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["expeditions"])
	})
}

func TestGetExpedition(t *testing.T) {
	app := fiber.New()
	app.Get("/expeditions/:id", handlers.GetExpedition)

	t.Run("Successfully get expedition by ID", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		expeditionID := uuid.New()
		cityID := uuid.New()

		expeditionRows := sqlmock.NewRows([]string{"id", "name", "address", "city_id", "area", "phone1", "phone2", "email", "website", "created_at", "updated_at"}).
			AddRow(expeditionID, "JNE", "Jl. Test", cityID, "Area 1", "021-111", nil, nil, nil, time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "expeditions" WHERE id = $1`)).
			WithArgs(expeditionID.String()).
			WillReturnRows(expeditionRows)

		cityRows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
			AddRow(cityID, "Jakarta", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cities" WHERE "cities"."id" = $1`)).
			WithArgs(cityID).
			WillReturnRows(cityRows)

		req := httptest.NewRequest("GET", "/expeditions/"+expeditionID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["expedition"])
	})

	t.Run("Expedition not found", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		expeditionID := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "expeditions" WHERE id = $1`)).
			WithArgs(expeditionID.String()).
			WillReturnError(gorm.ErrRecordNotFound)

		req := httptest.NewRequest("GET", "/expeditions/"+expeditionID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Expedition not found", response["error"])
	})
}

func TestCreateExpedition(t *testing.T) {
	app := fiber.New()
	app.Post("/expeditions", handlers.CreateExpedition)

	t.Run("Successfully create expedition", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		cityID := uuid.New()
		reqBody := models.Expedition{
			Name:    "JNE",
			Address: "Jl. Test",
			CityID:  &cityID,
			Phone1:  "021-111",
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "expeditions"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
				AddRow(uuid.New(), time.Now(), time.Now()))
		mock.ExpectCommit()

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/expeditions", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	})

	t.Run("Invalid request body", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/expeditions", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Invalid request body", response["error"])
	})
}

func TestUpdateExpedition(t *testing.T) {
	app := fiber.New()
	app.Put("/expeditions/:id", handlers.UpdateExpedition)

	t.Run("Successfully update expedition", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		expeditionID := uuid.New()
		cityID := uuid.New()

		expeditionRows := sqlmock.NewRows([]string{"id", "name", "address", "city_id", "area", "phone1", "phone2", "email", "website", "created_at", "updated_at"}).
			AddRow(expeditionID, "JNE", "Jl. Test", cityID, "Area 1", "021-111", nil, nil, nil, time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "expeditions" WHERE id = $1`)).
			WithArgs(expeditionID.String()).
			WillReturnRows(expeditionRows)

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "expeditions" SET .+ WHERE "id" = .+`).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		reqBody := models.Expedition{
			ID:      expeditionID,
			Name:    "JNE Updated",
			Address: "Jl. Test Updated",
			Phone1:  "021-111",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("PUT", "/expeditions/"+expeditionID.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("Expedition not found", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		expeditionID := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "expeditions" WHERE id = $1`)).
			WithArgs(expeditionID.String()).
			WillReturnError(gorm.ErrRecordNotFound)

		reqBody := models.Expedition{
			Name:    "JNE",
			Address: "Jl. Test",
			Phone1:  "021-111",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("PUT", "/expeditions/"+expeditionID.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Expedition not found", response["error"])
	})
}

func TestDeleteExpedition(t *testing.T) {
	db, mock, err := testutil.SetupMockDB()
	assert.NoError(t, err)
	defer testutil.CloseMockDB(db)

	app := fiber.New()
	app.Delete("/expeditions/:id", handlers.DeleteExpedition)

	t.Run("Successfully delete expedition", func(t *testing.T) {
		expeditionID := uuid.New()

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "expeditions" WHERE id = $1`)).
			WithArgs(expeditionID.String()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		req := httptest.NewRequest("DELETE", "/expeditions/"+expeditionID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Expedition deleted successfully", response["message"])
	})

	t.Run("Expedition not found", func(t *testing.T) {
		expeditionID := uuid.New()

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "expeditions" WHERE id = $1`)).
			WithArgs(expeditionID.String()).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		req := httptest.NewRequest("DELETE", "/expeditions/"+expeditionID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Expedition not found", response["error"])
	})
}
