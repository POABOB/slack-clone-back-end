package simple

import (
	"errors"
	"github.com/POABOB/slack-clone-back-end/pkg/auth"
	"testing"
	"time"

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

// setupTestJWTManager initializes a JWTManager instance with test-specific configurations.
func setupTestJWTManager() *JWTManager {
	return NewJWTManager(&config.JWTConfig{
		SecretKey: []byte(SecretKey),
		ExpiresIn: ExpiresIn,
	})
}

// setupTestJWTManagerExpiresFast creates a JWTManager with a test secret key and an expiry duration set to a minimal value.
func setupTestJWTManagerExpiresFast() *JWTManager {
	return NewJWTManager(&config.JWTConfig{
		SecretKey: []byte(SecretKey),
		ExpiresIn: ExpiresInFast,
	})
}

// getDefaultClaims creates and returns a default set of claims for testing purposes, including user ID, email, and username.
func getDefaultClaims() *DefaultClaims {
	userID := uint(1)
	email := "test@example.com"
	username := "testUser"

	claims := NewDefaultClaims(userID, email, username)
	now := time.Now()
	claims.SetRegisteredClaims(jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(ExpiresIn) * time.Millisecond)),
		NotBefore: jwt.NewNumericDate(now),
		IssuedAt:  jwt.NewNumericDate(now),
	})
	return claims
}

// TestJWTManager tests all JWT manager functionality
func TestJWTManager(t *testing.T) {
	t.Run("Generate and Validate Token", func(t *testing.T) {
		// Setup
		jwtManager := setupTestJWTManager()
		claims := getDefaultClaims()

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
		assert.True(t, validatedClaims.GetRegisteredClaims().ExpiresAt.After(time.Now()))
	})

	t.Run("Token Validation Scenarios", func(t *testing.T) {
		// Setup
		jwtManager := setupTestJWTManager()
		claims := getDefaultClaims()

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
			SecretKey: []byte("wrong-secret-key"),
			ExpiresIn: ExpiresIn,
		}
		wrongManager := NewJWTManager(wrongCfg)
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

	t.Run("Claims Functionality", func(t *testing.T) {
		// Setup
		claims := getDefaultClaims()

		// 測試有效的 claims
		assert.True(t, claims.GetRegisteredClaims().ExpiresAt.After(time.Now()))

		// 測試過期的 claims
		expiredClaims := getDefaultClaims()
		expiredClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(-24 * time.Hour))
		assert.True(t, expiredClaims.ExpiresAt.Before(time.Now()))
	})

	t.Run("Token Expiration", func(t *testing.T) {
		// Setup
		jwtManager := setupTestJWTManager()
		claims := getDefaultClaims()

		// 生成 token
		token, err := jwtManager.GenerateToken(claims)
		require.NoError(t, err)

		// 驗證 token
		validatedClaims, err := jwtManager.ValidateToken(token)
		require.NoError(t, err)

		// 檢查過期時間是否在合理範圍內
		expectedExpiration := time.Now().Add(24 * time.Hour)
		assert.True(t, validatedClaims.GetRegisteredClaims().ExpiresAt.After(time.Now()))
		assert.True(t, validatedClaims.GetRegisteredClaims().ExpiresAt.Before(expectedExpiration.Add(time.Minute)))
		assert.True(t, validatedClaims.GetRegisteredClaims().ExpiresAt.After(expectedExpiration.Add(-time.Minute)))
	})

	t.Run("Refresh Token", func(t *testing.T) {
		// Setup
		jwtManager := setupTestJWTManager()
		jwtManagerExpiresFast := setupTestJWTManagerExpiresFast()
		claims := getDefaultClaims()

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

	t.Run("JWT Manager Configuration", func(t *testing.T) {
		// 測試默認過期時間
		cfg := &config.JWTConfig{
			SecretKey: []byte(SecretKey),
		}
		manager := NewJWTManager(cfg)
		assert.Equal(t, 24*60*60*1000, manager.GetExpiresIn())

		// 測試自定義過期時間
		customExpiresIn := 12 * 60 * 60 * 1000 // 12 hours
		cfg = &config.JWTConfig{
			SecretKey: []byte(SecretKey),
			ExpiresIn: customExpiresIn,
		}
		manager = NewJWTManager(cfg)
		assert.Equal(t, customExpiresIn, manager.GetExpiresIn())
	})
}
