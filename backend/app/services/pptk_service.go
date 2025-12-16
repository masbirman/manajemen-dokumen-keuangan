package services

import (
	"errors"
	"strings"

	"dokumen-keuangan/app/models"
	"dokumen-keuangan/app/repositories"

	"github.com/google/uuid"
)

var (
	ErrPPTKNotFound        = errors.New("PPTK not found")
	ErrPPTKNIPExists       = errors.New("PPTK with this NIP already exists")
	ErrPPTKNIPRequired     = errors.New("NIP is required")
	ErrPPTKNamaRequired    = errors.New("Nama is required")
	ErrPPTKUnitKerjaRequired = errors.New("Unit Kerja is required")
	ErrPPTKUnitKerjaNotFound = errors.New("Unit Kerja not found")
)

// PPTKService handles PPTK business logic
type PPTKService struct {
	repo           *repositories.PPTKRepository
	unitKerjaRepo  *repositories.UnitKerjaRepository
}

// NewPPTKService creates a new PPTKService instance
func NewPPTKService() *PPTKService {
	return &PPTKService{
		repo:          repositories.NewPPTKRepository(),
		unitKerjaRepo: repositories.NewUnitKerjaRepository(),
	}
}

// CreatePPTKInput represents input for creating a PPTK
type CreatePPTKInput struct {
	NIP         string    `json:"nip"`
	Nama        string    `json:"nama"`
	UnitKerjaID uuid.UUID `json:"unit_kerja_id"`
}

// UpdatePPTKInput represents input for updating a PPTK
type UpdatePPTKInput struct {
	NIP         *string    `json:"nip"`
	Nama        *string    `json:"nama"`
	UnitKerjaID *uuid.UUID `json:"unit_kerja_id"`
	AvatarPath  *string    `json:"avatar_path"`
	IsActive    *bool      `json:"is_active"`
}


// Validate validates the create input
func (i *CreatePPTKInput) Validate() map[string]string {
	errors := make(map[string]string)
	
	if strings.TrimSpace(i.NIP) == "" {
		errors["nip"] = "NIP is required"
	}
	if strings.TrimSpace(i.Nama) == "" {
		errors["nama"] = "Nama is required"
	}
	if i.UnitKerjaID == uuid.Nil {
		errors["unit_kerja_id"] = "Unit Kerja is required"
	}
	
	return errors
}

// Create creates a new PPTK
func (s *PPTKService) Create(input *CreatePPTKInput) (*models.PPTK, error) {
	// Validate input
	if validationErrors := input.Validate(); len(validationErrors) > 0 {
		if _, ok := validationErrors["nip"]; ok {
			return nil, ErrPPTKNIPRequired
		}
		if _, ok := validationErrors["nama"]; ok {
			return nil, ErrPPTKNamaRequired
		}
		if _, ok := validationErrors["unit_kerja_id"]; ok {
			return nil, ErrPPTKUnitKerjaRequired
		}
	}

	// Check if NIP already exists
	exists, err := s.repo.ExistsByNIP(strings.TrimSpace(input.NIP), nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrPPTKNIPExists
	}

	// Verify Unit Kerja exists
	_, err = s.unitKerjaRepo.FindByID(input.UnitKerjaID)
	if err != nil {
		return nil, ErrPPTKUnitKerjaNotFound
	}

	// Create PPTK
	pptk := &models.PPTK{
		ID:          uuid.New(),
		NIP:         strings.TrimSpace(input.NIP),
		Nama:        strings.TrimSpace(input.Nama),
		UnitKerjaID: input.UnitKerjaID,
		IsActive:    true,
	}

	if err := s.repo.Create(pptk); err != nil {
		return nil, err
	}

	// Reload with UnitKerja relationship
	return s.repo.FindByID(pptk.ID)
}

// GetByID retrieves a PPTK by ID
func (s *PPTKService) GetByID(id uuid.UUID) (*models.PPTK, error) {
	pptk, err := s.repo.FindByID(id)
	if err != nil {
		return nil, ErrPPTKNotFound
	}
	return pptk, nil
}

// GetAll retrieves all PPTK with pagination
func (s *PPTKService) GetAll(page, pageSize int, search string) (*repositories.PaginationResult, error) {
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

// GetAllWithFilter retrieves all PPTK with pagination and unit kerja filter
func (s *PPTKService) GetAllWithFilter(page, pageSize int, search, unitKerjaID string) (*repositories.PaginationResult, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	var unitKerjaUUID *uuid.UUID
	if unitKerjaID != "" {
		parsed, err := uuid.Parse(unitKerjaID)
		if err == nil {
			unitKerjaUUID = &parsed
		}
	}

	return s.repo.GetAllWithFilter(page, pageSize, search, unitKerjaUUID)
}

// GetAllActive retrieves all active PPTK (for dropdowns)
func (s *PPTKService) GetAllActive() ([]models.PPTK, error) {
	return s.repo.GetAllActive()
}

// GetByUnitKerja retrieves all PPTK for a specific Unit Kerja
func (s *PPTKService) GetByUnitKerja(unitKerjaID uuid.UUID) ([]models.PPTK, error) {
	// Verify Unit Kerja exists
	_, err := s.unitKerjaRepo.FindByID(unitKerjaID)
	if err != nil {
		return nil, ErrPPTKUnitKerjaNotFound
	}

	return s.repo.GetByUnitKerja(unitKerjaID)
}


// Update updates a PPTK
func (s *PPTKService) Update(id uuid.UUID, input *UpdatePPTKInput) (*models.PPTK, error) {
	// Find existing PPTK
	pptk, err := s.repo.FindByID(id)
	if err != nil {
		return nil, ErrPPTKNotFound
	}

	// Update fields if provided
	if input.NIP != nil {
		nip := strings.TrimSpace(*input.NIP)
		if nip == "" {
			return nil, ErrPPTKNIPRequired
		}
		// Check if new NIP already exists (excluding current record)
		exists, err := s.repo.ExistsByNIP(nip, &id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrPPTKNIPExists
		}
		pptk.NIP = nip
	}

	if input.Nama != nil {
		nama := strings.TrimSpace(*input.Nama)
		if nama == "" {
			return nil, ErrPPTKNamaRequired
		}
		pptk.Nama = nama
	}

	if input.UnitKerjaID != nil {
		// Verify Unit Kerja exists
		_, err := s.unitKerjaRepo.FindByID(*input.UnitKerjaID)
		if err != nil {
			return nil, ErrPPTKUnitKerjaNotFound
		}
		pptk.UnitKerjaID = *input.UnitKerjaID
	}

	if input.AvatarPath != nil {
		pptk.AvatarPath = input.AvatarPath
	}

	if input.IsActive != nil {
		pptk.IsActive = *input.IsActive
	}

	if err := s.repo.Update(pptk); err != nil {
		return nil, err
	}

	// Reload with UnitKerja relationship
	return s.repo.FindByID(pptk.ID)
}

// Delete deletes a PPTK by ID
func (s *PPTKService) Delete(id uuid.UUID) error {
	// Check if PPTK exists
	_, err := s.repo.FindByID(id)
	if err != nil {
		return ErrPPTKNotFound
	}

	return s.repo.Delete(id)
}
