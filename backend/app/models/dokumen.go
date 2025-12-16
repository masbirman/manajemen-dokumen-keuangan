package models

import (
	"time"

	"github.com/google/uuid"
)

// Dokumen represents a financial document
type Dokumen struct {
	ID              uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	NomorDokumen    string     `gorm:"size:255" json:"nomor_dokumen"`
	TanggalDokumen  *time.Time `gorm:"type:date" json:"tanggal_dokumen"`
	UnitKerjaID     uuid.UUID  `gorm:"type:uuid;not null" json:"unit_kerja_id"`
	PPTKID          uuid.UUID  `gorm:"type:uuid;not null" json:"pptk_id"`
	JenisDokumenID  uuid.UUID  `gorm:"type:uuid;not null" json:"jenis_dokumen_id"`
	SumberDanaID    uuid.UUID  `gorm:"type:uuid;not null" json:"sumber_dana_id"`
	Nilai           float64    `gorm:"type:decimal(15,2);not null" json:"nilai"`
	Uraian          string     `gorm:"type:text;not null" json:"uraian"`
	FilePath        string     `gorm:"not null;size:500" json:"file_path"`
	CreatedBy       uuid.UUID  `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt       time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	UnitKerja    UnitKerja    `gorm:"foreignKey:UnitKerjaID" json:"unit_kerja,omitempty"`
	PPTK         PPTK         `gorm:"foreignKey:PPTKID" json:"pptk,omitempty"`
	JenisDokumen JenisDokumen `gorm:"foreignKey:JenisDokumenID" json:"jenis_dokumen,omitempty"`
	SumberDana   SumberDana   `gorm:"foreignKey:SumberDanaID" json:"sumber_dana,omitempty"`
	Creator      User         `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

// TableName returns the table name for Dokumen
func (Dokumen) TableName() string {
	return "dokumen"
}
