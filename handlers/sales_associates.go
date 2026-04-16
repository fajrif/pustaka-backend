package handlers

import (
	"github.com/gofiber/fiber/v2"
	"pustaka-backend/config"
	"pustaka-backend/helpers"
	"pustaka-backend/models"
)

type CreateSalesAssociateRequest struct {
	Code            string  `json:"code"`
	Name            string  `json:"name"`
	NoKtp           *string `json:"no_ktp"`
	Description     *string `json:"description"`
	Address         string  `json:"address"`
	CityID          *string `json:"city_id"`
	Area            *string `json:"area"`
	Phone1          string  `json:"phone1"`
	Phone2          *string `json:"phone2"`
	Email           *string `json:"email"`
	Website         *string `json:"website"`
	JenisPembayaran *string `json:"jenis_pembayaran"`
	JoinDate        *string `json:"join_date"`
	EndJoinDate     *string `json:"end_join_date"`
	Discount        float64 `json:"discount"`
	PhotoUrl        *string `json:"photo_url"`
	FileUrl         *string `json:"file_url"`
}

type UpdateSalesAssociateRequest struct {
	Code            *string  `json:"code"`
	Name            *string  `json:"name"`
	NoKtp           *string  `json:"no_ktp"`
	Description     *string  `json:"description"`
	Address         *string  `json:"address"`
	CityID          *string  `json:"city_id"`
	Area            *string  `json:"area"`
	Phone1          *string  `json:"phone1"`
	Phone2          *string  `json:"phone2"`
	Email           *string  `json:"email"`
	Website         *string  `json:"website"`
	JenisPembayaran *string  `json:"jenis_pembayaran"`
	JoinDate        *string  `json:"join_date"`
	EndJoinDate     *string  `json:"end_join_date"`
	Discount        *float64 `json:"discount"`
	PhotoUrl        *string  `json:"photo_url"`
	FileUrl         *string  `json:"file_url"`
}

// GetAllSalesAssociates godoc
// @Summary Get all sales associates
// @Description Retrieve all sales associates with their related city information
// @Tags SalesAssociates
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param search query string false "Search by code, name, or description"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items per page (default: 20)"
// @Success 200 {object} map[string]interface{} "List of all sales associates with pagination"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/sales-associates [get]
func GetAllSalesAssociates(c *fiber.Ctx) error {
	var salesAssociates []models.SalesAssociate

	// Get pagination parameters
	pagination := helpers.GetPaginationParams(c)

	query := config.DB.Order("created_at DESC")
	queryCount := config.DB.Model(&models.SalesAssociate{})

	// add params for not using pagination
	if c.Query("all") == "true" {
		pagination.Limit = -1 // No limit
		pagination.Offset = 0 // No offset
	}

	// Filter search
	if searchQuery := c.Query("search"); searchQuery != "" {
		// Wrap string search with wildcard SQL LIKE
		searchTerm := "%" + searchQuery + "%"
		cond := "sales_associates.code ILIKE ? OR sales_associates.name ILIKE ? OR sales_associates.description ILIKE ?"
		args := []interface{}{searchTerm, searchTerm, searchTerm}

		query = query.Where(cond, args...)
		queryCount = queryCount.Where(cond, args...)
	}

	// Apply pagination and fetch data
	if err := query.Offset(pagination.Offset).Limit(pagination.Limit).Preload("City").Find(&salesAssociates).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch all sales associates",
		})
	}

	// Create pagination response
	response, err := helpers.CreatePaginationResponse(queryCount, salesAssociates, "sales_associates", pagination.Page, pagination.Limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create pagination response",
		})
	}

	return c.JSON(response)
}

// GetSalesAssociate godoc
// @Summary Get a sales associate by ID
// @Description Retrieve a single sales associate by its ID with related city information
// @Tags SalesAssociates
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "SalesAssociate ID (UUID)"
// @Success 200 {object} map[string]interface{} "SalesAssociate details"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "SalesAssociate not found"
// @Router /api/sales-associates/{id} [get]
func GetSalesAssociate(c *fiber.Ctx) error {
	id := c.Params("id")

	var salesAssociate models.SalesAssociate
	if err := config.DB.Preload("City").Where("id = ?", id).First(&salesAssociate).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "SalesAssociate not found",
		})
	}

	return c.JSON(fiber.Map{
		"sales_associate": salesAssociate,
	})
}

// CreateSalesAssociate godoc
// @Summary Create a new sales associate
// @Description Create a new sales associate entry
// @Tags SalesAssociates
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateSalesAssociateRequest true "SalesAssociate details"
// @Success 201 {object} models.SalesAssociate "Created sales associate"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/sales-associates [post]
func CreateSalesAssociate(c *fiber.Ctx) error {
	var req CreateSalesAssociateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	joinDate, err := helpers.ParseDateString(req.JoinDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid join_date format. Use YYYY-MM-DD",
		})
	}
	if joinDate == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "join_date is required",
		})
	}

	endJoinDate, _ := helpers.ParseDateString(req.EndJoinDate)

	cityID := helpers.ParseUUIDPtr(req.CityID)

	salesAssociate := models.SalesAssociate{
		Code:            req.Code,
		Name:            req.Name,
		NoKtp:           req.NoKtp,
		Description:     req.Description,
		Address:         req.Address,
		CityID:          cityID,
		Area:            req.Area,
		Phone1:          req.Phone1,
		Phone2:          req.Phone2,
		Email:           req.Email,
		Website:         req.Website,
		JenisPembayaran: "T",
		JoinDate:        *joinDate,
		EndJoinDate:     endJoinDate,
		Discount:        req.Discount,
		PhotoUrl:        req.PhotoUrl,
		FileUrl:         req.FileUrl,
	}

	if req.JenisPembayaran != nil {
		salesAssociate.JenisPembayaran = *req.JenisPembayaran
	}

	if err := config.DB.Create(&salesAssociate).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create sales associate",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(salesAssociate)
}

// UpdateSalesAssociate godoc
// @Summary Update a sales associate
// @Description Update an existing sales associate by ID
// @Tags SalesAssociates
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "SalesAssociate ID (UUID)"
// @Param request body UpdateSalesAssociateRequest true "Updated sales associate details"
// @Success 200 {object} models.SalesAssociate "Updated sales associate"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "SalesAssociate not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/sales-associates/{id} [put]
func UpdateSalesAssociate(c *fiber.Ctx) error {
	id := c.Params("id")

	var salesAssociate models.SalesAssociate
	if err := config.DB.Where("id = ?", id).First(&salesAssociate).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "SalesAssociate not found",
		})
	}

	var req UpdateSalesAssociateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	updates := make(map[string]interface{})

	if req.Code != nil {
		updates["code"] = *req.Code
	}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.NoKtp != nil {
		updates["no_ktp"] = *req.NoKtp
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Address != nil {
		updates["address"] = *req.Address
	}
	if req.CityID != nil {
		cityUUID := helpers.ParseUUIDPtr(req.CityID)
		if cityUUID != nil {
			updates["city_id"] = *cityUUID
		}
	}
	if req.Area != nil {
		updates["area"] = *req.Area
	}
	if req.Phone1 != nil {
		updates["phone1"] = *req.Phone1
	}
	if req.Phone2 != nil {
		updates["phone2"] = *req.Phone2
	}
	if req.Email != nil {
		updates["email"] = *req.Email
	}
	if req.Website != nil {
		updates["website"] = *req.Website
	}
	if req.JenisPembayaran != nil {
		updates["jenis_pembayaran"] = *req.JenisPembayaran
	}
	if req.JoinDate != nil {
		joinDate, err := helpers.ParseDateString(req.JoinDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid join_date format. Use YYYY-MM-DD",
			})
		}
		if joinDate != nil {
			updates["join_date"] = *joinDate
		}
	}
	if req.EndJoinDate != nil {
		endJoinDate, _ := helpers.ParseDateString(req.EndJoinDate)
		updates["end_join_date"] = endJoinDate
	}
	if req.Discount != nil {
		updates["discount"] = *req.Discount
	}
	if req.PhotoUrl != nil {
		updates["photo_url"] = *req.PhotoUrl
	}
	if req.FileUrl != nil {
		updates["file_url"] = *req.FileUrl
	}

	if err := config.DB.Model(&salesAssociate).Updates(updates).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update sales associate",
		})
	}

	config.DB.Where("id = ?", id).First(&salesAssociate)
	return c.JSON(salesAssociate)
}

// DeleteSalesAssociate godoc
// @Summary Delete a sales associate
// @Description Delete a sales associate by ID
// @Tags SalesAssociates
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "SalesAssociate ID (UUID)"
// @Success 200 {object} map[string]interface{} "SalesAssociate deleted successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "SalesAssociate not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/sales-associates/{id} [delete]
func DeleteSalesAssociate(c *fiber.Ctx) error {
	id := c.Params("id")

	result := config.DB.Delete(&models.SalesAssociate{}, "id = ?", id)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete sales associate",
		})
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "SalesAssociate not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "SalesAssociate deleted successfully",
	})
}
