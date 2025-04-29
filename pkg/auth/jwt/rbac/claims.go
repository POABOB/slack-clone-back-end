package rbac

import (
	"github.com/POABOB/slack-clone-back-end/pkg/auth/jwt/simple"
)

// RBACClaims RBAC JWT 聲明
type RBACClaims struct {
	*simple.DefaultClaims
	Role        string              `json:"role"`
	Permissions map[string]struct{} `json:"permissions"`
}

// NewRBACClaims 創建新的 RBAC 聲明
func NewRBACClaims(userID uint, email, username, role string, permissions []string) *RBACClaims {
	// 將權限列表轉換為 map
	permMap := make(map[string]struct{})
	for _, p := range permissions {
		permMap[p] = struct{}{}
	}

	return &RBACClaims{
		DefaultClaims: &simple.DefaultClaims{
			UserID:   userID,
			Email:    email,
			Username: username,
		},
		Role:        role,
		Permissions: permMap,
	}
}

// GetRole 獲取用戶角色
func (c *RBACClaims) GetRole() string {
	return c.Role
}

// GetPermissions 獲取用戶權限列表
func (c *RBACClaims) GetPermissions() []string {
	permissions := make([]string, 0, len(c.Permissions))
	for p := range c.Permissions {
		permissions = append(permissions, p)
	}
	return permissions
}

// HasPermission 檢查用戶是否具有特定權限
func (c *RBACClaims) HasPermission(permission string) bool {
	_, exists := c.Permissions[permission]
	return exists
}

// HasAnyPermission 檢查用戶是否具有任意一個權限
func (c *RBACClaims) HasAnyPermission(permissions ...string) bool {
	for _, p := range permissions {
		if _, exists := c.Permissions[p]; exists {
			return true
		}
	}
	return false
}

// HasAllPermissions 檢查用戶是否具有所有權限
func (c *RBACClaims) HasAllPermissions(permissions ...string) bool {
	for _, p := range permissions {
		if _, exists := c.Permissions[p]; !exists {
			return false
		}
	}
	return true
}
