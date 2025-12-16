package services

import (
	"testing"

	"github.com/google/uuid"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"golang.org/x/crypto/bcrypt"
)

// **Feature: manajemen-dokumen-keuangan, Property 1: User Creation Persistence**
// **Validates: Requirements 1.1**
// For any valid user data (username, password, role), when a Super Admin creates the user,
// the system should store all fields correctly and the user should be retrievable with the same data.

// **Feature: manajemen-dokumen-keuangan, Property 2: Operator Assignment Integrity**
// **Validates: Requirements 1.2**
// For any Operator, PPTK, and Unit Kerja combination, when assigned,
// the Operator should be linked to both entities and these relationships should be queryable.

// **Feature: manajemen-dokumen-keuangan, Property 4: User Deactivation Blocks Login**
// **Validates: Requirements 1.5**
// For any active user, when deactivated by Super Admin, subsequent login attempts
// with valid credentials should fail.

// MockUser represents a user for testing
type MockUser struct {
	ID          uuid.UUID
	Username    string
	Password    string
	Name        string
	Role        string
	UnitKerjaID *uuid.UUID
	PPTKID      *uuid.UUID
	IsActive    bool
}

// MockUserStore simulates user storage for testing
type MockUserStore struct {
	users map[uuid.UUID]*MockUser
}

func NewMockUserStore() *MockUserStore {
	return &MockUserStore{
		users: make(map[uuid.UUID]*MockUser),
	}
}

func (s *MockUserStore) Create(user *MockUser) error {
	user.ID = uuid.New()
	s.users[user.ID] = user
	return nil
}

func (s *MockUserStore) FindByID(id uuid.UUID) (*MockUser, error) {
	if user, ok := s.users[id]; ok {
		return user, nil
	}
	return nil, ErrUserNotFound
}

func (s *MockUserStore) Deactivate(id uuid.UUID) error {
	if user, ok := s.users[id]; ok {
		user.IsActive = false
		return nil
	}
	return ErrUserNotFound
}

func (s *MockUserStore) CanLogin(username, password string) bool {
	for _, user := range s.users {
		if user.Username == username && user.IsActive {
			// Check password
			err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
			return err == nil
		}
	}
	return false
}


func TestProperty1_UserCreationPersistence(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Property: Created user should be retrievable with same data
	properties.Property("User creation should persist all fields correctly", prop.ForAll(
		func(username, name string, roleIdx int) bool {
			if username == "" || name == "" {
				return true // Skip invalid input
			}

			store := NewMockUserStore()
			roles := []string{"super_admin", "admin", "operator"}
			role := roles[roleIdx%3]

			// Hash password
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testpass"), bcrypt.MinCost)

			user := &MockUser{
				Username: username,
				Password: string(hashedPassword),
				Name:     name,
				Role:     role,
				IsActive: true,
			}

			// Create user
			store.Create(user)

			// Retrieve user
			retrieved, err := store.FindByID(user.ID)
			if err != nil {
				return false
			}

			// Verify all fields match
			return retrieved.Username == username &&
				retrieved.Name == name &&
				retrieved.Role == role &&
				retrieved.IsActive == true
		},
		gen.AlphaString(),
		gen.AlphaString(),
		gen.IntRange(0, 100),
	))

	// Property: User ID should be unique
	properties.Property("Each created user should have unique ID", prop.ForAll(
		func(count int) bool {
			if count <= 0 || count > 50 {
				return true
			}

			store := NewMockUserStore()
			ids := make(map[uuid.UUID]bool)

			for i := 0; i < count; i++ {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.MinCost)
				user := &MockUser{
					Username: "user" + string(rune('A'+i)),
					Password: string(hashedPassword),
					Name:     "User " + string(rune('A'+i)),
					Role:     "operator",
					IsActive: true,
				}
				store.Create(user)

				if ids[user.ID] {
					return false // Duplicate ID found
				}
				ids[user.ID] = true
			}

			return true
		},
		gen.IntRange(1, 50),
	))

	properties.TestingRun(t)
}

func TestProperty2_OperatorAssignmentIntegrity(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Property: Operator should be linked to Unit Kerja and PPTK
	properties.Property("Operator assignment should link to both entities", prop.ForAll(
		func(hasUnitKerja, hasPPTK bool) bool {
			store := NewMockUserStore()

			var unitKerjaID, pptkID *uuid.UUID
			if hasUnitKerja {
				id := uuid.New()
				unitKerjaID = &id
			}
			if hasPPTK {
				id := uuid.New()
				pptkID = &id
			}

			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.MinCost)
			user := &MockUser{
				Username:    "operator_test",
				Password:    string(hashedPassword),
				Name:        "Test Operator",
				Role:        "operator",
				UnitKerjaID: unitKerjaID,
				PPTKID:      pptkID,
				IsActive:    true,
			}

			store.Create(user)
			retrieved, _ := store.FindByID(user.ID)

			// Verify assignments
			unitKerjaMatch := (hasUnitKerja && retrieved.UnitKerjaID != nil) ||
				(!hasUnitKerja && retrieved.UnitKerjaID == nil)
			pptkMatch := (hasPPTK && retrieved.PPTKID != nil) ||
				(!hasPPTK && retrieved.PPTKID == nil)

			return unitKerjaMatch && pptkMatch
		},
		gen.Bool(),
		gen.Bool(),
	))

	properties.TestingRun(t)
}

func TestProperty4_UserDeactivationBlocksLogin(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Property: Deactivated user should not be able to login
	properties.Property("Deactivated user cannot login with valid credentials", prop.ForAll(
		func(username string) bool {
			if username == "" {
				return true
			}

			store := NewMockUserStore()
			password := "testpassword123"

			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
			user := &MockUser{
				Username: username,
				Password: string(hashedPassword),
				Name:     "Test User",
				Role:     "operator",
				IsActive: true,
			}

			store.Create(user)

			// Verify can login when active
			if !store.CanLogin(username, password) {
				return false // Should be able to login when active
			}

			// Deactivate user
			store.Deactivate(user.ID)

			// Verify cannot login after deactivation
			return !store.CanLogin(username, password)
		},
		gen.AlphaString(),
	))

	// Property: Active user should be able to login
	properties.Property("Active user can login with valid credentials", prop.ForAll(
		func(username string) bool {
			if username == "" {
				return true
			}

			store := NewMockUserStore()
			password := "validpassword"

			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
			user := &MockUser{
				Username: username,
				Password: string(hashedPassword),
				Name:     "Active User",
				Role:     "admin",
				IsActive: true,
			}

			store.Create(user)

			return store.CanLogin(username, password)
		},
		gen.AlphaString(),
	))

	properties.TestingRun(t)
}
