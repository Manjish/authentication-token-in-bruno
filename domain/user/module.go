package user

import (
    "go.uber.org/fx"
)

var Module = fx.Module("user",
	fx.Options(
		fx.Provide(
			NewService,
			NewController,
      		NewRepository,
			NewRoute,
		),
		fx.Invoke(RegisterRoute),
	),
)
