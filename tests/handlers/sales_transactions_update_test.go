package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"os"
	"path/filepath"
	"pustaka-backend/config"
	"pustaka-backend/handlers"
	"pustaka-backend/models"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// setupIntegrationTest connects to the real database for integration testing
func setupIntegrationTest(t *testing.T) {
	// Load .env file from project root
	// Try to find .env by walking up from the test directory
	envPath := findEnvFile()
	if envPath != "" {
		if err := godotenv.Load(envPath); err != nil {
			t.Logf("Warning: Could not load .env file: %v", err)
		}
	}

	// Build DSN from environment variables
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		getEnvOrDefault("DB_HOST", "localhost"),
		getEnvOrDefault("DB_PORT", "5432"),
		getEnvOrDefault("DB_USER", "postgres"),
		getEnvOrDefault("DB_PASSWORD", "postgres"),
		getEnvOrDefault("DB_NAME", "pustaka_db"),
		getEnvOrDefault("DB_SSLMODE", "disable"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skipf("Skipping integration test: unable to connect to database: %v", err)
	}

	// Set the global config.DB for handlers to use
	config.DB = db

	// Verify database has required data
	var count int64
	db.Model(&models.SalesAssociate{}).Count(&count)
	if count == 0 {
		t.Skip("Skipping integration test: no sales associates in database. Run migrations and seed data first.")
	}
}

// findEnvFile looks for .env file starting from current directory and walking up
func findEnvFile() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	for {
		envPath := filepath.Join(dir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			return envPath
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}

// getEnvOrDefault returns the environment variable value or a default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// cleanupTestData removes test data from database
func cleanupTestData(transactionID uuid.UUID) {
	config.DB.Exec("DELETE FROM sales_transaction_installments WHERE transaction_id = ?", transactionID)
	config.DB.Exec("DELETE FROM sales_transaction_items WHERE transaction_id = ?", transactionID)
	config.DB.Exec("DELETE FROM sales_transactions WHERE id = ?", transactionID)
}

func TestUpdateSalesTransaction_Scenario1_UpdateExpedition(t *testing.T) {
	setupIntegrationTest(t)

	app := fiber.New()
	app.Post("/sales-transactions", handlers.CreateSalesTransaction)
	app.Put("/sales-transactions/:id", handlers.UpdateSalesTransaction)
	app.Get("/sales-transactions/:id", handlers.GetSalesTransaction)

	// Fetch a sales associate from DB for testing
	var salesAssociate models.SalesAssociate
	if err := config.DB.First(&salesAssociate).Error; err != nil {
		t.Skipf("Skipping test: no sales associate found in database: %v", err)
	}

	// Fetch an expedition from DB for testing
	var expedition models.Expedition
	if err := config.DB.First(&expedition).Error; err != nil {
		t.Skipf("Skipping test: no expedition found in database: %v", err)
	}

	// Fetch a book from DB for testing
	var book models.Book
	if err := config.DB.First(&book).Error; err != nil {
		t.Skipf("Skipping test: no book found in database: %v", err)
	}

	// Step 1: Create transaction without expedition
	t.Run("Step 1: Create transaction without expedition", func(t *testing.T) {
		createReq := handlers.CreateTransactionRequest{
			SalesAssociateID: salesAssociate.ID.String(),
			PaymentType:      "T",
			TransactionDate:  time.Now(),
			ExpeditionPrice:  0,
			Items: []handlers.CreateTransactionItemRequest{
				{
					BookID:   book.ID.String(),
					Quantity: 1,
				},
			},
		}

		bodyBytes, _ := json.Marshal(createReq)
		req := httptest.NewRequest("POST", "/sales-transactions", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

		var createdTransaction models.SalesTransaction
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &createdTransaction)

		assert.NotEqual(t, uuid.Nil, createdTransaction.ID)

		// fmt.Printf("Transaction: %+v\n", createdTransaction)

		assert.Nil(t, createdTransaction.ExpeditionID)
		assert.Equal(t, float64(0), createdTransaction.ExpeditionPrice)

		// Store for cleanup
		defer cleanupTestData(createdTransaction.ID)

		// Step 2: Update with expedition
		t.Run("Step 2: Update with expedition", func(t *testing.T) {
			expeditionID := expedition.ID.String()
			expeditionPrice := float64(20000)

			updateReq := handlers.UpdateTransactionRequest{
				ExpeditionID:    &expeditionID,
				ExpeditionPrice: &expeditionPrice,
			}

			bodyBytes, _ := json.Marshal(updateReq)
			req := httptest.NewRequest("PUT", fmt.Sprintf("/sales-transactions/%s", createdTransaction.ID), bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			resp, _ := app.Test(req)

			assert.Equal(t, fiber.StatusOK, resp.StatusCode)

			var updatedTransaction models.SalesTransaction
			respBody, _ := io.ReadAll(resp.Body)
			json.Unmarshal(respBody, &updatedTransaction)

			assert.NotNil(t, updatedTransaction.ExpeditionID)
			assert.Equal(t, expedition.ID, *updatedTransaction.ExpeditionID)
			assert.Equal(t, float64(20000), updatedTransaction.ExpeditionPrice)

			// Verify total_amount includes expedition price
			expectedTotal := book.Price + 20000
			assert.Equal(t, expectedTotal, updatedTransaction.TotalAmount)

			// Step 3: Update to remove expedition (set to nil and 0)
			t.Run("Step 3: Remove expedition", func(t *testing.T) {
				emptyExpeditionID := ""
				zeroExpeditionPrice := float64(0)

				updateReq := handlers.UpdateTransactionRequest{
					ExpeditionID:    &emptyExpeditionID,
					ExpeditionPrice: &zeroExpeditionPrice,
				}

				bodyBytes, _ := json.Marshal(updateReq)
				req := httptest.NewRequest("PUT", fmt.Sprintf("/sales-transactions/%s", createdTransaction.ID), bytes.NewReader(bodyBytes))
				req.Header.Set("Content-Type", "application/json")
				resp, _ := app.Test(req)

				assert.Equal(t, fiber.StatusOK, resp.StatusCode)

				var finalTransaction models.SalesTransaction
				respBody, _ := io.ReadAll(resp.Body)
				json.Unmarshal(respBody, &finalTransaction)

				assert.Nil(t, finalTransaction.ExpeditionID)
				assert.Equal(t, float64(0), finalTransaction.ExpeditionPrice)

				// Verify total_amount no longer includes expedition price
				expectedTotal := book.Price
				assert.Equal(t, expectedTotal, finalTransaction.TotalAmount)
			})
		})
	})
}

func TestUpdateSalesTransaction_Scenario2_UpdateItems(t *testing.T) {
	setupIntegrationTest(t)

	app := fiber.New()
	app.Post("/sales-transactions", handlers.CreateSalesTransaction)
	app.Put("/sales-transactions/:id", handlers.UpdateSalesTransaction)
	app.Get("/sales-transactions/:id", handlers.GetSalesTransaction)

	// Fetch a sales associate from DB for testing
	var salesAssociate models.SalesAssociate
	if err := config.DB.First(&salesAssociate).Error; err != nil {
		t.Skipf("Skipping test: no sales associate found in database: %v", err)
	}

	// Fetch books from DB for testing (need at least 2)
	var books []models.Book
	if err := config.DB.Limit(2).Find(&books).Error; err != nil || len(books) < 2 {
		t.Skipf("Skipping test: need at least 2 books in database")
	}
	bookX := books[0]
	bookY := books[1]

	// Step 1: Create transaction with 1 item
	t.Run("Step 1: Create transaction with 1 item", func(t *testing.T) {
		createReq := handlers.CreateTransactionRequest{
			SalesAssociateID: salesAssociate.ID.String(),
			PaymentType:      "T",
			TransactionDate:  time.Now(),
			ExpeditionPrice:  0,
			Items: []handlers.CreateTransactionItemRequest{
				{
					BookID:   bookX.ID.String(),
					Quantity: 1,
				},
			},
		}

		bodyBytes, _ := json.Marshal(createReq)
		req := httptest.NewRequest("POST", "/sales-transactions", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

		var createdTransaction models.SalesTransaction
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &createdTransaction)

		assert.NotEqual(t, uuid.Nil, createdTransaction.ID)
		assert.Len(t, createdTransaction.Items, 1)
		assert.Equal(t, bookX.ID, createdTransaction.Items[0].BookID)
		assert.Equal(t, 1, createdTransaction.Items[0].Quantity)

		// Store for cleanup
		defer cleanupTestData(createdTransaction.ID)

		// Step 2: Update status only, items should remain the same
		t.Run("Step 2: Update status only", func(t *testing.T) {
			status := 2
			updateReq := handlers.UpdateTransactionRequest{
				Status: &status,
				Items: []handlers.CreateTransactionItemRequest{
					{
						BookID:   bookX.ID.String(),
						Quantity: 1,
					},
				},
			}

			bodyBytes, _ := json.Marshal(updateReq)
			req := httptest.NewRequest("PUT", fmt.Sprintf("/sales-transactions/%s", createdTransaction.ID), bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			resp, _ := app.Test(req)

			assert.Equal(t, fiber.StatusOK, resp.StatusCode)

			var updatedTransaction models.SalesTransaction
			respBody, _ := io.ReadAll(resp.Body)
			json.Unmarshal(respBody, &updatedTransaction)

			assert.Equal(t, 2, updatedTransaction.Status)
			assert.Len(t, updatedTransaction.Items, 1)
			assert.Equal(t, bookX.ID, updatedTransaction.Items[0].BookID)
			assert.Equal(t, 1, updatedTransaction.Items[0].Quantity)

			// Step 3: Update quantity of existing item
			t.Run("Step 3: Update quantity of existing item", func(t *testing.T) {
				status := 1
				updateReq := handlers.UpdateTransactionRequest{
					Status: &status,
					Items: []handlers.CreateTransactionItemRequest{
						{
							BookID:   bookX.ID.String(),
							Quantity: 2, // Changed from 1 to 2
						},
					},
				}

				bodyBytes, _ := json.Marshal(updateReq)
				req := httptest.NewRequest("PUT", fmt.Sprintf("/sales-transactions/%s", createdTransaction.ID), bytes.NewReader(bodyBytes))
				req.Header.Set("Content-Type", "application/json")
				resp, _ := app.Test(req)

				assert.Equal(t, fiber.StatusOK, resp.StatusCode)

				var updatedTransaction models.SalesTransaction
				respBody, _ := io.ReadAll(resp.Body)
				json.Unmarshal(respBody, &updatedTransaction)

				assert.Equal(t, 1, updatedTransaction.Status)
				assert.Len(t, updatedTransaction.Items, 1)
				assert.Equal(t, bookX.ID, updatedTransaction.Items[0].BookID)
				assert.Equal(t, 2, updatedTransaction.Items[0].Quantity)

				// Verify subtotal was recalculated
				expectedSubtotal := bookX.Price * 2
				assert.Equal(t, expectedSubtotal, updatedTransaction.Items[0].Subtotal)

				// Step 4: Add additional item
				t.Run("Step 4: Add additional item", func(t *testing.T) {
					status := 1
					updateReq := handlers.UpdateTransactionRequest{
						Status: &status,
						Items: []handlers.CreateTransactionItemRequest{
							{
								BookID:   bookX.ID.String(),
								Quantity: 2,
							},
							{
								BookID:   bookY.ID.String(),
								Quantity: 1, // New item
							},
						},
					}

					bodyBytes, _ := json.Marshal(updateReq)
					req := httptest.NewRequest("PUT", fmt.Sprintf("/sales-transactions/%s", createdTransaction.ID), bytes.NewReader(bodyBytes))
					req.Header.Set("Content-Type", "application/json")
					resp, _ := app.Test(req)

					assert.Equal(t, fiber.StatusOK, resp.StatusCode)

					var finalTransaction models.SalesTransaction
					respBody, _ := io.ReadAll(resp.Body)
					json.Unmarshal(respBody, &finalTransaction)

					assert.Equal(t, 1, finalTransaction.Status)
					assert.Len(t, finalTransaction.Items, 2)

					// Verify both items exist with correct quantities
					itemsMap := make(map[string]models.SalesTransactionItem)
					for _, item := range finalTransaction.Items {
						itemsMap[item.BookID.String()] = item
					}

					// Check book X
					assert.Contains(t, itemsMap, bookX.ID.String())
					assert.Equal(t, 2, itemsMap[bookX.ID.String()].Quantity)

					// Check book Y
					assert.Contains(t, itemsMap, bookY.ID.String())
					assert.Equal(t, 1, itemsMap[bookY.ID.String()].Quantity)

					// Verify total amount
					expectedTotal := (bookX.Price * 2) + (bookY.Price * 1)
					assert.Equal(t, expectedTotal, finalTransaction.TotalAmount)
				})
			})
		})
	})
}

func TestUpdateSalesTransaction_DuplicateBookIDValidation(t *testing.T) {
	setupIntegrationTest(t)

	app := fiber.New()
	app.Post("/sales-transactions", handlers.CreateSalesTransaction)
	app.Put("/sales-transactions/:id", handlers.UpdateSalesTransaction)

	// Fetch a sales associate and book from DB
	var salesAssociate models.SalesAssociate
	if err := config.DB.First(&salesAssociate).Error; err != nil {
		t.Skipf("Skipping test: no sales associate found in database: %v", err)
	}

	var book models.Book
	if err := config.DB.First(&book).Error; err != nil {
		t.Skipf("Skipping test: no book found in database: %v", err)
	}

	// Create a transaction
	createReq := handlers.CreateTransactionRequest{
		SalesAssociateID: salesAssociate.ID.String(),
		PaymentType:      "T",
		TransactionDate:  time.Now(),
		ExpeditionPrice:  0,
		Items: []handlers.CreateTransactionItemRequest{
			{
				BookID:   book.ID.String(),
				Quantity: 1,
			},
		},
	}

	bodyBytes, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/sales-transactions", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	var createdTransaction models.SalesTransaction
	respBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &createdTransaction)

	defer cleanupTestData(createdTransaction.ID)

	// Try to update with duplicate book_id in request
	t.Run("Should reject duplicate book_id in request", func(t *testing.T) {
		updateReq := handlers.UpdateTransactionRequest{
			Items: []handlers.CreateTransactionItemRequest{
				{
					BookID:   book.ID.String(),
					Quantity: 1,
				},
				{
					BookID:   book.ID.String(), // Duplicate!
					Quantity: 2,
				},
			},
		}

		bodyBytes, _ := json.Marshal(updateReq)
		req := httptest.NewRequest("PUT", fmt.Sprintf("/sales-transactions/%s", createdTransaction.ID), bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &response)

		assert.Contains(t, response["error"], "Duplicate book_id")
	})
}
