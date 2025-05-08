package auth

import (
	"github.com/POABOB/slack-clone-back-end/services/user-service/internal/domain/user"
)

// LoginRequest 登入結構體
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登入響應 VO
type LoginResponse struct {
	Token string    `json:"token"`
	User  *UserInfo `json:"user"`
}

// UserInfo 用戶資訊 VO
type UserInfo struct {
	ID          uint     `json:"id"`
	Email       string   `json:"email"`
	Username    string   `json:"username"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
}

// TokenResponse Token 響應
type TokenResponse struct {
	Token string `json:"token"`
}

// NewLoginResponse 創建新的登入響應
func NewLoginResponse(token string, user *UserInfo) *LoginResponse {
	return &LoginResponse{
		Token: token,
		User:  user,
	}
}

// NewUserInfo 創建新的用戶資訊
func NewUserInfo(id uint, email, username, role string, permissions []string) *UserInfo {
	return &UserInfo{
		ID:          id,
		Email:       email,
		Username:    username,
		Role:        role,
		Permissions: permissions,
	}
}

// NewTokenResponse 創建新的刷新 Token 響應
func NewTokenResponse(token string) *TokenResponse {
	return &TokenResponse{
		Token: token,
	}
}

// AuthService 驗證邏輯介面
type AuthService interface {
	Register(user *user.User) error
	Login(email, password string) (string, error)
	GenerateToken(user *user.User) (string, error)
}
