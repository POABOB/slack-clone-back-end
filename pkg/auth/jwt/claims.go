package auth

import "github.com/golang-jwt/jwt/v5"

// BaseClaims 基礎 JWT 聲明介面
type BaseClaims interface {
	// GetUserID 獲取用戶 ID
	GetUserID() uint
	// GetEmail 獲取用戶郵箱
	GetEmail() string
	// GetUsername 獲取用戶名
	GetUsername() string
	// GetRegisteredClaims 獲取 jwt 預設 payload
	GetRegisteredClaims() jwt.RegisteredClaims
	// SetRegisteredClaims 設定 jwt 預設 payload
	SetRegisteredClaims(claims jwt.RegisteredClaims)
}
