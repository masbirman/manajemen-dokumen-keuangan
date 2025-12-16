package middleware

import (
	"dokumen-keuangan/app/models"

	"github.com/gofiber/fiber/v2"
)

// RoleHierarchy defines the privilege level for each role
var RoleHierarchy = map[models.UserRole]int{
	models.RoleOperator:   1,
	models.RoleAdmin:      2,
	models.RoleSuperAdmin: 3,
}

// RequireRole creates a middleware that checks if user has the required role or higher
func RequireRole(requiredRole models.UserRole) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user role from context (set by AuthMiddleware)
		userRoleInterface := c.Locals("userRole")
		if userRoleInterface == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authentication required",
			})
		}

		userRole := userRoleInterface.(models.UserRole)
		
		// Check role hierarchy
		userLevel, userExists := RoleHierarchy[userRole]
		requiredLevel, requiredExists := RoleHierarchy[requiredRole]

		if !userExists || !requiredExists {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Invalid role configuration",
			})
		}

		if userLevel < requiredLevel {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied. Required role: " + string(requiredRole),
			})
		}

		return c.Next()
	}
}

// RequireExactRole creates a middleware that checks if user has exactly the specified role
func RequireExactRole(role models.UserRole) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRoleInterface := c.Locals("userRole")
		if userRoleInterface == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authentication required",
			})
		}

		userRole := userRoleInterface.(models.UserRole)

		if userRole != role {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied. Required role: " + string(role),
			})
		}

		return c.Next()
	}
}

// RequireAnyRole creates a middleware that checks if user has any of the specified roles
func RequireAnyRole(roles ...models.UserRole) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRoleInterface := c.Locals("userRole")
		if userRoleInterface == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authentication required",
			})
		}

		userRole := userRoleInterface.(models.UserRole)

		for _, role := range roles {
			if userRole == role {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied. Insufficient permissions",
		})
	}
}

// IsSuperAdmin checks if the current user is a super admin
func IsSuperAdmin(c *fiber.Ctx) bool {
	userRoleInterface := c.Locals("userRole")
	if userRoleInterface == nil {
		return false
	}
	return userRoleInterface.(models.UserRole) == models.RoleSuperAdmin
}

// IsAdminOrAbove checks if the current user is admin or super admin
func IsAdminOrAbove(c *fiber.Ctx) bool {
	userRoleInterface := c.Locals("userRole")
	if userRoleInterface == nil {
		return false
	}
	userRole := userRoleInterface.(models.UserRole)
	return userRole == models.RoleAdmin || userRole == models.RoleSuperAdmin
}
