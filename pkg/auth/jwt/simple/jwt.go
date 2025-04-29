package simple

import (
	"errors"
	"fmt"
	"github.com/POABOB/slack-clone-back-end/pkg/auth/jwt"
	"github.com/POABOB/slack-clone-back-end/pkg/jwt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// JWTManager JWT 管理器
type JWTManager struct {
	secretKey string
	expiresIn int
}

// NewJWTManager 創建新的 JWT 管理器
func NewJWTManager(secretKey string, expiresIn int) *JWTManager {
	return &JWTManager{
		secretKey: secretKey,
		expiresIn: expiresIn,
	}
}

// GenerateToken 生成 JWT token
func (m *JWTManager) GenerateToken(claims auth.auth) (string, error) {
	// 設置過期時間
	claims.SetRegisteredClaims(jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(m.GetExpiresIn()) * time.Hour)),
		NotBefore: jwt.NewNumericDate(time.Now()),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})

	// 創建 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims.(*DefaultClaims))
	return token.SignedString([]byte(m.GetSecretKey()))
}

// ValidateToken 驗證 JWT token
func (m *JWTManager) ValidateToken(tokenString string) (auth.BaseClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &DefaultClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.GetSecretKey()), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, auth.ErrExpiredToken
		}
		return nil, auth.ErrInvalidToken
	}

	claims, ok := token.Claims.(*DefaultClaims)
	if !ok || !token.Valid {
		return nil, auth.ErrInvalidToken
	}

	return claims, nil
}

// RefreshToken 刷新 JWT token
func (m *JWTManager) RefreshToken(tokenString string) (string, error) {
	claims, err := m.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	return m.GenerateToken(claims)
}

// GetExpiresIn 獲取過期時間
func (m *JWTManager) GetExpiresIn() int {
	return m.expiresIn
}

// GetSecretKey 獲取私鑰
func (m *JWTManager) GetSecretKey() string {
	return m.secretKey
}
