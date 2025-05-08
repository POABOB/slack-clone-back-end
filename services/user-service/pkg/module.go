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
