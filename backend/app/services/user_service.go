package services

import (
	"errors"
	"math/rand"
	"strings"
	"time"

	"dokumen-keuangan/app/models"
	"dokumen-keuangan/app/repositories"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrUserUsernameExists   = errors.New("username already exists")
	ErrUserUsernameRequired = errors.New("username is required")
	ErrUserPasswordRequired = errors.New("password is required")
	ErrUserNameRequired     = errors.New("name is required")
	ErrUserRoleRequired     = errors.New("role is required")
	ErrUserInvalidRole      = errors.New("invalid role")
)

// UserService handles user business logic
type UserService struct {
	repo          *repositories.UserRepository
	unitKerjaRepo *repositories.UnitKerjaRepository
	pptkRepo      *repositories.PPTKRepository
}

// NewUserService creates a new UserService instance
func NewUserService() *UserService {
	return &UserService{
		repo:          repositories.NewUserRepository(),
		unitKerjaRepo: repositories.NewUnitKerjaRepository(),
		pptkRepo:      repositories.NewPPTKRepository(),
	}
}

// CreateUserInput represents input for creating a user
type CreateUserInput struct {
	Username    string      `json:"username"`
	Password    string      `json:"password"`
	Name        string      `json:"name"`
	Role        string      `json:"role"`
	UnitKerjaID *uuid.UUID  `json:"unit_kerja_id"`
	PPTKID      *uuid.UUID  `json:"pptk_id"`
	PPTKIDs     []uuid.UUID `json:"pptk_ids"` // Multiple PPTK support
}

// UpdateUserInput represents input for updating a user
type UpdateUserInput struct {
	Username    *string     `json:"username"`
	Password    *string     `json:"password"`
	Name        *string     `json:"name"`
	Role        *string     `json:"role"`
	UnitKerjaID *uuid.UUID  `json:"unit_kerja_id"`
	PPTKID      *uuid.UUID  `json:"pptk_id"`
	PPTKIDs     []uuid.UUID `json:"pptk_ids"` // Multiple PPTK support
	AvatarPath  *string     `json:"avatar_path"`
	IsActive    *bool       `json:"is_active"`
}


// Validate validates the create input
func (i *CreateUserInput) Validate() map[string]string {
	errors := make(map[string]string)

	if strings.TrimSpace(i.Username) == "" {
		errors["username"] = "Username is required"
	}
	if strings.TrimSpace(i.Password) == "" {
		errors["password"] = "Password is required"
	}
	if strings.TrimSpace(i.Name) == "" {
		errors["name"] = "Name is required"
	}
	if strings.TrimSpace(i.Role) == "" {
		errors["role"] = "Role is required"
	} else if !isValidRole(i.Role) {
		errors["role"] = "Invalid role. Must be super_admin, admin, or operator"
	}

	return errors
}

func isValidRole(role string) bool {
	return role == string(models.RoleSuperAdmin) || role == string(models.RoleAdmin) || role == string(models.RoleOperator)
}

// Create creates a new user
func (s *UserService) Create(input *CreateUserInput) (*models.User, error) {
	// Validate input
	if validationErrors := input.Validate(); len(validationErrors) > 0 {
		if _, ok := validationErrors["username"]; ok {
			return nil, ErrUserUsernameRequired
		}
		if _, ok := validationErrors["password"]; ok {
			return nil, ErrUserPasswordRequired
		}
		if _, ok := validationErrors["name"]; ok {
			return nil, ErrUserNameRequired
		}
		if _, ok := validationErrors["role"]; ok {
			if input.Role == "" {
				return nil, ErrUserRoleRequired
			}
			return nil, ErrUserInvalidRole
		}
	}

	// Check if username already exists
	exists, err := s.repo.ExistsByUsername(strings.TrimSpace(input.Username), nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUserUsernameExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Validate Unit Kerja if provided
	if input.UnitKerjaID != nil {
		_, err := s.unitKerjaRepo.FindByID(*input.UnitKerjaID)
		if err != nil {
			return nil, errors.New("unit kerja not found")
		}
	}

	// Validate PPTK if provided
	if input.PPTKID != nil {
		_, err := s.pptkRepo.FindByID(*input.PPTKID)
		if err != nil {
			return nil, errors.New("PPTK not found")
		}
	}

	// Create user
	user := &models.User{
		ID:          uuid.New(),
		Username:    strings.TrimSpace(input.Username),
		Password:    string(hashedPassword),
		Name:        strings.TrimSpace(input.Name),
		Role:        models.UserRole(input.Role),
		UnitKerjaID: input.UnitKerjaID,
		PPTKID:      input.PPTKID,
		IsActive:    true,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	// Set multiple PPTKs if provided
	if len(input.PPTKIDs) > 0 {
		if err := s.repo.SetUserPPTKs(user.ID, input.PPTKIDs); err != nil {
			return nil, err
		}
	}

	// Reload with relationships
	return s.repo.FindByID(user.ID)
}

// GetByID retrieves a user by ID
func (s *UserService) GetByID(id uuid.UUID) (*models.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// GetAll retrieves all users with pagination
func (s *UserService) GetAll(page, pageSize int, search string) (*repositories.PaginationResult, error) {
	return s.GetAllWithFilter(page, pageSize, search, "")
}

// GetAllWithFilter retrieves all users with pagination and role filter
func (s *UserService) GetAllWithFilter(page, pageSize int, search, role string) (*repositories.PaginationResult, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	return s.repo.GetAllWithFilter(page, pageSize, search, role)
}


// Update updates a user
func (s *UserService) Update(id uuid.UUID, input *UpdateUserInput) (*models.User, error) {
	// Find existing user
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Update fields if provided
	if input.Username != nil {
		username := strings.TrimSpace(*input.Username)
		if username == "" {
			return nil, ErrUserUsernameRequired
		}
		// Check if new username already exists (excluding current record)
		exists, err := s.repo.ExistsByUsername(username, &id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrUserUsernameExists
		}
		user.Username = username
	}

	if input.Password != nil && *input.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*input.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashedPassword)
	}

	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		if name == "" {
			return nil, ErrUserNameRequired
		}
		user.Name = name
	}

	if input.Role != nil {
		if !isValidRole(*input.Role) {
			return nil, ErrUserInvalidRole
		}
		user.Role = models.UserRole(*input.Role)
	}

	if input.UnitKerjaID != nil {
		_, err := s.unitKerjaRepo.FindByID(*input.UnitKerjaID)
		if err != nil {
			return nil, errors.New("unit kerja not found")
		}
		user.UnitKerjaID = input.UnitKerjaID
	}

	if input.PPTKID != nil {
		_, err := s.pptkRepo.FindByID(*input.PPTKID)
		if err != nil {
			return nil, errors.New("PPTK not found")
		}
		user.PPTKID = input.PPTKID
	}

	if input.AvatarPath != nil {
		user.AvatarPath = input.AvatarPath
	}

	if input.IsActive != nil {
		user.IsActive = *input.IsActive
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	// Update multiple PPTKs if provided (even empty array to clear)
	if input.PPTKIDs != nil {
		if err := s.repo.SetUserPPTKs(user.ID, input.PPTKIDs); err != nil {
			return nil, err
		}
	}

	// Reload with relationships
	return s.repo.FindByID(user.ID)
}

// Delete deactivates a user by ID
func (s *UserService) Delete(id uuid.UUID) error {
	// Check if user exists
	_, err := s.repo.FindByID(id)
	if err != nil {
		return ErrUserNotFound
	}

	return s.repo.Delete(id)
}

// Deactivate deactivates a user (same as Delete but more explicit)
func (s *UserService) Deactivate(id uuid.UUID) error {
	return s.Delete(id)
}

// Activate activates a user
func (s *UserService) Activate(id uuid.UUID) (*models.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	user.IsActive = true
	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return s.repo.FindByID(user.ID)
}

// Count returns the total number of users
func (s *UserService) Count() (int64, error) {
	result, err := s.repo.GetAllWithFilter(1, 1, "", "")
	if err != nil {
		return 0, err
	}
	return result.Total, nil
}

// GenerateRandomPassword generates a random password with given length
func (s *UserService) GenerateRandomPassword(length int) string {
	const charset = "abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	password := make([]byte, length)
	for i := range password {
		password[i] = charset[r.Intn(len(charset))]
	}
	return string(password)
}
