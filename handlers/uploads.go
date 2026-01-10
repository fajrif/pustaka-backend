package handlers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"pustaka-backend/config"
	"pustaka-backend/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// UploadResourceField godoc
// @Summary Upload file for a resource
// @Description Upload image or file for a specific resource and field
// @Tags Upload
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param resource path string true "Resource name (users, books, publishers, expeditions, sales-associates)"
// @Param field path string true "Field name (photo, image, file, logo)"
// @Param id path string true "Resource ID (UUID)"
// @Param file formData file true "File to upload"
// @Success 200 {object} map[string]interface{} "Upload successful with file URL"
// @Failure 400 {object} map[string]interface{} "Invalid request or validation error"
// @Failure 404 {object} map[string]interface{} "Resource not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/upload/{resource}/{field}/{id} [post]
func UploadResourceField(c *fiber.Ctx) error {
	resource := c.Params("resource")
	field := c.Params("field")
	id := c.Params("id")

	// Validate resource and field combination
	validCombinations := map[string][]string{
		"users":            {"photo"},
		"books":            {"image", "file"},
		"publishers":       {"logo", "file"},
		"expeditions":      {"logo", "file"},
		"sales-associates": {"photo", "file"},
		"billers":          {"logo"},
	}

	validFields, resourceExists := validCombinations[resource]
	if !resourceExists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid resource. Valid resources: users, books, publishers, expeditions, sales-associates, billers",
		})
	}

	fieldValid := false
	for _, vf := range validFields {
		if vf == field {
			fieldValid = true
			break
		}
	}

	if !fieldValid {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("Invalid field '%s' for resource '%s'. Valid fields: %s", field, resource, strings.Join(validFields, ", ")),
		})
	}

	// Retrieve uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No file uploaded",
		})
	}

	// Validate file type and size
	fileExt := strings.ToLower(filepath.Ext(file.Filename))
	var maxSize int64

	// Determine allowed extensions and max size based on field type
	isImageField := field == "photo" || field == "image" || field == "logo"
	isFileField := field == "file"

	if isImageField {
		// Image fields: JPEG, PNG, max 5MB
		if fileExt != ".jpg" && fileExt != ".jpeg" && fileExt != ".png" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid file type. Only JPEG and PNG images are allowed for this field",
			})
		}
		maxSize = 5 * 1024 * 1024 // 5MB
	} else if isFileField {
		// File fields: PDF, max 10MB
		if fileExt != ".pdf" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid file type. Only PDF files are allowed for this field",
			})
		}
		maxSize = 10 * 1024 * 1024 // 10MB
	}

	if file.Size > maxSize {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("File size exceeds maximum allowed size of %dMB", maxSize/(1024*1024)),
		})
	}

	// Create directory structure: uploads/:resource/:field/:id
	uploadDir := filepath.Join("uploads", resource, field, id)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create upload directory",
		})
	}

	// Generate unique filename to avoid conflicts
	timestamp := time.Now().Unix()
	uniqueFilename := fmt.Sprintf("%d_%s%s", timestamp, uuid.New().String()[:8], fileExt)
	filePath := filepath.Join(uploadDir, uniqueFilename)

	// Save the file
	if err := c.SaveFile(file, filePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save file",
		})
	}

	// Construct file URL (without protocol and host for flexibility)
	fileURL := "/" + strings.ReplaceAll(filePath, "\\", "/")

	// Update database based on resource type
	// This is transactional - if DB update fails, we delete the uploaded file
	err = updateResourceFileURL(resource, field, id, fileURL)
	if err != nil {
		// Remove uploaded file if database update fails
		os.Remove(filePath)

		if err.Error() == "resource not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": fmt.Sprintf("%s with ID %s not found", resource, id),
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update database",
		})
	}

	return c.JSON(fiber.Map{
		"message":  "File uploaded successfully",
		"file_url": fileURL,
		"resource": resource,
		"field":    field,
		"id":       id,
	})
}

// updateResourceFileURL updates the file URL in the database for the specified resource
func updateResourceFileURL(resource, field, id, fileURL string) error {
	// Map field names to database column names
	fieldColumnMap := map[string]string{
		"photo": "photo_url",
		"image": "image_url",
		"file":  "file_url",
		"logo":  "logo_url",
	}

	columnName := fieldColumnMap[field]

	var err error

	switch resource {
	case "users":
		var user models.User
		if err := config.DB.Where("id = ?", id).First(&user).Error; err != nil {
			return fmt.Errorf("resource not found")
		}
		err = config.DB.Model(&user).Update(columnName, fileURL).Error

	case "books":
		var book models.Book
		if err := config.DB.Where("id = ?", id).First(&book).Error; err != nil {
			return fmt.Errorf("resource not found")
		}
		err = config.DB.Model(&book).Update(columnName, fileURL).Error

	case "publishers":
		var publisher models.Publisher
		if err := config.DB.Where("id = ?", id).First(&publisher).Error; err != nil {
			return fmt.Errorf("resource not found")
		}
		err = config.DB.Model(&publisher).Update(columnName, fileURL).Error

	case "expeditions":
		var expedition models.Expedition
		if err := config.DB.Where("id = ?", id).First(&expedition).Error; err != nil {
			return fmt.Errorf("resource not found")
		}
		err = config.DB.Model(&expedition).Update(columnName, fileURL).Error

	case "sales-associates":
		var salesAssociate models.SalesAssociate
		if err := config.DB.Where("id = ?", id).First(&salesAssociate).Error; err != nil {
			return fmt.Errorf("resource not found")
		}
		err = config.DB.Model(&salesAssociate).Update(columnName, fileURL).Error

	case "billers":
		var biller models.Biller
		if err := config.DB.Where("id = ?", id).First(&biller).Error; err != nil {
			return fmt.Errorf("resource not found")
		}
		err = config.DB.Model(&biller).Update(columnName, fileURL).Error

	default:
		return fmt.Errorf("invalid resource")
	}

	return err
}

// DeleteResourceField godoc
// @Summary Delete file for a resource
// @Description Delete uploaded file and set database field to null
// @Tags Upload
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param resource path string true "Resource name (users, books, publishers, expeditions, sales-associates)"
// @Param field path string true "Field name (photo, image, file, logo)"
// @Param id path string true "Resource ID (UUID)"
// @Success 200 {object} map[string]interface{} "File deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request or validation error"
// @Failure 404 {object} map[string]interface{} "Resource or file not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/upload/{resource}/{field}/{id} [delete]
func DeleteResourceField(c *fiber.Ctx) error {
	resource := c.Params("resource")
	field := c.Params("field")
	id := c.Params("id")

	// Validate resource and field combination
	validCombinations := map[string][]string{
		"users":            {"photo"},
		"books":            {"image", "file"},
		"publishers":       {"logo", "file"},
		"expeditions":      {"logo", "file"},
		"sales-associates": {"photo", "file"},
		"billers":          {"logo"},
	}

	validFields, resourceExists := validCombinations[resource]
	if !resourceExists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid resource. Valid resources: users, books, publishers, expeditions, sales-associates, billers",
		})
	}

	fieldValid := false
	for _, vf := range validFields {
		if vf == field {
			fieldValid = true
			break
		}
	}

	if !fieldValid {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("Invalid field '%s' for resource '%s'. Valid fields: %s", field, resource, strings.Join(validFields, ", ")),
		})
	}

	// Get current file URL from database
	currentFileURL, err := getResourceFileURL(resource, field, id)
	if err != nil {
		if err.Error() == "resource not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": fmt.Sprintf("%s with ID %s not found", resource, id),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve resource",
		})
	}

	// Check if there's a file to delete
	if currentFileURL == nil || *currentFileURL == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No file found for this resource field",
		})
	}

	// Delete the physical file
	// Remove leading slash if present
	filePath := strings.TrimPrefix(*currentFileURL, "/")
	if err := os.Remove(filePath); err != nil {
		// Log the error but continue with database update
		// File might have been manually deleted
		fmt.Printf("Warning: Failed to delete file %s: %v\n", filePath, err)
	}

	// Update database field to NULL
	err = updateResourceFileURL(resource, field, id, "")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update database",
		})
	}

	return c.JSON(fiber.Map{
		"message":  "File deleted successfully",
		"resource": resource,
		"field":    field,
		"id":       id,
	})
}

// getResourceFileURL retrieves the current file URL from the database for the specified resource
func getResourceFileURL(resource, field, id string) (*string, error) {
	// Map field names to database column names
	fieldColumnMap := map[string]string{
		"photo": "photo_url",
		"image": "image_url",
		"file":  "file_url",
		"logo":  "logo_url",
	}

	columnName := fieldColumnMap[field]

	switch resource {
	case "users":
		var user models.User
		if err := config.DB.Where("id = ?", id).First(&user).Error; err != nil {
			return nil, fmt.Errorf("resource not found")
		}
		return user.PhotoUrl, nil

	case "books":
		var book models.Book
		if err := config.DB.Where("id = ?", id).First(&book).Error; err != nil {
			return nil, fmt.Errorf("resource not found")
		}
		if columnName == "image_url" {
			return book.ImageUrl, nil
		}
		return book.FileUrl, nil

	case "publishers":
		var publisher models.Publisher
		if err := config.DB.Where("id = ?", id).First(&publisher).Error; err != nil {
			return nil, fmt.Errorf("resource not found")
		}
		if columnName == "logo_url" {
			return publisher.LogoUrl, nil
		}
		return publisher.FileUrl, nil

	case "expeditions":
		var expedition models.Expedition
		if err := config.DB.Where("id = ?", id).First(&expedition).Error; err != nil {
			return nil, fmt.Errorf("resource not found")
		}
		if columnName == "logo_url" {
			return expedition.LogoUrl, nil
		}
		return expedition.FileUrl, nil

	case "sales-associates":
		var salesAssociate models.SalesAssociate
		if err := config.DB.Where("id = ?", id).First(&salesAssociate).Error; err != nil {
			return nil, fmt.Errorf("resource not found")
		}
		if columnName == "photo_url" {
			return salesAssociate.PhotoUrl, nil
		}
		return salesAssociate.FileUrl, nil

	case "billers":
		var biller models.Biller
		if err := config.DB.Where("id = ?", id).First(&biller).Error; err != nil {
			return nil, fmt.Errorf("resource not found")
		}
		return biller.LogoUrl, nil

	default:
		return nil, fmt.Errorf("invalid resource")
	}
}
