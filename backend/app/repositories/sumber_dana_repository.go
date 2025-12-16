package repositories

import (
	"dokumen-keuangan/app/models"
	"dokumen-keuangan/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SumberDanaRepository handles sumber dana data access
type SumberDanaRepository struct {
	db *gorm.DB
}

// NewSumberDanaRepository creates a new SumberDanaRepository instance
func NewSumberDanaRepository() *SumberDanaRepository {
	return &SumberDanaRepository{
		db: database.GetDB(),
	}
}

// Create creates a new sumber dana
func (r *SumberDanaRepository) Create(sumberDana *models.SumberDana) error {
	if sumberDana.ID == uuid.Nil {
		sumberDana.ID = uuid.New()
	}
	return r.db.Create(sumberDana).Error
}

// FindByID finds a sumber dana by ID
func (r *SumberDanaRepository) FindByID(id uuid.UUID) (*models.SumberDana, error) {
	var sumberDana models.SumberDana
	err := r.db.Where("id = ?", id).First(&sumberDana).Error
	if err != nil {
		return nil, err
	}
	return &sumberDana, nil
}

// FindByKode finds a sumber dana by kode
func (r *SumberDanaRepository) FindByKode(kode string) (*models.SumberDana, error) {
	var sumberDana models.SumberDana
	err := r.db.Where("kode = ?", kode).First(&sumberDana).Error
	if err != nil {
		return nil, err
	}
	return &sumberDana, nil
}


// GetAll retrieves all sumber dana with pagination
func (r *SumberDanaRepository) GetAll(page, pageSize int, search string) (*PaginationResult, error) {
	var sumberDanas []models.SumberDana
	var total int64

	query := r.db.Model(&models.SumberDana{})

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
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&sumberDanas).Error; err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &PaginationResult{
		Data:       sumberDanas,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// GetAllActive retrieves all active sumber dana (for dropdowns)
func (r *SumberDanaRepository) GetAllActive() ([]models.SumberDana, error) {
	var sumberDanas []models.SumberDana
	err := r.db.Where("is_active = ?", true).Order("nama ASC").Find(&sumberDanas).Error
	if err != nil {
		return nil, err
	}
	return sumberDanas, nil
}

// Update updates a sumber dana
func (r *SumberDanaRepository) Update(sumberDana *models.SumberDana) error {
	return r.db.Save(sumberDana).Error
}

// Delete deletes a sumber dana by ID
func (r *SumberDanaRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.SumberDana{}, "id = ?", id).Error
}

// ExistsByKode checks if a sumber dana with the given kode exists (excluding a specific ID)
func (r *SumberDanaRepository) ExistsByKode(kode string, excludeID *uuid.UUID) (bool, error) {
	var count int64
	query := r.db.Model(&models.SumberDana{}).Where("kode = ?", kode)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// CountDocumentsBySumberDanaID counts documents referencing this sumber dana
func (r *SumberDanaRepository) CountDocumentsBySumberDanaID(id uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&models.Dokumen{}).Where("sumber_dana_id = ?", id).Count(&count).Error
	return count, err
}
