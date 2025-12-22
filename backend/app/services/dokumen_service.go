package services

import (
	"errors"
	"mime/multipart"
	"strings"
	"time"

	"dokumen-keuangan/app/models"
	"dokumen-keuangan/app/repositories"

	"github.com/google/uuid"
)

var (
	ErrDokumenNotFound          = errors.New("dokumen not found")
	ErrDokumenUnitKerjaRequired = errors.New("unit kerja is required")
	ErrDokumenPPTKRequired      = errors.New("PPTK is required")
	ErrDokumenJenisRequired     = errors.New("jenis dokumen is required")
	ErrDokumenSumberDanaRequired = errors.New("sumber dana is required")
	ErrDokumenNilaiRequired     = errors.New("nilai is required")
	ErrDokumenUraianRequired    = errors.New("uraian is required")
	ErrDokumenFileRequired      = errors.New("file is required")
)

// DokumenService handles dokumen business logic
type DokumenService struct {
	repo          *repositories.DokumenRepository
	unitKerjaRepo *repositories.UnitKerjaRepository
	pptkRepo      *repositories.PPTKRepository
	jenisDokRepo  *repositories.JenisDokumenRepository
	sumberDanaRepo *repositories.SumberDanaRepository
	fileService   *FileService
}

// NewDokumenService creates a new DokumenService instance
func NewDokumenService() *DokumenService {
	return &DokumenService{
		repo:          repositories.NewDokumenRepository(),
		unitKerjaRepo: repositories.NewUnitKerjaRepository(),
		pptkRepo:      repositories.NewPPTKRepository(),
		jenisDokRepo:  repositories.NewJenisDokumenRepository(),
		sumberDanaRepo: repositories.NewSumberDanaRepository(),
		fileService:   NewFileService(),
	}
}

// CreateDokumenInput represents input for creating a dokumen
type CreateDokumenInput struct {
	NomorDokumen   string    `json:"nomor_dokumen"`
	TanggalDokumen string    `json:"tanggal_dokumen"`
	UnitKerjaID    uuid.UUID `json:"unit_kerja_id"`
	PPTKID         uuid.UUID `json:"pptk_id"`
	JenisDokumenID uuid.UUID `json:"jenis_dokumen_id"`
	SumberDanaID   uuid.UUID `json:"sumber_dana_id"`
	Nilai          float64   `json:"nilai"`
	Uraian         string    `json:"uraian"`
	NomorKwitansi  string    `json:"nomor_kwitansi"`
}

// UpdateDokumenInput represents input for updating a dokumen
type UpdateDokumenInput struct {
	NomorDokumen   string     `json:"nomor_dokumen"`
	TanggalDokumen string     `json:"tanggal_dokumen"`
	UnitKerjaID    *uuid.UUID `json:"unit_kerja_id"`
	PPTKID         *uuid.UUID `json:"pptk_id"`
	JenisDokumenID *uuid.UUID `json:"jenis_dokumen_id"`
	SumberDanaID   *uuid.UUID `json:"sumber_dana_id"`
	Nilai          *float64   `json:"nilai"`
	Uraian         string     `json:"uraian"`
	NomorKwitansi  string     `json:"nomor_kwitansi"`
}

// DokumenFilterInput represents filter input for listing dokumen
type DokumenFilterInput struct {
	UnitKerjaID *uuid.UUID `json:"unit_kerja_id"`
	PPTKID      *uuid.UUID `json:"pptk_id"`
	StartDate   *string    `json:"start_date"`
	EndDate     *string    `json:"end_date"`
}

// Validate validates the create input
func (i *CreateDokumenInput) Validate() map[string]string {
	errs := make(map[string]string)

	if i.UnitKerjaID == uuid.Nil {
		errs["unit_kerja_id"] = "Unit Kerja is required"
	}
	if i.PPTKID == uuid.Nil {
		errs["pptk_id"] = "PPTK is required"
	}
	if i.JenisDokumenID == uuid.Nil {
		errs["jenis_dokumen_id"] = "Jenis Dokumen is required"
	}
	if i.SumberDanaID == uuid.Nil {
		errs["sumber_dana_id"] = "Sumber Dana is required"
	}
	if i.Nilai <= 0 {
		errs["nilai"] = "Nilai must be greater than 0"
	}
	if strings.TrimSpace(i.Uraian) == "" {
		errs["uraian"] = "Uraian is required"
	}

	return errs
}

// Create creates a new dokumen with file upload
func (s *DokumenService) Create(input *CreateDokumenInput, file *multipart.FileHeader, createdBy uuid.UUID) (*models.Dokumen, error) {
	// Validate input
	if validationErrors := input.Validate(); len(validationErrors) > 0 {
		if _, ok := validationErrors["unit_kerja_id"]; ok {
			return nil, ErrDokumenUnitKerjaRequired
		}
		if _, ok := validationErrors["pptk_id"]; ok {
			return nil, ErrDokumenPPTKRequired
		}
		if _, ok := validationErrors["jenis_dokumen_id"]; ok {
			return nil, ErrDokumenJenisRequired
		}
		if _, ok := validationErrors["sumber_dana_id"]; ok {
			return nil, ErrDokumenSumberDanaRequired
		}
		if _, ok := validationErrors["nilai"]; ok {
			return nil, ErrDokumenNilaiRequired
		}
		if _, ok := validationErrors["uraian"]; ok {
			return nil, ErrDokumenUraianRequired
		}
	}

	// Validate file
	if file == nil {
		return nil, ErrDokumenFileRequired
	}

	// Validate file type (PDF only)
	if !ValidatePDFFile(file.Header.Get("Content-Type"), file.Filename) {
		return nil, ErrInvalidFileType
	}

	// Validate references exist
	if _, err := s.unitKerjaRepo.FindByID(input.UnitKerjaID); err != nil {
		return nil, errors.New("unit kerja not found")
	}
	if _, err := s.pptkRepo.FindByID(input.PPTKID); err != nil {
		return nil, errors.New("PPTK not found")
	}
	if _, err := s.jenisDokRepo.FindByID(input.JenisDokumenID); err != nil {
		return nil, errors.New("jenis dokumen not found")
	}
	if _, err := s.sumberDanaRepo.FindByID(input.SumberDanaID); err != nil {
		return nil, errors.New("sumber dana not found")
	}

	// Upload file
	filePath, err := s.fileService.UploadDocument(file)
	if err != nil {
		return nil, err
	}

	// Parse tanggal dokumen
	var tanggalDokumen *time.Time
	if input.TanggalDokumen != "" {
		t, err := time.Parse("2006-01-02", input.TanggalDokumen)
		if err == nil {
			tanggalDokumen = &t
		}
	}

	// Create dokumen
	dokumen := &models.Dokumen{
		ID:             uuid.New(),
		NomorDokumen:   input.NomorDokumen,
		TanggalDokumen: tanggalDokumen,
		UnitKerjaID:    input.UnitKerjaID,
		PPTKID:         input.PPTKID,
		JenisDokumenID: input.JenisDokumenID,
		SumberDanaID:   input.SumberDanaID,
		Nilai:          input.Nilai,
		Uraian:         strings.TrimSpace(input.Uraian),
		NomorKwitansi:  input.NomorKwitansi,
		FilePath:       filePath,
		CreatedBy:      createdBy,
	}

	if err := s.repo.Create(dokumen); err != nil {
		return nil, err
	}

	// Reload with relationships
	return s.repo.FindByID(dokumen.ID)
}

// GetByID retrieves a dokumen by ID
func (s *DokumenService) GetByID(id uuid.UUID) (*models.Dokumen, error) {
	dokumen, err := s.repo.FindByID(id)
	if err != nil {
		return nil, ErrDokumenNotFound
	}
	return dokumen, nil
}

// GetAll retrieves all dokumen with pagination and filters
func (s *DokumenService) GetAll(page, pageSize int, input *DokumenFilterInput, userRole models.UserRole, userID uuid.UUID, year int) (*repositories.PaginationResult, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	filter := &repositories.DokumenFilter{
		Year: year,
	}

	// Apply role-based filtering
	if userRole == models.RoleOperator {
		// Operator can only see their own documents
		filter.CreatedBy = &userID
	}

	// Apply additional filters from input
	if input != nil {
		if input.UnitKerjaID != nil {
			filter.UnitKerjaID = input.UnitKerjaID
		}
		if input.PPTKID != nil {
			filter.PPTKID = input.PPTKID
		}
		if input.StartDate != nil && *input.StartDate != "" {
			startDate, err := time.Parse("2006-01-02", *input.StartDate)
			if err == nil {
				filter.StartDate = &startDate
			}
		}
		if input.EndDate != nil && *input.EndDate != "" {
			endDate, err := time.Parse("2006-01-02", *input.EndDate)
			if err == nil {
				filter.EndDate = &endDate
			}
		}
	}

	return s.repo.GetAll(page, pageSize, filter)
}

// GetFilePath returns the full path to a dokumen file
func (s *DokumenService) GetFilePath(id uuid.UUID) (string, error) {
	dokumen, err := s.repo.FindByID(id)
	if err != nil {
		return "", ErrDokumenNotFound
	}
	return s.fileService.GetFilePath(dokumen.FilePath), nil
}

// CanAccessDokumen checks if a user can access a specific dokumen
func (s *DokumenService) CanAccessDokumen(dokumenID uuid.UUID, userRole string, userID uuid.UUID) bool {
	// Admin and SuperAdmin can access all documents
	if userRole == string(models.RoleSuperAdmin) || userRole == string(models.RoleAdmin) {
		return true
	}

	// Operator can only access their own documents
	dokumen, err := s.repo.FindByID(dokumenID)
	if err != nil {
		return false
	}

	return dokumen.CreatedBy == userID
}

// Update updates a dokumen
func (s *DokumenService) Update(id uuid.UUID, input *UpdateDokumenInput, file *multipart.FileHeader) (*models.Dokumen, error) {
	// Get existing dokumen
	dokumen, err := s.repo.FindByID(id)
	if err != nil {
		return nil, ErrDokumenNotFound
	}

	// Update fields if provided
	if input.NomorDokumen != "" {
		dokumen.NomorDokumen = input.NomorDokumen
	}
	if input.TanggalDokumen != "" {
		t, err := time.Parse("2006-01-02", input.TanggalDokumen)
		if err == nil {
			dokumen.TanggalDokumen = &t
		}
	}
	if input.UnitKerjaID != nil {
		dokumen.UnitKerjaID = *input.UnitKerjaID
	}
	if input.PPTKID != nil {
		dokumen.PPTKID = *input.PPTKID
	}
	if input.JenisDokumenID != nil {
		dokumen.JenisDokumenID = *input.JenisDokumenID
	}
	if input.SumberDanaID != nil {
		dokumen.SumberDanaID = *input.SumberDanaID
	}
	if input.Nilai != nil {
		dokumen.Nilai = *input.Nilai
	}
	if input.Uraian != "" {
		dokumen.Uraian = strings.TrimSpace(input.Uraian)
	}
	if input.NomorKwitansi != "" {
		dokumen.NomorKwitansi = input.NomorKwitansi
	}

	// Handle file upload if provided
	if file != nil {
		// Validate file type (PDF only)
		if !ValidatePDFFile(file.Header.Get("Content-Type"), file.Filename) {
			return nil, ErrInvalidFileType
		}

		// Delete old file
		if dokumen.FilePath != "" {
			s.fileService.DeleteFile(dokumen.FilePath)
		}

		// Upload new file
		filePath, err := s.fileService.UploadDocument(file)
		if err != nil {
			return nil, err
		}
		dokumen.FilePath = filePath
	}

	if err := s.repo.Update(dokumen); err != nil {
		return nil, err
	}

	// Reload with relationships
	return s.repo.FindByID(dokumen.ID)
}

// Delete deletes a dokumen by ID
func (s *DokumenService) Delete(id uuid.UUID) error {
	dokumen, err := s.repo.FindByID(id)
	if err != nil {
		return ErrDokumenNotFound
	}

	// Delete the file
	if dokumen.FilePath != "" {
		s.fileService.DeleteFile(dokumen.FilePath)
	}

	return s.repo.Delete(id)
}

// CountByYear counts documents by year
func (s *DokumenService) CountByYear(year int, userRole models.UserRole, userID uuid.UUID) int64 {
	filter := &repositories.DokumenFilter{
		Year: year,
	}

	if userRole == models.RoleOperator {
		filter.CreatedBy = &userID
	}

	res, err := s.repo.GetAll(1, 1, filter)
	if err != nil {
		return 0
	}
	return res.Total
}

// GetRecent gets recent documents for dashboard
func (s *DokumenService) GetRecent(limit int, year int, userRole models.UserRole, userID uuid.UUID) ([]models.Dokumen, error) {
	filter := &repositories.DokumenFilter{
		Year: year,
	}

	if userRole == models.RoleOperator {
		filter.CreatedBy = &userID
	}

	res, err := s.repo.GetAll(1, limit, filter)
	if err != nil {
		return nil, err
	}

	// Type assertion
	if dokumens, ok := res.Data.([]models.Dokumen); ok {
		return dokumens, nil
	}
	return nil, errors.New("failed to cast result data")
}
