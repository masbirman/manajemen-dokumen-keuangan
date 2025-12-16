package controllers

import (
	"strconv"

	"dokumen-keuangan/app/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// JenisDokumenController handles jenis dokumen endpoints
type JenisDokumenController struct {
	service *services.JenisDokumenService
}

// NewJenisDokumenController creates a new JenisDokumenController instance
func NewJenisDokumenController() *JenisDokumenController {
	return &JenisDokumenController{
		service: services.NewJenisDokumenService(),
	}
}

// GetAll retrieves all jenis dokumen with pagination
// GET /api/jenis-dokumen
func (c *JenisDokumenController) GetAll(ctx *fiber.Ctx) error {
	// Parse query parameters
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.Query("page_size", "10"))
	search := ctx.Query("search", "")

	result, err := c.service.GetAll(page, pageSize, search)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve jenis dokumen",
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

// GetAllActive retrieves all active jenis dokumen (for dropdowns)
// GET /api/jenis-dokumen/active
func (c *JenisDokumenController) GetAllActive(ctx *fiber.Ctx) error {
	jenisDokumens, err := c.service.GetAllActive()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve jenis dokumen",
		})
	}

	return ctx.JSON(fiber.Map{
		"data": jenisDokumens,
	})
}


// GetByID retrieves a jenis dokumen by ID
// GET /api/jenis-dokumen/:id
func (c *JenisDokumenController) GetByID(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid jenis dokumen ID",
		})
	}

	jenisDokumen, err := c.service.GetByID(id)
	if err != nil {
		if err == services.ErrJenisDokumenNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Jenis dokumen not found",
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve jenis dokumen",
		})
	}

	return ctx.JSON(fiber.Map{
		"data": jenisDokumen,
	})
}

// Create creates a new jenis dokumen
// POST /api/jenis-dokumen
func (c *JenisDokumenController) Create(ctx *fiber.Ctx) error {
	var input services.CreateJenisDokumenInput
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

	jenisDokumen, err := c.service.Create(&input)
	if err != nil {
		switch err {
		case services.ErrJenisDokumenKodeExists:
			return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Jenis dokumen with this kode already exists",
			})
		case services.ErrJenisDokumenKodeRequired:
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"errors": fiber.Map{"kode": "Kode is required"},
			})
		case services.ErrJenisDokumenNamaRequired:
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"errors": fiber.Map{"nama": "Nama is required"},
			})
		default:
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create jenis dokumen",
			})
		}
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Jenis dokumen created successfully",
		"data":    jenisDokumen,
	})
}


// Update updates a jenis dokumen
// PUT /api/jenis-dokumen/:id
func (c *JenisDokumenController) Update(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid jenis dokumen ID",
		})
	}

	var input services.UpdateJenisDokumenInput
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	jenisDokumen, err := c.service.Update(id, &input)
	if err != nil {
		switch err {
		case services.ErrJenisDokumenNotFound:
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Jenis dokumen not found",
			})
		case services.ErrJenisDokumenKodeExists:
			return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Jenis dokumen with this kode already exists",
			})
		case services.ErrJenisDokumenKodeRequired:
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"errors": fiber.Map{"kode": "Kode is required"},
			})
		case services.ErrJenisDokumenNamaRequired:
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"errors": fiber.Map{"nama": "Nama is required"},
			})
		default:
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update jenis dokumen",
			})
		}
	}

	return ctx.JSON(fiber.Map{
		"message": "Jenis dokumen updated successfully",
		"data":    jenisDokumen,
	})
}

// Delete deletes a jenis dokumen
// DELETE /api/jenis-dokumen/:id
func (c *JenisDokumenController) Delete(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid jenis dokumen ID",
		})
	}

	err = c.service.Delete(id)
	if err != nil {
		if err == services.ErrJenisDokumenNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Jenis dokumen not found",
			})
		}
		// Check for referential integrity error
		if refErr, ok := err.(*services.ErrJenisDokumenReferenced); ok {
			return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": refErr.Error(),
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete jenis dokumen",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "Jenis dokumen deleted successfully",
	})
}
