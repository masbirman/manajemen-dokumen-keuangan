package repositories

import (
	"dokumen-keuangan/app/models"
	"dokumen-keuangan/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PaginationResult holds paginated results
type PaginationResult struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

// UnitKerjaRepository handles unit kerja data access
type UnitKerjaRepository struct {
	db *gorm.DB
}

// NewUnitKerjaRepository creates a new UnitKerjaRepository instance
func NewUnitKerjaRepository() *UnitKerjaRepository {
	return &UnitKerjaRepository{
		db: database.GetDB(),
	}
}

// Create creates a new unit kerja
func (r *UnitKerjaRepository) Create(unitKerja *models.UnitKerja) error {
	if unitKerja.ID == uuid.Nil {
		unitKerja.ID = uuid.New()
	}
	return r.db.Create(unitKerja).Error
}

// FindByID finds a unit kerja by ID
func (r *UnitKerjaRepository) FindByID(id uuid.UUID) (*models.UnitKerja, error) {
	var unitKerja models.UnitKerja
	err := r.db.Where("id = ?", id).First(&unitKerja).Error
	if err != nil {
		return nil, err
	}
	return &unitKerja, nil
}

// FindByKode finds a unit kerja by kode
func (r *UnitKerjaRepository) FindByKode(kode string) (*models.UnitKerja, error) {
	var unitKerja models.UnitKerja
	err := r.db.Where("kode = ?", kode).First(&unitKerja).Error
	if err != nil {
		return nil, err
	}
	return &unitKerja, nil
}


// GetAll retrieves all unit kerja with pagination
func (r *UnitKerjaRepository) GetAll(page, pageSize int, search string) (*PaginationResult, error) {
	var unitKerjas []models.UnitKerja
	var total int64

	query := r.db.Model(&models.UnitKerja{})

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
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&unitKerjas).Error; err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &PaginationResult{
		Data:       unitKerjas,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// GetAllActive retrieves all active unit kerja (for dropdowns)
func (r *UnitKerjaRepository) GetAllActive() ([]models.UnitKerja, error) {
	var unitKerjas []models.UnitKerja
	err := r.db.Where("is_active = ?", true).Order("nama ASC").Find(&unitKerjas).Error
	if err != nil {
		return nil, err
	}
	return unitKerjas, nil
}

// Update updates a unit kerja
func (r *UnitKerjaRepository) Update(unitKerja *models.UnitKerja) error {
	return r.db.Save(unitKerja).Error
}

// Delete deletes a unit kerja by ID
func (r *UnitKerjaRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.UnitKerja{}, "id = ?", id).Error
}

// ExistsByKode checks if a unit kerja with the given kode exists (excluding a specific ID)
func (r *UnitKerjaRepository) ExistsByKode(kode string, excludeID *uuid.UUID) (bool, error) {
	var count int64
	query := r.db.Model(&models.UnitKerja{}).Where("kode = ?", kode)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
