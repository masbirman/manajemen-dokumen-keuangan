package services

import (
	"errors"
	"fmt"
	"strings"

	"dokumen-keuangan/app/models"
	"dokumen-keuangan/app/repositories"

	"github.com/google/uuid"
)

var (
	ErrSumberDanaNotFound     = errors.New("sumber dana not found")
	ErrSumberDanaKodeExists   = errors.New("sumber dana with this kode already exists")
	ErrSumberDanaKodeRequired = errors.New("kode is required")
	ErrSumberDanaNamaRequired = errors.New("nama is required")
)

// ErrSumberDanaReferenced is returned when trying to delete a sumber dana that is referenced by documents
type ErrSumberDanaReferenced struct {
	Count int64
}

func (e *ErrSumberDanaReferenced) Error() string {
	return fmt.Sprintf("cannot delete: referenced by %d documents", e.Count)
}

// SumberDanaService handles sumber dana business logic
type SumberDanaService struct {
	repo *repositories.SumberDanaRepository
}

// NewSumberDanaService creates a new SumberDanaService instance
func NewSumberDanaService() *SumberDanaService {
	return &SumberDanaService{
		repo: repositories.NewSumberDanaRepository(),
	}
}

// CreateSumberDanaInput represents input for creating a sumber dana
type CreateSumberDanaInput struct {
	Kode string `json:"kode"`
	Nama string `json:"nama"`
}

// UpdateSumberDanaInput represents input for updating a sumber dana
type UpdateSumberDanaInput struct {
	Kode     *string `json:"kode"`
	Nama     *string `json:"nama"`
	IsActive *bool   `json:"is_active"`
}


// Validate validates the create input
func (i *CreateSumberDanaInput) Validate() map[string]string {
	errors := make(map[string]string)

	if strings.TrimSpace(i.Kode) == "" {
		errors["kode"] = "Kode is required"
	}
	if strings.TrimSpace(i.Nama) == "" {
		errors["nama"] = "Nama is required"
	}

	return errors
}

// Create creates a new sumber dana
func (s *SumberDanaService) Create(input *CreateSumberDanaInput) (*models.SumberDana, error) {
	// Validate input
	if validationErrors := input.Validate(); len(validationErrors) > 0 {
		if _, ok := validationErrors["kode"]; ok {
			return nil, ErrSumberDanaKodeRequired
		}
		if _, ok := validationErrors["nama"]; ok {
			return nil, ErrSumberDanaNamaRequired
		}
	}

	// Check if kode already exists
	exists, err := s.repo.ExistsByKode(strings.TrimSpace(input.Kode), nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrSumberDanaKodeExists
	}

	// Create sumber dana
	sumberDana := &models.SumberDana{
		ID:       uuid.New(),
		Kode:     strings.TrimSpace(input.Kode),
		Nama:     strings.TrimSpace(input.Nama),
		IsActive: true,
	}

	if err := s.repo.Create(sumberDana); err != nil {
		return nil, err
	}

	return sumberDana, nil
}

// GetByID retrieves a sumber dana by ID
func (s *SumberDanaService) GetByID(id uuid.UUID) (*models.SumberDana, error) {
	sumberDana, err := s.repo.FindByID(id)
	if err != nil {
		return nil, ErrSumberDanaNotFound
	}
	return sumberDana, nil
}

// GetAll retrieves all sumber dana with pagination
func (s *SumberDanaService) GetAll(page, pageSize int, search string) (*repositories.PaginationResult, error) {
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

// GetAllActive retrieves all active sumber dana (for dropdowns)
func (s *SumberDanaService) GetAllActive() ([]models.SumberDana, error) {
	return s.repo.GetAllActive()
}


// Update updates a sumber dana
func (s *SumberDanaService) Update(id uuid.UUID, input *UpdateSumberDanaInput) (*models.SumberDana, error) {
	// Find existing sumber dana
	sumberDana, err := s.repo.FindByID(id)
	if err != nil {
		return nil, ErrSumberDanaNotFound
	}

	// Update fields if provided
	if input.Kode != nil {
		kode := strings.TrimSpace(*input.Kode)
		if kode == "" {
			return nil, ErrSumberDanaKodeRequired
		}
		// Check if new kode already exists (excluding current record)
		exists, err := s.repo.ExistsByKode(kode, &id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrSumberDanaKodeExists
		}
		sumberDana.Kode = kode
	}

	if input.Nama != nil {
		nama := strings.TrimSpace(*input.Nama)
		if nama == "" {
			return nil, ErrSumberDanaNamaRequired
		}
		sumberDana.Nama = nama
	}

	if input.IsActive != nil {
		sumberDana.IsActive = *input.IsActive
	}

	if err := s.repo.Update(sumberDana); err != nil {
		return nil, err
	}

	return sumberDana, nil
}

// Delete deletes a sumber dana by ID
// Returns ErrSumberDanaReferenced if the sumber dana is referenced by documents
func (s *SumberDanaService) Delete(id uuid.UUID) error {
	// Check if sumber dana exists
	_, err := s.repo.FindByID(id)
	if err != nil {
		return ErrSumberDanaNotFound
	}

	// Check referential integrity - count documents referencing this sumber dana
	count, err := s.repo.CountDocumentsBySumberDanaID(id)
	if err != nil {
		return err
	}
	if count > 0 {
		return &ErrSumberDanaReferenced{Count: count}
	}

	return s.repo.Delete(id)
}
