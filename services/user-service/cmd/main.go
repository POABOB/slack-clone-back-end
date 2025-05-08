package main

import (
	"context"
	"fmt"
	configlib "github.com/POABOB/slack-clone-back-end/pkg/config"
	"github.com/POABOB/slack-clone-back-end/services/user-service/config"
	"github.com/POABOB/slack-clone-back-end/services/user-service/internal"
	"github.com/POABOB/slack-clone-back-end/services/user-service/internal/router"
	"github.com/POABOB/slack-clone-back-end/services/user-service/pkg"
	"go.uber.org/fx"
	"log"
)

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
func main() {
	// TODO 依賴注入 JWT AUTH Service
	// TODO 熔斷、超時、重試
	app := fx.New(
		config.Module,
		pkg.AuthModule,
		pkg.PostgresqlModule,
		internal.Module,
		router.Module,
		// 加上 Setup 和 HTTP Server 啟動
		fx.Invoke(
			func(r *router.Router) {
				r.Setup() // 設定路由
			},
			StartHTTPServer, // 啟動 Gin server（用 fx.Lifecycle）
		),
	)
	app.Run()
}

// StartHTTPServer 開啟 HTTP 服務
func StartHTTPServer(lc fx.Lifecycle, r *router.Router, cfg *configlib.ServerConfig) {
	engine := r.Engine() // 取得 gin.Engine
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := engine.Run(addr); err != nil {
					log.Fatalf("failed to run server: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Shutting down...")
			return nil
		},
	})
}
