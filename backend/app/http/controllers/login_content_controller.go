package controllers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"dokumen-keuangan/app/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// LoginContentController handles login content endpoints
type LoginContentController struct {
	service *services.LoginContentService
}

// NewLoginContentController creates a new LoginContentController instance
func NewLoginContentController() *LoginContentController {
	return &LoginContentController{
		service: services.NewLoginContentService(),
	}
}

// GetAll retrieves all login contents
// GET /api/login-content
func (c *LoginContentController) GetAll(ctx *fiber.Ctx) error {
	contents, err := c.service.GetAll()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengambil data",
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    contents,
	})
}

// GetActive retrieves the currently active login content (public)
// GET /api/public/login-content/active
func (c *LoginContentController) GetActive(ctx *fiber.Ctx) error {
	content, err := c.service.GetActive()
	if err != nil {
		return ctx.JSON(fiber.Map{
			"success": false,
			"data":    nil,
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    content,
	})
}

// Create creates a new login content
// POST /api/login-content
func (c *LoginContentController) Create(ctx *fiber.Ctx) error {
	var input services.CreateLoginContentInput
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Data tidak valid",
		})
	}

	if input.Title == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Judul wajib diisi",
		})
	}

	content, err := c.service.Create(&input)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal menyimpan data",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Konten berhasil ditambahkan",
		"data":    content,
	})
}

// Update updates a login content
// PUT /api/login-content/:id
func (c *LoginContentController) Update(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID tidak valid",
		})
	}

	var input services.CreateLoginContentInput
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Data tidak valid",
		})
	}

	content, err := c.service.Update(id, &input)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Konten tidak ditemukan",
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Konten berhasil diupdate",
		"data":    content,
	})
}

// Delete deletes a login content
// DELETE /api/login-content/:id
func (c *LoginContentController) Delete(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID tidak valid",
		})
	}

	// Get content first to delete image
	content, _ := c.service.GetByID(id)
	if content != nil && content.ImageURL != "" {
		os.Remove("./storage/app/public" + strings.TrimPrefix(content.ImageURL, "/api/files"))
	}

	if err := c.service.Delete(id); err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Konten tidak ditemukan",
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Konten berhasil dihapus",
	})
}

// UploadImage uploads image for a login content
// POST /api/login-content/:id/image
func (c *LoginContentController) UploadImage(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID tidak valid",
		})
	}

	content, err := c.service.GetByID(id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Konten tidak ditemukan",
		})
	}

	file, err := ctx.FormFile("image")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "File tidak ditemukan",
		})
	}

	// Validate file type
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
	if !allowedExts[ext] {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format file tidak didukung",
		})
	}

	// Validate file size (max 2MB)
	if file.Size > 2*1024*1024 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Ukuran file maksimal 2MB",
		})
	}

	// Delete old image if exists
	if content.ImageURL != "" {
		os.Remove("./storage/app/public" + strings.TrimPrefix(content.ImageURL, "/api/files"))
	}

	// Create directory
	storageDir := "./storage/app/public/login-content"
	os.MkdirAll(storageDir, 0755)

	// Generate filename
	filename := fmt.Sprintf("login_content_%s_%d%s", id.String()[:8], time.Now().Unix(), ext)
	filePath := filepath.Join(storageDir, filename)

	// Save file
	if err := ctx.SaveFile(file, filePath); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal menyimpan file",
		})
	}

	// Update image URL
	imageURL := fmt.Sprintf("/api/files/login-content/%s", filename)
	if err := c.service.UpdateImageURL(id, imageURL); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal update URL gambar",
		})
	}

	return ctx.JSON(fiber.Map{
		"success":   true,
		"message":   "Gambar berhasil diupload",
		"image_url": imageURL,
	})
}
