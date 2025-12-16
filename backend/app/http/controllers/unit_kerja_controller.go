package controllers

import (
	"strconv"

	"dokumen-keuangan/app/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// UnitKerjaController handles unit kerja endpoints
type UnitKerjaController struct {
	service      *services.UnitKerjaService
	excelService *services.ExcelService
}

// NewUnitKerjaController creates a new UnitKerjaController instance
func NewUnitKerjaController() *UnitKerjaController {
	return &UnitKerjaController{
		service:      services.NewUnitKerjaService(),
		excelService: services.NewExcelService(),
	}
}

// GetAll retrieves all unit kerja with pagination
// GET /api/unit-kerja
func (c *UnitKerjaController) GetAll(ctx *fiber.Ctx) error {
	// Parse query parameters
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.Query("page_size", "10"))
	search := ctx.Query("search", "")

	result, err := c.service.GetAll(page, pageSize, search)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve unit kerja",
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

// GetAllActive retrieves all active unit kerja (for dropdowns)
// GET /api/unit-kerja/active
func (c *UnitKerjaController) GetAllActive(ctx *fiber.Ctx) error {
	unitKerjas, err := c.service.GetAllActive()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve unit kerja",
		})
	}

	return ctx.JSON(fiber.Map{
		"data": unitKerjas,
	})
}


// GetByID retrieves a unit kerja by ID
// GET /api/unit-kerja/:id
func (c *UnitKerjaController) GetByID(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid unit kerja ID",
		})
	}

	unitKerja, err := c.service.GetByID(id)
	if err != nil {
		if err == services.ErrUnitKerjaNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Unit kerja not found",
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve unit kerja",
		})
	}

	return ctx.JSON(fiber.Map{
		"data": unitKerja,
	})
}

// Create creates a new unit kerja
// POST /api/unit-kerja
func (c *UnitKerjaController) Create(ctx *fiber.Ctx) error {
	var input services.CreateUnitKerjaInput
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

	unitKerja, err := c.service.Create(&input)
	if err != nil {
		switch err {
		case services.ErrUnitKerjaKodeExists:
			return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Unit kerja with this kode already exists",
			})
		case services.ErrUnitKerjaKodeRequired:
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"errors": fiber.Map{"kode": "Kode is required"},
			})
		case services.ErrUnitKerjaNamaRequired:
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"errors": fiber.Map{"nama": "Nama is required"},
			})
		default:
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create unit kerja",
			})
		}
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Unit kerja created successfully",
		"data":    unitKerja,
	})
}

// Update updates a unit kerja
// PUT /api/unit-kerja/:id
func (c *UnitKerjaController) Update(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid unit kerja ID",
		})
	}

	var input services.UpdateUnitKerjaInput
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	unitKerja, err := c.service.Update(id, &input)
	if err != nil {
		switch err {
		case services.ErrUnitKerjaNotFound:
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Unit kerja not found",
			})
		case services.ErrUnitKerjaKodeExists:
			return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Unit kerja with this kode already exists",
			})
		case services.ErrUnitKerjaKodeRequired:
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"errors": fiber.Map{"kode": "Kode is required"},
			})
		case services.ErrUnitKerjaNamaRequired:
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"errors": fiber.Map{"nama": "Nama is required"},
			})
		default:
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update unit kerja",
			})
		}
	}

	return ctx.JSON(fiber.Map{
		"message": "Unit kerja updated successfully",
		"data":    unitKerja,
	})
}

// Delete deletes a unit kerja
// DELETE /api/unit-kerja/:id
func (c *UnitKerjaController) Delete(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid unit kerja ID",
		})
	}

	err = c.service.Delete(id)
	if err != nil {
		if err == services.ErrUnitKerjaNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Unit kerja not found",
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete unit kerja",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "Unit kerja deleted successfully",
	})
}


// GetTemplate downloads Excel template for Unit Kerja import
// GET /api/unit-kerja/template
func (c *UnitKerjaController) GetTemplate(ctx *fiber.Ctx) error {
	buf, err := c.excelService.GenerateUnitKerjaTemplate()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate template",
		})
	}

	ctx.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Set("Content-Disposition", "attachment; filename=template_unit_kerja.xlsx")

	return ctx.Send(buf.Bytes())
}

// Import imports Unit Kerja data from Excel file
// POST /api/unit-kerja/import
func (c *UnitKerjaController) Import(ctx *fiber.Ctx) error {
	// Get uploaded file
	file, err := ctx.FormFile("file")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File is required",
		})
	}

	// Validate file extension
	if file.Header.Get("Content-Type") != "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" &&
		file.Header.Get("Content-Type") != "application/vnd.ms-excel" {
		// Also check by filename extension
		filename := file.Filename
		if len(filename) < 5 || (filename[len(filename)-5:] != ".xlsx" && filename[len(filename)-4:] != ".xls") {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid file type. Please upload an Excel file (.xlsx or .xls)",
			})
		}
	}

	// Open file
	f, err := file.Open()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to read file",
		})
	}
	defer f.Close()

	// Import data
	result, err := c.excelService.ImportUnitKerja(f)
	if err != nil {
		switch err {
		case services.ErrInvalidExcelFile:
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid Excel file",
			})
		case services.ErrEmptyExcelFile:
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Excel file is empty or has no data rows",
			})
		default:
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to import data",
			})
		}
	}

	return ctx.JSON(fiber.Map{
		"message": "Import completed",
		"data":    result,
	})
}
