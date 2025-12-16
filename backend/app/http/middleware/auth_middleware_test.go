package middleware

import (
	"testing"
	"time"

	"dokumen-keuangan/app/models"
	"dokumen-keuangan/app/services"
	"dokumen-keuangan/config"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// **Feature: manajemen-dokumen-keuangan, Property 24: Session Expiration Handling**
// **Validates: Requirements 10.5**
// For any expired JWT token, API requests should fail with 401 status indicating session expiration.
func TestProperty24_SessionExpirationHandling(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	cfg := &config.Config{
		AppName:       "TestApp",
		JWTSecret:     "test-secret-key-for-session-test",
		JWTTTL:        60,
		JWTRefreshTTL: 20160,
	}

	authService := services.NewAuthService(cfg)

	properties.Property("expired tokens should fail validation", prop.ForAll(
		func(username string, hoursExpired int) bool {
			// Create an expired token
			user := &models.User{
				ID:       uuid.New(),
				Username: username,
				Role:     models.RoleOperator,
				IsActive: true,
			}

			// Generate expired token manually
			expiredToken := createExpiredToken(user, cfg, time.Duration(hoursExpired)*time.Hour)

			// Validate should fail
			_, err := authService.ValidateToken(expiredToken)
			return err != nil
		},
		gen.AlphaString().SuchThat(func(s string) bool {
			return len(s) >= 1 && len(s) <= 50
		}),
		gen.IntRange(1, 100), // Hours expired (1 to 100 hours ago)
	))

	properties.Property("valid non-expired tokens should pass validation", prop.ForAll(
		func(username string, role string) bool {
			user := &models.User{
				ID:       uuid.New(),
				Username: username,
				Role:     models.UserRole(role),
				IsActive: true,
			}

			// Generate valid token
			tokenPair, err := authService.GenerateTokenPair(user)
			if err != nil {
				return false
			}

			// Validate should succeed
			claims, err := authService.ValidateToken(tokenPair.AccessToken)
			if err != nil {
				return false
			}

			return claims.UserID == user.ID && claims.Type == "access"
		},
		gen.AlphaString().SuchThat(func(s string) bool {
			return len(s) >= 1 && len(s) <= 50
		}),
		gen.OneConstOf("super_admin", "admin", "operator"),
	))

	properties.Property("token expiration time should be respected", prop.ForAll(
		func(username string) bool {
			user := &models.User{
				ID:       uuid.New(),
				Username: username,
				Role:     models.RoleOperator,
				IsActive: true,
			}

			// Generate token
			tokenPair, err := authService.GenerateTokenPair(user)
			if err != nil {
				return false
			}

			// Validate and check expiration
			claims, err := authService.ValidateToken(tokenPair.AccessToken)
			if err != nil {
				return false
			}

			// Token should expire in the future
			return claims.ExpiresAt.Time.After(time.Now())
		},
		gen.AlphaString().SuchThat(func(s string) bool {
			return len(s) >= 1 && len(s) <= 50
		}),
	))

	properties.TestingRun(t)
}

// createExpiredToken creates a JWT token that expired hoursAgo hours ago
func createExpiredToken(user *models.User, cfg *config.Config, hoursAgo time.Duration) string {
	now := time.Now()
	claims := services.JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		Type:     "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(-hoursAgo)),
			IssuedAt:  jwt.NewNumericDate(now.Add(-hoursAgo - time.Hour)),
			NotBefore: jwt.NewNumericDate(now.Add(-hoursAgo - time.Hour)),
			Issuer:    cfg.AppName,
			Subject:   user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(cfg.JWTSecret))
	return tokenString
}
