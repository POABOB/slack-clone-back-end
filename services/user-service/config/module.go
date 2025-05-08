package config

import (
	configlib "github.com/POABOB/slack-clone-back-end/pkg/config"
	"go.uber.org/fx"
)

// Module 依賴注入統一管理
var Module = fx.Module("config",
	fx.Provide(
		func() (*configlib.Config, error) {
			return configlib.LoadConfig()
		},
		func(cfg *configlib.Config) *configlib.ServerConfig { return &cfg.Server },
		func(cfg *configlib.Config) *configlib.RouterConfig { return &cfg.Router },
		func(cfg *configlib.Config) *configlib.JWTConfig { return &cfg.JWT },
		func(cfg *configlib.Config) *configlib.DatabaseConfig { return &cfg.Database },
		func(cfg *configlib.Config) *configlib.RedisConfig { return &cfg.Redis },
	),
)
