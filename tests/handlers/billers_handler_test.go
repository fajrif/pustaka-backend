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

func TestGetAllBillers(t *testing.T) {
	db, mock, err := testutil.SetupMockDB()
	assert.NoError(t, err)
	defer testutil.CloseMockDB(db)

	app := fiber.New()
	app.Get("/billers", handlers.GetAllBillers)

	t.Run("Successfully get all billers", func(t *testing.T) {
		billerID1 := uuid.New()
		billerID2 := uuid.New()
		cityID1 := uuid.New()
		cityID2 := uuid.New()

		billerRows := sqlmock.NewRows([]string{"id", "code", "name", "description", "npwp", "address", "city_id", "phone1", "phone2", "fax", "email", "website", "logo_url", "created_at", "updated_at"}).
			AddRow(billerID1, "BIL001", "Biller A", "Test Description", "12.345.678.9-012.345", "Jl. Test 1", cityID1, "021-111", "021-112", "021-113", "biller_a@example.com", "www.biller-a.com", nil, time.Now(), time.Now()).
			AddRow(billerID2, "BIL002", "Biller B", nil, "98.765.432.1-098.765", "Jl. Test 2", cityID2, "021-222", nil, nil, nil, nil, nil, time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "billers" ORDER BY created_at DESC`)).
			WillReturnRows(billerRows)

		cityRows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
			AddRow(cityID1, "Jakarta", time.Now(), time.Now()).
			AddRow(cityID2, "Bandung", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cities" WHERE "cities"."id" IN`)).
			WithArgs(cityID1, cityID2).
			WillReturnRows(cityRows)

		countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "billers"`)).
			WillReturnRows(countRows)

		req := httptest.NewRequest("GET", "/billers", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["billers"])
		assert.NotNil(t, response["pagination"])
	})

	t.Run("Search filter by code", func(t *testing.T) {
		db2, mock2, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db2)

		billerID := uuid.New()
		cityID := uuid.New()
		description := "Test Description"

		billerRows := sqlmock.NewRows([]string{"id", "code", "name", "description", "npwp", "address", "city_id", "phone1", "phone2", "fax", "email", "website", "logo_url", "created_at", "updated_at"}).
			AddRow(billerID, "BIL001", "Biller A", &description, "12.345.678.9-012.345", "Jl. Test", cityID, "021-111", nil, nil, nil, nil, nil, time.Now(), time.Now())

		mock2.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "billers" WHERE billers.code ILIKE $1 OR billers.name ILIKE $2 OR billers.description ILIKE $3 ORDER BY created_at DESC`)).
			WithArgs("%BIL001%", "%BIL001%", "%BIL001%").
			WillReturnRows(billerRows)

		cityRows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
			AddRow(cityID, "Jakarta", time.Now(), time.Now())

		mock2.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cities" WHERE "cities"."id" = $1`)).
			WithArgs(cityID).
			WillReturnRows(cityRows)

		countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
		mock2.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "billers"`)).
			WillReturnRows(countRows)

		req := httptest.NewRequest("GET", "/billers?search=BIL001", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["billers"])
		assert.NotNil(t, response["pagination"])
	})

	t.Run("Search filter by name", func(t *testing.T) {
		db3, mock3, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db3)

		billerID := uuid.New()
		cityID := uuid.New()

		billerRows := sqlmock.NewRows([]string{"id", "code", "name", "description", "npwp", "address", "city_id", "phone1", "phone2", "fax", "email", "website", "logo_url", "created_at", "updated_at"}).
			AddRow(billerID, "BIL001", "BillerA", nil, "12.345.678.9-012.345", "Jl. Test", cityID, "021-111", nil, nil, nil, nil, nil, time.Now(), time.Now())

		mock3.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "billers" WHERE billers.code ILIKE $1 OR billers.name ILIKE $2 OR billers.description ILIKE $3 ORDER BY created_at DESC`)).
			WithArgs("%BillerA%", "%BillerA%", "%BillerA%").
			WillReturnRows(billerRows)

		cityRows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
			AddRow(cityID, "Jakarta", time.Now(), time.Now())

		mock3.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cities" WHERE "cities"."id" = $1`)).
			WithArgs(cityID).
			WillReturnRows(cityRows)

		countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
		mock3.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "billers"`)).
			WillReturnRows(countRows)

		req := httptest.NewRequest("GET", "/billers?search=BillerA", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["billers"])
		assert.NotNil(t, response["pagination"])
	})

	t.Run("Search filter by description", func(t *testing.T) {
		db4, mock4, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db4)

		billerID := uuid.New()
		cityID := uuid.New()

		billerRows := sqlmock.NewRows([]string{"id", "code", "name", "description", "npwp", "address", "city_id", "phone1", "phone2", "fax", "email", "website", "logo_url", "created_at", "updated_at"}).
			AddRow(billerID, "BIL001", "BillerA", "Yang di Ampera", "12.345.678.9-012.345", "Jl. Test", cityID, "021-111", nil, nil, nil, nil, nil, time.Now(), time.Now())

		mock4.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "billers" WHERE billers.code ILIKE $1 OR billers.name ILIKE $2 OR billers.description ILIKE $3 ORDER BY created_at DESC`)).
			WithArgs("%ampera%", "%ampera%", "%ampera%").
			WillReturnRows(billerRows)

		cityRows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
			AddRow(cityID, "Jakarta", time.Now(), time.Now())

		mock4.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cities" WHERE "cities"."id" = $1`)).
			WithArgs(cityID).
			WillReturnRows(cityRows)

		countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
		mock4.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "billers"`)).
			WillReturnRows(countRows)

		req := httptest.NewRequest("GET", "/billers?search=ampera", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["billers"])
		assert.NotNil(t, response["pagination"])
	})
}

func TestGetBiller(t *testing.T) {
	app := fiber.New()
	app.Get("/billers/:id", handlers.GetBiller)

	t.Run("Successfully get biller by ID", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		billerID := uuid.New()
		cityID := uuid.New()

		billerRows := sqlmock.NewRows([]string{"id", "code", "name", "description", "npwp", "address", "city_id", "phone1", "phone2", "fax", "email", "website", "logo_url", "created_at", "updated_at"}).
			AddRow(billerID, "BIL001", "Biller A", "Test Description", "12.345.678.9-012.345", "Jl. Test", cityID, "021-111", nil, nil, nil, nil, nil, time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "billers" WHERE id = $1`)).
			WithArgs(billerID.String()).
			WillReturnRows(billerRows)

		cityRows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
			AddRow(cityID, "Jakarta", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cities" WHERE "cities"."id" = $1`)).
			WithArgs(cityID).
			WillReturnRows(cityRows)

		req := httptest.NewRequest("GET", "/billers/"+billerID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["biller"])
	})

	t.Run("Biller not found", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		billerID := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "billers" WHERE id = $1`)).
			WithArgs(billerID.String()).
			WillReturnError(gorm.ErrRecordNotFound)

		req := httptest.NewRequest("GET", "/billers/"+billerID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Biller not found", response["error"])
	})
}

func TestCreateBiller(t *testing.T) {
	app := fiber.New()
	app.Post("/billers", handlers.CreateBiller)

	t.Run("Successfully create biller", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		cityID := uuid.New()
		name := "Biller A"
		reqBody := models.Biller{
			Code:    "BIL001",
			Name:    &name,
			NPWP:    "12.345.678.9-012.345",
			Address: "Jl. Test",
			CityID:  &cityID,
			Phone1:  "021-111",
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "billers"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
				AddRow(uuid.New(), time.Now(), time.Now()))
		mock.ExpectCommit()

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/billers", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	})

	t.Run("Invalid request body", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/billers", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Invalid request body", response["error"])
	})
}

func TestUpdateBiller(t *testing.T) {
	app := fiber.New()
	app.Put("/billers/:id", handlers.UpdateBiller)

	t.Run("Successfully update biller", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		billerID := uuid.New()
		cityID := uuid.New()

		billerRows := sqlmock.NewRows([]string{"id", "code", "name", "description", "npwp", "address", "city_id", "phone1", "phone2", "fax", "email", "website", "logo_url", "created_at", "updated_at"}).
			AddRow(billerID, "BIL001", "Biller A", "Test Description", "12.345.678.9-012.345", "Jl. Test", cityID, "021-111", nil, nil, nil, nil, nil, time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "billers" WHERE id = $1`)).
			WithArgs(billerID.String()).
			WillReturnRows(billerRows)

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "billers" SET .+ WHERE "id" = .+`).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		name := "Biller A Updated"
		reqBody := models.Biller{
			ID:      billerID,
			Code:    "BIL001",
			Name:    &name,
			NPWP:    "12.345.678.9-012.345",
			Address: "Jl. Test Updated",
			Phone1:  "021-111",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("PUT", "/billers/"+billerID.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("Biller not found", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		billerID := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "billers" WHERE id = $1`)).
			WithArgs(billerID.String()).
			WillReturnError(gorm.ErrRecordNotFound)

		name := "Biller A"
		reqBody := models.Biller{
			Code:    "BIL001",
			Name:    &name,
			NPWP:    "12.345.678.9-012.345",
			Address: "Jl. Test",
			Phone1:  "021-111",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("PUT", "/billers/"+billerID.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Biller not found", response["error"])
	})
}

func TestDeleteBiller(t *testing.T) {
	db, mock, err := testutil.SetupMockDB()
	assert.NoError(t, err)
	defer testutil.CloseMockDB(db)

	app := fiber.New()
	app.Delete("/billers/:id", handlers.DeleteBiller)

	t.Run("Successfully delete biller", func(t *testing.T) {
		billerID := uuid.New()

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "billers" WHERE id = $1`)).
			WithArgs(billerID.String()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		req := httptest.NewRequest("DELETE", "/billers/"+billerID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Biller deleted successfully", response["message"])
	})

	t.Run("Biller not found", func(t *testing.T) {
		billerID := uuid.New()

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "billers" WHERE id = $1`)).
			WithArgs(billerID.String()).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		req := httptest.NewRequest("DELETE", "/billers/"+billerID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Biller not found", response["error"])
	})
}
