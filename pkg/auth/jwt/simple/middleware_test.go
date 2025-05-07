package simple

import (
	"errors"
	"github.com/POABOB/slack-clone-back-end/pkg/auth"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestJWTAuthMiddleware tests the JWT authentication middleware.
func TestJWTAuthMiddleware(t *testing.T) {
	// Setup
	jwtManager := setupTestJWTManager()
	claims := getDefaultClaims()
	token, err := jwtManager.GenerateToken(claims)
	require.NoError(t, err)

	// 創建測試用的 gin context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	// 測試中間件
	handler := JWTAuthMiddleware(jwtManager)
	handler(c)

	// 驗證 context 中的值
	userID, exists := c.Get("user_id")
	assert.True(t, exists)
	assert.Equal(t, claims.GetUserID(), userID)

	email, exists := c.Get("email")
	assert.True(t, exists)
	assert.Equal(t, claims.GetEmail(), email)

	username, exists := c.Get("username")
	assert.True(t, exists)
	assert.Equal(t, claims.GetUsername(), username)

	// 測試無 token 的情況
	c.Request.Header.Del("Authorization")
	handler(c)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// 測試無效 token 的情況
	c.Request.Header.Set("Authorization", "Bearer invalid-token")
	handler(c)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestJWTAuthMiddlewareErrorCases tests various error cases in the JWT middleware.
func TestJWTAuthMiddlewareErrorCases(t *testing.T) {
	// Setup
	jwtManager := setupTestJWTManager()
	gin.SetMode(gin.TestMode)

	// 測試無效的 Bearer token 格式
	t.Run("Invalid Bearer Format", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Request.Header.Set("Authorization", "InvalidFormat token123")

		handler := JWTAuthMiddleware(jwtManager)
		handler(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Len(t, c.Errors, 1)
		assert.True(t, errors.Is(c.Errors.Last().Err, auth.ErrInvalidToken))
	})

	// 測試過期 token
	t.Run("Expired Token", func(t *testing.T) {
		// 使用快速過期的 JWTManager 生成 token
		fastManager := setupTestJWTManagerExpiresFast()
		claims := getDefaultClaims()
		token, err := fastManager.GenerateToken(claims)
		require.NoError(t, err)

		// 等待 token 過期
		time.Sleep(time.Second)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Request.Header.Set("Authorization", "Bearer "+token)

		handler := JWTAuthMiddleware(jwtManager)
		handler(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Len(t, c.Errors, 1)
		assert.True(t, errors.Is(c.Errors.Last().Err, auth.ErrExpiredToken))
	})

	// 測試無效的 token 格式
	t.Run("Invalid Token Format", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Request.Header.Set("Authorization", "Bearer invalid.token.format")

		handler := JWTAuthMiddleware(jwtManager)
		handler(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Len(t, c.Errors, 1)
		assert.True(t, errors.Is(c.Errors.Last().Err, auth.ErrInvalidToken))
	})

	// 測試缺少 Authorization header
	t.Run("Missing Authorization Header", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)

		handler := JWTAuthMiddleware(jwtManager)
		handler(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Len(t, c.Errors, 1)
		assert.True(t, errors.Is(c.Errors.Last().Err, auth.ErrInvalidToken))
	})
}
