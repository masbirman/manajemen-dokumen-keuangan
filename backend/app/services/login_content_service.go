package services

import (
	"errors"
	"time"

	"dokumen-keuangan/app/models"
	"dokumen-keuangan/database"

	"github.com/google/uuid"
)

var (
	ErrLoginContentNotFound = errors.New("login content not found")
)

// LoginContentService handles login content business logic
type LoginContentService struct{}

// NewLoginContentService creates a new LoginContentService instance
func NewLoginContentService() *LoginContentService {
	return &LoginContentService{}
}

// CreateLoginContentInput represents input for creating login content
type CreateLoginContentInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageWidth  int    `json:"image_width"`
	TitleSize   int    `json:"title_size"`
	DescSize    int    `json:"desc_size"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
}

// GetAll retrieves all login contents
func (s *LoginContentService) GetAll() ([]models.LoginContent, error) {
	var contents []models.LoginContent
	err := database.DB.Order("start_date DESC").Find(&contents).Error
	return contents, err
}

// GetByID retrieves a login content by ID
func (s *LoginContentService) GetByID(id uuid.UUID) (*models.LoginContent, error) {
	var content models.LoginContent
	if err := database.DB.First(&content, "id = ?", id).Error; err != nil {
		return nil, ErrLoginContentNotFound
	}
	return &content, nil
}

// GetActive retrieves the currently active login content
func (s *LoginContentService) GetActive() (*models.LoginContent, error) {
	var content models.LoginContent
	today := time.Now().Format("2006-01-02")

	err := database.DB.Where("is_active = ? AND start_date <= ? AND end_date >= ?", true, today, today).
		Order("start_date DESC").
		First(&content).Error

	if err != nil {
		return nil, ErrLoginContentNotFound
	}
	return &content, nil
}

// Create creates a new login content
func (s *LoginContentService) Create(input *CreateLoginContentInput) (*models.LoginContent, error) {
	startDate, _ := time.Parse("2006-01-02", input.StartDate)
	endDate, _ := time.Parse("2006-01-02", input.EndDate)

	content := models.LoginContent{
		Title:       input.Title,
		Description: input.Description,
		ImageWidth:  input.ImageWidth,
		TitleSize:   input.TitleSize,
		DescSize:    input.DescSize,
		StartDate:   startDate,
		EndDate:     endDate,
		IsActive:    true,
	}

	// Set defaults
	if content.ImageWidth == 0 {
		content.ImageWidth = 400
	}
	if content.TitleSize == 0 {
		content.TitleSize = 28
	}
	if content.DescSize == 0 {
		content.DescSize = 16
	}

	if err := database.DB.Create(&content).Error; err != nil {
		return nil, err
	}
	return &content, nil
}

// Update updates a login content
func (s *LoginContentService) Update(id uuid.UUID, input *CreateLoginContentInput) (*models.LoginContent, error) {
	content, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	startDate, _ := time.Parse("2006-01-02", input.StartDate)
	endDate, _ := time.Parse("2006-01-02", input.EndDate)

	content.Title = input.Title
	content.Description = input.Description
	content.ImageWidth = input.ImageWidth
	content.TitleSize = input.TitleSize
	content.DescSize = input.DescSize
	content.StartDate = startDate
	content.EndDate = endDate

	if err := database.DB.Save(content).Error; err != nil {
		return nil, err
	}
	return content, nil
}

// Delete deletes a login content
func (s *LoginContentService) Delete(id uuid.UUID) error {
	content, err := s.GetByID(id)
	if err != nil {
		return err
	}
	return database.DB.Delete(content).Error
}

// UpdateImageURL updates the image URL for a login content
func (s *LoginContentService) UpdateImageURL(id uuid.UUID, imageURL string) error {
	return database.DB.Model(&models.LoginContent{}).Where("id = ?", id).Update("image_url", imageURL).Error
}
