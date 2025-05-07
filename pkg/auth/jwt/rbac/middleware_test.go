package rbac

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/POABOB/slack-clone-back-end/pkg/auth"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRouter() (*gin.Engine, *RBACJWTManager) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	jwtManager := setupTestRBACJWTManager()
	return router, jwtManager
}

func TestRBACMiddleware(t *testing.T) {
	// Setup
	jwtManager := setupTestRBACJWTManager()
	claims := getDefaultRBACClaims()
	token, err := jwtManager.GenerateToken(claims)
	require.NoError(t, err)

	t.Run("Successfully set user info", func(t *testing.T) {
		// 創建測試用的 gin context
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Request.Header.Set("Authorization", "Bearer "+token)

		// 測試中間件
		handler := RBACMiddleware(jwtManager)
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

		role, exists := c.Get("role")
		assert.True(t, exists)
		assert.Equal(t, claims.GetRole(), role)

		permissions, exists := c.Get("permissions")
		assert.True(t, exists)
		assert.ElementsMatch(t, claims.GetPermissions(), permissions)
	})

	t.Run("No token provided", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)

		handler := RBACMiddleware(jwtManager)
		handler(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Len(t, c.Errors, 1)
		assert.True(t, errors.Is(c.Errors.Last().Err, auth.ErrInvalidToken))
	})

	t.Run("Invalid token", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Request.Header.Set("Authorization", "Bearer invalid-token")

		handler := RBACMiddleware(jwtManager)
		handler(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Len(t, c.Errors, 1)
		assert.True(t, errors.Is(c.Errors.Last().Err, auth.ErrInvalidToken))
	})
}

func TestRequireRole(t *testing.T) {
	// Setup
	jwtManager := setupTestRBACJWTManager()
	claims := getDefaultRBACClaims()
	token, err := jwtManager.GenerateToken(claims)
	require.NoError(t, err)

	t.Run("Role matches requirement", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Request.Header.Set("Authorization", "Bearer "+token)

		// 先設置 RBAC 中間件
		rbacHandler := RBACMiddleware(jwtManager)
		rbacHandler(c)

		// 測試角色符合要求
		roleHandler := RequireRole("user")
		roleHandler(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Role does not match requirement", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Request.Header.Set("Authorization", "Bearer "+token)

		// 先設置 RBAC 中間件
		rbacHandler := RBACMiddleware(jwtManager)
		rbacHandler(c)

		// 測試角色不符合要求
		roleHandler := RequireRole("admin")
		roleHandler(c)
		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Len(t, c.Errors, 1)
		assert.True(t, errors.Is(c.Errors.Last().Err, auth.ErrForbidden))
	})
}

func TestRequirePermission(t *testing.T) {
	// Setup
	jwtManager := setupTestRBACJWTManager()
	claims := getDefaultRBACClaims()
	token, err := jwtManager.GenerateToken(claims)
	require.NoError(t, err)

	t.Run("Has required permission", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Request.Header.Set("Authorization", "Bearer "+token)

		// 先設置 RBAC 中間件
		rbacHandler := RBACMiddleware(jwtManager)
		rbacHandler(c)

		// 測試具有所需權限
		permHandler := RequirePermission("user:read")
		permHandler(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Does not have required permission", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Request.Header.Set("Authorization", "Bearer "+token)

		// 先設置 RBAC 中間件
		rbacHandler := RBACMiddleware(jwtManager)
		rbacHandler(c)

		// 測試不具有所需權限
		permHandler := RequirePermission("admin:read")
		permHandler(c)
		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Len(t, c.Errors, 1)
		assert.True(t, errors.Is(c.Errors.Last().Err, auth.ErrForbidden))
	})
}

func TestRequireAnyPermission(t *testing.T) {
	// Setup
	jwtManager := setupTestRBACJWTManager()
	claims := getDefaultRBACClaims()
	token, err := jwtManager.GenerateToken(claims)
	require.NoError(t, err)

	t.Run("Has any of required permissions", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Request.Header.Set("Authorization", "Bearer "+token)

		// 先設置 RBAC 中間件
		rbacHandler := RBACMiddleware(jwtManager)
		rbacHandler(c)

		// 測試具有任一所需權限
		permHandler := RequireAnyPermission("user:read", "admin:read")
		permHandler(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Does not have any of required permissions", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Request.Header.Set("Authorization", "Bearer "+token)

		// 先設置 RBAC 中間件
		rbacHandler := RBACMiddleware(jwtManager)
		rbacHandler(c)

		// 測試不具有任何所需權限
		permHandler := RequireAnyPermission("admin:read", "admin:write")
		permHandler(c)
		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Len(t, c.Errors, 1)
		assert.True(t, errors.Is(c.Errors.Last().Err, auth.ErrForbidden))
	})
}

func TestRequireAllPermissions(t *testing.T) {
	// Setup
	jwtManager := setupTestRBACJWTManager()
	claims := getDefaultRBACClaims()
	token, err := jwtManager.GenerateToken(claims)
	require.NoError(t, err)

	t.Run("Has all required permissions", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Request.Header.Set("Authorization", "Bearer "+token)

		// 先設置 RBAC 中間件
		rbacHandler := RBACMiddleware(jwtManager)
		rbacHandler(c)

		// 測試具有所有所需權限
		permHandler := RequireAllPermissions("user:read", "user:insert")
		permHandler(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Missing some required permissions", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Request.Header.Set("Authorization", "Bearer "+token)

		// 先設置 RBAC 中間件
		rbacHandler := RBACMiddleware(jwtManager)
		rbacHandler(c)

		// 測試缺少部分所需權限
		permHandler := RequireAllPermissions("user:read", "admin:read")
		permHandler(c)
		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Len(t, c.Errors, 1)
		assert.True(t, errors.Is(c.Errors.Last().Err, auth.ErrForbidden))
	})
}

func TestExtractPermissions(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Permissions not set in context", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)

		perms, ok := extractPermissions(c)
		assert.False(t, ok)
		assert.Nil(t, perms)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Len(t, c.Errors, 1)
		assert.True(t, errors.Is(c.Errors.Last().Err, auth.ErrUnauthorized))
	})

	t.Run("Invalid permissions type in context", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Set("permissions", "invalid-type") // 設置錯誤類型的權限

		perms, ok := extractPermissions(c)
		assert.False(t, ok)
		assert.Nil(t, perms)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Len(t, c.Errors, 1)
		assert.True(t, errors.Is(c.Errors.Last().Err, auth.ErrInvalidToken))
	})
}

func TestExtractRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Role not set in context", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)

		role, ok := extractRole(c)
		assert.False(t, ok)
		assert.Empty(t, role)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Len(t, c.Errors, 1)
		assert.True(t, errors.Is(c.Errors.Last().Err, auth.ErrUnauthorized))
	})

	t.Run("Invalid role type in context", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Set("role", 123) // 設置錯誤類型的角色

		role, ok := extractRole(c)
		assert.False(t, ok)
		assert.Empty(t, role)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Len(t, c.Errors, 1)
		assert.True(t, errors.Is(c.Errors.Last().Err, auth.ErrInvalidToken))
	})
}

func TestRequirePermissionErrorCases(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Permissions not set in context", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)

		handler := RequirePermission("user:read")
		handler(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Len(t, c.Errors, 1)
		assert.True(t, errors.Is(c.Errors.Last().Err, auth.ErrUnauthorized))
	})

	t.Run("Invalid permissions type in context", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Set("permissions", "invalid-type")

		handler := RequirePermission("user:read")
		handler(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Len(t, c.Errors, 1)
		assert.True(t, errors.Is(c.Errors.Last().Err, auth.ErrInvalidToken))
	})
}

func TestRequireRoleErrorCases(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Role not set in context", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)

		handler := RequireRole("user")
		handler(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Len(t, c.Errors, 1)
		assert.True(t, errors.Is(c.Errors.Last().Err, auth.ErrUnauthorized))
	})

	t.Run("Invalid role type in context", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Set("role", 123)

		handler := RequireRole("user")
		handler(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Len(t, c.Errors, 1)
		assert.True(t, errors.Is(c.Errors.Last().Err, auth.ErrInvalidToken))
	})
}
