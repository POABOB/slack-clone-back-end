package rbac

import (
	"errors"
	"testing"
	"time"

	"github.com/POABOB/slack-clone-back-end/pkg/auth"
	"github.com/POABOB/slack-clone-back-end/pkg/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	ExpiresIn     = 24 * 60 * 60 * 1000
	ExpiresInFast = 1
	SecretKey     = "test-secret-key"
)

// setupTestRBACJWTManager initializes a RBACJWTManager instance with test-specific configurations.
func setupTestRBACJWTManager() *RBACJWTManager {
	return NewRBACJWTManager(&config.JWTConfig{
		SecretKey: SecretKey,
		ExpiresIn: ExpiresIn,
	})
}

// setupTestRBACJWTManagerExpiresFast creates a RBACJWTManager with a test secret key and an expiry duration set to a minimal value.
func setupTestRBACJWTManagerExpiresFast() *RBACJWTManager {
	return NewRBACJWTManager(&config.JWTConfig{
		SecretKey: SecretKey,
		ExpiresIn: ExpiresInFast,
	})
}

// getDefaultRBACClaims creates and returns a default set of RBAC claims for testing purposes.
func getDefaultRBACClaims() *RBACClaims {
	userID := uint(1)
	email := "test@example.com"
	username := "testUser"
	roles := "user"
	permissions := []string{"user:read", "user:insert", "user:update", "user:delete"}

	claims := NewRBACClaims(userID, email, username, roles, permissions)
	now := time.Now()
	claims.SetRegisteredClaims(jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(ExpiresIn) * time.Millisecond)),
		NotBefore: jwt.NewNumericDate(now),
		IssuedAt:  jwt.NewNumericDate(now),
	})
	return claims
}

// TestRBACJWTManager tests all RBAC JWT manager functionality
func TestRBACJWTManager(t *testing.T) {
	t.Run("Generate and Validate Token", func(t *testing.T) {
		// Setup
		jwtManager := setupTestRBACJWTManager()
		claims := getDefaultRBACClaims()

		// 測試生成 token
		token, err := jwtManager.GenerateToken(claims)
		require.NoError(t, err)
		assert.NotEmpty(t, token)

		// 測試驗證 token
		validatedClaims, err := jwtManager.ValidateToken(token)
		require.NoError(t, err)
		assert.Equal(t, claims.UserID, validatedClaims.GetUserID())
		assert.Equal(t, claims.Email, validatedClaims.GetEmail())
		assert.Equal(t, claims.Username, validatedClaims.GetUsername())

		// 驗證 RBAC 特定的字段
		rbacClaims, ok := validatedClaims.(*RBACClaims)
		require.True(t, ok)
		assert.Equal(t, claims.Role, rbacClaims.Role)
		assert.Equal(t, claims.Permissions, rbacClaims.Permissions)
	})

	t.Run("Token Validation Scenarios", func(t *testing.T) {
		// Setup
		jwtManager := setupTestRBACJWTManager()
		claims := getDefaultRBACClaims()

		// 生成有效的 token
		validToken, err := jwtManager.GenerateToken(claims)
		require.NoError(t, err)

		// 測試驗證有效 token
		validatedClaims, err := jwtManager.ValidateToken(validToken)
		require.NoError(t, err)
		assert.Equal(t, claims.UserID, validatedClaims.GetUserID())
		assert.Equal(t, claims.Email, validatedClaims.GetEmail())
		assert.Equal(t, claims.Username, validatedClaims.GetUsername())

		// 測試使用錯誤的密鑰
		wrongCfg := &config.JWTConfig{
			SecretKey: "wrong-secret-key",
			ExpiresIn: ExpiresIn,
		}
		wrongManager := NewRBACJWTManager(wrongCfg)
		_, err = wrongManager.ValidateToken(validToken)
		assert.Error(t, err)

		// 測試過期的 token
		expiredToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJlbWFpbCI6InRlc3RAZXhhbXBsZS5jb20iLCJ1c2VybmFtZSI6InRlc3R1c2VyIiwiZXhwIjoxNTE2MjM5MDIyfQ.8tGhDxH4z3P1K6I9E2X5L8M1N4P7R0T3V6Y9"
		_, err = jwtManager.ValidateToken(expiredToken)
		assert.Error(t, err)

		// 測試無效的 token 格式
		_, err = jwtManager.ValidateToken("invalid-token")
		assert.Error(t, err)
	})

	t.Run("Refresh Token", func(t *testing.T) {
		// Setup
		jwtManager := setupTestRBACJWTManager()
		jwtManagerExpiresFast := setupTestRBACJWTManagerExpiresFast()
		claims := getDefaultRBACClaims()

		// 生成初始 token
		originalToken, err := jwtManager.GenerateToken(claims)
		require.NoError(t, err)

		// 驗證原始 token 的過期時間
		originalClaims, err := jwtManager.ValidateToken(originalToken)
		require.NoError(t, err)
		originalExpiresAt := originalClaims.GetRegisteredClaims().ExpiresAt

		// 等待一小段時間，確保時間戳不同
		time.Sleep(time.Second)

		// 測試刷新 token
		refreshedToken, err := jwtManager.RefreshToken(originalToken)
		require.NoError(t, err)
		assert.NotEmpty(t, refreshedToken)
		assert.NotEqual(t, originalToken, refreshedToken)

		// 驗證刷新後的 token
		refreshedClaims, err := jwtManager.ValidateToken(refreshedToken)
		require.NoError(t, err)
		assert.Equal(t, claims.UserID, refreshedClaims.GetUserID())
		assert.Equal(t, claims.Email, refreshedClaims.GetEmail())
		assert.Equal(t, claims.Username, refreshedClaims.GetUsername())

		// 驗證過期時間是否更新
		refreshedExpiresAt := refreshedClaims.GetRegisteredClaims().ExpiresAt
		assert.True(t, refreshedExpiresAt.After(originalExpiresAt.Time))
		assert.True(t, refreshedExpiresAt.After(time.Now()))

		// 驗證 RBAC 特定的字段是否保持不變
		rbacClaims, ok := refreshedClaims.(*RBACClaims)
		require.True(t, ok)
		assert.Equal(t, claims.Role, rbacClaims.Role)
		assert.Equal(t, claims.Permissions, rbacClaims.Permissions)

		// 測試使用無效 token 刷新
		_, err = jwtManager.RefreshToken("invalid-token")
		assert.Error(t, err)
		assert.True(t, errors.Is(err, auth.ErrInvalidToken))

		// 測試使用過期 token 刷新
		expiredToken, err := jwtManagerExpiresFast.GenerateToken(claims)
		require.NoError(t, err)
		time.Sleep(time.Second) // 等待 token 過期
		_, err = jwtManager.RefreshToken(expiredToken)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, auth.ErrExpiredToken))
	})
}
