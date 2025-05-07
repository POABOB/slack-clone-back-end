package rbac

import (
	"github.com/POABOB/slack-clone-back-end/pkg/auth/jwt/simple"
)

// RBACClaims RBAC JWT 聲明
type RBACClaims struct {
	*simple.DefaultClaims
	Role        string              `json:"role"`
	// Permissions 使用 map 作為底層資料結構，key 為權限字串，value 為空結構體
	// 使用 map 而非 slice 的原因：
	// 1. 快速查詢：O(1) 時間複雜度檢查權限是否存在
	// 2. 自動去重：相同的權限只會存在一次
	// 3. 記憶體效率：使用空結構體作為 value 不佔用額外空間
	Permissions map[string]struct{} `json:"permissions"`
}

// NewRBACClaims 創建新的 RBAC 聲明
func NewRBACClaims(userID uint, email, username, role string, permissions []string) *RBACClaims {
	// 將權限列表轉換為 map
	// 使用 map 作為內部存儲結構，提供 O(1) 的權限查詢效率
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
// 將內部存儲的 map 轉換回 []string 形式返回
// 這樣做的好處：
// 1. 對外提供更直觀的切片介面
// 2. 保持與輸入參數類型一致
// 3. 方便 JSON 序列化和反序列化
func (c *RBACClaims) GetPermissions() []string {
	permissions := make([]string, 0, len(c.Permissions))
	for p := range c.Permissions {
		permissions = append(permissions, p)
	}
	return permissions
}
