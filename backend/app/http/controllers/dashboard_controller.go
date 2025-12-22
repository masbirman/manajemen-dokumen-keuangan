package controllers

import (
	"dokumen-keuangan/app/http/middleware"
	"dokumen-keuangan/app/models"
	"dokumen-keuangan/app/services"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type DashboardController struct {
	dokumenService      *services.DokumenService
	unitKerjaService    *services.UnitKerjaService
	jenisDokumenService *services.JenisDokumenService
	userService         *services.UserService
}

func NewDashboardController() *DashboardController {
	return &DashboardController{
		dokumenService:      services.NewDokumenService(),
		unitKerjaService:    services.NewUnitKerjaService(),
		jenisDokumenService: services.NewJenisDokumenService(),
		userService:         services.NewUserService(),
	}
}

// GetStats returns dashboard statistics
// GET /api/dashboard/stats
func (c *DashboardController) GetStats(ctx *fiber.Ctx) error {
	// Get year from header
	yearHeader := ctx.Get("X-Tahun-Anggaran")
	var year int
	if yearHeader != "" {
		year, _ = strconv.Atoi(yearHeader)
	}
	if year == 0 {
		year = time.Now().Year() // Default to current year
	}

	userRole := middleware.GetUserRoleFromContext(ctx)
	userIDStr := middleware.GetUserIDFromContext(ctx)
	userID, _ := uuid.Parse(userIDStr)

	// Get stats from services (passing year filter)
	// We need to extend services/repositories to support CountByYear or similar
	
	// Placeholder implementation - will need actual service methods
	totalDokumen := c.dokumenService.CountByYear(year, models.UserRole(userRole), userID)
	totalUnit, _ := c.unitKerjaService.Count() // Unit limits might not be year dependent, or maybe they are? Usually not.
	totalJenis, _ := c.jenisDokumenService.Count()
	totalUser, _ := c.userService.Count()
	
	recentDokumen, _ := c.dokumenService.GetRecent(5, year, models.UserRole(userRole), userID)

	return ctx.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"total_pelimpahan":  totalDokumen,
			"total_unit":        totalUnit,
			"total_jenis":       totalJenis,
			"total_user":        totalUser,
			"recent_pelimpahan": recentDokumen,
		},
	})
}
