package router

import (
	"github.com/POABOB/slack-clone-back-end/pkg/config"
	"github.com/POABOB/slack-clone-back-end/services/user-service/internal/handler/http"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Router 路由管理器
type Router struct {
	engine *gin.Engine
	config *config.RouterConfig

	// 其他處理器...
	userHandler *handler.UserHandler
	authHandler *handler.AuthHandler
}

// NewRouter 創建新的路由管理器
func NewRouter(engine *gin.Engine, config *config.RouterConfig, userHandler *handler.UserHandler,
	authHandler *handler.AuthHandler) *Router {
	return &Router{
		engine:      engine,
		config:      config,
		userHandler: userHandler,
		authHandler: authHandler,
	}
}

// Setup 設置所有路由
func (r *Router) Setup() {
	// API 版本分組
	v1 := r.engine.Group("/api/v1")
	{
		// 設置各個模組的路由
		r.authHandler.RegisterRoutes(v1)
		r.userHandler.RegisterRoutes(v1)
	}
	url := ginSwagger.URL("/swagger/doc.json") // The url pointing to API definition
	r.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
}

// Engine returns the underlying gin.Engine instance used by the Router.
func (r *Router) Engine() *gin.Engine {
	return r.engine
}
