package jwt

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
	GetSecretKey() []byte
}
