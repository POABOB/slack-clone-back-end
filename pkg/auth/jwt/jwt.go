package auth

import (
	"errors"
)

var (
	// ErrInvalidToken 無效的 token
	ErrInvalidToken = errors.New("invalid token")
	// ErrExpiredToken 過期的 token
	ErrExpiredToken = errors.New("token has expired")
	// ErrUnauthorized 沒有 Authorization Header
	ErrUnauthorized = errors.New("unauthorized")
	// ErrForbidden Forbidden
	ErrForbidden = errors.New("forbidden")
)

// TokenManager Token 管理介面
type TokenManager interface {
	// GenerateToken 生成 JWT token
	GenerateToken(claims BaseClaims) (string, error)
	// ValidateToken 驗證 JWT token
	ValidateToken(tokenString string) (BaseClaims, error)
	// RefreshToken 刷新 JWT token
	RefreshToken(tokenString string) (string, error)
	// GetExpiresIn 獲取過期時間
	GetExpiresIn() int
	// GetSecretKey 獲取私鑰
	GetSecretKey() string
}
