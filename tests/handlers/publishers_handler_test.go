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

func TestGetAllPublishers(t *testing.T) {
	db, mock, err := testutil.SetupMockDB()
	assert.NoError(t, err)
	defer testutil.CloseMockDB(db)

	app := fiber.New()
	app.Get("/publishers", handlers.GetAllPublishers)

	t.Run("Successfully get all publishers", func(t *testing.T) {
		publisherID1 := uuid.New()
		publisherID2 := uuid.New()
		cityID1 := uuid.New()
		cityID2 := uuid.New()

		publisherRows := sqlmock.NewRows([]string{"id", "code", "name", "description", "address", "city_id", "area", "phone1", "phone2", "email", "website", "created_at", "updated_at"}).
			AddRow(publisherID1, "PUB001", "Gramedia", "Test Description", "Jl. Test 1", cityID1, "Area 1", "021-111", "021-112", "gramedia@example.com", "www.gramedia.com", time.Now(), time.Now()).
			AddRow(publisherID2, "PUB002", "Erlangga", nil, "Jl. Test 2", cityID2, "Area 2", "021-222", nil, nil, nil, time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "publishers" ORDER BY created_at DESC`)).
			WillReturnRows(publisherRows)

		cityRows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
			AddRow(cityID1, "Jakarta", time.Now(), time.Now()).
			AddRow(cityID2, "Bandung", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cities" WHERE "cities"."id" IN`)).
			WithArgs(cityID1, cityID2).
			WillReturnRows(cityRows)

		req := httptest.NewRequest("GET", "/publishers", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["publishers"])
	})

	t.Run("Search filter by code", func(t *testing.T) {
		db2, mock2, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db2)

		publisherID := uuid.New()
		cityID := uuid.New()
		description := "Test Description"

		publisherRows := sqlmock.NewRows([]string{"id", "code", "name", "description", "address", "city_id", "area", "phone1", "phone2", "email", "website", "created_at", "updated_at"}).
			AddRow(publisherID, "PUB001", "Publisher A", &description, "Jl. Test", cityID, "Area 1", "021-111", nil, nil, nil, time.Now(), time.Now())

		mock2.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "publishers" WHERE publishers.code ILIKE $1 OR publishers.name ILIKE $2 OR publishers.description ILIKE $3 ORDER BY created_at DESC`)).
			WithArgs("%PUB001%", "%PUB001%", "%PUB001%").
			WillReturnRows(publisherRows)

		cityRows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
			AddRow(cityID, "Jakarta", time.Now(), time.Now())

		mock2.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cities" WHERE "cities"."id" = $1`)).
			WithArgs(cityID).
			WillReturnRows(cityRows)

		req := httptest.NewRequest("GET", "/publishers?search=PUB001", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["publishers"])
	})

	t.Run("Search filter by name", func(t *testing.T) {
		db3, mock3, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db3)

		publisherID := uuid.New()
		cityID := uuid.New()

		publisherRows := sqlmock.NewRows([]string{"id", "code", "name", "description", "address", "city_id", "area", "phone1", "phone2", "email", "website", "created_at", "updated_at"}).
			AddRow(publisherID, "PUB001", "PublisherA", nil, "Jl. Test", cityID, "Area 1", "021-111", nil, nil, nil, time.Now(), time.Now())

		mock3.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "publishers" WHERE publishers.code ILIKE $1 OR publishers.name ILIKE $2 OR publishers.description ILIKE $3 ORDER BY created_at DESC`)).
			WithArgs("%PublisherA%", "%PublisherA%", "%PublisherA%").
			WillReturnRows(publisherRows)

		cityRows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
			AddRow(cityID, "Jakarta", time.Now(), time.Now())

		mock3.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cities" WHERE "cities"."id" = $1`)).
			WithArgs(cityID).
			WillReturnRows(cityRows)

		req := httptest.NewRequest("GET", "/publishers?search=PublisherA", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["publishers"])
	})
}

func TestGetPublisher(t *testing.T) {
	app := fiber.New()
	app.Get("/publishers/:id", handlers.GetPublisher)

	t.Run("Successfully get publisher by ID", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		publisherID := uuid.New()
		cityID := uuid.New()

		publisherRows := sqlmock.NewRows([]string{"id", "code", "name", "description", "address", "city_id", "area", "phone1", "phone2", "email", "website", "created_at", "updated_at"}).
			AddRow(publisherID, "PUB001", "Gramedia", "Test Description", "Jl. Test", cityID, "Area 1", "021-111", nil, nil, nil, time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "publishers" WHERE id = $1`)).
			WithArgs(publisherID.String()).
			WillReturnRows(publisherRows)

		cityRows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
			AddRow(cityID, "Jakarta", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cities" WHERE "cities"."id" = $1`)).
			WithArgs(cityID).
			WillReturnRows(cityRows)

		req := httptest.NewRequest("GET", "/publishers/"+publisherID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["publisher"])
	})

	t.Run("Publisher not found", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		publisherID := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "publishers" WHERE id = $1`)).
			WithArgs(publisherID.String()).
			WillReturnError(gorm.ErrRecordNotFound)

		req := httptest.NewRequest("GET", "/publishers/"+publisherID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Publisher not found", response["error"])
	})
}

func TestCreatePublisher(t *testing.T) {
	app := fiber.New()
	app.Post("/publishers", handlers.CreatePublisher)

	t.Run("Successfully create publisher", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		cityID := uuid.New()
		reqBody := models.Publisher{
			Name:    "Gramedia",
			Address: "Jl. Test",
			CityID:  &cityID,
			Phone1:  "021-111",
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "publishers"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
				AddRow(uuid.New(), time.Now(), time.Now()))
		mock.ExpectCommit()

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/publishers", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	})

	t.Run("Invalid request body", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/publishers", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Invalid request body", response["error"])
	})
}

func TestUpdatePublisher(t *testing.T) {
	app := fiber.New()
	app.Put("/publishers/:id", handlers.UpdatePublisher)

	t.Run("Successfully update publisher", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		publisherID := uuid.New()
		cityID := uuid.New()

		publisherRows := sqlmock.NewRows([]string{"id", "code", "name", "description", "address", "city_id", "area", "phone1", "phone2", "email", "website", "created_at", "updated_at"}).
			AddRow(publisherID, "PUB001", "Gramedia", "Test Description", "Jl. Test", cityID, "Area 1", "021-111", nil, nil, nil, time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "publishers" WHERE id = $1`)).
			WithArgs(publisherID.String()).
			WillReturnRows(publisherRows)

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "publishers" SET .+ WHERE "id" = .+`).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		reqBody := models.Publisher{
			ID:      publisherID,
			Name:    "Gramedia Updated",
			Address: "Jl. Test Updated",
			Phone1:  "021-111",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("PUT", "/publishers/"+publisherID.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("Publisher not found", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		publisherID := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "publishers" WHERE id = $1`)).
			WithArgs(publisherID.String()).
			WillReturnError(gorm.ErrRecordNotFound)

		reqBody := models.Publisher{
			Name:    "Gramedia",
			Address: "Jl. Test",
			Phone1:  "021-111",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("PUT", "/publishers/"+publisherID.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Publisher not found", response["error"])
	})
}

func TestDeletePublisher(t *testing.T) {
	db, mock, err := testutil.SetupMockDB()
	assert.NoError(t, err)
	defer testutil.CloseMockDB(db)

	app := fiber.New()
	app.Delete("/publishers/:id", handlers.DeletePublisher)

	t.Run("Successfully delete publisher", func(t *testing.T) {
		publisherID := uuid.New()

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "publishers" WHERE id = $1`)).
			WithArgs(publisherID.String()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		req := httptest.NewRequest("DELETE", "/publishers/"+publisherID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Publisher deleted successfully", response["message"])
	})

	t.Run("Publisher not found", func(t *testing.T) {
		publisherID := uuid.New()

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "publishers" WHERE id = $1`)).
			WithArgs(publisherID.String()).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		req := httptest.NewRequest("DELETE", "/publishers/"+publisherID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Publisher not found", response["error"])
	})
}
