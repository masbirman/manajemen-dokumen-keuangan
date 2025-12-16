package controllers

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"dokumen-keuangan/app/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// PetunjukController handles petunjuk endpoints
type PetunjukController struct {
	service *services.PetunjukService
}

// NewPetunjukController creates a new PetunjukController instance
func NewPetunjukController() *PetunjukController {
	return &PetunjukController{
		service: services.NewPetunjukService(),
	}
}

// GetAll retrieves all petunjuk with pagination
// GET /api/petunjuk
func (c *PetunjukController) GetAll(ctx *fiber.Ctx) error {
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.Query("page_size", "10"))
	search := ctx.Query("search", "")

	result, err := c.service.GetAll(page, pageSize, search)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve petunjuk",
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

// GetByHalaman retrieves petunjuk for a specific page
// GET /api/petunjuk/halaman/:halaman
func (c *PetunjukController) GetByHalaman(ctx *fiber.Ctx) error {
	halaman := ctx.Params("halaman")
	
	petunjuks, err := c.service.GetByHalaman(halaman)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve petunjuk",
		})
	}

	return ctx.JSON(fiber.Map{
		"data": petunjuks,
	})
}

// GetByID retrieves a petunjuk by ID
// GET /api/petunjuk/:id
func (c *PetunjukController) GetByID(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid petunjuk ID",
		})
	}

	petunjuk, err := c.service.GetByID(id)
	if err != nil {
		if err == services.ErrPetunjukNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Petunjuk not found",
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve petunjuk",
		})
	}

	return ctx.JSON(fiber.Map{
		"data": petunjuk,
	})
}

// Create creates a new petunjuk
// POST /api/petunjuk
func (c *PetunjukController) Create(ctx *fiber.Ctx) error {
	var input services.CreatePetunjukInput
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if validationErrors := input.Validate(); len(validationErrors) > 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": validationErrors,
		})
	}

	petunjuk, err := c.service.Create(&input)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create petunjuk",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Petunjuk created successfully",
		"data":    petunjuk,
	})
}

// Update updates a petunjuk
// PUT /api/petunjuk/:id
func (c *PetunjukController) Update(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid petunjuk ID",
		})
	}

	var input services.UpdatePetunjukInput
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	petunjuk, err := c.service.Update(id, &input)
	if err != nil {
		if err == services.ErrPetunjukNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Petunjuk not found",
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update petunjuk",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "Petunjuk updated successfully",
		"data":    petunjuk,
	})
}

// Delete deletes a petunjuk
// DELETE /api/petunjuk/:id
func (c *PetunjukController) Delete(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid petunjuk ID",
		})
	}

	err = c.service.Delete(id)
	if err != nil {
		if err == services.ErrPetunjukNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Petunjuk not found",
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete petunjuk",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "Petunjuk deleted successfully",
	})
}

// UploadImage handles image upload for petunjuk
// POST /api/petunjuk/upload-image
func (c *PetunjukController) UploadImage(ctx *fiber.Ctx) error {
	file, err := ctx.FormFile("file")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File is required",
		})
	}

	// Validate file type
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	if !allowedExts[ext] {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid file type. Allowed: jpg, jpeg, png, gif, webp",
		})
	}

	// Validate file size (max 5MB)
	if file.Size > 5*1024*1024 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File size must be less than 5MB",
		})
	}

	// Create storage directory if not exists
	storageDir := "./storage/app/public/petunjuk"
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create storage directory",
		})
	}

	// Generate unique filename
	filename := fmt.Sprintf("petunjuk_%s_%d%s", uuid.New().String()[:8], time.Now().Unix(), ext)
	filePath := filepath.Join(storageDir, filename)

	// Save the file
	if err := ctx.SaveFile(file, filePath); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save file",
		})
	}

	// Generate URL
	fileURL := fmt.Sprintf("/api/files/petunjuk/%s", filename)

	return ctx.JSON(fiber.Map{
		"message": "Image uploaded successfully",
		"url":     fileURL,
	})
}
