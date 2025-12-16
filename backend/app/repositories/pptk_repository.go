package repositories

import (
	"dokumen-keuangan/app/models"
	"dokumen-keuangan/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PPTKRepository handles PPTK data access
type PPTKRepository struct {
	db *gorm.DB
}

// NewPPTKRepository creates a new PPTKRepository instance
func NewPPTKRepository() *PPTKRepository {
	return &PPTKRepository{
		db: database.GetDB(),
	}
}

// Create creates a new PPTK
func (r *PPTKRepository) Create(pptk *models.PPTK) error {
	if pptk.ID == uuid.Nil {
		pptk.ID = uuid.New()
	}
	return r.db.Create(pptk).Error
}

// FindByID finds a PPTK by ID with UnitKerja preloaded
func (r *PPTKRepository) FindByID(id uuid.UUID) (*models.PPTK, error) {
	var pptk models.PPTK
	err := r.db.Preload("UnitKerja").Where("id = ?", id).First(&pptk).Error
	if err != nil {
		return nil, err
	}
	return &pptk, nil
}

// FindByNIP finds a PPTK by NIP
func (r *PPTKRepository) FindByNIP(nip string) (*models.PPTK, error) {
	var pptk models.PPTK
	err := r.db.Preload("UnitKerja").Where("nip = ?", nip).First(&pptk).Error
	if err != nil {
		return nil, err
	}
	return &pptk, nil
}


// GetAll retrieves all PPTK with pagination
func (r *PPTKRepository) GetAll(page, pageSize int, search string) (*PaginationResult, error) {
	return r.GetAllWithFilter(page, pageSize, search, nil)
}

// GetAllWithFilter retrieves all PPTK with pagination and unit kerja filter
func (r *PPTKRepository) GetAllWithFilter(page, pageSize int, search string, unitKerjaID *uuid.UUID) (*PaginationResult, error) {
	var pptks []models.PPTK
	var total int64

	query := r.db.Model(&models.PPTK{}).Preload("UnitKerja")

	// Apply search filter if provided
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("nip ILIKE ? OR nama ILIKE ?", searchPattern, searchPattern)
	}

	// Apply unit kerja filter if provided
	if unitKerjaID != nil {
		query = query.Where("unit_kerja_id = ?", *unitKerjaID)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Fetch paginated results
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&pptks).Error; err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &PaginationResult{
		Data:       pptks,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// GetAllActive retrieves all active PPTK (for dropdowns)
func (r *PPTKRepository) GetAllActive() ([]models.PPTK, error) {
	var pptks []models.PPTK
	err := r.db.Preload("UnitKerja").Where("is_active = ?", true).Order("nama ASC").Find(&pptks).Error
	if err != nil {
		return nil, err
	}
	return pptks, nil
}

// GetByUnitKerja retrieves all PPTK for a specific Unit Kerja
func (r *PPTKRepository) GetByUnitKerja(unitKerjaID uuid.UUID) ([]models.PPTK, error) {
	var pptks []models.PPTK
	err := r.db.Preload("UnitKerja").Where("unit_kerja_id = ? AND is_active = ?", unitKerjaID, true).Order("nama ASC").Find(&pptks).Error
	if err != nil {
		return nil, err
	}
	return pptks, nil
}

// Update updates a PPTK
func (r *PPTKRepository) Update(pptk *models.PPTK) error {
	return r.db.Save(pptk).Error
}

// Delete deletes a PPTK by ID
func (r *PPTKRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.PPTK{}, "id = ?", id).Error
}

// ExistsByNIP checks if a PPTK with the given NIP exists (excluding a specific ID)
func (r *PPTKRepository) ExistsByNIP(nip string, excludeID *uuid.UUID) (bool, error) {
	var count int64
	query := r.db.Model(&models.PPTK{}).Where("nip = ?", nip)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
