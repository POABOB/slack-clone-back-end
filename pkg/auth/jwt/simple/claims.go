package simple

import (
	"github.com/golang-jwt/jwt/v5"
)

// DefaultClaims 默認 JWT 聲明實現
type DefaultClaims struct {
	UserID   uint   `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// NewDefaultClaims 創建新的 JWT
func NewDefaultClaims(userID uint, email, username string) *DefaultClaims {
	return &DefaultClaims{
		UserID:   userID,
		Email:    email,
		Username: username,
	}
}

// GetUserID 獲取 UserID
func (c *DefaultClaims) GetUserID() uint {
	return c.UserID
}

// GetEmail 獲取 Email
func (c *DefaultClaims) GetEmail() string {
	return c.Email
}

// GetUsername 獲取 Username
func (c *DefaultClaims) GetUsername() string {
	return c.Username
}

// GetRegisteredClaims 獲取 jwt 預設 payload
func (c *DefaultClaims) GetRegisteredClaims() jwt.RegisteredClaims {
	return c.RegisteredClaims
}

// SetRegisteredClaims 設定 jwt 預設 payload
func (c *DefaultClaims) SetRegisteredClaims(claims jwt.RegisteredClaims) { c.RegisteredClaims = claims }
