package services

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"dokumen-keuangan/app/models"
	"dokumen-keuangan/app/repositories"
	"dokumen-keuangan/config"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrUserInactive       = errors.New("user account is inactive")
	ErrInvalidToken       = errors.New("invalid or expired token")
)

// TokenPair represents access and refresh tokens
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// JWTClaims represents the claims in a JWT token
type JWTClaims struct {
	UserID   uuid.UUID       `json:"user_id"`
	Username string          `json:"username"`
	Role     models.UserRole `json:"role"`
	Type     string          `json:"type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

// AuthService handles authentication logic
type AuthService struct {
	userRepo  *repositories.UserRepository
	config    *config.Config
	jwtSecret []byte
}

// NewAuthService creates a new AuthService instance
func NewAuthService(cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo:  repositories.NewUserRepository(),
		config:    cfg,
		jwtSecret: []byte(cfg.JWTSecret),
	}
}


// Login authenticates a user and returns token pair
func (s *AuthService) Login(username, password string) (*TokenPair, *models.User, error) {
	// Find user by username
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return nil, nil, ErrInvalidCredentials
	}

	// Check if user is active
	if !user.IsActive {
		return nil, nil, ErrUserInactive
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, nil, ErrInvalidCredentials
	}

	// Generate tokens
	tokenPair, err := s.GenerateTokenPair(user)
	if err != nil {
		return nil, nil, err
	}

	return tokenPair, user, nil
}

// GenerateTokenPair generates access and refresh tokens for a user
func (s *AuthService) GenerateTokenPair(user *models.User) (*TokenPair, error) {
	// Generate access token
	accessToken, err := s.generateToken(user, "access", time.Duration(s.config.JWTTTL)*time.Minute)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, err := s.generateToken(user, "refresh", time.Duration(s.config.JWTRefreshTTL)*time.Minute)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.config.JWTTTL * 60), // in seconds
		TokenType:    "Bearer",
	}, nil
}

// generateToken creates a JWT token
func (s *AuthService) generateToken(user *models.User, tokenType string, expiry time.Duration) (string, error) {
	now := time.Now()
	claims := JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		Type:     tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    s.config.AppName,
			Subject:   user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// ValidateToken validates a JWT token and returns the claims
func (s *AuthService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

// RefreshTokens generates new token pair from a valid refresh token
func (s *AuthService) RefreshTokens(refreshToken string) (*TokenPair, *models.User, error) {
	// Validate refresh token
	claims, err := s.ValidateToken(refreshToken)
	if err != nil {
		return nil, nil, err
	}

	// Check if it's a refresh token
	if claims.Type != "refresh" {
		return nil, nil, ErrInvalidToken
	}

	// Get user from database
	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return nil, nil, ErrInvalidToken
	}

	// Check if user is still active
	if !user.IsActive {
		return nil, nil, ErrUserInactive
	}

	// Generate new token pair
	tokenPair, err := s.GenerateTokenPair(user)
	if err != nil {
		return nil, nil, err
	}

	return tokenPair, user, nil
}

// GetUserFromToken retrieves user from a valid access token
func (s *AuthService) GetUserFromToken(tokenString string) (*models.User, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	// Check if it's an access token
	if claims.Type != "access" {
		return nil, ErrInvalidToken
	}

	return s.userRepo.FindByID(claims.UserID)
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// VerifyPassword verifies a password against a hash
func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// UpdateProfile updates the user's profile
func (s *AuthService) UpdateProfile(userID uuid.UUID, name, username, password string) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	user.Name = name
	user.Username = username
	// user.Email = email // Removed as User model likely doesn't support it or user rejected it

	if password != "" {
		hashed, err := HashPassword(password)
		if err != nil {
			return nil, err
		}
		user.Password = hashed
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateAvatar updates the user's avatar
func (s *AuthService) UpdateAvatar(userID uuid.UUID, file *multipart.FileHeader, ctx *fiber.Ctx) (string, error) {
	// Validate file
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
	if !allowedExts[ext] {
		return "", errors.New("invalid file type")
	}

	if file.Size > 2*1024*1024 { // 2MB
		return "", errors.New("file too large (max 2MB)")
	}

	// Create directory
	dir := "./storage/app/public/avatars"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	// Save file
	filename := fmt.Sprintf("%s_%d%s", userID.String(), time.Now().Unix(), ext)
	path := filepath.Join(dir, filename)
	if err := ctx.SaveFile(file, path); err != nil {
		return "", err
	}

	// Update user
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return "", err
	}

	publicPath := fmt.Sprintf("/api/files/avatars/%s", filename)
	user.AvatarPath = &publicPath

	if err := s.userRepo.Update(user); err != nil {
		return "", err
	}

	return publicPath, nil
}
