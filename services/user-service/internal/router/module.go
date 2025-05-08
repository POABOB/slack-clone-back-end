package router

import (
	configlib "github.com/POABOB/slack-clone-back-end/pkg/config"
	"go.uber.org/fx"
)

// Module 依賴注入統一管理
var Module = fx.Module("router",
	fx.Provide(
		configlib.NewGinEngine,
		NewRouter,
	),
)
