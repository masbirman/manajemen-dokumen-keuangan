package middleware

import (
	"fmt"
	"strings"

	"dokumen-keuangan/app/models"
	"dokumen-keuangan/app/services"
	"dokumen-keuangan/config"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// AuthMiddleware validates JWT tokens and attaches user to context
func AuthMiddleware(cfg *config.Config) fiber.Handler {
	authService := services.NewAuthService(cfg)

	return func(c *fiber.Ctx) error {
		// Get Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authentication required",
			})
		}

		// Extract token from "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header format",
			})
		}

		// Validate token
		claims, err := authService.ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		// Check if it's an access token
		if claims.Type != "access" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token type",
			})
		}

		// Attach user info to context
		c.Locals("userID", claims.UserID)
		c.Locals("username", claims.Username)
		c.Locals("userRole", claims.Role)

		return c.Next()
	}
}

// GetUserIDFromContext retrieves user ID from fiber context
func GetUserIDFromContext(c *fiber.Ctx) string {
	if userID := c.Locals("userID"); userID != nil {
		switch id := userID.(type) {
		case string:
			return id
		case uuid.UUID:
			return id.String()
		default:
			return fmt.Sprintf("%v", id)
		}
	}
	return ""
}

// GetUserRoleFromContext retrieves user role from fiber context
func GetUserRoleFromContext(c *fiber.Ctx) string {
	if role := c.Locals("userRole"); role != nil {
		// Handle both models.UserRole and string types
		switch r := role.(type) {
		case string:
			return r
		case models.UserRole:
			return string(r)
		default:
			return fmt.Sprintf("%v", r)
		}
	}
	return ""
}
