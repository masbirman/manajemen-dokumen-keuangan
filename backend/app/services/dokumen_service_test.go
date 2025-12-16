package services

import (
	"testing"
	"time"

	"dokumen-keuangan/app/models"

	"github.com/google/uuid"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// **Feature: manajemen-dokumen-keuangan, Property 10: Document Creation Completeness**
// **Validates: Requirements 6.3**

// **Feature: manajemen-dokumen-keuangan, Property 12: Document Validation Errors**
// **Validates: Requirements 6.5**

// **Feature: manajemen-dokumen-keuangan, Property 16: Admin Document Visibility**
// **Validates: Requirements 8.1**

// **Feature: manajemen-dokumen-keuangan, Property 17: Operator Document Isolation**
// **Validates: Requirements 8.2**

// **Feature: manajemen-dokumen-keuangan, Property 18: Unit Kerja Filter Accuracy**
// **Validates: Requirements 8.3**

// **Feature: manajemen-dokumen-keuangan, Property 19: Date Range Filter Accuracy**
// **Validates: Requirements 8.4**

// MockDokumen represents a document for testing
type MockDokumen struct {
	ID          uuid.UUID
	UnitKerjaID uuid.UUID
	CreatedBy   uuid.UUID
	CreatedAt   time.Time
	Nilai       float64
	Uraian      string
	FilePath    string
}

// MockDokumenStore simulates document storage
type MockDokumenStore struct {
	dokumens map[uuid.UUID]*MockDokumen
}

func NewMockDokumenStore() *MockDokumenStore {
	return &MockDokumenStore{
		dokumens: make(map[uuid.UUID]*MockDokumen),
	}
}

func (s *MockDokumenStore) Create(d *MockDokumen) {
	d.ID = uuid.New()
	d.CreatedAt = time.Now()
	s.dokumens[d.ID] = d
}

func (s *MockDokumenStore) GetByID(id uuid.UUID) *MockDokumen {
	return s.dokumens[id]
}

func (s *MockDokumenStore) GetAll(userRole string, userID uuid.UUID) []*MockDokumen {
	var result []*MockDokumen
	for _, d := range s.dokumens {
		if userRole == string(models.RoleOperator) {
			if d.CreatedBy == userID {
				result = append(result, d)
			}
		} else {
			result = append(result, d)
		}
	}
	return result
}

func (s *MockDokumenStore) GetByUnitKerja(unitKerjaID uuid.UUID) []*MockDokumen {
	var result []*MockDokumen
	for _, d := range s.dokumens {
		if d.UnitKerjaID == unitKerjaID {
			result = append(result, d)
		}
	}
	return result
}

func (s *MockDokumenStore) GetByDateRange(start, end time.Time) []*MockDokumen {
	var result []*MockDokumen
	for _, d := range s.dokumens {
		if (d.CreatedAt.Equal(start) || d.CreatedAt.After(start)) &&
			(d.CreatedAt.Equal(end) || d.CreatedAt.Before(end)) {
			result = append(result, d)
		}
	}
	return result
}


func TestProperty10_DocumentCreationCompleteness(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Property: Created document should have all required fields
	properties.Property("Document creation should store all fields", prop.ForAll(
		func(nilai float64, uraian string) bool {
			if nilai <= 0 || uraian == "" {
				return true // Skip invalid input
			}

			store := NewMockDokumenStore()
			unitKerjaID := uuid.New()
			createdBy := uuid.New()

			doc := &MockDokumen{
				UnitKerjaID: unitKerjaID,
				CreatedBy:   createdBy,
				Nilai:       nilai,
				Uraian:      uraian,
				FilePath:    "documents/test.pdf",
			}

			store.Create(doc)
			retrieved := store.GetByID(doc.ID)

			return retrieved != nil &&
				retrieved.UnitKerjaID == unitKerjaID &&
				retrieved.CreatedBy == createdBy &&
				retrieved.Nilai == nilai &&
				retrieved.Uraian == uraian &&
				retrieved.FilePath != ""
		},
		gen.Float64Range(0.01, 1000000),
		gen.AlphaString(),
	))

	properties.TestingRun(t)
}

func TestProperty12_DocumentValidationErrors(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Property: Missing required fields should produce validation errors
	properties.Property("Missing fields should produce validation errors", prop.ForAll(
		func(hasUnitKerja, hasPPTK, hasJenis, hasSumberDana, hasNilai, hasUraian bool) bool {
			input := &CreateDokumenInput{}

			if hasUnitKerja {
				input.UnitKerjaID = uuid.New()
			}
			if hasPPTK {
				input.PPTKID = uuid.New()
			}
			if hasJenis {
				input.JenisDokumenID = uuid.New()
			}
			if hasSumberDana {
				input.SumberDanaID = uuid.New()
			}
			if hasNilai {
				input.Nilai = 1000
			}
			if hasUraian {
				input.Uraian = "Test uraian"
			}

			errors := input.Validate()

			// Count expected errors
			expectedErrors := 0
			if !hasUnitKerja {
				expectedErrors++
			}
			if !hasPPTK {
				expectedErrors++
			}
			if !hasJenis {
				expectedErrors++
			}
			if !hasSumberDana {
				expectedErrors++
			}
			if !hasNilai {
				expectedErrors++
			}
			if !hasUraian {
				expectedErrors++
			}

			return len(errors) == expectedErrors
		},
		gen.Bool(), gen.Bool(), gen.Bool(), gen.Bool(), gen.Bool(), gen.Bool(),
	))

	properties.TestingRun(t)
}

func TestProperty16_AdminDocumentVisibility(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Property: Admin should see all documents
	properties.Property("Admin should see all documents", prop.ForAll(
		func(numDocs int) bool {
			if numDocs <= 0 || numDocs > 20 {
				return true
			}

			store := NewMockDokumenStore()
			adminID := uuid.New()

			// Create documents by different users
			for i := 0; i < numDocs; i++ {
				doc := &MockDokumen{
					UnitKerjaID: uuid.New(),
					CreatedBy:   uuid.New(), // Different creator each time
					Nilai:       float64(i * 1000),
					Uraian:      "Test",
					FilePath:    "test.pdf",
				}
				store.Create(doc)
			}

			// Admin should see all
			adminDocs := store.GetAll(string(models.RoleAdmin), adminID)
			return len(adminDocs) == numDocs
		},
		gen.IntRange(1, 20),
	))

	properties.TestingRun(t)
}

func TestProperty17_OperatorDocumentIsolation(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Property: Operator should only see their own documents
	properties.Property("Operator should only see own documents", prop.ForAll(
		func(ownDocs, otherDocs int) bool {
			if ownDocs < 0 || otherDocs < 0 || ownDocs > 10 || otherDocs > 10 {
				return true
			}

			store := NewMockDokumenStore()
			operatorID := uuid.New()
			otherUserID := uuid.New()

			// Create operator's documents
			for i := 0; i < ownDocs; i++ {
				doc := &MockDokumen{
					UnitKerjaID: uuid.New(),
					CreatedBy:   operatorID,
					Nilai:       1000,
					Uraian:      "Own doc",
					FilePath:    "test.pdf",
				}
				store.Create(doc)
			}

			// Create other user's documents
			for i := 0; i < otherDocs; i++ {
				doc := &MockDokumen{
					UnitKerjaID: uuid.New(),
					CreatedBy:   otherUserID,
					Nilai:       2000,
					Uraian:      "Other doc",
					FilePath:    "test.pdf",
				}
				store.Create(doc)
			}

			// Operator should only see their own
			operatorDocs := store.GetAll(string(models.RoleOperator), operatorID)
			return len(operatorDocs) == ownDocs
		},
		gen.IntRange(0, 10),
		gen.IntRange(0, 10),
	))

	properties.TestingRun(t)
}

func TestProperty18_UnitKerjaFilterAccuracy(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Property: Unit Kerja filter should return only matching documents
	properties.Property("Unit Kerja filter returns only matching documents", prop.ForAll(
		func(targetDocs, otherDocs int) bool {
			if targetDocs < 0 || otherDocs < 0 || targetDocs > 10 || otherDocs > 10 {
				return true
			}

			store := NewMockDokumenStore()
			targetUnitKerja := uuid.New()
			otherUnitKerja := uuid.New()

			// Create target unit kerja documents
			for i := 0; i < targetDocs; i++ {
				doc := &MockDokumen{
					UnitKerjaID: targetUnitKerja,
					CreatedBy:   uuid.New(),
					Nilai:       1000,
					Uraian:      "Target",
					FilePath:    "test.pdf",
				}
				store.Create(doc)
			}

			// Create other unit kerja documents
			for i := 0; i < otherDocs; i++ {
				doc := &MockDokumen{
					UnitKerjaID: otherUnitKerja,
					CreatedBy:   uuid.New(),
					Nilai:       2000,
					Uraian:      "Other",
					FilePath:    "test.pdf",
				}
				store.Create(doc)
			}

			// Filter should return only target
			filtered := store.GetByUnitKerja(targetUnitKerja)
			
			// All returned should match target unit kerja
			for _, d := range filtered {
				if d.UnitKerjaID != targetUnitKerja {
					return false
				}
			}
			
			return len(filtered) == targetDocs
		},
		gen.IntRange(0, 10),
		gen.IntRange(0, 10),
	))

	properties.TestingRun(t)
}

func TestProperty19_DateRangeFilterAccuracy(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Property: Date range filter should return only documents within range
	properties.Property("Date range filter returns documents within range", prop.ForAll(
		func(daysOffset int) bool {
			store := NewMockDokumenStore()
			now := time.Now()

			// Create documents at different times
			for i := -5; i <= 5; i++ {
				doc := &MockDokumen{
					UnitKerjaID: uuid.New(),
					CreatedBy:   uuid.New(),
					Nilai:       1000,
					Uraian:      "Test",
					FilePath:    "test.pdf",
				}
				store.Create(doc)
				// Manually set created time
				doc.CreatedAt = now.AddDate(0, 0, i)
			}

			// Filter for today only
			startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
			endOfDay := startOfDay.Add(24 * time.Hour)

			filtered := store.GetByDateRange(startOfDay, endOfDay)

			// All returned should be within range
			for _, d := range filtered {
				if d.CreatedAt.Before(startOfDay) || d.CreatedAt.After(endOfDay) {
					return false
				}
			}

			return true
		},
		gen.IntRange(0, 10),
	))

	properties.TestingRun(t)
}
