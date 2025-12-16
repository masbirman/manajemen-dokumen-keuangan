package repositories

import (
	"time"

	"dokumen-keuangan/app/models"
	"dokumen-keuangan/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DokumenFilter represents filter options for document queries
type DokumenFilter struct {
	UnitKerjaID *uuid.UUID
	PPTKID      *uuid.UUID
	CreatedBy   *uuid.UUID
	StartDate   *time.Time
	EndDate     *time.Time
}

// DokumenRepository handles dokumen data access
type DokumenRepository struct {
	db *gorm.DB
}

// NewDokumenRepository creates a new DokumenRepository instance
func NewDokumenRepository() *DokumenRepository {
	return &DokumenRepository{
		db: database.GetDB(),
	}
}

// Create creates a new dokumen
func (r *DokumenRepository) Create(dokumen *models.Dokumen) error {
	if dokumen.ID == uuid.Nil {
		dokumen.ID = uuid.New()
	}
	return r.db.Create(dokumen).Error
}

// FindByID finds a dokumen by ID with all relationships
func (r *DokumenRepository) FindByID(id uuid.UUID) (*models.Dokumen, error) {
	var dokumen models.Dokumen
	err := r.db.Preload("UnitKerja").
		Preload("PPTK").
		Preload("JenisDokumen").
		Preload("SumberDana").
		Preload("Creator").
		Where("id = ?", id).
		First(&dokumen).Error
	if err != nil {
		return nil, err
	}
	return &dokumen, nil
}


// GetAll retrieves all dokumen with pagination and filters
func (r *DokumenRepository) GetAll(page, pageSize int, filter *DokumenFilter) (*PaginationResult, error) {
	var dokumens []models.Dokumen
	var total int64

	query := r.db.Model(&models.Dokumen{}).
		Preload("UnitKerja").
		Preload("PPTK").
		Preload("JenisDokumen").
		Preload("SumberDana").
		Preload("Creator")

	// Apply filters
	if filter != nil {
		if filter.UnitKerjaID != nil {
			query = query.Where("unit_kerja_id = ?", *filter.UnitKerjaID)
		}
		if filter.PPTKID != nil {
			query = query.Where("pptk_id = ?", *filter.PPTKID)
		}
		if filter.CreatedBy != nil {
			query = query.Where("created_by = ?", *filter.CreatedBy)
		}
		if filter.StartDate != nil {
			query = query.Where("tanggal_dokumen >= ?", *filter.StartDate)
		}
		if filter.EndDate != nil {
			// Add 1 day to include the end date
			endDate := filter.EndDate.Add(24 * time.Hour)
			query = query.Where("tanggal_dokumen < ?", endDate)
		}
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Fetch paginated results
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&dokumens).Error; err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &PaginationResult{
		Data:       dokumens,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// GetByCreator retrieves all dokumen created by a specific user
func (r *DokumenRepository) GetByCreator(creatorID uuid.UUID, page, pageSize int) (*PaginationResult, error) {
	filter := &DokumenFilter{
		CreatedBy: &creatorID,
	}
	return r.GetAll(page, pageSize, filter)
}

// Update updates a dokumen
func (r *DokumenRepository) Update(dokumen *models.Dokumen) error {
	return r.db.Save(dokumen).Error
}

// Delete deletes a dokumen by ID
func (r *DokumenRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Dokumen{}, "id = ?", id).Error
}
