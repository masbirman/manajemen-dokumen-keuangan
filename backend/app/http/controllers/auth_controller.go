package controllers

import (
	"strings"

	"dokumen-keuangan/app/services"
	"dokumen-keuangan/config"

	"github.com/gofiber/fiber/v2"
)

// AuthController handles authentication endpoints
type AuthController struct {
	authService *services.AuthService
}

// NewAuthController creates a new AuthController instance
func NewAuthController(cfg *config.Config) *AuthController {
	return &AuthController{
		authService: services.NewAuthService(cfg),
	}
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RefreshRequest represents the refresh token request body
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// Login handles user login
// POST /api/auth/login
func (c *AuthController) Login(ctx *fiber.Ctx) error {
	var req LoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if req.Username == "" || req.Password == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": fiber.Map{
				"username": condMsg(req.Username == "", "Username is required"),
				"password": condMsg(req.Password == "", "Password is required"),
			},
		})
	}

	// Attempt login
	tokenPair, user, err := c.authService.Login(req.Username, req.Password)
	if err != nil {
		switch err {
		case services.ErrInvalidCredentials:
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid username or password",
			})
		case services.ErrUserInactive:
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User account is inactive",
			})
		default:
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error",
			})
		}
	}

	return ctx.JSON(fiber.Map{
		"message": "Login successful",
		"data": fiber.Map{
			"user":  user,
			"token": tokenPair,
		},
	})
}

// Logout handles user logout
// POST /api/auth/logout
func (c *AuthController) Logout(ctx *fiber.Ctx) error {
	// In a stateless JWT implementation, logout is handled client-side
	// by removing the token. Server-side, we just acknowledge the request.
	// For a more robust implementation, you could maintain a token blacklist.
	return ctx.JSON(fiber.Map{
		"message": "Logout successful",
	})
}

// Refresh handles token refresh
// POST /api/auth/refresh
func (c *AuthController) Refresh(ctx *fiber.Ctx) error {
	var req RefreshRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.RefreshToken == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Refresh token is required",
		})
	}

	tokenPair, user, err := c.authService.RefreshTokens(req.RefreshToken)
	if err != nil {
		switch err {
		case services.ErrInvalidToken:
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired refresh token",
			})
		case services.ErrUserInactive:
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User account is inactive",
			})
		default:
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error",
			})
		}
	}

	return ctx.JSON(fiber.Map{
		"message": "Token refreshed successfully",
		"data": fiber.Map{
			"user":  user,
			"token": tokenPair,
		},
	})
}

// Me returns the current authenticated user
// GET /api/auth/me
func (c *AuthController) Me(ctx *fiber.Ctx) error {
	// Get token from Authorization header
	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization header is required",
		})
	}

	// Extract token from "Bearer <token>"
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid authorization header format",
		})
	}

	user, err := c.authService.GetUserFromToken(tokenString)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}

	return ctx.JSON(fiber.Map{
		"data": user,
	})
}

// condMsg returns the message if condition is true, otherwise empty string
func condMsg(cond bool, msg string) string {
	if cond {
		return msg
	}
	return ""
}
