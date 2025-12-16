package controllers

import (
	"strconv"

	"dokumen-keuangan/app/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// UserController handles user endpoints
type UserController struct {
	service *services.UserService
}

// NewUserController creates a new UserController instance
func NewUserController() *UserController {
	return &UserController{
		service: services.NewUserService(),
	}
}

// GetAll retrieves all users with pagination
// GET /api/users
func (c *UserController) GetAll(ctx *fiber.Ctx) error {
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.Query("page_size", "10"))
	search := ctx.Query("search", "")
	role := ctx.Query("role", "")

	result, err := c.service.GetAllWithFilter(page, pageSize, search, role)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve users",
		})
	}

	return ctx.JSON(fiber.Map{
		"data":        result.Data,
		"total":       result.Total,
		"page":        result.Page,
		"page_size":   result.PageSize,
		"total_pages": result.TotalPages,
	})
}

// GetByID retrieves a user by ID
// GET /api/users/:id
func (c *UserController) GetByID(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	user, err := c.service.GetByID(id)
	if err != nil {
		if err == services.ErrUserNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve user",
		})
	}

	return ctx.JSON(fiber.Map{
		"data": user,
	})
}


// Create creates a new user
// POST /api/users
func (c *UserController) Create(ctx *fiber.Ctx) error {
	var input services.CreateUserInput
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate input
	if validationErrors := input.Validate(); len(validationErrors) > 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": validationErrors,
		})
	}

	user, err := c.service.Create(&input)
	if err != nil {
		switch err {
		case services.ErrUserUsernameExists:
			return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Username already exists",
			})
		case services.ErrUserUsernameRequired:
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"errors": fiber.Map{"username": "Username is required"},
			})
		case services.ErrUserPasswordRequired:
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"errors": fiber.Map{"password": "Password is required"},
			})
		case services.ErrUserNameRequired:
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"errors": fiber.Map{"name": "Name is required"},
			})
		case services.ErrUserRoleRequired:
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"errors": fiber.Map{"role": "Role is required"},
			})
		case services.ErrUserInvalidRole:
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"errors": fiber.Map{"role": "Invalid role"},
			})
		default:
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
		"data":    user,
	})
}

// Update updates a user
// PUT /api/users/:id
func (c *UserController) Update(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	var input services.UpdateUserInput
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	user, err := c.service.Update(id, &input)
	if err != nil {
		switch err {
		case services.ErrUserNotFound:
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		case services.ErrUserUsernameExists:
			return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Username already exists",
			})
		case services.ErrUserInvalidRole:
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"errors": fiber.Map{"role": "Invalid role"},
			})
		default:
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	return ctx.JSON(fiber.Map{
		"message": "User updated successfully",
		"data":    user,
	})
}

// Delete deactivates a user
// DELETE /api/users/:id
func (c *UserController) Delete(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	err = c.service.Delete(id)
	if err != nil {
		if err == services.ErrUserNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to deactivate user",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "User deactivated successfully",
	})
}

// Activate activates a user
// POST /api/users/:id/activate
func (c *UserController) Activate(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	user, err := c.service.Activate(id)
	if err != nil {
		if err == services.ErrUserNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to activate user",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "User activated successfully",
		"data":    user,
	})
}


// UploadAvatar uploads avatar for a user
// POST /api/users/:id/avatar
func (c *UserController) UploadAvatar(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Get user
	user, err := c.service.GetByID(id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Get uploaded file
	file, err := ctx.FormFile("avatar")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Avatar file is required",
		})
	}

	// Upload file
	fileService := services.NewFileService()
	avatarPath, err := fileService.UploadAvatar(file, "users", id)
	if err != nil {
		switch err {
		case services.ErrInvalidFileType:
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid file type. Allowed: JPG, PNG, GIF, WebP",
			})
		case services.ErrFileTooLarge:
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "File too large. Maximum size: 5MB",
			})
		default:
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to upload avatar",
			})
		}
	}

	// Update user avatar path
	updatedUser, err := c.service.Update(id, &services.UpdateUserInput{
		AvatarPath: &avatarPath,
	})
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user avatar",
		})
	}

	_ = user // Suppress unused variable warning

	return ctx.JSON(fiber.Map{
		"message": "Avatar uploaded successfully",
		"data":    updatedUser,
	})
}
