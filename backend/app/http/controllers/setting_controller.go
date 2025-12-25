package controllers

import (
	"encoding/json"
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

// GetCountdownSettings retrieves countdown settings
// GET /api/settings/countdown
func (c *SettingController) GetCountdownSettings(ctx *fiber.Ctx) error {
	settingsMap, err := c.service.GetSettingsMap()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve settings",
		})
	}

	// Default values
	active := false
	if val, ok := settingsMap["countdown_active"]; ok && *val == "true" {
		active = true
	}

	data := fiber.Map{
		"active":      active,
		"title":       "",
		"description": "",
		"target_date": "",
	}

	if val, ok := settingsMap["countdown_title"]; ok {
		data["title"] = *val
	}
	if val, ok := settingsMap["countdown_description"]; ok {
		data["description"] = *val
	}
	if val, ok := settingsMap["countdown_target_date"]; ok {
		data["target_date"] = *val
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

// GetBrandingSettings retrieves branding settings
// GET /api/settings/branding
func (c *SettingController) GetBrandingSettings(ctx *fiber.Ctx) error {
	settingsMap, err := c.service.GetSettingsMap()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve settings",
		})
	}

	data := fiber.Map{
		"app_name":     "Sistem Pelimpahan",
		"app_subtitle": "Dana UP/GU",
		"logo_url":     "",
	}

	if val, ok := settingsMap["login_title"]; ok {
		data["app_name"] = *val
	}
	if val, ok := settingsMap["login_subtitle"]; ok {
		data["app_subtitle"] = *val
	}
	if val, ok := settingsMap["login_logo_url"]; ok {
		data["logo_url"] = *val
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    data,
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

// GetLockStatus - Check if year is locked
// GET /api/settings/lock-status
func (c *SettingController) GetLockStatus(ctx *fiber.Ctx) error {
	tahunAnggaran := ctx.Get("X-Tahun-Anggaran")
	if tahunAnggaran == "" {
		tahunAnggaran = "2025"
	}

	settingsMap, err := c.service.GetSettingsMap()
	if err != nil {
		return ctx.JSON(fiber.Map{
			"success": true,
			"data": fiber.Map{
				"locked":        false,
				"tahun":         tahunAnggaran,
				"locked_at":     nil,
				"locked_reason": "",
			},
		})
	}

	key := "tahun_dikunci_" + tahunAnggaran
	if val, ok := settingsMap[key]; ok && val != nil {
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(*val), &data); err == nil {
			return ctx.JSON(fiber.Map{
				"success": true,
				"data":    data,
			})
		}
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"locked":        false,
			"tahun":         tahunAnggaran,
			"locked_at":     nil,
			"locked_reason": "",
		},
	})
}

// ToggleLock - Lock/Unlock year (super_admin only)
// POST /api/settings/toggle-lock
func (c *SettingController) ToggleLock(ctx *fiber.Ctx) error {
	var req struct {
		Locked bool   `json:"locked"`
		Reason string `json:"reason"`
	}

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request",
		})
	}

	tahunAnggaran := ctx.Get("X-Tahun-Anggaran")
	if tahunAnggaran == "" {
		tahunAnggaran = "2025"
	}

	key := "tahun_dikunci_" + tahunAnggaran
	data := map[string]interface{}{
		"locked":        req.Locked,
		"tahun":         tahunAnggaran,
		"locked_at":     nil,
		"locked_reason": req.Reason,
	}
	if req.Locked {
		data["locked_at"] = time.Now().Format("2006-01-02 15:04:05")
	}
	jsonData, _ := json.Marshal(data)
	jsonStr := string(jsonData)

	input := &services.UpdateSettingInput{
		Settings: map[string]*string{
			key: &jsonStr,
		},
	}
	if _, err := c.service.Update(input); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal menyimpan pengaturan",
		})
	}

	action := "dibuka"
	if req.Locked {
		action = "dikunci"
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": fmt.Sprintf("Tahun Anggaran %s berhasil %s", tahunAnggaran, action),
		"data":    data,
	})
}

// IsYearLocked - Helper function to check if year is locked
func IsYearLocked(tahunAnggaran string) bool {
	settingService := services.NewSettingService()
	settingsMap, err := settingService.GetSettingsMap()
	if err != nil {
		return false
	}

	key := "tahun_dikunci_" + tahunAnggaran
	if val, ok := settingsMap[key]; ok && val != nil {
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(*val), &data); err == nil {
			locked, ok := data["locked"].(bool)
			return ok && locked
		}
	}
	return false
}
