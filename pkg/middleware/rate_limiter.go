package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/POABOB/slack-clone-back-end/pkg/logger"

	"github.com/gin-gonic/gin"
)

// RateLimitConfig 速率限制配置
type RateLimitConfig struct {
	// 每秒請求數
	RequestsPerSecond int
	// 突發請求數
	Burst int
}

// TODO Each IP
// RateLimiter 速率限制中間件
func RateLimiter(config RateLimitConfig) gin.HandlerFunc {
	// 使用令牌桶算法實現速率限制
	var (
		tokens     = config.Burst
		lastUpdate = time.Now()
		mu         sync.Mutex
	)

	return func(c *gin.Context) {
		mu.Lock()
		defer mu.Unlock()

		// 計算時間差並更新令牌
		now := time.Now()
		elapsed := now.Sub(lastUpdate)
		lastUpdate = now

		// 根據時間差添加令牌
		newTokens := int(elapsed.Seconds() * float64(config.RequestsPerSecond))
		if newTokens > 0 {
			tokens = min(tokens+newTokens, config.Burst)
		}

		// 檢查是否有足夠的令牌
		if tokens <= 0 {
			logger.Warn("rate limit exceeded",
				logger.String("path", c.Request.URL.Path),
				logger.String("method", c.Request.Method),
				logger.String("ip", c.ClientIP()),
			)

			c.JSON(http.StatusTooManyRequests, ErrorResponse{
				Code:    http.StatusTooManyRequests,
				Message: "rate limit exceeded",
			})
			c.Abort()
			return
		}

		// 消耗一個令牌
		tokens--
		c.Next()
	}
}

// min 返回兩個整數中的較小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
