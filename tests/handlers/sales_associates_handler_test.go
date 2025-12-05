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

func TestGetAllSalesAssociates(t *testing.T) {
	db, mock, err := testutil.SetupMockDB()
	assert.NoError(t, err)
	defer testutil.CloseMockDB(db)

	app := fiber.New()
	app.Get("/sales-associates", handlers.GetAllSalesAssociates)

	t.Run("Successfully get all sales associates", func(t *testing.T) {
		salesAssociateID1 := uuid.New()
		salesAssociateID2 := uuid.New()
		cityID1 := uuid.New()
		cityID2 := uuid.New()
		joinDate := time.Now()
		discount := 10.5

		salesAssociateRows := sqlmock.NewRows([]string{"id", "code", "name", "description", "address", "city_id", "area", "phone1", "phone2", "email", "website", "jenis_pembayaran", "join_date", "end_join_date", "discount", "created_at", "updated_at"}).
			AddRow(salesAssociateID1, "SA001", "PT Distributor A", "Test Description", "Jl. Test 1", cityID1, "Area 1", "021-111", "021-112", "distributora@example.com", "www.distributora.com", "T", joinDate, nil, discount, time.Now(), time.Now()).
			AddRow(salesAssociateID2, "SA002", "PT Distributor B", nil, "Jl. Test 2", cityID2, "Area 2", "021-222", nil, nil, nil, "T", joinDate, nil, discount, time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "sales_associates" ORDER BY created_at DESC`)).
			WillReturnRows(salesAssociateRows)

		cityRows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
			AddRow(cityID1, "Jakarta", time.Now(), time.Now()).
			AddRow(cityID2, "Bandung", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cities" WHERE "cities"."id" IN`)).
			WithArgs(cityID1, cityID2).
			WillReturnRows(cityRows)

		countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "sales_associates"`)).
			WillReturnRows(countRows)

		req := httptest.NewRequest("GET", "/sales-associates", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["sales_associates"])
		assert.NotNil(t, response["pagination"])
	})

	t.Run("Search filter by code", func(t *testing.T) {
		db2, mock2, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db2)

		salesAssociateID := uuid.New()
		cityID := uuid.New()
		description := "Test Description"
		joinDate := time.Now()

		salesAssociateRows := sqlmock.NewRows([]string{"id", "code", "name", "description", "address", "city_id", "area", "phone1", "phone2", "email", "website", "jenis_pembayaran", "join_date", "end_join_date", "discount", "created_at", "updated_at"}).
			AddRow(salesAssociateID, "SA001", "Sales A", &description, "Jl. Test", cityID, "Area 1", "021-111", nil, nil, nil, "T", joinDate, nil, 10.0, time.Now(), time.Now())

		mock2.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "sales_associates" WHERE sales_associates.code ILIKE $1 OR sales_associates.name ILIKE $2 OR sales_associates.description ILIKE $3 ORDER BY created_at DESC`)).
			WithArgs("%SA001%", "%SA001%", "%SA001%").
			WillReturnRows(salesAssociateRows)

		cityRows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
			AddRow(cityID, "Jakarta", time.Now(), time.Now())

		mock2.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cities" WHERE "cities"."id" = $1`)).
			WithArgs(cityID).
			WillReturnRows(cityRows)

		countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
		mock2.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "sales_associates"`)).
			WillReturnRows(countRows)

		req := httptest.NewRequest("GET", "/sales-associates?search=SA001", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["sales_associates"])
		assert.NotNil(t, response["pagination"])
	})

	t.Run("Search filter by name", func(t *testing.T) {
		db3, mock3, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db3)

		salesAssociateID := uuid.New()
		cityID := uuid.New()
		joinDate := time.Now()

		salesAssociateRows := sqlmock.NewRows([]string{"id", "code", "name", "description", "address", "city_id", "area", "phone1", "phone2", "email", "website", "jenis_pembayaran", "join_date", "end_join_date", "discount", "created_at", "updated_at"}).
			AddRow(salesAssociateID, "SA001", "SalesA", nil, "Jl. Test", cityID, "Area 1", "021-111", nil, nil, nil, "T", joinDate, nil, 10.0, time.Now(), time.Now())

		mock3.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "sales_associates" WHERE sales_associates.code ILIKE $1 OR sales_associates.name ILIKE $2 OR sales_associates.description ILIKE $3 ORDER BY created_at DESC`)).
			WithArgs("%SalesA%", "%SalesA%", "%SalesA%").
			WillReturnRows(salesAssociateRows)

		cityRows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
			AddRow(cityID, "Jakarta", time.Now(), time.Now())

		mock3.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cities" WHERE "cities"."id" = $1`)).
			WithArgs(cityID).
			WillReturnRows(cityRows)

		countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
		mock3.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "sales_associates"`)).
			WillReturnRows(countRows)

		req := httptest.NewRequest("GET", "/sales-associates?search=SalesA", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["sales_associates"])
		assert.NotNil(t, response["pagination"])
	})
}

func TestGetSalesAssociate(t *testing.T) {
	app := fiber.New()
	app.Get("/sales-associates/:id", handlers.GetSalesAssociate)

	t.Run("Successfully get sales associate by ID", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		salesAssociateID := uuid.New()
		cityID := uuid.New()
		joinDate := time.Now()
		discount := 10.5

		salesAssociateRows := sqlmock.NewRows([]string{"id", "code", "name", "description", "address", "city_id", "area", "phone1", "phone2", "email", "website", "jenis_pembayaran", "join_date", "end_join_date", "discount", "created_at", "updated_at"}).
			AddRow(salesAssociateID, "SA001", "PT Distributor A", "Test Description", "Jl. Test", cityID, "Area 1", "021-111", nil, nil, nil, "T", joinDate, nil, discount, time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "sales_associates" WHERE id = $1`)).
			WithArgs(salesAssociateID.String()).
			WillReturnRows(salesAssociateRows)

		cityRows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
			AddRow(cityID, "Jakarta", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cities" WHERE "cities"."id" = $1`)).
			WithArgs(cityID).
			WillReturnRows(cityRows)

		req := httptest.NewRequest("GET", "/sales-associates/"+salesAssociateID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotNil(t, response["sales_associate"])
	})

	t.Run("Sales associate not found", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		salesAssociateID := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "sales_associates" WHERE id = $1`)).
			WithArgs(salesAssociateID.String()).
			WillReturnError(gorm.ErrRecordNotFound)

		req := httptest.NewRequest("GET", "/sales-associates/"+salesAssociateID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "SalesAssociate not found", response["error"])
	})
}

func TestCreateSalesAssociate(t *testing.T) {
	app := fiber.New()
	app.Post("/sales-associates", handlers.CreateSalesAssociate)

	t.Run("Successfully create sales associate", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		cityID := uuid.New()
		joinDate := time.Now()
		discount := 10.5
		jenisPembayaran := "T"

		reqBody := models.SalesAssociate{
			Name:            "PT Distributor A",
			Address:         "Jl. Test",
			CityID:          &cityID,
			Phone1:          "021-111",
			JenisPembayaran: jenisPembayaran,
			JoinDate:        joinDate,
			Discount:        discount,
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "sales_associates"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
				AddRow(uuid.New(), time.Now(), time.Now()))
		mock.ExpectCommit()

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/sales-associates", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	})

	t.Run("Invalid request body", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/sales-associates", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Invalid request body", response["error"])
	})
}

func TestUpdateSalesAssociate(t *testing.T) {
	app := fiber.New()
	app.Put("/sales-associates/:id", handlers.UpdateSalesAssociate)

	t.Run("Successfully update sales associate", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		salesAssociateID := uuid.New()
		cityID := uuid.New()
		joinDate := time.Now()
		discount := 10.5
		jenisPembayaran := "T"

		salesAssociateRows := sqlmock.NewRows([]string{"id", "code", "name", "description", "address", "city_id", "area", "phone1", "phone2", "email", "website", "jenis_pembayaran", "join_date", "end_join_date", "discount", "created_at", "updated_at"}).
			AddRow(salesAssociateID, "SA001", "PT Distributor A", "Test Description", "Jl. Test", cityID, "Area 1", "021-111", nil, nil, nil, "T", joinDate, nil, discount, time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "sales_associates" WHERE id = $1`)).
			WithArgs(salesAssociateID.String()).
			WillReturnRows(salesAssociateRows)

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "sales_associates" SET .+ WHERE "id" = .+`).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		reqBody := models.SalesAssociate{
			ID:              salesAssociateID,
			Name:            "PT Distributor A Updated",
			Address:         "Jl. Test Updated",
			Phone1:          "021-111",
			JenisPembayaran: jenisPembayaran,
			JoinDate:        joinDate,
			Discount:        discount,
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("PUT", "/sales-associates/"+salesAssociateID.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("Sales associate not found", func(t *testing.T) {
		db, mock, err := testutil.SetupMockDB()
		assert.NoError(t, err)
		defer testutil.CloseMockDB(db)

		salesAssociateID := uuid.New()
		joinDate := time.Now()
		discount := 10.5
		jenisPembayaran := "T"

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "sales_associates" WHERE id = $1`)).
			WithArgs(salesAssociateID.String()).
			WillReturnError(gorm.ErrRecordNotFound)

		reqBody := models.SalesAssociate{
			Name:            "PT Distributor A",
			Address:         "Jl. Test",
			Phone1:          "021-111",
			JenisPembayaran: jenisPembayaran,
			JoinDate:        joinDate,
			Discount:        discount,
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("PUT", "/sales-associates/"+salesAssociateID.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "SalesAssociate not found", response["error"])
	})
}

func TestDeleteSalesAssociate(t *testing.T) {
	db, mock, err := testutil.SetupMockDB()
	assert.NoError(t, err)
	defer testutil.CloseMockDB(db)

	app := fiber.New()
	app.Delete("/sales-associates/:id", handlers.DeleteSalesAssociate)

	t.Run("Successfully delete sales associate", func(t *testing.T) {
		salesAssociateID := uuid.New()

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "sales_associates" WHERE id = $1`)).
			WithArgs(salesAssociateID.String()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		req := httptest.NewRequest("DELETE", "/sales-associates/"+salesAssociateID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "SalesAssociate deleted successfully", response["message"])
	})

	t.Run("Sales associate not found", func(t *testing.T) {
		salesAssociateID := uuid.New()

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "sales_associates" WHERE id = $1`)).
			WithArgs(salesAssociateID.String()).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		req := httptest.NewRequest("DELETE", "/sales-associates/"+salesAssociateID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "SalesAssociate not found", response["error"])
	})
}
