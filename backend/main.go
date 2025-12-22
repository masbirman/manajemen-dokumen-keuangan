package main

import (
	"log"

	"dokumen-keuangan/config"
	"dokumen-keuangan/database"
	"dokumen-keuangan/database/seeders"
	"dokumen-keuangan/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Connect to database
	if err := database.Connect(cfg); err != nil {
		log.Printf("Warning: Database connection failed: %v", err)
		log.Println("API will start without database connection")
	} else {
		// Run migrations
		if err := database.RunMigrations(); err != nil {
			log.Printf("Warning: Migration failed: %v", err)
		}

		// Run seeders
		if err := seeders.RunSeeders(); err != nil {
			log.Printf("Warning: Seeder failed: %v", err)
		}
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:   cfg.AppName,
		BodyLimit: 50 * 1024 * 1024, // 50MB limit for large content (images in rich text)
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173, http://localhost:80, https://dokumen.keudisdiksulteng.web.id, http://dokumen.keudisdiksulteng.web.id",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Requested-With, X-Tahun-Anggaran",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS, PATCH",
		AllowCredentials: true,
	}))

	// Health check endpoint
	app.Get("/api/health", func(c *fiber.Ctx) error {
		dbStatus := "connected"
		if database.GetDB() == nil {
			dbStatus = "disconnected"
		} else {
			sqlDB, err := database.GetDB().DB()
			if err != nil || sqlDB.Ping() != nil {
				dbStatus = "disconnected"
			}
		}

		return c.JSON(fiber.Map{
			"status":   "ok",
			"message":  "Dokumen Keuangan API is running",
			"database": dbStatus,
		})
	})

	// Setup routes
	routes.SetupRoutes(app, cfg)

	// Start server
	addr := cfg.AppHost + ":" + cfg.AppPort
	log.Printf("Server starting on %s", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
