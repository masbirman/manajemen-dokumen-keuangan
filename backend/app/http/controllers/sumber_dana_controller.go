package controllers

import (
	"strconv"

	"dokumen-keuangan/app/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// SumberDanaController handles sumber dana endpoints
type SumberDanaController struct {
	service *services.SumberDanaService
}

// NewSumberDanaController creates a new SumberDanaController instance
func NewSumberDanaController() *SumberDanaController {
	return &SumberDanaController{
		service: services.NewSumberDanaService(),
	}
}

// GetAll retrieves all sumber dana with pagination
// GET /api/sumber-dana
func (c *SumberDanaController) GetAll(ctx *fiber.Ctx) error {
	// Parse query parameters
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.Query("page_size", "10"))
	search := ctx.Query("search", "")

	result, err := c.service.GetAll(page, pageSize, search)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve sumber dana",
		})
	}

	return ctx.JSON(fiber.Map{
		"data":        result.Data,
		"total":       result.Total,
		"page":        result.Page,
		"page_size":   result.PageSize,
		"total_pages": result.TotalPages,
	})
}

// GetAllActive retrieves all active sumber dana (for dropdowns)
// GET /api/sumber-dana/active
func (c *SumberDanaController) GetAllActive(ctx *fiber.Ctx) error {
	sumberDanas, err := c.service.GetAllActive()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve sumber dana",
		})
	}

	return ctx.JSON(fiber.Map{
		"data": sumberDanas,
	})
}


// GetByID retrieves a sumber dana by ID
// GET /api/sumber-dana/:id
func (c *SumberDanaController) GetByID(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid sumber dana ID",
		})
	}

	sumberDana, err := c.service.GetByID(id)
	if err != nil {
		if err == services.ErrSumberDanaNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Sumber dana not found",
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve sumber dana",
		})
	}

	return ctx.JSON(fiber.Map{
		"data": sumberDana,
	})
}

// Create creates a new sumber dana
// POST /api/sumber-dana
func (c *SumberDanaController) Create(ctx *fiber.Ctx) error {
	var input services.CreateSumberDanaInput
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate input
	if validationErrors := input.Validate(); len(validationErrors) > 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": validationErrors,
		})
	}

	sumberDana, err := c.service.Create(&input)
	if err != nil {
		switch err {
		case services.ErrSumberDanaKodeExists:
			return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Sumber dana with this kode already exists",
			})
		case services.ErrSumberDanaKodeRequired:
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"errors": fiber.Map{"kode": "Kode is required"},
			})
		case services.ErrSumberDanaNamaRequired:
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"errors": fiber.Map{"nama": "Nama is required"},
			})
		default:
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create sumber dana",
			})
		}
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Sumber dana created successfully",
		"data":    sumberDana,
	})
}


// Update updates a sumber dana
// PUT /api/sumber-dana/:id
func (c *SumberDanaController) Update(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid sumber dana ID",
		})
	}

	var input services.UpdateSumberDanaInput
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	sumberDana, err := c.service.Update(id, &input)
	if err != nil {
		switch err {
		case services.ErrSumberDanaNotFound:
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Sumber dana not found",
			})
		case services.ErrSumberDanaKodeExists:
			return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Sumber dana with this kode already exists",
			})
		case services.ErrSumberDanaKodeRequired:
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"errors": fiber.Map{"kode": "Kode is required"},
			})
		case services.ErrSumberDanaNamaRequired:
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"errors": fiber.Map{"nama": "Nama is required"},
			})
		default:
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update sumber dana",
			})
		}
	}

	return ctx.JSON(fiber.Map{
		"message": "Sumber dana updated successfully",
		"data":    sumberDana,
	})
}

// Delete deletes a sumber dana
// DELETE /api/sumber-dana/:id
func (c *SumberDanaController) Delete(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid sumber dana ID",
		})
	}

	err = c.service.Delete(id)
	if err != nil {
		if err == services.ErrSumberDanaNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Sumber dana not found",
			})
		}
		// Check for referential integrity error
		if refErr, ok := err.(*services.ErrSumberDanaReferenced); ok {
			return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": refErr.Error(),
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete sumber dana",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "Sumber dana deleted successfully",
	})
}
