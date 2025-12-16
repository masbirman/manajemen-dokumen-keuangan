package controllers

import (
	"strconv"

	"dokumen-keuangan/app/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// PPTKController handles PPTK endpoints
type PPTKController struct {
	service      *services.PPTKService
	excelService *services.ExcelService
}

// NewPPTKController creates a new PPTKController instance
func NewPPTKController() *PPTKController {
	return &PPTKController{
		service:      services.NewPPTKService(),
		excelService: services.NewExcelService(),
	}
}

// GetAll retrieves all PPTK with pagination
// GET /api/pptk
func (c *PPTKController) GetAll(ctx *fiber.Ctx) error {
	// Parse query parameters
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.Query("page_size", "10"))
	search := ctx.Query("search", "")
	unitKerjaID := ctx.Query("unit_kerja_id", "")

	result, err := c.service.GetAllWithFilter(page, pageSize, search, unitKerjaID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve PPTK",
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

// GetAllActive retrieves all active PPTK (for dropdowns)
// GET /api/pptk/active
func (c *PPTKController) GetAllActive(ctx *fiber.Ctx) error {
	pptks, err := c.service.GetAllActive()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve PPTK",
		})
	}

	return ctx.JSON(fiber.Map{
		"data": pptks,
	})
}


// GetByUnitKerja retrieves all PPTK for a specific Unit Kerja
// GET /api/pptk/by-unit-kerja/:unitKerjaId
func (c *PPTKController) GetByUnitKerja(ctx *fiber.Ctx) error {
	unitKerjaIDParam := ctx.Params("unitKerjaId")
	unitKerjaID, err := uuid.Parse(unitKerjaIDParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Unit Kerja ID",
		})
	}

	pptks, err := c.service.GetByUnitKerja(unitKerjaID)
	if err != nil {
		if err == services.ErrPPTKUnitKerjaNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Unit Kerja not found",
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve PPTK",
		})
	}

	return ctx.JSON(fiber.Map{
		"data": pptks,
	})
}

// GetByID retrieves a PPTK by ID
// GET /api/pptk/:id
func (c *PPTKController) GetByID(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid PPTK ID",
		})
	}

	pptk, err := c.service.GetByID(id)
	if err != nil {
		if err == services.ErrPPTKNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "PPTK not found",
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve PPTK",
		})
	}

	return ctx.JSON(fiber.Map{
		"data": pptk,
	})
}

// Create creates a new PPTK
// POST /api/pptk
func (c *PPTKController) Create(ctx *fiber.Ctx) error {
	var input services.CreatePPTKInput
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

	pptk, err := c.service.Create(&input)
	if err != nil {
		switch err {
		case services.ErrPPTKNIPExists:
			return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "PPTK with this NIP already exists",
			})
		case services.ErrPPTKNIPRequired:
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"errors": fiber.Map{"nip": "NIP is required"},
			})
		case services.ErrPPTKNamaRequired:
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"errors": fiber.Map{"nama": "Nama is required"},
			})
		case services.ErrPPTKUnitKerjaRequired:
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"errors": fiber.Map{"unit_kerja_id": "Unit Kerja is required"},
			})
		case services.ErrPPTKUnitKerjaNotFound:
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"errors": fiber.Map{"unit_kerja_id": "Unit Kerja not found"},
			})
		default:
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create PPTK",
			})
		}
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "PPTK created successfully",
		"data":    pptk,
	})
}


// Update updates a PPTK
// PUT /api/pptk/:id
func (c *PPTKController) Update(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid PPTK ID",
		})
	}

	var input services.UpdatePPTKInput
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	pptk, err := c.service.Update(id, &input)
	if err != nil {
		switch err {
		case services.ErrPPTKNotFound:
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "PPTK not found",
			})
		case services.ErrPPTKNIPExists:
			return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "PPTK with this NIP already exists",
			})
		case services.ErrPPTKNIPRequired:
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"errors": fiber.Map{"nip": "NIP is required"},
			})
		case services.ErrPPTKNamaRequired:
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"errors": fiber.Map{"nama": "Nama is required"},
			})
		case services.ErrPPTKUnitKerjaNotFound:
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"errors": fiber.Map{"unit_kerja_id": "Unit Kerja not found"},
			})
		default:
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update PPTK",
			})
		}
	}

	return ctx.JSON(fiber.Map{
		"message": "PPTK updated successfully",
		"data":    pptk,
	})
}

// Delete deletes a PPTK
// DELETE /api/pptk/:id
func (c *PPTKController) Delete(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid PPTK ID",
		})
	}

	err = c.service.Delete(id)
	if err != nil {
		if err == services.ErrPPTKNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "PPTK not found",
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete PPTK",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "PPTK deleted successfully",
	})
}


// GetTemplate downloads Excel template for PPTK import
// GET /api/pptk/template
func (c *PPTKController) GetTemplate(ctx *fiber.Ctx) error {
	buf, err := c.excelService.GeneratePPTKTemplate()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate template",
		})
	}

	ctx.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Set("Content-Disposition", "attachment; filename=template_pptk.xlsx")

	return ctx.Send(buf.Bytes())
}

// Import imports PPTK data from Excel file
// POST /api/pptk/import
func (c *PPTKController) Import(ctx *fiber.Ctx) error {
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
	result, err := c.excelService.ImportPPTK(f)
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


// UploadAvatar uploads avatar for a PPTK
// POST /api/pptk/:id/avatar
func (c *PPTKController) UploadAvatar(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid PPTK ID",
		})
	}

	// Get PPTK
	pptk, err := c.service.GetByID(id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "PPTK not found",
		})
	}

	// Get uploaded file
	file, err := ctx.FormFile("avatar")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Avatar file is required",
		})
	}

	// Upload file
	fileService := services.NewFileService()
	avatarPath, err := fileService.UploadAvatar(file, "pptk", id)
	if err != nil {
		switch err {
		case services.ErrInvalidFileType:
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid file type. Allowed: JPG, PNG, GIF, WebP",
			})
		case services.ErrFileTooLarge:
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "File too large. Maximum size: 5MB",
			})
		default:
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to upload avatar",
			})
		}
	}

	// Update PPTK avatar path
	updatedPPTK, err := c.service.Update(id, &services.UpdatePPTKInput{
		AvatarPath: &avatarPath,
	})
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update PPTK avatar",
		})
	}

	_ = pptk // Suppress unused variable warning

	return ctx.JSON(fiber.Map{
		"message": "Avatar uploaded successfully",
		"data":    updatedPPTK,
	})
}
