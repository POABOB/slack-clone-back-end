package internal

import (
	handler "github.com/POABOB/slack-clone-back-end/services/user-service/internal/handler/http"
	repository "github.com/POABOB/slack-clone-back-end/services/user-service/internal/repository/postgresql"
	"github.com/POABOB/slack-clone-back-end/services/user-service/internal/service"
	"go.uber.org/fx"
)

var Module = fx.Module("internal",
	fx.Provide(
		repository.NewUserRepository,
		service.NewUserService,
		handler.NewUserHandler,

		service.NewAuthService,
		handler.NewAuthHandler,
	),
)
