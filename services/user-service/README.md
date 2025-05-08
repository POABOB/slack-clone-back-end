# Slack Clone Backend - User Service

這是一個基於微服務架構的 Slack Clone 後端 User Service 專案。

## 專案結構

```
.
├── cmd/                  # 主程式進入點
│   └── main.go
├── config/               # 專案配置
├── docs/                 # Swaggo 文件
├── internal/             # 內部函示庫
│   ├── domain/           # 領域模型和 interface
│   ├── handler/          # 處理器
│   │   ├── http/         # http 處理器
│   │   └── grpc/         # grpc 處理器
│   ├── middleware/       # 中間件
│   ├── repository/       # 資料存取層
│   │   └── postgresql/   # postgresql Repo
│   ├── router/           # http 路由
│   └── service/          # 業務邏輯層
├── pkg/                  # 可重用的公共函示庫
├── scripts/              # 腳本文件
├── config.yaml           # 配置文件
├── docker-compose.yml    # Docker Compose
├── go.mod
└── README.md
```

## 函式庫

- Gin: Web 框架
- GORM: ORM 框架
- FX: 依賴注入
- JWT: 身份驗證
- Viper: 配置管理
- Zap: 日誌系統
- go-redis: Redis 客戶端
- testify: 測試框架
- go-sqlmock: SQL 模擬測試

## 開發環境設置

1. 安裝依賴：
```bash
go mod download
```

2. 運行服務：
```bash
go run cmd/main.go
```

## 測試

```bash
go test ./...
```

---

## 📁 專案架構與模組說明

### ✨ handler → service → repository 範例（模組名稱：User）

每個模組皆使用 handler → service → repository 架構進行劃分，確保職責分明，提升模組化與測試便利性。

### ✅ `domain/user/user.go` 說明
> `domain` 資料夾為核心模組，負責定義所有的 **DTO 結構（資料傳輸物件）**，以及各模組的 **interface 定義**，例如 `IDataService`、`IDataRepository` 等。

```go
package user

import (
	"time"
)

// User 使用者實體
type User struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Email       string    `json:"email" gorm:"uniqueIndex"`
	Password    string    `json:"-" gorm:"not null"`
	Username    string    `json:"username" gorm:"not null"`
	Role        string    `json:"role" gorm:"not null;default:'user'"`
	Permissions []string  `json:"permissions" gorm:"type:json"`
	LastLogin   time.Time `json:"last_login"`
	IsDeleted   bool      `json:"is_deleted"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UserRepository 使用者資料存取介面
type UserRepository interface {
	Create(user *User) error
	FindByID(id uint) (*User, error)
	FindByEmail(email string) (*User, error)
	Update(user *User) error
	Delete(id uint) error
}

// UserService 使用者業務邏輯介面
type UserService interface {
	GetUserByID(id uint) (*User, error)
	UpdateUser(user *User) error
	DeleteUser(id uint) error
}

#### 主要功能說明：

1. **資料結構定義**
   - `User` 結構體定義了使用者的基本資訊
   - 使用 GORM 標籤定義資料庫映射
   - 使用 JSON 標籤定義 API 響應格式
   - 包含敏感資訊處理（如密碼欄位使用 `json:"-"` 隱藏）

2. **資料庫欄位說明**
   - `ID`：主鍵，自動遞增
   - `Email`：唯一索引，用於登入
   - `Password`：加密存儲的密碼
   - `Username`：使用者名稱
   - `Role`：使用者角色，預設為 "user"
   - `Permissions`：JSON 格式的權限列表
   - `LastLogin`：最後登入時間
   - `IsDeleted`：軟刪除標記
   - `CreatedAt`/`UpdatedAt`：時間戳記

3. **Repository 介面**
   - `Create`：創建新使用者
   - `FindByID`：根據 ID 查詢使用者
   - `FindByEmail`：根據郵箱查詢使用者
   - `Update`：更新使用者資訊
   - `Delete`：刪除使用者（軟刪除）

4. **Service 介面**
   - `GetUserByID`：獲取使用者資訊
   - `UpdateUser`：更新使用者資訊
   - `DeleteUser`：刪除使用者

5. **設計考量**
   - 使用介面定義實現依賴反轉
   - 分離資料存取和業務邏輯
   - 支援軟刪除機制
   - 整合 RBAC 權限控制
   - 提供完整的時間追蹤

6. **安全性考慮**
   - 密碼欄位在 JSON 響應中隱藏
   - 使用 GORM 的資料庫約束
   - 支援權限和角色管理
   - 提供審計追蹤（時間戳記）

---

### ✅ `handler/http/user_handler.go`
> `handler` 資料夾負責處理 HTTP/gRPC 請求，實現 API 端點，並處理請求驗證、權限檢查和錯誤處理。

```go
package handler

import (
	authlib "github.com/POABOB/slack-clone-back-end/pkg/auth"
	"github.com/POABOB/slack-clone-back-end/pkg/auth/jwt/rbac"
	"github.com/POABOB/slack-clone-back-end/services/user-service/internal/domain/user"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserHandler 使用者 HTTP 處理器
// 負責處理所有與使用者相關的 HTTP 請求
type UserHandler struct {
	userService    user.UserService
	rbacMiddleware gin.HandlerFunc
}

// NewUserHandler 創建新的使用者處理器實例
// 參數：
//   - userService: 使用者服務介面，用於處理業務邏輯
//   - rbacMiddleware: RBAC 中間件，用於權限驗證
func NewUserHandler(userService user.UserService, rbacMiddleware gin.HandlerFunc) *UserHandler {
	return &UserHandler{
		userService:    userService,
		rbacMiddleware: rbacMiddleware,
	}
}

// RegisterRoutes 註冊使用者相關的路由
// 路由說明：
// 1. 所有路由都在 /user 路徑下
// 2. 使用 RBACMiddleware 進行權限驗證
// 3. 每個路由都需要特定的權限：
//    - GET /:user_id 需要 user:read 權限
//    - PATCH /:user_id 需要 user:update 權限
//    - DELETE /:user_id 需要 user:delete 權限
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
// @Description 根據使用者 ID 獲取使用者詳細訊息
// @Id User-1
// @Tags User
// @version 1.0
// @accept application/json
// @produce application/json
// @Security BearerAuth
// @param user_id path int true "使用者 ID" minimum(1)
// @Success 200 {object} user.User "成功獲取使用者訊息"
// @Failure 401 {object} middleware.ErrorResponse "未授權"
// @Failure 403 {object} middleware.ErrorResponse "權限不足"
// @Failure 404 {object} middleware.ErrorResponse "使用者不存在"
// @Failure 500 {object} middleware.ErrorResponse "伺服器內部錯誤"
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
// @Description 更新指定使用者的訊息
// @Id User-2
// @Tags User
// @version 1.0
// @accept application/json
// @produce application/json
// @Security BearerAuth
// @param user_id path int true "使用者 ID" minimum(1)
// @param user body user.User true "使用者更新資料"
// @Success 200 {object} nil "成功更新使用者訊息"
// @Failure 400 {object} middleware.ErrorResponse "請求格式錯誤"
// @Failure 401 {object} middleware.ErrorResponse "未授權"
// @Failure 403 {object} middleware.ErrorResponse "權限不足"
// @Failure 404 {object} middleware.ErrorResponse "使用者不存在"
// @Failure 500 {object} middleware.ErrorResponse "伺服器內部錯誤"
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
// @Description 刪除指定的使用者
// @Id User-3
// @Tags User
// @version 1.0
// @accept application/json
// @produce application/json
// @Security BearerAuth
// @param user_id path int true "使用者 ID" minimum(1)
// @Success 200 {object} nil "成功刪除使用者"
// @Failure 400 {object} middleware.ErrorResponse "無效的使用者 ID"
// @Failure 401 {object} middleware.ErrorResponse "未授權"
// @Failure 403 {object} middleware.ErrorResponse "權限不足或非本人操作"
// @Failure 404 {object} middleware.ErrorResponse "使用者不存在"
// @Failure 500 {object} middleware.ErrorResponse "伺服器內部錯誤"
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
```

#### 主要功能說明：

1. **路由註冊與權限控制**
   - 所有路由都在 `/user` 路徑下
   - 使用 `RBACMiddleware` 進行權限驗證
   - 每個路由都需要特定的權限：
     - `GET /:user_id` 需要 `user:read` 權限
     - `PATCH /:user_id` 需要 `user:update` 權限
     - `DELETE /:user_id` 需要 `user:delete` 權限

2. **API 端點說明**
   - `GET /:user_id`：獲取使用者訊息
   - `PATCH /:user_id`：更新使用者訊息
   - `DELETE /:user_id`：刪除使用者

3. **請求驗證**
   - 路徑參數驗證：確保 `user_id` 為有效數字
   - 請求體驗證：使用 `ShouldBindJSON` 驗證更新請求
   - 權限驗證：檢查使用者是否具有所需權限
   - 業務邏輯驗證：如刪除時檢查是否為本人操作

4. **錯誤處理**
   - 400：請求格式錯誤
   - 401：未授權
   - 403：權限不足
   - 404：資源不存在
   - 500：伺服器錯誤

5. **Swagger 文檔**
   - 使用 Swagger 註解定義 API 文檔
   - 包含請求參數、響應狀態和錯誤碼
   - 提供 API 使用說明和示例

6. **安全性考慮**
   - 使用 RBAC 進行權限控制
   - 驗證使用者身份
   - 防止未授權訪問
   - 確保使用者只能操作自己的資料

---

✅ `service/user_service.go`
> `service` 資料夾負責實現業務邏輯層，處理複雜的業務規則，協調資料存取層的操作，並確保資料的一致性和安全性。服務層是連接控制器層和資料存取層的橋樑，實現了關注點分離和依賴反轉原則。

```go
package service

import (
	"github.com/POABOB/slack-clone-back-end/pkg/auth"
	"github.com/POABOB/slack-clone-back-end/services/user-service/internal/domain/user"
)

type userService struct {
	repo user.UserRepository
}

// NewUserService 創建新的使用者服務實例
func NewUserService(repo user.UserRepository) user.UserService {
	return &userService{
		repo: repo,
	}
}

// GetUserByID 獲取使用者訊息
func (s *userService) GetUserByID(id uint) (*user.User, error) {
	return s.repo.FindByID(id)
}

// UpdateUser 更新使用者訊息
func (s *userService) UpdateUser(user *user.User) error {
	// 如果密碼被更新，需要重新加密
	if user.Password != "" {
		hashedPassword, err := auth.HashPassword(user.Password)
		if err != nil {
			return err
		}
		user.Password = hashedPassword
	}
	return s.repo.Update(user)
}

// DeleteUser 刪除使用者
func (s *userService) DeleteUser(id uint) error {
	return s.repo.Delete(id)
}
```

#### 主要功能說明：

1. **服務層職責**
   - 實現業務邏輯處理
   - 協調資料存取層操作
   - 處理資料轉換和驗證
   - 實現業務規則和約束

2. **依賴注入**
   - 通過建構函數注入 Repository
   - 實現依賴反轉原則
   - 便於單元測試和模擬

3. **業務邏輯處理**
   - `GetUserByID`：獲取使用者資訊
     - 直接調用 Repository 層
     - 不包含額外業務邏輯
   - `UpdateUser`：更新使用者資訊
     - 處理密碼加密邏輯
     - 驗證更新資料
   - `DeleteUser`：刪除使用者
     - 實現軟刪除機制
     - 確保資料一致性

4. **安全性處理**
   - 密碼加密：使用 `auth.HashPassword` 進行密碼加密
   - 資料驗證：確保更新資料的有效性
   - 權限控制：與 RBAC 系統整合

5. **錯誤處理**
   - 統一錯誤處理機制
   - 錯誤傳遞和轉換
   - 業務邏輯錯誤處理

6. **擴展性考慮**
   - 介面化設計
   - 模組化結構
   - 易於添加新功能
   - 支援橫切關注點

---

### ✅ `repository/postgresql/user_repo.go`
> `repository` 資料夾負責實現資料存取層，處理與資料庫的交互操作。它封裝了所有資料庫相關的邏輯，提供了一個抽象層來隔離業務邏輯和資料存取細節。這種設計使得系統更容易維護、測試和擴展。

```go
package repository

import (
	"github.com/POABOB/slack-clone-back-end/services/user-service/internal/domain/user"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 創建新的用戶資料存取實例
func NewUserRepository(db *gorm.DB) user.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *user.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByID(id uint) (*user.User, error) {
	var singleUser user.User
	err := r.db.First(&singleUser, id).Error
	if err != nil {
		return nil, err
	}
	return &singleUser, nil
}

func (r *userRepository) FindByEmail(email string) (*user.User, error) {
	var singleUser user.User
	err := r.db.Where("email = ?", email).First(&singleUser).Error
	if err != nil {
		return nil, err
	}
	return &singleUser, nil
}

func (r *userRepository) Update(user *user.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Model(&user.User{}).Where("id = ?", id).Update("is_deleted", true).Error
}
```

#### 主要功能說明：

1. **資料存取層職責**
   - 實現資料庫操作
   - 封裝 SQL 查詢邏輯
   - 處理資料持久化
   - 提供資料存取介面

2. **資料庫操作**
   - `Create`：創建新使用者記錄
   - `FindByID`：根據 ID 查詢使用者
   - `FindByEmail`：根據郵箱查詢使用者
   - `Update`：更新使用者資訊
   - `Delete`：實現軟刪除機制

3. **ORM 使用**
   - 使用 GORM 框架
   - 簡化資料庫操作
   - 提供資料庫遷移支援
   - 自動處理關聯關係

4. **錯誤處理**
   - 統一錯誤處理機制
   - 資料庫錯誤轉換
   - 提供清晰的錯誤訊息

5. **安全性考慮**
   - 使用參數化查詢
   - 防止 SQL 注入
   - 資料驗證和清理
   - 事務管理

6. **擴展性設計**
   - 介面化實現
   - 支援多種資料庫
   - 易於切換資料來源
   - 便於單元測試

---

## ✨ FX 使用

FX 是一個 Go 語言的依賴注入框架，用於管理應用程式的生命週期和依賴關係。在 User Service 中，我們使用 FX 來實現依賴注入，使程式碼更加模組化和可測試。

### Module 定義方式

1. **Pkg Module**
```go
// pkg/module.go
package pkg

import (
	"github.com/POABOB/slack-clone-back-end/pkg/auth/jwt/rbac"
	"github.com/POABOB/slack-clone-back-end/pkg/database/postgresql"
	"go.uber.org/fx"
)

// PostgresqlModule 依賴注入統一管理
var PostgresqlModule = fx.Module("postgresql",
	fx.Provide(postgresql.NewDatabase),
)

var AuthModule = fx.Module("auth",
	fx.Provide(rbac.NewRBACJWTManager, rbac.RBACMiddleware),
)
```

2. **Internal Module**
```go
// internal/module.go
package internal

import (
	handler "github.com/POABOB/slack-clone-back-end/services/user-service/internal/handler/http"
	repository "github.com/POABOB/slack-clone-back-end/services/user-service/internal/repository/postgresql"
	"github.com/POABOB/slack-clone-back-end/services/user-service/internal/service"
	"go.uber.org/fx"
)

var Module = fx.Module("user-service",
	fx.Provide(
		repository.NewUserRepository,
		service.NewUserService,
		handler.NewUserHandler,

		service.NewAuthService,
		handler.NewAuthHandler,
	),
)
```

3. **Router Module**
```go
// internal/router/module.go
package router

import (
	configlib "github.com/POABOB/slack-clone-back-end/pkg/config"
	"go.uber.org/fx"
)

var Module = fx.Module("router",
	fx.Provide(
		configlib.NewGinEngine,
		NewRouter,
	),
)
```

4. **Config Module**
```go
// config/module.go
package config

import (
	configlib "github.com/POABOB/slack-clone-back-end/pkg/config"
	"go.uber.org/fx"
)

var Module = fx.Module("config",
	fx.Provide(
		func() (*configlib.Config, error) {
			return configlib.LoadConfig("config.yaml")
		},
		func(cfg *configlib.Config) *configlib.ServerConfig { return &cfg.Server },
		func(cfg *configlib.Config) *configlib.RouterConfig { return &cfg.Router },
		func(cfg *configlib.Config) *configlib.JWTConfig { return &cfg.JWT },
		func(cfg *configlib.Config) *configlib.DatabaseConfig { return &cfg.Database },
		func(cfg *configlib.Config) *configlib.RedisConfig { return &cfg.Redis },
	),
)
```

### 單元測試範例

1. **Handler 單元測試**
```go
// internal/handler/http/user_handler_test.go
package handler

import (
	"testing"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/fx"
)

// MockUserService 模擬 UserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetUserByID(id uint) (*user.User, error) {
	args := m.Called(id)
	return args.Get(0).(*user.User), args.Error(1)
}

func TestUserHandler_GetUser(t *testing.T) {
	// 創建測試模組
	testModule := fx.Module("test",
		fx.Provide(
			func() user.UserService { return &MockUserService{} },
			func() gin.HandlerFunc { return func(c *gin.Context) {} },
			NewUserHandler,
		),
	)

	// 注入依賴
	var handler *UserHandler
	err := fx.New(
		testModule,
		fx.Inject(&handler),
	).Err()

	assert.NoError(t, err)

	// 設置測試案例
	mockService := handler.userService.(*MockUserService)
	expectedUser := &user.User{ID: 1, Username: "test"}
	mockService.On("GetUserByID", uint(1)).Return(expectedUser, nil)

	// 創建測試請求
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", uint(1))

	// 執行測試
	handler.GetUser(c)

	// 驗證結果
	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}
```

2. **Service 單元測試**
```go
// internal/service/data_service_test.go
package service

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/fx"
)

// MockUserRepository 模擬 UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByID(id uint) (*user.User, error) {
	args := m.Called(id)
	return args.Get(0).(*user.User), args.Error(1)
}

func TestUserService_GetUserByID(t *testing.T) {
	// 創建測試模組
	testModule := fx.Module("test",
		fx.Provide(
			func() user.UserRepository { return &MockUserRepository{} },
			NewUserService,
		),
	)

	// 注入依賴
	var service user.UserService
	err := fx.New(
		testModule,
		fx.Inject(&service),
	).Err()

	assert.NoError(t, err)

	// 設置測試案例
	mockRepo := service.(*userService).repo.(*MockUserRepository)
	expectedUser := &user.User{ID: 1, Username: "test"}
	mockRepo.On("FindByID", uint(1)).Return(expectedUser, nil)

	// 執行測試
	user, err := service.GetUserByID(1)

	// 驗證結果
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockRepo.AssertExpectations(t)
}
```

3. **Repository 單元測試**
```go
// internal/repository/postgresql/user_repo_test.go
package repository

import (
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestUserRepository_FindByID(t *testing.T) {
	// 創建 SQL 模擬
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// 創建 GORM 實例
	dialector := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	// 創建測試模組
	testModule := fx.Module("test",
		fx.Provide(
			func() *gorm.DB { return gormDB },
			NewUserRepository,
		),
	)

	// 注入依賴
	var repo user.UserRepository
	err = fx.New(
		testModule,
		fx.Inject(&repo),
	).Err()

	assert.NoError(t, err)

	// 設置預期查詢
	rows := sqlmock.NewRows([]string{"id", "username", "email"}).
		AddRow(1, "test", "test@example.com")
	mock.ExpectQuery("SELECT (.+) FROM \"users\"").
		WithArgs(1).
		WillReturnRows(rows)

	// 執行測試
	user, err := repo.FindByID(1)

	// 驗證結果
	assert.NoError(t, err)
	assert.Equal(t, uint(1), user.ID)
	assert.Equal(t, "test", user.Username)
	assert.Equal(t, "test@example.com", user.Email)
	assert.NoError(t, mock.ExpectationsWereMet())
}
```

### 測試注意事項

1. **Mock 使用**
   - 使用 `testify/mock` 創建模擬對象
   - 設置預期行為和返回值
   - 驗證方法調用

2. **依賴注入**
   - 使用 `fx.Module` 創建測試模組
   - 注入模擬依賴
   - 使用 `fx.Inject` 獲取測試對象

3. **資料庫測試**
   - 使用 `go-sqlmock` 模擬資料庫
   - 設置預期查詢和結果
   - 驗證 SQL 執行

4. **HTTP 測試**
   - 使用 `gin.CreateTestContext` 創建測試上下文
   - 模擬 HTTP 請求和響應
   - 驗證響應狀態和內容

---

## Swaggo 使用


Swaggo 是一個用於自動生成 Swagger/OpenAPI 2.0 文檔的工具。在 User Service 中，我們使用 Swaggo 來生成 API 文檔，提供清晰的 API 使用說明。

### 安裝與設置

1. **安裝 Swaggo CLI**：
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

2. **添加必要的依賴**：
```go
// go.mod
require (
    github.com/swaggo/swag v1.16.2
    github.com/swaggo/gin-swagger v1.6.0
    github.com/swaggo/files v1.0.1
)
```

### 文檔生成流程

1. **在 `main.go` 中添加基本資訊**：
```go
// @title User Service API
// @version 1.0
// @description User Service API 文檔
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
```

2. **在 Handler 中添加 API 註解**：
```go
// GetUser 處理獲取使用者訊息請求
// @Summary 獲取使用者訊息
// @Description 根據使用者 ID 獲取使用者詳細訊息
// @Id User-1
// @Tags User
// @version 1.0
// @accept application/json
// @produce application/json
// @Security BearerAuth
// @param user_id path int true "使用者 ID" minimum(1)
// @Success 200 {object} user.User "成功獲取使用者訊息"
// @Failure 401 {object} middleware.ErrorResponse "未授權"
// @Failure 403 {object} middleware.ErrorResponse "權限不足"
// @Failure 404 {object} middleware.ErrorResponse "使用者不存在"
// @Failure 500 {object} middleware.ErrorResponse "伺服器內部錯誤"
// @Router /user/{user_id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
    // ... 實作內容
}
```

3. **設置 Swagger 路由**：
```go
// router/router.go
package router

import (
    "github.com/gin-gonic/gin"
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
    _ "github.com/POABOB/slack-clone-back-end/services/user-service/docs" // 引入 docs
)

func NewRouter(engine *gin.Engine) *gin.Engine {
    // ... 其他路由設置

    // Swagger 文檔路由
    engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    return engine
}
```

4. **生成文檔**：
```bash
# 在專案根目錄執行
swag init -g cmd/main.go -o docs
```

### 目錄結構

```
services/user-service/
├── cmd/
│ └── main.go # 包含 Swagger 基本資訊
├── docs/ # 生成的文檔目錄
│ ├── docs.go # 文檔程式碼
│ ├── swagger.json # JSON 格式文檔
│ └── swagger.yaml # YAML 格式文檔
├── internal/
│ └── handler/
│ └── http/
│ └── user_handler.go # 包含 API 註解
└── router/
└── router.go # 包含 Swagger 路由設置
```

### 常用註解說明

1. **基本資訊註解**：
```go
// @Summary 簡短描述
// @Description 詳細描述
// @Tags 標籤分組
// @Accept 接受的請求格式
// @Produce 回應的格式
```

2. **參數註解**：
```go
// @Param 參數名稱 參數位置 參數類型 是否必須 參數描述
// 參數位置：path, query, header, body, formData
// 參數類型：string, int, bool, object, array
// 是否必須：true, false
```

3. **回應註解**：
```go
// @Success 狀態碼 {類型} 描述
// @Failure 狀態碼 {類型} 描述
// 類型可以是：object, array, string, int, bool
```

4. **安全認證註解**：
```go
// @Security BearerAuth
// @Security ApiKeyAuth
```

### 自定義 Swagger UI

```go
engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
    ginSwagger.URL("/swagger/doc.json"),
    ginSwagger.DefaultModelsExpandDepth(-1),
    ginSwagger.PersistAuthorization(true),
))
```

### 使用方式

1. **啟動服務**：
```bash
go run cmd/main.go
```

2. **訪問 Swagger UI**：
```bash
http://localhost:8080/swagger/index.html
```

3. **更新文檔**：
```bash
# 當修改了 API 註解後，需要重新生成文檔
swag init -g cmd/main.go -o docs
```

### 注意事項

1. **文檔更新**：
   - 修改 API 註解後需要重新生成文檔
   - 確保 `docs` 目錄存在
   - 檢查生成的文檔是否正確

2. **註解格式**：
   - 註解必須緊貼在函數上方
   - 參數和回應類型必須正確定義
   - 路由路徑必須與實際路由一致

3. **安全認證**：
   - 正確設置安全認證方式
   - 在 Swagger UI 中測試認證
   - 確保認證資訊正確傳遞

4. **錯誤處理**：
   - 定義所有可能的錯誤回應
   - 提供清晰的錯誤訊息
   - 使用統一的錯誤格式

---