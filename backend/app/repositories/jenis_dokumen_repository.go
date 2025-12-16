package repositories

import (
	"dokumen-keuangan/app/models"
	"dokumen-keuangan/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// JenisDokumenRepository handles jenis dokumen data access
type JenisDokumenRepository struct {
	db *gorm.DB
}

// NewJenisDokumenRepository creates a new JenisDokumenRepository instance
func NewJenisDokumenRepository() *JenisDokumenRepository {
	return &JenisDokumenRepository{
		db: database.GetDB(),
	}
}

// Create creates a new jenis dokumen
func (r *JenisDokumenRepository) Create(jenisDokumen *models.JenisDokumen) error {
	if jenisDokumen.ID == uuid.Nil {
		jenisDokumen.ID = uuid.New()
	}
	return r.db.Create(jenisDokumen).Error
}

// FindByID finds a jenis dokumen by ID
func (r *JenisDokumenRepository) FindByID(id uuid.UUID) (*models.JenisDokumen, error) {
	var jenisDokumen models.JenisDokumen
	err := r.db.Where("id = ?", id).First(&jenisDokumen).Error
	if err != nil {
		return nil, err
	}
	return &jenisDokumen, nil
}

// FindByKode finds a jenis dokumen by kode
func (r *JenisDokumenRepository) FindByKode(kode string) (*models.JenisDokumen, error) {
	var jenisDokumen models.JenisDokumen
	err := r.db.Where("kode = ?", kode).First(&jenisDokumen).Error
	if err != nil {
		return nil, err
	}
	return &jenisDokumen, nil
}


// GetAll retrieves all jenis dokumen with pagination
func (r *JenisDokumenRepository) GetAll(page, pageSize int, search string) (*PaginationResult, error) {
	var jenisDokumens []models.JenisDokumen
	var total int64

	query := r.db.Model(&models.JenisDokumen{})

	// Apply search filter if provided
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("kode ILIKE ? OR nama ILIKE ?", searchPattern, searchPattern)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Fetch paginated results
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&jenisDokumens).Error; err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &PaginationResult{
		Data:       jenisDokumens,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// GetAllActive retrieves all active jenis dokumen (for dropdowns)
func (r *JenisDokumenRepository) GetAllActive() ([]models.JenisDokumen, error) {
	var jenisDokumens []models.JenisDokumen
	err := r.db.Where("is_active = ?", true).Order("nama ASC").Find(&jenisDokumens).Error
	if err != nil {
		return nil, err
	}
	return jenisDokumens, nil
}

// Update updates a jenis dokumen
func (r *JenisDokumenRepository) Update(jenisDokumen *models.JenisDokumen) error {
	return r.db.Save(jenisDokumen).Error
}

// Delete deletes a jenis dokumen by ID
func (r *JenisDokumenRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.JenisDokumen{}, "id = ?", id).Error
}

// ExistsByKode checks if a jenis dokumen with the given kode exists (excluding a specific ID)
func (r *JenisDokumenRepository) ExistsByKode(kode string, excludeID *uuid.UUID) (bool, error) {
	var count int64
	query := r.db.Model(&models.JenisDokumen{}).Where("kode = ?", kode)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// CountDocumentsByJenisDokumenID counts documents referencing this jenis dokumen
func (r *JenisDokumenRepository) CountDocumentsByJenisDokumenID(id uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&models.Dokumen{}).Where("jenis_dokumen_id = ?", id).Count(&count).Error
	return count, err
}
