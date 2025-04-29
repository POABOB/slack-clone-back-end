# Slack Clone Backend

這是一個基於微服務架構的 Slack Clone 後端專案。

## 微服務專案結構

```
.
├── pkg/                    # 共用套件
│   ├── auth/              # 身份驗證相關
│   ├── config/            # 配置管理
│   ├── database/          # 資料庫連接
│   ├── logger/            # 日誌系統
│   ├── redis/             # Redis 客戶端
│   └── utils/             # 通用工具
│
├── services/              # 微服務
│   └── user-service/      # 用戶服務
│       ├── cmd/          # 主程式進入點
│       ├── internal/     # 內部包
│       │   ├── domain/   # 領域模型
│       │   ├── handler/  # HTTP 處理器
│       │   ├── repository/ # 資料存取層
│       │   └── service/  # 業務邏輯層
│       └── config/      # 服務配置
│
├── scripts/               # 腳本文件
└── README.md
```

## 共用套件說明

### pkg/auth
- JWT 相關功能
- 權限驗證
- 用戶認證

### pkg/config
- 配置管理
- 環境變數處理
- 配置驗證

### pkg/database
- 資料庫連接池
- 資料庫配置
- 遷移工具

### pkg/logger
- 日誌配置
- 日誌格式化
- 日誌輸出

### pkg/redis
- Redis 連接池
- Redis 配置
- 快取工具

### pkg/utils
- 通用工具函數
- 錯誤處理
- 驗證工具

## 微服務說明

### user-service
- 用戶管理
- 身份驗證
- 權限控制

## 開發環境設置

1. 安裝依賴：
```bash
go mod download
```

2. 運行服務：
```bash
# 運行用戶服務
go run services/user-service/cmd/main.go
```

## 測試

```bash
go test ./...
``` 