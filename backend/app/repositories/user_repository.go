package repositories

import (
	"dokumen-keuangan/app/models"
	"dokumen-keuangan/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRepository handles user data access
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository instance
func NewUserRepository() *UserRepository {
	return &UserRepository{
		db: database.GetDB(),
	}
}

// FindByUsername finds a user by username
func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByID finds a user by ID
func (r *UserRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.Preload("UnitKerja").Preload("PPTK").Preload("PPTKList").Preload("PPTKList.PPTK").Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindActiveByUsername finds an active user by username
func (r *UserRepository) FindActiveByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ? AND is_active = ?", username, true).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}


// Create creates a new user
func (r *UserRepository) Create(user *models.User) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	return r.db.Create(user).Error
}

// GetAll retrieves all users with pagination
func (r *UserRepository) GetAll(page, pageSize int, search string) (*PaginationResult, error) {
	return r.GetAllWithFilter(page, pageSize, search, "")
}

// GetAllWithFilter retrieves all users with pagination and role filter
func (r *UserRepository) GetAllWithFilter(page, pageSize int, search, role string) (*PaginationResult, error) {
	var users []models.User
	var total int64

	query := r.db.Model(&models.User{}).Preload("UnitKerja").Preload("PPTK").Preload("PPTKList").Preload("PPTKList.PPTK")

	// Apply search filter if provided
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("username ILIKE ? OR name ILIKE ?", searchPattern, searchPattern)
	}

	// Apply role filter if provided
	if role != "" {
		query = query.Where("role = ?", role)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Fetch paginated results
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &PaginationResult{
		Data:       users,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// Update updates a user
func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// Delete deletes a user by ID (soft delete by setting is_active = false)
func (r *UserRepository) Delete(id uuid.UUID) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Update("is_active", false).Error
}

// ExistsByUsername checks if a user with the given username exists (excluding a specific ID)
func (r *UserRepository) ExistsByUsername(username string, excludeID *uuid.UUID) (bool, error) {
	var count int64
	query := r.db.Model(&models.User{}).Where("username = ?", username)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}


// GetUserPPTKs retrieves all PPTK IDs assigned to a user
func (r *UserRepository) GetUserPPTKs(userID uuid.UUID) ([]uuid.UUID, error) {
	var userPPTKs []models.UserPPTK
	err := r.db.Where("user_id = ?", userID).Find(&userPPTKs).Error
	if err != nil {
		return nil, err
	}
	
	pptkIDs := make([]uuid.UUID, len(userPPTKs))
	for i, up := range userPPTKs {
		pptkIDs[i] = up.PPTKID
	}
	return pptkIDs, nil
}

// SetUserPPTKs sets the PPTK assignments for a user (replaces existing)
func (r *UserRepository) SetUserPPTKs(userID uuid.UUID, pptkIDs []uuid.UUID) error {
	// Delete existing assignments
	if err := r.db.Where("user_id = ?", userID).Delete(&models.UserPPTK{}).Error; err != nil {
		return err
	}
	
	// Create new assignments
	for _, pptkID := range pptkIDs {
		userPPTK := models.UserPPTK{
			ID:     uuid.New(),
			UserID: userID,
			PPTKID: pptkID,
		}
		if err := r.db.Create(&userPPTK).Error; err != nil {
			return err
		}
	}
	
	return nil
}

// DeleteUserPPTKs deletes all PPTK assignments for a user
func (r *UserRepository) DeleteUserPPTKs(userID uuid.UUID) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.UserPPTK{}).Error
}
