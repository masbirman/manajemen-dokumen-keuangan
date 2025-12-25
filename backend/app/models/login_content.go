package models

import (
	"time"

	"github.com/google/uuid"
)

// LoginContent represents scheduled content for the login page
type LoginContent struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Title       string    `gorm:"size:255;not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	ImageURL    string    `gorm:"size:500" json:"image_url"`
	ImageWidth  int       `gorm:"default:400" json:"image_width"`
	TitleSize   int       `gorm:"default:28" json:"title_size"`
	DescSize    int       `gorm:"default:16" json:"desc_size"`
	StartDate   time.Time `gorm:"type:date;not null" json:"start_date"`
	EndDate     time.Time `gorm:"type:date;not null" json:"end_date"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName returns the table name for LoginContent
func (LoginContent) TableName() string {
	return "login_contents"
}
