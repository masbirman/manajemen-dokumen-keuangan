package controllers

import (
	"os"
	"strconv"

	"dokumen-keuangan/app/http/middleware"
	"dokumen-keuangan/app/models"
	"dokumen-keuangan/app/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// DokumenController handles dokumen endpoints
type DokumenController struct {
	service *services.DokumenService
}

// NewDokumenController creates a new DokumenController instance
func NewDokumenController() *DokumenController {
	return &DokumenController{
		service: services.NewDokumenService(),
	}
}

// GetAll retrieves all dokumen with pagination and filters
// GET /api/dokumen
func (c *DokumenController) GetAll(ctx *fiber.Ctx) error {
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.Query("page_size", "10"))

	// Get filter parameters
	var filter services.DokumenFilterInput
	
	if unitKerjaID := ctx.Query("unit_kerja_id"); unitKerjaID != "" {
		id, err := uuid.Parse(unitKerjaID)
		if err == nil {
			filter.UnitKerjaID = &id
		}
	}

	if pptkID := ctx.Query("pptk_id"); pptkID != "" {
		id, err := uuid.Parse(pptkID)
		if err == nil {
			filter.PPTKID = &id
		}
	}
	
	if startDate := ctx.Query("start_date"); startDate != "" {
		filter.StartDate = &startDate
	}
	
	if endDate := ctx.Query("end_date"); endDate != "" {
		filter.EndDate = &endDate
	}

	// Get user info from context
	userRole := middleware.GetUserRoleFromContext(ctx)
	userIDStr := middleware.GetUserIDFromContext(ctx)
	userID, _ := uuid.Parse(userIDStr)

	// Get Year from Header
	yearHeader := ctx.Get("X-Tahun-Anggaran")
	var year int
	if yearHeader != "" {
		year, _ = strconv.Atoi(yearHeader)
	}

	result, err := c.service.GetAll(page, pageSize, &filter, models.UserRole(userRole), userID, year)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve dokumen",
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


// GetByID retrieves a dokumen by ID
// GET /api/dokumen/:id
func (c *DokumenController) GetByID(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid dokumen ID",
		})
	}

	// Check access permission
	userRole := middleware.GetUserRoleFromContext(ctx)
	userIDStr := middleware.GetUserIDFromContext(ctx)
	userID, _ := uuid.Parse(userIDStr)

	if !c.service.CanAccessDokumen(id, userRole, userID) {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied",
		})
	}

	dokumen, err := c.service.GetByID(id)
	if err != nil {
		if err == services.ErrDokumenNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Dokumen not found",
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve dokumen",
		})
	}

	return ctx.JSON(fiber.Map{
		"data": dokumen,
	})
}

// Create creates a new dokumen
// POST /api/dokumen
func (c *DokumenController) Create(ctx *fiber.Ctx) error {
	// Parse form data
	var input services.CreateDokumenInput

	// Get nomor_dokumen and tanggal_dokumen
	input.NomorDokumen = ctx.FormValue("nomor_dokumen")
	input.TanggalDokumen = ctx.FormValue("tanggal_dokumen")

	// Get UUIDs from form
	if unitKerjaID := ctx.FormValue("unit_kerja_id"); unitKerjaID != "" {
		id, err := uuid.Parse(unitKerjaID)
		if err == nil {
			input.UnitKerjaID = id
		}
	}
	if pptkID := ctx.FormValue("pptk_id"); pptkID != "" {
		id, err := uuid.Parse(pptkID)
		if err == nil {
			input.PPTKID = id
		}
	}
	if jenisDokumenID := ctx.FormValue("jenis_dokumen_id"); jenisDokumenID != "" {
		id, err := uuid.Parse(jenisDokumenID)
		if err == nil {
			input.JenisDokumenID = id
		}
	}
	if sumberDanaID := ctx.FormValue("sumber_dana_id"); sumberDanaID != "" {
		id, err := uuid.Parse(sumberDanaID)
		if err == nil {
			input.SumberDanaID = id
		}
	}

	// Get nilai
	if nilai := ctx.FormValue("nilai"); nilai != "" {
		n, err := strconv.ParseFloat(nilai, 64)
		if err == nil {
			input.Nilai = n
		}
	}

	// Get uraian
	input.Uraian = ctx.FormValue("uraian")

	// Validate input
	if validationErrors := input.Validate(); len(validationErrors) > 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": validationErrors,
		})
	}

	// Get file
	file, err := ctx.FormFile("file")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File is required",
		})
	}

	// Get creator ID from context
	userIDStr := middleware.GetUserIDFromContext(ctx)
	userID, _ := uuid.Parse(userIDStr)

	dokumen, err := c.service.Create(&input, file, userID)
	if err != nil {
		switch err {
		case services.ErrInvalidFileType:
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid file type. Only PDF files are allowed",
			})
		case services.ErrFileTooLarge:
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "File too large. Maximum size: 20MB",
			})
		default:
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Dokumen created successfully",
		"data":    dokumen,
	})
}

// GetFile downloads the dokumen file
// GET /api/dokumen/:id/file
func (c *DokumenController) GetFile(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid dokumen ID",
		})
	}

	// Check access permission
	userRole := middleware.GetUserRoleFromContext(ctx)
	userIDStr := middleware.GetUserIDFromContext(ctx)
	userID, _ := uuid.Parse(userIDStr)

	if !c.service.CanAccessDokumen(id, userRole, userID) {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied",
		})
	}

	filePath, err := c.service.GetFilePath(id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Dokumen not found",
		})
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "File not found",
		})
	}

	return ctx.SendFile(filePath)
}

// Update updates a dokumen
// PUT /api/dokumen/:id
func (c *DokumenController) Update(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid dokumen ID",
		})
	}

	// Check access permission
	userRole := middleware.GetUserRoleFromContext(ctx)
	userIDStr := middleware.GetUserIDFromContext(ctx)
	userID, _ := uuid.Parse(userIDStr)

	if !c.service.CanAccessDokumen(id, userRole, userID) {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied",
		})
	}

	// Parse form data
	var input services.UpdateDokumenInput

	input.NomorDokumen = ctx.FormValue("nomor_dokumen")
	input.TanggalDokumen = ctx.FormValue("tanggal_dokumen")
	input.Uraian = ctx.FormValue("uraian")

	// Get UUIDs from form
	if unitKerjaID := ctx.FormValue("unit_kerja_id"); unitKerjaID != "" {
		uid, err := uuid.Parse(unitKerjaID)
		if err == nil {
			input.UnitKerjaID = &uid
		}
	}
	if pptkID := ctx.FormValue("pptk_id"); pptkID != "" {
		uid, err := uuid.Parse(pptkID)
		if err == nil {
			input.PPTKID = &uid
		}
	}
	if jenisDokumenID := ctx.FormValue("jenis_dokumen_id"); jenisDokumenID != "" {
		uid, err := uuid.Parse(jenisDokumenID)
		if err == nil {
			input.JenisDokumenID = &uid
		}
	}
	if sumberDanaID := ctx.FormValue("sumber_dana_id"); sumberDanaID != "" {
		uid, err := uuid.Parse(sumberDanaID)
		if err == nil {
			input.SumberDanaID = &uid
		}
	}

	// Get nilai
	if nilai := ctx.FormValue("nilai"); nilai != "" {
		n, err := strconv.ParseFloat(nilai, 64)
		if err == nil {
			input.Nilai = &n
		}
	}

	// Get file (optional for update)
	file, _ := ctx.FormFile("file")

	dokumen, err := c.service.Update(id, &input, file)
	if err != nil {
		switch err {
		case services.ErrDokumenNotFound:
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Dokumen not found",
			})
		case services.ErrInvalidFileType:
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid file type. Only PDF files are allowed",
			})
		default:
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	return ctx.JSON(fiber.Map{
		"message": "Dokumen updated successfully",
		"data":    dokumen,
	})
}

// Delete deletes a dokumen by ID
// DELETE /api/dokumen/:id
func (c *DokumenController) Delete(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid dokumen ID",
		})
	}

	// Check access permission
	userRole := middleware.GetUserRoleFromContext(ctx)
	userIDStr := middleware.GetUserIDFromContext(ctx)
	userID, _ := uuid.Parse(userIDStr)

	if !c.service.CanAccessDokumen(id, userRole, userID) {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied",
		})
	}

	err = c.service.Delete(id)
	if err != nil {
		if err == services.ErrDokumenNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Dokumen not found",
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete dokumen",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "Dokumen deleted successfully",
	})
}
