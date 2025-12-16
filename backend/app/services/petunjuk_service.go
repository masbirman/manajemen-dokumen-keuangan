package services

import (
	"errors"
	"strings"

	"dokumen-keuangan/app/models"
	"dokumen-keuangan/app/repositories"

	"github.com/google/uuid"
)

var (
	ErrPetunjukNotFound      = errors.New("petunjuk not found")
	ErrPetunjukJudulRequired = errors.New("judul is required")
	ErrPetunjukKontenRequired = errors.New("konten is required")
	ErrPetunjukHalamanRequired = errors.New("halaman is required")
)

// PetunjukService handles petunjuk business logic
type PetunjukService struct {
	repo *repositories.PetunjukRepository
}

// NewPetunjukService creates a new PetunjukService instance
func NewPetunjukService() *PetunjukService {
	return &PetunjukService{
		repo: repositories.NewPetunjukRepository(),
	}
}

// CreatePetunjukInput represents input for creating a petunjuk
type CreatePetunjukInput struct {
	Judul    string `json:"judul"`
	Konten   string `json:"konten"`
	Halaman  string `json:"halaman"`
	Urutan   int    `json:"urutan"`
	IsActive bool   `json:"is_active"`
}

// UpdatePetunjukInput represents input for updating a petunjuk
type UpdatePetunjukInput struct {
	Judul    *string `json:"judul"`
	Konten   *string `json:"konten"`
	Halaman  *string `json:"halaman"`
	Urutan   *int    `json:"urutan"`
	IsActive *bool   `json:"is_active"`
}

// Validate validates the create input
func (i *CreatePetunjukInput) Validate() map[string]string {
	errs := make(map[string]string)
	if strings.TrimSpace(i.Judul) == "" {
		errs["judul"] = "Judul is required"
	}
	if strings.TrimSpace(i.Konten) == "" {
		errs["konten"] = "Konten is required"
	}
	if strings.TrimSpace(i.Halaman) == "" {
		errs["halaman"] = "Halaman is required"
	}
	return errs
}

// Create creates a new petunjuk
func (s *PetunjukService) Create(input *CreatePetunjukInput) (*models.Petunjuk, error) {
	if validationErrors := input.Validate(); len(validationErrors) > 0 {
		if _, ok := validationErrors["judul"]; ok {
			return nil, ErrPetunjukJudulRequired
		}
		if _, ok := validationErrors["konten"]; ok {
			return nil, ErrPetunjukKontenRequired
		}
		if _, ok := validationErrors["halaman"]; ok {
			return nil, ErrPetunjukHalamanRequired
		}
	}

	petunjuk := &models.Petunjuk{
		ID:       uuid.New(),
		Judul:    strings.TrimSpace(input.Judul),
		Konten:   strings.TrimSpace(input.Konten),
		Halaman:  strings.TrimSpace(input.Halaman),
		Urutan:   input.Urutan,
		IsActive: input.IsActive,
	}

	if err := s.repo.Create(petunjuk); err != nil {
		return nil, err
	}

	return petunjuk, nil
}

// GetByID retrieves a petunjuk by ID
func (s *PetunjukService) GetByID(id uuid.UUID) (*models.Petunjuk, error) {
	petunjuk, err := s.repo.FindByID(id)
	if err != nil {
		return nil, ErrPetunjukNotFound
	}
	return petunjuk, nil
}

// GetAll retrieves all petunjuk with pagination
func (s *PetunjukService) GetAll(page, pageSize int, search string) (*repositories.PaginationResult, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return s.repo.GetAll(page, pageSize, search)
}

// GetByHalaman retrieves all active petunjuk for a specific page
func (s *PetunjukService) GetByHalaman(halaman string) ([]models.Petunjuk, error) {
	return s.repo.GetByHalaman(halaman)
}

// Update updates a petunjuk
func (s *PetunjukService) Update(id uuid.UUID, input *UpdatePetunjukInput) (*models.Petunjuk, error) {
	petunjuk, err := s.repo.FindByID(id)
	if err != nil {
		return nil, ErrPetunjukNotFound
	}

	if input.Judul != nil {
		petunjuk.Judul = strings.TrimSpace(*input.Judul)
	}
	if input.Konten != nil {
		petunjuk.Konten = strings.TrimSpace(*input.Konten)
	}
	if input.Halaman != nil {
		petunjuk.Halaman = strings.TrimSpace(*input.Halaman)
	}
	if input.Urutan != nil {
		petunjuk.Urutan = *input.Urutan
	}
	if input.IsActive != nil {
		petunjuk.IsActive = *input.IsActive
	}

	if err := s.repo.Update(petunjuk); err != nil {
		return nil, err
	}

	return petunjuk, nil
}

// Delete deletes a petunjuk
func (s *PetunjukService) Delete(id uuid.UUID) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return ErrPetunjukNotFound
	}
	return s.repo.Delete(id)
}
