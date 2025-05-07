package jwt

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/POABOB/slack-clone-back-end/pkg/auth"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

// 創建一個測試用的 claims 結構
type testClaims struct {
	UserID            uint              `json:"user_id"`
	Email             string            `json:"email"`
	Username          string            `json:"username"`
	RegisteredClaims  jwt.RegisteredClaims
}

func (c *testClaims) GetUserID() uint {
	return c.UserID
}

func (c *testClaims) GetEmail() string {
	return c.Email
}

func (c *testClaims) GetUsername() string {
	return c.Username
}

func (c *testClaims) GetRegisteredClaims() jwt.RegisteredClaims {
	return c.RegisteredClaims
}

func (c *testClaims) SetRegisteredClaims(claims jwt.RegisteredClaims) {
	c.RegisteredClaims = claims
}

// 創建一個模擬的 TokenManager 用於測試
type mockTokenManager struct {
	validateTokenFunc func(token string) (BaseClaims, error)
	generateTokenFunc func(claims BaseClaims) (string, error)
	refreshTokenFunc  func(tokenString string) (string, error)
	expiresIn         int
	secretKey         []byte
}

func (m *mockTokenManager) ValidateToken(token string) (BaseClaims, error) {
	return m.validateTokenFunc(token)
}

func (m *mockTokenManager) GenerateToken(claims BaseClaims) (string, error) {
	return m.generateTokenFunc(claims)
}

func (m *mockTokenManager) RefreshToken(tokenString string) (string, error) {
	return m.refreshTokenFunc(tokenString)
}

func (m *mockTokenManager) GetExpiresIn() int {
	return m.expiresIn
}

func (m *mockTokenManager) GetSecretKey() []byte {
	return m.secretKey
}

func TestNewJWTMiddleware(t *testing.T) {
	// 定義一個簡單的 ClaimsHandler 用於測試
	handlerCalled := false
	testHandler := func(c *gin.Context, claims BaseClaims) {
		handlerCalled = true
		// 驗證 claims 是否正確
		assert.Equal(t, claims.GetUserID(), uint(1))
		assert.Equal(t, claims.GetEmail(), "test@example.com")
		assert.Equal(t, claims.GetUsername(), "testuser")
	}

	t.Run("Successfully process valid token", func(t *testing.T) {
		// 重置 handlerCalled 標誌
		handlerCalled = false

		// 創建模擬的 TokenManager
		mockManager := &mockTokenManager{
			validateTokenFunc: func(token string) (BaseClaims, error) {
				return &testClaims{
					UserID:   uint(1),
					Email:    "test@example.com",
					Username: "testuser",
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
						IssuedAt:  jwt.NewNumericDate(time.Now()),
						NotBefore: jwt.NewNumericDate(time.Now()),
					},
				}, nil
			},
			generateTokenFunc: func(claims BaseClaims) (string, error) {
				return "test-token", nil
			},
			refreshTokenFunc: func(tokenString string) (string, error) {
				return "refreshed-token", nil
			},
			expiresIn: 3600,
			secretKey: []byte("test-secret"),
		}

		// 創建測試用的 gin context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Request.Header.Set("Authorization", "Bearer valid-token")

		// 測試中間件
		middleware := NewJWTMiddleware(mockManager, testHandler)
		middleware(c)

		// 驗證結果
		assert.True(t, handlerCalled)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Empty(t, c.Errors)
	})

	t.Run("No Authorization header", func(t *testing.T) {
		// 重置 handlerCalled 標誌
		handlerCalled = false

		// 創建模擬的 TokenManager
		mockManager := &mockTokenManager{
			validateTokenFunc: func(token string) (BaseClaims, error) {
				return nil, nil
			},
			generateTokenFunc: func(claims BaseClaims) (string, error) {
				return "test-token", nil
			},
			refreshTokenFunc: func(tokenString string) (string, error) {
				return "refreshed-token", nil
			},
			expiresIn: 3600,
			secretKey: []byte("test-secret"),
		}

		// 創建測試用的 gin context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)

		// 測試中間件
		middleware := NewJWTMiddleware(mockManager, testHandler)
		middleware(c)

		// 驗證結果
		assert.False(t, handlerCalled)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Len(t, c.Errors, 1)
		assert.True(t, errors.Is(c.Errors.Last().Err, auth.ErrInvalidToken))
	})

	t.Run("Invalid Authorization header format", func(t *testing.T) {
		// 重置 handlerCalled 標誌
		handlerCalled = false

		// 創建模擬的 TokenManager
		mockManager := &mockTokenManager{
			validateTokenFunc: func(token string) (BaseClaims, error) {
				return nil, nil
			},
			generateTokenFunc: func(claims BaseClaims) (string, error) {
				return "test-token", nil
			},
			refreshTokenFunc: func(tokenString string) (string, error) {
				return "refreshed-token", nil
			},
			expiresIn: 3600,
			secretKey: []byte("test-secret"),
		}

		// 創建測試用的 gin context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Request.Header.Set("Authorization", "InvalidFormat")

		// 測試中間件
		middleware := NewJWTMiddleware(mockManager, testHandler)
		middleware(c)

		// 驗證結果
		assert.False(t, handlerCalled)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Len(t, c.Errors, 1)
		assert.True(t, errors.Is(c.Errors.Last().Err, auth.ErrInvalidToken))
	})

	t.Run("Invalid token", func(t *testing.T) {
		// 重置 handlerCalled 標誌
		handlerCalled = false

		// 創建模擬的 TokenManager
		mockManager := &mockTokenManager{
			validateTokenFunc: func(token string) (BaseClaims, error) {
				return nil, auth.ErrInvalidToken
			},
			generateTokenFunc: func(claims BaseClaims) (string, error) {
				return "test-token", nil
			},
			refreshTokenFunc: func(tokenString string) (string, error) {
				return "refreshed-token", nil
			},
			expiresIn: 3600,
			secretKey: []byte("test-secret"),
		}

		// 創建測試用的 gin context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Request.Header.Set("Authorization", "Bearer invalid-token")

		// 測試中間件
		middleware := NewJWTMiddleware(mockManager, testHandler)
		middleware(c)

		// 驗證結果
		assert.False(t, handlerCalled)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Len(t, c.Errors, 1)
		assert.True(t, errors.Is(c.Errors.Last().Err, auth.ErrInvalidToken))
	})

	t.Run("Expired token", func(t *testing.T) {
		// 重置 handlerCalled 標誌
		handlerCalled = false

		// 創建模擬的 TokenManager
		mockManager := &mockTokenManager{
			validateTokenFunc: func(token string) (BaseClaims, error) {
				return nil, auth.ErrExpiredToken
			},
			generateTokenFunc: func(claims BaseClaims) (string, error) {
				return "test-token", nil
			},
			refreshTokenFunc: func(tokenString string) (string, error) {
				return "refreshed-token", nil
			},
			expiresIn: 3600,
			secretKey: []byte("test-secret"),
		}

		// 創建測試用的 gin context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Request.Header.Set("Authorization", "Bearer expired-token")

		// 測試中間件
		middleware := NewJWTMiddleware(mockManager, testHandler)
		middleware(c)

		// 驗證結果
		assert.False(t, handlerCalled)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Len(t, c.Errors, 1)
		assert.True(t, errors.Is(c.Errors.Last().Err, auth.ErrExpiredToken))
	})
}
