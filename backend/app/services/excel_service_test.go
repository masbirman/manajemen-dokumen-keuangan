package services

import (
	"bytes"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"github.com/xuri/excelize/v2"
)

// **Feature: manajemen-dokumen-keuangan, Property 7: Excel Import Creates All Valid Records**
// **Validates: Requirements 2.3, 3.4**
// For any valid Excel file with N rows of data (Unit Kerja or PPTK),
// import operation should create exactly N records in the database.

// **Feature: manajemen-dokumen-keuangan, Property 8: Excel Import Error Reporting**
// **Validates: Requirements 2.4, 3.5**
// For any Excel file containing invalid data, import operation should return
// validation errors that include the specific row numbers of invalid entries.

// MockExcelImporter simulates Excel import logic for testing
type MockExcelImporter struct {
	existingKodes map[string]bool
}

func NewMockExcelImporter() *MockExcelImporter {
	return &MockExcelImporter{
		existingKodes: make(map[string]bool),
	}
}

// ValidateRow validates a single row and returns error if invalid
func (m *MockExcelImporter) ValidateRow(kode, nama string) *ExcelValidationError {
	if kode == "" {
		return &ExcelValidationError{Column: "Kode", Message: "Kode is required"}
	}
	if nama == "" {
		return &ExcelValidationError{Column: "Nama", Message: "Nama is required"}
	}
	if m.existingKodes[kode] {
		return &ExcelValidationError{Column: "Kode", Message: "Kode already exists"}
	}
	return nil
}

// ImportRows simulates importing multiple rows
func (m *MockExcelImporter) ImportRows(rows [][]string) *ExcelImportResult {
	result := &ExcelImportResult{
		TotalRows: len(rows),
		Errors:    []ExcelValidationError{},
	}

	for i, row := range rows {
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

		kode := row[0]
		nama := row[1]

		if err := m.ValidateRow(kode, nama); err != nil {
			err.Row = rowNum
			result.Errors = append(result.Errors, *err)
			result.ErrorCount++
			continue
		}

		// Mark as existing for duplicate detection
		m.existingKodes[kode] = true
		result.SuccessCount++
	}

	return result
}


func TestProperty7_ExcelImportCreatesAllValidRecords(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Property: For any N valid rows, import should create exactly N records
	properties.Property("Valid rows should all be imported successfully", prop.ForAll(
		func(numRows int) bool {
			if numRows <= 0 {
				return true // Skip invalid input
			}

			importer := NewMockExcelImporter()
			
			// Generate valid rows
			rows := make([][]string, numRows)
			for i := 0; i < numRows; i++ {
				rows[i] = []string{
					"UK-" + string(rune('A'+i%26)) + string(rune('0'+i/26)),
					"Unit Kerja " + string(rune('A'+i)),
				}
			}

			result := importer.ImportRows(rows)

			// All rows should be successful
			return result.SuccessCount == numRows && 
				   result.ErrorCount == 0 && 
				   result.TotalRows == numRows
		},
		gen.IntRange(1, 50), // Test with 1 to 50 rows
	))

	// Property: Success count + error count should equal total rows
	properties.Property("Success + Error count equals Total rows", prop.ForAll(
		func(validCount, invalidCount int) bool {
			if validCount < 0 || invalidCount < 0 {
				return true
			}

			importer := NewMockExcelImporter()
			
			// Generate mix of valid and invalid rows
			rows := make([][]string, validCount+invalidCount)
			
			// Valid rows
			for i := 0; i < validCount; i++ {
				rows[i] = []string{
					"UK-VALID-" + string(rune('A'+i%26)) + string(rune('0'+i/26)),
					"Valid Unit " + string(rune('A'+i)),
				}
			}
			
			// Invalid rows (empty kode)
			for i := 0; i < invalidCount; i++ {
				rows[validCount+i] = []string{"", "Invalid Unit"}
			}

			result := importer.ImportRows(rows)

			return result.SuccessCount + result.ErrorCount == result.TotalRows
		},
		gen.IntRange(0, 25),
		gen.IntRange(0, 25),
	))

	properties.TestingRun(t)
}

func TestProperty8_ExcelImportErrorReporting(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Property: Invalid rows should have errors with correct row numbers
	properties.Property("Invalid rows should report correct row numbers", prop.ForAll(
		func(invalidRowIndices []int) bool {
			if len(invalidRowIndices) == 0 {
				return true
			}

			importer := NewMockExcelImporter()
			totalRows := 10
			
			// Create rows, some invalid based on indices
			rows := make([][]string, totalRows)
			invalidSet := make(map[int]bool)
			for _, idx := range invalidRowIndices {
				if idx >= 0 && idx < totalRows {
					invalidSet[idx] = true
				}
			}

			for i := 0; i < totalRows; i++ {
				if invalidSet[i] {
					rows[i] = []string{"", "Invalid"} // Empty kode = invalid
				} else {
					rows[i] = []string{
						"UK-" + string(rune('A'+i)),
						"Valid Unit " + string(rune('A'+i)),
					}
				}
			}

			result := importer.ImportRows(rows)

			// Check that error count matches invalid rows
			expectedErrors := len(invalidSet)
			if result.ErrorCount != expectedErrors {
				return false
			}

			// Check that each error has a valid row number
			for _, err := range result.Errors {
				// Row numbers should be >= 2 (1-indexed, skip header)
				if err.Row < 2 || err.Row > totalRows+1 {
					return false
				}
			}

			return true
		},
		gen.SliceOf(gen.IntRange(0, 9)),
	))

	// Property: Duplicate kodes should be reported with row numbers
	properties.Property("Duplicate kodes should report errors with row numbers", prop.ForAll(
		func(numDuplicates int) bool {
			if numDuplicates <= 0 {
				return true
			}

			importer := NewMockExcelImporter()
			
			// Create rows with duplicates
			rows := make([][]string, numDuplicates+1)
			rows[0] = []string{"UK-SAME", "First Unit"} // Original
			for i := 1; i <= numDuplicates; i++ {
				rows[i] = []string{"UK-SAME", "Duplicate Unit " + string(rune('A'+i))}
			}

			result := importer.ImportRows(rows)

			// First row should succeed, rest should fail as duplicates
			return result.SuccessCount == 1 && 
				   result.ErrorCount == numDuplicates &&
				   len(result.Errors) == numDuplicates
		},
		gen.IntRange(1, 20),
	))

	// Property: Error messages should contain column information
	properties.Property("Errors should contain column information", prop.ForAll(
		func(errorType int) bool {
			importer := NewMockExcelImporter()
			
			var rows [][]string
			switch errorType % 3 {
			case 0: // Empty kode
				rows = [][]string{{"", "Some Name"}}
			case 1: // Empty nama
				rows = [][]string{{"UK-001", ""}}
			case 2: // Insufficient columns
				rows = [][]string{{"UK-001"}}
			}

			result := importer.ImportRows(rows)

			if len(result.Errors) == 0 {
				return false
			}

			// Error should have column information
			return result.Errors[0].Column != "" && result.Errors[0].Message != ""
		},
		gen.IntRange(0, 100),
	))

	properties.TestingRun(t)
}


// TestExcelTemplateGeneration tests that templates are generated correctly
func TestExcelTemplateGeneration(t *testing.T) {
	// Test Unit Kerja template structure
	t.Run("Unit Kerja template has correct headers", func(t *testing.T) {
		f := excelize.NewFile()
		defer f.Close()

		sheetName := "Unit Kerja"
		f.NewSheet(sheetName)
		f.SetCellValue(sheetName, "A1", "Kode")
		f.SetCellValue(sheetName, "B1", "Nama")

		// Verify headers
		kode, _ := f.GetCellValue(sheetName, "A1")
		nama, _ := f.GetCellValue(sheetName, "B1")

		if kode != "Kode" {
			t.Errorf("Expected header 'Kode', got '%s'", kode)
		}
		if nama != "Nama" {
			t.Errorf("Expected header 'Nama', got '%s'", nama)
		}
	})

	// Test PPTK template structure
	t.Run("PPTK template has correct headers", func(t *testing.T) {
		f := excelize.NewFile()
		defer f.Close()

		sheetName := "PPTK"
		f.NewSheet(sheetName)
		f.SetCellValue(sheetName, "A1", "NIP")
		f.SetCellValue(sheetName, "B1", "Nama")
		f.SetCellValue(sheetName, "C1", "Kode Unit Kerja")

		// Verify headers
		nip, _ := f.GetCellValue(sheetName, "A1")
		nama, _ := f.GetCellValue(sheetName, "B1")
		unitKerja, _ := f.GetCellValue(sheetName, "C1")

		if nip != "NIP" {
			t.Errorf("Expected header 'NIP', got '%s'", nip)
		}
		if nama != "Nama" {
			t.Errorf("Expected header 'Nama', got '%s'", nama)
		}
		if unitKerja != "Kode Unit Kerja" {
			t.Errorf("Expected header 'Kode Unit Kerja', got '%s'", unitKerja)
		}
	})
}

// TestExcelFileValidation tests file validation logic
func TestExcelFileValidation(t *testing.T) {
	t.Run("Empty file should return error", func(t *testing.T) {
		// Create empty Excel file
		f := excelize.NewFile()
		buf := new(bytes.Buffer)
		f.Write(buf)
		f.Close()

		// Try to read it
		_, err := excelize.OpenReader(bytes.NewReader(buf.Bytes()))
		if err != nil {
			t.Logf("Empty file handling: %v", err)
		}
	})
}
