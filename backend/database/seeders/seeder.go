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

// seedSuperAdmin creates or updates the default super admin user
func seedSuperAdmin() error {
	var existingUser models.User
	result := database.DB.Where("username = ?", "superadmin").First(&existingUser)

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if result.Error == nil {
		// User exists, update password
		existingUser.Password = string(hashedPassword)
		if err := database.DB.Save(&existingUser).Error; err != nil {
			return err
		}
		log.Println("Super Admin password updated successfully")
		log.Println("Username: superadmin")
		log.Println("Password: admin123")
		return nil
	}

	// User doesn't exist, create new
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
