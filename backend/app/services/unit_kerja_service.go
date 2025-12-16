package services

import (
	"errors"
	"strings"

	"dokumen-keuangan/app/models"
	"dokumen-keuangan/app/repositories"

	"github.com/google/uuid"
)

var (
	ErrUnitKerjaNotFound     = errors.New("unit kerja not found")
	ErrUnitKerjaKodeExists   = errors.New("unit kerja with this kode already exists")
	ErrUnitKerjaKodeRequired = errors.New("kode is required")
	ErrUnitKerjaNamaRequired = errors.New("nama is required")
)

// UnitKerjaService handles unit kerja business logic
type UnitKerjaService struct {
	repo *repositories.UnitKerjaRepository
}

// NewUnitKerjaService creates a new UnitKerjaService instance
func NewUnitKerjaService() *UnitKerjaService {
	return &UnitKerjaService{
		repo: repositories.NewUnitKerjaRepository(),
	}
}

// CreateInput represents input for creating a unit kerja
type CreateUnitKerjaInput struct {
	Kode string `json:"kode"`
	Nama string `json:"nama"`
}

// UpdateInput represents input for updating a unit kerja
type UpdateUnitKerjaInput struct {
	Kode     *string `json:"kode"`
	Nama     *string `json:"nama"`
	IsActive *bool   `json:"is_active"`
}

// Validate validates the create input
func (i *CreateUnitKerjaInput) Validate() map[string]string {
	errors := make(map[string]string)
	
	if strings.TrimSpace(i.Kode) == "" {
		errors["kode"] = "Kode is required"
	}
	if strings.TrimSpace(i.Nama) == "" {
		errors["nama"] = "Nama is required"
	}
	
	return errors
}


// Create creates a new unit kerja
func (s *UnitKerjaService) Create(input *CreateUnitKerjaInput) (*models.UnitKerja, error) {
	// Validate input
	if validationErrors := input.Validate(); len(validationErrors) > 0 {
		if _, ok := validationErrors["kode"]; ok {
			return nil, ErrUnitKerjaKodeRequired
		}
		if _, ok := validationErrors["nama"]; ok {
			return nil, ErrUnitKerjaNamaRequired
		}
	}

	// Check if kode already exists
	exists, err := s.repo.ExistsByKode(strings.TrimSpace(input.Kode), nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUnitKerjaKodeExists
	}

	// Create unit kerja
	unitKerja := &models.UnitKerja{
		ID:       uuid.New(),
		Kode:     strings.TrimSpace(input.Kode),
		Nama:     strings.TrimSpace(input.Nama),
		IsActive: true,
	}

	if err := s.repo.Create(unitKerja); err != nil {
		return nil, err
	}

	return unitKerja, nil
}

// GetByID retrieves a unit kerja by ID
func (s *UnitKerjaService) GetByID(id uuid.UUID) (*models.UnitKerja, error) {
	unitKerja, err := s.repo.FindByID(id)
	if err != nil {
		return nil, ErrUnitKerjaNotFound
	}
	return unitKerja, nil
}

// GetAll retrieves all unit kerja with pagination
func (s *UnitKerjaService) GetAll(page, pageSize int, search string) (*repositories.PaginationResult, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	return s.repo.GetAll(page, pageSize, search)
}

// GetAllActive retrieves all active unit kerja (for dropdowns)
func (s *UnitKerjaService) GetAllActive() ([]models.UnitKerja, error) {
	return s.repo.GetAllActive()
}

// Update updates a unit kerja
func (s *UnitKerjaService) Update(id uuid.UUID, input *UpdateUnitKerjaInput) (*models.UnitKerja, error) {
	// Find existing unit kerja
	unitKerja, err := s.repo.FindByID(id)
	if err != nil {
		return nil, ErrUnitKerjaNotFound
	}

	// Update fields if provided
	if input.Kode != nil {
		kode := strings.TrimSpace(*input.Kode)
		if kode == "" {
			return nil, ErrUnitKerjaKodeRequired
		}
		// Check if new kode already exists (excluding current record)
		exists, err := s.repo.ExistsByKode(kode, &id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrUnitKerjaKodeExists
		}
		unitKerja.Kode = kode
	}

	if input.Nama != nil {
		nama := strings.TrimSpace(*input.Nama)
		if nama == "" {
			return nil, ErrUnitKerjaNamaRequired
		}
		unitKerja.Nama = nama
	}

	if input.IsActive != nil {
		unitKerja.IsActive = *input.IsActive
	}

	if err := s.repo.Update(unitKerja); err != nil {
		return nil, err
	}

	return unitKerja, nil
}

// Delete deletes a unit kerja by ID
func (s *UnitKerjaService) Delete(id uuid.UUID) error {
	// Check if unit kerja exists
	_, err := s.repo.FindByID(id)
	if err != nil {
		return ErrUnitKerjaNotFound
	}

	return s.repo.Delete(id)
}
