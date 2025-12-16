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

// SettingController handles setting endpoints
type SettingController struct {
	service *services.SettingService
}

// NewSettingController creates a new SettingController instance
func NewSettingController() *SettingController {
	return &SettingController{
		service: services.NewSettingService(),
	}
}

// GetAll retrieves all settings
// GET /api/settings
func (c *SettingController) GetAll(ctx *fiber.Ctx) error {
	settings, err := c.service.GetAll()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve settings",
		})
	}

	// Convert to map for easier frontend consumption
	settingsMap, _ := c.service.GetSettingsMap()

	return ctx.JSON(fiber.Map{
		"data": settings,
		"map":  settingsMap,
	})
}

// Update updates settings
// PUT /api/settings
func (c *SettingController) Update(ctx *fiber.Ctx) error {
	var input services.UpdateSettingInput
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	settings, err := c.service.Update(&input)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update settings",
		})
	}

	settingsMap, _ := c.service.GetSettingsMap()

	return ctx.JSON(fiber.Map{
		"message": "Settings updated successfully",
		"data":    settings,
		"map":     settingsMap,
	})
}

// GetLoginSettings retrieves login page settings (public endpoint)
// GET /api/public/login-settings
func (c *SettingController) GetLoginSettings(ctx *fiber.Ctx) error {
	settingsMap, err := c.service.GetSettingsMap()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve settings",
		})
	}

	// Only return login page related settings
	loginKeys := []string{
		"login_logo_url",
		"login_logo_size",
		"login_title",
		"login_subtitle",
		"login_info_title",
		"login_info_content",
		"login_bg_color",
		"login_accent_color",
		"login_font_family",
		"login_title_size",
		"login_subtitle_size",
	}

	loginSettings := make(map[string]*string)
	for _, key := range loginKeys {
		if val, ok := settingsMap[key]; ok {
			loginSettings[key] = val
		}
	}

	return ctx.JSON(fiber.Map{
		"data": loginSettings,
	})
}

// UploadLogo uploads login page logo
// POST /api/settings/upload-logo
func (c *SettingController) UploadLogo(ctx *fiber.Ctx) error {
	// Get the file from the request
	file, err := ctx.FormFile("logo")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Logo file is required",
		})
	}

	// Validate file type
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".svg": true, ".webp": true}
	if !allowedExts[ext] {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid file type. Allowed: jpg, jpeg, png, gif, svg, webp",
		})
	}

	// Validate file size (max 2MB)
	if file.Size > 2*1024*1024 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File size must be less than 2MB",
		})
	}

	// Create storage directory if not exists
	storageDir := "./storage/app/public/logos"
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create storage directory",
		})
	}

	// Generate unique filename
	filename := fmt.Sprintf("logo_%s_%d%s", uuid.New().String()[:8], time.Now().Unix(), ext)
	filePath := filepath.Join(storageDir, filename)

	// Save the file
	if err := ctx.SaveFile(file, filePath); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save file",
		})
	}

	// Generate URL for the logo
	logoURL := fmt.Sprintf("/api/files/logos/%s", filename)

	// Update the login_logo_url setting
	input := &services.UpdateSettingInput{
		Settings: map[string]*string{
			"login_logo_url": &logoURL,
		},
	}
	if _, err := c.service.Update(input); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update logo setting",
		})
	}

	return ctx.JSON(fiber.Map{
		"message":  "Logo uploaded successfully",
		"logo_url": logoURL,
	})
}
