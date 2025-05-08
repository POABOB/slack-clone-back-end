package handler

import (
	authlib "github.com/POABOB/slack-clone-back-end/pkg/auth"
	"github.com/POABOB/slack-clone-back-end/pkg/auth/jwt/rbac"
	"github.com/POABOB/slack-clone-back-end/services/user-service/internal/domain/user"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService    user.UserService
	rbacMiddleware gin.HandlerFunc
}

// NewUserHandler 創建新的使用者處理器實例
func NewUserHandler(userService user.UserService, rbacMiddleware gin.HandlerFunc) *UserHandler {
	return &UserHandler{
		userService:    userService,
		rbacMiddleware: rbacMiddleware,
	}
}

// TODO 針對單一路由進行速率限制
// RegisterRoutes sets up the user-related routes on the provided RouterGroup with RBAC middleware and permission checks.
func (h *UserHandler) RegisterRoutes(e *gin.RouterGroup) {
	userGroup := e.Group("/user")
	userGroup.Use(h.rbacMiddleware)
	{
		userGroup.GET("/:user_id", rbac.RequirePermission("user:read"), h.GetUser)
		userGroup.PATCH("/:user_id", rbac.RequirePermission("user:update"), h.UpdateUser)
		userGroup.DELETE("/:user_id", rbac.RequirePermission("user:delete"), h.DeleteUser)
	}
}

// GetUser 處理獲取使用者訊息請求
// @Summary 獲取使用者訊息
// @Id User-1
// @Tags User
// @version 1.0
// @accept application/json
// @produce application/json
// @Security BearerAuth
// @param user_id path int true "使用者 ID"
// @Success 200 {objects} user.User
// @Failure 500 {objects} middleware.ErrorResponse
// @Router /api/v1/user/{user_id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.GetUint("user_id")
	singleUser, err := h.userService.GetUserByID(id)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePrivate)
		return
	}

	c.JSON(http.StatusOK, singleUser)
}

// UpdateUser 處理更新使用者訊息請求
// @Summary 更新使用者訊息
// @Id User-2
// @Tags User
// @version 1.0
// @accept application/json
// @produce application/json
// @Security BearerAuth
// @param user_id path int true "使用者 ID"
// @Success 200 {objects} nil
// @Failure 500 {objects} middleware.ErrorResponse
// @Router /api/v1/user/{user_id} [patch]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	var singleUser user.User
	if err := c.ShouldBindJSON(&singleUser); err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	if err := h.userService.UpdateUser(&singleUser); err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePrivate)
		return
	}

	c.JSON(http.StatusOK, nil)
}

// DeleteUser 處理刪除使用者請求
// @Summary 刪除使用者
// @Id User-3
// @Tags User
// @version 1.0
// @accept application/json
// @produce application/json
// @Security BearerAuth
// @param user_id path int true "使用者 ID"
// @Success 200 {objects} nil
// @Failure 500 {objects} middleware.ErrorResponse
// @Router /api/v1/user/{user_id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.GetUint("user_id")
	pathIDString := c.Param("user_id")
	pathID, err := strconv.ParseUint(pathIDString, 10, 64)
	if err != nil {
		_ = c.Error(authlib.ErrInvalidID).SetType(gin.ErrorTypePublic)
		return
	}

	// 確保提交人是自己
	if pathID != uint64(id) {
		_ = c.Error(authlib.ErrForbidden).SetType(gin.ErrorTypePublic)
		return
	}

	if err := h.userService.DeleteUser(id); err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePrivate)
		return
	}

	c.JSON(http.StatusOK, nil)
}
