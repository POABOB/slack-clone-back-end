package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/POABOB/slack-clone-back-end/pkg/config"
	"github.com/go-redis/redis/v8"
)

var (
	// Client 全局 Redis 客戶端實例
	Client *redis.Client
)

// InitRedis 初始化 Redis 連接
func InitRedis(cfg *config.RedisConfig) error {
	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// 測試連接
	ctx := context.Background()
	if err := Client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect redis: %v", err)
	}

	return nil
}

// GetClient 獲取 Redis 客戶端實例
func GetClient() *redis.Client {
	return Client
}

// Close 關閉 Redis 連接
func Close() error {
	if Client != nil {
		return Client.Close()
	}
	return nil
}

// Set 設置鍵值對
func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return Client.Set(ctx, key, value, expiration).Err()
}

// Get 獲取值
func Get(ctx context.Context, key string) (string, error) {
	return Client.Get(ctx, key).Result()
}

// Delete 刪除鍵
func Delete(ctx context.Context, key string) error {
	return Client.Del(ctx, key).Err()
}

// Exists 檢查鍵是否存在
func Exists(ctx context.Context, key string) (bool, error) {
	result, err := Client.Exists(ctx, key).Result()
	return result > 0, err
}
