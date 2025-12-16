package models

import (
	"time"

	"github.com/google/uuid"
)

// UserRole represents the role of a user
type UserRole string

const (
	RoleSuperAdmin UserRole = "super_admin"
	RoleAdmin      UserRole = "admin"
	RoleOperator   UserRole = "operator"
)

// User represents a system user
type User struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Username    string     `gorm:"uniqueIndex;not null;size:100" json:"username"`
	Password    string     `gorm:"not null;size:255" json:"-"`
	Name        string     `gorm:"not null;size:255" json:"name"`
	Role        UserRole   `gorm:"not null;size:20" json:"role"`
	UnitKerjaID *uuid.UUID `gorm:"type:uuid" json:"unit_kerja_id"`
	PPTKID      *uuid.UUID `gorm:"type:uuid" json:"pptk_id"`
	AvatarPath  *string    `gorm:"size:500" json:"avatar_path"`
	IsActive    bool       `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	UnitKerja *UnitKerja  `gorm:"foreignKey:UnitKerjaID" json:"unit_kerja,omitempty"`
	PPTK      *PPTK       `gorm:"foreignKey:PPTKID" json:"pptk,omitempty"`
	PPTKList  []UserPPTK  `gorm:"foreignKey:UserID" json:"pptk_list,omitempty"`
}

// TableName returns the table name for User
func (User) TableName() string {
	return "users"
}

// IsAdminOrAbove checks if user has admin or super_admin role
func (u *User) IsAdminOrAbove() bool {
	return u.Role == RoleAdmin || u.Role == RoleSuperAdmin
}

// IsSuperAdmin checks if user has super_admin role
func (u *User) IsSuperAdmin() bool {
	return u.Role == RoleSuperAdmin
}
