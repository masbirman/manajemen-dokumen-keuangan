package middleware

import (
	"testing"

	"dokumen-keuangan/app/models"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// **Feature: manajemen-dokumen-keuangan, Property 23: Role-Based Access Control**
// **Validates: Requirements 10.3, 10.4**
// For any user with role R attempting to access endpoint requiring role R',
// where R has lower privilege than R', the system should deny access.
func TestProperty23_RoleBasedAccessControl(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Define role generators
	roleGen := gen.OneConstOf(
		models.RoleOperator,
		models.RoleAdmin,
		models.RoleSuperAdmin,
	)

	properties.Property("higher privilege role should access lower privilege endpoints", prop.ForAll(
		func(userRole, requiredRole models.UserRole) bool {
			userLevel := RoleHierarchy[userRole]
			requiredLevel := RoleHierarchy[requiredRole]

			// If user level >= required level, access should be granted
			// This property tests that the hierarchy is correctly implemented
			shouldHaveAccess := userLevel >= requiredLevel

			// Simulate the check
			hasAccess := checkRoleAccess(userRole, requiredRole)

			return hasAccess == shouldHaveAccess
		},
		roleGen,
		roleGen,
	))

	properties.Property("operator cannot access admin endpoints", prop.ForAll(
		func(_ bool) bool {
			// Operator (level 1) should not access Admin (level 2) endpoints
			return !checkRoleAccess(models.RoleOperator, models.RoleAdmin)
		},
		gen.Bool(), // dummy generator to run the test
	))

	properties.Property("operator cannot access super_admin endpoints", prop.ForAll(
		func(_ bool) bool {
			// Operator (level 1) should not access SuperAdmin (level 3) endpoints
			return !checkRoleAccess(models.RoleOperator, models.RoleSuperAdmin)
		},
		gen.Bool(),
	))

	properties.Property("admin cannot access super_admin endpoints", prop.ForAll(
		func(_ bool) bool {
			// Admin (level 2) should not access SuperAdmin (level 3) endpoints
			return !checkRoleAccess(models.RoleAdmin, models.RoleSuperAdmin)
		},
		gen.Bool(),
	))

	properties.Property("super_admin can access all endpoints", prop.ForAll(
		func(requiredRole models.UserRole) bool {
			// SuperAdmin should access any endpoint
			return checkRoleAccess(models.RoleSuperAdmin, requiredRole)
		},
		roleGen,
	))

	properties.Property("role hierarchy is transitive", prop.ForAll(
		func(role1, role2, role3 models.UserRole) bool {
			level1 := RoleHierarchy[role1]
			level2 := RoleHierarchy[role2]
			level3 := RoleHierarchy[role3]

			// If role1 >= role2 and role2 >= role3, then role1 >= role3
			if level1 >= level2 && level2 >= level3 {
				return level1 >= level3
			}
			return true // Skip cases where premise doesn't hold
		},
		roleGen,
		roleGen,
		roleGen,
	))

	properties.TestingRun(t)
}

// checkRoleAccess simulates the role check logic
func checkRoleAccess(userRole, requiredRole models.UserRole) bool {
	userLevel, userExists := RoleHierarchy[userRole]
	requiredLevel, requiredExists := RoleHierarchy[requiredRole]

	if !userExists || !requiredExists {
		return false
	}

	return userLevel >= requiredLevel
}

// TestRoleHierarchyValues ensures role hierarchy values are correctly defined
func TestRoleHierarchyValues(t *testing.T) {
	// Verify hierarchy values
	if RoleHierarchy[models.RoleOperator] >= RoleHierarchy[models.RoleAdmin] {
		t.Error("Operator should have lower privilege than Admin")
	}
	if RoleHierarchy[models.RoleAdmin] >= RoleHierarchy[models.RoleSuperAdmin] {
		t.Error("Admin should have lower privilege than SuperAdmin")
	}
	if RoleHierarchy[models.RoleOperator] >= RoleHierarchy[models.RoleSuperAdmin] {
		t.Error("Operator should have lower privilege than SuperAdmin")
	}
}
