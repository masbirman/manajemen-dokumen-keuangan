package models

import (
	"time"

	"github.com/google/uuid"
)

// UserPPTK represents the many-to-many relationship between users and PPTK
type UserPPTK struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	PPTKID    uuid.UUID `gorm:"type:uuid;not null" json:"pptk_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	// Relations
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	PPTK *PPTK `gorm:"foreignKey:PPTKID" json:"pptk,omitempty"`
}

// TableName returns the table name for UserPPTK
func (UserPPTK) TableName() string {
	return "user_pptk"
}
