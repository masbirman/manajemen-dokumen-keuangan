package models

import (
	"time"

	"github.com/google/uuid"
)

// SumberDana represents a funding source category
type SumberDana struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Kode      string    `gorm:"uniqueIndex;not null;size:50" json:"kode"`
	Nama      string    `gorm:"not null;size:255" json:"nama"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName returns the table name for SumberDana
func (SumberDana) TableName() string {
	return "sumber_dana"
}
