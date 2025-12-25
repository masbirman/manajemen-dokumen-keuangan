package routes

import (
	"dokumen-keuangan/app/http/controllers"
	"dokumen-keuangan/app/http/middleware"
	"dokumen-keuangan/app/models"
	"dokumen-keuangan/config"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configures all application routes
func SetupRoutes(app *fiber.App, cfg *config.Config) {
	// API group
	api := app.Group("/api")

	// Static files route for avatars and uploads
	api.Static("/files", "./storage/app/public")
	
	// Serve uploads directory from root to match frontend expectation
	app.Static("/uploads", "./storage/app/public/uploads")

	// Auth routes (public)
	setupAuthRoutes(api, cfg)

	// Public settings for login page (no auth required)
	settingController := controllers.NewSettingController()
	api.Get("/public/login-settings", settingController.GetLoginSettings)

	// Public login content for login page (no auth required)
	loginContentController := controllers.NewLoginContentController()
	api.Get("/public/login-content/active", loginContentController.GetActive)

	// Protected routes (require authentication)
	protected := api.Group("", middleware.AuthMiddleware(cfg))

	// Unit Kerja routes (Admin+)
	setupUnitKerjaRoutes(protected)

	// PPTK routes (Admin+)
	setupPPTKRoutes(protected)

	// Sumber Dana routes (Admin+)
	setupSumberDanaRoutes(protected)

	// Jenis Dokumen routes (Admin+)
	setupJenisDokumenRoutes(protected)

	// User routes (Super Admin only)
	setupUserRoutes(protected)

	// Dokumen routes (All authenticated users)
	setupDokumenRoutes(protected)

	// Dashboard routes (All authenticated users)
	setupDashboardRoutes(protected)

	// Settings routes (Super Admin only)
	setupSettingRoutes(protected)

	// Petunjuk routes
	setupPetunjukRoutes(protected)

	// Login Content routes (Super Admin only)
	setupLoginContentRoutes(protected)
}

// setupDashboardRoutes configures dashboard routes
func setupDashboardRoutes(api fiber.Router) {
	dashboardController := controllers.NewDashboardController()

	dashboard := api.Group("/dashboard")
	dashboard.Get("/stats", dashboardController.GetStats) // Allow all authenticated
}

// setupAuthRoutes configures authentication routes
func setupAuthRoutes(api fiber.Router, cfg *config.Config) {
	authController := controllers.NewAuthController(cfg)

	auth := api.Group("/auth")
	auth.Post("/login", authController.Login)
	auth.Post("/logout", authController.Logout)
	auth.Post("/refresh", authController.Refresh)
	auth.Get("/me", authController.Me)
	
	// Profile management
	auth.Put("/profile", authController.UpdateProfile)
	auth.Post("/profile/avatar", authController.UpdateAvatar)
}

// setupUnitKerjaRoutes configures unit kerja routes
func setupUnitKerjaRoutes(api fiber.Router) {
	unitKerjaController := controllers.NewUnitKerjaController()

	unitKerja := api.Group("/unit-kerja")

	// GET routes - accessible by Admin and above
	unitKerja.Get("/", middleware.RequireRole(models.RoleAdmin), unitKerjaController.GetAll)
	unitKerja.Get("/active", middleware.RequireRole(models.RoleOperator), unitKerjaController.GetAllActive)
	unitKerja.Get("/template", middleware.RequireRole(models.RoleAdmin), unitKerjaController.GetTemplate)
	unitKerja.Get("/:id", middleware.RequireRole(models.RoleAdmin), unitKerjaController.GetByID)

	// CUD routes - accessible by Admin and above
	unitKerja.Post("/", middleware.RequireRole(models.RoleAdmin), unitKerjaController.Create)
	unitKerja.Post("/import", middleware.RequireRole(models.RoleAdmin), unitKerjaController.Import)
	unitKerja.Put("/:id", middleware.RequireRole(models.RoleAdmin), unitKerjaController.Update)
	unitKerja.Delete("/:id", middleware.RequireRole(models.RoleAdmin), unitKerjaController.Delete)
}


// setupPPTKRoutes configures PPTK routes
func setupPPTKRoutes(api fiber.Router) {
	pptkController := controllers.NewPPTKController()

	pptk := api.Group("/pptk")

	// GET routes - accessible by Admin and above (except active and by-unit-kerja for Operator)
	pptk.Get("/", middleware.RequireRole(models.RoleAdmin), pptkController.GetAll)
	pptk.Get("/active", middleware.RequireRole(models.RoleOperator), pptkController.GetAllActive)
	pptk.Get("/template", middleware.RequireRole(models.RoleAdmin), pptkController.GetTemplate)
	pptk.Get("/by-unit-kerja/:unitKerjaId", middleware.RequireRole(models.RoleOperator), pptkController.GetByUnitKerja)
	pptk.Get("/:id", middleware.RequireRole(models.RoleAdmin), pptkController.GetByID)

	// CUD routes - accessible by Admin and above
	pptk.Post("/", middleware.RequireRole(models.RoleAdmin), pptkController.Create)
	pptk.Post("/import", middleware.RequireRole(models.RoleAdmin), pptkController.Import)
	pptk.Post("/:id/avatar", middleware.RequireRole(models.RoleAdmin), pptkController.UploadAvatar)
	pptk.Put("/:id", middleware.RequireRole(models.RoleAdmin), pptkController.Update)
	pptk.Delete("/:id", middleware.RequireRole(models.RoleAdmin), pptkController.Delete)
}


// setupSumberDanaRoutes configures sumber dana routes
func setupSumberDanaRoutes(api fiber.Router) {
	sumberDanaController := controllers.NewSumberDanaController()

	sumberDana := api.Group("/sumber-dana")

	// GET routes - accessible by Admin and above (except active for Operator)
	sumberDana.Get("/", middleware.RequireRole(models.RoleAdmin), sumberDanaController.GetAll)
	sumberDana.Get("/active", middleware.RequireRole(models.RoleOperator), sumberDanaController.GetAllActive)
	sumberDana.Get("/:id", middleware.RequireRole(models.RoleAdmin), sumberDanaController.GetByID)

	// CUD routes - accessible by Admin and above
	sumberDana.Post("/", middleware.RequireRole(models.RoleAdmin), sumberDanaController.Create)
	sumberDana.Put("/:id", middleware.RequireRole(models.RoleAdmin), sumberDanaController.Update)
	sumberDana.Delete("/:id", middleware.RequireRole(models.RoleAdmin), sumberDanaController.Delete)
}


// setupJenisDokumenRoutes configures jenis dokumen routes
func setupJenisDokumenRoutes(api fiber.Router) {
	jenisDokumenController := controllers.NewJenisDokumenController()

	jenisDokumen := api.Group("/jenis-dokumen")

	// GET routes - accessible by Admin and above (except active for Operator)
	jenisDokumen.Get("/", middleware.RequireRole(models.RoleAdmin), jenisDokumenController.GetAll)
	jenisDokumen.Get("/active", middleware.RequireRole(models.RoleOperator), jenisDokumenController.GetAllActive)
	jenisDokumen.Get("/:id", middleware.RequireRole(models.RoleAdmin), jenisDokumenController.GetByID)

	// CUD routes - accessible by Admin and above
	jenisDokumen.Post("/", middleware.RequireRole(models.RoleAdmin), jenisDokumenController.Create)
	jenisDokumen.Put("/:id", middleware.RequireRole(models.RoleAdmin), jenisDokumenController.Update)
	jenisDokumen.Delete("/:id", middleware.RequireRole(models.RoleAdmin), jenisDokumenController.Delete)
}


// setupUserRoutes configures user management routes
func setupUserRoutes(api fiber.Router) {
	userController := controllers.NewUserController()

	users := api.Group("/users")

	// All user routes require Super Admin role
	users.Get("/", middleware.RequireRole(models.RoleSuperAdmin), userController.GetAll)
	users.Get("/:id", middleware.RequireRole(models.RoleSuperAdmin), userController.GetByID)
	users.Post("/", middleware.RequireRole(models.RoleSuperAdmin), userController.Create)
	users.Put("/:id", middleware.RequireRole(models.RoleSuperAdmin), userController.Update)
	users.Delete("/:id", middleware.RequireRole(models.RoleSuperAdmin), userController.Delete)
	users.Post("/:id/activate", middleware.RequireRole(models.RoleSuperAdmin), userController.Activate)
	users.Post("/:id/avatar", middleware.RequireRole(models.RoleSuperAdmin), userController.UploadAvatar)
	
	// Reset password - Admin and Super Admin can reset
	users.Post("/:id/reset-password", middleware.RequireRole(models.RoleAdmin), userController.ResetPassword)
}





// setupSettingRoutes configures setting routes
func setupSettingRoutes(api fiber.Router) {
	settingController := controllers.NewSettingController()

	settings := api.Group("/settings")

	// Settings routes - Super Admin only
	settings.Get("/", middleware.RequireRole(models.RoleSuperAdmin), settingController.GetAll)
	settings.Put("/", middleware.RequireRole(models.RoleSuperAdmin), settingController.Update)
	settings.Get("/countdown", settingController.GetCountdownSettings) 
	settings.Get("/branding", settingController.GetBrandingSettings) // Allow all authenticated
	settings.Post("/upload-logo", middleware.RequireRole(models.RoleSuperAdmin), settingController.UploadLogo)
	
	// Lock status routes
	settings.Get("/lock-status", settingController.GetLockStatus) // All authenticated can check
	settings.Post("/toggle-lock", middleware.RequireRole(models.RoleSuperAdmin), settingController.ToggleLock)
}


// setupPetunjukRoutes configures petunjuk routes
func setupPetunjukRoutes(api fiber.Router) {
	petunjukController := controllers.NewPetunjukController()

	petunjuk := api.Group("/petunjuk")

	// Public route for getting petunjuk by halaman (all authenticated users)
	petunjuk.Get("/halaman/:halaman", middleware.RequireRole(models.RoleOperator), petunjukController.GetByHalaman)

	// Admin routes
	petunjuk.Get("/", middleware.RequireRole(models.RoleSuperAdmin), petunjukController.GetAll)
	petunjuk.Get("/:id", middleware.RequireRole(models.RoleSuperAdmin), petunjukController.GetByID)
	petunjuk.Post("/", middleware.RequireRole(models.RoleSuperAdmin), petunjukController.Create)
	petunjuk.Post("/upload-image", middleware.RequireRole(models.RoleSuperAdmin), petunjukController.UploadImage)
	petunjuk.Put("/:id", middleware.RequireRole(models.RoleSuperAdmin), petunjukController.Update)
	petunjuk.Delete("/:id", middleware.RequireRole(models.RoleSuperAdmin), petunjukController.Delete)
}

// setupDokumenRoutes configures dokumen routes
func setupDokumenRoutes(api fiber.Router) {
	dokumenController := controllers.NewDokumenController()

	dokumen := api.Group("/dokumen")

	// All authenticated users can access dokumen (with role-based filtering)
	dokumen.Get("/", middleware.RequireRole(models.RoleOperator), dokumenController.GetAll)
	dokumen.Get("/:id", middleware.RequireRole(models.RoleOperator), dokumenController.GetByID)
	dokumen.Get("/:id/file", middleware.RequireRole(models.RoleOperator), dokumenController.GetFile)
	dokumen.Post("/", middleware.RequireRole(models.RoleOperator), dokumenController.Create)
	dokumen.Put("/:id", middleware.RequireRole(models.RoleOperator), dokumenController.Update)
	dokumen.Delete("/:id", middleware.RequireRole(models.RoleOperator), dokumenController.Delete)
}

// setupLoginContentRoutes configures login content routes
func setupLoginContentRoutes(api fiber.Router) {
	loginContentController := controllers.NewLoginContentController()

	loginContent := api.Group("/login-content")

	// Login content routes - Super Admin only
	loginContent.Get("/", middleware.RequireRole(models.RoleSuperAdmin), loginContentController.GetAll)
	loginContent.Post("/", middleware.RequireRole(models.RoleSuperAdmin), loginContentController.Create)
	loginContent.Put("/:id", middleware.RequireRole(models.RoleSuperAdmin), loginContentController.Update)
	loginContent.Delete("/:id", middleware.RequireRole(models.RoleSuperAdmin), loginContentController.Delete)
	loginContent.Post("/:id/image", middleware.RequireRole(models.RoleSuperAdmin), loginContentController.UploadImage)
}

