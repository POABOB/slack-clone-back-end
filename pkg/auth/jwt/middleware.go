package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// ClaimsHandler 是一個自定義處理驗證後 claims 的函數型別
type ClaimsHandler func(c *gin.Context, claims BaseClaims)

// NewJWTMiddleware 回傳一個 JWT 驗證中間件
func NewJWTMiddleware(jwtManager TokenManager, handler ClaimsHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": ErrInvalidToken})
			_ = c.Error(ErrUnauthorized).SetType(gin.ErrorTypePrivate)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": ErrInvalidToken})
			_ = c.Error(ErrInvalidToken).SetType(gin.ErrorTypePrivate)
			return
		}

		claims, err := jwtManager.ValidateToken(parts[1])
		if err != nil {
			if errors.Is(err, ErrExpiredToken) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": ErrExpiredToken})
			} else {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": ErrInvalidToken})
			}
			_ = c.Error(err).SetType(gin.ErrorTypePrivate)
			c.Abort()
			return
		}

		// 呼叫自定義處理函數
		handler(c, claims)
		c.Next()
	}
}
