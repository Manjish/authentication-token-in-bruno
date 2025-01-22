package domain

import (
    "bruno_authentication/domain/user"
    "bruno_authentication/domain/middlewares"

    "go.uber.org/fx"
)

var Module = fx.Options(
	middlewares.Module,
	user.Module,
)
