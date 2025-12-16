package services

import (
	"testing"
	"time"

	"dokumen-keuangan/app/models"
	"dokumen-keuangan/config"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// **Feature: manajemen-dokumen-keuangan, Property 21: Authentication Success**
// **Validates: Requirements 10.1**
// For any user with valid credentials, login should succeed and return a valid JWT token.
func TestProperty21_AuthenticationSuccess(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	cfg := &config.Config{
		AppName:       "TestApp",
		JWTSecret:     "test-secret-key-for-testing",
		JWTTTL:        60,
		JWTRefreshTTL: 20160,
	}

	properties.Property("valid password should verify successfully", prop.ForAll(
		func(password string) bool {
			// Hash the password
			hashedPassword, err := HashPassword(password)
			if err != nil {
				return false
			}

			// Verify the password
			return VerifyPassword(password, hashedPassword)
		},
		gen.AlphaString().SuchThat(func(s string) bool {
			return len(s) >= 1 && len(s) <= 72 // bcrypt max length
		}),
	))

	properties.Property("JWT token generation should produce valid tokens", prop.ForAll(
		func(username string, role string) bool {
			authService := &AuthService{
				config:    cfg,
				jwtSecret: []byte(cfg.JWTSecret),
			}

			user := &models.User{
				ID:       uuid.New(),
				Username: username,
				Role:     models.UserRole(role),
				IsActive: true,
			}

			// Generate token pair
			tokenPair, err := authService.GenerateTokenPair(user)
			if err != nil {
				return false
			}

			// Validate access token
			accessClaims, err := authService.ValidateToken(tokenPair.AccessToken)
			if err != nil {
				return false
			}

			// Validate refresh token
			refreshClaims, err := authService.ValidateToken(tokenPair.RefreshToken)
			if err != nil {
				return false
			}

			// Check claims match user data
			return accessClaims.UserID == user.ID &&
				accessClaims.Username == user.Username &&
				accessClaims.Role == user.Role &&
				accessClaims.Type == "access" &&
				refreshClaims.Type == "refresh" &&
				tokenPair.TokenType == "Bearer"
		},
		gen.AlphaString().SuchThat(func(s string) bool {
			return len(s) >= 1 && len(s) <= 50
		}),
		gen.OneConstOf("super_admin", "admin", "operator"),
	))

	properties.TestingRun(t)
}

// **Feature: manajemen-dokumen-keuangan, Property 22: Authentication Failure**
// **Validates: Requirements 10.2**
// For any invalid credentials (wrong password), login should fail with appropriate error.
func TestProperty22_AuthenticationFailure(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("wrong password should not verify", prop.ForAll(
		func(password, wrongPassword string) bool {
			// Skip if passwords are the same
			if password == wrongPassword {
				return true
			}

			// Hash the correct password
			hashedPassword, err := HashPassword(password)
			if err != nil {
				return false
			}

			// Verify with wrong password should fail
			return !VerifyPassword(wrongPassword, hashedPassword)
		},
		gen.AlphaString().SuchThat(func(s string) bool {
			return len(s) >= 1 && len(s) <= 72
		}),
		gen.AlphaString().SuchThat(func(s string) bool {
			return len(s) >= 1 && len(s) <= 72
		}),
	))

	properties.Property("invalid token should fail validation", prop.ForAll(
		func(invalidToken string) bool {
			cfg := &config.Config{
				AppName:       "TestApp",
				JWTSecret:     "test-secret-key-for-testing",
				JWTTTL:        60,
				JWTRefreshTTL: 20160,
			}

			authService := &AuthService{
				config:    cfg,
				jwtSecret: []byte(cfg.JWTSecret),
			}

			// Invalid token should fail validation
			_, err := authService.ValidateToken(invalidToken)
			return err != nil
		},
		gen.AlphaString(),
	))

	properties.Property("token with wrong secret should fail validation", prop.ForAll(
		func(username string) bool {
			cfg1 := &config.Config{
				AppName:       "TestApp",
				JWTSecret:     "secret-key-1",
				JWTTTL:        60,
				JWTRefreshTTL: 20160,
			}
			cfg2 := &config.Config{
				AppName:       "TestApp",
				JWTSecret:     "secret-key-2",
				JWTTTL:        60,
				JWTRefreshTTL: 20160,
			}

			authService1 := &AuthService{
				config:    cfg1,
				jwtSecret: []byte(cfg1.JWTSecret),
			}
			authService2 := &AuthService{
				config:    cfg2,
				jwtSecret: []byte(cfg2.JWTSecret),
			}

			user := &models.User{
				ID:       uuid.New(),
				Username: username,
				Role:     models.RoleOperator,
				IsActive: true,
			}

			// Generate token with service 1
			tokenPair, err := authService1.GenerateTokenPair(user)
			if err != nil {
				return false
			}

			// Validate with service 2 (different secret) should fail
			_, err = authService2.ValidateToken(tokenPair.AccessToken)
			return err != nil
		},
		gen.AlphaString().SuchThat(func(s string) bool {
			return len(s) >= 1 && len(s) <= 50
		}),
	))

	properties.TestingRun(t)
}

// TestExpiredTokenValidation tests that expired tokens are rejected
func TestExpiredTokenValidation(t *testing.T) {
	cfg := &config.Config{
		AppName:       "TestApp",
		JWTSecret:     "test-secret-key",
		JWTTTL:        60,
		JWTRefreshTTL: 20160,
	}

	authService := &AuthService{
		config:    cfg,
		jwtSecret: []byte(cfg.JWTSecret),
	}

	user := &models.User{
		ID:       uuid.New(),
		Username: "testuser",
		Role:     models.RoleOperator,
		IsActive: true,
	}

	// Create an expired token manually
	now := time.Now()
	claims := JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		Type:     "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(-1 * time.Hour)), // Expired 1 hour ago
			IssuedAt:  jwt.NewNumericDate(now.Add(-2 * time.Hour)),
			NotBefore: jwt.NewNumericDate(now.Add(-2 * time.Hour)),
			Issuer:    cfg.AppName,
			Subject:   user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	expiredToken, err := token.SignedString(authService.jwtSecret)
	if err != nil {
		t.Fatalf("Failed to create expired token: %v", err)
	}

	// Validate expired token should fail
	_, err = authService.ValidateToken(expiredToken)
	if err == nil {
		t.Error("Expected error for expired token, got nil")
	}
}
