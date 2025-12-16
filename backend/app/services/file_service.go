package services

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

var (
	ErrInvalidFileType    = errors.New("invalid file type")
	ErrFileTooLarge       = errors.New("file too large")
	ErrFileUploadFailed   = errors.New("file upload failed")
	ErrFileNotFound       = errors.New("file not found")
)

// AllowedImageTypes contains allowed image MIME types
var AllowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

// AllowedDocumentTypes contains allowed document MIME types
var AllowedDocumentTypes = map[string]bool{
	"application/pdf": true,
}

// FileService handles file upload operations
type FileService struct {
	basePath string
}

// NewFileService creates a new FileService instance
func NewFileService() *FileService {
	basePath := os.Getenv("STORAGE_PATH")
	if basePath == "" {
		basePath = "storage/app/public"
	}
	return &FileService{
		basePath: basePath,
	}
}

// UploadAvatar uploads an avatar image
func (s *FileService) UploadAvatar(file *multipart.FileHeader, entityType string, entityID uuid.UUID) (string, error) {
	// Validate file type
	contentType := file.Header.Get("Content-Type")
	if !AllowedImageTypes[contentType] {
		// Also check by extension
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" && ext != ".webp" {
			return "", ErrInvalidFileType
		}
	}

	// Validate file size (max 5MB)
	if file.Size > 5*1024*1024 {
		return "", ErrFileTooLarge
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := uuid.New().String() + ext

	// Create directory path
	dirPath := filepath.Join(s.basePath, "avatars", entityType)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return "", ErrFileUploadFailed
	}

	// Full file path
	filePath := filepath.Join(dirPath, filename)

	// Open source file
	src, err := file.Open()
	if err != nil {
		return "", ErrFileUploadFailed
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return "", ErrFileUploadFailed
	}
	defer dst.Close()

	// Copy file
	if _, err := io.Copy(dst, src); err != nil {
		return "", ErrFileUploadFailed
	}

	// Return relative path for storage
	return filepath.Join("avatars", entityType, filename), nil
}


// UploadDocument uploads a PDF document
func (s *FileService) UploadDocument(file *multipart.FileHeader) (string, error) {
	// Validate file type
	contentType := file.Header.Get("Content-Type")
	if !AllowedDocumentTypes[contentType] {
		// Also check by extension
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".pdf" {
			return "", ErrInvalidFileType
		}
	}

	// Validate file size (max 20MB)
	if file.Size > 20*1024*1024 {
		return "", ErrFileTooLarge
	}

	// Generate unique filename
	filename := uuid.New().String() + ".pdf"

	// Create directory path
	dirPath := filepath.Join(s.basePath, "documents")
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return "", ErrFileUploadFailed
	}

	// Full file path
	filePath := filepath.Join(dirPath, filename)

	// Open source file
	src, err := file.Open()
	if err != nil {
		return "", ErrFileUploadFailed
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return "", ErrFileUploadFailed
	}
	defer dst.Close()

	// Copy file
	if _, err := io.Copy(dst, src); err != nil {
		return "", ErrFileUploadFailed
	}

	// Return relative path for storage
	return filepath.Join("documents", filename), nil
}

// GetFilePath returns the full path to a file
func (s *FileService) GetFilePath(relativePath string) string {
	return filepath.Join(s.basePath, relativePath)
}

// FileExists checks if a file exists
func (s *FileService) FileExists(relativePath string) bool {
	fullPath := s.GetFilePath(relativePath)
	_, err := os.Stat(fullPath)
	return err == nil
}

// DeleteFile deletes a file
func (s *FileService) DeleteFile(relativePath string) error {
	fullPath := s.GetFilePath(relativePath)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return ErrFileNotFound
	}
	return os.Remove(fullPath)
}

// ValidateImageFile validates if a file is a valid image
func ValidateImageFile(contentType, filename string) bool {
	if AllowedImageTypes[contentType] {
		return true
	}
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" || ext == ".webp"
}

// ValidatePDFFile validates if a file is a valid PDF
func ValidatePDFFile(contentType, filename string) bool {
	if AllowedDocumentTypes[contentType] {
		return true
	}
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".pdf"
}
