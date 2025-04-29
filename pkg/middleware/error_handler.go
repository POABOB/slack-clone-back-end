package middleware

import (
	"net/http"

	"github.com/POABOB/slack-clone-back-end/pkg/logger"

	"github.com/gin-gonic/gin"
)

// ErrorResponse 錯誤響應結構
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ErrorHandler 全局錯誤處理中間件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 檢查是否有錯誤
		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			// 記錄錯誤
			logger.Error("request error",
				logger.String("path", c.Request.URL.Path),
				logger.String("method", c.Request.Method),
				logger.Error(err),
			)

			// 根據錯誤類型返回適當的響應
			switch err.Type {
			case gin.ErrorTypeBind:
				c.JSON(http.StatusBadRequest, ErrorResponse{
					Code:    http.StatusBadRequest,
					Message: "invalid request parameters",
				})
			case gin.ErrorTypePrivate:
				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Code:    http.StatusInternalServerError,
					Message: "internal server error",
				})
			default:
				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Code:    http.StatusInternalServerError,
					Message: err.Error(),
				})
			}
		}
	}
}
