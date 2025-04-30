package handler

import (
	"net/http"

	auth "github.com/POABOB/slack-clone-back-end/pkg/auth/jwt"
	"github.com/POABOB/slack-clone-back-end/pkg/auth/jwt/rbac"
	"github.com/POABOB/slack-clone-back-end/user-service/internal/domain"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	jwtManager auth.TokenManager
}

func NewAuthHandler(jwtManager auth.TokenManager) *AuthHandler {
	return &AuthHandler{
		jwtManager: jwtManager,
	}
}

// TODO OAUTH 之類的登入、註冊

// Login 處理用戶登入
func (h *AuthHandler) Login(c *gin.Context) {
	var loginRequest domain.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	// TODO: 驗證用戶憑證
	// 這裡應該從資料庫中驗證用戶
	//if err != nil {
	//	_ = c.Error(err).SetType(gin.ErrorTypePrivate)
	//	return
	//}

	// 生成 token
	token, err := h.jwtManager.GenerateToken(rbac.NewRBACClaims(
		1,
		loginRequest.Email,
		"test_user",
		"admin",
		[]string{"read", "write", "delete"}),
	)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePrivate)
		return
	}

	c.JSON(http.StatusOK, domain.NewLoginResponse(token, domain.NewUserInfo(
		1,
		loginRequest.Email,
		"test_user",
		"admin",
		[]string{"read", "write", "delete"})),
	)
}

// RefreshToken 刷新 token
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// 從 context 中獲取用戶信息
	userId := c.MustGet("user_id").(uint)
	email := c.MustGet("email").(string)
	username := c.MustGet("username").(string)
	role := c.MustGet("role").(string)
	permissions := c.MustGet("permissions").([]string)

	// 生成新的 token
	token, err := h.jwtManager.GenerateToken(rbac.NewRBACClaims(
		userId,
		email,
		username,
		role,
		permissions),
	)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePrivate)
		return
	}

	c.JSON(http.StatusOK, domain.NewTokenResponse(token))
}

// GetUserInfo 獲取用戶信息
func (h *AuthHandler) GetUserInfo(c *gin.Context) {
	userId := c.MustGet("user_id").(uint)
	email := c.MustGet("email").(string)
	username := c.MustGet("username").(string)
	role := c.MustGet("role").(string)
	permissions := c.MustGet("permissions").([]string)

	userInfo := domain.NewUserInfo(userId, email, username, role, permissions)
	c.JSON(http.StatusOK, userInfo)
}
