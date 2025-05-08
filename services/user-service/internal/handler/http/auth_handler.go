package handler

import (
	"github.com/POABOB/slack-clone-back-end/services/user-service/internal/domain/auth"
	"github.com/POABOB/slack-clone-back-end/services/user-service/internal/domain/user"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthHandler struct {
	authService    auth.AuthService
	rbacMiddleware gin.HandlerFunc
}

func NewAuthHandler(authService auth.AuthService, rbacMiddleware gin.HandlerFunc) *AuthHandler {
	return &AuthHandler{
		authService:    authService,
		rbacMiddleware: rbacMiddleware,
	}
}

// RegisterRoutes sets up the auth-related routes on the provided RouterGroup with RBAC middleware and permission checks.
func (h *AuthHandler) RegisterRoutes(e *gin.RouterGroup) {
	authGroup := e.Group("/auth")

	authGroup.POST("/register", h.Register)
	authGroup.POST("/register", h.Login)
	authGroup.Use(h.rbacMiddleware)
	{
		authGroup.DELETE("/refresh", h.RefreshToken)
		authGroup.DELETE("/info", h.GetUserInfo)
	}
}

// TODO OAUTH 之類的登入、註冊
// TODO service、repo 層錯誤處理

// Register 處理使用者註冊請求
func (h *AuthHandler) Register(c *gin.Context) {
	var singleUser user.User
	if err := c.ShouldBindJSON(&singleUser); err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	if err := h.authService.Register(&singleUser); err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePrivate)
		return
	}

	c.JSON(http.StatusCreated, nil)
}

// Login 處理使用者登入
func (h *AuthHandler) Login(c *gin.Context) {
	var loginRequest auth.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	token, err := h.authService.Login(loginRequest.Email, loginRequest.Password)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePrivate)
		return
	}

	c.JSON(http.StatusOK, auth.NewTokenResponse(token))
}

// RefreshToken 刷新 token
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	userId := c.MustGet("user_id").(uint)
	email := c.MustGet("email").(string)
	username := c.MustGet("username").(string)
	role := c.MustGet("role").(string)
	permissions := c.MustGet("permissions").([]string)
	token, err := h.authService.GenerateToken(&user.User{
		ID:          userId,
		Email:       email,
		Username:    username,
		Role:        role,
		Permissions: permissions,
	})
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePrivate)
		return
	}

	c.JSON(http.StatusOK, auth.NewTokenResponse(token))
}

// GetUserInfo 獲取使用者訊息
func (h *AuthHandler) GetUserInfo(c *gin.Context) {
	userId := c.MustGet("user_id").(uint)
	email := c.MustGet("email").(string)
	username := c.MustGet("username").(string)
	role := c.MustGet("role").(string)
	permissions := c.MustGet("permissions").([]string)
	userInfo := auth.NewUserInfo(userId, email, username, role, permissions)

	c.JSON(http.StatusOK, userInfo)
}
