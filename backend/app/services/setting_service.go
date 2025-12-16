package services

import (
	"errors"

	"dokumen-keuangan/app/models"
	"dokumen-keuangan/app/repositories"
)

var (
	ErrSettingNotFound = errors.New("setting not found")
	ErrSettingKeyRequired = errors.New("setting key is required")
)

// SettingService handles setting business logic
type SettingService struct {
	repo *repositories.SettingRepository
}

// NewSettingService creates a new SettingService instance
func NewSettingService() *SettingService {
	return &SettingService{
		repo: repositories.NewSettingRepository(),
	}
}

// UpdateSettingInput represents input for updating settings
type UpdateSettingInput struct {
	Settings map[string]*string `json:"settings"`
}

// GetAll retrieves all settings
func (s *SettingService) GetAll() ([]models.Setting, error) {
	return s.repo.GetAll()
}

// GetByKey retrieves a setting by key
func (s *SettingService) GetByKey(key string) (*models.Setting, error) {
	if key == "" {
		return nil, ErrSettingKeyRequired
	}
	setting, err := s.repo.GetByKey(key)
	if err != nil {
		return nil, ErrSettingNotFound
	}
	return setting, nil
}

// Update updates multiple settings
func (s *SettingService) Update(input *UpdateSettingInput) ([]models.Setting, error) {
	if input == nil || len(input.Settings) == 0 {
		return s.repo.GetAll()
	}

	for key, value := range input.Settings {
		if key == "" {
			continue
		}
		if _, err := s.repo.Upsert(key, value); err != nil {
			return nil, err
		}
	}

	return s.repo.GetAll()
}

// GetSettingsMap returns settings as a map
func (s *SettingService) GetSettingsMap() (map[string]*string, error) {
	settings, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	result := make(map[string]*string)
	for _, setting := range settings {
		result[setting.Key] = setting.Value
	}
	return result, nil
}
