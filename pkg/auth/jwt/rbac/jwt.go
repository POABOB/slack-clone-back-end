package rbac

import (
	"github.com/POABOB/slack-clone-back-end/pkg/auth/jwt"
	"github.com/POABOB/slack-clone-back-end/pkg/auth/jwt/simple"
	"github.com/POABOB/slack-clone-back-end/pkg/jwt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// RBACJWTManager RBAC JWT 管理器
type RBACJWTManager struct {
	*simple.JWTManager
}

// NewRBACJWTManager 創建新的 RBAC JWT 管理器
func NewRBACJWTManager(secretKey string, expiresIn int) *RBACJWTManager {
	return &RBACJWTManager{
		JWTManager: simple.NewJWTManager(secretKey, expiresIn),
	}
}

// GenerateToken 生成 RBAC JWT token
func (m *RBACJWTManager) GenerateToken(claims auth.auth) (string, error) {
	// 設置過期時間
	claims.SetRegisteredClaims(jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(m.JWTManager.GetExpiresIn()) * time.Hour)),
		NotBefore: jwt.NewNumericDate(time.Now()),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})

	// 創建 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims.(*RBACClaims))
	return token.SignedString([]byte(m.JWTManager.GetSecretKey()))
}

// RefreshToken 刷新 RBAC JWT token
func (m *RBACJWTManager) RefreshToken(tokenString string) (string, error) {
	claims, err := m.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	return m.GenerateToken(claims)
}
