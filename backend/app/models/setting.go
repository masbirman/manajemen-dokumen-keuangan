package models

import (
	"time"

	"github.com/google/uuid"
)

// Setting represents a system configuration setting
type Setting struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Key       string    `gorm:"uniqueIndex;not null;size:100" json:"key"`
	Value     *string   `gorm:"type:text" json:"value"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName returns the table name for Setting
func (Setting) TableName() string {
	return "settings"
}
