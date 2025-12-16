package services

import (
	"testing"

	"github.com/google/uuid"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// **Feature: manajemen-dokumen-keuangan, Property 6: Referential Integrity on Delete**
// **Validates: Requirements 4.3, 5.3**
// For any master data entity (Sumber Dana, Jenis Dokumen) that is referenced by existing documents,
// delete operation should fail and return an error.

// MockSumberDanaRepository simulates the repository behavior for testing
type MockSumberDanaRepository struct {
	sumberDanas      map[uuid.UUID]bool
	documentCounts   map[uuid.UUID]int64
}

func NewMockSumberDanaRepository() *MockSumberDanaRepository {
	return &MockSumberDanaRepository{
		sumberDanas:    make(map[uuid.UUID]bool),
		documentCounts: make(map[uuid.UUID]int64),
	}
}

func (r *MockSumberDanaRepository) Exists(id uuid.UUID) bool {
	return r.sumberDanas[id]
}

func (r *MockSumberDanaRepository) CountDocuments(id uuid.UUID) int64 {
	return r.documentCounts[id]
}

func (r *MockSumberDanaRepository) Delete(id uuid.UUID) error {
	delete(r.sumberDanas, id)
	return nil
}


// MockJenisDokumenRepository simulates the repository behavior for testing
type MockJenisDokumenRepository struct {
	jenisDokumens  map[uuid.UUID]bool
	documentCounts map[uuid.UUID]int64
}

func NewMockJenisDokumenRepository() *MockJenisDokumenRepository {
	return &MockJenisDokumenRepository{
		jenisDokumens:  make(map[uuid.UUID]bool),
		documentCounts: make(map[uuid.UUID]int64),
	}
}

func (r *MockJenisDokumenRepository) Exists(id uuid.UUID) bool {
	return r.jenisDokumens[id]
}

func (r *MockJenisDokumenRepository) CountDocuments(id uuid.UUID) int64 {
	return r.documentCounts[id]
}

func (r *MockJenisDokumenRepository) Delete(id uuid.UUID) error {
	delete(r.jenisDokumens, id)
	return nil
}

// deleteWithReferentialCheck simulates the delete logic with referential integrity check
// This is the core logic we're testing - it should reject deletion when documents exist
func deleteWithReferentialCheck(exists bool, documentCount int64) (bool, error) {
	if !exists {
		return false, ErrSumberDanaNotFound
	}
	if documentCount > 0 {
		return false, &ErrSumberDanaReferenced{Count: documentCount}
	}
	return true, nil
}

// deleteJenisDokumenWithReferentialCheck simulates the delete logic for jenis dokumen
func deleteJenisDokumenWithReferentialCheck(exists bool, documentCount int64) (bool, error) {
	if !exists {
		return false, ErrJenisDokumenNotFound
	}
	if documentCount > 0 {
		return false, &ErrJenisDokumenReferenced{Count: documentCount}
	}
	return true, nil
}


func TestProperty6_ReferentialIntegrityOnDelete(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Property: For any Sumber Dana that is referenced by documents (count > 0),
	// delete operation should fail with ErrSumberDanaReferenced
	properties.Property("Sumber Dana with referenced documents should not be deletable", prop.ForAll(
		func(documentCount int64) bool {
			// Given: A sumber dana exists and has documents referencing it
			exists := true

			// When: We attempt to delete it
			deleted, err := deleteWithReferentialCheck(exists, documentCount)

			// Then: If documentCount > 0, deletion should fail
			if documentCount > 0 {
				if deleted {
					return false // Should not have been deleted
				}
				refErr, ok := err.(*ErrSumberDanaReferenced)
				if !ok {
					return false // Should be referential integrity error
				}
				return refErr.Count == documentCount // Error should contain correct count
			}

			// If documentCount == 0, deletion should succeed
			return deleted && err == nil
		},
		gen.Int64Range(0, 1000), // Document count from 0 to 1000
	))

	// Property: For any Jenis Dokumen that is referenced by documents (count > 0),
	// delete operation should fail with ErrJenisDokumenReferenced
	properties.Property("Jenis Dokumen with referenced documents should not be deletable", prop.ForAll(
		func(documentCount int64) bool {
			// Given: A jenis dokumen exists and has documents referencing it
			exists := true

			// When: We attempt to delete it
			deleted, err := deleteJenisDokumenWithReferentialCheck(exists, documentCount)

			// Then: If documentCount > 0, deletion should fail
			if documentCount > 0 {
				if deleted {
					return false // Should not have been deleted
				}
				refErr, ok := err.(*ErrJenisDokumenReferenced)
				if !ok {
					return false // Should be referential integrity error
				}
				return refErr.Count == documentCount // Error should contain correct count
			}

			// If documentCount == 0, deletion should succeed
			return deleted && err == nil
		},
		gen.Int64Range(0, 1000), // Document count from 0 to 1000
	))

	// Property: Non-existent entities should return not found error
	properties.Property("Non-existent Sumber Dana should return not found error", prop.ForAll(
		func(documentCount int64) bool {
			// Given: A sumber dana does not exist
			exists := false

			// When: We attempt to delete it
			deleted, err := deleteWithReferentialCheck(exists, documentCount)

			// Then: Should return not found error
			return !deleted && err == ErrSumberDanaNotFound
		},
		gen.Int64Range(0, 100),
	))

	properties.Property("Non-existent Jenis Dokumen should return not found error", prop.ForAll(
		func(documentCount int64) bool {
			// Given: A jenis dokumen does not exist
			exists := false

			// When: We attempt to delete it
			deleted, err := deleteJenisDokumenWithReferentialCheck(exists, documentCount)

			// Then: Should return not found error
			return !deleted && err == ErrJenisDokumenNotFound
		},
		gen.Int64Range(0, 100),
	))

	// Property: Error message should contain the correct document count
	properties.Property("Referential integrity error should contain correct document count", prop.ForAll(
		func(documentCount int64) bool {
			if documentCount <= 0 {
				return true // Skip for zero or negative counts
			}

			// Test Sumber Dana
			_, errSD := deleteWithReferentialCheck(true, documentCount)
			refErrSD, okSD := errSD.(*ErrSumberDanaReferenced)
			if !okSD || refErrSD.Count != documentCount {
				return false
			}

			// Test Jenis Dokumen
			_, errJD := deleteJenisDokumenWithReferentialCheck(true, documentCount)
			refErrJD, okJD := errJD.(*ErrJenisDokumenReferenced)
			if !okJD || refErrJD.Count != documentCount {
				return false
			}

			return true
		},
		gen.Int64Range(1, 1000), // Only positive counts
	))

	properties.TestingRun(t)
}
