package config

import (
	"github.com/POABOB/slack-clone-back-end/pkg/middleware"

	"github.com/gin-gonic/gin"
)

// RouterConfig 路由配置
type RouterConfig struct {
	// API 版本
	APIVersion string
	// 是否啟用 CORS
	EnableCORS bool
	// 是否啟用請求日誌
	EnableRequestLog bool
	// 是否啟用錯誤處理
	EnableErrorHandler bool
	// 是否啟用速率限制
	EnableRateLimit bool
	// 速率限制配置
	RateLimitConfig middleware.RateLimitConfig
}

// DefaultRouterConfig 返回默認路由配置
func DefaultRouterConfig() *RouterConfig {
	return &RouterConfig{
		APIVersion:         "v1",
		EnableCORS:         true,
		EnableRequestLog:   true,
		EnableErrorHandler: true,
		EnableRateLimit:    true,
		RateLimitConfig: middleware.RateLimitConfig{
			RequestsPerSecond: 100,
			Burst:             200,
		},
	}
}

// ApplyConfig 應用路由配置
func ApplyConfig(engine *gin.Engine, config *RouterConfig) {
	// 設置模式
	if config.EnableRequestLog {
		engine.Use(gin.Logger())
	}

	// 設置 CORS
	if config.EnableCORS {
		engine.Use(gin.Recovery())
	}

	// 設置錯誤處理
	if config.EnableErrorHandler {
		engine.Use(middleware.ErrorHandler())
	}

	// 設置速率限制
	if config.EnableRateLimit {
		engine.Use(middleware.RateLimiter(config.RateLimitConfig))
	}
}
