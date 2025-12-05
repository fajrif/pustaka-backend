package handlers

import (
	"fmt"
	"pustaka-backend/config"
	"pustaka-backend/models"
	// "pustaka-backend/helpers"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetAllMerkBuku godoc
// @Summary Get all merk buku
// @Description Retrieve all book brands (merk buku) ordered by creation date
// @Tags MerkBuku
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "List of all merk buku"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/merk-buku [get]
func GetAllMerkBuku(c *fiber.Ctx) error {
    var all_merk_buku []models.MerkBuku
    query := config.DB.Order("created_at DESC")

    if err := query.
							Preload("User").
							Find(&all_merk_buku).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to fetch all merk buku",
        })
    }

		// Always return consistent JSON shape
    return c.JSON(fiber.Map{
        "merk_buku": all_merk_buku,
    })
}

// GetMerkBuku godoc
// @Summary Get a merk buku by ID
// @Description Retrieve a single book brand by its ID
// @Tags MerkBuku
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "MerkBuku ID (UUID)"
// @Success 200 {object} map[string]interface{} "MerkBuku details"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "MerkBuku not found"
// @Router /api/merk-buku/{id} [get]
func GetMerkBuku(c *fiber.Ctx) error {
    id := c.Params("id")
    query := config.DB

    var merk_buku models.MerkBuku
    if err := query.
							Preload("User").
							Where("id = ?", id).First(&merk_buku).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "MerkBuku not found",
        })
    }

    return c.JSON(fiber.Map{
        "merk_buku": merk_buku,
    })
}

// CreateMerkBuku godoc
// @Summary Create a new merk buku
// @Description Create a new book brand
// @Tags MerkBuku
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.MerkBuku true "MerkBuku details"
// @Success 201 {object} models.MerkBuku "Created merk buku"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/merk-buku [post]
func CreateMerkBuku(c *fiber.Ctx) error {
    var merk_buku models.MerkBuku
    if err := c.BodyParser(&merk_buku); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }

		fmt.Printf("MerkBuku object: %+v\n", merk_buku)

    // if err := helpers.ValidateMerkBukuData(&merk_buku); err != nil {
    //     return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
    //         "error": err.Error(),
    //     })
    // }

    userID := c.Locals("userID").(uuid.UUID)
    merk_buku.UserID = userID

    if err := config.DB.Create(&merk_buku).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to create merk buku",
        })
    }

    return c.Status(fiber.StatusCreated).JSON(merk_buku)
}

// UpdateMerkBuku godoc
// @Summary Update a merk buku
// @Description Update an existing book brand by ID
// @Tags MerkBuku
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "MerkBuku ID (UUID)"
// @Param request body models.MerkBuku true "Updated merk buku details"
// @Success 200 {object} models.MerkBuku "Updated merk buku"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "MerkBuku not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/merk-buku/{id} [put]
func UpdateMerkBuku(c *fiber.Ctx) error {
    id := c.Params("id")

    var merk_buku models.MerkBuku
    if err := config.DB.Where("id = ?", id).First(&merk_buku).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "MerkBuku not found",
        })
    }

    if err := c.BodyParser(&merk_buku); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }

    userID := c.Locals("userID").(uuid.UUID)
    merk_buku.UserID = userID

		fmt.Printf("MerkBuku object: %+v\n", merk_buku)

    // if err := helpers.ValidateMerkBukuData(&merk_buku); err != nil {
    //     return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
    //         "error": err.Error(),
    //     })
    // }

    if err := config.DB.Model(&merk_buku).Updates(merk_buku).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to update merk buku",
        })
    }

    return c.JSON(merk_buku)
}

// DeleteMerkBuku godoc
// @Summary Delete a merk buku
// @Description Delete a book brand by ID
// @Tags MerkBuku
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "MerkBuku ID (UUID)"
// @Success 200 {object} map[string]interface{} "MerkBuku deleted successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "MerkBuku not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/merk-buku/{id} [delete]
func DeleteMerkBuku(c *fiber.Ctx) error {
    id := c.Params("id")

    result := config.DB.Delete(&models.MerkBuku{}, "id = ?", id)
    if result.Error != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to delete merk buku",
        })
    }

    if result.RowsAffected == 0 {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "MerkBuku not found",
        })
    }

    return c.JSON(fiber.Map{
        "message": "MerkBuku deleted successfully",
    })
}
