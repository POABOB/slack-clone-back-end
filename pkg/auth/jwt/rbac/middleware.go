package rbac

import (
	"fmt"
	"github.com/POABOB/slack-clone-back-end/pkg/auth/jwt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// TODO 重購 回傳 ERROR 與系統紀錄 ERROR

// RBACMiddleware RBAC 中間件
func RBACMiddleware(jwtManager *RBACJWTManager) gin.HandlerFunc {
	return auth.NewJWTMiddleware(jwtManager, func(c *gin.Context, claims auth.auth) {
		c.Set("user_id", claims.GetUserID())
		c.Set("email", claims.GetEmail())
		c.Set("username", claims.GetUsername())

		if rbacClaims, ok := claims.(*RBACClaims); ok {
			c.Set("role", rbacClaims.GetRole())
			c.Set("permissions", rbacClaims.GetPermissions())
		}
	})
}

// RequireRole 檢查用戶是否具有特定角色
func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, ok := extractRole(c)
		if !ok {
			return
		}

		if userRole != role {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": auth.ErrForbidden})
			_ = c.Error(auth.ErrForbidden).SetType(gin.ErrorTypePrivate)
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequirePermission 檢查用戶是否具有特定權限
func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		perms, ok := extractPermissions(c)
		if !ok {
			return
		}

		if _, ok := perms[permission]; !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": auth.ErrForbidden})
			_ = c.Error(auth.ErrForbidden).SetType(gin.ErrorTypePrivate)
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireAnyPermission 檢查用戶是否具有任意一個權限
func RequireAnyPermission(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		perms, ok := extractPermissions(c)
		if !ok {
			return
		}

		for _, required := range permissions {
			if _, ok := perms[required]; ok {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": auth.ErrForbidden})
		_ = c.Error(auth.ErrForbidden).SetType(gin.ErrorTypePrivate)
		c.Abort()
	}
}

// RequireAllPermissions 檢查用戶是否具有所有權限
func RequireAllPermissions(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		perms, ok := extractPermissions(c)
		if !ok {
			return
		}

		missingPerms := make([]string, 0)
		for _, required := range permissions {
			if _, ok := perms[required]; !ok {
				missingPerms = append(missingPerms, required)
			}
		}
		if len(missingPerms) > 0 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": auth.ErrForbidden})
			_ = c.Error(fmt.Errorf("缺少以下權限: %v", missingPerms)).SetType(gin.ErrorTypePrivate)
			c.Abort()
			return
		}
		c.Next()
	}
}

// extractRole 獲取權限資訊
func extractPermissions(c *gin.Context) (map[string]struct{}, bool) {
	raw, exists := c.Get("permissions")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": auth.ErrUnauthorized})
		_ = c.Error(auth.ErrUnauthorized).SetType(gin.ErrorTypePrivate)
		return nil, false
	}

	perms, ok := raw.(map[string]struct{})
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": auth.ErrInvalidToken})
		_ = c.Error(auth.ErrInvalidToken).SetType(gin.ErrorTypePrivate)
		return nil, false
	}
	return perms, true
}

// extractRole 獲取角色資訊
func extractRole(c *gin.Context) (string, bool) {
	raw, exists := c.Get("role")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": auth.ErrUnauthorized})
		_ = c.Error(auth.ErrUnauthorized).SetType(gin.ErrorTypePrivate)
		c.Abort()
		return "", false
	}

	role, ok := raw.(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": auth.ErrInvalidToken})
		_ = c.Error(auth.ErrInvalidToken).SetType(gin.ErrorTypePrivate)
		return "", false
	}
	return role, true
}
