package repositories

import (
	"dokumen-keuangan/app/models"
	"dokumen-keuangan/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SettingRepository handles setting data access
type SettingRepository struct {
	db *gorm.DB
}

// NewSettingRepository creates a new SettingRepository instance
func NewSettingRepository() *SettingRepository {
	return &SettingRepository{
		db: database.GetDB(),
	}
}

// GetAll retrieves all settings
func (r *SettingRepository) GetAll() ([]models.Setting, error) {
	var settings []models.Setting
	err := r.db.Order("key ASC").Find(&settings).Error
	return settings, err
}

// GetByKey retrieves a setting by key
func (r *SettingRepository) GetByKey(key string) (*models.Setting, error) {
	var setting models.Setting
	err := r.db.Where("key = ?", key).First(&setting).Error
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

// Upsert creates or updates a setting
func (r *SettingRepository) Upsert(key string, value *string) (*models.Setting, error) {
	var setting models.Setting
	err := r.db.Where("key = ?", key).First(&setting).Error
	
	if err == gorm.ErrRecordNotFound {
		// Create new
		setting = models.Setting{
			ID:    uuid.New(),
			Key:   key,
			Value: value,
		}
		if err := r.db.Create(&setting).Error; err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	} else {
		// Update existing
		setting.Value = value
		if err := r.db.Save(&setting).Error; err != nil {
			return nil, err
		}
	}
	
	return &setting, nil
}

// Delete deletes a setting by key
func (r *SettingRepository) Delete(key string) error {
	return r.db.Where("key = ?", key).Delete(&models.Setting{}).Error
}
