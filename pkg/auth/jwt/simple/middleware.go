package simple

import (
	"github.com/POABOB/slack-clone-back-end/pkg/auth/jwt"
	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware JWT 驗證中間件
func JWTAuthMiddleware(jwtManager *JWTManager) gin.HandlerFunc {
	return jwt.NewJWTMiddleware(jwtManager, func(c *gin.Context, claims jwt.BaseClaims) {
		c.Set("user_id", claims.GetUserID())
		c.Set("email", claims.GetEmail())
		c.Set("username", claims.GetUsername())
	})
}
