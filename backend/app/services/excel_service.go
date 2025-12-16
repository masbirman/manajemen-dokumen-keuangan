package services

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	"dokumen-keuangan/app/models"
	"dokumen-keuangan/app/repositories"

	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

var (
	ErrInvalidExcelFile    = errors.New("invalid Excel file")
	ErrEmptyExcelFile      = errors.New("Excel file is empty")
	ErrInvalidExcelFormat  = errors.New("invalid Excel format")
)

// ExcelValidationError represents a validation error with row number
type ExcelValidationError struct {
	Row     int    `json:"row"`
	Column  string `json:"column"`
	Message string `json:"message"`
}

// ExcelImportResult represents the result of an Excel import operation
type ExcelImportResult struct {
	TotalRows    int                    `json:"total_rows"`
	SuccessCount int                    `json:"success_count"`
	ErrorCount   int                    `json:"error_count"`
	Errors       []ExcelValidationError `json:"errors,omitempty"`
}

// ExcelService handles Excel import/export operations
type ExcelService struct {
	unitKerjaRepo *repositories.UnitKerjaRepository
	pptkRepo      *repositories.PPTKRepository
}

// NewExcelService creates a new ExcelService instance
func NewExcelService() *ExcelService {
	return &ExcelService{
		unitKerjaRepo: repositories.NewUnitKerjaRepository(),
		pptkRepo:      repositories.NewPPTKRepository(),
	}
}


// GenerateUnitKerjaTemplate generates an Excel template for Unit Kerja import
func (s *ExcelService) GenerateUnitKerjaTemplate() (*bytes.Buffer, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Unit Kerja"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	// Set headers
	headers := []string{"Kode", "Nama"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	// Style headers
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"4472C4"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	f.SetCellStyle(sheetName, "A1", "B1", headerStyle)

	// Set column widths
	f.SetColWidth(sheetName, "A", "A", 20)
	f.SetColWidth(sheetName, "B", "B", 40)

	// Add example data
	f.SetCellValue(sheetName, "A2", "UK-001")
	f.SetCellValue(sheetName, "B2", "Dinas Pendidikan")
	f.SetCellValue(sheetName, "A3", "UK-002")
	f.SetCellValue(sheetName, "B3", "Dinas Kesehatan")

	// Write to buffer
	buf := new(bytes.Buffer)
	if err := f.Write(buf); err != nil {
		return nil, err
	}

	return buf, nil
}

// ImportUnitKerja imports Unit Kerja data from Excel file
func (s *ExcelService) ImportUnitKerja(reader io.Reader) (*ExcelImportResult, error) {
	// Read file into buffer
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, reader); err != nil {
		return nil, ErrInvalidExcelFile
	}

	f, err := excelize.OpenReader(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return nil, ErrInvalidExcelFile
	}
	defer f.Close()

	// Get first sheet
	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return nil, ErrEmptyExcelFile
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, ErrInvalidExcelFile
	}

	if len(rows) < 2 {
		return nil, ErrEmptyExcelFile
	}

	result := &ExcelImportResult{
		TotalRows: len(rows) - 1, // Exclude header
		Errors:    []ExcelValidationError{},
	}

	// Process data rows (skip header)
	for i, row := range rows[1:] {
		rowNum := i + 2 // Excel row number (1-indexed, skip header)

		if len(row) < 2 {
			result.Errors = append(result.Errors, ExcelValidationError{
				Row:     rowNum,
				Column:  "All",
				Message: "Row has insufficient columns",
			})
			result.ErrorCount++
			continue
		}

		kode := strings.TrimSpace(row[0])
		nama := strings.TrimSpace(row[1])

		// Validate kode
		if kode == "" {
			result.Errors = append(result.Errors, ExcelValidationError{
				Row:     rowNum,
				Column:  "Kode",
				Message: "Kode is required",
			})
			result.ErrorCount++
			continue
		}

		// Validate nama
		if nama == "" {
			result.Errors = append(result.Errors, ExcelValidationError{
				Row:     rowNum,
				Column:  "Nama",
				Message: "Nama is required",
			})
			result.ErrorCount++
			continue
		}

		// Check if kode already exists
		exists, err := s.unitKerjaRepo.ExistsByKode(kode, nil)
		if err != nil {
			result.Errors = append(result.Errors, ExcelValidationError{
				Row:     rowNum,
				Column:  "Kode",
				Message: "Database error checking kode",
			})
			result.ErrorCount++
			continue
		}
		if exists {
			result.Errors = append(result.Errors, ExcelValidationError{
				Row:     rowNum,
				Column:  "Kode",
				Message: fmt.Sprintf("Kode '%s' already exists", kode),
			})
			result.ErrorCount++
			continue
		}

		// Create Unit Kerja
		unitKerja := &models.UnitKerja{
			ID:       uuid.New(),
			Kode:     kode,
			Nama:     nama,
			IsActive: true,
		}

		if err := s.unitKerjaRepo.Create(unitKerja); err != nil {
			result.Errors = append(result.Errors, ExcelValidationError{
				Row:     rowNum,
				Column:  "All",
				Message: "Failed to create record",
			})
			result.ErrorCount++
			continue
		}

		result.SuccessCount++
	}

	return result, nil
}


// GeneratePPTKTemplate generates an Excel template for PPTK import
func (s *ExcelService) GeneratePPTKTemplate() (*bytes.Buffer, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "PPTK"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	// Set headers - NIP is optional (will be auto-generated if empty)
	headers := []string{"NIP (Opsional)", "Nama", "Kode Unit Kerja"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	// Style headers
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"4472C4"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	f.SetCellStyle(sheetName, "A1", "C1", headerStyle)

	// Set column widths
	f.SetColWidth(sheetName, "A", "A", 25)
	f.SetColWidth(sheetName, "B", "B", 40)
	f.SetColWidth(sheetName, "C", "C", 20)

	// Add instructions row first (row 2)
	noteStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Italic: true, Color: "666666"},
		Alignment: &excelize.Alignment{Horizontal: "left"},
	})
	f.SetCellValue(sheetName, "A2", "Catatan: NIP bersifat opsional. Jika dikosongkan, sistem akan generate NIP otomatis. Hapus baris ini sebelum import.")
	f.MergeCell(sheetName, "A2", "C2")
	f.SetCellStyle(sheetName, "A2", "C2", noteStyle)

	// Add example data - showing both with NIP and without NIP (starting from row 3)
	f.SetCellValue(sheetName, "A3", "198501012010011001")
	f.SetCellValue(sheetName, "B3", "Budi Santoso")
	f.SetCellValue(sheetName, "C3", "UK-001")
	f.SetCellValue(sheetName, "A4", "") // NIP kosong - akan di-generate otomatis
	f.SetCellValue(sheetName, "B4", "Siti Rahayu")
	f.SetCellValue(sheetName, "C4", "UK-002")

	// Add reference sheet for Unit Kerja
	refSheetName := "Referensi Unit Kerja"
	_, err = f.NewSheet(refSheetName)
	if err != nil {
		return nil, err
	}

	// Add Unit Kerja reference data
	f.SetCellValue(refSheetName, "A1", "Kode")
	f.SetCellValue(refSheetName, "B1", "Nama")
	f.SetCellStyle(refSheetName, "A1", "B1", headerStyle)
	f.SetColWidth(refSheetName, "A", "A", 20)
	f.SetColWidth(refSheetName, "B", "B", 40)

	// Get all active Unit Kerja for reference
	unitKerjas, err := s.unitKerjaRepo.GetAllActive()
	if err == nil {
		for i, uk := range unitKerjas {
			rowNum := i + 2
			f.SetCellValue(refSheetName, fmt.Sprintf("A%d", rowNum), uk.Kode)
			f.SetCellValue(refSheetName, fmt.Sprintf("B%d", rowNum), uk.Nama)
		}
	}

	// Write to buffer
	buf := new(bytes.Buffer)
	if err := f.Write(buf); err != nil {
		return nil, err
	}

	return buf, nil
}

// ImportPPTK imports PPTK data from Excel file
func (s *ExcelService) ImportPPTK(reader io.Reader) (*ExcelImportResult, error) {
	// Read file into buffer
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, reader); err != nil {
		return nil, ErrInvalidExcelFile
	}

	f, err := excelize.OpenReader(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return nil, ErrInvalidExcelFile
	}
	defer f.Close()

	// Get first sheet
	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return nil, ErrEmptyExcelFile
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, ErrInvalidExcelFile
	}

	if len(rows) < 2 {
		return nil, ErrEmptyExcelFile
	}

	result := &ExcelImportResult{
		TotalRows: len(rows) - 1, // Exclude header
		Errors:    []ExcelValidationError{},
	}

	// Process data rows (skip header)
	for i, row := range rows[1:] {
		rowNum := i + 2 // Excel row number (1-indexed, skip header)

		// Need at least 2 columns (Nama and Kode Unit Kerja), NIP is optional
		if len(row) < 2 {
			result.Errors = append(result.Errors, ExcelValidationError{
				Row:     rowNum,
				Column:  "All",
				Message: "Row has insufficient columns (need at least Nama, Kode Unit Kerja)",
			})
			result.ErrorCount++
			continue
		}

		// Get values - handle different column arrangements
		var nip, nama, kodeUnitKerja string
		if len(row) >= 3 {
			nip = strings.TrimSpace(row[0])
			nama = strings.TrimSpace(row[1])
			kodeUnitKerja = strings.TrimSpace(row[2])
		} else {
			// If only 2 columns, assume Nama and Kode Unit Kerja
			nama = strings.TrimSpace(row[0])
			kodeUnitKerja = strings.TrimSpace(row[1])
		}

		// Validate nama
		if nama == "" {
			result.Errors = append(result.Errors, ExcelValidationError{
				Row:     rowNum,
				Column:  "Nama",
				Message: "Nama is required",
			})
			result.ErrorCount++
			continue
		}

		// Validate kode unit kerja
		if kodeUnitKerja == "" {
			result.Errors = append(result.Errors, ExcelValidationError{
				Row:     rowNum,
				Column:  "Kode Unit Kerja",
				Message: "Kode Unit Kerja is required",
			})
			result.ErrorCount++
			continue
		}

		// If NIP is empty, generate a unique one
		if nip == "" {
			nip = fmt.Sprintf("PPTK-%s", uuid.New().String()[:8])
		}

		// Check if NIP already exists
		exists, err := s.pptkRepo.ExistsByNIP(nip, nil)
		if err != nil {
			result.Errors = append(result.Errors, ExcelValidationError{
				Row:     rowNum,
				Column:  "NIP",
				Message: "Database error checking NIP",
			})
			result.ErrorCount++
			continue
		}
		if exists {
			result.Errors = append(result.Errors, ExcelValidationError{
				Row:     rowNum,
				Column:  "NIP",
				Message: fmt.Sprintf("NIP '%s' already exists", nip),
			})
			result.ErrorCount++
			continue
		}

		// Find Unit Kerja by kode
		unitKerja, err := s.unitKerjaRepo.FindByKode(kodeUnitKerja)
		if err != nil {
			result.Errors = append(result.Errors, ExcelValidationError{
				Row:     rowNum,
				Column:  "Kode Unit Kerja",
				Message: fmt.Sprintf("Unit Kerja with kode '%s' not found", kodeUnitKerja),
			})
			result.ErrorCount++
			continue
		}

		// Create PPTK
		pptk := &models.PPTK{
			ID:          uuid.New(),
			NIP:         nip,
			Nama:        nama,
			UnitKerjaID: unitKerja.ID,
			IsActive:    true,
		}

		if err := s.pptkRepo.Create(pptk); err != nil {
			result.Errors = append(result.Errors, ExcelValidationError{
				Row:     rowNum,
				Column:  "All",
				Message: fmt.Sprintf("Failed to create record: %v", err),
			})
			result.ErrorCount++
			continue
		}

		result.SuccessCount++
	}

	return result, nil
}
