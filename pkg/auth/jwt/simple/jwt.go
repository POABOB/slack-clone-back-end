package simple

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/POABOB/slack-clone-back-end/pkg/auth"
	jwtlib "github.com/POABOB/slack-clone-back-end/pkg/auth/jwt"
	"github.com/POABOB/slack-clone-back-end/pkg/config"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// JWTManager JWT 管理器
type JWTManager struct {
	secretKey []byte
	expiresIn int
}

// NewJWTManager 創建新的 JWT 管理器
func NewJWTManager(cfg *config.JWTConfig) *JWTManager {
	if cfg.ExpiresIn == 0 {
		cfg.ExpiresIn = 1 * 24 * 60 * 60 * 1000
	}
	return &JWTManager{
		secretKey: cfg.SecretKey,
		expiresIn: cfg.ExpiresIn,
	}
}

// GenerateToken 生成 JWT token
func (m *JWTManager) GenerateToken(claims jwtlib.BaseClaims) (string, error) {
	// 設置過期時間
	now := time.Now()
	claims.SetRegisteredClaims(jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(m.GetExpiresIn()) * time.Millisecond)),
		NotBefore: jwt.NewNumericDate(now),
		IssuedAt:  jwt.NewNumericDate(now),
		ID:        GenerateRandomID(), // 每次 Refresh 都可以產生不一樣的 Token，不會被 Date 所侷限
	})

	// 創建 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims.(*DefaultClaims))
	return token.SignedString(m.GetSecretKey())
}

// ValidateToken 驗證 JWT token
func (m *JWTManager) ValidateToken(tokenString string) (jwtlib.BaseClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &DefaultClaims{}, func(token *jwt.Token) (interface{}, error) {
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
func (m *JWTManager) GetSecretKey() []byte {
	return m.secretKey
}

// GenerateRandomID 生成一個隨機的 ID
func GenerateRandomID() string {
	b := make([]byte, 8) // 使用 8 字節，生成 11 字符的 base64 字符串
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
