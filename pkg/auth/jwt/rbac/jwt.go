package rbac

import (
	"errors"
	"fmt"
	"github.com/POABOB/slack-clone-back-end/pkg/auth"
	jwtlib "github.com/POABOB/slack-clone-back-end/pkg/auth/jwt"
	"github.com/POABOB/slack-clone-back-end/pkg/auth/jwt/simple"
	"github.com/POABOB/slack-clone-back-end/pkg/config"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// RBACJWTManager RBAC JWT 管理器
type RBACJWTManager struct {
	*simple.JWTManager
}

// NewRBACJWTManager 創建新的 RBAC JWT 管理器
func NewRBACJWTManager(cfg *config.JWTConfig) *RBACJWTManager {
	return &RBACJWTManager{
		JWTManager: simple.NewJWTManager(cfg),
	}
}

// GenerateToken 生成 RBAC JWT token
func (m *RBACJWTManager) GenerateToken(claims jwtlib.BaseClaims) (string, error) {
	// 設置過期時間
	now := time.Now()
	claims.SetRegisteredClaims(jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(m.JWTManager.GetExpiresIn()) * time.Millisecond)),
		NotBefore: jwt.NewNumericDate(now),
		IssuedAt:  jwt.NewNumericDate(now),
		ID:        simple.GenerateRandomID(), // 每次 Refresh 都可以產生不一樣的 Token，不會被 Date 所侷限
	})

	// 創建 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims.(*RBACClaims))
	return token.SignedString(m.JWTManager.GetSecretKey())
}

// ValidateToken 驗證 JWT token
func (m *RBACJWTManager) ValidateToken(tokenString string) (jwtlib.BaseClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RBACClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.GetSecretKey(), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, auth.ErrExpiredToken
		}
		return nil, auth.ErrInvalidToken
	}

	claims, ok := token.Claims.(*RBACClaims)
	if !ok || !token.Valid {
		return nil, auth.ErrInvalidToken
	}

	return claims, nil
}

// RefreshToken 刷新 RBAC JWT token
func (m *RBACJWTManager) RefreshToken(tokenString string) (string, error) {
	claims, err := m.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	return m.GenerateToken(claims)
}
