package models

import (
	"time"

	"github.com/google/uuid"
)

// Petunjuk represents guidance/information content
type Petunjuk struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Judul     string    `gorm:"not null;size:255" json:"judul"`
	Konten    string    `gorm:"type:text;not null" json:"konten"`
	Halaman   string    `gorm:"not null;size:100" json:"halaman"`
	Urutan    int       `gorm:"default:0" json:"urutan"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName returns the table name for Petunjuk
func (Petunjuk) TableName() string {
	return "petunjuk"
}

// Available halaman values
const (
	HalamanInputDokumen = "input_dokumen"
	HalamanListDokumen  = "list_dokumen"
	HalamanDashboard    = "dashboard"
)
