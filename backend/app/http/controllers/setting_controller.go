package controllers

import (
	"dokumen-keuangan/app/services"

	"github.com/gofiber/fiber/v2"
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
