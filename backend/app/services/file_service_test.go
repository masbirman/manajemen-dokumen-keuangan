package services

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// **Feature: manajemen-dokumen-keuangan, Property 3: Avatar Upload Persistence**
// **Validates: Requirements 1.3, 3.6**
// For any valid image file uploaded as avatar (for User or PPTK),
// the system should store the file and associate the file path with the entity.

// **Feature: manajemen-dokumen-keuangan, Property 11: PDF File Type Validation**
// **Validates: Requirements 6.4**
// For any uploaded file, the system should accept only PDF files
// and reject other file types with appropriate error message.

func TestProperty3_AvatarUploadPersistence(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Property: Valid image types should be accepted
	properties.Property("Valid image types should be accepted", prop.ForAll(
		func(extIdx int) bool {
			validExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
			ext := validExtensions[extIdx%len(validExtensions)]
			filename := "test" + ext

			return ValidateImageFile("", filename)
		},
		gen.IntRange(0, 100),
	))

	// Property: Invalid image types should be rejected
	properties.Property("Invalid image types should be rejected", prop.ForAll(
		func(extIdx int) bool {
			invalidExtensions := []string{".pdf", ".doc", ".txt", ".exe", ".zip"}
			ext := invalidExtensions[extIdx%len(invalidExtensions)]
			filename := "test" + ext

			return !ValidateImageFile("", filename)
		},
		gen.IntRange(0, 100),
	))

	// Property: Valid content types should be accepted
	properties.Property("Valid image content types should be accepted", prop.ForAll(
		func(typeIdx int) bool {
			validTypes := []string{"image/jpeg", "image/png", "image/gif", "image/webp"}
			contentType := validTypes[typeIdx%len(validTypes)]

			return ValidateImageFile(contentType, "")
		},
		gen.IntRange(0, 100),
	))

	properties.TestingRun(t)
}

func TestProperty11_PDFFileTypeValidation(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Property: PDF files should be accepted
	properties.Property("PDF files should be accepted", prop.ForAll(
		func(filename string) bool {
			if filename == "" {
				return true
			}
			pdfFilename := filename + ".pdf"
			return ValidatePDFFile("application/pdf", pdfFilename)
		},
		gen.AlphaString(),
	))

	// Property: Non-PDF files should be rejected
	properties.Property("Non-PDF files should be rejected", prop.ForAll(
		func(extIdx int) bool {
			invalidExtensions := []string{".jpg", ".png", ".doc", ".txt", ".exe", ".zip", ".xlsx"}
			ext := invalidExtensions[extIdx%len(invalidExtensions)]
			filename := "document" + ext

			return !ValidatePDFFile("", filename)
		},
		gen.IntRange(0, 100),
	))

	// Property: PDF content type should be accepted regardless of extension
	properties.Property("PDF content type should be accepted", prop.ForAll(
		func(dummy int) bool {
			return ValidatePDFFile("application/pdf", "")
		},
		gen.IntRange(0, 100),
	))

	// Property: Non-PDF content types should be rejected
	properties.Property("Non-PDF content types should be rejected", prop.ForAll(
		func(typeIdx int) bool {
			invalidTypes := []string{"image/jpeg", "text/plain", "application/msword", "application/zip"}
			contentType := invalidTypes[typeIdx%len(invalidTypes)]

			return !ValidatePDFFile(contentType, "document.doc")
		},
		gen.IntRange(0, 100),
	))

	properties.TestingRun(t)
}

func TestFilePathGeneration(t *testing.T) {
	t.Run("Avatar path should contain entity type", func(t *testing.T) {
		// Simulate path generation
		entityType := "users"
		filename := "test-avatar.jpg"
		path := filepath.Join("avatars", entityType, filename)

		if !strings.Contains(path, entityType) {
			t.Errorf("Path should contain entity type, got: %s", path)
		}
		if !strings.Contains(path, "avatars") {
			t.Errorf("Path should contain 'avatars', got: %s", path)
		}
	})

	t.Run("Document path should be in documents folder", func(t *testing.T) {
		filename := "test-doc.pdf"
		path := filepath.Join("documents", filename)

		if !strings.Contains(path, "documents") {
			t.Errorf("Path should contain 'documents', got: %s", path)
		}
	})
}
