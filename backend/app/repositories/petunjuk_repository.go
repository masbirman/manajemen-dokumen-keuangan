package repositories

import (
	"dokumen-keuangan/app/models"
	"dokumen-keuangan/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PetunjukRepository handles petunjuk data access
type PetunjukRepository struct {
	db *gorm.DB
}

// NewPetunjukRepository creates a new PetunjukRepository instance
func NewPetunjukRepository() *PetunjukRepository {
	return &PetunjukRepository{
		db: database.GetDB(),
	}
}

// Create creates a new petunjuk
func (r *PetunjukRepository) Create(petunjuk *models.Petunjuk) error {
	if petunjuk.ID == uuid.Nil {
		petunjuk.ID = uuid.New()
	}
	return r.db.Create(petunjuk).Error
}

// FindByID finds a petunjuk by ID
func (r *PetunjukRepository) FindByID(id uuid.UUID) (*models.Petunjuk, error) {
	var petunjuk models.Petunjuk
	err := r.db.Where("id = ?", id).First(&petunjuk).Error
	if err != nil {
		return nil, err
	}
	return &petunjuk, nil
}

// GetAll retrieves all petunjuk with pagination
func (r *PetunjukRepository) GetAll(page, pageSize int, search string) (*PaginationResult, error) {
	var petunjuks []models.Petunjuk
	var total int64

	query := r.db.Model(&models.Petunjuk{})

	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("judul ILIKE ? OR konten ILIKE ?", searchPattern, searchPattern)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (page - 1) * pageSize

	if err := query.Order("halaman ASC, urutan ASC").Offset(offset).Limit(pageSize).Find(&petunjuks).Error; err != nil {
		return nil, err
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &PaginationResult{
		Data:       petunjuks,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// GetByHalaman retrieves all active petunjuk for a specific page
func (r *PetunjukRepository) GetByHalaman(halaman string) ([]models.Petunjuk, error) {
	var petunjuks []models.Petunjuk
	err := r.db.Where("halaman = ? AND is_active = ?", halaman, true).
		Order("urutan ASC").
		Find(&petunjuks).Error
	if err != nil {
		return nil, err
	}
	return petunjuks, nil
}

// Update updates a petunjuk
func (r *PetunjukRepository) Update(petunjuk *models.Petunjuk) error {
	return r.db.Save(petunjuk).Error
}

// Delete deletes a petunjuk by ID
func (r *PetunjukRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Petunjuk{}, "id = ?", id).Error
}
