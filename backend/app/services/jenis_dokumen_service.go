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
	ErrJenisDokumenNotFound     = errors.New("jenis dokumen not found")
	ErrJenisDokumenKodeExists   = errors.New("jenis dokumen with this kode already exists")
	ErrJenisDokumenKodeRequired = errors.New("kode is required")
	ErrJenisDokumenNamaRequired = errors.New("nama is required")
)

// ErrJenisDokumenReferenced is returned when trying to delete a jenis dokumen that is referenced by documents
type ErrJenisDokumenReferenced struct {
	Count int64
}

func (e *ErrJenisDokumenReferenced) Error() string {
	return fmt.Sprintf("cannot delete: referenced by %d documents", e.Count)
}

// JenisDokumenService handles jenis dokumen business logic
type JenisDokumenService struct {
	repo *repositories.JenisDokumenRepository
}

// NewJenisDokumenService creates a new JenisDokumenService instance
func NewJenisDokumenService() *JenisDokumenService {
	return &JenisDokumenService{
		repo: repositories.NewJenisDokumenRepository(),
	}
}

// CreateJenisDokumenInput represents input for creating a jenis dokumen
type CreateJenisDokumenInput struct {
	Kode string `json:"kode"`
	Nama string `json:"nama"`
}

// UpdateJenisDokumenInput represents input for updating a jenis dokumen
type UpdateJenisDokumenInput struct {
	Kode     *string `json:"kode"`
	Nama     *string `json:"nama"`
	IsActive *bool   `json:"is_active"`
}


// Validate validates the create input
func (i *CreateJenisDokumenInput) Validate() map[string]string {
	errors := make(map[string]string)

	if strings.TrimSpace(i.Kode) == "" {
		errors["kode"] = "Kode is required"
	}
	if strings.TrimSpace(i.Nama) == "" {
		errors["nama"] = "Nama is required"
	}

	return errors
}

// Create creates a new jenis dokumen
func (s *JenisDokumenService) Create(input *CreateJenisDokumenInput) (*models.JenisDokumen, error) {
	// Validate input
	if validationErrors := input.Validate(); len(validationErrors) > 0 {
		if _, ok := validationErrors["kode"]; ok {
			return nil, ErrJenisDokumenKodeRequired
		}
		if _, ok := validationErrors["nama"]; ok {
			return nil, ErrJenisDokumenNamaRequired
		}
	}

	// Check if kode already exists
	exists, err := s.repo.ExistsByKode(strings.TrimSpace(input.Kode), nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrJenisDokumenKodeExists
	}

	// Create jenis dokumen
	jenisDokumen := &models.JenisDokumen{
		ID:       uuid.New(),
		Kode:     strings.TrimSpace(input.Kode),
		Nama:     strings.TrimSpace(input.Nama),
		IsActive: true,
	}

	if err := s.repo.Create(jenisDokumen); err != nil {
		return nil, err
	}

	return jenisDokumen, nil
}

// GetByID retrieves a jenis dokumen by ID
func (s *JenisDokumenService) GetByID(id uuid.UUID) (*models.JenisDokumen, error) {
	jenisDokumen, err := s.repo.FindByID(id)
	if err != nil {
		return nil, ErrJenisDokumenNotFound
	}
	return jenisDokumen, nil
}

// GetAll retrieves all jenis dokumen with pagination
func (s *JenisDokumenService) GetAll(page, pageSize int, search string) (*repositories.PaginationResult, error) {
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

// GetAllActive retrieves all active jenis dokumen (for dropdowns)
func (s *JenisDokumenService) GetAllActive() ([]models.JenisDokumen, error) {
	return s.repo.GetAllActive()
}


// Update updates a jenis dokumen
func (s *JenisDokumenService) Update(id uuid.UUID, input *UpdateJenisDokumenInput) (*models.JenisDokumen, error) {
	// Find existing jenis dokumen
	jenisDokumen, err := s.repo.FindByID(id)
	if err != nil {
		return nil, ErrJenisDokumenNotFound
	}

	// Update fields if provided
	if input.Kode != nil {
		kode := strings.TrimSpace(*input.Kode)
		if kode == "" {
			return nil, ErrJenisDokumenKodeRequired
		}
		// Check if new kode already exists (excluding current record)
		exists, err := s.repo.ExistsByKode(kode, &id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrJenisDokumenKodeExists
		}
		jenisDokumen.Kode = kode
	}

	if input.Nama != nil {
		nama := strings.TrimSpace(*input.Nama)
		if nama == "" {
			return nil, ErrJenisDokumenNamaRequired
		}
		jenisDokumen.Nama = nama
	}

	if input.IsActive != nil {
		jenisDokumen.IsActive = *input.IsActive
	}

	if err := s.repo.Update(jenisDokumen); err != nil {
		return nil, err
	}

	return jenisDokumen, nil
}

// Delete deletes a jenis dokumen by ID
// Returns ErrJenisDokumenReferenced if the jenis dokumen is referenced by documents
func (s *JenisDokumenService) Delete(id uuid.UUID) error {
	// Check if jenis dokumen exists
	_, err := s.repo.FindByID(id)
	if err != nil {
		return ErrJenisDokumenNotFound
	}

	// Check referential integrity - count documents referencing this jenis dokumen
	count, err := s.repo.CountDocumentsByJenisDokumenID(id)
	if err != nil {
		return err
	}
	if count > 0 {
		return &ErrJenisDokumenReferenced{Count: count}
	}

	return s.repo.Delete(id)
}
