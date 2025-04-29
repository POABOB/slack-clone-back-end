package database

import (
	"fmt"

	"github.com/POABOB/slack-clone-back-end/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	// DB 全局資料庫實例
	DB *gorm.DB
)

// InitDatabase 初始化資料庫連接
func InitDatabase(cfg *config.DatabaseConfig) error {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect database: %v", err)
	}

	// 設置連接池
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %v", err)
	}

	// 設置最大閒置連接數
	sqlDB.SetMaxIdleConns(10)
	// 設置最大打開連接數
	sqlDB.SetMaxOpenConns(100)

	return nil
}

// GetDB 獲取資料庫實例
func GetDB() *gorm.DB {
	return DB
}

// Close 關閉資料庫連接
func Close() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
