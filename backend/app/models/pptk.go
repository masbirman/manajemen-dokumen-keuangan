package models

import (
	"time"

	"github.com/google/uuid"
)

// PPTK represents a Pejabat Pelaksana Teknis Kegiatan
type PPTK struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	NIP         string     `gorm:"column:nip;uniqueIndex;not null;size:50" json:"nip"`
	Nama        string     `gorm:"not null;size:255" json:"nama"`
	Jabatan     string     `gorm:"size:255" json:"jabatan"`
	UnitKerjaID uuid.UUID  `gorm:"type:uuid;not null" json:"unit_kerja_id"`
	AvatarPath  *string    `gorm:"size:500" json:"avatar_path"`
	IsActive    bool       `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	UnitKerja UnitKerja `gorm:"foreignKey:UnitKerjaID" json:"unit_kerja,omitempty"`
}

// TableName returns the table name for PPTK
func (PPTK) TableName() string {
	return "pptk"
}
