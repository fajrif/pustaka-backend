package handlers_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"os"
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
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Request structs for testing
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func TestRegister(t *testing.T) {
	db, mock, err := testutil.SetupMockDB()
	assert.NoError(t, err)
	defer testutil.CloseMockDB(db)

	app := fiber.New()
	app.Post("/register", handlers.Register)

	t.Run("Successful registration", func(t *testing.T) {
		reqBody := RegisterRequest{
			Email:    "newuser@example.com",
			Password: "password123",
			FullName: "New User",
		}

		// Mock: Check if user exists (should return error/not found)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1`)).
			WithArgs("newuser@example.com").
			WillReturnError(gorm.ErrRecordNotFound)

		// Mock: Create user
		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email", "full_name", "role", "created_at", "updated_at"}).
				AddRow(uuid.New(), "newuser@example.com", "New User", "user", time.Now(), time.Now()))
		mock.ExpectCommit()

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "User registered successfully", response["message"])
		assert.NotNil(t, response["user"])
	})

	t.Run("Invalid request body", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/register", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Invalid request body", response["error"])
	})

	t.Run("Email already exists", func(t *testing.T) {
		reqBody := RegisterRequest{
			Email:    "existing@example.com",
			Password: "password123",
			FullName: "Existing User",
		}

		// Mock: User already exists
		existingUser := models.User{
			ID:    uuid.New(),
			Email: "existing@example.com",
		}

		rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "full_name", "role", "created_at", "updated_at"}).
			AddRow(existingUser.ID, existingUser.Email, "hashedpass", "Existing User", "user", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1`)).
			WithArgs("existing@example.com").
			WillReturnRows(rows)

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusConflict, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Email already exists", response["error"])
	})
}

func TestLogin(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	os.Setenv("JWT_EXPIRE_HOURS", "24")
	defer os.Unsetenv("JWT_SECRET")
	defer os.Unsetenv("JWT_EXPIRE_HOURS")

	db, mock, err := testutil.SetupMockDB()
	assert.NoError(t, err)
	defer testutil.CloseMockDB(db)

	app := fiber.New()
	app.Post("/login", handlers.Login)

	t.Run("Successful login", func(t *testing.T) {
		// Create a hashed password
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

		userID := uuid.New()
		rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "full_name", "role", "created_at", "updated_at"}).
			AddRow(userID, "user@example.com", string(hashedPassword), "Test User", "user", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1`)).
			WithArgs("user@example.com").
			WillReturnRows(rows)

		reqBody := LoginRequest{
			Email:    "user@example.com",
			Password: "password123",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.NotEmpty(t, response["token"])
		assert.NotNil(t, response["user"])
	})

	t.Run("Invalid request body", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/login", bytes.NewReader([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("User not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1`)).
			WithArgs("nonexistent@example.com").
			WillReturnError(gorm.ErrRecordNotFound)

		reqBody := LoginRequest{
			Email:    "nonexistent@example.com",
			Password: "password123",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Invalid email or password", response["error"])
	})

	t.Run("Invalid password", func(t *testing.T) {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)

		userID := uuid.New()
		rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "full_name", "role", "created_at", "updated_at"}).
			AddRow(userID, "user@example.com", string(hashedPassword), "Test User", "user", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1`)).
			WithArgs("user@example.com").
			WillReturnRows(rows)

		reqBody := LoginRequest{
			Email:    "user@example.com",
			Password: "wrongpassword",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "Invalid email or password", response["error"])
	})
}

func TestGetMe(t *testing.T) {
	db, mock, err := testutil.SetupMockDB()
	assert.NoError(t, err)
	defer testutil.CloseMockDB(db)

	app := fiber.New()

	// Add middleware to set userID in locals
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", "550e8400-e29b-41d4-a716-446655440000")
		return c.Next()
	})

	app.Get("/me", handlers.GetMe)

	t.Run("Successfully get user profile", func(t *testing.T) {
		userID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
		rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "full_name", "role", "created_at", "updated_at"}).
			AddRow(userID, "user@example.com", "hashedpass", "Test User", "user", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1`)).
			WithArgs("550e8400-e29b-41d4-a716-446655440000").
			WillReturnRows(rows)

		req := httptest.NewRequest("GET", "/me", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "user@example.com", response["email"])
		assert.Equal(t, "Test User", response["full_name"])
		assert.Equal(t, "user", response["role"])
	})

	t.Run("User not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1`)).
			WithArgs("550e8400-e29b-41d4-a716-446655440000").
			WillReturnError(gorm.ErrRecordNotFound)

		req := httptest.NewRequest("GET", "/me", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Equal(t, "User not found", response["error"])
	})
}

func TestUpdateMe(t *testing.T) {
	db, mock, err := testutil.SetupMockDB()
	assert.NoError(t, err)
	defer testutil.CloseMockDB(db)

	app := fiber.New()

	// Add middleware to set userID in locals
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", "550e8400-e29b-41d4-a716-446655440000")
		return c.Next()
	})

	app.Put("/me", handlers.UpdateMe)

	t.Run("Successfully update user profile", func(t *testing.T) {
		userID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")

		// Mock: Find user
		rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "full_name", "role", "created_at", "updated_at"}).
			AddRow(userID, "user@example.com", "hashedpass", "Old Name", "user", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1`)).
			WithArgs("550e8400-e29b-41d4-a716-446655440000").
			WillReturnRows(rows)

		// Mock: Update user
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		updateData := map[string]interface{}{
			"full_name": "Updated Name",
		}

		body, _ := json.Marshal(updateData)
		req := httptest.NewRequest("PUT", "/me", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("Invalid request body", func(t *testing.T) {
		userID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")

		// Mock: Find user (this happens before body parsing in the handler)
		rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "full_name", "role", "created_at", "updated_at"}).
			AddRow(userID, "user@example.com", "hashedpass", "Test User", "user", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1`)).
			WithArgs("550e8400-e29b-41d4-a716-446655440000").
			WillReturnRows(rows)

		req := httptest.NewRequest("PUT", "/me", bytes.NewReader([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("User not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1`)).
			WithArgs("550e8400-e29b-41d4-a716-446655440000").
			WillReturnError(gorm.ErrRecordNotFound)

		updateData := map[string]interface{}{
			"full_name": "Updated Name",
		}

		body, _ := json.Marshal(updateData)
		req := httptest.NewRequest("PUT", "/me", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})
}
