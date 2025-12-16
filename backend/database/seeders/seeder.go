package seeders

import (
	"log"

	"dokumen-keuangan/app/models"
	"dokumen-keuangan/database"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// RunSeeders executes all database seeders
func RunSeeders() error {
	if database.DB == nil {
		log.Println("Database not connected, skipping seeders")
		return nil
	}

	if err := seedSuperAdmin(); err != nil {
		return err
	}

	log.Println("All seeders completed successfully")
	return nil
}

// seedSuperAdmin creates the default super admin user if not exists
func seedSuperAdmin() error {
	var count int64
	database.DB.Model(&models.User{}).Where("role = ?", models.RoleSuperAdmin).Count(&count)

	if count > 0 {
		log.Println("Super Admin already exists, skipping seeder")
		return nil
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	superAdmin := models.User{
		ID:       uuid.New(),
		Username: "superadmin",
		Password: string(hashedPassword),
		Name:     "Super Administrator",
		Role:     models.RoleSuperAdmin,
		IsActive: true,
	}

	if err := database.DB.Create(&superAdmin).Error; err != nil {
		return err
	}

	log.Println("Super Admin user created successfully")
	log.Println("Username: superadmin")
	log.Println("Password: admin123")
	return nil
}
